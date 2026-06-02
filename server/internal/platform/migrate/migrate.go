package migrate

import (
	"database/sql"
	"fmt"
	"io/fs"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const fullSchemaAndSeedFile = "full_schema_and_seed.sql"

// Run 执行数据库迁移；遇到脏锁时自动强制解锁后重试。
func Run(driver, dsn string, migrationsFS fs.FS, log *zap.Logger) error {
	sub, err := fs.Sub(migrationsFS, "migrations/"+driver)
	if err != nil {
		return fmt.Errorf("open migrations/%s: %w", driver, err)
	}

	hasVersioned, err := hasVersionedMigrations(sub)
	if err != nil {
		return fmt.Errorf("inspect migrations/%s: %w", driver, err)
	}
	if !hasVersioned {
		return runFullSchemaAndSeed(driver, dsn, sub, log)
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
		// 脏迁移通常是上次中断导致，强制解锁后重试自动恢复。
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

func hasVersionedMigrations(migrations fs.FS) (bool, error) {
	entries, err := fs.ReadDir(migrations, ".")
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if isVersionedUpMigration(entry.Name()) {
			return true, nil
		}
	}
	return false, nil
}

func isVersionedUpMigration(name string) bool {
	if !strings.HasSuffix(name, ".up.sql") || len(name) < len("000001_x.up.sql") {
		return false
	}
	if name[6] != '_' {
		return false
	}
	for i := 0; i < 6; i++ {
		if name[i] < '0' || name[i] > '9' {
			return false
		}
	}
	return true
}

func runFullSchemaAndSeed(driver, dsn string, migrations fs.FS, log *zap.Logger) error {
	sqlText, err := fs.ReadFile(migrations, fullSchemaAndSeedFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", fullSchemaAndSeedFile, err)
	}

	dbDriver, dbDSN, err := fullSchemaSQLDriver(driver, dsn)
	if err != nil {
		return err
	}

	db, err := sql.Open(dbDriver, dbDSN)
	if err != nil {
		return fmt.Errorf("open database for full initialization: %w", err)
	}
	defer db.Close()

	initialized, partial, err := fullInitializationState(db, driver)
	if err != nil {
		return fmt.Errorf("check full initialization state: %w", err)
	}
	if initialized {
		log.Info("database full initialization up to date", zap.String("driver", driver))
		return nil
	}
	if partial {
		return fmt.Errorf("database has existing schema but built-in seed is incomplete; clear the database or import %s manually", fullSchemaAndSeedFile)
	}

	if _, err := db.Exec(string(sqlText)); err != nil {
		return fmt.Errorf("run %s: %w", fullSchemaAndSeedFile, err)
	}

	log.Info("database full initialization applied", zap.String("driver", driver))
	return nil
}

func fullSchemaSQLDriver(driver, dsn string) (string, string, error) {
	switch driver {
	case "postgres":
		return "postgres", dsn, nil
	case "mysql":
		if strings.Contains(dsn, "?") {
			return "mysql", dsn + "&multiStatements=true", nil
		}
		return "mysql", dsn + "?multiStatements=true", nil
	default:
		return "", "", fmt.Errorf("unsupported database driver: %s", driver)
	}
}

func fullInitializationState(db *sql.DB, driver string) (initialized bool, partial bool, err error) {
	totalTables, err := countTables(db, driver, "")
	if err != nil {
		return false, false, err
	}
	if totalTables == 0 {
		return false, false, nil
	}

	requiredTables := []string{"sys_role", "sys_menu", "casbin_rule"}
	for _, table := range requiredTables {
		exists, err := tableExists(db, driver, table)
		if err != nil {
			return false, false, err
		}
		if !exists {
			return false, true, nil
		}
	}

	var roleCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM sys_role WHERE code = 'super_admin'").Scan(&roleCount); err != nil {
		return false, false, err
	}
	var menuCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM sys_menu WHERE code IN ('iam', 'system', 'audit')").Scan(&menuCount); err != nil {
		return false, false, err
	}
	var policyCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM casbin_rule WHERE ptype = 'p' AND v0 = 'super_admin'").Scan(&policyCount); err != nil {
		return false, false, err
	}

	return roleCount > 0 && menuCount >= 3 && policyCount > 0, true, nil
}

func tableExists(db *sql.DB, driver, table string) (bool, error) {
	count, err := countTables(db, driver, table)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func countTables(db *sql.DB, driver, table string) (int, error) {
	var tableCount int
	var tableQuery string
	switch driver {
	case "postgres":
		tableQuery = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'"
	case "mysql":
		tableQuery = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()"
	default:
		return 0, fmt.Errorf("unsupported database driver: %s", driver)
	}
	if table != "" {
		switch driver {
		case "postgres":
			tableQuery += " AND table_name = $1"
		case "mysql":
			tableQuery += " AND table_name = ?"
		}
		if err := db.QueryRow(tableQuery, table).Scan(&tableCount); err != nil {
			return 0, err
		}
		return tableCount, nil
	}

	if err := db.QueryRow(tableQuery).Scan(&tableCount); err != nil {
		return 0, err
	}
	return tableCount, nil
}
