package attachment

import (
	attachmentapi "ez-admin-gin/server/internal/modules/system/attachment/api"
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
	attachmentapi.RegisterRoutes(group, attachmentapi.RouteOptions{
		DB:     opts.DB,
		Upload: opts.Upload,
		Log:    opts.Log,
	})
}
