package attachment

import (
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

// ListQuery 表示附件中心分页查询参数。
type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Category string `form:"category"`
	BizType  string `form:"biz_type"`
	Ext      string `form:"ext"`
	Status   int    `form:"status"`
}

// CreateRequest 表示上传并创建附件的表单参数。
type CreateRequest struct {
	DisplayName string                       `form:"display_name" json:"display_name"`
	Category    string                       `form:"category" json:"category"`
	BizType     string                       `form:"biz_type" json:"biz_type"`
	Status      model.SystemAttachmentStatus `form:"status" json:"status"`
	Remark      string                       `form:"remark" json:"remark"`
}

// UpdateRequest 表示修改附件元数据请求体。
type UpdateRequest struct {
	DisplayName string                       `json:"display_name"`
	Category    string                       `json:"category"`
	BizType     string                       `json:"biz_type"`
	Status      model.SystemAttachmentStatus `json:"status"`
	Remark      string                       `json:"remark"`
}

// UpdateStatusRequest 表示单独修改附件状态请求体。
type UpdateStatusRequest struct {
	Status model.SystemAttachmentStatus `json:"status"`
}

// Response 表示附件中心返回结构。
type Response struct {
	ID           uint                         `json:"id"`
	FileID       uint                         `json:"file_id"`
	DisplayName  string                       `json:"display_name"`
	Category     string                       `json:"category"`
	BizType      string                       `json:"biz_type"`
	OriginalName string                       `json:"original_name"`
	FileName     string                       `json:"file_name"`
	Ext          string                       `json:"ext"`
	MimeType     string                       `json:"mime_type"`
	Size         int64                        `json:"size"`
	URL          string                       `json:"url"`
	UploaderID   uint                         `json:"uploader_id"`
	Status       model.SystemAttachmentStatus `json:"status"`
	Remark       string                       `json:"remark"`
	CreatedAt    time.Time                    `json:"created_at"`
	UpdatedAt    time.Time                    `json:"updated_at"`
}

// ListResponse 表示附件分页结果。
type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
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

// NormalizeCreateRequest 收口上传附件参数。
func NormalizeCreateRequest(req CreateRequest, fallbackDisplayName string) (CreateRequest, error) {
	displayName, err := normalizeDisplayName(req.DisplayName, fallbackDisplayName)
	if err != nil {
		return CreateRequest{}, err
	}
	status, err := normalizeStatus(req.Status, true)
	if err != nil {
		return CreateRequest{}, err
	}
	category, err := normalizeCategory(req.Category)
	if err != nil {
		return CreateRequest{}, err
	}
	bizType, err := normalizeBizType(req.BizType)
	if err != nil {
		return CreateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateRequest{}, err
	}

	req.DisplayName = displayName
	req.Category = category
	req.BizType = bizType
	req.Status = status
	req.Remark = remark
	return req, nil
}

// NormalizeUpdateRequest 收口附件元数据编辑参数。
func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	displayName, err := normalizeDisplayName(req.DisplayName, "")
	if err != nil {
		return UpdateRequest{}, err
	}
	status, err := normalizeStatus(req.Status, false)
	if err != nil {
		return UpdateRequest{}, err
	}
	category, err := normalizeCategory(req.Category)
	if err != nil {
		return UpdateRequest{}, err
	}
	bizType, err := normalizeBizType(req.BizType)
	if err != nil {
		return UpdateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return UpdateRequest{}, err
	}

	req.DisplayName = displayName
	req.Category = category
	req.BizType = bizType
	req.Status = status
	req.Remark = remark
	return req, nil
}

// NormalizeStatusFilter 把状态查询参数转换成附件状态。
func NormalizeStatusFilter(value int) (*model.SystemAttachmentStatus, error) {
	if value == 0 {
		return nil, nil
	}

	status := model.SystemAttachmentStatus(value)
	if !ValidStatus(status) {
		return nil, apperror.BadRequest("附件状态不正确")
	}

	return &status, nil
}

// ParseAttachmentID 解析路径参数中的附件 ID。
func ParseAttachmentID(value string) (uint, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil || id == 0 {
		return 0, apperror.BadRequest("附件 ID 不正确")
	}
	return uint(id), nil
}

// ValidStatus 判断附件状态是否合法。
func ValidStatus(status model.SystemAttachmentStatus) bool {
	return status == model.SystemAttachmentStatusEnabled || status == model.SystemAttachmentStatusDisabled
}

// BuildResponse 把联合查询视图压成 API 返回结构。
func BuildResponse(item View) Response {
	return Response{
		ID:           item.ID,
		FileID:       item.FileID,
		DisplayName:  item.DisplayName,
		Category:     item.Category,
		BizType:      item.BizType,
		OriginalName: item.OriginalName,
		FileName:     item.FileName,
		Ext:          item.Ext,
		MimeType:     item.MimeType,
		Size:         item.Size,
		URL:          item.URL,
		UploaderID:   item.UploaderID,
		Status:       item.Status,
		Remark:       item.Remark,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}

func normalizeDisplayName(value string, fallback string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		value = strings.TrimSpace(fallback)
	}
	if value == "" {
		return "", apperror.BadRequest("附件名称不能为空")
	}
	if len(value) > 255 {
		return "", apperror.BadRequest("附件名称不能超过 255 个字符")
	}

	return value, nil
}

func normalizeCategory(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 64 {
		return "", apperror.BadRequest("附件分类不能超过 64 个字符")
	}
	return value, nil
}

func normalizeBizType(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 64 {
		return "", apperror.BadRequest("业务类型不能超过 64 个字符")
	}
	return value, nil
}

func normalizeRemark(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 255 {
		return "", apperror.BadRequest("备注不能超过 255 个字符")
	}
	return value, nil
}

func normalizeStatus(status model.SystemAttachmentStatus, allowDefault bool) (model.SystemAttachmentStatus, error) {
	if status == 0 && allowDefault {
		status = model.SystemAttachmentStatusEnabled
	}
	if !ValidStatus(status) {
		return 0, apperror.BadRequest("附件状态不正确")
	}
	return status, nil
}
