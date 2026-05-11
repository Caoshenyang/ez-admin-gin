package domain

import (
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Ext      string `form:"ext"`
	Status   int    `form:"status"`
}

type Response struct {
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

type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type SavedUploadedFile struct {
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

const (
	PermissionList   = "system:file:list"
	PermissionUpload = "system:file:upload"
)

func NormalizeExt(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))
	if ext == "" {
		return ""
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return ext
}

func NormalizeAllowedExts(exts []string) []string {
	result := make([]string, 0, len(exts))
	seen := make(map[string]struct{}, len(exts))

	for _, ext := range exts {
		ext = NormalizeExt(ext)
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

func ValidStatus(status model.SystemFileStatus) bool {
	return status == model.SystemFileStatusEnabled || status == model.SystemFileStatusDisabled
}

func BuildResponse(item model.SystemFile) Response {
	return Response{
		ID:           item.ID,
		Storage:      item.Storage,
		OriginalName: item.OriginalName,
		FileName:     item.FileName,
		Ext:          item.Ext,
		MimeType:     item.MimeType,
		Size:         item.Size,
		Sha256:       item.Sha256,
		Path:         item.Path,
		URL:          item.URL,
		UploaderID:   item.UploaderID,
		Status:       item.Status,
		Remark:       item.Remark,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}

// extToMIME maps file extensions to their expected MIME type prefixes.
// Used to cross-validate uploaded file content against declared extension.
var extToMIME = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".pdf":  "application/pdf",
	".txt":  "text/plain",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
}

// ExpectedMIME returns the expected MIME type for a given extension, or empty string if unknown.
func ExpectedMIME(ext string) string {
	return extToMIME[NormalizeExt(ext)]
}

func ValidateAllowedExt(ext string, allowedExts []string) error {
	ext = NormalizeExt(ext)
	if ext == "" {
		return errorsx.BadRequest("文件后缀不能为空")
	}

	for _, allowed := range allowedExts {
		if ext == NormalizeExt(allowed) {
			return nil
		}
	}

	return errorsx.BadRequest("不支持上传该文件类型")
}
