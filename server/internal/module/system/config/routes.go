package config

import (
	goredis "github.com/redis/go-redis/v9"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总系统配置模块的路由依赖。
type RouteOptions struct {
	DB    *gorm.DB
	Redis *goredis.Client
	Log   *zap.Logger
}

// RegisterRoutes 注册系统配置模块路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo, opts.Redis, opts.Log)
	handler := NewHandler(service, opts.Log)

	group.GET("/configs", handler.List)
	group.POST("/configs", handler.Create)
	group.POST("/configs/:id/update", handler.Update)
	group.POST("/configs/:id/status", handler.UpdateStatus)
	group.GET("/configs/value/:key", handler.Value)
}
