package domain

import (
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

type ListQuery struct {
	Keyword string `form:"keyword"`
	Status  int    `form:"status"`
}

type CreateRequest struct {
	Code   string           `json:"code"`
	Name   string           `json:"name"`
	Sort   int              `json:"sort"`
	Status model.PostStatus `json:"status"`
	Remark string           `json:"remark"`
}

type UpdateRequest struct {
	Code   string           `json:"code"`
	Name   string           `json:"name"`
	Sort   int              `json:"sort"`
	Status model.PostStatus `json:"status"`
	Remark string           `json:"remark"`
}

type UpdateStatusRequest struct {
	Status model.PostStatus `json:"status"`
}

type Response struct {
	ID        uint             `json:"id"`
	Code      string           `json:"code"`
	Name      string           `json:"name"`
	Sort      int              `json:"sort"`
	Status    model.PostStatus `json:"status"`
	Remark    string           `json:"remark"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type Entity = model.Post

const (
	PermissionList         = "system:post:list"
	PermissionCreate       = "system:post:create"
	PermissionUpdate       = "system:post:update"
	PermissionUpdateStatus = "system:post:status"
)

func NormalizeInput(code string, name string, sort int, status model.PostStatus, remark string) (string, string, int, model.PostStatus, string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", "", 0, 0, "", errorsx.BadRequest("岗位编码不能为空")
	}
	if len(code) > 64 {
		return "", "", 0, 0, "", errorsx.BadRequest("岗位编码不能超过 64 个字符")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return "", "", 0, 0, "", errorsx.BadRequest("岗位名称不能为空")
	}
	if len(name) > 64 {
		return "", "", 0, 0, "", errorsx.BadRequest("岗位名称不能超过 64 个字符")
	}

	if status == 0 {
		status = model.PostStatusEnabled
	}
	if !ValidStatus(status) {
		return "", "", 0, 0, "", errorsx.BadRequest("岗位状态不正确")
	}

	remark = strings.TrimSpace(remark)
	if len(remark) > 255 {
		return "", "", 0, 0, "", errorsx.BadRequest("备注不能超过 255 个字符")
	}

	return code, name, sort, status, remark, nil
}

func ValidStatus(status model.PostStatus) bool {
	return status == model.PostStatusEnabled || status == model.PostStatusDisabled
}

func BuildResponse(item model.Post) Response {
	return Response{
		ID:        item.ID,
		Code:      item.Code,
		Name:      item.Name,
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
