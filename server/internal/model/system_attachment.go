package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemAttachmentStatus 表示附件中心记录状态。
type SystemAttachmentStatus int

const (
	// SystemAttachmentStatusEnabled 表示附件可正常使用。
	SystemAttachmentStatusEnabled SystemAttachmentStatus = 1
	// SystemAttachmentStatusDisabled 表示附件已停用。
	SystemAttachmentStatusDisabled SystemAttachmentStatus = 2
)

// SystemAttachment 是附件中心记录模型。
type SystemAttachment struct {
	ID          uint                   `gorm:"primaryKey" json:"id"`
	FileID      uint                   `gorm:"not null;uniqueIndex" json:"file_id"`
	DisplayName string                 `gorm:"size:255;not null" json:"display_name"`
	Category    string                 `gorm:"size:64;not null;default:'';index" json:"category"`
	BizType     string                 `gorm:"size:64;not null;default:'';index" json:"biz_type"`
	UploaderID  uint                   `gorm:"not null;default:0;index" json:"uploader_id"`
	Status      SystemAttachmentStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark      string                 `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `gorm:"index" json:"-"`
}

// TableName 固定附件中心表名。
func (SystemAttachment) TableName() string {
	return "sys_attachment"
}
