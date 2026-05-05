package followup

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

// Service 负责 CRM 客户跟进的业务规则与事务边界。
type Service struct {
	db   *gorm.DB
	repo *Repository
}

// NewService 创建 CRM 客户跟进服务。
func NewService(db *gorm.DB, repo *Repository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

// List 返回当前数据范围内的客户跟进分页结果。
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

// ListCustomerOptions 返回当前数据范围内可选客户。
func (s *Service) ListCustomerOptions(actor datascope.Actor, keyword string, limit int) ([]CustomerOption, error) {
	return s.repo.ListCustomerOptions(actor, keyword, NormalizeCustomerOptionLimit(limit))
}

// Create 创建客户跟进。
func (s *Service) Create(actor datascope.Actor, req CreateRequest) (Response, error) {
	req, err := NormalizeCreateRequest(req)
	if err != nil {
		return Response{}, err
	}

	var created Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		customer, err := s.repo.FindCustomerInScope(tx, actor, req.CustomerID)
		if err != nil {
			return err
		}
		if customer.Status != model.CustomerStatusEnabled {
			return apperror.BadRequest("当前客户已停用，不能继续创建跟进")
		}

		created = Entity{
			CustomerID:   customer.ID,
			DepartmentID: customer.DepartmentID,
			OwnerUserID:  customer.OwnerUserID,
			FollowType:   req.FollowType,
			Subject:      req.Subject,
			Content:      req.Content,
			Result:       req.Result,
			NextFollowAt: req.NextFollowAt,
			Status:       req.Status,
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

// Update 编辑客户跟进。
func (s *Service) Update(actor datascope.Actor, followUpID uint, req UpdateRequest) (Response, error) {
	req, err := NormalizeUpdateRequest(req)
	if err != nil {
		return Response{}, err
	}

	var updatedID uint
	err = s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByIDInScope(tx, actor, followUpID)
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

// UpdateStatus 单独修改客户跟进状态。
func (s *Service) UpdateStatus(actor datascope.Actor, followUpID uint, status model.CustomerFollowUpStatus) error {
	status, err := normalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByIDInScope(tx, actor, followUpID)
		if err != nil {
			return err
		}
		return s.repo.UpdateStatus(tx, &item, status)
	})
}
