package system

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MenuAdminHandler 负责后台菜单管理接口。
type MenuAdminHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewMenuAdminHandler 创建菜单管理 Handler。
func NewMenuAdminHandler(db *gorm.DB, log *zap.Logger) *MenuAdminHandler {
	return &MenuAdminHandler{
		db:  db,
		log: log,
	}
}

type createMenuRequest struct {
	ParentID  uint             `json:"parent_id"`
	Type      model.MenuType   `json:"type"`
	Code      string           `json:"code"`
	Title     string           `json:"title"`
	Path      string           `json:"path"`
	Component string           `json:"component"`
	Icon      string           `json:"icon"`
	Sort      int              `json:"sort"`
	Status    model.MenuStatus `json:"status"`
	Remark    string           `json:"remark"`
}

type updateMenuRequest struct {
	ParentID  uint             `json:"parent_id"`
	Type      model.MenuType   `json:"type"`
	Title     string           `json:"title"`
	Path      string           `json:"path"`
	Component string           `json:"component"`
	Icon      string           `json:"icon"`
	Sort      int              `json:"sort"`
	Status    model.MenuStatus `json:"status"`
	Remark    string           `json:"remark"`
}

type updateMenuStatusRequest struct {
	Status model.MenuStatus `json:"status"`
}

type menuAdminResponse struct {
	ID        uint                `json:"id"`
	ParentID  uint                `json:"parent_id"`
	Type      model.MenuType      `json:"type"`
	Code      string              `json:"code"`
	Title     string              `json:"title"`
	Path      string              `json:"path"`
	Component string              `json:"component"`
	Icon      string              `json:"icon"`
	Sort      int                 `json:"sort"`
	Status    model.MenuStatus    `json:"status"`
	Remark    string              `json:"remark"`
	Children  []menuAdminResponse `json:"children,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type menuAdminNode struct {
	menuAdminResponse
	children []*menuAdminNode
}

// Tree 返回完整菜单树。
func (h *MenuAdminHandler) Tree(c *gin.Context) {
	var menus []model.Menu
	if err := h.db.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		response.Error(c, apperror.Internal("查询菜单树失败", err), h.log)
		return
	}

	response.Success(c, buildMenuAdminTree(menus))
}

// Create 创建菜单、目录或按钮。
func (h *MenuAdminHandler) Create(c *gin.Context) {
	var req createMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	menu, err := normalizeCreateMenuRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureMenuCodeAvailable(tx, menu.Code); err != nil {
			return err
		}

		if err := ensureParentMenuUsable(tx, menu.ParentID, menu.Type, 0); err != nil {
			return err
		}

		return tx.Create(&menu).Error
	})
	if err != nil {
		writeMenuError(c, err, "创建菜单失败", h.log)
		return
	}

	response.Success(c, buildMenuAdminResponse(menu))
}

// Update 编辑菜单基础信息。
func (h *MenuAdminHandler) Update(c *gin.Context) {
	menuID, ok := menuIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	update, err := normalizeUpdateMenuRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var menu model.Menu
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&menu, menuID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("菜单不存在")
			}
			return err
		}

		if err := ensureParentMenuUsable(tx, update.ParentID, update.Type, menuID); err != nil {
			return err
		}

		if err := tx.Model(&menu).Updates(map[string]any{
			"parent_id": update.ParentID,
			"type":      update.Type,
			"title":     update.Title,
			"path":      update.Path,
			"component": update.Component,
			"icon":      update.Icon,
			"sort":      update.Sort,
			"status":    update.Status,
			"remark":    update.Remark,
		}).Error; err != nil {
			return err
		}

		menu.ParentID = update.ParentID
		menu.Type = update.Type
		menu.Title = update.Title
		menu.Path = update.Path
		menu.Component = update.Component
		menu.Icon = update.Icon
		menu.Sort = update.Sort
		menu.Status = update.Status
		menu.Remark = update.Remark
		return nil
	})
	if err != nil {
		writeMenuError(c, err, "更新菜单失败", h.log)
		return
	}

	response.Success(c, buildMenuAdminResponse(menu))
}

// UpdateStatus 修改菜单状态。
func (h *MenuAdminHandler) UpdateStatus(c *gin.Context) {
	menuID, ok := menuIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateMenuStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if !validMenuStatus(req.Status) {
		response.Error(c, apperror.BadRequest("菜单状态不正确"), h.log)
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var menu model.Menu
		if err := tx.First(&menu, menuID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("菜单不存在")
			}
			return err
		}

		return tx.Model(&menu).Update("status", req.Status).Error
	})
	if err != nil {
		writeMenuError(c, err, "更新菜单状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     menuID,
		"status": req.Status,
	})
}

// Delete 删除菜单。
func (h *MenuAdminHandler) Delete(c *gin.Context) {
	menuID, ok := menuIDParam(c, h.log)
	if !ok {
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var menu model.Menu
		if err := tx.First(&menu, menuID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("菜单不存在")
			}
			return err
		}

		if err := ensureMenuCanDelete(tx, menuID); err != nil {
			return err
		}

		return tx.Delete(&menu).Error
	})
	if err != nil {
		writeMenuError(c, err, "删除菜单失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id": menuID,
	})
}

func normalizeCreateMenuRequest(req createMenuRequest) (model.Menu, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return model.Menu{}, apperror.BadRequest("菜单编码不能为空")
	}
	if len(code) > 128 {
		return model.Menu{}, apperror.BadRequest("菜单编码不能超过 128 个字符")
	}

	title, path, component, icon, status, remark, err := normalizeMenuFields(
		req.Type,
		req.Title,
		req.Path,
		req.Component,
		req.Icon,
		req.Status,
		req.Remark,
	)
	if err != nil {
		return model.Menu{}, err
	}

	return model.Menu{
		ParentID:  req.ParentID,
		Type:      req.Type,
		Code:      code,
		Title:     title,
		Path:      path,
		Component: component,
		Icon:      icon,
		Sort:      req.Sort,
		Status:    status,
		Remark:    remark,
	}, nil
}

func normalizeUpdateMenuRequest(req updateMenuRequest) (model.Menu, error) {
	title, path, component, icon, status, remark, err := normalizeMenuFields(
		req.Type,
		req.Title,
		req.Path,
		req.Component,
		req.Icon,
		req.Status,
		req.Remark,
	)
	if err != nil {
		return model.Menu{}, err
	}

	return model.Menu{
		ParentID:  req.ParentID,
		Type:      req.Type,
		Title:     title,
		Path:      path,
		Component: component,
		Icon:      icon,
		Sort:      req.Sort,
		Status:    status,
		Remark:    remark,
	}, nil
}

func normalizeMenuFields(menuType model.MenuType, title string, path string, component string, icon string, status model.MenuStatus, remark string) (string, string, string, string, model.MenuStatus, string, error) {
	if !validMenuType(menuType) {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单类型不正确")
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单名称不能为空")
	}
	if len(title) > 64 {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单名称不能超过 64 个字符")
	}

	path = strings.TrimSpace(path)
	component = strings.TrimSpace(component)
	icon = strings.TrimSpace(icon)
	remark = strings.TrimSpace(remark)

	if len(path) > 255 {
		return "", "", "", "", 0, "", apperror.BadRequest("路由路径不能超过 255 个字符")
	}
	if len(component) > 255 {
		return "", "", "", "", 0, "", apperror.BadRequest("组件路径不能超过 255 个字符")
	}
	if len(icon) > 64 {
		return "", "", "", "", 0, "", apperror.BadRequest("图标标识不能超过 64 个字符")
	}
	if len(remark) > 255 {
		return "", "", "", "", 0, "", apperror.BadRequest("备注不能超过 255 个字符")
	}

	if status == 0 {
		status = model.MenuStatusEnabled
	}
	if !validMenuStatus(status) {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单状态不正确")
	}

	if menuType == model.MenuTypeMenu && path == "" {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单节点需要填写路由路径")
	}

	return title, path, component, icon, status, remark, nil
}

func ensureMenuCodeAvailable(db *gorm.DB, code string) error {
	var menu model.Menu
	err := db.Unscoped().Where("code = ?", code).First(&menu).Error
	if err == nil {
		return apperror.BadRequest("菜单编码已存在")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func ensureParentMenuUsable(db *gorm.DB, parentID uint, menuType model.MenuType, currentID uint) error {
	if parentID == 0 {
		if menuType != model.MenuTypeDirectory {
			return apperror.BadRequest("根节点只能是目录")
		}
		return nil
	}

	if parentID == currentID {
		return apperror.BadRequest("父级菜单不能选择自己")
	}

	var parent model.Menu
	if err := db.First(&parent, parentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.BadRequest("父级菜单不存在")
		}
		return err
	}

	if parent.Type == model.MenuTypeButton {
		return apperror.BadRequest("按钮下面不能再添加子节点")
	}

	if menuType == model.MenuTypeButton && parent.Type != model.MenuTypeMenu {
		return apperror.BadRequest("按钮只能挂在菜单下面")
	}

	return nil
}

func ensureMenuCanDelete(db *gorm.DB, menuID uint) error {
	var childCount int64
	if err := db.Model(&model.Menu{}).Where("parent_id = ?", menuID).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return apperror.BadRequest("请先删除子菜单")
	}

	var roleMenuCount int64
	if err := db.Model(&model.RoleMenu{}).Where("menu_id = ?", menuID).Count(&roleMenuCount).Error; err != nil {
		return err
	}
	if roleMenuCount > 0 {
		return apperror.BadRequest("菜单已分配给角色，不能删除")
	}

	return nil
}

func validMenuType(menuType model.MenuType) bool {
	return menuType == model.MenuTypeDirectory ||
		menuType == model.MenuTypeMenu ||
		menuType == model.MenuTypeButton
}

func validMenuStatus(status model.MenuStatus) bool {
	return status == model.MenuStatusEnabled || status == model.MenuStatusDisabled
}

func menuIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("菜单 ID 不正确"), log)
		return 0, false
	}

	return uint(id), true
}

func buildMenuAdminTree(menus []model.Menu) []menuAdminResponse {
	nodes := make(map[uint]*menuAdminNode, len(menus))

	for _, menu := range menus {
		nodes[menu.ID] = &menuAdminNode{
			menuAdminResponse: buildMenuAdminResponse(menu),
		}
	}

	roots := make([]*menuAdminNode, 0)
	for _, menu := range menus {
		node := nodes[menu.ID]
		if menu.ParentID == 0 {
			roots = append(roots, node)
			continue
		}

		parent, ok := nodes[menu.ParentID]
		if !ok {
			roots = append(roots, node)
			continue
		}

		parent.children = append(parent.children, node)
	}

	return menuAdminNodesToResponses(roots)
}

func menuAdminNodesToResponses(nodes []*menuAdminNode) []menuAdminResponse {
	result := make([]menuAdminResponse, 0, len(nodes))
	for _, node := range nodes {
		item := node.menuAdminResponse
		item.Children = menuAdminNodesToResponses(node.children)
		result = append(result, item)
	}

	return result
}

func buildMenuAdminResponse(menu model.Menu) menuAdminResponse {
	return menuAdminResponse{
		ID:        menu.ID,
		ParentID:  menu.ParentID,
		Type:      menu.Type,
		Code:      menu.Code,
		Title:     menu.Title,
		Path:      menu.Path,
		Component: menu.Component,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
		Status:    menu.Status,
		Remark:    menu.Remark,
		CreatedAt: menu.CreatedAt,
		UpdatedAt: menu.UpdatedAt,
	}
}

func writeMenuError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
