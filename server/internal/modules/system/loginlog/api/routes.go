package api

import (
	loginlogapp "ez-admin-gin/server/internal/modules/system/loginlog/application"
	loginloginfra "ez-admin-gin/server/internal/modules/system/loginlog/infra"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := loginloginfra.NewRepository(opts.DB)
	service := loginlogapp.NewService(repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/login-logs", handler.List)
}
