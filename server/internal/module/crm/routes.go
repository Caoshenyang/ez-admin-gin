package crm

import (
	"ez-admin-gin/server/internal/middleware"
	crmCustomerModule "ez-admin-gin/server/internal/module/crm/customer"
	crmFollowUpModule "ez-admin-gin/server/internal/module/crm/followup"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteOptions 汇总 CRM 模块路由依赖。
type RouteOptions struct {
	DB         *gorm.DB
	Log        *zap.Logger
	Token      *authnPlatform.Manager
	Permission *authzPlatform.Enforcer
}

// RegisterRoutes 注册 CRM 聚合路由。
func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	api := r.Group("/api/v1")
	crm := api.Group("/crm")
	crm.Use(middleware.Auth(opts.Token, opts.Log))
	crm.Use(middleware.LoadActor(opts.DB, opts.Log))
	crm.Use(middleware.OperationLog(opts.DB, opts.Log))
	crm.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))

	crmCustomerModule.RegisterRoutes(crm, crmCustomerModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
	crmFollowUpModule.RegisterRoutes(crm, crmFollowUpModule.RouteOptions{
		DB:  opts.DB,
		Log: opts.Log,
	})
}
