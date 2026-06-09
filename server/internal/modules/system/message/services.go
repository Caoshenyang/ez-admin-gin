package message

import (
	messageapp "ez-admin-gin/server/internal/modules/system/message/application"
	messageinfra "ez-admin-gin/server/internal/modules/system/message/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢消息提醒配置模块的依赖装配，让 routes 层保持统一入口。
func NewService(opts ServiceOptions) *messageapp.Service {
	repo := messageinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return messageapp.NewService(transactor, repo)
}
