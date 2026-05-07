package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemFileStatus 表示文件记录状态。
type SystemFileStatus int

const (
	// SystemFileStatusEnabled 表示文件可正常使用。
	SystemFileStatusEnabled SystemFileStatus = 1
	// SystemFileStatusDisabled 表示文件已停用。
	SystemFileStatusDisabled SystemFileStatus = 2
)

// SystemFile 是文件上传记录模型。
type SystemFile struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	Storage      string           `gorm:"size:32;not null;default:'local'" json:"storage"`
	OriginalName string           `gorm:"size:255;not null" json:"original_name"`
	FileName     string           `gorm:"size:255;not null" json:"file_name"`
	Ext          string           `gorm:"size:32;not null;default:'';index" json:"ext"`
	MimeType     string           `gorm:"size:128;not null;default:''" json:"mime_type"`
	Size         int64            `gorm:"not null;default:0" json:"size"`
	Sha256       string           `gorm:"size:64;not null;default:'';index" json:"sha256"`
	Path         string           `gorm:"size:500;not null" json:"path"`
	URL          string           `gorm:"column:url;size:500;not null" json:"url"`
	UploaderID   uint             `gorm:"not null;default:0;index" json:"uploader_id"`
	Status       SystemFileStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark       string           `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`
}

// TableName 固定文件记录表名。
func (SystemFile) TableName() string {
	return "sys_file"
}
