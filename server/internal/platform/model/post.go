package model

import (
	"time"

	"gorm.io/gorm"
)

// PostStatus 表示岗位状态。
type PostStatus int

const (
	// PostStatusEnabled 表示岗位可用。
	PostStatusEnabled PostStatus = 1
	// PostStatusDisabled 表示岗位已停用。
	PostStatusDisabled PostStatus = 2
)

// Post 是岗位表模型。
type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	Sort      int            `gorm:"not null;default:0" json:"sort"`
	Status    PostStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string         `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定岗位表名。
func (Post) TableName() string {
	return "sys_post"
}
