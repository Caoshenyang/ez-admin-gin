package dict

import (
	dictapi "ez-admin-gin/server/internal/modules/system/dict/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	service := NewService(ServiceOptions{DB: opts.DB})
	dictapi.RegisterRoutes(group, dictapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
