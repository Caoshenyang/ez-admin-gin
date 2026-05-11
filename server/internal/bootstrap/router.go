// Package bootstrap 负责服务启动阶段的路由注册、中间件装配和 Swagger 初始化。
package bootstrap

import (
	authModule "ez-admin-gin/server/internal/modules/auth"
	iamModule "ez-admin-gin/server/internal/modules/iam"
	setupModule "ez-admin-gin/server/internal/modules/setup"
	systemModule "ez-admin-gin/server/internal/modules/system"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	appLogger "ez-admin-gin/server/internal/platform/logger"
	platformMiddleware "ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouterOptions 汇总模块路由装配需要的依赖。
type RouterOptions struct {
	Config     *platformConfig.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      *goredis.Client
	Token      *authnPlatform.Manager
	Session    authnPlatform.SessionStore
	Permission *authzPlatform.Enforcer
}

// NewRouter 创建 Gin 引擎，并按模块聚合路由。
func NewRouter(opts RouterOptions) *gin.Engine {
	r := gin.New()
	r.Use(
		platformMiddleware.CORS(opts.Config.CORS, opts.Config.App.Env),
		platformMiddleware.SecurityHeaders(),
		platformMiddleware.RequestID(),
		appLogger.GinLogger(opts.Log),
		appLogger.GinRecovery(opts.Log),
	)

	if opts.Config.Upload.MaxSizeMB > 0 {
		r.MaxMultipartMemory = opts.Config.Upload.MaxSizeMB << 20
	}
	r.Static(opts.Config.Upload.PublicPath, opts.Config.Upload.Dir)

	if opts.Config.Swagger.Enabled {
		RegisterSwagger(r)
	}

	authModule.RegisterRoutes(r, authModule.RouteOptions{
		Config: opts.Config,
		Log:    opts.Log,
		DB:     opts.DB,
		Redis:  opts.Redis,
		Token:  opts.Token,
		Session: opts.Session,
	})
	setupModule.RegisterRoutes(r, setupModule.RouteOptions{
		Log: opts.Log,
		DB:  opts.DB,
	})
	iamModule.RegisterRoutes(r, iamModule.RouteOptions{
		Log:        opts.Log,
		DB:         opts.DB,
		Token:      opts.Token,
		Permission: opts.Permission,
	})
	systemModule.RegisterRoutes(r, systemModule.RouteOptions{
		Config:     opts.Config,
		Log:        opts.Log,
		DB:         opts.DB,
		Redis:      opts.Redis,
		Token:      opts.Token,
		Permission: opts.Permission,
	})

	return r
}
