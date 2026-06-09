package model

import "time"

// UserRole 是用户与角色的绑定关系。
type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:uk_sys_user_role_user_role;index:idx_sys_user_role_user_id" json:"user_id"`
	RoleID    uint      `gorm:"not null;uniqueIndex:uk_sys_user_role_user_role;index:idx_sys_user_role_role_id" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 固定用户角色关系表名。
func (UserRole) TableName() string {
	return "sys_user_role"
}
