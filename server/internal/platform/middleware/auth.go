package middleware

import (
	"context"
	"strings"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	currentUserIDKey   = "current_user_id"
	currentUsernameKey = "current_username"
)

// TokenBlacklistChecker checks if an access token has been revoked.
type TokenBlacklistChecker interface {
	IsBlacklisted(ctx context.Context, tokenString string) bool
}

// Auth 校验 Authorization 请求头，并把当前用户信息写入 Gin 上下文。
func Auth(tokenManager *authnPlatform.Manager, blacklist TokenBlacklistChecker, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			httpx.Error(c, errorsx.Unauthorized("请先登录"), log)
			c.Abort()
			return
		}

		claims, err := tokenManager.ParseAccessToken(tokenString)
		if err != nil {
			if log != nil {
				log.Warn("parse access token failed", zap.Error(err))
			}

			httpx.Error(c, errorsx.Unauthorized("登录已过期，请重新登录"), log)
			c.Abort()
			return
		}

		if blacklist != nil && blacklist.IsBlacklisted(c.Request.Context(), tokenString) {
			httpx.Error(c, errorsx.Unauthorized("登录已过期，请重新登录"), log)
			c.Abort()
			return
		}

		// 后续 Handler 可以从 Gin 上下文中取当前用户信息。
		c.Set(currentUserIDKey, claims.UserID)
		c.Set(currentUsernameKey, claims.Username)
		c.Next()
	}
}

// CurrentUserID 从 Gin 上下文中取当前用户 ID。
func CurrentUserID(c *gin.Context) (uint, bool) {
	value, ok := c.Get(currentUserIDKey)
	if !ok {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}

// CurrentUsername 从 Gin 上下文中取当前用户名。
func CurrentUsername(c *gin.Context) (string, bool) {
	value, ok := c.Get(currentUsernameKey)
	if !ok {
		return "", false
	}

	username, ok := value.(string)
	return username, ok
}

// RequireUserID 从 Gin 上下文中取当前用户 ID，未登录时自动写入 401 响应。
func RequireUserID(c *gin.Context, log *zap.Logger) (uint, bool) {
	userID, ok := CurrentUserID(c)
	if !ok {
		httpx.Error(c, errorsx.Unauthorized("请先登录"), log)
		return 0, false
	}
	return userID, true
}

// Username 从 Gin 上下文中取当前用户名，未登录时返回空字符串。
func Username(c *gin.Context) string {
	username, _ := CurrentUsername(c)
	return username
}

// bearerToken 解析 Authorization: Bearer <token>。
func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 {
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	if strings.TrimSpace(parts[1]) == "" {
		return "", false
	}

	return parts[1], true
}
