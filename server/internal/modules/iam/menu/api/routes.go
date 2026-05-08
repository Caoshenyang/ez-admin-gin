package api

import (
	menuapp "ez-admin-gin/server/internal/modules/iam/menu/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *menuapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/menus", handler.List)
	group.POST("/menus", handler.Create)
	group.POST("/menus/:id/update", handler.Update)
	group.POST("/menus/:id/status", handler.UpdateStatus)
	group.POST("/menus/:id/delete", handler.Delete)
}
