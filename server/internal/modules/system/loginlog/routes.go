package loginlog

import (
	loginlogapi "ez-admin-gin/server/internal/modules/system/loginlog/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	loginlogapi.RegisterRoutes(group, loginlogapi.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
