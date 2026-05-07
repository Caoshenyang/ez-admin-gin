package model

import (
	"time"

	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

// RoleStatus 表示角色状态。
type RoleStatus int

const (
	// RoleStatusEnabled 表示角色可以正常使用。
	RoleStatusEnabled RoleStatus = 1
	// RoleStatusDisabled 表示角色已被禁用。
	RoleStatusDisabled RoleStatus = 2
)

// Role 是后台角色表模型。
type Role struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Code      string          `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Name      string          `gorm:"size:64;not null" json:"name"`
	Sort      int             `gorm:"not null;default:0" json:"sort"`
	DataScope datascope.Scope `gorm:"size:32;not null;default:'self'" json:"data_scope"`
	Status    RoleStatus      `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string          `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 固定角色表名，避免后续调整命名策略时影响已有表。
func (Role) TableName() string {
	return "sys_role"
}
