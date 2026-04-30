package migrate

import (
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

// Run 执行数据库迁移。根据 driver 参数加载对应子目录下的 SQL 文件。
// 返回 nil 表示迁移成功（包括"已是最新"的情况）。
func Run(driver, dsn string, migrationsFS fs.FS, log *zap.Logger) error {
	sub, err := fs.Sub(migrationsFS, "migrations/"+driver)
	if err != nil {
		return fmt.Errorf("open migrations/%s: %w", driver, err)
	}

	source, err := iofs.New(sub, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, dsn)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		// 迁移失败时检查是否 dirty state（上次迁移中途崩溃），
		// 自动解锁后重试，避免需要手动修复 schema_migrations 表。
		version, dirty, vErr := m.Version()
		if vErr == nil && dirty {
			log.Warn("dirty migration detected, forcing unlock",
				zap.Uint("version", version))
			if forceErr := m.Force(int(version)); forceErr != nil {
				return fmt.Errorf("force unlock dirty migration: %w", forceErr)
			}
			err = m.Up()
		}
	}
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Info("database migrations up to date", zap.String("driver", driver))
	} else {
		log.Info("database migrations applied", zap.String("driver", driver))
	}

	return nil
}
