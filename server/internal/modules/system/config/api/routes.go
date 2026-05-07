package api

import (
	configapp "ez-admin-gin/server/internal/modules/system/config/application"
	configinfra "ez-admin-gin/server/internal/modules/system/config/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

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
	repo := configinfra.NewRepository(opts.DB)
	cache := configinfra.NewCache(opts.Redis)
	service := configapp.NewService(platformDatabase.NewTransactor(opts.DB), repo, cache, opts.Log)
	handler := NewHandler(service, opts.Log)

	group.GET("/configs", handler.List)
	group.POST("/configs", handler.Create)
	group.POST("/configs/:id/update", handler.Update)
	group.POST("/configs/:id/status", handler.UpdateStatus)
	group.GET("/configs/value/:key", handler.Value)
}
