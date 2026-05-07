package department

import (
	departmentapi "ez-admin-gin/server/internal/modules/iam/department/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	departmentapi.RegisterRoutes(group, departmentapi.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
