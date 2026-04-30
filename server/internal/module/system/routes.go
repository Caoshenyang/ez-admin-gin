package system

import (
	"ez-admin-gin/server/internal/config"
	legacySystemHandler "ez-admin-gin/server/internal/handler/system"
	"ez-admin-gin/server/internal/middleware"
	iamDepartmentModule "ez-admin-gin/server/internal/module/iam/department"
	iamMenuModule "ez-admin-gin/server/internal/module/iam/menu"
	iamPostModule "ez-admin-gin/server/internal/module/iam/post"
	iamRoleModule "ez-admin-gin/server/internal/module/iam/role"
	iamUserModule "ez-admin-gin/server/internal/module/iam/user"
	systemConfigModule "ez-admin-gin/server/internal/module/system/config"
	systemFileModule "ez-admin-gin/server/internal/module/system/file"
	systemLoginLogModule "ez-admin-gin/server/internal/module/system/loginlog"
	systemNoticeModule "ez-admin-gin/server/internal/module/system/notice"
	systemOperationLogModule "ez-admin-gin/server/internal/module/system/operationlog"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总系统模块路由依赖。
type RouteOptions struct {
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      *goredis.Client
	Token      *authnPlatform.Manager
	Permission *authzPlatform.Enforcer
}

// RegisterRoutes 注册系统模块路由。
func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	health := legacySystemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)

	r.GET("/health", health.Check)

	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.LoadActor(opts.DB, opts.Log))
	system.Use(middleware.OperationLog(opts.DB, opts.Log))
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))

	system.GET("/health", health.Check)
	iamUserModule.RegisterRoutes(system, iamUserModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	iamRoleModule.RegisterRoutes(system, iamRoleModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	iamDepartmentModule.RegisterRoutes(system, iamDepartmentModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	iamPostModule.RegisterRoutes(system, iamPostModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	iamMenuModule.RegisterRoutes(system, iamMenuModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	systemConfigModule.RegisterRoutes(system, systemConfigModule.RouteOptions{
		DB:    opts.DB,
		Redis: opts.Redis,
		Log:   opts.Log,
	})
	systemFileModule.RegisterRoutes(system, systemFileModule.RouteOptions{
		DB:     opts.DB,
		Upload: opts.Config.Upload,
		Log:    opts.Log,
	})
	systemOperationLogModule.RegisterRoutes(system, systemOperationLogModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	systemLoginLogModule.RegisterRoutes(system, systemLoginLogModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	systemNoticeModule.RegisterRoutes(system, systemNoticeModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
