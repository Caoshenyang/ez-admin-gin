package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总用户模块路由需要的依赖。
type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

// RegisterRoutes 注册用户模块路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/users", handler.List)
	group.POST("/users", handler.Create)
	group.POST("/users/:id/update", handler.Update)
	group.POST("/users/:id/status", handler.UpdateStatus)
	group.POST("/users/:id/roles", handler.UpdateRoles)
}
