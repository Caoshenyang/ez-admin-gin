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
	Username string `form:"username"`
	IP       string `form:"ip"`
	Status   int    `form:"status"`
}

type Response struct {
	ID        uint                 `json:"id"`
	UserID    uint                 `json:"user_id"`
	Username  string               `json:"username"`
	Status    model.LoginLogStatus `json:"status"`
	Message   string               `json:"message"`
	IP        string               `json:"ip"`
	UserAgent string               `json:"user_agent"`
	CreatedAt time.Time            `json:"created_at"`
}

type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type Entity = model.LoginLog

const PermissionList = "system:loginlog:list"

func NormalizeStatusFilter(value int) (*model.LoginLogStatus, error) {
	if value == 0 {
		return nil, nil
	}
	status := model.LoginLogStatus(value)
	if status != model.LoginLogStatusSuccess && status != model.LoginLogStatusFailed {
		return nil, errorsx.BadRequest("登录状态不正确")
	}
	return &status, nil
}

func NormalizeIP(value string) string {
	return strings.TrimSpace(value)
}

func NormalizeUsername(value string) string {
	return strings.TrimSpace(value)
}

func BuildResponse(item model.LoginLog) Response {
	return Response{
		ID:        item.ID,
		UserID:    item.UserID,
		Username:  item.Username,
		Status:    item.Status,
		Message:   item.Message,
		IP:        item.IP,
		UserAgent: item.UserAgent,
		CreatedAt: item.CreatedAt,
	}
}
