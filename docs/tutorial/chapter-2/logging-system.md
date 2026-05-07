---
title: 日志系统
description: "接入 Zap 日志，并按当前启动链和路由聚合方式完成请求日志挂载。"
---

# 日志系统

上一节已经把运行参数交给配置管理。这一节继续把结构化日志接进来，让服务启动、请求访问和异常恢复都有统一记录。

::: tip 本页先抓一件事
当前日志系统不是在 `main.go` 里直接 `gin.New()` 后手动拼装，而是由 `bootstrap/run.go` 创建日志对象，再由 `bootstrap/router.go` 把请求日志和恢复中间件挂到 Gin 引擎上。
:::

## 日志相关代码放在哪里

```text
server/
├─ configs/
│  └─ config.yaml
├─ main.go
└─ internal/
   ├─ bootstrap/
   │  ├─ run.go
   │  └─ router.go
   └─ platform/
      └─ logger/
         └─ logger.go
```

| 位置 | 用途 |
| --- | --- |
| `configs/config.yaml` | 定义日志级别、格式和文件位置 |
| `internal/platform/logger/logger.go` | 初始化 Zap，并提供 Gin 中间件 |
| `internal/bootstrap/run.go` | 启动时创建日志对象 |
| `internal/bootstrap/router.go` | 挂载 `GinLogger` 与 `GinRecovery` |

## 日志配置

当前配置文件里的日志段通常长这样：

```yaml
log:
  level: info
  format: console
  filename: logs/app.log
  max_size: 100
  max_backups: 7
  max_age: 30
  compress: false
```

字段含义和前面版本类似，但当前主线目录应理解为 `platform/logger`，而不是继续围绕旧的 `internal/logger` 教程结构思考。

## 日志对象在哪里创建

当前启动时由 `internal/bootstrap/run.go` 负责创建日志对象：

```go
log, err := platformLogger.New(cfg.Log)
if err != nil {
	stdlog.Fatalf("create logger: %v", err)
}
defer func() {
	_ = log.Sync()
}()
```

这一步仍然发生在数据库、Redis、鉴权初始化之前，因为后续所有失败场景都要依赖日志输出。

## 请求日志和恢复中间件挂在哪里

当前不是在入口文件里直接写：

```go
r := gin.New()
r.Use(appLogger.GinLogger(log), appLogger.GinRecovery(log))
```

而是由 `internal/bootstrap/router.go` 统一完成：

```go
func NewRouter(opts RouterOptions) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	// 省略上传、静态文件和模块路由注册
	return r
}
```

这样做的好处是：日志中间件和路由引擎的生命周期保持一致，不会再出现“教程里在 `main.go` 手写路由、后面又迁到别处”的结构漂移。

## `platform/logger` 负责什么

当前 `server/internal/platform/logger/logger.go` 的职责可以概括成三件事：

- 创建 `*zap.Logger`
- 提供 `GinLogger`
- 提供 `GinRecovery`

也就是说：

- 启动日志来自 `bootstrap/run.go`
- 请求日志和 panic 恢复来自 `bootstrap/router.go`
- 具体实现细节落在 `platform/logger`

## ✅ 启动并验证

在 `server/` 目录执行：

```bash
go run .
```

启动后控制台应该能看到类似日志：

```text
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

随后应新增一条 `http request` 日志。

## ✅ 验证文件日志

查看日志文件：

::: code-group

```powershell [Windows PowerShell]
Get-Content .\logs\app.log -Tail 5
```

```bash [macOS / Linux]
tail -n 5 logs/app.log
```

:::

如果能看到 `server started` 和 `http request`，说明控制台与文件输出都已经生效。

## 这一页需要同步的认知

1. 当前日志实现目录是 `internal/platform/logger`。
2. 当前日志对象在 `bootstrap/run.go` 中创建。
3. 当前 Gin 日志中间件在 `bootstrap/router.go` 中挂载。
4. 当前启动入口 `server/main.go` 不再自己承载 `gin.New()` 和中间件拼装。

下一节继续连接数据库：[数据库连接](./database-connection)。
