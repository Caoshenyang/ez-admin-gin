package application

import (
	postdomain "ez-admin-gin/server/internal/modules/iam/post/domain"
	postinfra "ez-admin-gin/server/internal/modules/iam/post/infra"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	repo *postinfra.Repository
}

func NewService(db *gorm.DB, repo *postinfra.Repository) *Service {
	return &Service{db: db, repo: repo}
}

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

func (s *Service) Create(req postdomain.CreateRequest) (postdomain.Response, error) {
	code, name, sortValue, status, remark, err := postdomain.NormalizeInput(req.Code, req.Name, req.Sort, req.Status, req.Remark)
	if err != nil {
		return postdomain.Response{}, err
	}

	var created postdomain.Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
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

func (s *Service) Update(postID uint, req postdomain.UpdateRequest) (postdomain.Response, error) {
	code, name, sortValue, status, remark, err := postdomain.NormalizeInput(req.Code, req.Name, req.Sort, req.Status, req.Remark)
	if err != nil {
		return postdomain.Response{}, err
	}

	var updated postdomain.Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
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

func (s *Service) UpdateStatus(postID uint, status model.PostStatus) error {
	if !postdomain.ValidStatus(status) {
		return errorsx.BadRequest("岗位状态不正确")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, postID)
		if err != nil {
			return err
		}

		return s.repo.UpdateStatus(tx, &item, status)
	})
}
