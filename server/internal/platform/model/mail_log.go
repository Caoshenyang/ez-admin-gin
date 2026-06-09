package model

import "time"

// MailLogStatus 表示邮件发送结果。
type MailLogStatus int

const (
	// MailLogStatusSuccess 表示邮件发送成功。
	MailLogStatusSuccess MailLogStatus = 1
	// MailLogStatusFailed 表示邮件发送失败。
	MailLogStatusFailed MailLogStatus = 2
)

// MailLog 是邮件发送日志表模型。
type MailLog struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	AccountID    uint          `gorm:"not null;default:0;index" json:"account_id"`
	AccountName  string        `gorm:"size:64;not null;default:''" json:"account_name"`
	TemplateID   uint          `gorm:"not null;default:0;index" json:"template_id"`
	TemplateCode string        `gorm:"size:64;not null;default:'';index" json:"template_code"`
	Subject      string        `gorm:"size:200;not null;default:''" json:"subject"`
	FromEmail    string        `gorm:"size:128;not null;default:''" json:"from_email"`
	ToEmails     string        `gorm:"type:text;not null" json:"to_emails"`
	CcEmails     string        `gorm:"type:text;not null" json:"cc_emails"`
	BccEmails    string        `gorm:"type:text;not null" json:"bcc_emails"`
	Status       MailLogStatus `gorm:"type:smallint;not null;index" json:"status"`
	ErrorMessage string        `gorm:"size:500;not null;default:''" json:"error_message"`
	CreatedAt    time.Time     `json:"created_at"`
}

// TableName 固定邮件发送日志表名。
func (MailLog) TableName() string {
	return "sys_mail_log"
}
