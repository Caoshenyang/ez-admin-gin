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
│  ├─ bootstrap/
│  │  └─ bootstrap.go
│  ├─ handler/
│  │  └─ auth/
│  │     └─ login.go
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
| `internal/bootstrap/bootstrap.go` | 启动时创建默认管理员 |
| `internal/handler/auth/login.go` | 实现登录接口的请求解析、用户查询和密码校验 |
| `internal/router/router.go` | 注册 `/api/v1/auth/login` |
| `main.go` | 服务启动时执行初始化逻辑 |

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

字段说明：

| 字段 | 说明 |
| --- | --- |
| `ID` | 用户记录主键，由数据库自增生成 |
| `Username` | 登录用户名，设置唯一索引 |
| `PasswordHash` | 密码哈希，不通过 JSON 返回 |
| `Nickname` | 管理台展示名称 |
| `Status` | 用户状态，`1` 表示启用，`2` 表示禁用 |
| `CreatedAt` | 创建时间 |
| `UpdatedAt` | 更新时间 |
| `DeletedAt` | 逻辑删除时间，删除后默认不会被普通查询查出 |

完整建表语句可以看参考手册：[数据库建表语句：`sys_user`](../../reference/database-ddl#sys-user)。

::: details 主键是怎么生成的
本项目默认使用数据库自增 BIGINT 主键。创建用户时不需要在代码里给 `ID` 赋值，数据库会自动生成，GORM 会在 `Create` 成功后把生成的主键回填到结构体中。

`username` 这类业务标识单独做唯一字段，不和主键混用。
:::

::: details 为什么表名叫 `sys_user`
后台底座里通常会有用户、角色、菜单、日志等系统表。这里先用 `sys_` 前缀把系统表区分出来，后续业务表可以使用自己的模块前缀。

表名使用单数形式，是因为表对应的是一种实体模型，表里的多行数据才是集合。后续 `role`、`menu` 这类系统表也会沿用这个约定。
:::

::: details 为什么 `CreatedAt` 和 `UpdatedAt` 没有写 `gorm` 标签
这是 GORM 的内置约定。字段名叫 `CreatedAt`、`UpdatedAt`，类型是 `time.Time` 时，GORM 会自动映射成 `created_at`、`updated_at`，并在创建和更新数据时维护时间。

本项目约定时间字段由应用代码维护，不依赖数据库默认函数或触发器。如果后续直接写初始化 SQL，也要显式写入 `created_at` 和 `updated_at`。

只有需要改默认行为时，才需要额外写 `gorm` 标签。例如修改字段类型、改列名、设置索引、设置非空约束，或者像 `DeletedAt` 这样需要告诉 GORM “这是逻辑删除字段”。
:::

::: details 为什么 `username` 使用普通唯一索引
`username` 是账号身份标识，本节默认使用普通唯一索引。也就是说，即使账号被逻辑删除，相同用户名也不允许再次创建。

这样做更适合后台账号：历史记录、审计日志和操作者身份不会因为用户名复用而产生歧义。关于逻辑删除与唯一索引的更多背景，可以看：[逻辑删除与唯一索引冲突](../../reference/logical-delete-and-unique-index)。
:::

## 🛠️ 执行用户表建表 SQL

本节的表结构通过 SQL 建表脚本准备。先打开参考手册中的 [`sys_user` 建表语句](../../reference/database-ddl#sys-user)，按当前数据库类型选择对应版本执行。

当前本地环境使用 PostgreSQL，可以进入数据库客户端后粘贴 PostgreSQL 标签页中的 SQL：

```bash
# 在项目根目录执行，进入本地 PostgreSQL
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin
```

执行完成后，可以在 `psql` 中确认表已经创建：

```sql
-- 查看 sys_user 表结构
\d+ sys_user
```

::: warning ⚠️ 先建表，再启动后端
后面的启动初始化只负责创建默认管理员，不负责创建 `sys_user` 表。如果跳过建表 SQL，服务启动时查询 `sys_user` 会失败。
:::

## 🛠️ 创建启动初始化

创建 `server/internal/bootstrap/bootstrap.go`。这是新增文件，直接完整写入即可。

```go
package bootstrap

import (
	"errors"
	"fmt"

	"ez-admin-gin/server/internal/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "EzAdmin@123456"
)

// Run 执行服务启动时必须完成的初始化动作。
func Run(db *gorm.DB, log *zap.Logger) error {
	if err := seedDefaultAdmin(db, log); err != nil {
		return fmt.Errorf("seed default admin: %w", err)
	}

	return nil
}

// seedDefaultAdmin 创建本地起步用的默认管理员。
func seedDefaultAdmin(db *gorm.DB, log *zap.Logger) error {
	var user model.User
	// Unscoped 会把已逻辑删除记录也查出来，避免重复创建同名默认账号。
	err := db.Unscoped().Where("username = ?", defaultAdminUsername).First(&user).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash default admin password: %w", err)
	}

	user = model.User{
		Username:     defaultAdminUsername,
		PasswordHash: string(passwordHash),
		Nickname:     "系统管理员",
		Status:       model.UserStatusEnabled,
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	log.Info("default admin user created", zap.String("username", defaultAdminUsername))

	return nil
}
```

::: warning ⚠️ 默认管理员只用于本地起步
默认账号是 `admin`，默认密码是 `EzAdmin@123456`。后续接入系统配置和初始化脚本时，应把默认密码改成可配置或首次启动后强制修改。
:::

::: details 为什么查询默认管理员时用了 `Unscoped`
本节的用户名使用普通唯一索引，默认不允许逻辑删除后复用。如果 `admin` 曾经被逻辑删除，普通查询会查不到它，但数据库中的唯一索引仍然会拦截同名账号。

所以初始化默认管理员时使用 `Unscoped` 把历史记录也查出来，避免启动时重复创建同名账号。
:::

::: details 为什么建表 SQL 放在参考手册
表结构属于数据库迁移内容，放在 SQL 脚本里更清晰：字段注释、表注释、索引和跨数据库差异都能明确看到。

启动初始化只处理“运行时需要准备的数据”，比如默认管理员。这样表结构变更和业务初始化不会混在一起。
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

## 🛠️ 注册登录路由

修改 `server/internal/router/router.go`。这一处重点看两个变化：

- 新增登录 Handler 的 import。
- 注册 `/api/v1/auth/login`。

先调整 import：

```go
import (
	"ez-admin-gin/server/internal/config"
	authHandler "ez-admin-gin/server/internal/handler/auth" // [!code ++]
	systemHandler "ez-admin-gin/server/internal/handler/system"
	appLogger "ez-admin-gin/server/internal/logger"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)
```

在 `New` 中注册认证路由：

```go
func New(opts Options) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	registerSystemRoutes(r, opts)
	registerAuthRoutes(r, opts) // [!code ++]

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
```

::: details 为什么登录接口放在 `/api/v1/auth/login`
登录属于认证能力，不属于系统健康检查，所以单独放在 `auth` 分组下。后续刷新 Token、退出登录、获取当前用户信息，也可以继续放在这个分组里。
:::

## 🛠️ 在启动入口执行初始化

修改 `server/main.go`。这一处重点看两个变化：

- 引入 `internal/bootstrap`。
- 数据库连接成功后执行 `bootstrap.Run`。

先调整 import：

```go
import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/bootstrap" // [!code ++]
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/router"

	"go.uber.org/zap"
)
```

在数据库连接成功后，增加初始化调用：

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

	// 数据库表已通过 SQL 准备好，这里只创建默认管理员。
	if err := bootstrap.Run(db, log); err != nil { // [!code ++]
		log.Fatal("bootstrap application", zap.Error(err)) // [!code ++]
	} // [!code ++]
```

::: warning ⚠️ 初始化要放在数据库连接之后
`bootstrap.Run` 会查询并写入默认管理员，必须在数据库连接成功后执行。如果放在数据库连接前，会没有可用的 `db` 对象。
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
INFO	default admin user created	{"username": "admin"}
INFO	server started	{"addr": ":8080", "env": "dev"}
```

如果默认管理员已经存在，后续启动不会重复创建。

## ✅ 验证用户表和默认管理员

打开另一个终端，在项目根目录执行：

```bash
# 查看默认管理员是否已经写入数据库
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, username, nickname, status, deleted_at from sys_user;"
```

应该看到类似结果：

```text
 id | username |  nickname  | status | deleted_at
----+----------+------------+--------+------------
  1 | admin    | 系统管理员 |      1 |
```

::: details 如果提示 `relation "sys_user" does not exist`
说明用户表还没有创建。先回到 [`sys_user` 建表语句](../../reference/database-ddl#sys-user)，执行对应数据库版本的 SQL，然后重新启动服务。
:::

## ✅ 验证登录成功

保持后端服务运行，调用登录接口：

::: code-group

```powershell [Windows PowerShell]
# 使用默认管理员登录
$body = @{
  username = "admin"
  password = "EzAdmin@123456"
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/auth/login `
  -ContentType "application/json" `
  -Body $body
```

```bash [macOS / Linux]
# 使用默认管理员登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"EzAdmin@123456"}'
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
    "nickname": "系统管理员"
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

::: details 修改默认密码后没有生效
默认管理员只在账号不存在时创建。如果 `admin` 已经存在，修改 `defaultAdminPassword` 不会自动覆盖已有密码。

本地验证时可以物理删除默认管理员后重新启动服务：

```bash
# 只用于本地重置默认管理员，重启服务后会重新创建
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "delete from sys_user where username = 'admin';"
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

## ✅ 确认 Git 状态

回到项目根目录：

::: code-group

```powershell [Windows PowerShell]
# 回到项目根目录后查看本节改动
Set-Location ..
git status
```

```bash [macOS / Linux]
# 回到项目根目录后查看本节改动
cd ..
git status
```

:::

应该能看到本节新增或修改的文件：

```text
server/internal/bootstrap/bootstrap.go
server/internal/handler/auth/login.go
server/internal/model/user.go
server/internal/router/router.go
server/main.go
server/go.mod
server/go.sum
```

下一节会在登录成功后签发 Token：[JWT 认证](./jwt-auth)。
