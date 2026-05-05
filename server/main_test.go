package main

import (
	"io/fs"
	"reflect"
	"sort"
	"testing"
)

func TestEmbeddedMigrationsStayInSyncAcrossDrivers(t *testing.T) {
	mysqlFiles := readMigrationFiles(t, "migrations/mysql")
	postgresFiles := readMigrationFiles(t, "migrations/postgres")

	if !reflect.DeepEqual(mysqlFiles, postgresFiles) {
		t.Fatalf("expected mysql and postgres migration files to stay in sync\nmysql=%v\npostgres=%v", mysqlFiles, postgresFiles)
	}

	required := []string{
		"000010_phase4_followup_schema.down.sql",
		"000010_phase4_followup_schema.up.sql",
		"000011_phase4_followup_seed_data.down.sql",
		"000011_phase4_followup_seed_data.up.sql",
	}
	for _, name := range required {
		if !containsFile(mysqlFiles, name) {
			t.Fatalf("expected embedded migrations to include %s, got %v", name, mysqlFiles)
		}
	}
}

func readMigrationFiles(t *testing.T, dir string) []string {
	t.Helper()

	entries, err := fs.ReadDir(migrationsFS, dir)
	if err != nil {
		t.Fatalf("read embedded migrations in %s: %v", dir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	return files
}

func containsFile(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
