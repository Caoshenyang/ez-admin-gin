package menu

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Service 负责菜单模块的业务规则、树结构组装和事务边界。
type Service struct {
	db   *gorm.DB
	repo *Repository
}

// NewService 创建菜单服务。
func NewService(db *gorm.DB, repo *Repository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

// List 返回完整菜单树。
func (s *Service) List() ([]Response, error) {
	items, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	return buildTree(items), nil
}

// Create 创建目录、菜单或按钮。
func (s *Service) Create(req CreateRequest) (Response, error) {
	req, err := NormalizeCreateRequest(req)
	if err != nil {
		return Response{}, err
	}

	created := Entity{
		ParentID:  req.ParentID,
		Type:      req.Type,
		Code:      req.Code,
		Title:     req.Title,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.repo.CodeExists(tx, req.Code)
		if err != nil {
			return err
		}
		if exists {
			return apperror.BadRequest("菜单编码已存在")
		}
		if err := s.repo.ParentUsable(tx, req.ParentID, req.Type, 0); err != nil {
			return err
		}

		return s.repo.Create(tx, &created)
	})
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(created), nil
}

// Update 编辑菜单基础信息。
func (s *Service) Update(menuID uint, req UpdateRequest) (Response, error) {
	req, err := NormalizeUpdateRequest(req)
	if err != nil {
		return Response{}, err
	}

	var updated Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, menuID)
		if err != nil {
			return err
		}
		if err := s.repo.ParentUsable(tx, req.ParentID, req.Type, menuID); err != nil {
			return err
		}
		if err := s.repo.UpdateBase(tx, &item, req); err != nil {
			return err
		}

		updated = item
		return nil
	})
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(updated), nil
}

// UpdateStatus 单独修改菜单状态。
func (s *Service) UpdateStatus(menuID uint, status model.MenuStatus) error {
	if !ValidStatus(status) {
		return apperror.BadRequest("菜单状态不正确")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, menuID)
		if err != nil {
			return err
		}

		return s.repo.UpdateStatus(tx, &item, status)
	})
}

// Delete 删除菜单节点。
func (s *Service) Delete(menuID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, menuID)
		if err != nil {
			return err
		}
		if err := s.repo.CanDelete(tx, menuID); err != nil {
			return err
		}

		return s.repo.Delete(tx, &item)
	})
}

// buildTree 在服务层统一完成树结构组装，避免 Handler 处理领域结构细节。
func buildTree(items []model.Menu) []Response {
	type responseNode struct {
		response Response
		children []*responseNode
	}

	nodes := make(map[uint]*responseNode, len(items))
	roots := make([]*responseNode, 0)

	for _, item := range items {
		nodes[item.ID] = &responseNode{response: BuildResponse(item)}
	}

	for _, item := range items {
		node := nodes[item.ID]
		if parent, ok := nodes[item.ParentID]; ok {
			parent.children = append(parent.children, node)
			continue
		}
		roots = append(roots, node)
	}

	var toResponses func(nodes []*responseNode) []Response
	toResponses = func(nodes []*responseNode) []Response {
		result := make([]Response, 0, len(nodes))
		for _, node := range nodes {
			item := node.response
			item.Children = toResponses(node.children)
			result = append(result, item)
		}
		return result
	}

	return toResponses(roots)
}
