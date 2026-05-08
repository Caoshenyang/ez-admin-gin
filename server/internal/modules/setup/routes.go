package setup

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	Log *zap.Logger
	DB  *gorm.DB
}

func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	service := NewAppService(ServiceOptions{DB: opts.DB})
	h := newSetupHandler(service, opts.Log)

	api := r.Group("/api/v1")
	setupGroup := api.Group("/setup")
	setupGroup.POST("/init", h.Init)
}
