package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemDictStatus 表示字典类型和字典项状态。
type SystemDictStatus int

const (
	// SystemDictStatusEnabled 表示字典处于启用状态。
	SystemDictStatusEnabled SystemDictStatus = 1
	// SystemDictStatusDisabled 表示字典处于禁用状态。
	SystemDictStatusDisabled SystemDictStatus = 2
)

// SystemDictType 是字典类型表模型。
type SystemDictType struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	Code      string           `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Name      string           `gorm:"size:64;not null" json:"name"`
	Sort      int              `gorm:"not null;default:0" json:"sort"`
	Status    SystemDictStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string           `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}

// TableName 固定字典类型表名。
func (SystemDictType) TableName() string {
	return "sys_dict_type"
}
