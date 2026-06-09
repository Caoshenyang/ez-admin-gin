package attachment

import (
	attachmentapi "ez-admin-gin/server/internal/modules/system/attachment/api"
	platformConfig "ez-admin-gin/server/internal/platform/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB            *gorm.DB
	Upload        platformConfig.UploadConfig
	RuntimeConfig *platformConfig.RuntimeStore
	Log           *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	service := NewService(ServiceOptions{
		DB:            opts.DB,
		Upload:        opts.Upload,
		RuntimeConfig: opts.RuntimeConfig,
		Log:           opts.Log,
	})
	attachmentapi.RegisterRoutes(group, attachmentapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
