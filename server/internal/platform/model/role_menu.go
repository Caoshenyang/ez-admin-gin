package model

import "time"

// RoleMenu 是角色和菜单的绑定关系。
type RoleMenu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoleID    uint      `gorm:"not null;uniqueIndex:uk_sys_role_menu_role_menu;index:idx_sys_role_menu_role_id" json:"role_id"`
	MenuID    uint      `gorm:"not null;uniqueIndex:uk_sys_role_menu_role_menu;index:idx_sys_role_menu_menu_id" json:"menu_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 固定角色菜单关系表名。
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}
