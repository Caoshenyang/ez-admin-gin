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
