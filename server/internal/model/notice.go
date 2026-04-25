package model

import (
	"time"

	"gorm.io/gorm"
)

// NoticeStatus 表示公告状态。
type NoticeStatus int

const (
	// NoticeStatusEnabled 表示公告可见。
	NoticeStatusEnabled NoticeStatus = 1
	// NoticeStatusDisabled 表示公告已隐藏。
	NoticeStatusDisabled NoticeStatus = 2
)

// Notice 是公告表模型。
type Notice struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Title     string        `gorm:"size:128;not null" json:"title"`
	Content   string        `gorm:"type:text;not null" json:"content"`
	Sort      int           `gorm:"not null;default:0" json:"sort"`
	Status    NoticeStatus  `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string        `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定公告表名。
func (Notice) TableName() string {
	return "sys_notice"
}
