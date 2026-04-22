package model

import (
	"time"

	"gorm.io/gorm"
)

// MenuType 表示菜单节点类型。
type MenuType int

const (
	// MenuTypeDirectory 表示目录节点。
	MenuTypeDirectory MenuType = 1
	// MenuTypeMenu 表示可访问页面。
	MenuTypeMenu MenuType = 2
	// MenuTypeButton 表示页面内按钮或操作点。
	MenuTypeButton MenuType = 3
)

// MenuStatus 表示菜单状态。
type MenuStatus int

const (
	// MenuStatusEnabled 表示菜单正常启用。
	MenuStatusEnabled MenuStatus = 1
	// MenuStatusDisabled 表示菜单已禁用。
	MenuStatusDisabled MenuStatus = 2
)

// Menu 是后台菜单和按钮权限模型。
type Menu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ParentID  uint           `gorm:"not null;default:0;index" json:"parent_id"`
	Type      MenuType       `gorm:"type:smallint;not null" json:"type"`
	Code      string         `gorm:"size:128;not null;uniqueIndex" json:"code"`
	Title     string         `gorm:"size:64;not null" json:"title"`
	Path      string         `gorm:"size:255;not null;default:''" json:"path"`
	Component string         `gorm:"size:255;not null;default:''" json:"component"`
	Icon      string         `gorm:"size:64;not null;default:''" json:"icon"`
	Sort      int            `gorm:"not null;default:0" json:"sort"`
	Status    MenuStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string         `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定菜单表名。
func (Menu) TableName() string {
	return "sys_menu"
}
