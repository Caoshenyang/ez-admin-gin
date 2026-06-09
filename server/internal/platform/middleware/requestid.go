package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

const headerKey = "X-Request-ID"

// RequestID 为每个请求生成或传递唯一 ID，注入 context 和响应头。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(headerKey)
		if id == "" {
			id = uuid.New().String()
		}

		c.Set(string(requestIDKey), id)
		c.Header(headerKey, id)

		ctx := context.WithValue(c.Request.Context(), requestIDKey, id)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetRequestID 从 context 中取出 request ID，供日志等下游使用。
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
