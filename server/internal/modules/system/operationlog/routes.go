package operationlog

import (
	operationlogapi "ez-admin-gin/server/internal/modules/system/operationlog/api"

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
	operationlogapi.RegisterRoutes(group, operationlogapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
