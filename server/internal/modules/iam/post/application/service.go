// Package application 实现岗位的业务逻辑：列表查询、CRUD 和状态切换。
package application

import (
	"context"

	postdomain "ez-admin-gin/server/internal/modules/iam/post/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Service 提供岗位的业务操作服务。
type Service struct {
	tx   PostTransactor
	repo PostRepository
}

func NewService(tx PostTransactor, repo PostRepository) *Service {
	return &Service{tx: tx, repo: repo}
}

// List 根据查询条件返回岗位列表。
func (s *Service) List(query postdomain.ListQuery) ([]postdomain.Response, error) {
	items, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	result := make([]postdomain.Response, 0, len(items))
	for _, item := range items {
		result = append(result, postdomain.BuildResponse(item))
	}

	return result, nil
}

// Create 校验编码唯一性后创建岗位。
func (s *Service) Create(req postdomain.CreateRequest) (postdomain.Response, error) {
	code, name, sortValue, status, remark, err := postdomain.NormalizeInput(req.Code, req.Name, req.Sort, req.Status, req.Remark)
	if err != nil {
		return postdomain.Response{}, err
	}

	var created postdomain.Entity
	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		exists, err := s.repo.CodeExists(tx, code, 0)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("岗位编码已存在")
		}

		created = postdomain.Entity{
			Code:   code,
			Name:   name,
			Sort:   sortValue,
			Status: status,
			Remark: remark,
		}
		return s.repo.Create(tx, &created)
	})
	if err != nil {
		return postdomain.Response{}, err
	}

	return postdomain.BuildResponse(created), nil
}

// Update 更新岗位信息。
func (s *Service) Update(postID uint, req postdomain.UpdateRequest) (postdomain.Response, error) {
	code, name, sortValue, status, remark, err := postdomain.NormalizeInput(req.Code, req.Name, req.Sort, req.Status, req.Remark)
	if err != nil {
		return postdomain.Response{}, err
	}

	var updated postdomain.Entity
	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, postID)
		if err != nil {
			return err
		}

		exists, err := s.repo.CodeExists(tx, code, postID)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("岗位编码已存在")
		}

		if err := s.repo.Update(tx, &item, code, name, sortValue, status, remark); err != nil {
			return err
		}

		updated = item
		return nil
	})
	if err != nil {
		return postdomain.Response{}, err
	}

	return postdomain.BuildResponse(updated), nil
}

// UpdateStatus 切换岗位的启用/禁用状态。
func (s *Service) UpdateStatus(postID uint, status model.PostStatus) error {
	if !postdomain.ValidStatus(status) {
		return errorsx.BadRequest("岗位状态不正确")
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, postID)
		if err != nil {
			return err
		}

		return s.repo.UpdateStatus(tx, &item, status)
	})
}
