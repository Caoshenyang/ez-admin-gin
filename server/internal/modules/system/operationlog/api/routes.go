package api

import (
	operationlogapp "ez-admin-gin/server/internal/modules/system/operationlog/application"
	operationloginfra "ez-admin-gin/server/internal/modules/system/operationlog/infra"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := operationloginfra.NewRepository(opts.DB)
	service := operationlogapp.NewService(repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/operation-logs", handler.List)
}
