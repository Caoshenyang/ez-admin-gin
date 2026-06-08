package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const defaultRefreshTokenTTL = 7 * 24 * time.Hour

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

// SwaggerConfig 控制 Swagger UI 是否对外注册。
type SwaggerConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// AppConfig 描述应用名称和运行环境。
type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

// ServerConfig 描述 HTTP 服务监听地址。
type ServerConfig struct {
	Addr string `mapstructure:"addr"`
}

// DatabaseConfig 描述数据库连接与连接池配置。
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

// RedisConfig 描述 Redis 连接与连接池配置。
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	MaxRetries   int    `mapstructure:"max_retries"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	PoolSize     int    `mapstructure:"pool_size"`
}

// AuthConfig 描述 JWT 和刷新令牌配置。
type AuthConfig struct {
	JWTSecret       string `mapstructure:"jwt_secret"`
	AccessTokenTTL  int    `mapstructure:"access_token_ttl"`
	RefreshTokenTTL int    `mapstructure:"refresh_token_ttl"`
	Issuer          string `mapstructure:"issuer"`
}

// LogConfig 描述 zap 与日志轮转配置。
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// UploadConfig 描述本地上传目录、访问路径和文件限制。
type UploadConfig struct {
	Dir         string   `mapstructure:"dir"`
	PublicPath  string   `mapstructure:"public_path"`
	MaxSizeMB   int64    `mapstructure:"max_size_mb"`
	AllowedExts []string `mapstructure:"allowed_exts"`
}

// CORSConfig 描述允许跨域访问的前端来源。
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// RateLimitConfig 描述登录接口的 IP 限流和账号锁定策略。
type RateLimitConfig struct {
	LoginMaxRequests      int `mapstructure:"login_max_requests"`
	LoginWindowSec        int `mapstructure:"login_window_sec"`
	LoginLockoutThreshold int `mapstructure:"login_lockout_threshold"`
	LoginLockoutSec       int `mapstructure:"login_lockout_sec"`
}

// Load 从 configs/config.yaml 和 EZ_* 环境变量加载服务端配置。
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

	if strings.Contains(c.Auth.JWTSecret, "change-me") || strings.Contains(c.Auth.JWTSecret, "dev-secret") {
		return fmt.Errorf("production: auth.jwt_secret contains a default value; set EZ_AUTH_JWT_SECRET to a secure random string")
	}

	if len(c.CORS.AllowedOrigins) == 0 {
		return fmt.Errorf("production: cors.allowed_origins is empty; set EZ_CORS_ALLOWED_ORIGINS to your frontend domain(s)")
	}

	if c.Swagger.Enabled {
		return fmt.Errorf("production: swagger.enabled must be false; disable EZ_SWAGGER_ENABLED")
	}

	if c.Upload.MaxSizeMB > 50 {
		return fmt.Errorf("production: upload.max_size_mb (%d) exceeds safety limit of 50 MB", c.Upload.MaxSizeMB)
	}

	return nil
}

// AccessTokenDuration 返回访问令牌有效期。
func (c AuthConfig) AccessTokenDuration() time.Duration {
	return time.Duration(c.AccessTokenTTL) * time.Second
}

// RefreshTokenDuration 返回刷新令牌有效期；非正数时使用安全兜底值。
func (c AuthConfig) RefreshTokenDuration() time.Duration {
	if c.RefreshTokenTTL <= 0 {
		return defaultRefreshTokenTTL
	}
	return time.Duration(c.RefreshTokenTTL) * time.Second
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
	v.SetDefault("rate_limit.login_lockout_threshold", 5)
	v.SetDefault("rate_limit.login_lockout_sec", 300)
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
		"rate_limit.login_lockout_threshold",
		"rate_limit.login_lockout_sec",
	}

	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}
