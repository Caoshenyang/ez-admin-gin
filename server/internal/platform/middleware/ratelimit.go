package middleware

import (
	"context"
	"fmt"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	goredis "github.com/redis/go-redis/v9"
	"github.com/gin-gonic/gin"
)

// LoginRateLimit 返回一个基于 Redis 滑动窗口的 Gin 限流中间件。
// 按 IP 限制登录接口的请求频率。
func LoginRateLimit(rdb *goredis.Client, maxRequests int, windowSec int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:login:%s", ip)

		ctx := c.Request.Context()
		now := time.Now().UnixNano()
		windowStart := now - int64(windowSec)*int64(time.Second)

		pipe := rdb.Pipeline()
		pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", windowStart))
		pipe.ZCard(ctx, key)
		pipe.ZAdd(ctx, key, goredis.Z{Score: float64(now), Member: now})
		pipe.Expire(ctx, key, time.Duration(windowSec)*time.Second)

		results, err := pipe.Exec(ctx)
		if err != nil {
			c.Next()
			return
		}

		count := results[1].(*goredis.IntCmd).Val()
		if count >= int64(maxRequests) {
			c.Header("Retry-After", fmt.Sprintf("%d", windowSec))
			httpx.Error(c, errorsx.TooManyRequests("登录请求过于频繁，请稍后再试"), nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

const lockoutFailKeyPrefix = "lockout:fail:"
const lockoutBlockedKeyPrefix = "lockout:blocked:"

// IsUsernameLocked 检查用户名是否因连续登录失败被临时锁定。
// 当 maxFailures <= 0 时（未配置锁定阈值），始终返回 false。
func IsUsernameLocked(ctx context.Context, rdb *goredis.Client, username string) bool {
	if rdb == nil || username == "" {
		return false
	}
	val, err := rdb.Exists(ctx, lockoutBlockedKeyPrefix+username).Result()
	if err != nil {
		return false
	}
	return val > 0
}

// RecordLoginFailure 记录一次登录失败。连续失败达到 maxFailures 次后，锁定用户名 lockoutSec 秒。
// 当 maxFailures <= 0 时（未配置锁定阈值），不做任何操作。
func RecordLoginFailure(ctx context.Context, rdb *goredis.Client, username string, maxFailures int, lockoutSec int) {
	if rdb == nil || username == "" || maxFailures <= 0 {
		return
	}

	failKey := lockoutFailKeyPrefix + username
	blockedKey := lockoutBlockedKeyPrefix + username

	count, err := rdb.Incr(ctx, failKey).Result()
	if err != nil {
		return
	}

	// Set/refresh TTL on the failure counter.
	rdb.Expire(ctx, failKey, time.Duration(lockoutSec*2)*time.Second)

	// Once the threshold is reached, set the block key and reset the counter.
	if count >= int64(maxFailures) {
		rdb.Set(ctx, blockedKey, "1", time.Duration(lockoutSec)*time.Second)
		rdb.Del(ctx, failKey)
	}
}

// ClearLoginFailures 清除用户名的失败计数（登录成功时调用）。
func ClearLoginFailures(ctx context.Context, rdb *goredis.Client, username string) {
	if rdb == nil || username == "" {
		return
	}
	rdb.Del(ctx, lockoutFailKeyPrefix+username)
}
