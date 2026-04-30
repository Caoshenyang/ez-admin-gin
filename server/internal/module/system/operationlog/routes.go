package operationlog

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总操作日志模块的路由依赖。
type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

// RegisterRoutes 注册操作日志模块路由。
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/operation-logs", handler.List)
}
