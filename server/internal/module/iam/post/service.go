package post

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	repo *Repository
}

func NewService(db *gorm.DB, repo *Repository) *Service {
	return &Service{db: db, repo: repo}
}

func (s *Service) List(query ListQuery) ([]Response, error) {
	items, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	result := make([]Response, 0, len(items))
	for _, item := range items {
		result = append(result, BuildResponse(item))
	}

	return result, nil
}

func (s *Service) Create(req CreateRequest) (Response, error) {
	code, name, sortValue, status, remark, err := NormalizeInput(req.Code, req.Name, req.Sort, req.Status, req.Remark)
	if err != nil {
		return Response{}, err
	}

	var created Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.repo.CodeExists(tx, code, 0)
		if err != nil {
			return err
		}
		if exists {
			return apperror.BadRequest("岗位编码已存在")
		}

		created = Entity{
			Code:   code,
			Name:   name,
			Sort:   sortValue,
			Status: status,
			Remark: remark,
		}
		return s.repo.Create(tx, &created)
	})
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(created), nil
}

func (s *Service) Update(postID uint, req UpdateRequest) (Response, error) {
	code, name, sortValue, status, remark, err := NormalizeInput(req.Code, req.Name, req.Sort, req.Status, req.Remark)
	if err != nil {
		return Response{}, err
	}

	var updated Entity
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
			return apperror.BadRequest("岗位编码已存在")
		}

		if err := s.repo.Update(tx, &item, code, name, sortValue, status, remark); err != nil {
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

func (s *Service) UpdateStatus(postID uint, status model.PostStatus) error {
	if !ValidStatus(status) {
		return apperror.BadRequest("岗位状态不正确")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, postID)
		if err != nil {
			return err
		}

		return s.repo.UpdateStatus(tx, &item, status)
	})
}
