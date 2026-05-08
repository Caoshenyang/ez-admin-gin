package modulekit

import (
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"
	"ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProtectedSystemGroupOptions struct {
	Log        *zap.Logger
	DB         *gorm.DB
	Token      *authnPlatform.Manager
	Permission *authzPlatform.Enforcer
}

func NewProtectedSystemGroup(r *gin.Engine, opts ProtectedSystemGroupOptions) *gin.RouterGroup {
	api := r.Group("/api/v1")
	system := api.Group("/system")
	// 这里收口后台受保护分组的公共中间件，避免 IAM/System 顶层路由重复拷贝一份。
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.LoadActor(opts.DB, opts.Log))
	system.Use(middleware.OperationLog(opts.DB, opts.Log))
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))
	return system
}
