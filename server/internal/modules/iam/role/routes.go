package role

import (
	roleapi "ez-admin-gin/server/internal/modules/iam/role/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	roleapi.RegisterRoutes(group, roleapi.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
