package api

import (
	noticeapp "ez-admin-gin/server/internal/modules/system/notice/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *noticeapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/notices", handler.List)
	group.POST("/notices", handler.Create)
	group.POST("/notices/:id/update", handler.Update)
	group.POST("/notices/:id/status", handler.UpdateStatus)
	group.POST("/notices/:id/delete", handler.Delete)
}
