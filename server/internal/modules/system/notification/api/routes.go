package api

import (
	notiapp "ez-admin-gin/server/internal/modules/system/notification/application"
	notiws "ez-admin-gin/server/internal/modules/system/notification/ws"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *notiapp.Service
	Hub     *notiws.Hub
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Hub, opts.Log)

	group.GET("/notifications", handler.List)
	group.GET("/notifications/unread-count", handler.UnreadCount)
	group.POST("/notifications/mark-read", handler.MarkRead)
	group.POST("/notifications/mark-all-read", handler.MarkAllRead)
	// WebSocket 端点注册在系统路由组之外（不走 Auth 中间件链）
}
