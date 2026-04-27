package database

import (
	"fmt"
	"time"

	"ez-admin-gin/server/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// New 创建数据库连接，并完成连接池设置和连通性检查。
func New(cfg config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	dialector, err := openDialector(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		// Warn 级别可以记录慢查询和潜在问题，同时避免本地开发日志过多。
		Logger: gormLogger.Default.LogMode(gormLogger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql database: %w", err)
	}

	// 连接池参数从配置文件读取，便于不同环境单独调整。
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Info(
		"database connected",
		zap.String("driver", cfg.Driver),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name),
	)

	return db, nil
}

// Ping 用于健康检查，确认数据库当前仍然可连接。
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

// Close 关闭底层数据库连接池。
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql database: %w", err)
	}

	return sqlDB.Close()
}

// DSN 返回当前驱动对应的连接字符串。
func DSN(cfg config.DatabaseConfig) (string, error) {
	switch cfg.Driver {
	case "postgres":
		return dsnPostgres(cfg), nil
	case "mysql":
		return dsnMySQL(cfg), nil
	default:
		return "", fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

// openDialector 根据配置返回对应的 GORM Dialector。
func openDialector(cfg config.DatabaseConfig) (gorm.Dialector, error) {
	switch cfg.Driver {
	case "postgres":
		return postgres.Open(dsnPostgres(cfg)), nil
	case "mysql":
		return mysql.Open(dsnMySQL(cfg)), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

// dsnPostgres 把配置转换成 PostgreSQL 连接字符串。
func dsnPostgres(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)
}

// dsnMySQL 把配置转换成 MySQL 连接字符串。
func dsnMySQL(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}
