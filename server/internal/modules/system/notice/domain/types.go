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
	Status   int    `form:"status"`
}

type CreateRequest struct {
	Title   string             `json:"title"`
	Content string             `json:"content"`
	Sort    int                `json:"sort"`
	Status  model.NoticeStatus `json:"status"`
	Remark  string             `json:"remark"`
}

type UpdateRequest struct {
	Title   string             `json:"title"`
	Content string             `json:"content"`
	Sort    int                `json:"sort"`
	Status  model.NoticeStatus `json:"status"`
	Remark  string             `json:"remark"`
}

type UpdateStatusRequest struct {
	Status model.NoticeStatus `json:"status"`
}

type Response struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	Sort      int                `json:"sort"`
	Status    model.NoticeStatus `json:"status"`
	Remark    string             `json:"remark"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type Entity = model.Notice

const (
	PermissionList         = "system:notice:list"
	PermissionCreate       = "system:notice:create"
	PermissionUpdate       = "system:notice:update"
	PermissionUpdateStatus = "system:notice:update_status"
)

func NormalizeCreateRequest(req CreateRequest) (CreateRequest, error) {
	title, err := normalizeTitle(req.Title)
	if err != nil {
		return CreateRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, true)
	if err != nil {
		return CreateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateRequest{}, err
	}

	req.Title = title
	req.Content = strings.TrimSpace(req.Content)
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	title, err := normalizeTitle(req.Title)
	if err != nil {
		return UpdateRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, false)
	if err != nil {
		return UpdateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return UpdateRequest{}, err
	}

	req.Title = title
	req.Content = strings.TrimSpace(req.Content)
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeStatus(status model.NoticeStatus, allowDefault bool) (model.NoticeStatus, error) {
	if status == 0 && allowDefault {
		status = model.NoticeStatusEnabled
	}
	if status != model.NoticeStatusEnabled && status != model.NoticeStatusDisabled {
		return 0, errorsx.BadRequest("公告状态不正确")
	}
	return status, nil
}

func NormalizeStatusFilter(value int) (*model.NoticeStatus, error) {
	if value == 0 {
		return nil, nil
	}
	status := model.NoticeStatus(value)
	if status != model.NoticeStatusEnabled && status != model.NoticeStatusDisabled {
		return nil, errorsx.BadRequest("公告状态不正确")
	}
	return &status, nil
}

func normalizeTitle(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest("公告标题不能为空")
	}
	if len(value) > 128 {
		return "", errorsx.BadRequest("公告标题不能超过 128 个字符")
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

func BuildResponse(item model.Notice) Response {
	return Response{
		ID:        item.ID,
		Title:     item.Title,
		Content:   item.Content,
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
