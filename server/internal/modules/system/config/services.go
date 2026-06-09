package config

import (
	configapp "ez-admin-gin/server/internal/modules/system/config/application"
	configinfra "ez-admin-gin/server/internal/modules/system/config/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB    *gorm.DB
	Redis *goredis.Client
	Log   *zap.Logger
}

// NewService 收拢系统配置模块的依赖装配，保持 routes 层只关注模块入口。
func NewService(opts ServiceOptions) *configapp.Service {
	repo := configinfra.NewRepository(opts.DB)
	cache := configinfra.NewCache(opts.Redis)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return configapp.NewService(transactor, repo, cache, opts.Log)
}
