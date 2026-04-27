---
title: Casbin 权限控制
description: "接入 Casbin，用角色编码判断接口访问权限。"
---

# Casbin 权限控制

前面已经完成登录、Token 认证和用户角色关系。这一节接入 Casbin，把“某个角色能不能访问某个接口”交给策略表维护。

::: tip 🎯 本节目标
完成后，`/health` 仍然是公开探针接口；`/api/v1/system/health` 会变成受保护接口，需要登录，并且当前用户角色必须有对应权限策略。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ configs/
│  └─ rbac_model.conf
├─ internal/
│  ├─ middleware/
│  │  └─ permission.go
│  ├─ model/
│  │  └─ casbin_rule.go
│  ├─ permission/
│  │  └─ enforcer.go
│  └─ router/
│     └─ router.go
├─ main.go
├─ go.mod
└─ go.sum
```

| 位置 | 用途 |
| --- | --- |
| `configs/rbac_model.conf` | 定义 Casbin 的请求、策略和匹配规则 |
| `internal/model/casbin_rule.go` | 定义 Casbin 策略表结构，供初始化使用 |
| `internal/permission/enforcer.go` | 创建 Casbin Enforcer，并关闭自动建表 |
| `internal/middleware/permission.go` | 根据当前用户角色判断接口权限 |
| `internal/router/router.go` | 给受保护接口挂载权限中间件 |
| `main.go` | 创建权限 Enforcer 并传给路由 |

## 权限判断方式

本节先使用角色编码作为 Casbin 的主体：

```text
sub = 角色编码，例如 super_admin
obj = 接口路径，例如 /api/v1/system/health
act = 请求方法，例如 GET
```

默认策略会长这样：

```text
p, super_admin, /api/v1/system/health, GET
```

含义是：`super_admin` 角色可以用 `GET` 访问 `/api/v1/system/health`。

::: warning ⚠️ Casbin 只做判断，不替代业务校验
Casbin 负责回答“这个角色能不能访问这个接口”。用户是否存在、角色是否启用、角色是否绑定到用户，仍然由数据库和业务逻辑处理。
:::

## 🛠️ 安装 Casbin 依赖

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

安装 Casbin 和 GORM 适配器：

```bash
# 安装 Casbin 与 GORM 策略存储适配器
go get github.com/casbin/casbin/v3@latest
go get github.com/casbin/gorm-adapter/v3@latest
```

依赖资料入口：

| 依赖 | 用途 | 资料 |
| --- | --- | --- |
| `github.com/casbin/casbin/v3` | 权限模型和策略判断 | [Go 包文档](https://pkg.go.dev/github.com/casbin/casbin/v3) |
| `github.com/casbin/gorm-adapter/v3` | 从数据库加载 Casbin 策略 | [Go 包文档](https://pkg.go.dev/github.com/casbin/gorm-adapter/v3) |

::: warning ⚠️ 继续使用 SQL 建表
`gorm-adapter` 默认会尝试自动建表。本节会在代码中关闭它的自动迁移能力，表结构仍然以参考手册中的 SQL 为准。
:::

## 先创建数据表

本节新增 `casbin_rule`，用于保存 Casbin 接口权限策略。

::: tip 建表 SQL
字段说明、表名约定、唯一索引和 PostgreSQL / MySQL 建表语句统一放在参考手册：[数据库建表语句 - `casbin_rule`](../../reference/database-ddl#casbin-rule)。
:::

## 🛠️ 创建 Casbin 模型文件

创建 `server/configs/rbac_model.conf`。这是新增文件，直接完整写入即可。

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

配置含义：

| 配置 | 说明 |
| --- | --- |
| `r = sub, obj, act` | 请求由角色编码、接口路径、请求方法组成 |
| `p = sub, obj, act` | 策略也由角色编码、接口路径、请求方法组成 |
| `keyMatch2` | 支持路径匹配，后续可以匹配 `/api/v1/users/:id` |
| `act == "*"` | 允许某条策略匹配全部 HTTP 方法 |

::: details 为什么没有在 Casbin 里写用户和角色关系
用户和角色关系已经放在 `sys_user_role` 里。权限判断时，中间件先根据当前用户查出角色编码，再把角色编码交给 Casbin 判断。

这样可以避免在 `sys_user_role` 和 Casbin `g` 策略里重复维护同一份用户角色关系。
:::

## 🛠️ 创建 Casbin 策略模型

创建 `server/internal/model/casbin_rule.go`。这是新增文件，直接完整写入即可。

```go
package model

// CasbinRule 是 Casbin gorm-adapter 使用的策略表模型。
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Ptype string `gorm:"size:100;not null;default:''" json:"ptype"`
	V0    string `gorm:"size:100;not null;default:''" json:"v0"`
	V1    string `gorm:"size:100;not null;default:''" json:"v1"`
	V2    string `gorm:"size:100;not null;default:''" json:"v2"`
	V3    string `gorm:"size:100;not null;default:''" json:"v3"`
	V4    string `gorm:"size:100;not null;default:''" json:"v4"`
	V5    string `gorm:"size:100;not null;default:''" json:"v5"`
}

// TableName 固定 Casbin 策略表名。
func (CasbinRule) TableName() string {
	return "casbin_rule"
}
```

## 🛠️ 创建 Enforcer

创建 `server/internal/permission/enforcer.go`。这是新增文件，直接完整写入即可。

```go
package permission

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Enforcer 包装 Casbin 权限判断能力。
type Enforcer struct {
	inner *casbin.Enforcer
}

// NewEnforcer 创建权限判断器，并从数据库加载策略。
func NewEnforcer(db *gorm.DB, modelPath string) (*Enforcer, error) {
	// 本项目统一使用 SQL 建表，不让 gorm-adapter 自动迁移表结构。
	gormadapter.TurnOffAutoMigrate(db)

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("create casbin adapter: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("create casbin enforcer: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load casbin policy: %w", err)
	}

	return &Enforcer{
		inner: enforcer,
	}, nil
}

// Enforce 判断角色是否允许访问某个接口。
func (e *Enforcer) Enforce(sub string, obj string, act string) (bool, error) {
	allowed, err := e.inner.Enforce(sub, obj, act)
	if err != nil {
		return false, fmt.Errorf("enforce permission: %w", err)
	}

	return allowed, nil
}
```

::: details 为什么这里还要包装一层
业务代码只需要知道“能不能访问”，不需要到处直接依赖 Casbin 的具体类型。后续如果要加缓存、重新加载策略、日志统计，也可以放在这个包里。
:::

## 🛠️ 创建权限中间件

创建 `server/internal/middleware/permission.go`。这是新增文件，直接完整写入即可。

```go
package middleware

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/permission"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Permission 根据当前用户角色判断接口访问权限。
func Permission(db *gorm.DB, enforcer *permission.Enforcer, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := CurrentUserID(c)
		if !ok {
			response.Error(c, apperror.Unauthorized("请先登录"), log)
			c.Abort()
			return
		}

		roleCodes, err := currentRoleCodes(db, userID)
		if err != nil {
			response.Error(c, apperror.Internal("权限校验失败", err), log)
			c.Abort()
			return
		}

		if len(roleCodes) == 0 {
			response.Error(c, apperror.Forbidden("没有权限访问"), log)
			c.Abort()
			return
		}

		obj := c.FullPath()
		if obj == "" {
			obj = c.Request.URL.Path
		}
		act := c.Request.Method

		for _, roleCode := range roleCodes {
			allowed, err := enforcer.Enforce(roleCode, obj, act)
			if err != nil {
				response.Error(c, apperror.Internal("权限校验失败", err), log)
				c.Abort()
				return
			}

			if allowed {
				c.Next()
				return
			}
		}

		response.Error(c, apperror.Forbidden("没有权限访问"), log)
		c.Abort()
	}
}

// currentRoleCodes 查询当前用户拥有的启用角色编码。
func currentRoleCodes(db *gorm.DB, userID uint) ([]string, error) {
	var roleCodes []string
	err := db.
		Table("sys_role AS r").
		Select("r.code").
		Joins("JOIN sys_user_role AS ur ON ur.role_id = r.id").
		Where("ur.user_id = ?", userID).
		Where("r.status = ?", model.RoleStatusEnabled).
		Where("r.deleted_at IS NULL").
		Pluck("r.code", &roleCodes).Error
	if err != nil {
		return nil, err
	}

	return roleCodes, nil
}
```

::: details 为什么用 `c.FullPath()`
`c.FullPath()` 返回路由注册时的路径。例如后续有 `/api/v1/users/:id`，它会返回带 `:id` 的模板路径，而不是某个具体 ID。

这样策略可以写成一条规则匹配一类接口。
:::

::: tip 📌 默认接口权限初始化
默认接口权限策略通过数据库迁移文件自动创建，不需要在代码中手动初始化。当服务启动时，会执行 `server/migrations/{postgres,mysql}/000002_seed_data.up.sql` 迁移文件，创建超级管理员角色、系统菜单和权限策略。

这样可以确保权限策略在服务启动时就已经准备就绪，不需要通过代码手动写入。
:::

::: warning ⚠️ 策略初始化后需要重新加载
本节的 Enforcer 在服务启动时加载策略。所以新增或修改 `casbin_rule` 后，当前服务进程不会自动感知。

现在先通过重启服务重新加载策略。后续做权限管理接口时，再补“保存策略后重新加载”的流程。
:::

## 🛠️ 在启动入口创建 Enforcer

修改 `server/main.go`。这一处重点看两个变化：

- 引入 `internal/permission`。
- 创建 Enforcer，并传给路由。

先调整 import：

```go
import (
	"ez-admin-gin/server/internal/bootstrap"
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	"ez-admin-gin/server/internal/permission" // [!code ++]
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/router"
	"ez-admin-gin/server/internal/token"

	"go.uber.org/zap"
)
```

在创建路由前，增加权限 Enforcer：

```go
	// 权限判断器负责根据角色策略判断接口访问权限。
	permissionEnforcer, err := permission.NewEnforcer(db, "configs/rbac_model.conf") // [!code ++]
	if err != nil { // [!code ++]
		log.Fatal("create permission enforcer", zap.Error(err)) // [!code ++]
	} // [!code ++]

	// 路由注册交给 internal/router，main.go 只保留启动流程。
	r := router.New(router.Options{
		Config:     cfg,
		Log:        log,
		DB:         db,
		Redis:      redisClient,
		Token:      tokenManager,
		Permission: permissionEnforcer, // [!code ++]
	})
```

## 🛠️ 给接口挂载权限中间件

修改 `server/internal/router/router.go`。这一处重点看两个变化：

- `Options` 增加 `Permission` 字段。
- `/api/v1/system/health` 增加认证和权限校验。

先调整 import：

```go
import (
	"ez-admin-gin/server/internal/config"
	authHandler "ez-admin-gin/server/internal/handler/auth"
	systemHandler "ez-admin-gin/server/internal/handler/system"
	appLogger "ez-admin-gin/server/internal/logger"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/permission" // [!code ++]
	"ez-admin-gin/server/internal/token"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)
```

更新 `Options`：

```go
type Options struct {
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      *goredis.Client
	Token      *token.Manager
	Permission *permission.Enforcer // [!code ++]
}
```

更新 `registerSystemRoutes`：

```go
// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)

	// /health 通常给部署探针和本地快速验证使用，保持公开访问。
	r.GET("/health", health.Check)

	// /api/v1/system/health 作为后台接口，需要登录并通过权限校验。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log)) // [!code ++]
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log)) // [!code ++]
	system.GET("/health", health.Check)
}
```

::: details 为什么 `/health` 仍然公开
部署探针和本地快速检查通常不应该依赖登录状态，所以根路径 `/health` 保持公开。

`/api/v1/system/health` 属于后台接口分组，用它来验证认证和权限链路更合适。
:::

## ✅ 整理依赖并启动

整理依赖：

```bash
# 在 server/ 目录执行
go mod tidy
```

确认数据库和 Redis 正在运行：

```bash
# 在项目根目录执行，确认本地依赖服务处于运行状态
docker compose -f deploy/compose.local.yml ps
```

回到 `server/` 目录启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

第一次启动后，控制台应该能看到类似日志：

```text
INFO	database migrations applied
INFO	server started	{"addr": ":8080", "env": "dev"}
```

## ✅ 创建管理员账号

服务启动后，先通过初始化接口创建管理员账号：

```bash
# 创建管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

## ✅ 验证策略已经写入

打开另一个终端，在项目根目录执行：

```bash
# 查看默认接口权限策略
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ptype, v0, v1, v2 from casbin_rule;"
```

应该看到类似结果，包含系统默认的权限策略：

```text
 ptype |     v0      |           v1            | v2
-------+-------------+-------------------------+-----
 p     | super_admin | /api/v1/system/health   | GET
 p     | super_admin | /api/v1/auth/login      | POST
 p     | super_admin | /api/v1/setup/init      | POST
```

## ✅ 验证公开健康检查仍然可访问

```bash
# 不需要 Token，仍然可以访问
curl -i http://localhost:8080/health
```

应该看到 HTTP 状态码为 `200`。

## ✅ 验证后台健康检查需要 Token 和权限

先不带 Token 请求：

::: code-group

```powershell [Windows PowerShell]
try {
  Invoke-RestMethod `
    -Method Get `
    -Uri http://localhost:8080/api/v1/system/health
} catch {
  $_.ErrorDetails.Message
}
```

```bash [macOS / Linux]
curl -i http://localhost:8080/api/v1/system/health
```

:::

应该看到 HTTP 状态码为 `401`，响应体类似：

```json
{
  "code": 40100,
  "message": "请先登录"
}
```

再登录获取 Token，并携带 Token 访问：

::: code-group

```powershell [Windows PowerShell]
$body = @{
  username = "admin"
  password = "YourPassword123"
} | ConvertTo-Json

$login = Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/auth/login `
  -ContentType "application/json" `
  -Body $body

$token = $login.data.access_token

Invoke-RestMethod `
  -Method Get `
  -Uri http://localhost:8080/api/v1/system/health `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123"}' | jq -r '.data.access_token')

curl -X GET http://localhost:8080/api/v1/system/health \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该看到统一成功响应。

::: details 怎么验证权限不足
本节默认管理员拥有 `super_admin` 角色，并且已经给这个角色写入健康检查权限，所以正常访问会成功。

如果要验证 `403`，可以临时把 `casbin_rule` 中这条策略删除或改成其他路径，然后重启服务，再携带同一个 Token 请求 `/api/v1/system/health`。验证后记得恢复策略。
:::

## 常见问题

::: details 启动时报 `relation "casbin_rule" does not exist`
说明 Casbin 策略表还没有创建。先执行 [`casbin_rule` 建表语句](../../reference/database-ddl#casbin-rule)，再重新启动服务。
:::

::: details 修改了 `casbin_rule`，权限没有立即变化
本节启动时会执行 `LoadPolicy`。服务运行期间直接改数据库，内存中的 Enforcer 不会自动刷新。

现在先重启服务让策略重新加载。后续做权限管理接口时，再补重新加载策略的代码路径。
:::

下一节会继续设计菜单和按钮权限：[菜单权限设计](./menu-permission)。
