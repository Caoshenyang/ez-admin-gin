package config

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/platform/model"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	RuntimeKeyRateLimitLoginMaxRequests      = "rate_limit:login_max_requests"
	RuntimeKeyRateLimitLoginWindowSec        = "rate_limit:login_window_sec"
	RuntimeKeyRateLimitLoginLockoutThreshold = "rate_limit:login_lockout_threshold"
	RuntimeKeyRateLimitLoginLockoutSec       = "rate_limit:login_lockout_sec"
	RuntimeKeyUploadMaxSizeMB                = "upload:max_size_mb"
	RuntimeKeyUploadAllowedExts              = "upload:allowed_exts"

	runtimeConfigCachePrefix = "sys_config:"
	runtimeConfigCacheTTL    = time.Hour
)

// RuntimeStore reads the small whitelist of business settings that may be
// managed from sys_config at runtime. Startup and secret settings stay in
// config.yaml / EZ_* env vars.
type RuntimeStore struct {
	db       *gorm.DB
	redis    *goredis.Client
	fallback *Config
	log      *zap.Logger
}

func NewRuntimeStore(db *gorm.DB, redis *goredis.Client, fallback *Config, log *zap.Logger) *RuntimeStore {
	return &RuntimeStore{db: db, redis: redis, fallback: fallback, log: log}
}

func (s *RuntimeStore) RateLimitConfig(ctx context.Context) RateLimitConfig {
	cfg := RateLimitConfig{
		LoginMaxRequests:      10,
		LoginWindowSec:        60,
		LoginLockoutThreshold: 5,
		LoginLockoutSec:       300,
	}
	if s != nil && s.fallback != nil {
		cfg = s.fallback.RateLimit
	}

	if s == nil {
		return cfg
	}
	cfg.LoginMaxRequests = s.intValue(ctx, RuntimeKeyRateLimitLoginMaxRequests, cfg.LoginMaxRequests, 1, 10000)
	cfg.LoginWindowSec = s.intValue(ctx, RuntimeKeyRateLimitLoginWindowSec, cfg.LoginWindowSec, 1, 86400)
	cfg.LoginLockoutThreshold = s.intValue(ctx, RuntimeKeyRateLimitLoginLockoutThreshold, cfg.LoginLockoutThreshold, 1, 100)
	cfg.LoginLockoutSec = s.intValue(ctx, RuntimeKeyRateLimitLoginLockoutSec, cfg.LoginLockoutSec, 1, 86400)
	return cfg
}

func (s *RuntimeStore) UploadConfig(ctx context.Context) UploadConfig {
	cfg := UploadConfig{
		Dir:        "uploads",
		PublicPath: "/uploads",
		MaxSizeMB:  10,
		AllowedExts: []string{
			".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".txt", ".docx", ".xlsx",
		},
	}
	if s != nil && s.fallback != nil {
		cfg = s.fallback.Upload
	}

	if s == nil {
		return cfg
	}
	cfg.MaxSizeMB = s.int64Value(ctx, RuntimeKeyUploadMaxSizeMB, cfg.MaxSizeMB, 1, 50)
	cfg.AllowedExts = s.stringListValue(ctx, RuntimeKeyUploadAllowedExts, cfg.AllowedExts)
	return cfg
}

func (s *RuntimeStore) intValue(ctx context.Context, key string, fallback int, min int, max int) int {
	raw, ok := s.value(ctx, key)
	if !ok {
		return fallback
	}

	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value < min || value > max {
		s.warnInvalid(key, raw, "integer out of allowed range")
		return fallback
	}
	return value
}

func (s *RuntimeStore) int64Value(ctx context.Context, key string, fallback int64, min int64, max int64) int64 {
	raw, ok := s.value(ctx, key)
	if !ok {
		return fallback
	}

	value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || value < min || value > max {
		s.warnInvalid(key, raw, "integer out of allowed range")
		return fallback
	}
	return value
}

func (s *RuntimeStore) stringListValue(ctx context.Context, key string, fallback []string) []string {
	raw, ok := s.value(ctx, key)
	if !ok {
		return cloneStrings(fallback)
	}

	values := splitConfigList(raw)
	if len(values) == 0 {
		s.warnInvalid(key, raw, "empty list")
		return cloneStrings(fallback)
	}
	return values
}

func (s *RuntimeStore) value(ctx context.Context, key string) (string, bool) {
	if s == nil || strings.TrimSpace(key) == "" {
		return "", false
	}

	cacheKey := runtimeConfigCachePrefix + key
	if s.redis != nil {
		value, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			return value, true
		}
		if err != nil && !errors.Is(err, goredis.Nil) {
			s.warn("read runtime config cache failed", key, err)
		}
	}

	if s.db == nil {
		return "", false
	}

	var item model.SystemConfig
	err := s.db.
		Where("config_key = ?", key).
		Where("status = ?", model.SystemConfigStatusEnabled).
		First(&item).
		Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.warn("read runtime config from database failed", key, err)
		}
		return "", false
	}

	if s.redis != nil {
		if err := s.redis.Set(ctx, cacheKey, item.Value, runtimeConfigCacheTTL).Err(); err != nil {
			s.warn("write runtime config cache failed", key, err)
		}
	}
	return item.Value, true
}

func (s *RuntimeStore) warnInvalid(key string, value string, reason string) {
	if s == nil || s.log == nil {
		return
	}
	s.log.Warn(
		"ignore invalid runtime config value",
		zap.String("key", key),
		zap.String("value", value),
		zap.String("reason", reason),
	)
}

func (s *RuntimeStore) warn(message string, key string, err error) {
	if s == nil || s.log == nil {
		return
	}
	s.log.Warn(message, zap.String("key", key), zap.Error(err))
}

func splitConfigList(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || r == '\t' || r == ' '
	})

	values := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		values = append(values, value)
	}
	return values
}

func cloneStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	return append([]string(nil), values...)
}
