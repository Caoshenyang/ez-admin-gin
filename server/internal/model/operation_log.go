package model

import "time"

// OperationLog 是后台操作日志模型。
type OperationLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;default:0;index" json:"user_id"`
	Username     string    `gorm:"size:64;not null;default:'';index" json:"username"`
	Method       string    `gorm:"size:10;not null;index" json:"method"`
	Path         string    `gorm:"size:255;not null;index" json:"path"`
	RoutePath    string    `gorm:"size:255;not null;default:'';index" json:"route_path"`
	Query        string    `gorm:"size:1000;not null;default:''" json:"query"`
	IP           string    `gorm:"column:ip;size:64;not null;default:''" json:"ip"`
	UserAgent    string    `gorm:"size:500;not null;default:''" json:"user_agent"`
	StatusCode   int       `gorm:"not null;default:0;index" json:"status_code"`
	LatencyMs    int64     `gorm:"not null;default:0" json:"latency_ms"`
	Success      bool      `gorm:"not null;default:true;index" json:"success"`
	ErrorMessage string    `gorm:"size:500;not null;default:''" json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 固定操作日志表名。
func (OperationLog) TableName() string {
	return "sys_operation_log"
}
