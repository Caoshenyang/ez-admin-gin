package apiresource

import (
	apiresourceapi "ez-admin-gin/server/internal/modules/iam/apiresource/api"

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
	apiresourceapi.RegisterRoutes(group, apiresourceapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
