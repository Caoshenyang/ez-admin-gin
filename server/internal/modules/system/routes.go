package system

import (
	"ez-admin-gin/server/internal/modules/modulekit"
	attachmentmodule "ez-admin-gin/server/internal/modules/system/attachment"
	configmodule "ez-admin-gin/server/internal/modules/system/config"
	dictmodule "ez-admin-gin/server/internal/modules/system/dict"
	filemodule "ez-admin-gin/server/internal/modules/system/file"
	loginlogmodule "ez-admin-gin/server/internal/modules/system/loginlog"
	noticemodule "ez-admin-gin/server/internal/modules/system/notice"
	operationlogmodule "ez-admin-gin/server/internal/modules/system/operationlog"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"
	platformConfig "ez-admin-gin/server/internal/platform/config"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	Config     *platformConfig.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      *goredis.Client
	Token      *authnPlatform.Manager
	Permission *authzPlatform.Enforcer
}

func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	health := newHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)

	r.GET("/health", health.Check)

	system := modulekit.NewProtectedSystemGroup(r, modulekit.ProtectedSystemGroupOptions{
		Log:        opts.Log,
		DB:         opts.DB,
		Token:      opts.Token,
		Permission: opts.Permission,
	})

	system.GET("/health", health.Check)
	configmodule.RegisterRoutes(system, configmodule.RouteOptions{
		DB:    opts.DB,
		Redis: opts.Redis,
		Log:   opts.Log,
	})
	dictmodule.RegisterRoutes(system, dictmodule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	attachmentmodule.RegisterRoutes(system, attachmentmodule.RouteOptions{
		DB:     opts.DB,
		Upload: opts.Config.Upload,
		Log:    opts.Log,
	})
	filemodule.RegisterRoutes(system, filemodule.RouteOptions{
		DB:     opts.DB,
		Upload: opts.Config.Upload,
		Log:    opts.Log,
	})
	operationlogmodule.RegisterRoutes(system, operationlogmodule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	loginlogmodule.RegisterRoutes(system, loginlogmodule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	noticemodule.RegisterRoutes(system, noticemodule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
