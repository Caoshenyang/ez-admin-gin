package file

import (
	"ez-admin-gin/server/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总文件模块的路由依赖。
type RouteOptions struct {
	DB     *gorm.DB
	Upload config.UploadConfig
	Log    *zap.Logger
}

// RegisterRoutes 注册文件模块路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo, opts.Upload, opts.Log)
	handler := NewHandler(service, opts.Log)

	group.GET("/files", handler.List)
	group.POST("/files", handler.Upload)
}
