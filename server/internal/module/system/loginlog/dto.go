package loginlog

import (
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

// ListQuery 表示登录日志分页查询参数。
type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `form:"username"`
	IP       string `form:"ip"`
	Status   int    `form:"status"`
}

// Response 表示登录日志对象返回结构。
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

// ListResponse 表示登录日志分页结果。
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

// NormalizeStatusFilter 把状态查询参数转换成登录日志状态。
func NormalizeStatusFilter(value int) (*model.LoginLogStatus, error) {
	if value == 0 {
		return nil, nil
	}

	status := model.LoginLogStatus(value)
	if status != model.LoginLogStatusSuccess && status != model.LoginLogStatusFailed {
		return nil, apperror.BadRequest("登录状态不正确")
	}

	return &status, nil
}

// NormalizeIP 统一收口 IP 查询参数。
func NormalizeIP(value string) string {
	return strings.TrimSpace(value)
}

// NormalizeUsername 统一收口用户名查询参数。
func NormalizeUsername(value string) string {
	return strings.TrimSpace(value)
}

// BuildResponse 把模型对象压成 API 返回结构。
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
