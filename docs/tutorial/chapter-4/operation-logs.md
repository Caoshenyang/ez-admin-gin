---
title: 操作日志
description: "记录后台用户的关键写操作，为审计追踪和问题排查提供基础。"
---

# 操作日志

前面已经有用户、角色、菜单、配置和文件上传能力。现在补齐操作日志：记录后台用户什么时候调用了哪个管理接口、请求是否成功、耗时多久。

::: tip 🎯 本节目标
完成后，后台的 `POST` 写操作会自动写入操作日志；`super_admin` 可以查询操作日志列表。
:::

## 先确定记录范围

本节默认记录后台管理接口中的写操作：

| 方法 | 是否记录 | 说明 |
| --- | --- | --- |
| `GET` | 否 | 查询请求太频繁，容易把日志刷满 |
| `POST` | 是 | 创建、编辑、状态变更、删除、上传文件等写操作 |

::: warning ⚠️ 操作日志不要全量保存请求体
请求体里可能包含密码、Token、文件内容或其他敏感信息。本节只记录方法、路径、查询参数、用户、IP、状态码和耗时，不把请求体原样写入数据库。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
docs/
└─ reference/
   └─ database-ddl.md

server/
├─ internal/
│  ├─ handler/
│  │  └─ system/
│  │     └─ operation_logs.go
│  ├─ middleware/
│  │  └─ operation_log.go
│  ├─ model/
│  │  └─ operation_log.go
│  └─ router/
│     └─ router.go
└─ migrations/
   ├─ postgres/
   │  └─ 000002_seed_data.up.sql
   └─ mysql/
      └─ 000002_seed_data.up.sql
```

| 位置 | 用途 |
| --- | --- |
| `docs/reference/database-ddl.md` | 补充 `sys_operation_log` 建表语句 |
| `internal/model/operation_log.go` | 定义操作日志模型 |
| `internal/middleware/operation_log.go` | 请求结束后自动写日志 |
| `internal/handler/system/operation_logs.go` | 提供操作日志查询接口 |
| `internal/router/router.go` | 注册日志中间件和查询路由 |
| `migrations/{postgres,mysql}/000002_seed_data.up.sql` | 初始化操作日志权限和菜单 |

## 先创建数据表

本节新增 `sys_operation_log`，用于保存后台用户的关键写操作审计记录。

`sys_operation_log` 表保存后台用户的写操作审计记录，不做逻辑删除。字段和索引详情见 [数据库建表语句 - `sys_operation_log`](/reference/database-ddl#sys-operation-log)。

## 🛠️ 创建操作日志模型

创建 `server/internal/model/operation_log.go`。这是新增文件，直接完整写入即可。

```go
package model

import "time"

// OperationLog 是后台操作日志模型。
type OperationLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;default:0;index" json:"user_id"`
	Username     string    `gorm:"size:64;not null;default:'';index" json:"username"`
	Method       string    `gorm:"size:10;not null;index" json:"method"`
	Path         string    `gorm:"size:255;not null;index" json:"path"`
	RoutePath    string    `gorm:"size:255;not null;default:'';index" json:"route_path"`
	Query        string    `gorm:"size:1000;not null;default:''" json:"query"`
	IP           string    `gorm:"column:ip;size:64;not null;default:''" json:"ip"`
	UserAgent    string    `gorm:"size:500;not null;default:''" json:"user_agent"`
	StatusCode   int       `gorm:"not null;default:0;index" json:"status_code"`
	LatencyMs    int64     `gorm:"not null;default:0" json:"latency_ms"`
	Success      bool      `gorm:"not null;default:true;index" json:"success"`
	ErrorMessage string    `gorm:"size:500;not null;default:''" json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 固定操作日志表名。
func (OperationLog) TableName() string {
	return "sys_operation_log"
}
```

## 🛠️ 创建操作日志中间件

创建 `server/internal/middleware/operation_log.go`。这是新增文件，直接完整写入即可。

::: details `server/internal/middleware/operation_log.go` — 操作日志中间件

```go
package middleware

import (
	"net/http"
	"strings"
	"time"

	"ez-admin-gin/server/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const maxOperationLogTextLength = 500

// OperationLog 在请求结束后记录后台写操作。
func OperationLog(db *gorm.DB, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		if shouldSkipOperationLog(c) {
			return
		}

		userID, _ := CurrentUserID(c)
		username, _ := CurrentUsername(c)
		statusCode := c.Writer.Status()

		record := model.OperationLog{
			UserID:       userID,
			Username:     username,
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			RoutePath:    c.FullPath(),
			Query:        truncateOperationLogText(c.Request.URL.RawQuery, 1000),
			IP:           c.ClientIP(),
			UserAgent:    truncateOperationLogText(c.Request.UserAgent(), maxOperationLogTextLength),
			StatusCode:   statusCode,
			LatencyMs:    time.Since(start).Milliseconds(),
			Success:      statusCode < http.StatusBadRequest,
			ErrorMessage: operationErrorMessage(c, statusCode),
		}

		if err := db.Create(&record).Error; err != nil && log != nil {
			log.Warn("create operation log failed", zap.Error(err))
		}
	}
}

func shouldSkipOperationLog(c *gin.Context) bool {
	method := c.Request.Method
	if method != http.MethodPost {
		return true
	}

	// 静态资源和未匹配到路由的请求不作为后台操作记录。
	if c.FullPath() == "" {
		return true
	}

	return false
}

func operationErrorMessage(c *gin.Context, statusCode int) string {
	if len(c.Errors) > 0 {
		return truncateOperationLogText(c.Errors.Last().Error(), maxOperationLogTextLength)
	}

	if statusCode >= http.StatusBadRequest {
		return http.StatusText(statusCode)
	}

	return ""
}

func truncateOperationLogText(value string, maxLength int) string {
	value = strings.TrimSpace(value)
	if len(value) <= maxLength {
		return value
	}

	return value[:maxLength]
}
```

:::

::: details 为什么中间件只记录请求结束后的结果
操作日志需要知道接口最终是成功还是失败、状态码是多少、耗时多久。这些信息只有在 `c.Next()` 执行完后才能拿到。
:::

## 🛠️ 创建操作日志查询接口

创建 `server/internal/handler/system/operation_logs.go`。这是新增文件，直接完整写入即可。

::: details `server/internal/handler/system/operation_logs.go` — 操作日志查询接口

```go
package system

import (
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OperationLogHandler 负责操作日志查询接口。
type OperationLogHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewOperationLogHandler 创建操作日志 Handler。
func NewOperationLogHandler(db *gorm.DB, log *zap.Logger) *OperationLogHandler {
	return &OperationLogHandler{
		db:  db,
		log: log,
	}
}

type operationLogListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `form:"username"`
	Method   string `form:"method"`
	Path     string `form:"path"`
	Success  string `form:"success"`
}

type operationLogResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	RoutePath    string    `json:"route_path"`
	Query        string    `json:"query"`
	IP           string    `json:"ip"`
	UserAgent    string    `json:"user_agent"`
	StatusCode   int       `json:"status_code"`
	LatencyMs    int64     `json:"latency_ms"`
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

type operationLogListResponse struct {
	Items    []operationLogResponse `json:"items"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

// List 返回操作日志分页列表。
func (h *OperationLogHandler) List(c *gin.Context) {
	var query operationLogListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeOperationLogPage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.OperationLog{})

	username := strings.TrimSpace(query.Username)
	if username != "" {
		queryDB = queryDB.Where("username = ?", username)
	}

	method := strings.ToUpper(strings.TrimSpace(query.Method))
	if method != "" {
		queryDB = queryDB.Where("method = ?", method)
	}

	path := strings.TrimSpace(query.Path)
	if path != "" {
		queryDB = queryDB.Where("path LIKE ?", "%"+path+"%")
	}

	if query.Success != "" {
		success, ok := parseOperationLogSuccess(query.Success)
		if !ok {
			response.Error(c, apperror.BadRequest("成功状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("success = ?", success)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询操作日志总数失败", err), h.log)
		return
	}

	var logs []model.OperationLog
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		response.Error(c, apperror.Internal("查询操作日志列表失败", err), h.log)
		return
	}

	items := make([]operationLogResponse, 0, len(logs))
	for _, item := range logs {
		items = append(items, buildOperationLogResponse(item))
	}

	response.Success(c, operationLogListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
```

:::

继续在同一个 `operation_logs.go` 中追加下面的辅助函数：

::: details `server/internal/handler/system/operation_logs.go` — 辅助函数

```go
func normalizeOperationLogPage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

func parseOperationLogSuccess(value string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1":
		return true, true
	case "false", "0":
		return false, true
	default:
		return false, false
	}
}

func buildOperationLogResponse(item model.OperationLog) operationLogResponse {
	return operationLogResponse{
		ID:           item.ID,
		UserID:       item.UserID,
		Username:     item.Username,
		Method:       item.Method,
		Path:         item.Path,
		RoutePath:    item.RoutePath,
		Query:        item.Query,
		IP:           item.IP,
		UserAgent:    item.UserAgent,
		StatusCode:   item.StatusCode,
		LatencyMs:    item.LatencyMs,
		Success:      item.Success,
		ErrorMessage: item.ErrorMessage,
		CreatedAt:    item.CreatedAt,
	}
}
```

:::

## 🛠️ 注册中间件和路由

修改 `server/internal/router/router.go`。本次要改两处：

- 新增操作日志 Handler
- 在系统路由分组中挂载操作日志中间件和查询接口

::: details `server/internal/router/router.go` — 挂载操作日志中间件与路由

```go
// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
	users := systemHandler.NewUserHandler(opts.DB, opts.Log)
	roles := systemHandler.NewRoleHandler(opts.DB, opts.Log)
	menus := systemHandler.NewMenuAdminHandler(opts.DB, opts.Log)
	configs := systemHandler.NewSystemConfigHandler(opts.DB, opts.Redis, opts.Log)
	files := systemHandler.NewFileHandler(opts.DB, opts.Config.Upload, opts.Log)
	operationLogs := systemHandler.NewOperationLogHandler(opts.DB, opts.Log) // [!code ++]

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.OperationLog(opts.DB, opts.Log)) // [!code ++]
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))
	system.GET("/health", health.Check)
	system.GET("/users", users.List)
	system.POST("/users", users.Create)
	system.POST("/users/:id/update", users.Update)
	system.POST("/users/:id/status", users.UpdateStatus)
	system.POST("/users/:id/roles", users.UpdateRoles)
	system.GET("/roles", roles.List)
	system.POST("/roles", roles.Create)
	system.POST("/roles/:id/update", roles.Update)
	system.POST("/roles/:id/status", roles.UpdateStatus)
	system.POST("/roles/:id/permissions", roles.UpdatePermissions)
	system.POST("/roles/:id/menus", roles.UpdateMenus)
	system.GET("/menus", menus.Tree)
	system.POST("/menus", menus.Create)
	system.POST("/menus/:id/update", menus.Update)
	system.POST("/menus/:id/status", menus.UpdateStatus)
	system.POST("/menus/:id/delete", menus.Delete)
	system.GET("/configs", configs.List)
	system.POST("/configs", configs.Create)
	system.POST("/configs/:id/update", configs.Update)
	system.POST("/configs/:id/status", configs.UpdateStatus)
	system.GET("/configs/value/:key", configs.Value)
	system.GET("/files", files.List)
	system.POST("/files", files.Upload)
	system.GET("/operation-logs", operationLogs.List) // [!code ++]
}
```

:::

::: details 为什么操作日志中间件放在权限中间件前面
顺序是：先认证，再进入操作日志中间件，再执行权限校验。这样即使某个已登录用户请求被权限拦截，写操作也能留下失败记录。
:::

## 🛠️ 初始化操作日志权限和菜单

操作日志的权限和菜单已经在数据库迁移文件中初始化。迁移文件会在服务启动时自动执行，创建操作日志相关的权限策略和菜单数据。

::: tip 💡 权限和菜单初始化
- 权限策略：在 `migrations/{postgres,mysql}/000002_seed_data.up.sql` 中插入操作日志接口的 Casbin 规则
- 菜单数据：在同一迁移文件中插入操作日志菜单和按钮
- 角色菜单绑定：在同一迁移文件中绑定 `super_admin` 角色到操作日志菜单
:::

::: warning ⚠️ 继续确认只有一处 `return menus, nil`
操作日志菜单要接在文件管理菜单后面。不要在中间提前 `return`，否则后面的菜单不会初始化，也不会授权给 `super_admin`。
:::

## ✅ 启动并观察初始化日志

本节没有新增第三方依赖，可以直接启动：

```bash
# 在 server/ 目录启动服务
go run .
```

第一次启动后，控制台应该能看到类似日志：

```text
INFO	default permission created	{"role_code": "super_admin", "path": "/api/v1/system/operation-logs", "method": "GET"}
INFO	default menu created	{"menu_code": "system:operation-log"}
INFO	default role menu bound	{"role_id": 1, "menu_id": 23}
```

## ✅ 验证权限和菜单数据

先确认操作日志查询权限已经写入：

```bash
# 查看操作日志相关接口权限
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ptype, v0, v1, v2 from casbin_rule where v1 like '/api/v1/system/operation-logs%' order by v1, v2;"
```

应该能看到 `GET /api/v1/system/operation-logs`。

再确认操作日志菜单已经写入：

```bash
# 查看操作日志菜单和按钮
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, parent_id, type, code, title from sys_menu where code like 'system:operation-log%' order by sort, id;"
```

应该能看到 `system:operation-log` 和 `system:operation-log:list`。

## ✅ 验证操作日志写入

先登录拿到 Token：

::: code-group

```powershell [Windows PowerShell]
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
```

```bash [macOS / Linux]
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"EzAdmin@123456"}' | jq -r '.data.access_token')
```

:::

调用一个写操作。这里继续使用文件上传接口，方便稳定触发一条 `POST` 操作日志：

::: code-group

```powershell [Windows PowerShell]
Set-Content -Path .\operation-log-test.txt -Value "operation log test" -Encoding UTF8

curl.exe -X POST "http://localhost:8080/api/v1/system/files?source=operation-log-test" `
  -H "Authorization: Bearer $token" `
  -F "file=@.\operation-log-test.txt"
```

```bash [macOS / Linux]
echo "operation log test" > operation-log-test.txt

curl -X POST "http://localhost:8080/api/v1/system/files?source=operation-log-test" \
  -H "Authorization: Bearer ${TOKEN}" \
  -F "file=@./operation-log-test.txt"
```

:::

上传成功后，查询数据库中的操作日志：

```bash
# 查看最近的操作日志
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, user_id, username, method, path, query, status_code, success, latency_ms, created_at from sys_operation_log order by id desc limit 5;"
```

应该能看到一条 `POST /api/v1/system/files` 记录，并且：

- `username` 是 `admin`
- `query` 是 `source=operation-log-test`
- `status_code` 是 `200`
- `success` 是 `true`
- `latency_ms` 大于或等于 `0`

## ✅ 验证查询接口

调用操作日志列表接口：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri "http://localhost:8080/api/v1/system/operation-logs?page=1&page_size=10" `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl "http://localhost:8080/api/v1/system/operation-logs?page=1&page_size=10" \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该能看到包含刚才上传文件操作的分页结果。

也可以按方法筛选：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri "http://localhost:8080/api/v1/system/operation-logs?method=POST&page=1&page_size=10" `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl "http://localhost:8080/api/v1/system/operation-logs?method=POST&page=1&page_size=10" \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

## 常见问题

::: details 为什么查询用户列表不会产生操作日志
本节默认不记录 `GET` 请求。查询类接口调用频率很高，如果全部记录，日志表会很快膨胀，也会掩盖真正重要的写操作。
:::

::: details 为什么操作日志里没有请求体
请求体可能包含密码、Token、文件内容或业务敏感字段。操作日志先记录“谁、什么时间、请求了哪个接口、结果如何”。如果某个模块确实需要更细的审计字段，建议在业务层单独记录脱敏后的关键字段。
:::

::: details 操作失败时也会记录吗
会。操作日志中间件放在权限校验之前、认证之后。只要用户已经登录，并且请求方法是 `POST`，即使后续权限失败或接口返回错误，也会记录一条失败日志。
:::

下一节继续补齐登录行为审计：[登录日志](./login-logs)。
