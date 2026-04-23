package model

import "time"

// LoginLogStatus 表示登录结果。
type LoginLogStatus int

const (
	// LoginLogStatusSuccess 表示登录成功。
	LoginLogStatusSuccess LoginLogStatus = 1
	// LoginLogStatusFailed 表示登录失败。
	LoginLogStatusFailed LoginLogStatus = 2
)

// LoginLog 是后台登录日志模型。
type LoginLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;default:0;index" json:"user_id"`
	Username  string         `gorm:"size:64;not null;default:'';index" json:"username"`
	Status    LoginLogStatus `gorm:"type:smallint;not null;index" json:"status"`
	Message   string         `gorm:"size:255;not null;default:''" json:"message"`
	IP        string         `gorm:"column:ip;size:64;not null;default:'';index" json:"ip"`
	UserAgent string         `gorm:"size:500;not null;default:''" json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`
}

// TableName 固定登录日志表名。
func (LoginLog) TableName() string {
	return "sys_login_log"
}
