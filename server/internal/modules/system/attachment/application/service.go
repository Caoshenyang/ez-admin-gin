package application

import (
	"context"
	"errors"
	"mime/multipart"

	attachmentdomain "ez-admin-gin/server/internal/modules/system/attachment/domain"
	attachmentinfra "ez-admin-gin/server/internal/modules/system/attachment/infra"
	fileapp "ez-admin-gin/server/internal/modules/system/file/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Service struct {
	db          *gorm.DB
	repo        *attachmentinfra.Repository
	fileService *fileapp.Service
}

func NewService(db *gorm.DB, repo *attachmentinfra.Repository, fileService *fileapp.Service) *Service {
	return &Service{db: db, repo: repo, fileService: fileService}
}

func (s *Service) List(query attachmentdomain.ListQuery) (attachmentdomain.ListResponse, error) {
	if _, err := attachmentdomain.NormalizeStatusFilter(query.Status); err != nil {
		return attachmentdomain.ListResponse{}, err
	}

	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	items, total, err := s.repo.List(query, page, pageSize)
	if err != nil {
		return attachmentdomain.ListResponse{}, err
	}

	result := make([]attachmentdomain.Response, 0, len(items))
	for _, item := range items {
		result = append(result, attachmentdomain.BuildResponse(item))
	}

	return attachmentdomain.ListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) CreateByUpload(ctx context.Context, uploaderID uint, fileHeader *multipart.FileHeader, req attachmentdomain.CreateRequest) (attachmentdomain.Response, error) {
	uploaded, err := s.fileService.UploadEntity(ctx, uploaderID, fileHeader)
	if err != nil {
		return attachmentdomain.Response{}, err
	}

	req, err = attachmentdomain.NormalizeCreateRequest(req, uploaded.OriginalName)
	if err != nil {
		s.fileService.CleanupUploadedFile(uploaded)
		return attachmentdomain.Response{}, err
	}

	item := model.SystemAttachment{
		FileID:      uploaded.ID,
		DisplayName: req.DisplayName,
		Category:    req.Category,
		BizType:     req.BizType,
		UploaderID:  uploaderID,
		Status:      req.Status,
		Remark:      req.Remark,
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(tx, &item)
	}); err != nil {
		s.fileService.CleanupUploadedFile(uploaded)
		return attachmentdomain.Response{}, errorsx.Internal("创建附件记录失败", err)
	}

	view, err := s.repo.FindViewByID(item.ID)
	if err != nil {
		return attachmentdomain.Response{}, err
	}

	return attachmentdomain.BuildResponse(view), nil
}

func (s *Service) Update(id uint, req attachmentdomain.UpdateRequest) (attachmentdomain.Response, error) {
	req, err := attachmentdomain.NormalizeUpdateRequest(req)
	if err != nil {
		return attachmentdomain.Response{}, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorsx.NotFound("附件不存在")
			}
			return err
		}

		return s.repo.UpdateBase(tx, &item, req)
	})
	if err != nil {
		return attachmentdomain.Response{}, err
	}

	view, err := s.repo.FindViewByID(id)
	if err != nil {
		return attachmentdomain.Response{}, err
	}

	return attachmentdomain.BuildResponse(view), nil
}

func (s *Service) UpdateStatus(id uint, status model.SystemAttachmentStatus) error {
	if !attachmentdomain.ValidStatus(status) {
		return errorsx.BadRequest("附件状态不正确")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorsx.NotFound("附件不存在")
			}
			return err
		}

		return s.repo.UpdateStatus(tx, &item, status)
	})
}
