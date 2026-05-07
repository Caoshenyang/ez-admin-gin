package notice

import (
	noticeapi "ez-admin-gin/server/internal/modules/system/notice/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	noticeapi.RegisterRoutes(group, noticeapi.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
