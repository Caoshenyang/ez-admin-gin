package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 汇总整个服务端会读取的配置段。
type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Upload    UploadConfig    `mapstructure:"upload"`
	Log       LogConfig       `mapstructure:"log"`
	Swagger   SwaggerConfig   `mapstructure:"swagger"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type SwaggerConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type ServerConfig struct {
	Addr string `mapstructure:"addr"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	MaxRetries   int    `mapstructure:"max_retries"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	PoolSize     int    `mapstructure:"pool_size"`
}

type AuthConfig struct {
	JWTSecret       string `mapstructure:"jwt_secret"`
	AccessTokenTTL  int    `mapstructure:"access_token_ttl"`
	RefreshTokenTTL int    `mapstructure:"refresh_token_ttl"`
	Issuer          string `mapstructure:"issuer"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type UploadConfig struct {
	Dir         string   `mapstructure:"dir"`
	PublicPath  string   `mapstructure:"public_path"`
	MaxSizeMB   int64    `mapstructure:"max_size_mb"`
	AllowedExts []string `mapstructure:"allowed_exts"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type RateLimitConfig struct {
	LoginMaxRequests        int `mapstructure:"login_max_requests"`
	LoginWindowSec          int `mapstructure:"login_window_sec"`
	LoginMaxAccountAttempts int `mapstructure:"login_max_account_attempts"`
	LoginAccountLockoutSec  int `mapstructure:"login_account_lockout_sec"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	setDefaults(v)
	bindEnvs(v)

	v.SetEnvPrefix("EZ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// ValidateProduction 对生产环境必须检查的安全项做校验。
func (c *Config) ValidateProduction() error {
	if c.App.Env != "prod" {
		return nil
	}

	// JWT secret: must not contain default values, min 32 chars.
	secret := strings.TrimSpace(c.Auth.JWTSecret)
	if strings.Contains(secret, "change-me") || strings.Contains(secret, "dev-secret") {
		return fmt.Errorf("production: auth.jwt_secret contains a default value; please set EZ_AUTH_JWT_SECRET")
	}
	if len(secret) < 32 {
		return fmt.Errorf("production: auth.jwt_secret must be at least 32 characters")
	}

	// CORS: must have at least one non-wildcard origin.
	if len(c.CORS.AllowedOrigins) == 0 {
		return fmt.Errorf("production: cors.allowed_origins must be configured")
	}
	for _, o := range c.CORS.AllowedOrigins {
		if o == "*" {
			return fmt.Errorf("production: cors.allowed_origins must not contain wildcard '*'")
		}
	}

	// Swagger must be disabled.
	if c.Swagger.Enabled {
		return fmt.Errorf("production: swagger.enabled must be false")
	}

	// Database password must not be default.
	if c.Database.Password == "ez_admin_123456" {
		return fmt.Errorf("production: database.password must be changed from default")
	}

	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "ez-admin")
	v.SetDefault("app.env", "dev")
	v.SetDefault("server.addr", ":8080")
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "ez_admin")
	v.SetDefault("database.password", "ez_admin_123456")
	v.SetDefault("database.name", "ez_admin")
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "console")
	v.SetDefault("log.filename", "logs/app.log")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 7)
	v.SetDefault("log.max_age", 30)
	v.SetDefault("log.compress", false)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 50)
	v.SetDefault("database.conn_max_lifetime", 3600)
	v.SetDefault("redis.max_retries", 3)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("auth.jwt_secret", "ez-admin-dev-secret-change-me-please-32")
	v.SetDefault("auth.access_token_ttl", 7200)
	v.SetDefault("auth.refresh_token_ttl", 604800)
	v.SetDefault("auth.issuer", "ez-admin")
	v.SetDefault("upload.dir", "uploads")
	v.SetDefault("upload.public_path", "/uploads")
	v.SetDefault("upload.max_size_mb", 10)
	v.SetDefault("upload.allowed_exts", []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".txt", ".docx", ".xlsx"})
	v.SetDefault("swagger.enabled", true)
	v.SetDefault("cors.allowed_origins", []string{})
	v.SetDefault("rate_limit.login_max_requests", 10)
	v.SetDefault("rate_limit.login_window_sec", 60)
	v.SetDefault("rate_limit.login_max_account_attempts", 5)
	v.SetDefault("rate_limit.login_account_lockout_sec", 300)
}

func bindEnvs(v *viper.Viper) {
	keys := []string{
		"app.name",
		"app.env",
		"server.addr",
		"database.driver",
		"database.host",
		"database.port",
		"database.user",
		"database.password",
		"database.name",
		"redis.host",
		"redis.port",
		"redis.password",
		"redis.db",
		"log.level",
		"log.format",
		"log.filename",
		"log.max_size",
		"log.max_backups",
		"log.max_age",
		"log.compress",
		"database.max_idle_conns",
		"database.max_open_conns",
		"database.conn_max_lifetime",
		"redis.max_retries",
		"redis.min_idle_conns",
		"redis.pool_size",
		"auth.jwt_secret",
		"auth.access_token_ttl",
		"auth.refresh_token_ttl",
		"auth.issuer",
		"upload.dir",
		"upload.public_path",
		"upload.max_size_mb",
		"upload.allowed_exts",
		"swagger.enabled",
		"cors.allowed_origins",
		"rate_limit.login_max_requests",
		"rate_limit.login_window_sec",
		"rate_limit.login_max_account_attempts",
		"rate_limit.login_account_lockout_sec",
	}

	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}
