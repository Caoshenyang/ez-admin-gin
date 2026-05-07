---
title: 路由分组与健康检查
description: "按当前代码梳理启动入口、路由注册和健康检查的真实调用链。"
---

# 路由分组与健康检查

前面几节已经把配置、日志、数据库、Redis 和统一响应准备好了。这一节不再按早期版本那种“把路由直接写在 `main.go`”的方式展开，而是直接对齐当前仓库的真实结构：

`server/main.go` -> `internal/bootstrap/run.go` -> `internal/bootstrap/router.go` -> `internal/module/*/routes.go` -> `internal/module/system/health_handler.go`

::: tip 这页怎么读
如果你只想先抓主线，就先记住一句话：当前启动入口只负责把迁移文件系统交给 `bootstrap.MustRun`，真正的初始化编排在 `internal/bootstrap/run.go`，真正的路由聚合在 `internal/bootstrap/router.go`，真正的健康检查逻辑在 `internal/module/system/health_handler.go`。
:::

## 本节会看哪些文件

```text
server/
├─ main.go
└─ internal/
   ├─ app/
   │  ├─ handler.go
   │  └─ pagination.go
   ├─ bootstrap/
   │  ├─ run.go
   │  └─ router.go
   └─ module/
      └─ system/
         ├─ routes.go
         └─ health_handler.go
```

| 位置 | 当前职责 |
| --- | --- |
| `server/main.go` | 嵌入迁移文件，并调用 `bootstrap.MustRun(...)` |
| `internal/bootstrap/run.go` | 读取配置、初始化日志/数据库/Redis/鉴权、启动 HTTP 服务 |
| `internal/bootstrap/router.go` | 创建 Gin 引擎，挂载中间件，聚合模块路由 |
| `internal/module/system/routes.go` | 注册 `/health` 和 `/api/v1/system/*` 路由 |
| `internal/module/system/health_handler.go` | 执行数据库与 Redis 健康检查，返回统一响应 |
| `internal/app/*` | 给业务 Handler 提供统一错误写回、登录人提取、ID 参数解析等公共助手 |

::: warning 这一节不再以 `cmd/server/main.go` 为准
旧教程里常见的 `cmd/server/main.go`、`main.go` 里手写 `gin.New()`、或者在入口里直接注册 `/health`，都不是当前主线结构了。现在的主线以 `bootstrap / platform / module / app` 这一套分层为准。
:::

## 先看最薄的启动入口

当前 `server/main.go` 只保留了一件事：把嵌入的迁移文件和 RBAC 模型路径交给 `bootstrap.MustRun`。

```go
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
```

这和早期章节里“入口自己读配置、建数据库、建路由、再 `r.Run(...)`”的写法已经不一样了。现在入口文件的目标很明确：

- 暴露程序起点
- 提供迁移文件系统
- 不再承载初始化细节

这样做的直接好处是：入口不会随着功能增加而越来越胖，教程后面再加模块时，也不用反复改启动文件的主体结构。

## 真正的启动编排在 `bootstrap/run.go`

当前启动顺序由 `internal/bootstrap/run.go` 统一负责：

1. 读取配置
2. 初始化日志
3. 初始化数据库并执行迁移
4. 初始化 Redis
5. 初始化 Token 管理器和权限执行器
6. 创建路由引擎
7. 启动 HTTP 服务

核心代码可以先看这一段：

```go
func MustRun(migrationsFS fs.FS, rbacModelPath string) {
	cfg, err := platformConfig.Load()
	if err != nil {
		stdlog.Fatalf("load config: %v", err)
	}

	log, err := platformLogger.New(cfg.Log)
	if err != nil {
		stdlog.Fatalf("create logger: %v", err)
	}

	db, err := platformDatabase.New(cfg.Database, log)
	if err != nil {
		log.Fatal("connect database", zap.Error(err))
	}

	migrateDSN, err := platformDatabase.MigrateDSN(cfg.Database)
	if err != nil {
		log.Fatal("build migration dsn", zap.Error(err))
	}
	if err := platformMigrate.Run(cfg.Database.Driver, migrateDSN, migrationsFS, log); err != nil {
		log.Fatal("run database migrations", zap.Error(err))
	}

	redisClient, err := platformRedis.New(cfg.Redis, log)
	if err != nil {
		log.Fatal("connect redis", zap.Error(err))
	}

	tokenManager, err := authnPlatform.NewManager(cfg.Auth)
	if err != nil {
		log.Fatal("create token manager", zap.Error(err))
	}

	permissionEnforcer, err := authzPlatform.NewEnforcer(db, rbacModelPath)
	if err != nil {
		log.Fatal("create permission enforcer", zap.Error(err))
	}

	r := NewRouter(RouterOptions{
		Config:     cfg,
		Log:        log,
		DB:         db,
		Redis:      redisClient,
		Token:      tokenManager,
		Permission: permissionEnforcer,
	})

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

::: details 为什么这里要区分 `main.go` 和 `run.go`
`main.go` 是程序入口，越薄越好；`run.go` 是启动编排，适合承接“先做什么、后做什么、失败如何退出、资源如何关闭”这类流程。

这样后面要加 CLI 子命令、测试入口、不同运行模式时，也更容易复用这层编排，而不是把所有逻辑重新塞回入口文件。
:::

## 路由创建集中在 `bootstrap/router.go`

当前 `internal/bootstrap/router.go` 负责：

- 创建 Gin 引擎
- 挂载日志与恢复中间件
- 设置上传大小与静态目录
- 聚合 `auth`、`setup`、`system` 三个模块的路由注册

```go
func NewRouter(opts RouterOptions) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	if opts.Config.Upload.MaxSizeMB > 0 {
		r.MaxMultipartMemory = opts.Config.Upload.MaxSizeMB << 20
	}
	r.Static(opts.Config.Upload.PublicPath, opts.Config.Upload.Dir)

	authModule.RegisterRoutes(r, authModule.RouteOptions{
		Config: opts.Config,
		Log:    opts.Log,
		DB:     opts.DB,
		Redis:  opts.Redis,
		Token:  opts.Token,
	})
	setupModule.RegisterRoutes(r, setupModule.RouteOptions{
		Log: opts.Log,
		DB:  opts.DB,
	})
	systemModule.RegisterRoutes(r, systemModule.RouteOptions{
		Config:     opts.Config,
		Log:        opts.Log,
		DB:         opts.DB,
		Redis:      opts.Redis,
		Token:      opts.Token,
		Permission: opts.Permission,
	})

	return r
}
```

这里有两个要点很容易和旧教程混淆：

1. `router.go` 现在不再自己手写 `registerSystemRoutes` 这种局部函数，而是把职责交给各模块的 `RegisterRoutes`。
2. 健康检查虽然和系统模块相关，但它已经属于 `systemModule.RegisterRoutes(...)` 的一部分，不再由启动入口直接注册。

## 健康检查实际挂在哪里

当前健康检查的真实挂载点在 `server/internal/module/system/routes.go`：

```go
func RegisterRoutes(r *gin.Engine, opts RouteOptions) {
	health := newHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)

	r.GET("/health", health.Check)

	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.LoadActor(opts.DB, opts.Log))
	system.Use(middleware.OperationLog(opts.DB, opts.Log))
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))

	system.GET("/health", health.Check)
	// 省略其他 system / iam 路由注册
}
```

所以当前项目里这两个入口都存在：

| 路径 | 当前状态 | 说明 |
| --- | --- | --- |
| `/health` | 公开可访问 | 适合容器探针、本地验证、反向代理健康检查 |
| `/api/v1/system/health` | 挂在 `system` 组内 | 会经过系统分组上的中间件链 |

::: warning `/api/v1/system/health` 不再是“无条件公开调试接口”
旧版文档里容易把它描述成“和 `/health` 一样只是多了版本前缀”。但按当前代码，它位于 `system` 路由组内部，组上已经挂了认证、加载当前人、操作日志、权限校验等中间件。

也就是说，想理解当前项目的健康检查时，应该把 `/health` 看成探针入口，把 `/api/v1/system/health` 看成系统分组下的一个模块路由节点，而不是把两者当作完全等价的公开接口。
:::

## 健康检查 Handler 只做依赖探活

真正的探活逻辑在 `server/internal/module/system/health_handler.go`：

```go
func (h *healthHandler) Check(c *gin.Context) {
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

这层职责很单纯：

- 读取当前环境 `cfg.App.Env`
- 探活数据库
- 探活 Redis
- 通过统一响应包返回结果

它不做这些事情：

- 不创建数据库连接
- 不创建 Redis 客户端
- 不注册路由
- 不决定服务监听地址

这些工作分别在 `bootstrap/run.go` 和 `bootstrap/router.go` 已经完成。

## `internal/app/*` 在这条线里的位置

虽然健康检查本身没有用到 `internal/app/*` 的所有助手，但这一层已经是当前 Handler 编写的公共基座，后续业务接口都会沿用：

| 文件 | 作用 |
| --- | --- |
| `internal/app/handler.go` | `WriteError`、`CurrentActor`、`UintIDParam` 等 Handler 级公共助手 |
| `internal/app/pagination.go` | 分页参数规范化 |

后续写业务 Handler 时，可以优先复用这一层，而不是在每个模块里重复写“取当前登录人”“解析 ID”“兜底错误处理”。

## ✅ 启动并验证

启动服务：

```bash
# 在 server/ 目录
go run .
```

访问探针入口：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod http://localhost:8080/health
```

```bash [macOS / Linux]
curl http://localhost:8080/health
```

:::

成功时返回：

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

如果数据库或 Redis 不可用，会返回统一错误结构，例如：

```json
{
  "code": 50300,
  "message": "Redis 不可用"
}
```

## 这一节最该记住的结论

1. 当前入口不是 `cmd/server/main.go`，而是 `server/main.go`。
2. 当前启动编排不是“入口里自己 new 一切”，而是 `bootstrap.MustRun -> bootstrap.NewRouter`。
3. 当前路由聚合不是在入口里手写，而是按模块 `RegisterRoutes` 聚合。
4. 当前健康检查不是内联匿名函数，而是 `system/health_handler.go` 里的独立 Handler。
5. 当前公开探针入口是 `/health`；`/api/v1/system/health` 属于系统模块路由体系的一部分。

下一节继续回看前面几页时，如果看到“在 `main.go` 里直接注册 `/health`”这类说法，就要按这一页的主线来理解并校正。
