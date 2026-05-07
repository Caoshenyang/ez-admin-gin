---
title: 统一响应与错误处理
description: "定义后台 API 的统一响应格式，并对齐当前 Handler 层的错误处理边界。"
---

# 统一响应与错误处理

前面几节已经让后端能连接数据库和 Redis。现在把接口返回值统一起来：成功响应有固定结构，错误响应有固定错误码，后续管理台调用接口时就不用为每个接口猜格式。

::: tip 本页先抓主线
当前统一响应已经不再通过 `main.go` 里的匿名 `/health` 示例来承载，而是由 `response.Success` / `response.Error` 负责协议输出，由 `system/health_handler.go` 和各模块 Handler 负责调用；业务 Handler 里的共性错误处理则进一步沉淀在 `internal/app/*`。
:::

## 相关代码放在哪里

```text
server/
└─ internal/
   ├─ app/
   │  ├─ handler.go
   │  └─ pagination.go
   ├─ apperror/
   │  └─ apperror.go
   ├─ module/system/
   │  └─ health_handler.go
   └─ response/
      └─ response.go
```

| 位置 | 当前职责 |
| --- | --- |
| `internal/apperror/apperror.go` | 定义应用错误、错误码和 HTTP 状态码 |
| `internal/response/response.go` | 输出统一成功响应和错误响应 |
| `internal/module/system/health_handler.go` | 在健康检查里实际使用统一响应 |
| `internal/app/handler.go` | 给业务 Handler 提供统一错误写回、登录人提取、ID 参数解析等助手 |

## 响应格式

成功响应统一为：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "example": "value"
  }
}
```

错误响应统一为：

```json
{
  "code": 50300,
  "message": "Redis 不可用"
}
```

## 当前健康检查如何使用统一响应

当前不是在 `main.go` 里直接 `c.JSON(...)`，而是由 `system/health_handler.go` 负责：

```go
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
```

这意味着当前统一响应的真实落点已经从“教程演示版入口文件”切换成“模块 Handler + 公共响应包”。

## `internal/app/*` 在错误处理里的角色

随着模块变多，当前仓库已经把 Handler 层的共性逻辑收进了 `internal/app`：

```go
func WriteError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
```

还有另外两个很常用的助手：

- `CurrentActor`：从 `gin.Context` 提取当前登录人
- `UintIDParam`：统一解析路径里的 `uint` ID

也就是说，当前业务 Handler 的推荐写法已经从“每个模块手写一遍参数校验和错误分支”，变成了“优先复用 `app` 层助手 + `response` / `apperror` 协议层”。

## Handler 层以后怎么写

后续业务接口可以优先沿用这套边界：

```go
if username == "" {
	response.Error(c, apperror.BadRequest("用户名不能为空"), log)
	return
}

response.Success(c, gin.H{
	"id":       1,
	"username": username,
})
```

如果错误来自 service / repository，且可能既有应用错误也有底层错误，更适合用：

```go
app.WriteError(c, err, "创建用户失败", log)
```

## ✅ 验证成功响应

在 `server/` 目录启动服务：

```bash
go run .
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

应看到：

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

## ✅ 验证错误响应

保持服务运行，停掉 Redis 后再次访问 `/health`，应看到类似：

```json
{
  "code": 50300,
  "message": "Redis 不可用"
}
```

## 这一页需要同步的认知

1. 当前统一响应的真实使用场景是模块 Handler，不是入口文件里的旧示例。
2. 当前健康检查响应由 `system/health_handler.go` 输出。
3. 当前 Handler 共性逻辑已经沉淀到 `internal/app/*`。
4. 后续新增模块时，优先复用 `app`、`apperror`、`response` 这一套，不要回到旧式内联写法。

下一节继续看路由聚合与健康检查主线：[路由分组与健康检查](./routing-and-health)。
