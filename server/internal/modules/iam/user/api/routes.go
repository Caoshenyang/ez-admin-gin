package api

import (
	userapp "ez-admin-gin/server/internal/modules/iam/user/application"
	userinfra "ez-admin-gin/server/internal/modules/iam/user/infra"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := userinfra.NewRepository(opts.DB)
	service := userapp.NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/users", handler.List)
	group.POST("/users", handler.Create)
	group.POST("/users/:id/update", handler.Update)
	group.POST("/users/:id/status", handler.UpdateStatus)
	group.POST("/users/:id/roles", handler.UpdateRoles)
}
