package api

import (
	departmentapp "ez-admin-gin/server/internal/modules/iam/department/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *departmentapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/departments", handler.List)
	group.POST("/departments", handler.Create)
	group.POST("/departments/:id/update", handler.Update)
	group.POST("/departments/:id/status", handler.UpdateStatus)
}
