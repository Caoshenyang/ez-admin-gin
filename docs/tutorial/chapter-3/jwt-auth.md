---
title: JWT 认证
description: "生成和解析 JWT，让登录接口返回可以表达用户身份和过期时间的访问令牌。"
---

# JWT 认证

上一节已经能校验用户名和密码。这一节在登录成功后签发 `access_token`，让前端后续可以携带这个令牌访问需要登录的接口。

::: tip 🎯 本节目标
完成后，登录接口会返回 `access_token`、`token_type` 和 `expires_at`。Token 中会包含用户 ID、用户名、签发方和过期时间。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ configs/
│  └─ config.yaml
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  ├─ handler/
│  │  └─ auth/
│  │     └─ login.go
│  ├─ router/
│  │  └─ router.go
│  └─ token/
│     └─ jwt.go
├─ main.go
├─ go.mod
└─ go.sum
```

| 位置 | 用途 |
| --- | --- |
| `configs/config.yaml` | 增加 JWT 密钥、签发方和访问令牌有效期 |
| `internal/config/config.go` | 读取 `auth` 配置段 |
| `internal/token/jwt.go` | 封装 Token 生成和解析 |
| `internal/handler/auth/login.go` | 登录成功后签发 Token |
| `internal/router/router.go` | 把 Token 管理器传给登录 Handler |
| `main.go` | 创建 Token 管理器 |

## 🛠️ 安装 JWT 依赖

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

安装 JWT 依赖：

```bash
# 安装 JWT v5 版本
go get github.com/golang-jwt/jwt/v5@latest
```

依赖资料入口：

| 依赖 | 用途 | 资料 |
| --- | --- | --- |
| `github.com/golang-jwt/jwt/v5` | 生成和解析 JWT | [Go 包文档](https://pkg.go.dev/github.com/golang-jwt/jwt/v5) |

::: warning ⚠️ 使用 v5 导入路径
这里使用的是 `github.com/golang-jwt/jwt/v5`。不要写成旧的 `github.com/dgrijalva/jwt-go`，也不要漏掉最后的 `/v5`。
:::

## 🛠️ 增加认证配置

先修改 `server/configs/config.yaml`，在 `redis` 和 `log` 之间新增 `auth` 配置段：

```yaml
auth: # [!code ++]
  # JWT 签名密钥；本地开发可以先写在配置文件中，生产环境应改用环境变量覆盖。 # [!code ++]
  jwt_secret: "ez-admin-dev-secret-change-me-please-32" # [!code ++]
  # access_token 有效期，单位秒。这里先设置为 2 小时。 # [!code ++]
  access_token_ttl: 7200 # [!code ++]
  # Token 签发方，后续解析时会校验。 # [!code ++]
  issuer: "ez-admin" # [!code ++]
```

::: warning ⚠️ 生产环境不要使用示例密钥
`jwt_secret` 是签名密钥，泄露后别人就可能伪造 Token。本地开发可以先使用示例值，正式部署时要通过环境变量覆盖成更长、更随机的密钥。
:::

## 🛠️ 读取认证配置

修改 `server/internal/config/config.go`。这一处重点看三个变化：

- `Config` 增加 `Auth` 配置段。
- 新增 `AuthConfig` 结构体。
- 默认值和环境变量列表增加 `auth.*`。

先给 `Config` 增加字段：

```go
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Auth     AuthConfig     `mapstructure:"auth"` // [!code ++]
	Log      LogConfig      `mapstructure:"log"`
}
```

在 `RedisConfig` 后面新增：

```go
// AuthConfig 保存认证相关配置。
type AuthConfig struct {
	// JWTSecret 是 access_token 的签名密钥。
	JWTSecret string `mapstructure:"jwt_secret"`
	// AccessTokenTTL 是 access_token 有效期，单位秒。
	AccessTokenTTL int `mapstructure:"access_token_ttl"`
	// Issuer 是 Token 签发方。
	Issuer string `mapstructure:"issuer"`
}
```

在 `setDefaults` 中新增默认值：

```go
	v.SetDefault("auth.jwt_secret", "ez-admin-dev-secret-change-me-please-32") // [!code ++]
	v.SetDefault("auth.access_token_ttl", 7200) // [!code ++]
	v.SetDefault("auth.issuer", "ez-admin") // [!code ++]
```

在 `bindEnvs` 的 `keys` 中新增：

```go
		// 允许用 EZ_AUTH_JWT_SECRET 覆盖本地开发密钥。 // [!code ++]
		"auth.jwt_secret", // [!code ++]
		"auth.access_token_ttl", // [!code ++]
		"auth.issuer", // [!code ++]
```

::: details 为什么配置里用秒
配置文件里用整数秒更直观，也方便环境变量覆盖。代码里真正使用时，再转换成 `time.Duration`。
:::

## 🛠️ 创建 Token 管理器

创建 `server/internal/token/jwt.go`。这是新增文件，直接完整写入即可。

```go
package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"ez-admin-gin/server/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken 表示 Token 无效、过期或签名不正确。
	ErrInvalidToken = errors.New("invalid token")
)

// Claims 是写入 access_token 的业务载荷。
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Manager 负责生成和解析 access_token。
type Manager struct {
	secret         []byte
	issuer         string
	accessTokenTTL time.Duration
	now            func() time.Time
}

// NewManager 根据配置创建 Token 管理器。
func NewManager(cfg config.AuthConfig) (*Manager, error) {
	secret := strings.TrimSpace(cfg.JWTSecret)
	if len(secret) < 32 {
		return nil, fmt.Errorf("jwt secret must be at least 32 characters")
	}

	if cfg.AccessTokenTTL <= 0 {
		return nil, fmt.Errorf("access token ttl must be greater than 0")
	}

	issuer := strings.TrimSpace(cfg.Issuer)
	if issuer == "" {
		return nil, fmt.Errorf("jwt issuer cannot be empty")
	}

	return &Manager{
		secret:         []byte(secret),
		issuer:         issuer,
		accessTokenTTL: time.Duration(cfg.AccessTokenTTL) * time.Second,
		now:            time.Now,
	}, nil
}

// GenerateAccessToken 生成访问令牌，并返回令牌过期时间。
func (m *Manager) GenerateAccessToken(userID uint, username string) (string, time.Time, error) {
	now := m.now()
	expiresAt := now.Add(m.accessTokenTTL)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ParseAccessToken 解析并校验访问令牌。
func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}

			return m.secret, nil
		},
		jwt.WithIssuer(m.issuer),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
```

::: details 为什么要校验签名算法
解析 Token 时不能只看签名是否通过，还要确认算法就是系统预期的 `HS256`。否则错误配置或不安全的解析方式可能带来额外风险。
:::

::: details 为什么用了 `RegisteredClaims`
JWT 有一些标准字段，比如 `iss`、`sub`、`iat`、`exp`。`RegisteredClaims` 是 `jwt/v5` 推荐用来表达这些标准字段的结构。

本项目额外放了 `user_id` 和 `username`，方便后续中间件从 Token 中识别当前用户。
:::

## 🛠️ 在启动入口创建 Token 管理器

修改 `server/main.go`。这一处重点看两个变化：

- 引入 `internal/token`。
- 根据 `cfg.Auth` 创建 Token 管理器，并传给路由。

先调整 import：

```go
import (
	"ez-admin-gin/server/internal/bootstrap"
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/router"
	"ez-admin-gin/server/internal/token" // [!code ++]

	"go.uber.org/zap"
)
```

在创建路由前，增加 Token 管理器：

```go
	// Token 管理器负责签发和解析登录令牌。
	tokenManager, err := token.NewManager(cfg.Auth) // [!code ++]
	if err != nil { // [!code ++]
		log.Fatal("create token manager", zap.Error(err)) // [!code ++]
	} // [!code ++]

	// 路由注册交给 internal/router，main.go 只保留启动流程。
	r := router.New(router.Options{
		Config: cfg,
		Log:    log,
		DB:     db,
		Redis:  redisClient,
		Token:  tokenManager, // [!code ++]
	})
```

## 🛠️ 把 Token 管理器传给登录 Handler

修改 `server/internal/router/router.go`。这一处重点看两个变化：

- `Options` 增加 `Token` 字段。
- 创建登录 Handler 时传入 `opts.Token`。

先调整 import：

```go
import (
	"ez-admin-gin/server/internal/config"
	authHandler "ez-admin-gin/server/internal/handler/auth"
	systemHandler "ez-admin-gin/server/internal/handler/system"
	appLogger "ez-admin-gin/server/internal/logger"
	"ez-admin-gin/server/internal/token" // [!code ++]

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)
```

更新 `Options`：

```go
type Options struct {
	Config *config.Config
	Log    *zap.Logger
	DB     *gorm.DB
	Redis  *goredis.Client
	Token  *token.Manager // [!code ++]
}
```

更新 `registerAuthRoutes`：

```go
// registerAuthRoutes 注册认证相关路由。
func registerAuthRoutes(r *gin.Engine, opts Options) {
	login := authHandler.NewLoginHandler(opts.DB, opts.Log, opts.Token) // [!code focus]

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login", login.Login)
}
```

## 🛠️ 登录成功后返回 Token

修改 `server/internal/handler/auth/login.go`。这一处重点看三个变化：

- 引入 `internal/token`。
- `LoginHandler` 保存 `tokenManager`。
- 登录成功后生成 `access_token`。

先调整 import：

```go
import (
	"errors"
	"strings"
	"time" // [!code ++]

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"
	"ez-admin-gin/server/internal/token" // [!code ++]

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
```

更新 `LoginHandler` 和构造函数：

```go
// LoginHandler 负责登录相关接口。
type LoginHandler struct {
	db           *gorm.DB
	log          *zap.Logger
	tokenManager *token.Manager // [!code ++]
}

// NewLoginHandler 创建登录 Handler。
func NewLoginHandler(db *gorm.DB, log *zap.Logger, tokenManager *token.Manager) *LoginHandler {
	return &LoginHandler{
		db:           db,
		log:          log,
		tokenManager: tokenManager,
	}
}
```

更新登录响应结构：

```go
type loginResponse struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	AccessToken string `json:"access_token"` // [!code ++]
	TokenType   string `json:"token_type"` // [!code ++]
	ExpiresAt   string `json:"expires_at"` // [!code ++]
}
```

在密码校验通过后，生成 Token：

```go
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, apperror.Unauthorized("用户名或密码错误"), h.log)
		return
	}

	accessToken, expiresAt, err := h.tokenManager.GenerateAccessToken(user.ID, user.Username) // [!code ++]
	if err != nil { // [!code ++]
		response.Error(c, apperror.Internal("登录失败", err), h.log) // [!code ++]
		return // [!code ++]
	} // [!code ++]

	response.Success(c, loginResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		AccessToken: accessToken, // [!code ++]
		TokenType:   "Bearer", // [!code ++]
		ExpiresAt:   expiresAt.UTC().Format(time.RFC3339), // [!code ++]
	})
```

::: details 为什么返回 `token_type`
后续访问受保护接口时，通常会把 Token 放在请求头里：

```http
Authorization: Bearer <access_token>
```

`token_type` 返回 `Bearer`，可以让前端明确知道应该按哪种认证方式拼接请求头。
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

如果启动失败并提示 `jwt secret must be at least 32 characters`，说明 `auth.jwt_secret` 太短，改成更长的随机字符串后重新启动。

## ✅ 验证登录返回 Token

保持后端服务运行，调用登录接口：

::: code-group

```powershell [Windows PowerShell]
# 登录并保存响应
$body = @{
  username = "admin"
  password = "EzAdmin@123456"
} | ConvertTo-Json

$response = Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/auth/login `
  -ContentType "application/json" `
  -Body $body

$response.data
```

```bash [macOS / Linux]
# 登录并查看响应
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"EzAdmin@123456"}'
```

:::

应该看到类似结果：

```json
{
  "user_id": 1,
  "username": "admin",
  "nickname": "系统管理员",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_at": "2026-04-22T12:00:00Z"
}
```

::: details 怎么快速判断是不是 JWT
JWT 通常由三段组成，中间用 `.` 分隔：

```text
header.payload.signature
```

所以 `access_token` 看起来应该包含两个点。下一节会把它放进 `Authorization` 请求头，并由认证中间件真正解析校验。
:::

## 常见问题

::: details 提示 `jwt secret must be at least 32 characters`
说明 `auth.jwt_secret` 太短。把 `server/configs/config.yaml` 中的 `auth.jwt_secret` 改成至少 32 个字符的字符串。

生产环境推荐使用环境变量覆盖：

```powershell [Windows PowerShell]
# 当前终端临时覆盖 JWT 密钥
$env:EZ_AUTH_JWT_SECRET = "replace-with-a-long-random-secret-value"
```

```bash [macOS / Linux]
# 当前终端临时覆盖 JWT 密钥
export EZ_AUTH_JWT_SECRET="replace-with-a-long-random-secret-value"
```
:::

::: details 登录成功但没有 `access_token`
先确认 `login.go` 中已经调用了 `GenerateAccessToken`，并且 `loginResponse` 已经增加 `AccessToken`、`TokenType`、`ExpiresAt` 三个字段。
:::

::: details 提示 `no required module provides package github.com/golang-jwt/jwt/v5`
说明 JWT 依赖还没有加入当前 module。回到 `server/` 目录执行：

```bash
# 安装并整理依赖
go get github.com/golang-jwt/jwt/v5@latest
go mod tidy
```
:::

下一节会用这个 Token 保护接口：[认证中间件](./auth-middleware)。
