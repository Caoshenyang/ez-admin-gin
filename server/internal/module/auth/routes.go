package auth

import (
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/middleware"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总认证模块路由依赖。
type RouteOptions struct {
	Config *config.Config
	Log    *zap.Logger
	DB     *gorm.DB
	Redis  *goredis.Client
	Token  *authnPlatform.Manager
}

// RegisterRoutes 注册认证模块路由。
func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	loginService := NewLoginService(repo, opts.Token, opts.Log)
	meService := NewMeService()
	menuService := NewMenuService(repo)
	dashboardService := NewDashboardService(opts.Config, opts.DB, repo, opts.Redis, opts.Log)

	login := NewLoginHandler(loginService, opts.Log)
	me := NewMeHandler(meService, opts.Log)
	menus := NewMenuHandler(menuService, opts.Log)
	dashboard := NewDashboardHandler(dashboardService, opts.Log)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login", login.Login)

	protectedAuth := auth.Group("")
	protectedAuth.Use(middleware.Auth(opts.Token, opts.Log))
	protectedAuth.Use(middleware.LoadActor(opts.DB, opts.Log))
	protectedAuth.GET("/me", me.Me)
	protectedAuth.GET("/menus", menus.Menus)
	protectedAuth.GET("/dashboard", dashboard.Dashboard)
}
