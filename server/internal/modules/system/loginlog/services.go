package loginlog

import (
	loginlogapp "ez-admin-gin/server/internal/modules/system/loginlog/application"
	loginloginfra "ez-admin-gin/server/internal/modules/system/loginlog/infra"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢 loginlog 模块依赖装配，让 routes 层保持统一入口。
func NewService(opts ServiceOptions) *loginlogapp.Service {
	repo := loginloginfra.NewRepository(opts.DB)
	return loginlogapp.NewService(repo)
}
