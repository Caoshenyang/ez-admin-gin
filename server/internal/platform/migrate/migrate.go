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
		version, dirty, versionErr := m.Version()
		if versionErr == nil && dirty {
			log.Warn("dirty migration detected, forcing unlock", zap.Uint("version", version))
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
