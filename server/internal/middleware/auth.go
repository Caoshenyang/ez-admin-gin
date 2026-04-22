package middleware

import (
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/response"
	"ez-admin-gin/server/internal/token"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	currentUserIDKey   = "current_user_id"
	currentUsernameKey = "current_username"
)

// Auth 校验 Authorization 请求头，并把当前用户信息写入 Gin 上下文。
func Auth(tokenManager *token.Manager, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			response.Error(c, apperror.Unauthorized("请先登录"), log)
			c.Abort()
			return
		}

		claims, err := tokenManager.ParseAccessToken(tokenString)
		if err != nil {
			if log != nil {
				log.Warn("parse access token failed", zap.Error(err))
			}

			response.Error(c, apperror.Unauthorized("登录已过期，请重新登录"), log)
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
