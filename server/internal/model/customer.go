package model

import (
	"time"

	"gorm.io/gorm"
)

// CustomerStatus 表示客户状态。
type CustomerStatus int

const (
	// CustomerStatusEnabled 表示客户记录可正常使用。
	CustomerStatusEnabled CustomerStatus = 1
	// CustomerStatusDisabled 表示客户记录已停用。
	CustomerStatusDisabled CustomerStatus = 2
)

// Customer 是 CRM 客户档案表模型。
type Customer struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:128;not null" json:"name"`
	ContactName  string         `gorm:"size:64;not null;default:''" json:"contact_name"`
	Phone        string         `gorm:"size:32;not null;default:''" json:"phone"`
	Level        string         `gorm:"size:32;not null;default:'';index" json:"level"`
	Source       string         `gorm:"size:32;not null;default:'';index" json:"source"`
	DepartmentID uint           `gorm:"not null;default:0;index" json:"department_id"`
	OwnerUserID  uint           `gorm:"not null;default:0;index" json:"owner_user_id"`
	Status       CustomerStatus `gorm:"type:smallint;not null;default:1;index" json:"status"`
	Remark       string         `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定 CRM 客户表名。
func (Customer) TableName() string {
	return "sys_customer"
}
