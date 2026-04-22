package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemConfigStatus 表示系统配置状态。
type SystemConfigStatus int

const (
	// SystemConfigStatusEnabled 表示配置可用。
	SystemConfigStatusEnabled SystemConfigStatus = 1
	// SystemConfigStatusDisabled 表示配置已停用。
	SystemConfigStatusDisabled SystemConfigStatus = 2
)

// SystemConfig 是系统配置表模型。
type SystemConfig struct {
	ID        uint               `gorm:"primaryKey" json:"id"`
	GroupCode string             `gorm:"size:64;not null;index" json:"group_code"`
	ConfigKey string             `gorm:"column:config_key;size:128;not null;uniqueIndex" json:"key"`
	Name      string             `gorm:"size:64;not null" json:"name"`
	Value     string             `gorm:"type:text;not null" json:"value"`
	Sort      int                `gorm:"not null;default:0" json:"sort"`
	Status    SystemConfigStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string             `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
}

// TableName 固定系统配置表名。
func (SystemConfig) TableName() string {
	return "sys_config"
}
