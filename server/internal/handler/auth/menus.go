package auth

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MenuHandler 负责当前用户菜单相关接口。
type MenuHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewMenuHandler 创建菜单 Handler。
func NewMenuHandler(db *gorm.DB, log *zap.Logger) *MenuHandler {
	return &MenuHandler{
		db:  db,
		log: log,
	}
}

type menuResponse struct {
	ID        uint           `json:"id"`
	ParentID  uint           `json:"parent_id"`
	Type      model.MenuType `json:"type"`
	Code      string         `json:"code"`
	Title     string         `json:"title"`
	Path      string         `json:"path"`
	Component string         `json:"component"`
	Icon      string         `json:"icon"`
	Sort      int            `json:"sort"`
	Children  []menuResponse `json:"children,omitempty"`
}

type menuNode struct {
	menuResponse
	children []*menuNode
}

// Menus 返回当前登录用户可见的菜单树。
func (h *MenuHandler) Menus(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), h.log)
		return
	}

	var menus []model.Menu
	err := h.db.
		Table("sys_menu AS m").
		Select("DISTINCT m.*").
		Joins("JOIN sys_role_menu AS rm ON rm.menu_id = m.id").
		Joins("JOIN sys_user_role AS ur ON ur.role_id = rm.role_id").
		Joins("JOIN sys_role AS r ON r.id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Where("m.status = ?", model.MenuStatusEnabled).
		Where("r.status = ?", model.RoleStatusEnabled).
		Where("m.deleted_at IS NULL").
		Where("r.deleted_at IS NULL").
		Order("m.sort ASC, m.id ASC").
		Find(&menus).Error
	if err != nil {
		response.Error(c, apperror.Internal("查询菜单失败", err), h.log)
		return
	}

	response.Success(c, buildMenuTree(menus))
}

func buildMenuTree(menus []model.Menu) []menuResponse {
	nodes := make(map[uint]*menuNode, len(menus))

	for _, menu := range menus {
		nodes[menu.ID] = &menuNode{
			menuResponse: menuResponse{
				ID:        menu.ID,
				ParentID:  menu.ParentID,
				Type:      menu.Type,
				Code:      menu.Code,
				Title:     menu.Title,
				Path:      menu.Path,
				Component: menu.Component,
				Icon:      menu.Icon,
				Sort:      menu.Sort,
			},
		}
	}

	roots := make([]*menuNode, 0)
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

	return menuNodesToResponses(roots)
}

func menuNodesToResponses(nodes []*menuNode) []menuResponse {
	result := make([]menuResponse, 0, len(nodes))
	for _, node := range nodes {
		item := node.menuResponse
		item.Children = menuNodesToResponses(node.children)
		result = append(result, item)
	}

	return result
}
