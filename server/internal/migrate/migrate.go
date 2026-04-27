package migrate

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

// Run 执行数据库迁移。根据 driver 参数加载对应子目录下的 SQL 文件。
func Run(driver, dsn string, migrationsFS embed.FS, log *zap.Logger) error {
	sub, err := fs.Sub(migrationsFS, driver)
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

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Info("database migrations up to date", zap.String("driver", driver))
	} else {
		log.Info("database migrations applied", zap.String("driver", driver))
	}

	return nil
}
