package followup

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总 CRM 客户跟进模块路由依赖。
type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

// RegisterRoutes 注册 CRM 客户跟进模块路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/followups/customer-options", handler.ListCustomerOptions)
	group.GET("/followups", handler.List)
	group.POST("/followups", handler.Create)
	group.POST("/followups/:id/update", handler.Update)
	group.POST("/followups/:id/status", handler.UpdateStatus)
}
