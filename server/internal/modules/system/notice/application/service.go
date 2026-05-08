// Package application 实现公告的业务逻辑：分页列表、CRUD 和状态切换。
package application

import (
	"context"

	noticedomain "ez-admin-gin/server/internal/modules/system/notice/domain"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Service 封装公告的业务逻辑，包括列表查询、增删改和状态切换。
type Service struct {
	tx   NoticeTransactor
	repo NoticeRepository
}

func NewService(tx NoticeTransactor, repo NoticeRepository) *Service {
	return &Service{tx: tx, repo: repo}
}

// List 按关键词和状态分页查询公告列表。
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

// Create 创建公告并写入数据库。
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

	if err := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		return s.repo.Create(tx, &created)
	}); err != nil {
		return noticedomain.Response{}, err
	}

	return noticedomain.BuildResponse(created), nil
}

// Update 更新指定公告的基本信息。
func (s *Service) Update(noticeID uint, req noticedomain.UpdateRequest) (noticedomain.Response, error) {
	req, err := noticedomain.NormalizeUpdateRequest(req)
	if err != nil {
		return noticedomain.Response{}, err
	}

	var updated noticedomain.Entity
	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
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

// UpdateStatus 切换公告的启用/禁用状态。
func (s *Service) UpdateStatus(noticeID uint, status model.NoticeStatus) error {
	status, err := noticedomain.NormalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, noticeID)
		if err != nil {
			return err
		}
		return s.repo.UpdateStatus(tx, &item, status)
	})
}
