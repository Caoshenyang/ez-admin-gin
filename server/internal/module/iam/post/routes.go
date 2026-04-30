package post

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

	group.GET("/posts", handler.List)
	group.POST("/posts", handler.Create)
	group.POST("/posts/:id/update", handler.Update)
	group.POST("/posts/:id/status", handler.UpdateStatus)
}
