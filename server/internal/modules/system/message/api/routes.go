package api

import (
	messageapp "ez-admin-gin/server/internal/modules/system/message/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *messageapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/message-templates", handler.ListTemplates)
	group.POST("/message-templates", handler.CreateTemplate)
	group.POST("/message-templates/:id/update", handler.UpdateTemplate)
	group.POST("/message-templates/:id/status", handler.UpdateTemplateStatus)

	group.GET("/message-reminders", handler.ListReminders)
	group.POST("/message-reminders", handler.CreateReminder)
	group.POST("/message-reminders/:id/update", handler.UpdateReminder)
	group.POST("/message-reminders/:id/status", handler.UpdateReminderStatus)
}
