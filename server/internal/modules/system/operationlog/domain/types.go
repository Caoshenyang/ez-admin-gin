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
	Method   string `form:"method"`
	Path     string `form:"path"`
	Success  string `form:"success"`
}

type Response struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	RoutePath    string    `json:"route_path"`
	Query        string    `json:"query"`
	IP           string    `json:"ip"`
	UserAgent    string    `json:"user_agent"`
	StatusCode   int       `json:"status_code"`
	LatencyMs    int64     `json:"latency_ms"`
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type Entity = model.OperationLog

const PermissionList = "system:operationlog:list"

func NormalizeSuccessFilter(value string) (*bool, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1":
		result := true
		return &result, nil
	case "false", "0":
		result := false
		return &result, nil
	default:
		return nil, errorsx.BadRequest("成功状态不正确")
	}
}

func BuildResponse(item model.OperationLog) Response {
	return Response{
		ID:           item.ID,
		UserID:       item.UserID,
		Username:     item.Username,
		Method:       item.Method,
		Path:         item.Path,
		RoutePath:    item.RoutePath,
		Query:        item.Query,
		IP:           item.IP,
		UserAgent:    item.UserAgent,
		StatusCode:   item.StatusCode,
		LatencyMs:    item.LatencyMs,
		Success:      item.Success,
		ErrorMessage: item.ErrorMessage,
		CreatedAt:    item.CreatedAt,
	}
}
