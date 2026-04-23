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

	// 配置上传最大内存
	if opts.Config.Upload.MaxSizeMB > 0 {
		r.MaxMultipartMemory = opts.Config.Upload.MaxSizeMB << 20
	}
	// 配置静态文件服务
	r.Static(opts.Config.Upload.PublicPath, opts.Config.Upload.Dir)

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
	users := systemHandler.NewUserHandler(opts.DB, opts.Log)
	roles := systemHandler.NewRoleHandler(opts.DB, opts.Log)
	menus := systemHandler.NewMenuAdminHandler(opts.DB, opts.Log)
	configs := systemHandler.NewSystemConfigHandler(opts.DB, opts.Redis, opts.Log)
	files := systemHandler.NewFileHandler(opts.DB, opts.Config.Upload, opts.Log)
	operationLogs := systemHandler.NewOperationLogHandler(opts.DB, opts.Log)

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.OperationLog(opts.DB, opts.Log))
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))

	system.GET("/health", health.Check)
	system.GET("/users", users.List)
	system.POST("/users", users.Create)
	system.PUT("/users/:id", users.Update)
	system.PATCH("/users/:id/status", users.UpdateStatus)
	system.PUT("/users/:id/roles", users.UpdateRoles)
	system.GET("/roles", roles.List)
	system.POST("/roles", roles.Create)
	system.PUT("/roles/:id", roles.Update)
	system.PATCH("/roles/:id/status", roles.UpdateStatus)
	system.PUT("/roles/:id/permissions", roles.UpdatePermissions)
	system.PUT("/roles/:id/menus", roles.UpdateMenus)
	system.GET("/menus", menus.Tree)
	system.POST("/menus", menus.Create)
	system.PUT("/menus/:id", menus.Update)
	system.PATCH("/menus/:id/status", menus.UpdateStatus)
	system.DELETE("/menus/:id", menus.Delete)
	system.GET("/configs", configs.List)
	system.POST("/configs", configs.Create)
	system.PUT("/configs/:id", configs.Update)
	system.PATCH("/configs/:id/status", configs.UpdateStatus)
	system.GET("/configs/value/:key", configs.Value)
	system.GET("/files", files.List)
	system.POST("/files", files.Upload)
	system.GET("/operation-logs", operationLogs.List)

}
