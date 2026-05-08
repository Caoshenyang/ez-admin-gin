package application

import (
	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	"ez-admin-gin/server/internal/platform/model"
)

type menuNode struct {
	authdomain.MenuResponse
	children []*menuNode
}

// MenuService 提供当前用户可见菜单树的查询服务。
type MenuService struct {
	repo MenuRepository
}

func NewMenuService(repo MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

// Menus 查询当前用户可见的菜单树。
func (s *MenuService) Menus(userID uint) ([]authdomain.MenuResponse, error) {
	menus, err := s.repo.ListMenusByUserID(userID)
	if err != nil {
		return nil, err
	}

	return buildMenuTree(menus), nil
}

func buildMenuTree(menus []model.Menu) []authdomain.MenuResponse {
	nodes := make(map[uint]*menuNode, len(menus))

	for _, menu := range menus {
		nodes[menu.ID] = &menuNode{
			MenuResponse: authdomain.MenuResponse{
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

func menuNodesToResponses(nodes []*menuNode) []authdomain.MenuResponse {
	result := make([]authdomain.MenuResponse, 0, len(nodes))
	for _, node := range nodes {
		item := node.MenuResponse
		item.Children = menuNodesToResponses(node.children)
		result = append(result, item)
	}
	return result
}
