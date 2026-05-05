package customer

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

// Service 负责 CRM 客户的业务规则与事务边界。
type Service struct {
	db   *gorm.DB
	repo *Repository
}

// NewService 创建 CRM 客户服务。
func NewService(db *gorm.DB, repo *Repository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

// List 返回当前数据范围内的客户分页结果。
func (s *Service) List(actor datascope.Actor, query ListQuery) (ListResponse, error) {
	page, pageSize := NormalizePage(query.Page, query.PageSize)
	status, err := NormalizeStatusFilter(query.Status)
	if err != nil {
		return ListResponse{}, err
	}

	items, total, err := s.repo.List(actor, query, page, pageSize, status)
	if err != nil {
		return ListResponse{}, err
	}

	result := make([]Response, 0, len(items))
	for _, item := range items {
		result = append(result, BuildResponse(item))
	}

	return ListResponse{
		Items:    result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Create 创建客户。
func (s *Service) Create(actor datascope.Actor, req CreateRequest) (Response, error) {
	req, err := NormalizeCreateRequest(req)
	if err != nil {
		return Response{}, err
	}

	if actor.UserID == 0 {
		return Response{}, apperror.Unauthorized("请先登录")
	}

	var created Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.DepartmentUsable(tx, actor.DepartmentID); err != nil {
			return err
		}
		if err := s.repo.OwnerUsable(tx, actor.UserID); err != nil {
			return err
		}

		created = Entity{
			Name:         req.Name,
			ContactName:  req.ContactName,
			Phone:        req.Phone,
			Level:        req.Level,
			Source:       req.Source,
			DepartmentID: actor.DepartmentID,
			OwnerUserID:  actor.UserID,
			Status:       req.Status,
			Remark:       req.Remark,
		}

		return s.repo.Create(tx, &created)
	})
	if err != nil {
		return Response{}, err
	}

	view, err := s.repo.FindViewByID(nil, created.ID)
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(view), nil
}

// Update 编辑客户。
func (s *Service) Update(actor datascope.Actor, customerID uint, req UpdateRequest) (Response, error) {
	req, err := NormalizeUpdateRequest(req)
	if err != nil {
		return Response{}, err
	}

	var updatedID uint
	err = s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByIDInScope(tx, actor, customerID)
		if err != nil {
			return err
		}
		if err := s.repo.UpdateBase(tx, &item, req); err != nil {
			return err
		}

		updatedID = item.ID
		return nil
	})
	if err != nil {
		return Response{}, err
	}

	view, err := s.repo.FindViewByID(nil, updatedID)
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(view), nil
}

// UpdateStatus 单独修改客户状态。
func (s *Service) UpdateStatus(actor datascope.Actor, customerID uint, status model.CustomerStatus) error {
	status, err := normalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByIDInScope(tx, actor, customerID)
		if err != nil {
			return err
		}
		return s.repo.UpdateStatus(tx, &item, status)
	})
}
