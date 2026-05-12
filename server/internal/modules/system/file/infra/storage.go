// Package infra 实现文件的数据访问层。
package infra

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

// LocalStorage 实现本地磁盘文件存储。
type LocalStorage struct {
	cfg platformConfig.UploadConfig
}

func NewLocalStorage(cfg platformConfig.UploadConfig) *LocalStorage {
	return &LocalStorage{cfg: normalizeUploadConfig(cfg)}
}

// SaveUploadedFile 将上传文件按日期目录存储，同时计算 SHA-256 校验。
// 通过 http.DetectContentType 交叉验证文件实际内容与扩展名是否一致。
func (s *LocalStorage) SaveUploadedFile(fileHeader *multipart.FileHeader) (filedomain.SavedUploadedFile, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return filedomain.SavedUploadedFile{}, err
	}
	defer src.Close()

	ext := filedomain.NormalizeExt(filepath.Ext(fileHeader.Filename))
	if err := validateMIMEType(src, ext); err != nil {
		return filedomain.SavedUploadedFile{}, err
	}
	// Reset reader after MIME detection consumed the first 512 bytes.
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return filedomain.SavedUploadedFile{}, err
	}

	now := time.Now()
	dateDir := now.Format("20060102")
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

// validateMIMEType reads the first 512 bytes and cross-validates detected content
// type against the expected MIME type for the given extension.
func validateMIMEType(src io.ReadSeeker, ext string) error {
	expected := filedomain.ExpectedMIME(ext)
	if expected == "" {
		return nil // Unknown extensions skip MIME validation.
	}

	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read file header for MIME detection: %w", err)
	}

	detected := http.DetectContentType(buf[:n])
	if !mimeTypeMatch(detected, expected) {
		return errors.New("文件实际内容与扩展名不匹配，请检查文件是否被篡改")
	}

	return nil
}

// mimeTypeMatch checks if the detected MIME type is compatible with the expected one.
// For binary formats (Office documents), detected may be "application/zip" which is acceptable.
func mimeTypeMatch(detected, expected string) bool {
	if detected == expected {
		return true
	}
	// Office documents (.docx, .xlsx) are ZIP-based; http.DetectContentType reports "application/zip".
	if detected == "application/zip" && strings.HasPrefix(expected, "application/vnd.openxmlformats") {
		return true
	}
	return false
}

// Delete 删除指定路径的物理文件。
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
