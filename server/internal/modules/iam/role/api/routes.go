package api

import (
	roleapp "ez-admin-gin/server/internal/modules/iam/role/application"
	roleinfra "ez-admin-gin/server/internal/modules/iam/role/infra"
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
	repo := roleinfra.NewRepository(opts.DB)
	tx := database.NewTransactor(opts.DB)
	service := roleapp.NewService(tx, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/roles", handler.List)
	group.POST("/roles", handler.Create)
	group.POST("/roles/:id/update", handler.Update)
	group.POST("/roles/:id/status", handler.UpdateStatus)
	group.POST("/roles/:id/permissions", handler.UpdatePermissions)
	group.POST("/roles/:id/menus", handler.UpdateMenus)
}
