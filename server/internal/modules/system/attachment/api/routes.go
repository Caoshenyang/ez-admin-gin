package api

import (
	attachmentapp "ez-admin-gin/server/internal/modules/system/attachment/application"
	attachmentinfra "ez-admin-gin/server/internal/modules/system/attachment/infra"
	fileapp "ez-admin-gin/server/internal/modules/system/file/application"
	fileinfra "ez-admin-gin/server/internal/modules/system/file/infra"
	platformConfig "ez-admin-gin/server/internal/platform/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB     *gorm.DB
	Upload platformConfig.UploadConfig
	Log    *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := attachmentinfra.NewRepository(opts.DB)
	fileRepo := fileinfra.NewRepository(opts.DB)
	fileService := fileapp.NewService(opts.DB, fileRepo, opts.Upload, opts.Log)
	service := attachmentapp.NewService(opts.DB, repo, fileService)
	handler := NewHandler(service, opts.Log)

	group.GET("/attachments", handler.List)
	group.POST("/attachments", handler.Upload)
	group.POST("/attachments/:id/update", handler.Update)
	group.POST("/attachments/:id/status", handler.UpdateStatus)
}
