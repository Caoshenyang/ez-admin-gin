package middleware

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	platformConfig "ez-admin-gin/server/internal/platform/config"

	"github.com/gin-gonic/gin"
)

func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{name: "http localhost", origin: "http://localhost:5173", want: true},
		{name: "https localhost", origin: "https://localhost:5173", want: true},
		{name: "ipv4 loopback", origin: "http://127.0.0.1:5173", want: true},
		{name: "ipv6 loopback", origin: "http://[::1]:5173", want: true},
		{name: "remote origin", origin: "https://admin.example.com", want: false},
		{name: "unsupported scheme", origin: "ftp://localhost:21", want: false},
		{name: "missing scheme", origin: "localhost:5173", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLocalhost(tt.origin); got != tt.want {
				t.Fatalf("isLocalhost(%q) = %v, want %v", tt.origin, got, tt.want)
			}
		})
	}
}

func TestCORSAllowsHTTPSLocalhostInDev(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const origin = "https://localhost:5173"
	router := gin.New()
	router.Use(CORS(platformConfig.CORSConfig{}, "dev"))
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", origin)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != origin {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, origin)
	}
}

func TestAllowedWebSocketOriginPatterns(t *testing.T) {
	cfg := platformConfig.CORSConfig{
		AllowedOrigins: []string{
			"https://admin.example.com",
			"https://admin.example.com",
			"*",
		},
	}

	dev := AllowedWebSocketOriginPatterns(cfg, "dev")
	for _, pattern := range []string{"localhost:*", "127.0.0.1:*", "[::1]:*", "admin.example.com"} {
		if !slices.Contains(dev, pattern) {
			t.Fatalf("dev patterns = %v, want %q", dev, pattern)
		}
	}

	prod := AllowedWebSocketOriginPatterns(cfg, "prod")
	if slices.Contains(prod, "*") || slices.Contains(prod, "localhost:*") {
		t.Fatalf("prod patterns = %v, want no wildcard or dev localhost defaults", prod)
	}
	if got := len(prod); got != 1 {
		t.Fatalf("len(prod patterns) = %d, want 1: %v", got, prod)
	}
	if prod[0] != "admin.example.com" {
		t.Fatalf("prod[0] = %q, want admin.example.com", prod[0])
	}
}
