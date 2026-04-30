package notice

import (
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Service 负责公告的业务规则和事务边界。
type Service struct {
	db   *gorm.DB
	repo *Repository
}

// NewService 创建公告服务。
func NewService(db *gorm.DB, repo *Repository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

// List 返回公告分页结果。
func (s *Service) List(query ListQuery) (ListResponse, error) {
	page, pageSize := NormalizePage(query.Page, query.PageSize)
	status, err := NormalizeStatusFilter(query.Status)
	if err != nil {
		return ListResponse{}, err
	}

	items, total, err := s.repo.List(query, page, pageSize, status)
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

// Create 创建公告。
func (s *Service) Create(req CreateRequest) (Response, error) {
	req, err := NormalizeCreateRequest(req)
	if err != nil {
		return Response{}, err
	}

	created := Entity{
		Title:   req.Title,
		Content: req.Content,
		Sort:    req.Sort,
		Status:  req.Status,
		Remark:  req.Remark,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(tx, &created)
	}); err != nil {
		return Response{}, err
	}

	return BuildResponse(created), nil
}

// Update 编辑公告。
func (s *Service) Update(noticeID uint, req UpdateRequest) (Response, error) {
	req, err := NormalizeUpdateRequest(req)
	if err != nil {
		return Response{}, err
	}

	var updated Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, noticeID)
		if err != nil {
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

// UpdateStatus 单独修改公告状态。
func (s *Service) UpdateStatus(noticeID uint, status model.NoticeStatus) error {
	status, err := normalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, noticeID)
		if err != nil {
			return err
		}
		return s.repo.UpdateStatus(tx, &item, status)
	})
}
