package menu

import (
	menuapi "ez-admin-gin/server/internal/modules/iam/menu/api"

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
	menuapi.RegisterRoutes(group, menuapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
