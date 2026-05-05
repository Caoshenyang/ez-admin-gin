package attachment

import (
	"ez-admin-gin/server/internal/config"
	systemFileModule "ez-admin-gin/server/internal/module/system/file"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总附件中心模块依赖。
type RouteOptions struct {
	DB     *gorm.DB
	Upload config.UploadConfig
	Log    *zap.Logger
}

// RegisterRoutes 注册附件中心路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	fileRepo := systemFileModule.NewRepository(opts.DB)
	fileService := systemFileModule.NewService(opts.DB, fileRepo, opts.Upload, opts.Log)
	service := NewService(opts.DB, repo, fileService)
	handler := NewHandler(service, opts.Log)

	group.GET("/attachments", handler.List)
	group.POST("/attachments", handler.Upload)
	group.POST("/attachments/:id/update", handler.Update)
	group.POST("/attachments/:id/status", handler.UpdateStatus)
}
