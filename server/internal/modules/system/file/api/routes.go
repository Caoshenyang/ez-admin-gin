package api

import (
	fileapp "ez-admin-gin/server/internal/modules/system/file/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *fileapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/files", handler.List)
	group.POST("/files", handler.Upload)
}
