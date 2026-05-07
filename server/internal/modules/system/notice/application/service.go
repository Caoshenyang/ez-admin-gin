package application

import (
	noticedomain "ez-admin-gin/server/internal/modules/system/notice/domain"
	noticeinfra "ez-admin-gin/server/internal/modules/system/notice/infra"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	repo *noticeinfra.Repository
}

func NewService(db *gorm.DB, repo *noticeinfra.Repository) *Service {
	return &Service{db: db, repo: repo}
}

func (s *Service) List(query noticedomain.ListQuery) (noticedomain.ListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := noticedomain.NormalizeStatusFilter(query.Status)
	if err != nil {
		return noticedomain.ListResponse{}, err
	}

	items, total, err := s.repo.List(query, page, pageSize, status)
	if err != nil {
		return noticedomain.ListResponse{}, err
	}

	result := make([]noticedomain.Response, 0, len(items))
	for _, item := range items {
		result = append(result, noticedomain.BuildResponse(item))
	}

	return noticedomain.ListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) Create(req noticedomain.CreateRequest) (noticedomain.Response, error) {
	req, err := noticedomain.NormalizeCreateRequest(req)
	if err != nil {
		return noticedomain.Response{}, err
	}

	created := noticedomain.Entity{
		Title:   req.Title,
		Content: req.Content,
		Sort:    req.Sort,
		Status:  req.Status,
		Remark:  req.Remark,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(tx, &created)
	}); err != nil {
		return noticedomain.Response{}, err
	}

	return noticedomain.BuildResponse(created), nil
}

func (s *Service) Update(noticeID uint, req noticedomain.UpdateRequest) (noticedomain.Response, error) {
	req, err := noticedomain.NormalizeUpdateRequest(req)
	if err != nil {
		return noticedomain.Response{}, err
	}

	var updated noticedomain.Entity
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
		return noticedomain.Response{}, err
	}

	return noticedomain.BuildResponse(updated), nil
}

func (s *Service) UpdateStatus(noticeID uint, status model.NoticeStatus) error {
	status, err := noticedomain.NormalizeStatus(status, false)
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
