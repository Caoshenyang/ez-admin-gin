package notice

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总公告模块的路由依赖。
type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

// RegisterRoutes 注册公告模块路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/notices", handler.List)
	group.POST("/notices", handler.Create)
	group.POST("/notices/:id/update", handler.Update)
	group.POST("/notices/:id/status", handler.UpdateStatus)
}
