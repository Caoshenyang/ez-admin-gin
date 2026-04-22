package model

import (
	"time"

	"gorm.io/gorm"
)

// UserStatus 表示用户状态。
type UserStatus int

const (
	// UserStatusEnabled 表示用户可以正常登录。
	UserStatusEnabled UserStatus = 1
	// UserStatusDisabled 表示用户已被禁用。
	UserStatusDisabled UserStatus = 2
)

// User 是后台用户表模型。
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Nickname     string         `gorm:"size:64;not null;default:''" json:"nickname"`
	Status       UserStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定用户表名，避免后续调整命名策略时影响已有表。
func (User) TableName() string {
	return "sys_user"
}
