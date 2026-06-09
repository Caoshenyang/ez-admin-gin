package api

import (
	mailapp "ez-admin-gin/server/internal/modules/system/mail/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *mailapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/mail/accounts", handler.ListAccounts)
	group.POST("/mail/accounts", handler.CreateAccount)
	group.POST("/mail/accounts/:id/update", handler.UpdateAccount)
	group.POST("/mail/accounts/:id/status", handler.UpdateAccountStatus)
	group.POST("/mail/accounts/:id/delete", handler.DeleteAccount)
	group.POST("/mail/accounts/:id/test", handler.TestAccount)

	group.GET("/mail/templates", handler.ListTemplates)
	group.POST("/mail/templates", handler.CreateTemplate)
	group.POST("/mail/templates/:id/update", handler.UpdateTemplate)
	group.POST("/mail/templates/:id/status", handler.UpdateTemplateStatus)
	group.POST("/mail/templates/:id/delete", handler.DeleteTemplate)
	group.POST("/mail/templates/:id/render", handler.RenderTemplate)

	group.POST("/mail/send", handler.Send)
	group.GET("/mail/logs", handler.ListLogs)
}
