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
	service := NewService(ServiceOptions{DB: opts.DB})
	noticeapi.RegisterRoutes(group, noticeapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
