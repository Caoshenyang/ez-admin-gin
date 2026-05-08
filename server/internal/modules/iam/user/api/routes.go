package api

import (
	userapp "ez-admin-gin/server/internal/modules/iam/user/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *userapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/users", handler.List)
	group.POST("/users", handler.Create)
	group.POST("/users/:id/update", handler.Update)
	group.POST("/users/:id/status", handler.UpdateStatus)
	group.POST("/users/:id/roles", handler.UpdateRoles)
}
