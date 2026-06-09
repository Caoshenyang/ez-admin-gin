package mail

import (
	mailapp "ez-admin-gin/server/internal/modules/system/mail/application"
	mailinfra "ez-admin-gin/server/internal/modules/system/mail/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢邮件模块依赖装配，保持 routes 层只关注模块入口。
func NewService(opts ServiceOptions) *mailapp.Service {
	repo := mailinfra.NewRepository(opts.DB)
	sender := mailinfra.NewSMTPSender()
	transactor := platformDatabase.NewTransactor(opts.DB)
	return mailapp.NewService(transactor, repo, sender)
}
