package notification

import (
	notiapi "ez-admin-gin/server/internal/modules/system/notification/api"
	notiapp "ez-admin-gin/server/internal/modules/system/notification/application"

	authnPlatform "ez-admin-gin/server/internal/platform/authn"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteOptions struct {
	DB    *gorm.DB
	Redis *redis.Client
	Log   *zap.Logger
	Token *authnPlatform.Manager
}

func RegisterRoutes(group *gin.RouterGroup, engine *gin.Engine, opts RouteOptions) {
	service := NewService(ServiceOptions{DB: opts.DB, Redis: opts.Redis, Log: opts.Log})
	hub := NewHub(service, opts.Redis, opts.Log)

	notiapi.RegisterRoutes(group, notiapi.RouteOptions{
		Service: service,
		Hub:     hub,
		Log:     opts.Log,
	})

	// WebSocket 端点需要独立注册，因为它不走 Auth 中间件链
	registerWSRoute(engine, service, hub, opts.Token, opts.Log)
}

// registerWSRoute 在 system group 外部注册 WebSocket 端点。
// 因为浏览器 WebSocket API 不支持自定义 Header，WS 需要从 query param 解析 token。
func registerWSRoute(engine *gin.Engine, service *notiapp.Service, hub interface{ Run() }, token *authnPlatform.Manager, log *zap.Logger) {
	// 直接在 Engine 上注册，绕过 Auth / Permission 中间件链
	// WebSocket handler 内部自行从 query param 解析 token 认证
	wsHandler := notiapi.NewWSHandler(service, hub, token, log)
	engine.GET("/api/v1/system/notifications/ws", wsHandler.ServeWebSocket)
}
