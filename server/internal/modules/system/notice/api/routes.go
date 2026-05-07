package api

import (
	noticeapp "ez-admin-gin/server/internal/modules/system/notice/application"
	noticeinfra "ez-admin-gin/server/internal/modules/system/notice/infra"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := noticeinfra.NewRepository(opts.DB)
	service := noticeapp.NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/notices", handler.List)
	group.POST("/notices", handler.Create)
	group.POST("/notices/:id/update", handler.Update)
	group.POST("/notices/:id/status", handler.UpdateStatus)
}
