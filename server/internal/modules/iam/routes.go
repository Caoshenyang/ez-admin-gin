package iam

import (
	departmentmodule "ez-admin-gin/server/internal/modules/iam/department"
	menumodule "ez-admin-gin/server/internal/modules/iam/menu"
	postmodule "ez-admin-gin/server/internal/modules/iam/post"
	rolemodule "ez-admin-gin/server/internal/modules/iam/role"
	usermodule "ez-admin-gin/server/internal/modules/iam/user"
	"ez-admin-gin/server/internal/modules/modulekit"
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
	Blacklist  middleware.TokenBlacklistChecker
}

func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	system := modulekit.NewProtectedSystemGroup(r, modulekit.ProtectedSystemGroupOptions{
		Log:        opts.Log,
		DB:         opts.DB,
		Token:      opts.Token,
		Permission: opts.Permission,
		Blacklist:  opts.Blacklist,
	})

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
