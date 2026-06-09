package model

import "time"

// RoleDataScope 保存角色到自定义部门范围的绑定关系。
type RoleDataScope struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"not null;index" json:"role_id"`
	DepartmentID uint      `gorm:"not null;index" json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 固定角色数据范围关系表名。
func (RoleDataScope) TableName() string {
	return "sys_role_data_scope"
}
