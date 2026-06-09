package model

import (
	"time"

	"gorm.io/gorm"
)

// MailAccountStatus 表示系统邮箱账号状态。
type MailAccountStatus int

const (
	// MailAccountStatusEnabled 表示邮箱账号可用于发信。
	MailAccountStatusEnabled MailAccountStatus = 1
	// MailAccountStatusDisabled 表示邮箱账号已停用。
	MailAccountStatusDisabled MailAccountStatus = 2
)

// MailEncryption 表示 SMTP 连接加密方式。
type MailEncryption string

const (
	// MailEncryptionNone 表示不启用连接加密。
	MailEncryptionNone MailEncryption = "none"
	// MailEncryptionSSL 表示连接建立时直接使用 SSL/TLS。
	MailEncryptionSSL MailEncryption = "ssl"
	// MailEncryptionSTARTTLS 表示连接后升级到 STARTTLS。
	MailEncryptionSTARTTLS MailEncryption = "starttls"
)

// MailAccount 是系统邮箱账号配置表模型。
type MailAccount struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	Name        string            `gorm:"size:64;not null" json:"name"`
	Host        string            `gorm:"size:128;not null" json:"host"`
	Port        int               `gorm:"not null;default:25" json:"port"`
	Username    string            `gorm:"size:128;not null;default:''" json:"username"`
	Password    string            `gorm:"size:255;not null;default:''" json:"-"`
	FromEmail   string            `gorm:"size:128;not null" json:"from_email"`
	FromName    string            `gorm:"size:64;not null;default:''" json:"from_name"`
	Encryption  MailEncryption    `gorm:"size:16;not null;default:'none'" json:"encryption"`
	IsDefault   bool              `gorm:"not null;default:false;index" json:"is_default"`
	Status      MailAccountStatus `gorm:"type:smallint;not null;default:1;index" json:"status"`
	LastTestAt  *time.Time        `json:"last_test_at"`
	LastTestMsg string            `gorm:"size:255;not null;default:''" json:"last_test_msg"`
	Remark      string            `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
}

// TableName 固定系统邮箱账号表名。
func (MailAccount) TableName() string {
	return "sys_mail_account"
}
