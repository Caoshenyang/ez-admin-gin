package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	platformConfig "ez-admin-gin/server/internal/platform/config"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
)

// LoginRateLimit 返回一个基于 Redis 滑动窗口的 Gin 限流中间件。
// 按 IP 限制登录接口的请求频率。
func LoginRateLimit(rdb *goredis.Client, runtimeConfig *platformConfig.RuntimeStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			c.Next()
			return
		}

		cfg := runtimeConfig.RateLimitConfig(c.Request.Context())
		maxRequests := cfg.LoginMaxRequests
		windowSec := cfg.LoginWindowSec
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
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many login attempts, please try again later",
			})
			return
		}

		c.Next()
	}
}

// AccountLockChecker provides account-level login rate limiting via Redis.
type AccountLockChecker struct {
	rdb           *goredis.Client
	runtimeConfig *platformConfig.RuntimeStore
}

// NewAccountLockChecker creates a new AccountLockChecker.
func NewAccountLockChecker(rdb *goredis.Client, runtimeConfig *platformConfig.RuntimeStore) *AccountLockChecker {
	return &AccountLockChecker{rdb: rdb, runtimeConfig: runtimeConfig}
}

// IsLocked checks if the account is currently locked due to too many failed attempts.
func (a *AccountLockChecker) IsLocked(ctx context.Context, username string) bool {
	cfg := a.runtimeConfig.RateLimitConfig(ctx)
	if a.rdb == nil || cfg.LoginLockoutThreshold <= 0 {
		return false
	}
	key := fmt.Sprintf("ratelimit:login_account:%s", username)
	val, err := a.rdb.Get(ctx, key).Int()
	if err != nil {
		return false
	}
	return val >= cfg.LoginLockoutThreshold
}

// RecordFailure increments the failure counter for an account.
func (a *AccountLockChecker) RecordFailure(ctx context.Context, username string) {
	cfg := a.runtimeConfig.RateLimitConfig(ctx)
	if a.rdb == nil || cfg.LoginLockoutThreshold <= 0 {
		return
	}
	key := fmt.Sprintf("ratelimit:login_account:%s", username)
	a.rdb.Incr(ctx, key)
	a.rdb.Expire(ctx, key, time.Duration(cfg.LoginLockoutSec)*time.Second)
}

// ClearAttempts resets the failure counter for an account on successful login.
func (a *AccountLockChecker) ClearAttempts(ctx context.Context, username string) {
	if a.rdb == nil {
		return
	}
	key := fmt.Sprintf("ratelimit:login_account:%s", username)
	a.rdb.Del(ctx, key)
}
