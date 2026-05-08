package api

import (
	dictapp "ez-admin-gin/server/internal/modules/system/dict/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteOptions struct {
	Service *dictapp.Service
	Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	handler := NewHandler(opts.Service, opts.Log)

	group.GET("/dict-types", handler.ListTypes)
	group.POST("/dict-types", handler.CreateType)
	group.POST("/dict-types/:id/update", handler.UpdateType)
	group.POST("/dict-types/:id/status", handler.UpdateTypeStatus)

	group.GET("/dict-items", handler.ListItems)
	group.POST("/dict-items", handler.CreateItem)
	group.POST("/dict-items/:id/update", handler.UpdateItem)
	group.POST("/dict-items/:id/status", handler.UpdateItemStatus)
}
