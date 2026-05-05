package dict

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总字典模块路由依赖。
type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

// RegisterRoutes 注册字典模块路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo)
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
