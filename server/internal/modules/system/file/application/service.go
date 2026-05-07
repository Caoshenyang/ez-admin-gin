package application

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	filedomain "ez-admin-gin/server/internal/modules/system/file/domain"
	fileinfra "ez-admin-gin/server/internal/modules/system/file/infra"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	"ez-admin-gin/server/internal/platform/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	defaultUploadDir        = "uploads"
	defaultUploadPublicPath = "/uploads"
	defaultUploadMaxSizeMB  = 10
	localFileStorage        = "local"
)

var defaultUploadAllowedExts = []string{
	".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".txt", ".docx", ".xlsx",
}

type Service struct {
	db   *gorm.DB
	repo *fileinfra.Repository
	cfg  platformConfig.UploadConfig
	log  *zap.Logger
}

func NewService(db *gorm.DB, repo *fileinfra.Repository, cfg platformConfig.UploadConfig, log *zap.Logger) *Service {
	return &Service{
		db:   db,
		repo: repo,
		cfg:  normalizeUploadConfig(cfg),
		log:  log,
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

	saved, err := s.saveUploadedFile(fileHeader)
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

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(tx, &item)
	}); err != nil {
		_ = os.Remove(saved.AbsolutePath)
		return model.SystemFile{}, errorsx.Internal("保存文件记录失败", err)
	}

	return item, nil
}

func (s *Service) CleanupUploadedFile(item model.SystemFile) {
	_ = s.db.Where("id = ?", item.ID).Delete(&model.SystemFile{}).Error
	_ = os.Remove(item.Path)
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

func (s *Service) saveUploadedFile(fileHeader *multipart.FileHeader) (filedomain.SavedUploadedFile, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return filedomain.SavedUploadedFile{}, err
	}
	defer src.Close()

	now := time.Now()
	dateDir := now.Format("20060102")
	ext := filedomain.NormalizeExt(filepath.Ext(fileHeader.Filename))
	randomPart, err := randomHex(8)
	if err != nil {
		return filedomain.SavedUploadedFile{}, err
	}

	fileName := fmt.Sprintf("%s_%s%s", now.Format("20060102150405"), randomPart, ext)
	targetDir := filepath.Join(s.cfg.Dir, dateDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return filedomain.SavedUploadedFile{}, err
	}

	absolutePath := filepath.Join(targetDir, fileName)
	dst, err := os.OpenFile(absolutePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return filedomain.SavedUploadedFile{}, err
	}
	defer dst.Close()

	hasher := sha256.New()
	written, err := io.Copy(dst, io.TeeReader(src, hasher))
	if err != nil {
		_ = os.Remove(absolutePath)
		return filedomain.SavedUploadedFile{}, err
	}

	publicPath := normalizeUploadPublicPath(s.cfg.PublicPath)
	relativePath := filepath.ToSlash(filepath.Join(s.cfg.Dir, dateDir, fileName))
	url := publicPath + "/" + dateDir + "/" + fileName
	mimeType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return filedomain.SavedUploadedFile{
		OriginalName: filepath.Base(fileHeader.Filename),
		FileName:     fileName,
		Ext:          ext,
		MimeType:     mimeType,
		Size:         written,
		Sha256:       hex.EncodeToString(hasher.Sum(nil)),
		Path:         relativePath,
		URL:          url,
		AbsolutePath: absolutePath,
	}, nil
}

func normalizeUploadConfig(cfg platformConfig.UploadConfig) platformConfig.UploadConfig {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		cfg.Dir = defaultUploadDir
	}

	cfg.PublicPath = normalizeUploadPublicPath(cfg.PublicPath)
	if cfg.MaxSizeMB <= 0 {
		cfg.MaxSizeMB = defaultUploadMaxSizeMB
	}

	if len(cfg.AllowedExts) == 0 {
		cfg.AllowedExts = append([]string(nil), defaultUploadAllowedExts...)
	}
	cfg.AllowedExts = filedomain.NormalizeAllowedExts(cfg.AllowedExts)
	return cfg
}

func normalizeUploadPublicPath(publicPath string) string {
	publicPath = strings.TrimSpace(publicPath)
	if publicPath == "" {
		return defaultUploadPublicPath
	}
	if !strings.HasPrefix(publicPath, "/") {
		publicPath = "/" + publicPath
	}
	return strings.TrimRight(publicPath, "/")
}

func uploadMaxBytes(maxSizeMB int64) int64 {
	if maxSizeMB <= 0 {
		maxSizeMB = defaultUploadMaxSizeMB
	}
	return maxSizeMB * 1024 * 1024
}

func randomHex(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
