package file

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

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	defaultUploadDir        = "uploads"
	defaultUploadPublicPath = "/uploads"
	defaultUploadMaxSizeMB  = 10
	localFileStorage        = "local"
)

// Service 负责文件模块的业务规则、落盘流程和事务边界。
type Service struct {
	db   *gorm.DB
	repo *Repository
	cfg  config.UploadConfig
	log  *zap.Logger
}

// NewService 创建文件服务。
func NewService(db *gorm.DB, repo *Repository, cfg config.UploadConfig, log *zap.Logger) *Service {
	return &Service{
		db:   db,
		repo: repo,
		cfg:  normalizeUploadConfig(cfg),
		log:  log,
	}
}

// List 返回文件分页结果。
func (s *Service) List(query ListQuery) (ListResponse, error) {
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

// Upload 上传文件到本地目录，并写入文件记录。
func (s *Service) Upload(ctx context.Context, uploaderID uint, fileHeader *multipart.FileHeader) (Response, error) {
	if err := s.validateUploadFile(fileHeader); err != nil {
		return Response{}, err
	}

	saved, err := s.saveUploadedFile(fileHeader)
	if err != nil {
		return Response{}, apperror.Internal("保存文件失败", err)
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
		return Response{}, apperror.Internal("保存文件记录失败", err)
	}

	return BuildResponse(item), nil
}

func (s *Service) validateUploadFile(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return apperror.BadRequest("请选择要上传的文件")
	}
	if fileHeader.Size <= 0 {
		return apperror.BadRequest("不能上传空文件")
	}

	maxBytes := uploadMaxBytes(s.cfg.MaxSizeMB)
	if fileHeader.Size > maxBytes {
		return apperror.BadRequest("文件大小不能超过 " + strconv.FormatInt(s.cfg.MaxSizeMB, 10) + " MB")
	}

	ext := NormalizeExt(filepath.Ext(fileHeader.Filename))
	return ValidateAllowedExt(ext, s.cfg.AllowedExts)
}

func (s *Service) saveUploadedFile(fileHeader *multipart.FileHeader) (SavedUploadedFile, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return SavedUploadedFile{}, err
	}
	defer src.Close()

	now := time.Now()
	dateDir := now.Format("20060102")
	ext := NormalizeExt(filepath.Ext(fileHeader.Filename))
	randomPart, err := randomHex(8)
	if err != nil {
		return SavedUploadedFile{}, err
	}

	// 保存文件名由后端生成，避免重名和路径拼接风险。
	fileName := fmt.Sprintf("%s_%s%s", now.Format("20060102150405"), randomPart, ext)
	targetDir := filepath.Join(s.cfg.Dir, dateDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return SavedUploadedFile{}, err
	}

	absolutePath := filepath.Join(targetDir, fileName)
	dst, err := os.OpenFile(absolutePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return SavedUploadedFile{}, err
	}
	defer dst.Close()

	hasher := sha256.New()
	written, err := io.Copy(dst, io.TeeReader(src, hasher))
	if err != nil {
		_ = os.Remove(absolutePath)
		return SavedUploadedFile{}, err
	}

	publicPath := normalizeUploadPublicPath(s.cfg.PublicPath)
	relativePath := filepath.ToSlash(filepath.Join(s.cfg.Dir, dateDir, fileName))
	url := publicPath + "/" + dateDir + "/" + fileName
	mimeType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return SavedUploadedFile{
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

func normalizeUploadConfig(cfg config.UploadConfig) config.UploadConfig {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		cfg.Dir = defaultUploadDir
	}

	cfg.PublicPath = normalizeUploadPublicPath(cfg.PublicPath)
	if cfg.MaxSizeMB <= 0 {
		cfg.MaxSizeMB = defaultUploadMaxSizeMB
	}

	cfg.AllowedExts = NormalizeAllowedExts(cfg.AllowedExts)
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
