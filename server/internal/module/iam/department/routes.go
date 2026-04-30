package department

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/departments", handler.List)
	group.POST("/departments", handler.Create)
	group.POST("/departments/:id/update", handler.Update)
	group.POST("/departments/:id/status", handler.UpdateStatus)
}
