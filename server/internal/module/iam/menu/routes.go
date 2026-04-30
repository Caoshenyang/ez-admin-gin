package menu

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总菜单模块的路由依赖。
type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

// RegisterRoutes 注册菜单模块路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/menus", handler.List)
	group.POST("/menus", handler.Create)
	group.POST("/menus/:id/update", handler.Update)
	group.POST("/menus/:id/status", handler.UpdateStatus)
	group.POST("/menus/:id/delete", handler.Delete)
}
