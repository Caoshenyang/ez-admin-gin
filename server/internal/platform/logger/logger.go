package logger

import (
	legacyConfig "ez-admin-gin/server/internal/config"
	legacyLogger "ez-admin-gin/server/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// New 创建结构化日志对象。
func New(cfg legacyConfig.LogConfig) (*zap.Logger, error) {
	return legacyLogger.New(cfg)
}

// GinLogger 创建请求日志中间件。
func GinLogger(log *zap.Logger) gin.HandlerFunc {
	return legacyLogger.GinLogger(log)
}

// GinRecovery 创建 panic 恢复中间件。
func GinRecovery(log *zap.Logger) gin.HandlerFunc {
	return legacyLogger.GinRecovery(log)
}
