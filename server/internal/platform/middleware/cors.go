package middleware

import (
	"net/http"
	"strings"

	platformConfig "ez-admin-gin/server/internal/platform/config"

	"github.com/gin-gonic/gin"
)

// CORS 返回一个跨域中间件，根据配置决定允许的来源。
func CORS(cfg platformConfig.CORSConfig, env string) gin.HandlerFunc {
	allowAllLocalhost := env != "prod"
	isProd := env == "prod"

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		allowed := false

		if allowAllLocalhost && isLocalhost(origin) {
			allowed = true
		}

		if !allowed {
			for _, o := range cfg.AllowedOrigins {
				if isProd && o == "*" {
					continue
				}
				if o == origin || (!isProd && o == "*") {
					allowed = true
					break
				}
			}
		}

		if !allowed {
			c.Next()
			return
		}

		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isLocalhost(origin string) bool {
	return strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:")
}
