package api

import (
	fileapp "ez-admin-gin/server/internal/modules/system/file/application"
	fileinfra "ez-admin-gin/server/internal/modules/system/file/infra"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

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
	storage := fileinfra.NewLocalStorage(opts.Upload)
	service := fileapp.NewService(platformDatabase.NewTransactor(opts.DB), repo, storage, opts.Upload, opts.Log)
	handler := NewHandler(service, opts.Log)

	group.GET("/files", handler.List)
	group.POST("/files", handler.Upload)
}
