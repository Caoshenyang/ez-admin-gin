package bootstrap

import (
	"ez-admin-gin/server/internal/config"
	authModule "ez-admin-gin/server/internal/module/auth"
	setupModule "ez-admin-gin/server/internal/module/setup"
	systemModule "ez-admin-gin/server/internal/module/system"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"
	appLogger "ez-admin-gin/server/internal/platform/logger"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouterOptions 汇总模块路由装配需要的依赖。
type RouterOptions struct {
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      *goredis.Client
	Token      *authnPlatform.Manager
	Permission *authzPlatform.Enforcer
}

// NewRouter 创建 Gin 引擎，并按模块聚合路由。
func NewRouter(opts RouterOptions) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	if opts.Config.Upload.MaxSizeMB > 0 {
		r.MaxMultipartMemory = opts.Config.Upload.MaxSizeMB << 20
	}
	r.Static(opts.Config.Upload.PublicPath, opts.Config.Upload.Dir)

	authModule.RegisterRoutes(r, authModule.RouteOptions{
		Config: opts.Config,
		Log:    opts.Log,
		DB:     opts.DB,
		Redis:  opts.Redis,
		Token:  opts.Token,
	})
	setupModule.RegisterRoutes(r, setupModule.RouteOptions{
		Log: opts.Log,
		DB:  opts.DB,
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
