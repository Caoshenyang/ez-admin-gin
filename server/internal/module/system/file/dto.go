package file

import (
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

// ListQuery 表示文件分页查询参数。
type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Ext      string `form:"ext"`
	Status   int    `form:"status"`
}

// Response 表示文件记录返回结构。
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

// ListResponse 表示文件分页结果。
type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// SavedUploadedFile 保存落盘后的文件元信息。
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

// NormalizePage 统一分页边界。
func NormalizePage(page int, pageSize int) (int, int) {
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

// NormalizeExt 统一收敛文件后缀。
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

// NormalizeAllowedExts 对配置中的后缀白名单去重和清洗。
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

// ValidStatus 判断文件状态是否合法。
func ValidStatus(status model.SystemFileStatus) bool {
	return status == model.SystemFileStatusEnabled || status == model.SystemFileStatusDisabled
}

// BuildResponse 把模型对象压成 API 返回结构。
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

// ValidateAllowedExt 守住上传后缀白名单。
func ValidateAllowedExt(ext string, allowedExts []string) error {
	ext = NormalizeExt(ext)
	if ext == "" {
		return apperror.BadRequest("文件后缀不能为空")
	}

	for _, allowed := range allowedExts {
		if ext == NormalizeExt(allowed) {
			return nil
		}
	}

	return apperror.BadRequest("不支持上传该文件类型")
}
