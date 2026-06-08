package api

import (
	postapp "ez-admin-gin/server/internal/modules/iam/post/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *postapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/posts", handler.List)
	group.POST("/posts", handler.Create)
	group.POST("/posts/:id/update", handler.Update)
	group.POST("/posts/:id/status", handler.UpdateStatus)
	group.POST("/posts/:id/delete", handler.Delete)
}
