package application

import (
	maildomain "ez-admin-gin/server/internal/modules/system/mail/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type MailTransactor = database.Transactor

type MailRepository interface {
	ListAccounts(query maildomain.AccountListQuery, page int, pageSize int, status *model.MailAccountStatus) ([]maildomain.AccountEntity, int64, error)
	FindAccountByID(db *gorm.DB, accountID uint) (maildomain.AccountEntity, error)
	FindDefaultAccount() (maildomain.AccountEntity, error)
	FindEnabledAccountByID(accountID uint) (maildomain.AccountEntity, error)
	AccountNameExists(db *gorm.DB, name string, excludeID uint) (bool, error)
	CreateAccount(db *gorm.DB, item *maildomain.AccountEntity) error
	UpdateAccount(db *gorm.DB, item *maildomain.AccountEntity, req maildomain.UpdateAccountRequest) error
	UpdateAccountStatus(db *gorm.DB, item *maildomain.AccountEntity, status model.MailAccountStatus) error
	UpdateAccountTestResult(db *gorm.DB, item *maildomain.AccountEntity, message string) error
	ClearDefaultAccounts(db *gorm.DB, excludeID uint) error
	DeleteAccount(db *gorm.DB, item *maildomain.AccountEntity) error

	ListTemplates(query maildomain.TemplateListQuery, page int, pageSize int, status *model.MailTemplateStatus) ([]maildomain.TemplateEntity, int64, error)
	FindTemplateByID(db *gorm.DB, templateID uint) (maildomain.TemplateEntity, error)
	FindEnabledTemplateByCode(code string) (maildomain.TemplateEntity, error)
	TemplateCodeExists(db *gorm.DB, code string) (bool, error)
	CreateTemplate(db *gorm.DB, item *maildomain.TemplateEntity) error
	UpdateTemplate(db *gorm.DB, item *maildomain.TemplateEntity, req maildomain.UpdateTemplateRequest) error
	UpdateTemplateStatus(db *gorm.DB, item *maildomain.TemplateEntity, status model.MailTemplateStatus) error
	DeleteTemplate(db *gorm.DB, item *maildomain.TemplateEntity) error

	ListLogs(query maildomain.LogListQuery, page int, pageSize int, status *model.MailLogStatus) ([]maildomain.LogEntity, int64, error)
	CreateLog(db *gorm.DB, item *maildomain.LogEntity) error
}

type Sender interface {
	Send(account model.MailAccount, message Message) error
}

type Message struct {
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Content string
	IsHTML  bool
}
