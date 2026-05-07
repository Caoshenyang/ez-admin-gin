package api

import (
	authapp "ez-admin-gin/server/internal/modules/auth/application"
	authinfra "ez-admin-gin/server/internal/modules/auth/infra"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	"ez-admin-gin/server/internal/platform/middleware"

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
	repo := authinfra.NewRepository(opts.DB)

	login := NewLoginHandler(authapp.NewLoginService(repo, opts.Token, opts.Log), opts.Log)
	me := NewMeHandler(authapp.NewMeService(), opts.Log)
	account := NewAccountHandler(authapp.NewAccountService(opts.DB, repo), opts.Log)
	menus := NewMenuHandler(authapp.NewMenuService(repo), opts.Log)
	dashboard := NewDashboardHandler(
		authapp.NewDashboardService(opts.Config, opts.DB, repo, opts.Redis, opts.Log),
		opts.Log,
	)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login", login.Login)

	protectedAuth := auth.Group("")
	protectedAuth.Use(middleware.Auth(opts.Token, opts.Log))
	protectedAuth.Use(middleware.LoadActor(opts.DB, opts.Log))
	protectedAuth.GET("/me", me.Me)
	protectedAuth.GET("/account", account.Profile)
	protectedAuth.POST("/account/profile", account.UpdateProfile)
	protectedAuth.POST("/account/password", account.UpdatePassword)
	protectedAuth.GET("/menus", menus.Menus)
	protectedAuth.GET("/dashboard", dashboard.Dashboard)
}
