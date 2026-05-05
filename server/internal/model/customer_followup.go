package model

import (
	"time"

	"gorm.io/gorm"
)

// CustomerFollowUpStatus 表示客户跟进状态。
type CustomerFollowUpStatus int

const (
	// CustomerFollowUpStatusPending 表示待继续跟进。
	CustomerFollowUpStatusPending CustomerFollowUpStatus = 1
	// CustomerFollowUpStatusCompleted 表示本次跟进已完成。
	CustomerFollowUpStatusCompleted CustomerFollowUpStatus = 2
	// CustomerFollowUpStatusClosed 表示该条跟进已关闭。
	CustomerFollowUpStatusClosed CustomerFollowUpStatus = 3
)

// CustomerFollowUp 是 CRM 客户跟进记录模型。
type CustomerFollowUp struct {
	ID           uint                   `gorm:"primaryKey" json:"id"`
	CustomerID   uint                   `gorm:"not null;index" json:"customer_id"`
	DepartmentID uint                   `gorm:"not null;default:0;index" json:"department_id"`
	OwnerUserID  uint                   `gorm:"not null;default:0;index" json:"owner_user_id"`
	FollowType   string                 `gorm:"size:32;not null;default:'';index" json:"follow_type"`
	Subject      string                 `gorm:"size:128;not null" json:"subject"`
	Content      string                 `gorm:"size:1000;not null" json:"content"`
	Result       string                 `gorm:"size:255;not null;default:''" json:"result"`
	NextFollowAt *time.Time             `json:"next_follow_at"`
	Status       CustomerFollowUpStatus `gorm:"type:smallint;not null;default:1;index" json:"status"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    gorm.DeletedAt         `gorm:"index" json:"-"`
}

// TableName 固定 CRM 客户跟进表名。
func (CustomerFollowUp) TableName() string {
	return "sys_customer_followup"
}
