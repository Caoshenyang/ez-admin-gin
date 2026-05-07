package user

import (
	userapi "ez-admin-gin/server/internal/modules/iam/user/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	userapi.RegisterRoutes(group, userapi.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
