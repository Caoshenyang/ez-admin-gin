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
	dictapi.RegisterRoutes(group, dictapi.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
