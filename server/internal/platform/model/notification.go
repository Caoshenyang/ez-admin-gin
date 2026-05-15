package model

import (
	"database/sql"
	"time"
)

// NotificationType 表示通知类型。
type NotificationType int

const (
	NotificationTypeSystem   NotificationType = 1 // 系统通知
	NotificationTypeSecurity NotificationType = 2 // 安全通知
	NotificationTypeTask     NotificationType = 3 // 任务通知
	NotificationTypeMessage  NotificationType = 4 // 消息通知
)

// Notification 是通知表模型。
type Notification struct {
	ID        uint64          `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"not null;index:idx_notification_user_unread" json:"user_id"`
	Type      NotificationType `gorm:"type:smallint;not null;default:1" json:"type"`
	Title     string          `gorm:"size:128;not null;default:''" json:"title"`
	Content   string          `gorm:"type:text;not null;default:''" json:"content"`
	Extra     JSONMap         `gorm:"type:jsonb" json:"extra"`
	IsRead    bool            `gorm:"not null;default:false;index:idx_notification_user_unread" json:"is_read"`
	CreatedAt time.Time       `gorm:"not null;default:now()" json:"created_at"`
	ReadAt    sql.NullTime    `json:"read_at,omitempty"`
}

// TableName 固定通知表名。
func (Notification) TableName() string {
	return "sys_notification"
}

// JSONMap 用于映射 JSONB 字段。
type JSONMap map[string]any
