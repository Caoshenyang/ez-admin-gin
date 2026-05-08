package api

import (
	authservicekit "ez-admin-gin/server/internal/modules/auth/servicekit"
	"ez-admin-gin/server/internal/modules/modulekit"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"

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
