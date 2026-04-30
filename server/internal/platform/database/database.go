package database

import (
	legacyConfig "ez-admin-gin/server/internal/config"
	legacyDatabase "ez-admin-gin/server/internal/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// New 创建数据库连接。
func New(cfg legacyConfig.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	return legacyDatabase.New(cfg, log)
}

// Close 关闭数据库连接。
func Close(db *gorm.DB) error {
	return legacyDatabase.Close(db)
}

// MigrateDSN 生成迁移阶段使用的 DSN。
func MigrateDSN(cfg legacyConfig.DatabaseConfig) (string, error) {
	return legacyDatabase.MigrateDSN(cfg)
}
