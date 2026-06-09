// Package application 实现附件的业务逻辑：按业务关联查询附件列表。
package application

import (
	"context"
	"errors"
	"mime/multipart"

	attachmentdomain "ez-admin-gin/server/internal/modules/system/attachment/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Service 封装附件的业务逻辑，包括列表查询、上传创建、更新和状态切换。
type Service struct {
	tx          AttachmentTransactor
	repo        AttachmentRepository
	fileService FileAssetService
}

func NewService(tx AttachmentTransactor, repo AttachmentRepository, fileService FileAssetService) *Service {
	return &Service{tx: tx, repo: repo, fileService: fileService}
}

// List 按关键词、分类、业务类型、扩展名和状态分页查询附件列表。
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

// CreateByUpload 上传文件并创建附件记录，失败时自动回滚已上传的文件。
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

	if err := s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
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

// Update 更新指定附件的基本信息。
func (s *Service) Update(id uint, req attachmentdomain.UpdateRequest) (attachmentdomain.Response, error) {
	req, err := attachmentdomain.NormalizeUpdateRequest(req)
	if err != nil {
		return attachmentdomain.Response{}, err
	}

	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
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

// UpdateStatus 切换附件的启用/禁用状态。
func (s *Service) UpdateStatus(id uint, status model.SystemAttachmentStatus) error {
	if !attachmentdomain.ValidStatus(status) {
		return errorsx.BadRequest("附件状态不正确")
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
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
