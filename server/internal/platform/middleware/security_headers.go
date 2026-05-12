package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adds standard security response headers.
func SecurityHeaders(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-XSS-Protection", "0")
		c.Header("Vary", "Origin")

		if env == "prod" {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		c.Next()
	}
}
