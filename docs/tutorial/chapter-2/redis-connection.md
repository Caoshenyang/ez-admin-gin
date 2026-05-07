---
title: Redis 连接
description: "接入 Redis，并按当前启动链把 Redis 状态接入健康检查。"
---

# Redis 连接

数据库已经接进后端服务了。现在继续把 Redis 接进来，让后端具备缓存、会话、验证码和限流这类运行期能力的基础入口。

::: tip 本页先抓主线
当前 Redis 客户端不是在 `main.go` 里直接创建，也不是 `/health` 匿名函数里临时判断；它由 `bootstrap/run.go` 初始化，由 `platform/redis` 提供连接能力，由 `system/health_handler.go` 在健康检查时执行探活。
:::

## Redis 相关代码放在哪里

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
      └─ redis/
         └─ redis.go
```

| 位置 | 当前职责 |
| --- | --- |
| `configs/config.yaml` | Redis 地址、密码、DB 和连接池配置 |
| `internal/platform/redis/redis.go` | 创建和关闭 Redis 客户端 |
| `internal/bootstrap/run.go` | 启动时初始化 Redis |
| `internal/module/system/health_handler.go` | 健康检查时探活 Redis |

## Redis 配置

当前 Redis 配置通常包括：

```yaml
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  max_retries: 3
  min_idle_conns: 5
  pool_size: 10
```

## 启动时如何接入 Redis

当前 Redis 初始化位于 `internal/bootstrap/run.go`：

```go
redisClient, err := platformRedis.New(cfg.Redis, log)
if err != nil {
	log.Fatal("connect redis", zap.Error(err))
}
defer func() {
	if err := platformRedis.Close(redisClient); err != nil {
		log.Error("close redis", zap.Error(err))
	}
}()
```

这个位置很重要，因为它说明 Redis 已经和数据库一样，成为启动阶段的正式依赖，而不是后面某个接口第一次访问时才懒加载。

## 健康检查如何感知 Redis 状态

当前 Redis 健康检查在 `server/internal/module/system/health_handler.go`：

```go
if err := appRedis.Ping(h.redisClient); err != nil {
	h.log.Error("redis health check failed", zap.Error(err))
	response.Error(c, apperror.ServiceUnavailable("Redis 不可用", err), h.log)
	return
}
```

也就是说：

- Redis 是否可用，会直接影响健康检查结果
- `/health` 依然由独立 Handler 统一返回结果
- 不需要再在入口文件里内联一段“数据库成功后再判断 Redis”的匿名路由逻辑

## ✅ 启动并验证

确认 Redis 正在运行：

```bash
docker compose -f deploy/compose.local.yml ps
```

在 `server/` 目录启动服务：

```bash
go run .
```

成功时控制台应看到类似日志：

```text
INFO	redis connected	{"host": "localhost", "port": 6379, "db": 0}
INFO	server started	{"addr": ":8080", "env": "dev"}
```

访问健康检查：

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

1. 当前 Redis 实现目录是 `internal/platform/redis`。
2. 当前 Redis 初始化发生在 `bootstrap/run.go`。
3. 当前 Redis 探活发生在 `system/health_handler.go`。
4. 当前健康检查不再依赖入口文件里的旧式匿名路由示例。

下一节继续统一响应与错误处理：[统一响应与错误处理](./response-and-errors)。
