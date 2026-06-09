package config

import (
	"strings"
	"testing"
	"time"
)

func TestAuthConfigDurations(t *testing.T) {
	cfg := AuthConfig{
		AccessTokenTTL:  3600,
		RefreshTokenTTL: 7200,
	}

	if got := cfg.AccessTokenDuration(); got != time.Hour {
		t.Fatalf("AccessTokenDuration() = %s, want %s", got, time.Hour)
	}
	if got := cfg.RefreshTokenDuration(); got != 2*time.Hour {
		t.Fatalf("RefreshTokenDuration() = %s, want %s", got, 2*time.Hour)
	}

	cfg.RefreshTokenTTL = 0
	if got := cfg.RefreshTokenDuration(); got != defaultRefreshTokenTTL {
		t.Fatalf("RefreshTokenDuration() fallback = %s, want %s", got, defaultRefreshTokenTTL)
	}
}

func TestValidateProductionRejectsUnsafeDefaults(t *testing.T) {
	cfg := Config{
		App:     AppConfig{Env: "prod"},
		Auth:    AuthConfig{JWTSecret: "ez-admin-dev-secret-change-me-please-32"},
		CORS:    CORSConfig{AllowedOrigins: []string{"https://admin.example.com"}},
		Swagger: SwaggerConfig{Enabled: false},
		Upload:  UploadConfig{MaxSizeMB: 10},
	}

	if err := cfg.ValidateProduction(); err == nil {
		t.Fatal("ValidateProduction() error = nil, want unsafe default secret error")
	}
}

func TestValidateProductionAllowsSafeConfig(t *testing.T) {
	cfg := Config{
		App:     AppConfig{Env: "prod"},
		Auth:    AuthConfig{JWTSecret: strings.Repeat("a", 32)},
		CORS:    CORSConfig{AllowedOrigins: []string{"https://admin.example.com"}},
		Swagger: SwaggerConfig{Enabled: false},
		Upload:  UploadConfig{MaxSizeMB: 50},
	}

	if err := cfg.ValidateProduction(); err != nil {
		t.Fatalf("ValidateProduction() error = %v, want nil", err)
	}
}

func TestRuntimeStoreNilFallsBackToStaticConfig(t *testing.T) {
	static := Config{
		RateLimit: RateLimitConfig{
			LoginMaxRequests:      20,
			LoginWindowSec:        120,
			LoginLockoutThreshold: 6,
			LoginLockoutSec:       600,
		},
		Upload: UploadConfig{
			Dir:         "uploads",
			PublicPath:  "/uploads",
			MaxSizeMB:   12,
			AllowedExts: []string{".jpg", ".pdf"},
		},
	}

	store := NewRuntimeStore(nil, nil, &static, nil)

	rateLimit := store.RateLimitConfig(nil)
	if rateLimit.LoginMaxRequests != 20 || rateLimit.LoginWindowSec != 120 || rateLimit.LoginLockoutThreshold != 6 || rateLimit.LoginLockoutSec != 600 {
		t.Fatalf("RateLimitConfig() = %+v, want static fallback", rateLimit)
	}

	upload := store.UploadConfig(nil)
	if upload.MaxSizeMB != 12 || strings.Join(upload.AllowedExts, ",") != ".jpg,.pdf" {
		t.Fatalf("UploadConfig() = %+v, want static fallback", upload)
	}
}

func TestSplitConfigListTrimsAndDeduplicates(t *testing.T) {
	got := splitConfigList(" .jpg, .png\n.jpg; .pdf\t")
	want := ".jpg,.png,.pdf"
	if strings.Join(got, ",") != want {
		t.Fatalf("splitConfigList() = %q, want %q", strings.Join(got, ","), want)
	}
}
