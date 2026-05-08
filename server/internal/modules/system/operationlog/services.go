package operationlog

import (
	operationlogapp "ez-admin-gin/server/internal/modules/system/operationlog/application"
	operationloginfra "ez-admin-gin/server/internal/modules/system/operationlog/infra"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢 operationlog 模块依赖装配，让 routes 层保持统一入口。
func NewService(opts ServiceOptions) *operationlogapp.Service {
	repo := operationloginfra.NewRepository(opts.DB)
	return operationlogapp.NewService(repo)
}
