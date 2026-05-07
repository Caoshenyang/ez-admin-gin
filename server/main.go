package main

import (
	"embed"

	"ez-admin-gin/server/internal/bootstrap"
)

// @title           EZ Admin Gin API
// @version         1.0
// @description     EZ Admin 后台管理系统 API，包含认证、用户管理、角色管理、系统配置等功能。
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     输入 Bearer {token}

//go:embed migrations/postgres migrations/mysql
var migrationsFS embed.FS

func main() {
	bootstrap.MustRun(migrationsFS, "configs/rbac_model.conf")
}
