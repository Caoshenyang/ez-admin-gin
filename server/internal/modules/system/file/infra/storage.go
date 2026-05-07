package infra

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	filedomain "ez-admin-gin/server/internal/modules/system/file/domain"
	platformConfig "ez-admin-gin/server/internal/platform/config"
)

const (
	defaultUploadDir        = "uploads"
	defaultUploadPublicPath = "/uploads"
)

type LocalStorage struct {
	cfg platformConfig.UploadConfig
}

func NewLocalStorage(cfg platformConfig.UploadConfig) *LocalStorage {
	return &LocalStorage{cfg: normalizeUploadConfig(cfg)}
}

func (s *LocalStorage) SaveUploadedFile(fileHeader *multipart.FileHeader) (filedomain.SavedUploadedFile, error) {
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

	url := normalizeUploadPublicPath(s.cfg.PublicPath) + "/" + dateDir + "/" + fileName
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
		Path:         filepath.ToSlash(filepath.Join(s.cfg.Dir, dateDir, fileName)),
		URL:          url,
		AbsolutePath: absolutePath,
	}, nil
}

func (s *LocalStorage) Delete(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	return os.Remove(path)
}

func normalizeUploadConfig(cfg platformConfig.UploadConfig) platformConfig.UploadConfig {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		cfg.Dir = defaultUploadDir
	}

	cfg.PublicPath = normalizeUploadPublicPath(cfg.PublicPath)
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

func randomHex(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
