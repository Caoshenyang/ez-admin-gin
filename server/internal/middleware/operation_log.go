package middleware

import (
	"net/http"
	"strings"
	"time"

	"ez-admin-gin/server/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const maxOperationLogTextLength = 500

// OperationLog 在请求结束后记录后台写操作。
func OperationLog(db *gorm.DB, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		if shouldSkipOperationLog(c) {
			return
		}

		userID, _ := CurrentUserID(c)
		username, _ := CurrentUsername(c)
		statusCode := c.Writer.Status()

		record := model.OperationLog{
			UserID:       userID,
			Username:     username,
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			RoutePath:    c.FullPath(),
			Query:        truncateOperationLogText(c.Request.URL.RawQuery, 1000),
			IP:           c.ClientIP(),
			UserAgent:    truncateOperationLogText(c.Request.UserAgent(), maxOperationLogTextLength),
			StatusCode:   statusCode,
			LatencyMs:    time.Since(start).Milliseconds(),
			Success:      statusCode < http.StatusBadRequest,
			ErrorMessage: operationErrorMessage(c, statusCode),
		}

		if err := db.Create(&record).Error; err != nil && log != nil {
			log.Warn("create operation log failed", zap.Error(err))
		}
	}
}

func shouldSkipOperationLog(c *gin.Context) bool {
	method := c.Request.Method
	if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
		return true
	}

	// 静态资源和未匹配到路由的请求不作为后台操作记录。
	if c.FullPath() == "" {
		return true
	}

	return false
}

func operationErrorMessage(c *gin.Context, statusCode int) string {
	if len(c.Errors) > 0 {
		return truncateOperationLogText(c.Errors.Last().Error(), maxOperationLogTextLength)
	}

	if statusCode >= http.StatusBadRequest {
		return http.StatusText(statusCode)
	}

	return ""
}

func truncateOperationLogText(value string, maxLength int) string {
	value = strings.TrimSpace(value)
	if len(value) <= maxLength {
		return value
	}

	return value[:maxLength]
}
