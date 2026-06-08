package api

import (
	apiresourceapp "ez-admin-gin/server/internal/modules/iam/apiresource/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *apiresourceapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/apis", handler.List)
}
