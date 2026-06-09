package model

import (
	"time"

	"gorm.io/gorm"
)

// MessageConfigStatus 表示消息配置状态。
type MessageConfigStatus int

const (
	// MessageConfigStatusEnabled 表示配置已启用。
	MessageConfigStatusEnabled MessageConfigStatus = 1
	// MessageConfigStatusDisabled 表示配置已禁用。
	MessageConfigStatusDisabled MessageConfigStatus = 2
)

// MessageTemplateType 表示消息模板类型。
type MessageTemplateType int

const (
	// MessageTemplateTypeNotification 表示站内通知模板。
	MessageTemplateTypeNotification MessageTemplateType = 1
	// MessageTemplateTypeTodo 表示待办提醒模板。
	MessageTemplateTypeTodo MessageTemplateType = 2
	// MessageTemplateTypeAlert 表示告警提醒模板。
	MessageTemplateTypeAlert MessageTemplateType = 3
)

// MessageReceiverType 表示提醒接收人配置方式。
type MessageReceiverType int

const (
	// MessageReceiverTypeRole 表示按角色接收。
	MessageReceiverTypeRole MessageReceiverType = 1
	// MessageReceiverTypeUser 表示按指定用户接收。
	MessageReceiverTypeUser MessageReceiverType = 2
	// MessageReceiverTypeDepartment 表示按部门接收。
	MessageReceiverTypeDepartment MessageReceiverType = 3
	// MessageReceiverTypeInitiator 表示业务发起人接收。
	MessageReceiverTypeInitiator MessageReceiverType = 4
	// MessageReceiverTypeAssignee 表示业务负责人接收。
	MessageReceiverTypeAssignee MessageReceiverType = 5
)

// MessageTemplate 是消息模板表模型。
type MessageTemplate struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	Code      string              `gorm:"size:128;not null;uniqueIndex" json:"code"`
	Name      string              `gorm:"size:64;not null" json:"name"`
	Title     string              `gorm:"size:128;not null" json:"title"`
	Content   string              `gorm:"type:text;not null" json:"content"`
	Type      MessageTemplateType `gorm:"type:smallint;not null;default:1" json:"type"`
	Variables string              `gorm:"type:text;not null" json:"variables"`
	Sort      int                 `gorm:"not null;default:0" json:"sort"`
	Status    MessageConfigStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	IsSystem  bool                `gorm:"not null;default:false" json:"is_system"`
	Remark    string              `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"-"`
}

// TableName 固定消息模板表名。
func (MessageTemplate) TableName() string {
	return "sys_message_template"
}

// MessageReminder 是消息提醒规则表模型。
type MessageReminder struct {
	ID             uint                `gorm:"primaryKey" json:"id"`
	Code           string              `gorm:"size:128;not null;uniqueIndex" json:"code"`
	Name           string              `gorm:"size:64;not null" json:"name"`
	TriggerEvent   string              `gorm:"size:128;not null;index" json:"trigger_event"`
	TemplateID     uint                `gorm:"not null;index" json:"template_id"`
	Channels       string              `gorm:"size:128;not null;default:'notification'" json:"channels"`
	ReceiverType   MessageReceiverType `gorm:"type:smallint;not null;default:1" json:"receiver_type"`
	ReceiverValues string              `gorm:"type:text;not null" json:"receiver_values"`
	AdvanceMinutes int                 `gorm:"not null;default:0" json:"advance_minutes"`
	LinkURL        string              `gorm:"size:255;not null;default:''" json:"link_url"`
	Sort           int                 `gorm:"not null;default:0" json:"sort"`
	Status         MessageConfigStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	IsSystem       bool                `gorm:"not null;default:false" json:"is_system"`
	Remark         string              `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	DeletedAt      gorm.DeletedAt      `gorm:"index" json:"-"`
}

// TableName 固定消息提醒规则表名。
func (MessageReminder) TableName() string {
	return "sys_message_reminder"
}
