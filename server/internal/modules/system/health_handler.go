package system

import (
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	platformDatabase "ez-admin-gin/server/internal/platform/database"
	platformRedis "ez-admin-gin/server/internal/platform/redis"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type healthHandler struct {
	cfg         *platformConfig.Config
	db          *gorm.DB
	redisClient *goredis.Client
	log         *zap.Logger
}

func newHealthHandler(
	cfg *platformConfig.Config,
	db *gorm.DB,
	redisClient *goredis.Client,
	log *zap.Logger,
) *healthHandler {
	return &healthHandler{
		cfg:         cfg,
		db:          db,
		redisClient: redisClient,
		log:         log,
	}
}

// Liveness always returns 200 if the process is up. No dependency checks.
func (h *healthHandler) Liveness(c *gin.Context) {
	httpx.Success(c, gin.H{"status": "ok"})
}

// Readiness checks DB and Redis connectivity. Returns 503 if any dependency is down.
func (h *healthHandler) Readiness(c *gin.Context) {
	if err := platformDatabase.Ping(h.db); err != nil {
		h.log.Error("database health check failed", zap.Error(err))
		httpx.Error(c, errorsx.ServiceUnavailable("数据库不可用", err), h.log)
		return
	}

	if err := platformRedis.Ping(h.redisClient); err != nil {
		h.log.Error("redis health check failed", zap.Error(err))
		httpx.Error(c, errorsx.ServiceUnavailable("Redis 不可用", err), h.log)
		return
	}

	httpx.Success(c, gin.H{
		"status":   "ok",
		"database": "ok",
		"redis":    "ok",
	})
}

// Check is the legacy /health handler (equivalent to Readiness for backward compat).
func (h *healthHandler) Check(c *gin.Context) {
	h.Readiness(c)
}
