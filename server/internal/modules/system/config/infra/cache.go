// Package infra 实现系统配置的数据访问层。
package infra

import (
	"context"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// Cache 封装系统配置的 Redis 缓存读写操作。
type Cache struct {
	client *goredis.Client
}

func NewCache(client *goredis.Client) *Cache {
	return &Cache{client: client}
}

// Get 从 Redis 读取缓存值；Cache 为 nil 时静默返回未命中。
func (c *Cache) Get(ctx context.Context, key string) (string, bool, error) {
	if c == nil || c.client == nil {
		return "", false, nil
	}

	value, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, goredis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return value, true, nil
}

// Set 将键值对写入 Redis 并设置过期时间。
func (c *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Set(ctx, key, value, ttl).Err()
}

// Delete 从 Redis 中删除指定缓存键。
func (c *Cache) Delete(ctx context.Context, key string) error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Del(ctx, key).Err()
}
