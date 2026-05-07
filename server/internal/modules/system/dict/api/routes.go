package api

import (
	dictapp "ez-admin-gin/server/internal/modules/system/dict/application"
	dictinfra "ez-admin-gin/server/internal/modules/system/dict/infra"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := dictinfra.NewRepository(opts.DB)
	service := dictapp.NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/dict-types", handler.ListTypes)
	group.POST("/dict-types", handler.CreateType)
	group.POST("/dict-types/:id/update", handler.UpdateType)
	group.POST("/dict-types/:id/status", handler.UpdateTypeStatus)

	group.GET("/dict-items", handler.ListItems)
	group.POST("/dict-items", handler.CreateItem)
	group.POST("/dict-items/:id/update", handler.UpdateItem)
	group.POST("/dict-items/:id/status", handler.UpdateItemStatus)
}
