package api

import (
	departmentapp "ez-admin-gin/server/internal/modules/iam/department/application"
	departmentinfra "ez-admin-gin/server/internal/modules/iam/department/infra"
	"ez-admin-gin/server/internal/platform/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := departmentinfra.NewRepository(opts.DB)
	tx := database.NewTransactor(opts.DB)
	service := departmentapp.NewService(tx, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/departments", handler.List)
	group.POST("/departments", handler.Create)
	group.POST("/departments/:id/update", handler.Update)
	group.POST("/departments/:id/status", handler.UpdateStatus)
}
