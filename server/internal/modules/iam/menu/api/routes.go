package api

import (
	menuapp "ez-admin-gin/server/internal/modules/iam/menu/application"
	menuinfra "ez-admin-gin/server/internal/modules/iam/menu/infra"
	"ez-admin-gin/server/internal/platform/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := menuinfra.NewRepository(opts.DB)
	tx := database.NewTransactor(opts.DB)
	service := menuapp.NewService(tx, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/menus", handler.List)
	group.POST("/menus", handler.Create)
	group.POST("/menus/:id/update", handler.Update)
	group.POST("/menus/:id/status", handler.UpdateStatus)
	group.POST("/menus/:id/delete", handler.Delete)
}
