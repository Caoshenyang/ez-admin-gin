package main

import (
	"os"

	"ez-admin-gin/server/internal/bootstrap"
)

func main() {
	// cmd/server 作为 v2 入口，依赖运行时工作目录下的迁移文件。
	bootstrap.MustRun(os.DirFS("."), "configs/rbac_model.conf")
}
