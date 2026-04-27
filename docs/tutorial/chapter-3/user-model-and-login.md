---
title: 用户模型与登录接口
description: "设计用户模型，初始化默认管理员，并实现后台登录接口的基本验证链路。"
---

# 用户模型与登录接口

认证链路的第一步，是先让系统知道“用户是谁”。这一节先创建用户模型，启动时自动初始化一名默认管理员，并提供一个可以验证用户名和密码的登录接口。

::: tip 🎯 本节目标
完成后，数据库中会出现 `sys_user` 表；后端启动时会初始化默认管理员；访问 `/api/v1/auth/login` 可以完成用户名和密码校验。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ internal/
│  ├─ handler/
│  │  ├─ auth/
│  │  │  └─ login.go
│  │  └─ setup/
│  │     └─ setup.go
│  ├─ model/
│  │  └─ user.go
│  └─ router/
│     └─ router.go
├─ main.go
├─ go.mod
└─ go.sum
```

| 位置 | 用途 |
| --- | --- |
| `internal/model/user.go` | 定义用户表结构和用户状态 |
| `internal/handler/setup/setup.go` | 提供管理员初始化接口 |
| `internal/handler/auth/login.go` | 实现登录接口的请求解析、用户查询和密码校验 |
| `internal/router/router.go` | 注册 `/api/v1/auth/login` 和 `/api/v1/setup/init` |
| `main.go` | 服务启动时执行数据库迁移 |

## 🛠️ 安装密码哈希依赖

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

安装 `bcrypt` 所在的加密库：

```bash
# 安装密码哈希依赖
go get golang.org/x/crypto@latest
```

依赖资料入口：

| 依赖 | 用途 | 资料 |
| --- | --- | --- |
| `golang.org/x/crypto/bcrypt` | 对密码做不可逆哈希，并在登录时校验密码 | [Go 包文档](https://pkg.go.dev/golang.org/x/crypto/bcrypt) |

::: warning ⚠️ 密码不能明文保存
数据库里只能保存密码哈希，不能保存用户输入的原始密码。登录时也不是把明文密码查出来比较，而是用 `bcrypt.CompareHashAndPassword` 校验。
:::

## 先创建数据表

本节新增 `sys_user`，用于保存后台账号、密码哈希、用户状态和逻辑删除信息。

::: tip 建表 SQL
字段说明、索引设计、逻辑删除约定和 PostgreSQL / MySQL 建表语句统一放在参考手册：[数据库建表语句 - `sys_user`](../../reference/database-ddl#sys-user)。
:::

## 🛠️ 创建用户模型

创建 `server/internal/model/user.go`。这是新增文件，直接完整写入即可。

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// UserStatus 表示用户状态。
type UserStatus int

const (
	// UserStatusEnabled 表示用户可以正常登录。
	UserStatusEnabled UserStatus = 1
	// UserStatusDisabled 表示用户已被禁用。
	UserStatusDisabled UserStatus = 2
)

// User 是后台用户表模型。
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Nickname     string         `gorm:"size:64;not null;default:''" json:"nickname"`
	Status       UserStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定用户表名，避免后续调整命名策略时影响已有表。
func (User) TableName() string {
	return "sys_user"
}
```

## 🛠️ 创建管理员初始化接口

创建 `server/internal/handler/setup/setup.go`。这是新增文件，用于提供管理员初始化接口。

```go
package handler

import (
	"net/http"

	"ez-admin-gin/server/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SetupHandler 处理管理员一次性初始化。
type SetupHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewSetupHandler 创建 SetupHandler 实例。
func NewSetupHandler(db *gorm.DB, log *zap.Logger) *SetupHandler {
	return &SetupHandler{db: db, log: log}
}

// InitRequest 是管理员初始化接口的请求体。
type InitRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Nickname string `json:"nickname" binding:"required,min=1,max=64"`
}

// Init 创建第一个管理员账号并绑定到 super_admin 角色。
// POST /api/v1/setup/init
func (h *SetupHandler) Init(c *gin.Context) {
	// 检查是否已初始化（sys_user 是否有记录）
	var count int64
	if err := h.db.Model(&model.User{}).Count(&count).Error; err != nil {
		h.log.Error("check init status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查初始化状态失败"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "系统已初始化，不能重复执行"})
		return
	}

	var req InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// bcrypt 加密密码
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建管理员用户
	user := model.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Nickname:     req.Nickname,
		Status:       model.UserStatusEnabled,
	}
	if err := h.db.Create(&user).Error; err != nil {
		h.log.Error("create admin user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建管理员失败"})
		return
	}

	// 绑定到 super_admin 角色（ID=1）
	userRole := model.UserRole{
		UserID: user.ID,
		RoleID: 1,
	}
	if err := h.db.Create(&userRole).Error; err != nil {
		h.log.Error("bind admin role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "绑定管理员角色失败"})
		return
	}

	h.log.Info("admin user initialized", zap.String("username", req.Username))

	c.JSON(http.StatusOK, gin.H{
		"message":  "管理员账号创建成功",
		"user_id":  user.ID,
		"username": user.Username,
	})
}
```

::: tip 📌 管理员初始化接口
管理员账号通过 `/api/v1/setup/init` 接口创建，而不是在启动时自动生成。这样可以让用户设置自己的管理员账号和密码，更加安全。
:::

## 🛠️ 创建登录 Handler

创建 `server/internal/handler/auth/login.go`。这是新增文件，直接完整写入即可。

```go
package auth

import (
	"errors"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginHandler 负责登录相关接口。
type LoginHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewLoginHandler 创建登录 Handler。
func NewLoginHandler(db *gorm.DB, log *zap.Logger) *LoginHandler {
	return &LoginHandler{
		db:  db,
		log: log,
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// Login 校验用户名和密码。
func (h *LoginHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("用户名和密码不能为空"), h.log)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		response.Error(c, apperror.BadRequest("用户名和密码不能为空"), h.log)
		return
	}

	var user model.User
	// GORM 会自动过滤 deleted_at 不为空的记录。
	err := h.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, apperror.Unauthorized("用户名或密码错误"), h.log)
			return
		}

		h.log.Error("query login user failed", zap.Error(err))
		response.Error(c, apperror.Internal("登录失败", err), h.log)
		return
	}

	if user.Status != model.UserStatusEnabled {
		response.Error(c, apperror.Forbidden("用户已被禁用"), h.log)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, apperror.Unauthorized("用户名或密码错误"), h.log)
		return
	}

	response.Success(c, loginResponse{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	})
}
```

::: details 为什么用户名或密码错误返回同一句话
登录失败时不要区分“用户名不存在”和“密码错误”。如果提示得太细，攻击者可以借此枚举系统里有哪些账号。
:::

## 🛠️ 注册登录和初始化路由

修改 `server/internal/router/router.go`。这一处重点看三个变化：

- 新增登录 Handler 的 import。
- 新增 setup Handler 的 import。
- 注册 `/api/v1/auth/login` 和 `/api/v1/setup/init`。

先调整 import：

```go
import (
	"ez-admin-gin/server/internal/config"
	authHandler "ez-admin-gin/server/internal/handler/auth" // [!code ++]
	setupHandler "ez-admin-gin/server/internal/handler/setup" // [!code ++]
	systemHandler "ez-admin-gin/server/internal/handler/system"
	appLogger "ez-admin-gin/server/internal/logger"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)
```

在 `New` 中注册认证和初始化路由：

```go
func New(opts Options) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	registerSystemRoutes(r, opts)
	registerAuthRoutes(r, opts) // [!code ++]
	registerSetupRoutes(r, opts) // [!code ++]

	return r
}
```

继续在文件末尾新增：

```go
// registerAuthRoutes 注册认证相关路由。
func registerAuthRoutes(r *gin.Engine, opts Options) {
	login := authHandler.NewLoginHandler(opts.DB, opts.Log)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login", login.Login)
}

// registerSetupRoutes 注册初始化相关路由。
func registerSetupRoutes(r *gin.Engine, opts Options) {
	setup := setupHandler.NewSetupHandler(opts.DB, opts.Log)

	api := r.Group("/api/v1")
	setupGroup := api.Group("/setup")
	setupGroup.POST("/init", setup.Init)
}
```

::: details 为什么登录接口放在 `/api/v1/auth/login`
登录属于认证能力，不属于系统健康检查，所以单独放在 `auth` 分组下。后续刷新 Token、退出登录、获取当前用户信息，也可以继续放在这个分组里。
:::

## 🛠️ 在启动入口执行数据库迁移

修改 `server/main.go`。这一处重点看数据库迁移的执行：

在数据库连接成功后，执行数据库迁移：

```go
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

	// 执行数据库迁移，创建表结构和种子数据
	if err := database.Migrate(db, log); err != nil {
		log.Fatal("database migration", zap.Error(err))
	}
```

::: tip 📌 数据库迁移
服务启动时会自动执行数据库迁移，创建表结构和种子数据（包括超级管理员角色、菜单和权限）。管理员账号需要通过 `/api/v1/setup/init` 接口创建。
:::

## ✅ 整理依赖并启动

整理依赖：

```bash
# 整理新增依赖，更新 go.mod 和 go.sum
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

服务启动后，需要先通过初始化接口创建管理员账号：

::: code-group

```powershell [Windows PowerShell]
# 创建管理员账号
$body = @{
  username = "admin"
  password = "YourPassword123"
  nickname = "管理员"
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/setup/init `
  -ContentType "application/json" `
  -Body $body
```

```bash [macOS / Linux]
# 创建管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

:::

应该看到类似结果：

```json
{
  "message": "管理员账号创建成功",
  "user_id": 1,
  "username": "admin"
}
```

## ✅ 验证用户表和管理员账号

打开另一个终端，在项目根目录执行：

```bash
# 查看管理员账号是否已经写入数据库
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, username, nickname, status, deleted_at from sys_user;"
```

应该看到类似结果：

```text
 id | username | nickname | status | deleted_at
----+----------+----------+--------+------------
  1 | admin    | 管理员   |      1 |
```

::: details 如果提示 `relation "sys_user" does not exist`
说明用户表还没有创建。服务启动时会自动执行数据库迁移，创建表结构。如果迁移失败，查看服务启动日志获取详细信息。
:::

## ✅ 验证登录成功

保持后端服务运行，使用刚才创建的管理员账号登录：

::: code-group

```powershell [Windows PowerShell]
# 使用创建的管理员登录
$body = @{
  username = "admin"
  password = "YourPassword123"
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/auth/login `
  -ContentType "application/json" `
  -Body $body
```

```bash [macOS / Linux]
# 使用创建的管理员登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123"}'
```

:::

应该看到类似结果：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "user_id": 1,
    "username": "admin",
    "nickname": "管理员"
  }
}
```

::: info 为什么这一节还没有返回 Token
这一节先验证用户名和密码是否正确。下一节会在登录成功后签发 JWT，并让后续接口通过 Token 识别当前用户。
:::

## ✅ 验证登录失败

继续用错误密码请求登录：

::: code-group

```powershell [Windows PowerShell]
# 非 2xx 响应会进入 catch，这里直接打印响应体
$body = @{
  username = "admin"
  password = "wrong-password"
} | ConvertTo-Json

try {
  Invoke-RestMethod `
    -Method Post `
    -Uri http://localhost:8080/api/v1/auth/login `
    -ContentType "application/json" `
    -Body $body
} catch {
  $_.ErrorDetails.Message
}
```

```bash [macOS / Linux]
# -i 可以同时查看 HTTP 状态码和响应体
curl -i -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrong-password"}'
```

:::

应该能看到 HTTP 状态码为 `401`，响应体类似：

```json
{
  "code": 40100,
  "message": "用户名或密码错误"
}
```

## 常见问题

::: details 系统已初始化，不能重复执行
如果调用 `/api/v1/setup/init` 接口时收到这个错误，说明系统已经初始化过，管理员账号已经存在。

本地验证时可以物理删除管理员账号后重新创建：

```bash
# 只用于本地重置管理员账号
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "delete from sys_user where username = 'admin';"
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "delete from sys_user_role where user_id = 1;"
```
:::

::: details 提示 `用户名和密码不能为空`
检查请求体是否是 JSON，并确认请求头包含：

```http
Content-Type: application/json
```
:::

::: details 提示 `no required module provides package golang.org/x/crypto/bcrypt`
说明密码哈希依赖还没有加入当前 module。回到 `server/` 目录执行：

```bash
# 安装并整理依赖
go get golang.org/x/crypto@latest
go mod tidy
```
:::

下一节会在登录成功后签发 Token：[JWT 认证](./jwt-auth)。
