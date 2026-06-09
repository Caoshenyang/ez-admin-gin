package api

import (
	configapp "ez-admin-gin/server/internal/modules/system/config/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *configapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/configs", handler.List)
	group.POST("/configs", handler.Create)
	group.POST("/configs/:id/update", handler.Update)
	group.POST("/configs/:id/status", handler.UpdateStatus)
	group.POST("/configs/:id/delete", handler.Delete)
	group.GET("/configs/value/:key", handler.Value)
}
