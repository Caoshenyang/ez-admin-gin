package customer

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总 CRM 客户模块路由依赖。
type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

// RegisterRoutes 注册 CRM 客户模块路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/customers", handler.List)
	group.POST("/customers", handler.Create)
	group.POST("/customers/:id/update", handler.Update)
	group.POST("/customers/:id/status", handler.UpdateStatus)
}
