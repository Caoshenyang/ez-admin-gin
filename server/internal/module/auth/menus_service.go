package auth

import (
	"ez-admin-gin/server/internal/model"
)

type menuNode struct {
	MenuResponse
	children []*menuNode
}

// MenuService 负责当前登录用户菜单树组装。
type MenuService struct {
	repo *Repository
}

// NewMenuService 创建菜单服务。
func NewMenuService(repo *Repository) *MenuService {
	return &MenuService{repo: repo}
}

// Menus 返回当前登录用户可见菜单树。
func (s *MenuService) Menus(userID uint) ([]MenuResponse, error) {
	menus, err := s.repo.ListMenusByUserID(userID)
	if err != nil {
		return nil, err
	}

	return buildMenuTree(menus), nil
}

func buildMenuTree(menus []model.Menu) []MenuResponse {
	nodes := make(map[uint]*menuNode, len(menus))

	for _, menu := range menus {
		nodes[menu.ID] = &menuNode{
			MenuResponse: MenuResponse{
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

func menuNodesToResponses(nodes []*menuNode) []MenuResponse {
	result := make([]MenuResponse, 0, len(nodes))
	for _, node := range nodes {
		item := node.MenuResponse
		item.Children = menuNodesToResponses(node.children)
		result = append(result, item)
	}

	return result
}
