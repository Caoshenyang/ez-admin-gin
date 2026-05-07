package api

import (
	configapp "ez-admin-gin/server/internal/modules/system/config/application"
	configinfra "ez-admin-gin/server/internal/modules/system/config/infra"

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
	service := configapp.NewService(opts.DB, repo, opts.Redis, opts.Log)
	handler := NewHandler(service, opts.Log)

	group.GET("/configs", handler.List)
	group.POST("/configs", handler.Create)
	group.POST("/configs/:id/update", handler.Update)
	group.POST("/configs/:id/status", handler.UpdateStatus)
	group.GET("/configs/value/:key", handler.Value)
}
