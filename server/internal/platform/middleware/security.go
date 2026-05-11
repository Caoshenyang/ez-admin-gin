package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders 添加标准 Web 安全响应头。
//
//	 X-Content-Type-Options  — 阻止浏览器 MIME 嗅探
//	 X-Frame-Options         — 阻止 iframe 嵌入（防点击劫持）
//	 Referrer-Policy         — 控制Referer 泄漏
//	 Permissions-Policy      — 禁用不需要的浏览器 API
//	 Content-Security-Policy — 基础 CSP，限制资源加载来源
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")

		c.Next()
	}
}
