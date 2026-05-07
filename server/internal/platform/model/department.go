package model

import (
	"time"

	"gorm.io/gorm"
)

// DepartmentStatus 表示部门状态。
type DepartmentStatus int

const (
	// DepartmentStatusEnabled 表示部门可用。
	DepartmentStatusEnabled DepartmentStatus = 1
	// DepartmentStatusDisabled 表示部门已停用。
	DepartmentStatusDisabled DepartmentStatus = 2
)

// Department 是组织部门表模型。
type Department struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	ParentID     uint             `gorm:"not null;default:0;index" json:"parent_id"`
	Ancestors    string           `gorm:"size:500;not null;default:''" json:"ancestors"`
	Name         string           `gorm:"size:64;not null" json:"name"`
	Code         string           `gorm:"size:64;not null;uniqueIndex" json:"code"`
	LeaderUserID uint             `gorm:"not null;default:0;index" json:"leader_user_id"`
	Sort         int              `gorm:"not null;default:0" json:"sort"`
	Status       DepartmentStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark       string           `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`
}

// TableName 固定部门表名。
func (Department) TableName() string {
	return "sys_department"
}
