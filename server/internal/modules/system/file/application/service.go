package application

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	filedomain "ez-admin-gin/server/internal/modules/system/file/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	"ez-admin-gin/server/internal/platform/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	defaultUploadMaxSizeMB = 10
	localFileStorage       = "local"
)

var defaultUploadAllowedExts = []string{
	".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".txt", ".docx", ".xlsx",
}

type Service struct {
	tx      FileTransactor
	repo    FileRepository
	storage FileStorage
	cfg     platformConfig.UploadConfig
	log     *zap.Logger
}

func NewService(tx FileTransactor, repo FileRepository, storage FileStorage, cfg platformConfig.UploadConfig, log *zap.Logger) *Service {
	return &Service{
		tx:      tx,
		repo:    repo,
		storage: storage,
		cfg:     normalizeUploadConfig(cfg),
		log:     log,
	}
}

func (s *Service) List(query filedomain.ListQuery) (filedomain.ListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	items, total, err := s.repo.List(query, page, pageSize)
	if err != nil {
		return filedomain.ListResponse{}, err
	}

	result := make([]filedomain.Response, 0, len(items))
	for _, item := range items {
		result = append(result, filedomain.BuildResponse(item))
	}

	return filedomain.ListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) Upload(ctx context.Context, uploaderID uint, fileHeader *multipart.FileHeader) (filedomain.Response, error) {
	item, err := s.UploadEntity(ctx, uploaderID, fileHeader)
	if err != nil {
		return filedomain.Response{}, err
	}

	return filedomain.BuildResponse(item), nil
}

func (s *Service) UploadEntity(ctx context.Context, uploaderID uint, fileHeader *multipart.FileHeader) (model.SystemFile, error) {
	if err := s.validateUploadFile(fileHeader); err != nil {
		return model.SystemFile{}, err
	}

	saved, err := s.storage.SaveUploadedFile(fileHeader)
	if err != nil {
		return model.SystemFile{}, errorsx.Internal("保存文件失败", err)
	}

	item := model.SystemFile{
		Storage:      localFileStorage,
		OriginalName: saved.OriginalName,
		FileName:     saved.FileName,
		Ext:          saved.Ext,
		MimeType:     saved.MimeType,
		Size:         saved.Size,
		Sha256:       saved.Sha256,
		Path:         saved.Path,
		URL:          saved.URL,
		UploaderID:   uploaderID,
		Status:       model.SystemFileStatusEnabled,
		Remark:       "",
	}

	if err := s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		return s.repo.Create(tx, &item)
	}); err != nil {
		_ = s.storage.Delete(saved.AbsolutePath)
		return model.SystemFile{}, errorsx.Internal("保存文件记录失败", err)
	}

	return item, nil
}

func (s *Service) CleanupUploadedFile(item model.SystemFile) {
	_ = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		return s.repo.DeleteByID(tx, item.ID)
	})
	_ = s.storage.Delete(item.Path)
}

func (s *Service) validateUploadFile(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return errorsx.BadRequest("请选择要上传的文件")
	}
	if fileHeader.Size <= 0 {
		return errorsx.BadRequest("不能上传空文件")
	}

	maxBytes := uploadMaxBytes(s.cfg.MaxSizeMB)
	if fileHeader.Size > maxBytes {
		return errorsx.BadRequest("文件大小不能超过 " + strconv.FormatInt(s.cfg.MaxSizeMB, 10) + " MB")
	}

	ext := filedomain.NormalizeExt(filepath.Ext(fileHeader.Filename))
	return filedomain.ValidateAllowedExt(ext, s.cfg.AllowedExts)
}

func normalizeUploadConfig(cfg platformConfig.UploadConfig) platformConfig.UploadConfig {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		cfg.Dir = "uploads"
	}

	cfg.PublicPath = strings.TrimSpace(cfg.PublicPath)
	if cfg.PublicPath == "" {
		cfg.PublicPath = "/uploads"
	} else {
		if !strings.HasPrefix(cfg.PublicPath, "/") {
			cfg.PublicPath = "/" + cfg.PublicPath
		}
		cfg.PublicPath = strings.TrimRight(cfg.PublicPath, "/")
	}
	if cfg.MaxSizeMB <= 0 {
		cfg.MaxSizeMB = defaultUploadMaxSizeMB
	}

	if len(cfg.AllowedExts) == 0 {
		cfg.AllowedExts = append([]string(nil), defaultUploadAllowedExts...)
	}
	cfg.AllowedExts = filedomain.NormalizeAllowedExts(cfg.AllowedExts)
	return cfg
}

func uploadMaxBytes(maxSizeMB int64) int64 {
	if maxSizeMB <= 0 {
		maxSizeMB = defaultUploadMaxSizeMB
	}
	return maxSizeMB * 1024 * 1024
}
