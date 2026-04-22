package system

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// HealthHandler 负责系统健康检查。
type HealthHandler struct {
	cfg         *config.Config
	db          *gorm.DB
	redisClient *goredis.Client
	log         *zap.Logger
}

// NewHealthHandler 创建健康检查 Handler。
func NewHealthHandler(
	cfg *config.Config,
	db *gorm.DB,
	redisClient *goredis.Client,
	log *zap.Logger,
) *HealthHandler {
	return &HealthHandler{
		cfg:         cfg,
		db:          db,
		redisClient: redisClient,
		log:         log,
	}
}

// Check 返回当前服务、数据库和 Redis 的健康状态。
func (h *HealthHandler) Check(c *gin.Context) {
	if err := database.Ping(h.db); err != nil {
		h.log.Error("database health check failed", zap.Error(err))
		response.Error(c, apperror.ServiceUnavailable("数据库不可用", err), h.log)
		return
	}

	if err := appRedis.Ping(h.redisClient); err != nil {
		h.log.Error("redis health check failed", zap.Error(err))
		response.Error(c, apperror.ServiceUnavailable("Redis 不可用", err), h.log)
		return
	}

	response.Success(c, gin.H{
		"env":      h.cfg.App.Env,
		"database": "ok",
		"redis":    "ok",
	})
}
