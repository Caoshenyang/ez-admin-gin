package router

import (
	"ez-admin-gin/server/internal/config"
	authHandler "ez-admin-gin/server/internal/handler/auth"
	systemHandler "ez-admin-gin/server/internal/handler/system"
	appLogger "ez-admin-gin/server/internal/logger"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/permission"
	"ez-admin-gin/server/internal/token"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Options 汇总路由层需要依赖的对象。
type Options struct {
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      *goredis.Client
	Token      *token.Manager
	Permission *permission.Enforcer
}

// New 创建路由引擎，并统一注册中间件和路由分组。
func New(opts Options) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	registerSystemRoutes(r, opts)
	registerAuthRoutes(r, opts)

	return r
}

// registerAuthRoutes 注册认证相关路由。
func registerAuthRoutes(r *gin.Engine, opts Options) {
	login := authHandler.NewLoginHandler(opts.DB, opts.Log, opts.Token)
	me := authHandler.NewMeHandler(opts.Log)
	menus := authHandler.NewMenuHandler(opts.DB, opts.Log)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login", login.Login)

	protectedAuth := auth.Group("")
	protectedAuth.Use(middleware.Auth(opts.Token, opts.Log))
	protectedAuth.GET("/me", me.Me)
	protectedAuth.GET("/menus", menus.Menus)

}

// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))
	system.GET("/health", health.Check)
}
