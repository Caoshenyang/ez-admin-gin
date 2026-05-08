package file

import (
	fileapi "ez-admin-gin/server/internal/modules/system/file/api"
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
	service := NewService(ServiceOptions{
		DB:     opts.DB,
		Upload: opts.Upload,
		Log:    opts.Log,
	})
	fileapi.RegisterRoutes(group, fileapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
