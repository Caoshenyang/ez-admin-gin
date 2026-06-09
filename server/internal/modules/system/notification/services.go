package notification

import (
	notiapp "ez-admin-gin/server/internal/modules/system/notification/application"
	notiinfra "ez-admin-gin/server/internal/modules/system/notification/infra"
	notiws "ez-admin-gin/server/internal/modules/system/notification/ws"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB    *gorm.DB
	Redis *redis.Client
	Log   *zap.Logger
}

// NewService 收拢通知模块的依赖装配。
func NewService(opts ServiceOptions) *notiapp.Service {
	repo := notiinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return notiapp.NewService(transactor, repo)
}

// NewHub 创建 WebSocket Hub。
func NewHub(service *notiapp.Service, rdb *redis.Client, log *zap.Logger) *notiws.Hub {
	transport := notiws.NewRedisTransport(rdb, log)
	return notiws.NewHub(service, transport, log)
}
