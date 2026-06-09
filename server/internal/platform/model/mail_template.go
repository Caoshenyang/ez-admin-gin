package model

import (
	"time"

	"gorm.io/gorm"
)

// MailTemplateStatus 表示邮件模板状态。
type MailTemplateStatus int

const (
	// MailTemplateStatusEnabled 表示模板可用于发信。
	MailTemplateStatusEnabled MailTemplateStatus = 1
	// MailTemplateStatusDisabled 表示模板已停用。
	MailTemplateStatusDisabled MailTemplateStatus = 2
)

// MailTemplate 是邮件模板表模型。
type MailTemplate struct {
	ID        uint               `gorm:"primaryKey" json:"id"`
	Code      string             `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Name      string             `gorm:"size:64;not null" json:"name"`
	Subject   string             `gorm:"size:200;not null" json:"subject"`
	Content   string             `gorm:"type:text;not null" json:"content"`
	IsHTML    bool               `gorm:"column:is_html;not null;default:true" json:"is_html"`
	Variables string             `gorm:"type:text;not null" json:"variables"`
	Sort      int                `gorm:"not null;default:0" json:"sort"`
	Status    MailTemplateStatus `gorm:"type:smallint;not null;default:1;index" json:"status"`
	Remark    string             `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
}

// TableName 固定邮件模板表名。
func (MailTemplate) TableName() string {
	return "sys_mail_template"
}
