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

const superAdminRoleCode = "super_admin"

// RoleHandler 负责后台角色管理接口。
type RoleHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewRoleHandler 创建角色管理 Handler。
func NewRoleHandler(db *gorm.DB, log *zap.Logger) *RoleHandler {
	return &RoleHandler{
		db:  db,
		log: log,
	}
}

type roleListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

type createRoleRequest struct {
	Code   string           `json:"code"`
	Name   string           `json:"name"`
	Sort   int              `json:"sort"`
	Status model.RoleStatus `json:"status"`
	Remark string           `json:"remark"`
}

type updateRoleRequest struct {
	Name   string           `json:"name"`
	Sort   int              `json:"sort"`
	Status model.RoleStatus `json:"status"`
	Remark string           `json:"remark"`
}

type updateRoleStatusRequest struct {
	Status model.RoleStatus `json:"status"`
}

type rolePermissionItem struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type updateRolePermissionsRequest struct {
	Permissions []rolePermissionItem `json:"permissions"`
}

type updateRoleMenusRequest struct {
	MenuIDs []uint `json:"menu_ids"`
}

type roleResponse struct {
	ID          uint                 `json:"id"`
	Code        string               `json:"code"`
	Name        string               `json:"name"`
	Sort        int                  `json:"sort"`
	Status      model.RoleStatus     `json:"status"`
	Remark      string               `json:"remark"`
	Permissions []rolePermissionItem `json:"permissions"`
	MenuIDs     []uint               `json:"menu_ids"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type roleListResponse struct {
	Items    []roleResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// List 返回角色分页列表。
func (h *RoleHandler) List(c *gin.Context) {
	var query roleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeRolePage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.Role{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.RoleStatus(query.Status)
		if !validRoleStatus(status) {
			response.Error(c, apperror.BadRequest("角色状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询角色总数失败", err), h.log)
		return
	}

	var roles []model.Role
	if err := queryDB.
		Order("sort ASC, id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&roles).Error; err != nil {
		response.Error(c, apperror.Internal("查询角色列表失败", err), h.log)
		return
	}

	permissions, err := h.rolePermissions(roles)
	if err != nil {
		response.Error(c, apperror.Internal("查询角色接口权限失败", err), h.log)
		return
	}

	menuIDs, err := h.roleMenuIDs(roles)
	if err != nil {
		response.Error(c, apperror.Internal("查询角色菜单权限失败", err), h.log)
		return
	}

	items := make([]roleResponse, 0, len(roles))
	for _, role := range roles {
		items = append(items, buildRoleResponse(role, permissions[role.Code], menuIDs[role.ID]))
	}

	response.Success(c, roleListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Create 创建角色。
func (h *RoleHandler) Create(c *gin.Context) {
	var req createRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	code, name, status, remark, err := normalizeCreateRoleRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var created model.Role
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureRoleCodeAvailable(tx, code); err != nil {
			return err
		}

		role := model.Role{
			Code:   code,
			Name:   name,
			Sort:   req.Sort,
			Status: status,
			Remark: remark,
		}

		if err := tx.Create(&role).Error; err != nil {
			return err
		}

		created = role
		return nil
	})
	if err != nil {
		writeRoleError(c, err, "创建角色失败", h.log)
		return
	}

	response.Success(c, buildRoleResponse(created, nil, nil))
}

// Update 编辑角色基础信息。
func (h *RoleHandler) Update(c *gin.Context) {
	roleID, ok := roleIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	name, status, remark, err := normalizeUpdateRoleRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var role model.Role
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("角色不存在")
			}
			return err
		}

		if role.Code == superAdminRoleCode && status == model.RoleStatusDisabled {
			return apperror.BadRequest("不能禁用超级管理员角色")
		}

		if err := tx.Model(&role).Updates(map[string]any{
			"name":   name,
			"sort":   req.Sort,
			"status": status,
			"remark": remark,
		}).Error; err != nil {
			return err
		}

		role.Name = name
		role.Sort = req.Sort
		role.Status = status
		role.Remark = remark
		return nil
	})
	if err != nil {
		writeRoleError(c, err, "更新角色失败", h.log)
		return
	}

	response.Success(c, buildRoleResponse(role, nil, nil))
}

// UpdateStatus 修改角色状态。
func (h *RoleHandler) UpdateStatus(c *gin.Context) {
	roleID, ok := roleIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateRoleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if !validRoleStatus(req.Status) {
		response.Error(c, apperror.BadRequest("角色状态不正确"), h.log)
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var role model.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("角色不存在")
			}
			return err
		}

		if role.Code == superAdminRoleCode && req.Status == model.RoleStatusDisabled {
			return apperror.BadRequest("不能禁用超级管理员角色")
		}

		return tx.Model(&role).Update("status", req.Status).Error
	})
	if err != nil {
		writeRoleError(c, err, "更新角色状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     roleID,
		"status": req.Status,
	})
}

// UpdatePermissions 替换角色接口权限。
func (h *RoleHandler) UpdatePermissions(c *gin.Context) {
	roleID, ok := roleIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	permissions, err := normalizeRolePermissions(req.Permissions)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var role model.Role
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("角色不存在")
			}
			return err
		}

		if role.Code == superAdminRoleCode {
			return apperror.BadRequest("超级管理员角色权限不在这里修改")
		}

		return replaceRolePermissions(tx, role.Code, permissions)
	})
	if err != nil {
		writeRoleError(c, err, "更新角色接口权限失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":          roleID,
		"code":        role.Code,
		"permissions": permissions,
	})
}

// UpdateMenus 替换角色菜单权限。
func (h *RoleHandler) UpdateMenus(c *gin.Context) {
	roleID, ok := roleIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateRoleMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	menuIDs, err := normalizeMenuIDs(req.MenuIDs)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		var role model.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("角色不存在")
			}
			return err
		}

		if role.Code == superAdminRoleCode {
			return apperror.BadRequest("超级管理员菜单权限不在这里修改")
		}

		if err := ensureMenusUsable(tx, menuIDs); err != nil {
			return err
		}

		return replaceRoleMenus(tx, roleID, menuIDs)
	})
	if err != nil {
		writeRoleError(c, err, "更新角色菜单权限失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":       roleID,
		"menu_ids": menuIDs,
	})
}

func (h *RoleHandler) rolePermissions(roles []model.Role) (map[string][]rolePermissionItem, error) {
	result := make(map[string][]rolePermissionItem, len(roles))
	if len(roles) == 0 {
		return result, nil
	}

	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	var rows []model.CasbinRule
	if err := h.db.
		Where("ptype = ?", "p").
		Where("v0 IN ?", roleCodes).
		Order("v1 ASC, v2 ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.V0] = append(result[row.V0], rolePermissionItem{
			Path:   row.V1,
			Method: row.V2,
		})
	}

	return result, nil
}

func (h *RoleHandler) roleMenuIDs(roles []model.Role) (map[uint][]uint, error) {
	result := make(map[uint][]uint, len(roles))
	if len(roles) == 0 {
		return result, nil
	}

	roleIDs := make([]uint, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	var rows []model.RoleMenu
	if err := h.db.Where("role_id IN ?", roleIDs).Order("menu_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.RoleID] = append(result[row.RoleID], row.MenuID)
	}

	return result, nil
}

func normalizeCreateRoleRequest(req createRoleRequest) (string, string, model.RoleStatus, string, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return "", "", 0, "", apperror.BadRequest("角色编码不能为空")
	}
	if len(code) > 64 {
		return "", "", 0, "", apperror.BadRequest("角色编码不能超过 64 个字符")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "", "", 0, "", apperror.BadRequest("角色名称不能为空")
	}
	if len(name) > 64 {
		return "", "", 0, "", apperror.BadRequest("角色名称不能超过 64 个字符")
	}

	status := req.Status
	if status == 0 {
		status = model.RoleStatusEnabled
	}
	if !validRoleStatus(status) {
		return "", "", 0, "", apperror.BadRequest("角色状态不正确")
	}

	remark := strings.TrimSpace(req.Remark)
	if len(remark) > 255 {
		return "", "", 0, "", apperror.BadRequest("备注不能超过 255 个字符")
	}

	return code, name, status, remark, nil
}

func normalizeUpdateRoleRequest(req updateRoleRequest) (string, model.RoleStatus, string, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "", 0, "", apperror.BadRequest("角色名称不能为空")
	}
	if len(name) > 64 {
		return "", 0, "", apperror.BadRequest("角色名称不能超过 64 个字符")
	}

	if !validRoleStatus(req.Status) {
		return "", 0, "", apperror.BadRequest("角色状态不正确")
	}

	remark := strings.TrimSpace(req.Remark)
	if len(remark) > 255 {
		return "", 0, "", apperror.BadRequest("备注不能超过 255 个字符")
	}

	return name, req.Status, remark, nil
}

func normalizeRolePermissions(permissions []rolePermissionItem) ([]rolePermissionItem, error) {
	unique := make([]rolePermissionItem, 0, len(permissions))
	seen := make(map[string]struct{}, len(permissions))

	for _, item := range permissions {
		path := strings.TrimSpace(item.Path)
		method := strings.ToUpper(strings.TrimSpace(item.Method))
		if path == "" || method == "" {
			return nil, apperror.BadRequest("接口权限参数不正确")
		}

		key := path + " " + method
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		unique = append(unique, rolePermissionItem{
			Path:   path,
			Method: method,
		})
	}

	return unique, nil
}

func normalizeMenuIDs(menuIDs []uint) ([]uint, error) {
	unique := make([]uint, 0, len(menuIDs))
	seen := make(map[uint]struct{}, len(menuIDs))

	for _, menuID := range menuIDs {
		if menuID == 0 {
			return nil, apperror.BadRequest("菜单 ID 不正确")
		}
		if _, ok := seen[menuID]; ok {
			continue
		}

		seen[menuID] = struct{}{}
		unique = append(unique, menuID)
	}

	return unique, nil
}

func normalizeRolePage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

func validRoleStatus(status model.RoleStatus) bool {
	return status == model.RoleStatusEnabled || status == model.RoleStatusDisabled
}

func roleIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("角色 ID 不正确"), log)
		return 0, false
	}

	return uint(id), true
}

func ensureRoleCodeAvailable(db *gorm.DB, code string) error {
	var role model.Role
	err := db.Unscoped().Where("code = ?", code).First(&role).Error
	if err == nil {
		return apperror.BadRequest("角色编码已存在")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func ensureMenusUsable(db *gorm.DB, menuIDs []uint) error {
	if len(menuIDs) == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.Menu{}).
		Where("id IN ?", menuIDs).
		Where("status = ?", model.MenuStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count != int64(len(menuIDs)) {
		return apperror.BadRequest("菜单不存在或已禁用")
	}

	return nil
}

func replaceRolePermissions(db *gorm.DB, roleCode string, permissions []rolePermissionItem) error {
	if err := db.Where("ptype = ? AND v0 = ?", "p", roleCode).Delete(&model.CasbinRule{}).Error; err != nil {
		return err
	}

	if len(permissions) == 0 {
		return nil
	}

	rows := make([]model.CasbinRule, 0, len(permissions))
	for _, permission := range permissions {
		rows = append(rows, model.CasbinRule{
			Ptype: "p",
			V0:    roleCode,
			V1:    permission.Path,
			V2:    permission.Method,
		})
	}

	return db.Create(&rows).Error
}

func replaceRoleMenus(db *gorm.DB, roleID uint, menuIDs []uint) error {
	if err := db.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		return err
	}

	if len(menuIDs) == 0 {
		return nil
	}

	rows := make([]model.RoleMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		rows = append(rows, model.RoleMenu{
			RoleID: roleID,
			MenuID: menuID,
		})
	}

	return db.Create(&rows).Error
}

func buildRoleResponse(role model.Role, permissions []rolePermissionItem, menuIDs []uint) roleResponse {
	return roleResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Sort:        role.Sort,
		Status:      role.Status,
		Remark:      role.Remark,
		Permissions: permissions,
		MenuIDs:     menuIDs,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func writeRoleError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
