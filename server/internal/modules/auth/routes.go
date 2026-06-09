package auth

import (
	authapi "ez-admin-gin/server/internal/modules/auth/api"
	authservicekit "ez-admin-gin/server/internal/modules/auth/servicekit"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	platformConfig "ez-admin-gin/server/internal/platform/config"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	Config        *platformConfig.Config
	RuntimeConfig *platformConfig.RuntimeStore
	Log           *zap.Logger
	DB            *gorm.DB
	Redis         *goredis.Client
	Token         *authnPlatform.Manager
}

func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	var refreshStore *authnPlatform.RefreshTokenStore
	if opts.Token.RefreshStore() != nil {
		refreshStore = opts.Token.RefreshStore()
	}

	services := authservicekit.NewServices(authservicekit.ServiceOptions{
		Config:        opts.Config,
		RuntimeConfig: opts.RuntimeConfig,
		Log:           opts.Log,
		DB:            opts.DB,
		Redis:         opts.Redis,
		Token:         opts.Token,
		RefreshStore:  refreshStore,
	})
	authapi.RegisterRoutes(r, authapi.RouteOptions{
		Log:           opts.Log,
		DB:            opts.DB,
		Redis:         opts.Redis,
		Token:         opts.Token,
		Services:      services,
		RuntimeConfig: opts.RuntimeConfig,
		Env:           opts.Config.App.Env,
		Blacklist:     refreshStore,
	})
}
