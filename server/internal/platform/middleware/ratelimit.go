package middleware

import (
	"fmt"
	"net/http"
	"time"

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
			// Redis 不可用时放行请求，避免影响正常登录。
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
