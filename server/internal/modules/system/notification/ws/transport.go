package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const notificationChannel = "ez-admin:notifications"

// TransportMessage 是通过 Redis Pub/Sub 传输的消息。
type TransportMessage struct {
	Type   string `json:"type"`
	UserID uint   `json:"user_id"`
	Data   any    `json:"data"`
}

// Transport 定义消息分发接口。
type Transport interface {
	Publish(ctx context.Context, msg TransportMessage) error
	Subscribe() <-chan TransportMessage
}

// RedisTransport 基于 Redis Pub/Sub 的消息分发实现。
type RedisTransport struct {
	rdb *redis.Client
	log *zap.Logger
}

// NewRedisTransport 创建 Redis 传输层。
func NewRedisTransport(rdb *redis.Client, log *zap.Logger) *RedisTransport {
	return &RedisTransport{rdb: rdb, log: log}
}

// Publish 发布消息到 Redis channel。
func (t *RedisTransport) Publish(ctx context.Context, msg TransportMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal transport message: %w", err)
	}
	return t.rdb.Publish(ctx, notificationChannel, data).Err()
}

// Subscribe 订阅 Redis channel，返回消息 channel。
func (t *RedisTransport) Subscribe() <-chan TransportMessage {
	ch := make(chan TransportMessage, 64)

	go func() {
		defer close(ch)
		sub := t.rdb.Subscribe(context.Background(), notificationChannel)
		defer sub.Close()

		msgCh := sub.Channel()
		for msg := range msgCh {
			var tm TransportMessage
			if err := json.Unmarshal([]byte(msg.Payload), &tm); err != nil {
				t.log.Error("redis transport unmarshal error", zap.Error(err))
				continue
			}
			ch <- tm
		}
	}()

	return ch
}
