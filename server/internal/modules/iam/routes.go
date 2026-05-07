package iam

import (
	departmentmodule "ez-admin-gin/server/internal/modules/iam/department"
	menumodule "ez-admin-gin/server/internal/modules/iam/menu"
	postmodule "ez-admin-gin/server/internal/modules/iam/post"
	rolemodule "ez-admin-gin/server/internal/modules/iam/role"
	usermodule "ez-admin-gin/server/internal/modules/iam/user"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"
	"ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	Log        *zap.Logger
	DB         *gorm.DB
	Token      *authnPlatform.Manager
	Permission *authzPlatform.Enforcer
}

func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.LoadActor(opts.DB, opts.Log))
	system.Use(middleware.OperationLog(opts.DB, opts.Log))
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))

	usermodule.RegisterRoutes(system, usermodule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	rolemodule.RegisterRoutes(system, rolemodule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	departmentmodule.RegisterRoutes(system, departmentmodule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	postmodule.RegisterRoutes(system, postmodule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	menumodule.RegisterRoutes(system, menumodule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
