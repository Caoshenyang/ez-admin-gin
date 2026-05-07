package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	platformConfig "ez-admin-gin/server/internal/platform/config"
	platformDatabase "ez-admin-gin/server/internal/platform/database"
	platformMigrate "ez-admin-gin/server/internal/platform/migrate"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestEmbeddedMigrationsApplyOnRealDatabases(t *testing.T) {
	if os.Getenv("EZ_REAL_DB_MIGRATION") != "1" {
		t.Skip("set EZ_REAL_DB_MIGRATION=1 to run real database migration smoke tests")
	}

	drivers := resolveRealDBDrivers(t)
	for _, driver := range drivers {
		driver := driver
		t.Run(driver, func(t *testing.T) {
			cfg := loadRealDBConfig(t, driver)
			verifyRealDBMigrations(t, cfg)
		})
	}
}

func resolveRealDBDrivers(t *testing.T) []string {
	t.Helper()

	raw := strings.TrimSpace(os.Getenv("EZ_REAL_DB_DRIVERS"))
	if raw == "" {
		return []string{"postgres", "mysql"}
	}

	parts := strings.Split(raw, ",")
	drivers := make([]string, 0, len(parts))
	for _, part := range parts {
		driver := strings.TrimSpace(part)
		if driver == "" {
			continue
		}
		switch driver {
		case "postgres", "mysql":
			drivers = append(drivers, driver)
		default:
			t.Fatalf("unsupported EZ_REAL_DB_DRIVERS value %q", driver)
		}
	}
	if len(drivers) == 0 {
		t.Fatalf("EZ_REAL_DB_DRIVERS resolved to no supported drivers")
	}
	return drivers
}

func loadRealDBConfig(t *testing.T, driver string) platformConfig.DatabaseConfig {
	t.Helper()

	prefix := "EZ_REAL_DB_" + strings.ToUpper(driver) + "_"
	return platformConfig.DatabaseConfig{
		Driver:          driver,
		Host:            requireEnv(t, prefix+"HOST"),
		Port:            requireEnvInt(t, prefix+"PORT"),
		User:            requireEnv(t, prefix+"USER"),
		Password:        requireEnv(t, prefix+"PASSWORD"),
		Name:            requireEnv(t, prefix+"NAME"),
		MaxIdleConns:    2,
		MaxOpenConns:    4,
		ConnMaxLifetime: 60,
	}
}

func verifyRealDBMigrations(t *testing.T, cfg platformConfig.DatabaseConfig) {
	t.Helper()

	log := zap.NewNop()
	dsn, err := platformDatabase.MigrateDSN(cfg)
	if err != nil {
		t.Fatalf("build migrate dsn for %s: %v", cfg.Driver, err)
	}

	if err := platformMigrate.Run(cfg.Driver, dsn, migrationsFS, log); err != nil {
		t.Fatalf("apply migrations on %s: %v", cfg.Driver, err)
	}
	if err := platformMigrate.Run(cfg.Driver, dsn, migrationsFS, log); err != nil {
		t.Fatalf("re-apply migrations on %s should stay idempotent: %v", cfg.Driver, err)
	}

	db, err := platformDatabase.New(cfg, log)
	if err != nil {
		t.Fatalf("open %s after migrations: %v", cfg.Driver, err)
	}
	defer func() {
		if closeErr := platformDatabase.Close(db); closeErr != nil {
			t.Fatalf("close %s after migrations: %v", cfg.Driver, closeErr)
		}
	}()

	assertMigrationVersion(t, db, latestEmbeddedMigrationVersion(t, cfg.Driver))
	for _, table := range []string{
		"sys_user",
		"sys_role",
		"sys_department",
		"sys_dict_type",
		"sys_attachment",
		"sys_customer",
	} {
		assertTableAccessible(t, db, table)
	}
	assertSeedRowsExist(t, db, "sys_role")
	assertSeedRowsExist(t, db, "sys_menu")
}

func latestEmbeddedMigrationVersion(t *testing.T, driver string) int64 {
	t.Helper()

	files := readMigrationFiles(t, filepath.Join("migrations", driver))
	var latest int64
	for _, name := range files {
		versionPart, _, found := strings.Cut(name, "_")
		if !found {
			t.Fatalf("unexpected migration file name %q", name)
		}
		version, err := strconv.ParseInt(versionPart, 10, 64)
		if err != nil {
			t.Fatalf("parse migration version from %q: %v", name, err)
		}
		if version > latest {
			latest = version
		}
	}
	if latest == 0 {
		t.Fatalf("no embedded migrations found for %s", driver)
	}
	return latest
}

func assertMigrationVersion(t *testing.T, db *gorm.DB, expected int64) {
	t.Helper()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db for schema_migrations query: %v", err)
	}

	var version int64
	var dirty bool
	if err := sqlDB.QueryRow("SELECT version, dirty FROM schema_migrations LIMIT 1").Scan(&version, &dirty); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if dirty {
		t.Fatalf("expected schema_migrations dirty=false, got version=%d dirty=true", version)
	}
	if version != expected {
		t.Fatalf("expected migration version %d, got %d", expected, version)
	}
}

func assertTableAccessible(t *testing.T, db *gorm.DB, table string) {
	t.Helper()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db for table query %s: %v", table, err)
	}

	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if err := sqlDB.QueryRow(query).Scan(&count); err != nil {
		t.Fatalf("query table %s after migrations: %v", table, err)
	}
}

func assertSeedRowsExist(t *testing.T, db *gorm.DB, table string) {
	t.Helper()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db for seed query %s: %v", table, err)
	}

	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if err := sqlDB.QueryRow(query).Scan(&count); err != nil {
		t.Fatalf("query seed rows in %s: %v", table, err)
	}
	if count == 0 {
		t.Fatalf("expected seeded rows in %s, got 0", table)
	}
}

func requireEnv(t *testing.T, key string) string {
	t.Helper()

	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		t.Fatalf("required environment variable %s is not set", key)
	}
	return value
}

func requireEnvInt(t *testing.T, key string) int {
	t.Helper()

	value := requireEnv(t, key)
	parsed, err := strconv.Atoi(value)
	if err != nil {
		t.Fatalf("parse environment variable %s=%q as int: %v", key, value, err)
	}
	return parsed
}
