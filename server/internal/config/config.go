package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 汇总整个服务端会读取的配置段。
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
}

// AppConfig 保存应用自身信息。
type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

// ServerConfig 保存 HTTP 服务启动配置。
type ServerConfig struct {
	Addr string `mapstructure:"addr"`
}

// DatabaseConfig 保存数据库连接配置。
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	// MaxIdleConns 控制空闲连接数量。
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// MaxOpenConns 控制最大打开连接数量。
	MaxOpenConns int `mapstructure:"max_open_conns"`
	// ConnMaxLifetime 控制连接最长复用时间，单位秒。
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`
}

// RedisConfig 保存 Redis 连接配置。
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	// MaxRetries 控制 Redis 命令失败时的最大重试次数。
	MaxRetries int `mapstructure:"max_retries"`
	// MinIdleConns 控制最少保留多少个空闲连接。
	MinIdleConns int `mapstructure:"min_idle_conns"`
	// PoolSize 控制连接池大小。
	PoolSize int `mapstructure:"pool_size"`
}

// LogConfig 保存日志级别、格式和文件切割配置。
type LogConfig struct {
	// Level 控制输出哪些级别的日志。
	Level string `mapstructure:"level"`
	// Format 控制日志格式，支持 console 和 json。
	Format string `mapstructure:"format"`
	// Filename 是日志文件路径。为空时只输出到控制台。
	Filename string `mapstructure:"filename"`
	// MaxSize 是单个日志文件最大大小，单位 MB。
	MaxSize int `mapstructure:"max_size"`
	// MaxBackups 是最多保留的旧日志文件数量。
	MaxBackups int `mapstructure:"max_backups"`
	// MaxAge 是日志文件最多保留天数。
	MaxAge int `mapstructure:"max_age"`
	// Compress 控制是否压缩旧日志文件。
	Compress bool `mapstructure:"compress"`
}

// Load 读取配置文件，并把结果解析到 Config 结构体中。
func Load() (*Config, error) {
	v := viper.New()

	// 配置文件位置是 server/configs/config.yaml。
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	// 先设置默认值，再绑定环境变量。
	setDefaults(v)
	bindEnvs(v)

	// EZ_SERVER_ADDR 这类环境变量会覆盖 server.addr。
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

// setDefaults 设置兜底值，避免配置文件缺少字段时直接变成零值。
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "ez-admin")
	v.SetDefault("app.env", "dev")
	v.SetDefault("server.addr", ":8080")
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
	// 数据库连接池默认值适合本地开发和小型后台起步。
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 50)
	v.SetDefault("database.conn_max_lifetime", 3600)
	// Redis 连接池默认值适合本地开发和小型后台起步。
	v.SetDefault("redis.max_retries", 3)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.pool_size", 10)
}

// bindEnvs 让环境变量能稳定参与结构体解析。
func bindEnvs(v *viper.Viper) {
	keys := []string{
		"app.name",
		"app.env",
		"server.addr",
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
		// 允许用 EZ_DATABASE_MAX_OPEN_CONNS 这类环境变量覆盖连接池配置。
		"database.max_idle_conns",
		"database.max_open_conns",
		"database.conn_max_lifetime",
		// 允许用 EZ_REDIS_POOL_SIZE 这类环境变量覆盖 Redis 连接池配置。
		"redis.max_retries",
		"redis.min_idle_conns",
		"redis.pool_size",
	}

	for _, key := range keys {
		// BindEnv 返回错误通常来自 key 本身，这里的 key 是固定列表。
		_ = v.BindEnv(key)
	}
}
