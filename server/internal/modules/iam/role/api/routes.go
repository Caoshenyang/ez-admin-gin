package api

import (
	roleapp "ez-admin-gin/server/internal/modules/iam/role/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *roleapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/roles", handler.List)
	group.POST("/roles", handler.Create)
	group.POST("/roles/:id/update", handler.Update)
	group.POST("/roles/:id/status", handler.UpdateStatus)
	group.POST("/roles/:id/permissions", handler.UpdatePermissions)
	group.POST("/roles/:id/menus", handler.UpdateMenus)
	group.POST("/roles/:id/delete", handler.Delete)
}
