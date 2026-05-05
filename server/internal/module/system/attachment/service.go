package attachment

import (
	"context"
	"errors"
	"mime/multipart"
	"os"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	systemFileModule "ez-admin-gin/server/internal/module/system/file"

	"gorm.io/gorm"
)

// Service 负责附件中心的业务规则与文件上传复用。
type Service struct {
	db          *gorm.DB
	repo        *Repository
	fileService *systemFileModule.Service
}

// NewService 创建附件中心服务。
func NewService(db *gorm.DB, repo *Repository, fileService *systemFileModule.Service) *Service {
	return &Service{
		db:          db,
		repo:        repo,
		fileService: fileService,
	}
}

// List 返回附件中心分页结果。
func (s *Service) List(query ListQuery) (ListResponse, error) {
	if _, err := NormalizeStatusFilter(query.Status); err != nil {
		return ListResponse{}, err
	}

	page, pageSize := NormalizePage(query.Page, query.PageSize)
	items, total, err := s.repo.List(query, page, pageSize)
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

// CreateByUpload 复用底层文件上传链路创建附件中心记录。
func (s *Service) CreateByUpload(ctx context.Context, uploaderID uint, fileHeader *multipart.FileHeader, req CreateRequest) (Response, error) {
	uploaded, err := s.fileService.UploadEntity(ctx, uploaderID, fileHeader)
	if err != nil {
		return Response{}, err
	}

	req, err = NormalizeCreateRequest(req, uploaded.OriginalName)
	if err != nil {
		s.cleanupUploadedFile(uploaded)
		return Response{}, err
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
		s.cleanupUploadedFile(uploaded)
		return Response{}, apperror.Internal("创建附件记录失败", err)
	}

	view, err := s.repo.FindViewByID(item.ID)
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(view), nil
}

// Update 更新附件元数据。
func (s *Service) Update(id uint, req UpdateRequest) (Response, error) {
	req, err := NormalizeUpdateRequest(req)
	if err != nil {
		return Response{}, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("附件不存在")
			}
			return err
		}

		return s.repo.UpdateBase(tx, &item, req)
	})
	if err != nil {
		return Response{}, err
	}

	view, err := s.repo.FindViewByID(id)
	if err != nil {
		return Response{}, err
	}

	return BuildResponse(view), nil
}

// UpdateStatus 单独修改附件状态。
func (s *Service) UpdateStatus(id uint, status model.SystemAttachmentStatus) error {
	if !ValidStatus(status) {
		return apperror.BadRequest("附件状态不正确")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("附件不存在")
			}
			return err
		}

		return s.repo.UpdateStatus(tx, &item, status)
	})
}

func (s *Service) cleanupUploadedFile(item model.SystemFile) {
	_ = s.db.Where("id = ?", item.ID).Delete(&model.SystemFile{}).Error
	_ = os.Remove(item.Path)
}
