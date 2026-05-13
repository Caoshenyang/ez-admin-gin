package role

import (
	roleapi "ez-admin-gin/server/internal/modules/iam/role/api"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB       *gorm.DB
	Log      *zap.Logger
	Enforcer *authzPlatform.Enforcer
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	service := NewService(ServiceOptions{DB: opts.DB, Enforcer: opts.Enforcer})
	roleapi.RegisterRoutes(group, roleapi.RouteOptions{
		Service: service,
		Log:     opts.Log,
	})
}
