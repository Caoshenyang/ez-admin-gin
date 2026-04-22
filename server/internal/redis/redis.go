package redis

import (
	"context"
	"fmt"
	"time"

	"ez-admin-gin/server/internal/config"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var pingTimeout = 3 * time.Second

// New 创建 Redis 客户端，并在启动时完成连通性检查。
func New(cfg config.RedisConfig, log *zap.Logger) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
	})

	if err := Ping(client); err != nil {
		return nil, err
	}

	log.Info(
		"redis connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	return client, nil
}

// Ping 用于健康检查，确认 Redis 当前仍然可连接。
func Ping(client *goredis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	return nil
}

// Close 关闭 Redis 客户端。
func Close(client *goredis.Client) error {
	if err := client.Close(); err != nil {
		return fmt.Errorf("close redis: %w", err)
	}

	return nil
}
