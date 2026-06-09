package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemDictItem 是字典项表模型。
type SystemDictItem struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	TypeID    uint             `gorm:"column:type_id;not null;index:idx_sys_dict_item_type_id;uniqueIndex:uk_sys_dict_item_type_key" json:"type_id"`
	ItemKey   string           `gorm:"column:item_key;size:64;not null;uniqueIndex:uk_sys_dict_item_type_key" json:"item_key"`
	Label     string           `gorm:"size:64;not null" json:"label"`
	Value     string           `gorm:"size:255;not null" json:"value"`
	TagType   string           `gorm:"column:tag_type;size:32;not null;default:''" json:"tag_type"`
	Sort      int              `gorm:"not null;default:0" json:"sort"`
	Status    SystemDictStatus `gorm:"type:smallint;not null;default:1;index" json:"status"`
	Remark    string           `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}

// TableName 固定字典项表名。
func (SystemDictItem) TableName() string {
	return "sys_dict_item"
}
