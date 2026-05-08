package post

import (
	postapi "ez-admin-gin/server/internal/modules/iam/post/api"

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
	postapi.RegisterRoutes(group, postapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
