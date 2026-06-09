package model

import (
	"time"

	"gorm.io/gorm"
)

// APIStatus 表示接口元数据状态。
type APIStatus int

const (
	// APIStatusEnabled 表示接口权限正常启用。
	APIStatusEnabled APIStatus = 1
	// APIStatusDisabled 表示接口权限已禁用。
	APIStatusDisabled APIStatus = 2
)

// API 是后台接口权限元数据模型。
type API struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"size:128;not null;uniqueIndex" json:"code"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	Module    string         `gorm:"size:64;not null;default:'';index" json:"module"`
	Method    string         `gorm:"size:16;not null;uniqueIndex:uk_sys_api_method_path" json:"method"`
	Path      string         `gorm:"size:255;not null;uniqueIndex:uk_sys_api_method_path" json:"path"`
	Sort      int            `gorm:"not null;default:0" json:"sort"`
	Status    APIStatus      `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string         `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定接口元数据表名。
func (API) TableName() string {
	return "sys_api"
}
