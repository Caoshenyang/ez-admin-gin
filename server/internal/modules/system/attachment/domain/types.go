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
	Category string `form:"category"`
	BizType  string `form:"biz_type"`
	Ext      string `form:"ext"`
	Status   int    `form:"status"`
}

type CreateRequest struct {
	DisplayName string                       `form:"display_name" json:"display_name"`
	Category    string                       `form:"category" json:"category"`
	BizType     string                       `form:"biz_type" json:"biz_type"`
	Status      model.SystemAttachmentStatus `form:"status" json:"status"`
	Remark      string                       `form:"remark" json:"remark"`
}

type UpdateRequest struct {
	DisplayName string                       `json:"display_name"`
	Category    string                       `json:"category"`
	BizType     string                       `json:"biz_type"`
	Status      model.SystemAttachmentStatus `json:"status"`
	Remark      string                       `json:"remark"`
}

type UpdateStatusRequest struct {
	Status model.SystemAttachmentStatus `json:"status"`
}

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

type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type Entity = model.SystemAttachment

type View struct {
	ID           uint
	FileID       uint
	DisplayName  string
	Category     string
	BizType      string
	OriginalName string
	FileName     string
	Ext          string
	MimeType     string
	Size         int64
	URL          string
	UploaderID   uint
	Status       model.SystemAttachmentStatus
	Remark       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

const (
	PermissionList         = "system:attachment:list"
	PermissionUpload       = "system:attachment:upload"
	PermissionUpdate       = "system:attachment:update"
	PermissionUpdateStatus = "system:attachment:update_status"
)

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

func NormalizeStatusFilter(value int) (*model.SystemAttachmentStatus, error) {
	if value == 0 {
		return nil, nil
	}

	status := model.SystemAttachmentStatus(value)
	if !ValidStatus(status) {
		return nil, errorsx.BadRequest("附件状态不正确")
	}
	return &status, nil
}

func ValidStatus(status model.SystemAttachmentStatus) bool {
	return status == model.SystemAttachmentStatusEnabled || status == model.SystemAttachmentStatusDisabled
}

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
		return "", errorsx.BadRequest("附件名称不能为空")
	}
	if len(value) > 255 {
		return "", errorsx.BadRequest("附件名称不能超过 255 个字符")
	}

	return value, nil
}

func normalizeCategory(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 64 {
		return "", errorsx.BadRequest("附件分类不能超过 64 个字符")
	}
	return value, nil
}

func normalizeBizType(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 64 {
		return "", errorsx.BadRequest("业务类型不能超过 64 个字符")
	}
	return value, nil
}

func normalizeRemark(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 255 {
		return "", errorsx.BadRequest("备注不能超过 255 个字符")
	}
	return value, nil
}

func normalizeStatus(status model.SystemAttachmentStatus, allowDefault bool) (model.SystemAttachmentStatus, error) {
	if status == 0 && allowDefault {
		status = model.SystemAttachmentStatusEnabled
	}
	if !ValidStatus(status) {
		return 0, errorsx.BadRequest("附件状态不正确")
	}
	return status, nil
}
