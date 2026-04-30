package migrate

import (
	"io/fs"

	legacyMigrate "ez-admin-gin/server/internal/migrate"

	"go.uber.org/zap"
)

// Run 执行 SQL 迁移。
func Run(driver string, dsn string, migrationsFS fs.FS, log *zap.Logger) error {
	return legacyMigrate.Run(driver, dsn, migrationsFS, log)
}
