package model

import "time"

// RoleAPI 是角色和接口权限元数据的绑定关系。
type RoleAPI struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoleID    uint      `gorm:"not null;uniqueIndex:uk_sys_role_api_role_api;index:idx_sys_role_api_role_id" json:"role_id"`
	APIID     uint      `gorm:"not null;uniqueIndex:uk_sys_role_api_role_api;index:idx_sys_role_api_api_id" json:"api_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 固定角色接口关系表名。
func (RoleAPI) TableName() string {
	return "sys_role_api"
}
