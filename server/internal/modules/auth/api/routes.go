package api

import (
	authservicekit "ez-admin-gin/server/internal/modules/auth/servicekit"
	"ez-admin-gin/server/internal/modules/modulekit"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	platformMiddleware "ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	Log                 *zap.Logger
	DB                  *gorm.DB
	Redis               *goredis.Client
	Token               *authnPlatform.Manager
	Session             authnPlatform.SessionStore
	Services            authservicekit.Services
	RateLimitMax        int
	RateLimitSec        int
	RefreshTTLSec       int
	LockoutMaxFailures  int
	LockoutSec          int
}

func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	login := NewLoginHandler(opts.Services.Login, opts.RefreshTTLSec, opts.Redis, opts.LockoutMaxFailures, opts.LockoutSec, opts.Log)
	me := NewMeHandler(opts.Services.Me, opts.Log)
	account := NewAccountHandler(opts.Services.Account, opts.Log)
	menus := NewMenuHandler(opts.Services.Menu, opts.Log)
	dashboard := NewDashboardHandler(opts.Services.Dashboard, opts.Log)

	refreshService := opts.Services.Refresh

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login",
		platformMiddleware.LoginRateLimit(opts.Redis, opts.RateLimitMax, opts.RateLimitSec),
		login.Login,
	)

	auth.POST("/refresh", NewRefreshHandler(refreshService, opts.RefreshTTLSec, opts.Log).Refresh)
	auth.POST("/logout", NewLogoutHandler(refreshService, opts.Log).Logout)

	protectedAuth := modulekit.NewProtectedAuthGroup(auth, modulekit.ProtectedAuthGroupOptions{
		Log:   opts.Log,
		DB:    opts.DB,
		Token: opts.Token,
	})
	protectedAuth.GET("/me", me.Me)
	protectedAuth.GET("/account", account.Profile)
	protectedAuth.POST("/account/profile", account.UpdateProfile)
	protectedAuth.POST("/account/password", account.UpdatePassword)
	protectedAuth.GET("/menus", menus.Menus)
	protectedAuth.GET("/dashboard", dashboard.Dashboard)
}
