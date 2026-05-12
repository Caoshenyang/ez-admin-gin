package modulekit

import (
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	"ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProtectedAuthGroupOptions struct {
	Log        *zap.Logger
	DB         *gorm.DB
	Token      *authnPlatform.Manager
	Blacklist  middleware.TokenBlacklistChecker
}

func NewProtectedAuthGroup(parent *gin.RouterGroup, opts ProtectedAuthGroupOptions) *gin.RouterGroup {
	protected := parent.Group("")
	// auth 模块的受保护分组只需要认证和当前登录人装配，不挂权限和操作日志。
	protected.Use(middleware.Auth(opts.Token, opts.Blacklist, opts.Log))
	protected.Use(middleware.LoadActor(opts.DB, opts.Log))
	return protected
}
