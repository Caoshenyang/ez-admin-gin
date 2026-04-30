package main

import (
	"embed"

	"ez-admin-gin/server/internal/bootstrap"
)

//go:embed migrations/postgres migrations/mysql
var migrationsFS embed.FS

func main() {
	bootstrap.MustRun(migrationsFS, "configs/rbac_model.conf")
}
