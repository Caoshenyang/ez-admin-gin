package application

import (
	menudomain "ez-admin-gin/server/internal/modules/iam/menu/domain"
	menuinfra "ez-admin-gin/server/internal/modules/iam/menu/infra"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	repo *menuinfra.Repository
}

func NewService(db *gorm.DB, repo *menuinfra.Repository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

func (s *Service) List() ([]menudomain.Response, error) {
	items, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	return buildTree(items), nil
}

func (s *Service) Create(req menudomain.CreateRequest) (menudomain.Response, error) {
	req, err := menudomain.NormalizeCreateRequest(req)
	if err != nil {
		return menudomain.Response{}, err
	}

	created := menudomain.Entity{
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
			return errorsx.BadRequest("菜单编码已存在")
		}
		if err := s.repo.ParentUsable(tx, req.ParentID, req.Type, 0); err != nil {
			return err
		}

		return s.repo.Create(tx, &created)
	})
	if err != nil {
		return menudomain.Response{}, err
	}

	return menudomain.BuildResponse(created), nil
}

func (s *Service) Update(menuID uint, req menudomain.UpdateRequest) (menudomain.Response, error) {
	req, err := menudomain.NormalizeUpdateRequest(req)
	if err != nil {
		return menudomain.Response{}, err
	}

	var updated menudomain.Entity
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
		return menudomain.Response{}, err
	}

	return menudomain.BuildResponse(updated), nil
}

func (s *Service) UpdateStatus(menuID uint, status model.MenuStatus) error {
	if !menudomain.ValidStatus(status) {
		return errorsx.BadRequest("菜单状态不正确")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, menuID)
		if err != nil {
			return err
		}

		return s.repo.UpdateStatus(tx, &item, status)
	})
}

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

func buildTree(items []model.Menu) []menudomain.Response {
	type responseNode struct {
		response menudomain.Response
		children []*responseNode
	}

	nodes := make(map[uint]*responseNode, len(items))
	roots := make([]*responseNode, 0)

	for _, item := range items {
		nodes[item.ID] = &responseNode{response: menudomain.BuildResponse(item)}
	}

	for _, item := range items {
		node := nodes[item.ID]
		if parent, ok := nodes[item.ParentID]; ok {
			parent.children = append(parent.children, node)
			continue
		}
		roots = append(roots, node)
	}

	var toResponses func(nodes []*responseNode) []menudomain.Response
	toResponses = func(nodes []*responseNode) []menudomain.Response {
		result := make([]menudomain.Response, 0, len(nodes))
		for _, node := range nodes {
			item := node.response
			item.Children = toResponses(node.children)
			result = append(result, item)
		}
		return result
	}

	return toResponses(roots)
}
