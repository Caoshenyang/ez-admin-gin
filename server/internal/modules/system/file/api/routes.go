package api

import (
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
	repo := fileinfra.NewRepository(opts.DB)
	service := fileapp.NewService(opts.DB, repo, opts.Upload, opts.Log)
	handler := NewHandler(service, opts.Log)

	group.GET("/files", handler.List)
	group.POST("/files", handler.Upload)
}
