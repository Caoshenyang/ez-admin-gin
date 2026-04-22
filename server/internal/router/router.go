package router

import (
	"ez-admin-gin/server/internal/config"
	systemHandler "ez-admin-gin/server/internal/handler/system"
	appLogger "ez-admin-gin/server/internal/logger"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Options 汇总路由层需要依赖的对象。
type Options struct {
	Config *config.Config
	Log    *zap.Logger
	DB     *gorm.DB
	Redis  *goredis.Client
}

// New 创建路由引擎，并统一注册中间件和路由分组。
func New(opts Options) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	registerSystemRoutes(r, opts)

	return r
}

// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.GET("/health", health.Check)
}
