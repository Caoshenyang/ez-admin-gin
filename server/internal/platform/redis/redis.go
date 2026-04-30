package redis

import (
	legacyConfig "ez-admin-gin/server/internal/config"
	legacyRedis "ez-admin-gin/server/internal/redis"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// New 创建 Redis 客户端。
func New(cfg legacyConfig.RedisConfig, log *zap.Logger) (*goredis.Client, error) {
	return legacyRedis.New(cfg, log)
}

// Close 关闭 Redis 客户端。
func Close(client *goredis.Client) error {
	return legacyRedis.Close(client)
}
