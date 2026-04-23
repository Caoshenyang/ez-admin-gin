package system

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

const (
	defaultUploadDir        = "uploads"
	defaultUploadPublicPath = "/uploads"
	defaultUploadMaxSizeMB  = 10
	localFileStorage        = "local"
)

// FileHandler 负责后台文件上传接口。
type FileHandler struct {
	db  *gorm.DB
	cfg config.UploadConfig
	log *zap.Logger
}

// NewFileHandler 创建文件上传 Handler。
func NewFileHandler(db *gorm.DB, cfg config.UploadConfig, log *zap.Logger) *FileHandler {
	return &FileHandler{
		db:  db,
		cfg: normalizeUploadConfig(cfg),
		log: log,
	}
}

type fileListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Ext      string `form:"ext"`
	Status   int    `form:"status"`
}

type fileResponse struct {
	ID           uint                   `json:"id"`
	Storage      string                 `json:"storage"`
	OriginalName string                 `json:"original_name"`
	FileName     string                 `json:"file_name"`
	Ext          string                 `json:"ext"`
	MimeType     string                 `json:"mime_type"`
	Size         int64                  `json:"size"`
	Sha256       string                 `json:"sha256"`
	Path         string                 `json:"path"`
	URL          string                 `json:"url"`
	UploaderID   uint                   `json:"uploader_id"`
	Status       model.SystemFileStatus `json:"status"`
	Remark       string                 `json:"remark"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type fileListResponse struct {
	Items    []fileResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type savedUploadedFile struct {
	OriginalName string
	FileName     string
	Ext          string
	MimeType     string
	Size         int64
	Sha256       string
	Path         string
	URL          string
	AbsolutePath string
}

// List 返回文件分页列表。
func (h *FileHandler) List(c *gin.Context) {
	var query fileListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeFilePage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.SystemFile{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("original_name LIKE ? OR file_name LIKE ?", like, like)
	}

	ext := normalizeFileExt(query.Ext)
	if ext != "" {
		queryDB = queryDB.Where("ext = ?", ext)
	}

	if query.Status != 0 {
		status := model.SystemFileStatus(query.Status)
		if !validFileStatus(status) {
			response.Error(c, apperror.BadRequest("文件状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询文件总数失败", err), h.log)
		return
	}

	var files []model.SystemFile
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&files).Error; err != nil {
		response.Error(c, apperror.Internal("查询文件列表失败", err), h.log)
		return
	}

	items := make([]fileResponse, 0, len(files))
	for _, file := range files {
		items = append(items, buildFileResponse(file))
	}

	response.Success(c, fileListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Upload 上传文件到本地目录，并写入文件记录。
func (h *FileHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, apperror.BadRequest("请选择要上传的文件"), h.log)
		return
	}

	if err := h.validateUploadFile(fileHeader); err != nil {
		response.Error(c, err, h.log)
		return
	}

	saved, err := h.saveUploadedFile(fileHeader)
	if err != nil {
		response.Error(c, apperror.Internal("保存文件失败", err), h.log)
		return
	}

	uploaderID, _ := middleware.CurrentUserID(c)
	file := model.SystemFile{
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

	if err := h.db.Create(&file).Error; err != nil {
		_ = os.Remove(saved.AbsolutePath)
		response.Error(c, apperror.Internal("保存文件记录失败", err), h.log)
		return
	}

	response.Success(c, buildFileResponse(file))
}

func (h *FileHandler) validateUploadFile(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size <= 0 {
		return apperror.BadRequest("不能上传空文件")
	}

	maxBytes := uploadMaxBytes(h.cfg.MaxSizeMB)
	if fileHeader.Size > maxBytes {
		return apperror.BadRequest("文件大小不能超过 " + strconv.FormatInt(h.cfg.MaxSizeMB, 10) + " MB")
	}

	ext := normalizeFileExt(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		return apperror.BadRequest("文件后缀不能为空")
	}

	if !isAllowedFileExt(ext, h.cfg.AllowedExts) {
		return apperror.BadRequest("不支持上传该文件类型")
	}

	return nil
}

func (h *FileHandler) saveUploadedFile(fileHeader *multipart.FileHeader) (savedUploadedFile, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return savedUploadedFile{}, err
	}
	defer src.Close()

	now := time.Now()
	dateDir := now.Format("20060102")
	ext := normalizeFileExt(filepath.Ext(fileHeader.Filename))
	randomPart, err := randomHex(8)
	if err != nil {
		return savedUploadedFile{}, err
	}

	// 保存文件名由后端生成，避免重名和路径拼接风险。
	fileName := fmt.Sprintf("%s_%s%s", now.Format("20060102150405"), randomPart, ext)
	targetDir := filepath.Join(h.cfg.Dir, dateDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return savedUploadedFile{}, err
	}

	absolutePath := filepath.Join(targetDir, fileName)
	dst, err := os.OpenFile(absolutePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return savedUploadedFile{}, err
	}
	defer dst.Close()

	hasher := sha256.New()
	written, err := io.Copy(dst, io.TeeReader(src, hasher))
	if err != nil {
		_ = os.Remove(absolutePath)
		return savedUploadedFile{}, err
	}

	publicPath := normalizeUploadPublicPath(h.cfg.PublicPath)
	relativePath := filepath.ToSlash(filepath.Join(h.cfg.Dir, dateDir, fileName))
	url := publicPath + "/" + dateDir + "/" + fileName
	mimeType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return savedUploadedFile{
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

	cfg.AllowedExts = normalizeAllowedExts(cfg.AllowedExts)
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

func normalizeAllowedExts(exts []string) []string {
	result := make([]string, 0, len(exts))
	seen := make(map[string]struct{}, len(exts))

	for _, ext := range exts {
		ext = normalizeFileExt(ext)
		if ext == "" {
			continue
		}
		if _, ok := seen[ext]; ok {
			continue
		}

		seen[ext] = struct{}{}
		result = append(result, ext)
	}

	return result
}

func normalizeFileExt(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))
	if ext == "" {
		return ""
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return ext
}

func isAllowedFileExt(ext string, allowedExts []string) bool {
	ext = normalizeFileExt(ext)
	if ext == "" || len(allowedExts) == 0 {
		return false
	}

	for _, allowed := range allowedExts {
		if ext == normalizeFileExt(allowed) {
			return true
		}
	}

	return false
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

func normalizeFilePage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

func validFileStatus(status model.SystemFileStatus) bool {
	return status == model.SystemFileStatusEnabled || status == model.SystemFileStatusDisabled
}

func buildFileResponse(file model.SystemFile) fileResponse {
	return fileResponse{
		ID:           file.ID,
		Storage:      file.Storage,
		OriginalName: file.OriginalName,
		FileName:     file.FileName,
		Ext:          file.Ext,
		MimeType:     file.MimeType,
		Size:         file.Size,
		Sha256:       file.Sha256,
		Path:         file.Path,
		URL:          file.URL,
		UploaderID:   file.UploaderID,
		Status:       file.Status,
		Remark:       file.Remark,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
	}
}

func writeFileError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
