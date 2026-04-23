---
title: 认证中间件
description: "在请求进入业务接口前解析 access_token，并把当前用户信息写入 Gin 上下文。"
---

# 认证中间件

上一节登录接口已经能返回 `access_token`。这一节把 Token 校验前移到中间件中：请求进入受保护接口前，先从 `Authorization` 请求头中解析 Token，校验通过后再继续执行后续 Handler。

::: tip 🎯 本节目标
完成后，访问 `/api/v1/auth/me` 必须携带有效 Token；不带 Token、Token 格式错误或 Token 无效都会返回 `401`。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ internal/
│  ├─ handler/
│  │  └─ auth/
│  │     └─ me.go
│  ├─ middleware/
│  │  └─ auth.go
│  └─ router/
│     └─ router.go
```

| 位置 | 用途 |
| --- | --- |
| `internal/middleware/auth.go` | 解析 `Authorization` 请求头，校验 Token，并写入当前用户信息 |
| `internal/handler/auth/me.go` | 增加一个受保护接口，用来验证中间件是否生效 |
| `internal/router/router.go` | 给 `/api/v1/auth/me` 挂载认证中间件 |

## Header 约定

后续所有需要登录的接口，都使用下面这个请求头传递 Token：

```http
Authorization: Bearer <access_token>
```

其中：

| 部分 | 说明 |
| --- | --- |
| `Authorization` | HTTP 请求头名称 |
| `Bearer` | Token 类型，和登录接口返回的 `token_type` 对应 |
| `<access_token>` | 登录接口返回的 `access_token` |

::: warning ⚠️ `Bearer` 后面有一个空格
正确写法是 `Bearer <access_token>`。如果写成 `Bearer<access_token>`，或者只传 Token 不带 `Bearer`，中间件都会按未登录处理。
:::

## 🛠️ 创建认证中间件

创建 `server/internal/middleware/auth.go`。这是新增文件，直接完整写入即可。

```go
package middleware

import (
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/response"
	"ez-admin-gin/server/internal/token"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	currentUserIDKey   = "current_user_id"
	currentUsernameKey = "current_username"
)

// Auth 校验 Authorization 请求头，并把当前用户信息写入 Gin 上下文。
func Auth(tokenManager *token.Manager, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			response.Error(c, apperror.Unauthorized("请先登录"), log)
			c.Abort()
			return
		}

		claims, err := tokenManager.ParseAccessToken(tokenString)
		if err != nil {
			if log != nil {
				log.Warn("parse access token failed", zap.Error(err))
			}

			response.Error(c, apperror.Unauthorized("登录已过期，请重新登录"), log)
			c.Abort()
			return
		}

		// 后续 Handler 可以从 Gin 上下文中取当前用户信息。
		c.Set(currentUserIDKey, claims.UserID)
		c.Set(currentUsernameKey, claims.Username)
		c.Next()
	}
}

// CurrentUserID 从 Gin 上下文中取当前用户 ID。
func CurrentUserID(c *gin.Context) (uint, bool) {
	value, ok := c.Get(currentUserIDKey)
	if !ok {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}

// CurrentUsername 从 Gin 上下文中取当前用户名。
func CurrentUsername(c *gin.Context) (string, bool) {
	value, ok := c.Get(currentUsernameKey)
	if !ok {
		return "", false
	}

	username, ok := value.(string)
	return username, ok
}

// bearerToken 解析 Authorization: Bearer <token>。
func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 {
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	if strings.TrimSpace(parts[1]) == "" {
		return "", false
	}

	return parts[1], true
}
```

::: details 为什么用 `strings.Fields`
`strings.Fields` 会自动处理多余空格。比如 `Bearer   token` 也能被拆成两段，比手动按一个空格切分更稳。
:::

::: details 为什么把用户信息写进 Gin 上下文
Token 校验通过后，后续 Handler 不应该重复解析 Token。中间件把 `user_id`、`username` 写进上下文，后续接口就可以直接读取当前用户信息。
:::

## 🛠️ 创建当前用户接口

创建 `server/internal/handler/auth/me.go`。这是新增文件，直接完整写入即可。

```go
package auth

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MeHandler 负责当前用户相关接口。
type MeHandler struct {
	log *zap.Logger
}

// NewMeHandler 创建当前用户 Handler。
func NewMeHandler(log *zap.Logger) *MeHandler {
	return &MeHandler{
		log: log,
	}
}

type meResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// Me 返回当前登录用户的基础信息。
func (h *MeHandler) Me(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), h.log)
		return
	}

	username, _ := middleware.CurrentUsername(c)

	response.Success(c, meResponse{
		UserID:   userID,
		Username: username,
	})
}
```

::: details 为什么这里先不查数据库
本节重点是验证“Token 能被中间件解析，并把当前用户信息传给 Handler”。所以 `/api/v1/auth/me` 先返回 Token 里的基础信息。

后续进入用户管理、角色权限后，再根据需要查询数据库，补充昵称、角色、菜单等信息。
:::

## 🛠️ 注册受保护路由

修改 `server/internal/router/router.go`。这一处重点看三个变化：

- 新增 `internal/middleware` import。
- 创建 `MeHandler`。
- 给 `/api/v1/auth/me` 挂载认证中间件。

先调整 import：

```go
import (
	"ez-admin-gin/server/internal/config"
	authHandler "ez-admin-gin/server/internal/handler/auth"
	systemHandler "ez-admin-gin/server/internal/handler/system"
	appLogger "ez-admin-gin/server/internal/logger"
	"ez-admin-gin/server/internal/middleware" // [!code ++]
	"ez-admin-gin/server/internal/token"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)
```

更新 `registerAuthRoutes`：

```go
// registerAuthRoutes 注册认证相关路由。
func registerAuthRoutes(r *gin.Engine, opts Options) {
	login := authHandler.NewLoginHandler(opts.DB, opts.Log, opts.Token)
	me := authHandler.NewMeHandler(opts.Log) // [!code ++]

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login", login.Login)

	protectedAuth := auth.Group("") // [!code ++]
	protectedAuth.Use(middleware.Auth(opts.Token, opts.Log)) // [!code ++]
	protectedAuth.GET("/me", me.Me) // [!code ++]
}
```

::: details 为什么不把 `/login` 放进中间件
登录接口本身就是用来获取 Token 的。如果登录接口也要求先携带 Token，就会变成“必须先登录才能登录”。

所以 `/api/v1/auth/login` 是公开接口，`/api/v1/auth/me` 才是受保护接口。
:::

## ✅ 启动服务

确认数据库和 Redis 正在运行：

```bash
# 在项目根目录执行，确认本地依赖服务处于运行状态
docker compose -f deploy/compose.local.yml ps
```

进入 `server/` 目录启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

## ✅ 验证不带 Token 会失败

打开另一个终端，请求当前用户接口：

::: code-group

```powershell [Windows PowerShell]
# 不携带 Authorization 请求头
try {
  Invoke-RestMethod `
    -Method Get `
    -Uri http://localhost:8080/api/v1/auth/me
} catch {
  $_.ErrorDetails.Message
}
```

```bash [macOS / Linux]
# 不携带 Authorization 请求头
curl -i http://localhost:8080/api/v1/auth/me
```

:::

应该能看到 HTTP 状态码为 `401`，响应体类似：

```json
{
  "code": 40100,
  "message": "请先登录"
}
```

## ✅ 验证携带 Token 会成功

先登录拿到 Token，再请求 `/api/v1/auth/me`。

::: code-group

```powershell [Windows PowerShell]
# 登录并提取 access_token
$body = @{
  username = "admin"
  password = "EzAdmin@123456"
} | ConvertTo-Json

$login = Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/auth/login `
  -ContentType "application/json" `
  -Body $body

$token = $login.data.access_token

# 携带 Token 请求受保护接口
Invoke-RestMethod `
  -Method Get `
  -Uri http://localhost:8080/api/v1/auth/me `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
# 登录并提取 access_token，需要安装 jq
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"EzAdmin@123456"}' | jq -r '.data.access_token')

# 携带 Token 请求受保护接口
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该看到类似结果：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "user_id": 1,
    "username": "admin"
  }
}
```

::: info 为什么这里验证的是 `/auth/me`
`/auth/me` 足够小，只依赖认证中间件，不掺入角色、权限、菜单等后续内容。先用它确认 Token 能正确识别当前用户，下一步再继续做权限控制。
:::

## ✅ 验证错误 Token 会失败

继续请求当前用户接口，但故意传入错误 Token：

::: code-group

```powershell [Windows PowerShell]
# 携带错误 Token
try {
  Invoke-RestMethod `
    -Method Get `
    -Uri http://localhost:8080/api/v1/auth/me `
    -Headers @{ Authorization = "Bearer wrong-token" }
} catch {
  $_.ErrorDetails.Message
}
```

```bash [macOS / Linux]
# 携带错误 Token
curl -i -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer wrong-token"
```

:::

应该能看到 HTTP 状态码为 `401`，响应体类似：

```json
{
  "code": 40100,
  "message": "登录已过期，请重新登录"
}
```

## 常见问题

::: details 明明登录成功了，访问 `/auth/me` 还是提示 `请先登录`
优先检查 `Authorization` 请求头格式：

```http
Authorization: Bearer <access_token>
```

常见错误是漏写 `Bearer`、`Bearer` 后面没有空格，或者把整个登录响应当成 Token 传过去。
:::

::: details 提示 `登录已过期，请重新登录`
这表示请求头里有 Token，但解析或校验失败。常见原因：

- Token 被复制截断。
- 使用了旧服务签发的 Token，但重启后改过 `jwt_secret`。
- Token 已经过期。
- 请求头格式里混入了多余字符。
:::

::: details `CurrentUserID` 为什么返回 `(uint, bool)`
上下文里不一定存在当前用户信息。例如没有挂载认证中间件，或者中间件提前失败。返回 `bool` 可以让 Handler 明确判断“有没有取到当前用户”，避免直接类型断言导致 panic。
:::

下一节会继续设计角色、权限和关联表：[角色与权限模型](./rbac-model)。
