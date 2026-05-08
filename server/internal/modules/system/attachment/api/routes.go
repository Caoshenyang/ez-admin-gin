package api

import (
	attachmentapp "ez-admin-gin/server/internal/modules/system/attachment/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *attachmentapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/attachments", handler.List)
	group.POST("/attachments", handler.Upload)
	group.POST("/attachments/:id/update", handler.Update)
	group.POST("/attachments/:id/status", handler.UpdateStatus)
}
