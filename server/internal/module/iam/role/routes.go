package role

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/roles", handler.List)
	group.POST("/roles", handler.Create)
	group.POST("/roles/:id/update", handler.Update)
	group.POST("/roles/:id/status", handler.UpdateStatus)
	group.POST("/roles/:id/permissions", handler.UpdatePermissions)
	group.POST("/roles/:id/menus", handler.UpdateMenus)
}
