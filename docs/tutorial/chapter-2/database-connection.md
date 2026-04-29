---
title: 数据库连接
description: "使用 GORM 连接 PostgreSQL，并把数据库状态接入健康检查。"
---

# 数据库连接

前面已经准备好了 PostgreSQL 容器、配置文件和日志系统。现在把数据库连接接进后端服务，让服务启动时能连接数据库，并通过 `/health` 判断数据库是否可用。

::: tip 🎯 本节目标
完成后，服务启动时会连接 PostgreSQL，控制台会输出数据库连接日志，`/health` 接口会返回数据库状态。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ configs/
│  └─ config.yaml
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  └─ database/
│     └─ database.go
└─ main.go
```

| 位置 | 用途 |
| --- | --- |
| `configs/config.yaml` | 增加数据库连接池配置 |
| `internal/config/config.go` | 扩展数据库配置结构和环境变量绑定 |
| `internal/database/database.go` | 创建数据库连接、设置连接池、提供健康检查 |
| `main.go` | 启动时连接数据库，并把数据库状态接入 `/health` |

::: warning ⚠️ 先启动 PostgreSQL
继续之前，确认第一章的基础环境已经启动：

```bash
# 在项目根目录查看 PostgreSQL 和 Redis 状态
docker compose -f deploy/compose.local.yml ps
```

`postgres` 应该处于 `running` 或 `healthy` 状态。
:::

## 🛠️ 安装数据库依赖

进入 `server/` 目录：

::: code-group

```powershell [Windows PowerShell]
# 进入服务端目录
Set-Location .\server
```

```bash [macOS / Linux]
# 进入服务端目录
cd server
```

:::

安装 GORM、数据库驱动和 golang-migrate：

```bash
# 安装 ORM、PostgreSQL 驱动和 MySQL 驱动
go get gorm.io/gorm@latest gorm.io/driver/postgres@latest gorm.io/driver/mysql@latest

# 安装数据库迁移工具
go get github.com/golang-migrate/migrate/v4@latest
```

依赖资料入口：

| 依赖 | 用途 | 资料 |
| --- | --- | --- |
| `gorm.io/gorm` | ORM 主库，负责模型映射、查询、事务等能力 | [Go 包文档](https://pkg.go.dev/gorm.io/gorm) / [官方文档](https://gorm.io/docs/) |
| `gorm.io/driver/postgres` | GORM 的 PostgreSQL 驱动 | [Go 包文档](https://pkg.go.dev/gorm.io/driver/postgres) |
| `gorm.io/driver/mysql` | GORM 的 MySQL 驱动 | [Go 包文档](https://pkg.go.dev/gorm.io/driver/mysql) |
| `github.com/golang-migrate/migrate/v4` | 数据库迁移工具，管理建表和种子数据 | [GitHub](https://github.com/golang-migrate/migrate) / [数据库迁移工具选型](../../reference/migration-tool-selection) |

::: details 为什么先用 GORM
后台底座后续会频繁操作用户、角色、菜单、日志等表。GORM 可以先提供稳定的数据访问入口，后面需要复杂 SQL 时，也可以在局部直接写原生查询。
:::

## 🛠️ 扩展数据库配置

修改 `server/configs/config.yaml`，给 `database` 增加连接池配置：

```yaml
database:
  # 数据库配置和第一章的 Docker Compose 保持一致。
  host: localhost
  port: 5432
  user: ez_admin
  password: ez_admin_123456
  name: ez_admin
  # 空闲连接数，适合本地开发和小型后台的默认值。
  max_idle_conns: 10 # [!code ++]
  # 最大打开连接数，避免请求高峰时无限制创建连接。
  max_open_conns: 50 # [!code ++]
  # 连接最长复用时间，单位秒。
  conn_max_lifetime: 3600 # [!code ++]
```

字段含义：

| 字段 | 说明 |
| --- | --- |
| `max_idle_conns` | 保持多少个空闲连接，减少频繁创建连接 |
| `max_open_conns` | 最多允许多少个数据库连接 |
| `conn_max_lifetime` | 单个连接最多复用多久，单位秒 |

## 🛠️ 扩展配置结构

修改 `server/internal/config/config.go`。这一处有三个改动：

- 给 `DatabaseConfig` 增加连接池字段。
- 给数据库连接池补默认值。
- 给数据库连接池补环境变量绑定。

先修改 `DatabaseConfig`：

```go
// DatabaseConfig 保存数据库连接配置。
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	// MaxIdleConns 控制空闲连接数量。
	MaxIdleConns int `mapstructure:"max_idle_conns"` // [!code ++]
	// MaxOpenConns 控制最大打开连接数量。
	MaxOpenConns int `mapstructure:"max_open_conns"` // [!code ++]
	// ConnMaxLifetime 控制连接最长复用时间，单位秒。
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"` // [!code ++]
}
```

继续在 `setDefaults` 里增加默认值：

```go
// 数据库连接池默认值适合本地开发和小型后台起步。
v.SetDefault("database.max_idle_conns", 10) // [!code ++]
v.SetDefault("database.max_open_conns", 50) // [!code ++]
v.SetDefault("database.conn_max_lifetime", 3600) // [!code ++]
```

最后在 `bindEnvs` 的 `keys` 列表中增加：

```go
// 允许用 EZ_DATABASE_MAX_OPEN_CONNS 这类环境变量覆盖连接池配置。
"database.max_idle_conns", // [!code ++]
"database.max_open_conns", // [!code ++]
"database.conn_max_lifetime", // [!code ++]
```

常用环境变量示例：

| 配置项 | 环境变量 |
| --- | --- |
| `database.host` | `EZ_DATABASE_HOST` |
| `database.port` | `EZ_DATABASE_PORT` |
| `database.max_open_conns` | `EZ_DATABASE_MAX_OPEN_CONNS` |
| `database.conn_max_lifetime` | `EZ_DATABASE_CONN_MAX_LIFETIME` |

## 🛠️ 创建数据库包

::: details `server/internal/database/database.go` — 数据库连接与迁移

```go
package database

import (
	"fmt"
	"time"

	"ez-admin-gin/server/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// New 创建数据库连接，并完成连接池设置和连通性检查。
func New(cfg config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	dialector, err := openDialector(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		// Warn 级别可以记录慢查询和潜在问题，同时避免本地开发日志过多。
		Logger: gormLogger.Default.LogMode(gormLogger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql database: %w", err)
	}

	// 连接池参数从配置文件读取，便于不同环境单独调整。
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Info(
		"database connected",
		zap.String("driver", cfg.Driver),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name),
	)

	return db, nil
}

// Ping 用于健康检查，确认数据库当前仍然可连接。
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

// Close 关闭底层数据库连接池。
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql database: %w", err)
	}

	return sqlDB.Close()
}

// MigrateDSN 返回 golang-migrate 需要的连接字符串。
// GORM 和 golang-migrate 对 DSN 格式要求不同，所以分开生成。
func MigrateDSN(cfg config.DatabaseConfig) (string, error) {
	switch cfg.Driver {
	case "postgres":
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			url.PathEscape(cfg.User),
			url.PathEscape(cfg.Password),
			cfg.Host,
			cfg.Port,
			url.PathEscape(cfg.Name),
		), nil
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		), nil
	default:
		return "", fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

// openDialector 根据配置返回对应的 GORM Dialector。
func openDialector(cfg config.DatabaseConfig) (gorm.Dialector, error) {
	switch cfg.Driver {
	case "postgres":
		return postgres.Open(dsnPostgres(cfg)), nil
	case "mysql":
		return mysql.Open(dsnMySQL(cfg)), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

// dsnPostgres 把配置转换成 PostgreSQL 连接字符串。
func dsnPostgres(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)
}

// dsnMySQL 把配置转换成 MySQL 连接字符串。
func dsnMySQL(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}
```

:::

::: warning ⚠️ `sslmode=disable` 只适合本地开发
这里连接的是本地 Docker PostgreSQL，所以关闭 SSL。生产环境如果连接云数据库或跨网络数据库，应按数据库服务商要求配置 SSL。
:::

## 🛠️ 在启动入口接入数据库

修改 `server/main.go`。这一处重点看四个变化：

- 新增 `net/http`，用于返回健康检查状态码。
- 引入 `internal/database` 包。
- 服务启动时创建数据库连接，并在退出时关闭连接池。
- `/health` 增加数据库连通性检查。

::: warning ⚠️ 你可能会碰到复制 import 不成功，被自动删除
如果只先粘贴 `import` 里的新增依赖，GoLand 可能会因为“暂时未使用”把它自动删掉，看起来像是刚粘贴进去就没了。

更稳的跟做方式是直接替换下面这份完整 `main.go`，或者先粘贴真正使用数据库连接的代码，再让 GoLand 自动整理 import。
:::

::: details Q：怎么关闭 GoLand 自动删除未使用 import？
可以关，主要看你是哪个触发方式。

保存时自动删 import：

1. 打开 `File` -> `Settings` -> `Tools` -> `Actions on Save`
2. 找到 `Optimize imports`
3. 取消勾选

这是最常见的“保存后 import 被删”。

输入过程中自动优化 import：

1. 打开 `File` -> `Settings` -> `Go` -> `Imports`
2. 找到 `Optimize imports on the fly`
3. 如果已经勾选，就取消勾选

JetBrains 官方文档也提到，`Optimize Imports` 会移除未使用 import；自动保存时优化在 `Tools | Actions on Save` 中配置。参考：[GoLand Auto import 文档](https://www.jetbrains.com/help/go/creating-and-optimizing-imports.html)。

不过不用急着彻底关。Go 本身不允许未使用 import，IDE 帮你整理通常是好事。跟着本教程复制代码时，更稳的方式是直接替换完整 `main.go`，或者先粘贴下面真正使用依赖的代码，再让 GoLand 自动整理 import。
:::

::: details `server/main.go` — 完整版

```go
package main

import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"
	"net/http" // [!code ++]

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database" // [!code ++]
	appLogger "ez-admin-gin/server/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 先读取配置，日志和数据库初始化都依赖配置。
	cfg, err := config.Load()
	if err != nil {
		stdlog.Fatalf("load config: %v", err)
	}

	// 根据配置创建结构化日志对象。
	log, err := appLogger.New(cfg.Log)
	if err != nil {
		stdlog.Fatalf("create logger: %v", err)
	}
	defer func() {
		_ = log.Sync()
	}()

	// 启动时连接数据库；连接失败就直接终止服务。
	db, err := database.New(cfg.Database, log) // [!code ++]
	if err != nil {
		log.Fatal("connect database", zap.Error(err)) // [!code ++]
	}
	defer func() {
		if err := database.Close(db); err != nil { // [!code ++]
			log.Error("close database", zap.Error(err)) // [!code ++]
		}
	}()

	// 使用 gin.New()，再手动挂载自定义中间件。
	r := gin.New()
	r.Use(appLogger.GinLogger(log), appLogger.GinRecovery(log))

	r.GET("/health", func(c *gin.Context) {
		if err := database.Ping(db); err != nil { // [!code ++]
			log.Error("database health check failed", zap.Error(err)) // [!code ++]
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error", // [!code ++]
				"env":      cfg.App.Env,
				"database": "unavailable", // [!code ++]
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"env":      cfg.App.Env,
			"database": "ok", // [!code ++]
		})
	})

	// 服务启动日志记录关键运行参数。
	log.Info(
		"server started",
		zap.String("addr", cfg.Server.Addr),
		zap.String("env", cfg.App.Env),
	)

	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatal("run server", zap.Error(err))
	}
}
```

:::

::: details 为什么启动时就连接数据库
数据库是后台底座的核心依赖。启动时尽早连接，可以让配置错误、数据库未启动、账号密码错误这类问题尽快暴露，而不是等到第一次业务请求时才失败。
:::

## ✅ 整理依赖并启动

整理依赖：

```bash
# 整理新增依赖，更新 go.mod 和 go.sum
go mod tidy
```

确认 PostgreSQL 正在运行：

```bash
# 在项目根目录执行
docker compose -f deploy/compose.local.yml ps
```

回到 `server/` 目录启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

启动后，控制台应该能看到类似日志：

```text
INFO	database connected	{"host": "localhost", "port": 5432, "database": "ez_admin"}
INFO	database migrations applied
INFO	server started	{"addr": ":8080", "env": "dev"}
```

项目所有表（系统表 + 种子数据）都在 `server/migrations/` 下的 SQL 迁移文件中定义。启动时 golang-migrate 自动执行，不需要手动建表。完整的建表语句参考 [数据库建表语句](/reference/database-ddl)。

## ✅ 验证健康检查

访问健康检查接口：

::: code-group

```powershell [Windows PowerShell]
# 访问健康检查接口
Invoke-RestMethod http://localhost:8080/health
```

```bash [macOS / Linux]
# 访问健康检查接口
curl http://localhost:8080/health
```

:::

应该看到：

```json
{
  "database": "ok",
  "env": "dev",
  "status": "ok"
}
```

这说明后端服务和 PostgreSQL 都处于可用状态。

## ✅ 验证数据库不可用的情况

停止 PostgreSQL：

```bash
# 在项目根目录停止 PostgreSQL 容器
docker compose -f deploy/compose.local.yml stop postgres
```

再次启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

这时服务应该启动失败，并在日志里看到 `connect database` 相关错误。验证完成后，重新启动 PostgreSQL：

```bash
# 在项目根目录重新启动 PostgreSQL 容器
docker compose -f deploy/compose.local.yml start postgres
```

::: warning ⚠️ 验证失败场景后记得恢复服务
后续章节会继续依赖 PostgreSQL。验证数据库不可用的情况后，记得重新执行 `start postgres`，并确认容器状态恢复正常。
:::

## 常见问题

::: details 提示 `connection refused`
通常是 PostgreSQL 容器没有启动，或者端口不是 `5432`。先执行：

```bash
# 查看 PostgreSQL 容器状态和端口
docker compose -f deploy/compose.local.yml ps
```
:::

::: details 提示 `password authentication failed`
检查 `configs/config.yaml` 中的 `database.user`、`database.password`、`database.name`，确认它们和 `deploy/compose.local.yml` 保持一致。
:::

::: details 修改了数据库配置但没有生效
如果使用了环境变量，环境变量优先级高于 `config.yaml`。例如 `EZ_DATABASE_HOST` 会覆盖 `database.host`。
:::

下一节开始接入 Redis：[Redis 连接](./redis-connection)。
