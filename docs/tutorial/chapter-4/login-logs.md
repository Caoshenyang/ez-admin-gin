---
title: 登录日志
description: "记录登录成功和失败行为，为安全审计提供基础信息。"
---

# 登录日志

操作日志记录的是已登录用户的后台写操作。登录日志补齐另一类安全审计：谁尝试登录、是否成功、来自哪个 IP、使用了什么客户端。

::: tip 🎯 本节目标
完成后，登录成功、密码错误、用户禁用等登录结果都会写入 `sys_login_log`；`super_admin` 可以查询登录日志列表。
:::

## 登录日志和操作日志的区别

| 日志 | 记录对象 | 触发位置 |
| --- | --- | --- |
| 操作日志 | 已登录用户的后台写操作 | 受保护路由中间件 |
| 登录日志 | 登录成功和失败行为 | 登录接口内部 |

登录接口本身还没有通过认证中间件，所以登录日志不能像操作日志那样依赖 `CurrentUserID`。它需要在登录 Handler 里主动记录。

::: warning ⚠️ 登录日志不要记录密码
登录失败时也只记录用户名、IP、User-Agent 和失败原因，不记录明文密码，也不记录请求体。
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
│  │  ├─ auth/
│  │  │  └─ login.go
│  │  └─ system/
│  │     └─ login_logs.go
│  ├─ model/
│  │  └─ login_log.go
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
| `docs/reference/database-ddl.md` | 补充 `sys_login_log` 建表语句 |
| `internal/model/login_log.go` | 定义登录日志模型 |
| `internal/handler/auth/login.go` | 登录成功和失败时写入日志 |
| `internal/handler/system/login_logs.go` | 提供登录日志查询接口 |
| `internal/router/router.go` | 注册登录日志查询路由 |
| `migrations/{postgres,mysql}/000002_seed_data.up.sql` | 初始化登录日志权限和菜单 |

## 先创建数据表

本节新增 `sys_login_log`，用于保存后台登录成功和失败记录。

`sys_login_log` 表保存后台登录成功和失败记录，不做逻辑删除。字段和索引详情见 [数据库建表语句 - `sys_login_log`](/reference/database-ddl#sys-login-log)。

## 🛠️ 创建登录日志模型

::: details `server/internal/model/login_log.go` — 登录日志模型

```go
package model

import "time"

// LoginLogStatus 表示登录结果。
type LoginLogStatus int

const (
	// LoginLogStatusSuccess 表示登录成功。
	LoginLogStatusSuccess LoginLogStatus = 1
	// LoginLogStatusFailed 表示登录失败。
	LoginLogStatusFailed LoginLogStatus = 2
)

// LoginLog 是后台登录日志模型。
type LoginLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;default:0;index" json:"user_id"`
	Username  string         `gorm:"size:64;not null;default:'';index" json:"username"`
	Status    LoginLogStatus `gorm:"type:smallint;not null;index" json:"status"`
	Message   string         `gorm:"size:255;not null;default:''" json:"message"`
	IP        string         `gorm:"column:ip;size:64;not null;default:'';index" json:"ip"`
	UserAgent string         `gorm:"size:500;not null;default:''" json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`
}

// TableName 固定登录日志表名。
func (LoginLog) TableName() string {
	return "sys_login_log"
}
```

:::

## 🛠️ 修改登录接口写入日志

::: details `server/internal/handler/auth/login.go` — 完整版（含日志记录）

```go
package auth

import (
	"errors"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"
	"ez-admin-gin/server/internal/token"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginHandler 负责登录相关接口。
type LoginHandler struct {
	db           *gorm.DB
	log          *zap.Logger
	tokenManager *token.Manager
}

// NewLoginHandler 创建登录 Handler。
func NewLoginHandler(db *gorm.DB, log *zap.Logger, tokenManager *token.Manager) *LoginHandler {
	return &LoginHandler{
		db:           db,
		log:          log,
		tokenManager: tokenManager,
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
}

// Login 校验用户名和密码。
func (h *LoginHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.createLoginLog(c, 0, "", model.LoginLogStatusFailed, "用户名和密码不能为空")
		response.Error(c, apperror.BadRequest("用户名和密码不能为空"), h.log)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		h.createLoginLog(c, 0, req.Username, model.LoginLogStatusFailed, "用户名和密码不能为空")
		response.Error(c, apperror.BadRequest("用户名和密码不能为空"), h.log)
		return
	}

	var user model.User
	// GORM 会自动过滤 deleted_at 不为空的记录。
	err := h.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.createLoginLog(c, 0, req.Username, model.LoginLogStatusFailed, "用户名或密码错误")
			response.Error(c, apperror.Unauthorized("用户名或密码错误"), h.log)
			return
		}

		h.createLoginLog(c, 0, req.Username, model.LoginLogStatusFailed, "登录失败")
		h.log.Error("query login user failed", zap.Error(err))
		response.Error(c, apperror.Internal("登录失败", err), h.log)
		return
	}

	if user.Status != model.UserStatusEnabled {
		h.createLoginLog(c, user.ID, user.Username, model.LoginLogStatusFailed, "用户已被禁用")
		response.Error(c, apperror.Forbidden("用户已被禁用"), h.log)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.createLoginLog(c, user.ID, user.Username, model.LoginLogStatusFailed, "用户名或密码错误")
		response.Error(c, apperror.Unauthorized("用户名或密码错误"), h.log)
		return
	}

	accessToken, expiresAt, err := h.tokenManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		h.createLoginLog(c, user.ID, user.Username, model.LoginLogStatusFailed, "登录失败")
		response.Error(c, apperror.Internal("登录失败", err), h.log)
		return
	}

	h.createLoginLog(c, user.ID, user.Username, model.LoginLogStatusSuccess, "登录成功")
	response.Success(c, loginResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.UTC().Format(time.RFC3339),
	})
}

func (h *LoginHandler) createLoginLog(c *gin.Context, userID uint, username string, status model.LoginLogStatus, message string) {
	record := model.LoginLog{
		UserID:    userID,
		Username:  strings.TrimSpace(username),
		Status:    status,
		Message:   message,
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	if err := h.db.Create(&record).Error; err != nil && h.log != nil {
		h.log.Warn("create login log failed", zap.Error(err))
	}
}
```

:::

::: warning ⚠️ 登录日志写入失败不应该阻断登录
登录日志是审计能力，不是登录成功的前置条件。如果日志写入失败，记录服务端日志即可，不要让用户因为日志表短暂异常而无法登录。
:::

## 🛠️ 创建登录日志查询接口

::: details `server/internal/handler/system/login_logs.go` — 登录日志查询接口

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

// LoginLogHandler 负责登录日志查询接口。
type LoginLogHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewLoginLogHandler 创建登录日志 Handler。
func NewLoginLogHandler(db *gorm.DB, log *zap.Logger) *LoginLogHandler {
	return &LoginLogHandler{
		db:  db,
		log: log,
	}
}

type loginLogListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `form:"username"`
	IP       string `form:"ip"`
	Status   int    `form:"status"`
}

type loginLogResponse struct {
	ID        uint                 `json:"id"`
	UserID    uint                 `json:"user_id"`
	Username  string               `json:"username"`
	Status    model.LoginLogStatus `json:"status"`
	Message   string               `json:"message"`
	IP        string               `json:"ip"`
	UserAgent string               `json:"user_agent"`
	CreatedAt time.Time            `json:"created_at"`
}

type loginLogListResponse struct {
	Items    []loginLogResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// List 返回登录日志分页列表。
func (h *LoginLogHandler) List(c *gin.Context) {
	var query loginLogListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeLoginLogPage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.LoginLog{})

	username := strings.TrimSpace(query.Username)
	if username != "" {
		queryDB = queryDB.Where("username = ?", username)
	}

	ip := strings.TrimSpace(query.IP)
	if ip != "" {
		queryDB = queryDB.Where("ip = ?", ip)
	}

	if query.Status != 0 {
		status := model.LoginLogStatus(query.Status)
		if !validLoginLogStatus(status) {
			response.Error(c, apperror.BadRequest("登录状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询登录日志总数失败", err), h.log)
		return
	}

	var logs []model.LoginLog
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		response.Error(c, apperror.Internal("查询登录日志列表失败", err), h.log)
		return
	}

	items := make([]loginLogResponse, 0, len(logs))
	for _, item := range logs {
		items = append(items, buildLoginLogResponse(item))
	}

	response.Success(c, loginLogListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func normalizeLoginLogPage(page int, pageSize int) (int, int) {
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

func validLoginLogStatus(status model.LoginLogStatus) bool {
	return status == model.LoginLogStatusSuccess || status == model.LoginLogStatusFailed
}

func buildLoginLogResponse(item model.LoginLog) loginLogResponse {
	return loginLogResponse{
		ID:        item.ID,
		UserID:    item.UserID,
		Username:  item.Username,
		Status:    item.Status,
		Message:   item.Message,
		IP:        item.IP,
		UserAgent: item.UserAgent,
		CreatedAt: item.CreatedAt,
	}
}
```

:::

## 🛠️ 注册登录日志查询路由

::: details `server/internal/router/router.go` — 注册登录日志路由

```go
// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
	users := systemHandler.NewUserHandler(opts.DB, opts.Log)
	roles := systemHandler.NewRoleHandler(opts.DB, opts.Log)
	menus := systemHandler.NewMenuAdminHandler(opts.DB, opts.Log)
	configs := systemHandler.NewSystemConfigHandler(opts.DB, opts.Redis, opts.Log)
	files := systemHandler.NewFileHandler(opts.DB, opts.Config.Upload, opts.Log)
	operationLogs := systemHandler.NewOperationLogHandler(opts.DB, opts.Log)
	loginLogs := systemHandler.NewLoginLogHandler(opts.DB, opts.Log) // [!code ++]

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.OperationLog(opts.DB, opts.Log))
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
	system.GET("/operation-logs", operationLogs.List)
system.GET("/login-logs", loginLogs.List) // [!code ++]
}
```

:::

## 🛠️ 初始化登录日志权限和菜单

登录日志的权限和菜单已经在数据库迁移文件中初始化。迁移文件会在服务启动时自动执行，创建登录日志相关的权限策略和菜单数据。

::: tip 💡 权限和菜单初始化
- 权限策略：在 `migrations/{postgres,mysql}/000002_seed_data.up.sql` 中插入登录日志接口的 Casbin 规则
- 菜单数据：在同一迁移文件中插入登录日志菜单和按钮
- 角色菜单绑定：在同一迁移文件中绑定 `super_admin` 角色到登录日志菜单
:::

## ✅ 启动并观察初始化日志

本节没有新增第三方依赖，可以直接启动：

```bash
# 在 server/ 目录启动服务
go run .
```

第一次启动后，控制台应该能看到类似日志：

```text
INFO	default permission created	{"role_code": "super_admin", "path": "/api/v1/system/login-logs", "method": "GET"}
INFO	default menu created	{"menu_code": "system:login-log"}
INFO	default role menu bound	{"role_id": 1, "menu_id": 25}
```

## ✅ 验证登录日志写入

先发起一次失败登录：

::: code-group

```powershell [Windows PowerShell]
$body = @{
  username = "admin"
  password = "wrong-password"
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/auth/login `
  -ContentType "application/json" `
  -Body $body
```

```bash [macOS / Linux]
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrong-password"}'
```

:::

这次请求应该返回“用户名或密码错误”。

再发起一次成功登录：

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

然后查询数据库：

```bash
# 查看最近的登录日志
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, user_id, username, status, message, ip, created_at from sys_login_log order by id desc limit 5;"
```

应该至少能看到两条记录：

- `status = 1`，`message = 登录成功`
- `status = 2`，`message = 用户名或密码错误`

::: warning ⚠️ 失败登录不会保存密码
验证数据库时，只应该看到用户名、登录结果、IP、User-Agent 等信息，不应该看到明文密码。
:::

## ✅ 验证登录日志查询接口

调用登录日志列表接口：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri "http://localhost:8080/api/v1/system/login-logs?page=1&page_size=10" `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl "http://localhost:8080/api/v1/system/login-logs?page=1&page_size=10" \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该能看到包含成功和失败登录记录的分页结果。

也可以按状态筛选：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri "http://localhost:8080/api/v1/system/login-logs?status=2&page=1&page_size=10" `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl "http://localhost:8080/api/v1/system/login-logs?status=2&page=1&page_size=10" \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

`status = 2` 表示登录失败。

## 常见问题

::: details 为什么登录日志不是中间件实现
登录接口发生在认证之前，此时还没有当前登录用户上下文。直接在登录 Handler 中记录成功和失败结果，语义更明确，也更容易拿到具体失败原因。
:::

::: details 用户名不存在时为什么 `user_id` 是 `0`
如果用户名不存在，就没有对应的用户 ID。此时只记录请求传入的用户名、失败结果和来源信息。
:::

::: details 登录日志会不会被操作日志重复记录
不会。操作日志中间件挂在 `/api/v1/system` 受保护分组上，登录接口是 `/api/v1/auth/login`，两者不在同一个路由分组。
:::

到这里，第四章的通用系统模块主线就补齐了。下一章开始进入前端管理台：[第 5 章：前端管理台](../chapter-5/)。
