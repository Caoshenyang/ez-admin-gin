---
title: 数据库连接
description: "使用 GORM 连接 PostgreSQL，并按当前启动链把数据库状态接入健康检查。"
---

# 数据库连接

前面已经准备好了 PostgreSQL 容器、配置和日志系统。现在把数据库连接接进后端服务，让服务启动时完成数据库初始化、迁移执行，并把数据库状态纳入健康检查。

::: tip 本页先抓主线
当前数据库连接不是在 `main.go` 里直接创建，也不是由 `/health` 匿名函数临时探活；它由 `bootstrap/run.go` 初始化，由 `platform/database` 提供连接与迁移 DSN，由 `system/health_handler.go` 在健康检查时执行 `Ping`。
:::

## 数据库相关代码放在哪里

```text
server/
├─ configs/
│  └─ config.yaml
├─ main.go
└─ internal/
   ├─ bootstrap/
   │  └─ run.go
   ├─ module/system/
   │  └─ health_handler.go
   └─ platform/
      └─ database/
         └─ database.go
```

| 位置 | 当前职责 |
| --- | --- |
| `configs/config.yaml` | 数据库驱动、地址、账号和连接池配置 |
| `internal/platform/database/database.go` | 创建连接、关闭连接、生成迁移 DSN |
| `internal/bootstrap/run.go` | 启动时连接数据库并执行迁移 |
| `internal/module/system/health_handler.go` | 健康检查时探活数据库 |

## 数据库配置

当前数据库配置至少包括：

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  user: ez_admin
  password: ez_admin_123456
  name: ez_admin
  max_idle_conns: 10
  max_open_conns: 50
  conn_max_lifetime: 3600
```

其中 `driver` 现在已经是正式配置项，当前代码支持 `postgres` 和 `mysql`。

## 启动时如何接入数据库

当前数据库初始化位于 `internal/bootstrap/run.go`：

```go
db, err := platformDatabase.New(cfg.Database, log)
if err != nil {
	log.Fatal("connect database", zap.Error(err))
}
defer func() {
	if err := platformDatabase.Close(db); err != nil {
		log.Error("close database", zap.Error(err))
	}
}()

migrateDSN, err := platformDatabase.MigrateDSN(cfg.Database)
if err != nil {
	log.Fatal("build migration dsn", zap.Error(err))
}
if err := platformMigrate.Run(cfg.Database.Driver, migrateDSN, migrationsFS, log); err != nil {
	log.Fatal("run database migrations", zap.Error(err))
}
```

这一段说明了当前数据库主线的三个动作：

1. 先连库
2. 再生成迁移 DSN
3. 再执行迁移

所以文档里不应再把“数据库连接”和“迁移执行”拆成完全独立、互不相干的入口故事。

## 健康检查如何感知数据库状态

当前数据库健康检查在 `server/internal/module/system/health_handler.go`：

```go
if err := database.Ping(h.db); err != nil {
	h.log.Error("database health check failed", zap.Error(err))
	response.Error(c, apperror.ServiceUnavailable("数据库不可用", err), h.log)
	return
}
```

也就是说，现在 `/health` 的数据库状态来自独立 Handler，而不是 `main.go` 中的内联匿名函数。

## ✅ 启动并验证

确认 PostgreSQL 正在运行：

```bash
docker compose -f deploy/compose.local.yml ps
```

在 `server/` 目录启动服务：

```bash
go run .
```

成功时控制台应该看到类似日志：

```text
INFO	database connected	{"driver": "postgres", "host": "localhost", "port": 5432, "database": "ez_admin"}
INFO	server started	{"addr": ":8080", "env": "dev"}
```

访问健康检查接口：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod http://localhost:8080/health
```

```bash [macOS / Linux]
curl http://localhost:8080/health
```

:::

返回结构应类似：

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

## 这一页需要同步的认知

1. 当前数据库实现目录是 `internal/platform/database`。
2. 当前数据库初始化发生在 `bootstrap/run.go`，不是 `main.go` 直接 new。
3. 当前迁移执行已并入启动链，不要再把它写成独立的旧入口流程。
4. 当前健康检查的数据库探活发生在 `system/health_handler.go`，不是入口里的匿名 `/health`。

下一节继续接入 Redis：[Redis 连接](./redis-connection)。
