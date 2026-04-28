---
title: 路由分组与健康检查
description: "把路由注册从启动入口中拆出来，并保留可验证的健康检查接口。"
---

# 路由分组与健康检查

前面已经把配置、日志、数据库、Redis、统一响应都接进来了。现在继续整理路由结构：让 `main.go` 只负责初始化和启动服务，把具体路由注册放到独立包里。

::: tip 🎯 本节目标
完成后，健康检查会从独立 Handler 中返回；`main.go` 会变得更薄；后续新增登录、用户、权限接口时，只需要继续扩展路由包。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ internal/
│  ├─ handler/
│  │  └─ system/
│  │     └─ health.go
│  └─ router/
│     └─ router.go
└─ main.go
```

| 位置 | 用途 |
| --- | --- |
| `internal/handler/system/health.go` | 放健康检查处理函数 |
| `internal/router/router.go` | 统一创建路由引擎、挂载中间件、注册路由分组 |
| `main.go` | 移除内联路由，把路由创建交给 `router.New` |

::: info 本节不需要安装新依赖
这一节只做结构拆分，不需要执行 `go get`。
:::

## 路由结构先定下来

这一节先保留两个健康检查入口：

| 路径 | 用途 |
| --- | --- |
| `/health` | 给本地验证、部署探针、容器健康检查使用 |
| `/api/v1/system/health` | 放在接口版本分组下，方便管理台或调试工具统一访问 |

::: details 为什么要有 `/api/v1`
后台接口一旦被管理台调用，就会形成前后端约定。加上版本前缀后，后续如果接口结构发生不兼容变化，可以用 `/api/v2` 逐步迁移，而不是直接影响旧接口。

这一节只先把分组搭出来，不急着设计完整 API 规范。
:::

## 🛠️ 创建健康检查 Handler

::: details `server/internal/handler/system/health.go` — 健康检查 Handler

```go
package system

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/response"

	goredis "github.com/redis/go-redis/v9"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// HealthHandler 负责系统健康检查。
type HealthHandler struct {
	cfg         *config.Config
	db          *gorm.DB
	redisClient *goredis.Client
	log         *zap.Logger
}

// NewHealthHandler 创建健康检查 Handler。
func NewHealthHandler(
	cfg *config.Config,
	db *gorm.DB,
	redisClient *goredis.Client,
	log *zap.Logger,
) *HealthHandler {
	return &HealthHandler{
		cfg:         cfg,
		db:          db,
		redisClient: redisClient,
		log:         log,
	}
}

// Check 返回当前服务、数据库和 Redis 的健康状态。
func (h *HealthHandler) Check(c *gin.Context) {
	if err := database.Ping(h.db); err != nil {
		h.log.Error("database health check failed", zap.Error(err))
		response.Error(c, apperror.ServiceUnavailable("数据库不可用", err), h.log)
		return
	}

	if err := appRedis.Ping(h.redisClient); err != nil {
		h.log.Error("redis health check failed", zap.Error(err))
		response.Error(c, apperror.ServiceUnavailable("Redis 不可用", err), h.log)
		return
	}

	response.Success(c, gin.H{
		"env":      h.cfg.App.Env,
		"database": "ok",
		"redis":    "ok",
	})
}
```

:::

这个文件只做一件事：处理健康检查请求。它不负责创建路由，也不负责启动服务。

::: details 为什么要把 Handler 单独拆出来
`main.go` 适合放启动流程，不适合放越来越多的接口逻辑。把 Handler 拆出来后，后续新增用户、角色、菜单接口时，可以继续按模块放到不同目录里。
:::

## 🛠️ 创建路由包

::: details `server/internal/router/router.go` — 路由注册

```go
package router

import (
	"ez-admin-gin/server/internal/config"
	systemHandler "ez-admin-gin/server/internal/handler/system"
	appLogger "ez-admin-gin/server/internal/logger"

	goredis "github.com/redis/go-redis/v9"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Options 汇总路由层需要依赖的对象。
type Options struct {
	Config *config.Config
	Log    *zap.Logger
	DB     *gorm.DB
	Redis  *goredis.Client
}

// New 创建路由引擎，并统一注册中间件和路由分组。
func New(opts Options) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	registerSystemRoutes(r, opts)

	return r
}

// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.GET("/health", health.Check)
}
```

:::

这个包负责三件事：

- 创建路由引擎。
- 统一挂载日志和恢复中间件。
- 注册系统路由和后续 API 路由分组。

::: warning ⚠️ 不要把数据库和 Redis 连接写进路由包
路由包只接收已经初始化好的依赖。数据库、Redis、日志这些对象仍然在 `main.go` 中创建，这样启动流程会更清楚。
:::

## 🛠️ 改造启动入口

修改 `server/main.go`。这一处重点看三个变化：

- 移除 `internal/apperror`、`internal/response` 和 `github.com/gin-gonic/gin`。
- 新增 `internal/router`。
- 用 `router.New(...)` 创建路由引擎。

先调整 import：

```go
import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/apperror" // [!code --]
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/response" // [!code --]
	"ez-admin-gin/server/internal/router" // [!code ++]

	"github.com/gin-gonic/gin" // [!code --]
	"go.uber.org/zap"
)
```

再把原来 `main.go` 中创建路由和注册 `/health` 的代码，替换成下面这段：

```go
	// 路由注册交给 internal/router，main.go 只保留启动流程。
	r := router.New(router.Options{ // [!code ++]
		Config: cfg, // [!code ++]
		Log:    log, // [!code ++]
		DB:     db, // [!code ++]
		Redis:  redisClient, // [!code ++]
	}) // [!code ++]
```

原来的 `gin.New()`、中间件挂载、`r.GET("/health", ...)` 都会移动到 `internal/router` 和 `internal/handler/system` 中。

::: tip 更稳的修改方式
这一节删除的代码比较多。如果担心漏删旧的 `r.GET("/health", ...)`，可以直接替换完整 `main.go`，再让 IDE 自动整理 import。
:::

::: details `server/main.go` — 完整版

```go
package main

import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/router"

	"go.uber.org/zap"
)

func main() {
	// 先读取配置，日志、数据库、Redis 初始化都依赖配置。
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
	db, err := database.New(cfg.Database, log)
	if err != nil {
		log.Fatal("connect database", zap.Error(err))
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Error("close database", zap.Error(err))
		}
	}()

	// 启动时连接 Redis；连接失败就直接终止服务。
	redisClient, err := appRedis.New(cfg.Redis, log)
	if err != nil {
		log.Fatal("connect redis", zap.Error(err))
	}
	defer func() {
		if err := appRedis.Close(redisClient); err != nil {
			log.Error("close redis", zap.Error(err))
		}
	}()

	// 路由注册交给 internal/router，main.go 只保留启动流程。
	r := router.New(router.Options{
		Config: cfg,
		Log:    log,
		DB:     db,
		Redis:  redisClient,
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

## ✅ 启动并验证旧入口

整理依赖：

```bash
# 本节没有新增三方依赖，执行 tidy 只是确认 module 仍然干净
go mod tidy
```

确认 PostgreSQL 和 Redis 正在运行：

```bash
# 在项目根目录执行
docker compose -f deploy/compose.local.yml ps
```

回到 `server/` 目录启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

先访问原来的健康检查入口：

::: code-group

```powershell [Windows PowerShell]
# 验证保留的部署探针入口
Invoke-RestMethod http://localhost:8080/health
```

```bash [macOS / Linux]
# 验证保留的部署探针入口
curl http://localhost:8080/health
```

:::

应该看到类似结果：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "database": "ok",
    "env": "dev",
    "redis": "ok"
  }
}
```

## ✅ 验证分组后的新入口

继续访问版本分组下的健康检查入口：

::: code-group

```powershell [Windows PowerShell]
# 验证 /api/v1/system 分组下的健康检查
Invoke-RestMethod http://localhost:8080/api/v1/system/health
```

```bash [macOS / Linux]
# 验证 /api/v1/system 分组下的健康检查
curl http://localhost:8080/api/v1/system/health
```

:::

返回结构应该和 `/health` 一致：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "database": "ok",
    "env": "dev",
    "redis": "ok"
  }
}
```

::: warning ⚠️ 如果新入口返回 404
优先检查 `internal/router/router.go` 里是否注册了：

```go
api := r.Group("/api/v1")
system := api.Group("/system")
system.GET("/health", health.Check)
```

如果这里只注册了 `/health`，那么 `/api/v1/system/health` 会找不到路由。
:::

## 常见问题

::: details 为什么不直接把所有路由都写在 `main.go`
一开始写在 `main.go` 没问题，但接口变多后，启动流程、依赖初始化、路由注册、业务处理会混在一起。提前拆出 `router` 和 `handler`，后续扩展会更稳。
:::

::: details Handler 为什么要接收 `cfg`、`db`、`redisClient`、`log`
健康检查需要读取当前环境、检查数据库、检查 Redis，并在失败时记录日志。把这些对象通过构造函数传进去，比在 Handler 里重新创建连接更清楚。
:::

::: details 后续业务接口也都放在 `system` 目录吗
不是。`system` 只放系统级接口。后续用户、角色、菜单会按模块继续拆目录，避免所有 Handler 挤在一起。
:::

至此，第二章的后端基础设施已经有了清晰入口。下一章开始进入认证与权限。
