---
title: 统一响应与错误处理
description: "定义后台 API 的统一响应格式、错误码和错误处理边界。"
---

# 统一响应与错误处理

前面几节已经让后端能连接数据库和 Redis。现在把接口返回值统一起来：成功响应有固定结构，错误响应有固定错误码，后续管理台调用接口时就不用为每个接口单独猜格式。

::: tip 🎯 本节目标
完成后，`/health` 会返回统一响应格式；当数据库或 Redis 不可用时，也会返回统一错误结构。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ internal/
│  ├─ apperror/
│  │  └─ apperror.go
│  └─ response/
│     └─ response.go
└─ main.go
```

| 位置 | 用途 |
| --- | --- |
| `internal/apperror/apperror.go` | 定义应用错误、错误码和 HTTP 状态码 |
| `internal/response/response.go` | 定义统一成功响应和错误响应 |
| `main.go` | 把 `/health` 改成统一响应格式 |

::: info 本节不需要安装新依赖
这一节只使用标准库、Gin 和已经接入的 Zap，不需要执行 `go get`。
:::

## 响应格式长什么样

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

字段说明：

| 字段 | 说明 |
| --- | --- |
| `code` | 业务错误码。`0` 表示成功 |
| `message` | 给前端展示或排查用的简短信息 |
| `data` | 成功时返回的数据，错误时默认不返回 |

::: details 为什么第一版先保持简洁？
统一响应结构可以继续扩展，但第一版先保持 `code`、`message`、`data` 更合适。

随着功能增加，后续可以继续补充 `request_id`、`path`、`timestamp` 这类字段，用来关联日志、定位请求路径和记录错误时间。例如：

```json
{
  "code": 50300,
  "message": "Redis 不可用",
  "request_id": "01HV9Z3YJ8K2P7Q0M4R6T1A2BC",
  "path": "/health",
  "timestamp": "2026-04-22T10:30:00+08:00"
}
```

这一节先不加这些字段，是因为它们通常依赖请求追踪中间件、日志链路和前端拦截器一起配合。等后面接入全局中间件和认证流程时，再扩展响应结构会更自然。

无论字段多少，有一条边界始终不变：返回给前端的是可理解、可展示的信息；底层真实错误留在日志里排查。
:::

::: details 为什么不直接只看 HTTP 状态码
HTTP 状态码负责表达请求层面的成功或失败，例如 `200`、`400`、`503`。业务错误码负责表达应用内更细的错误类型，例如参数错误、权限不足、依赖服务不可用。

管理台前端通常会同时使用两者：HTTP 状态码用于拦截器判断请求是否成功，业务错误码用于展示更准确的提示。
::: 

## 错误码约定

先定义一组够用的基础错误码：

| 错误码 | 含义 | HTTP 状态码 |
| --- | --- | --- |
| `0` | 成功 | `200` |
| `40000` | 请求参数错误 | `400` |
| `40100` | 未登录或登录已过期 | `401` |
| `40300` | 没有权限 | `403` |
| `40400` | 资源不存在 | `404` |
| `50300` | 依赖服务不可用 | `503` |
| `50000` | 服务器内部错误 | `500` |

::: warning ⚠️ 错误码先少一点
这一节只定义基础错误码。后续做到登录、权限、菜单、文件上传时，再按模块补充更细的业务错误码。
:::

## 🛠️ 创建应用错误包

创建 `server/internal/apperror/apperror.go`。这是新增文件，直接完整写入即可。

```go
package apperror

import (
	"fmt"
	"net/http"
)

// Code 是业务响应码类型。
// Go 没有 Java 那样的 enum，通常用“自定义类型 + const 常量”表达枚举语义。
type Code int

const (
	// CodeSuccess 表示请求处理成功。
	CodeSuccess Code = 0
	// CodeBadRequest 表示请求参数错误。
	CodeBadRequest Code = 40000
	// CodeUnauthorized 表示未登录或登录已过期。
	CodeUnauthorized Code = 40100
	// CodeForbidden 表示没有权限访问资源。
	CodeForbidden Code = 40300
	// CodeNotFound 表示资源不存在。
	CodeNotFound Code = 40400
	// CodeServiceUnavailable 表示数据库、Redis 等依赖服务不可用。
	CodeServiceUnavailable Code = 50300
	// CodeInternal 表示服务器内部错误。
	CodeInternal Code = 50000
)

// Error 表示可以安全返回给前端的应用错误。
type Error struct {
	Code    Code
	Message string
	Status  int
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// New 创建一个不包裹底层错误的应用错误。
func New(status int, code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Wrap 创建一个包裹底层错误的应用错误。
func Wrap(err error, status int, code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

func BadRequest(message string) *Error {
	return New(http.StatusBadRequest, CodeBadRequest, message)
}

func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(message string) *Error {
	return New(http.StatusForbidden, CodeForbidden, message)
}

func NotFound(message string) *Error {
	return New(http.StatusNotFound, CodeNotFound, message)
}

func ServiceUnavailable(message string, err error) *Error {
	return Wrap(err, http.StatusServiceUnavailable, CodeServiceUnavailable, message)
}

func Internal(message string, err error) *Error {
	return Wrap(err, http.StatusInternalServerError, CodeInternal, message)
}
```

这个包负责把“给前端看的错误”和“底层真实错误”分开：

- `Message`：可以返回给前端。
- `Err`：保留给日志排查。
- `Status`：决定 HTTP 状态码。
- `Code`：决定业务错误码。

## 🛠️ 创建统一响应包

创建 `server/internal/response/response.go`。这是新增文件，直接完整写入即可。

```go
package response

import (
	"errors"
	"net/http"

	"ez-admin-gin/server/internal/apperror"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Body struct {
	Code    apperror.Code `json:"code"`
	Message string        `json:"message"`
	Data    any           `json:"data,omitempty"`
}

// Success 返回统一成功响应。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:    apperror.CodeSuccess,
		Message: "ok",
		Data:    data,
	})
}

// Error 返回统一错误响应。
func Error(c *gin.Context, err error, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, Body{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	// 未归类错误不把内部细节返回给前端，只记录到日志里。
	if log != nil {
		log.Error("unhandled error", zap.Error(err))
	}

	c.JSON(http.StatusInternalServerError, Body{
		Code:    apperror.CodeInternal,
		Message: "服务器内部错误",
	})
}
```

::: details 为什么 `data` 使用 `omitempty`
错误响应通常不需要 `data` 字段。加上 `omitempty` 后，`data` 为空时不会出现在 JSON 里，响应会更干净。
:::

## 🛠️ 改造健康检查响应

修改 `server/main.go`。这一处重点看三个变化：

- 引入 `internal/apperror` 和 `internal/response`。
- 数据库或 Redis 不可用时，使用 `response.Error` 返回统一错误。
- 成功时，使用 `response.Success` 返回统一成功结构。

```go
package main

import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/apperror" // [!code ++]
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/response" // [!code ++]

	"github.com/gin-gonic/gin"
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

	// 使用 gin.New()，再手动挂载自定义中间件。
	r := gin.New()
	r.Use(appLogger.GinLogger(log), appLogger.GinRecovery(log))

	r.GET("/health", func(c *gin.Context) {
		if err := database.Ping(db); err != nil {
			log.Error("database health check failed", zap.Error(err))
			response.Error(c, apperror.ServiceUnavailable("数据库不可用", err), log) // [!code ++]
			return
		}

		if err := appRedis.Ping(redisClient); err != nil {
			log.Error("redis health check failed", zap.Error(err))
			response.Error(c, apperror.ServiceUnavailable("Redis 不可用", err), log) // [!code ++]
			return
		}

		response.Success(c, gin.H{ // [!code ++]
			"env":      cfg.App.Env,
			"database": "ok",
			"redis":    "ok",
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

::: warning ⚠️ 健康检查返回结构会变化
从这一节开始，`/health` 不再直接返回 `status`、`database`、`redis`，而是统一包在 `code`、`message`、`data` 里。
:::

## ✅ 启动并验证成功响应

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

## ✅ 验证错误响应

保持后端服务运行，打开另一个终端停止 Redis：

```bash
# 在项目根目录停止 Redis 容器
docker compose -f deploy/compose.local.yml stop redis
```

再次访问健康检查接口：

::: code-group

```powershell [Windows PowerShell]
# 非 2xx 响应会进入 catch，这里直接打印响应体
try {
  Invoke-RestMethod http://localhost:8080/health
} catch {
  $_.ErrorDetails.Message
}
```

```bash [macOS / Linux]
# -i 可以同时查看 HTTP 状态码和响应体
curl -i http://localhost:8080/health
```

:::

应该能看到 HTTP 状态码为 `503`，响应体类似：

```json
{
  "code": 50300,
  "message": "Redis 不可用"
}
```

验证完成后，重新启动 Redis：

```bash
# 在项目根目录重新启动 Redis 容器
docker compose -f deploy/compose.local.yml start redis
```

::: warning ⚠️ 验证失败场景后记得恢复服务
后续章节还会继续依赖 Redis。验证错误响应后，记得重新执行 `start redis`，并确认容器状态恢复正常。
:::

## Handler 层以后怎么用

后续写业务接口时，可以按这个边界处理：

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

这条规则先记住：

- 参数不合法：返回 `apperror.BadRequest`
- 没登录：返回 `apperror.Unauthorized`
- 没权限：返回 `apperror.Forbidden`
- 找不到资源：返回 `apperror.NotFound`
- 数据库、Redis 等依赖不可用：返回 `apperror.ServiceUnavailable`
- 没有预料到的内部错误：返回 `apperror.Internal`

下一节开始整理路由结构：[路由分组与健康检查](./routing-and-health)。
