package api

import (
	{{ .Module }}app "ez-admin-gin/server/internal/modules/{{ .Group }}/{{ .Module }}/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *{{ .Module }}app.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/{{ .Module }}s", handler.List)
	group.POST("/{{ .Module }}s", handler.Create)
	group.POST("/{{ .Module }}s/:id/update", handler.Update)
{{- if .HasStatus }}
	group.POST("/{{ .Module }}s/:id/status", handler.UpdateStatus)
{{- end }}
}
