package middleware

import (
	"net/http"
	"net/url"
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

// AllowedWebSocketOriginPatterns 从 CORS 配置中提取 WebSocket OriginPatterns。
// dev 环境默认放行 localhost / 127.0.0.1 的任意端口，prod 只允许显式配置的前端域名。
func AllowedWebSocketOriginPatterns(cfg platformConfig.CORSConfig, env string) []string {
	patterns := make([]string, 0, len(cfg.AllowedOrigins)+2)
	seen := make(map[string]struct{})

	add := func(pattern string) {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			return
		}
		if _, ok := seen[pattern]; ok {
			return
		}
		seen[pattern] = struct{}{}
		patterns = append(patterns, pattern)
	}

	if env != "prod" {
		add("localhost:*")
		add("127.0.0.1:*")
	}

	for _, origin := range cfg.AllowedOrigins {
		host := originHost(origin)
		if host == "" || host == "*" {
			continue
		}
		add(host)
	}

	return patterns
}

func originHost(origin string) string {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return ""
	}

	parsed, err := url.Parse(origin)
	if err == nil && parsed.Host != "" {
		return parsed.Host
	}

	return strings.TrimPrefix(origin, "http://")
}
