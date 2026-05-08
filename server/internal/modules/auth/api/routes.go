package api

import (
	authservicekit "ez-admin-gin/server/internal/modules/auth/servicekit"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	"ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	Log      *zap.Logger
	DB       *gorm.DB
	Token    *authnPlatform.Manager
	Services authservicekit.Services
}

func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	login := NewLoginHandler(opts.Services.Login, opts.Log)
	me := NewMeHandler(opts.Services.Me, opts.Log)
	account := NewAccountHandler(opts.Services.Account, opts.Log)
	menus := NewMenuHandler(opts.Services.Menu, opts.Log)
	dashboard := NewDashboardHandler(opts.Services.Dashboard, opts.Log)

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
