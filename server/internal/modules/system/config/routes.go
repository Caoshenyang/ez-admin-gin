package config

import (
	configapi "ez-admin-gin/server/internal/modules/system/config/api"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB    *gorm.DB
	Redis *goredis.Client
	Log   *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	service := NewService(ServiceOptions{
		DB:    opts.DB,
		Redis: opts.Redis,
		Log:   opts.Log,
	})
	configapi.RegisterRoutes(group, configapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
