package domain

import (
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

type ListQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
	Type     int `form:"type"`
	IsRead   int `form:"is_read"` // 0=all, 1=unread, 2=read
}

type MarkReadRequest struct {
	IDs []uint64 `json:"ids"`
}

type Response struct {
	ID        uint64                   `json:"id"`
	Type      model.NotificationType   `json:"type"`
	Title     string                   `json:"title"`
	Content   string                   `json:"content"`
	Extra     model.JSONMap            `json:"extra"`
	IsRead    bool                     `json:"is_read"`
	CreatedAt time.Time                `json:"created_at"`
	ReadAt    *time.Time               `json:"read_at"`
}

type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

type Entity = model.Notification

const (
	PermissionList       = "system:notification:list"
	PermissionMarkRead   = "system:notification:mark_read"
)

func NormalizeTypeFilter(value int) (*model.NotificationType, error) {
	if value == 0 {
		return nil, nil
	}
	t := model.NotificationType(value)
	if t < model.NotificationTypeSystem || t > model.NotificationTypeMessage {
		return nil, errorsx.BadRequest("通知类型不正确")
	}
	return &t, nil
}

func NormalizeReadFilter(value int) (*bool, error) {
	if value == 0 {
		return nil, nil
	}
	isRead := value == 2
	return &isRead, nil
}

func NormalizeCreateRequest(title, content string, notificationType model.NotificationType) (string, string, model.NotificationType, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", "", 0, errorsx.BadRequest("通知标题不能为空")
	}
	if len(title) > 128 {
		return "", "", 0, errorsx.BadRequest("通知标题不能超过 128 个字符")
	}
	if notificationType < model.NotificationTypeSystem || notificationType > model.NotificationTypeMessage {
		return "", "", 0, errorsx.BadRequest("通知类型不正确")
	}
	return title, strings.TrimSpace(content), notificationType, nil
}

func BuildResponse(item model.Notification) Response {
	resp := Response{
		ID:        item.ID,
		Type:      item.Type,
		Title:     item.Title,
		Content:   item.Content,
		Extra:     item.Extra,
		IsRead:    item.IsRead,
		CreatedAt: item.CreatedAt,
	}
	if item.ReadAt.Valid {
		resp.ReadAt = &item.ReadAt.Time
	}
	return resp
}
