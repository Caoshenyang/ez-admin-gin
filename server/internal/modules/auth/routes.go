package auth

import (
	authapi "ez-admin-gin/server/internal/modules/auth/api"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	platformConfig "ez-admin-gin/server/internal/platform/config"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	Config *platformConfig.Config
	Log    *zap.Logger
	DB     *gorm.DB
	Redis  *goredis.Client
	Token  *authnPlatform.Manager
}

func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	authapi.RegisterRoutes(r, authapi.RouteOptions{
		Config: opts.Config,
		Log:    opts.Log,
		DB:     opts.DB,
		Redis:  opts.Redis,
		Token:  opts.Token,
	})
}
