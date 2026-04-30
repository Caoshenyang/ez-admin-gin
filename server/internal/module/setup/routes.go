package setup

import (
	legacySetupHandler "ez-admin-gin/server/internal/handler/setup"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总初始化模块路由依赖。
type RouteOptions struct {
	Log *zap.Logger
	DB  *gorm.DB
}

// RegisterRoutes 注册系统初始化路由。
func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	setup := legacySetupHandler.NewSetupHandler(opts.DB, opts.Log)

	api := r.Group("/api/v1")
	setupGroup := api.Group("/setup")
	setupGroup.POST("/init", setup.Init)
}
