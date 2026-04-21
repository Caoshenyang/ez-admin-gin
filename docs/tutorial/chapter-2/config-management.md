---
title: 配置管理
description: "设计后端配置文件、环境变量和不同运行环境的加载方式。"
---

# 配置管理

第一章已经让后端服务可以启动。现在把端口、运行环境、数据库和 Redis 这些信息从代码里移到配置里，后续接入日志、数据库和缓存时，就不用把参数写死在业务代码中。

::: tip 🎯 本节目标
完成后，后端会从 `configs/config.yaml` 读取配置，并允许用环境变量覆盖配置值。
:::

## 配置放在哪里

本节会新增下面这些文件：

```text
server/
├─ configs/
│  └─ config.yaml
└─ internal/
   └─ config/
      └─ config.go
```

| 位置 | 用途 |
| --- | --- |
| `configs/config.yaml` | 本地开发默认配置 |
| `internal/config/config.go` | 读取配置、设置默认值、支持环境变量覆盖 |

::: warning ⚠️ 不要把生产密钥写进配置文件
`config.yaml` 只放本地开发能公开的默认值。生产环境的密码、密钥、访问令牌，应该通过环境变量或部署平台的 Secret 管理。
:::

## 🛠️ 安装配置库

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

安装 Viper：

```bash
# 安装配置读取依赖
go get github.com/spf13/viper@latest
```

Viper 用来读取 YAML 配置，并支持环境变量覆盖。这里先只用它处理基础配置，不引入更复杂的配置中心。

依赖资料入口：

| 依赖 | 用途 | 资料 |
| --- | --- | --- |
| `github.com/spf13/viper` | 读取配置文件和环境变量 | [Go 包文档](https://pkg.go.dev/github.com/spf13/viper) / [项目仓库](https://github.com/spf13/viper) |

## 🛠️ 创建配置文件

创建 `server/configs/config.yaml`：

```yaml
app:
  # 应用名称，后续日志和部署信息会用到。
  name: ez-admin
  # 当前运行环境，本地开发先使用 dev。
  env: dev

server:
  # HTTP 服务监听地址。
  addr: ":8080"

database:
  # 数据库配置和第一章的 Docker Compose 保持一致。
  host: localhost
  port: 5432
  user: ez_admin
  password: ez_admin_123456
  name: ez_admin

redis:
  # Redis 配置和第一章的 Docker Compose 保持一致。
  host: localhost
  port: 6379
  password: ""
  db: 0
```

这些值和第一章的 `deploy/compose.local.yml` 保持一致，后续数据库和 Redis 连接会直接复用。

## 🛠️ 创建配置加载代码

创建 `server/internal/config/config.go`：

```go
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 汇总整个服务端会读取的配置段。
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

// AppConfig 保存应用自身信息。
type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

// ServerConfig 保存 HTTP 服务启动配置。
type ServerConfig struct {
	Addr string `mapstructure:"addr"`
}

// DatabaseConfig 保存数据库连接配置。
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

// RedisConfig 保存 Redis 连接配置。
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Load 读取配置文件，并把结果解析到 Config 结构体中。
func Load() (*Config, error) {
	v := viper.New()

	// 配置文件位置是 server/configs/config.yaml。
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	// 先设置默认值，再绑定环境变量。
	setDefaults(v)
	bindEnvs(v)

	// EZ_SERVER_ADDR 这类环境变量会覆盖 server.addr。
	v.SetEnvPrefix("EZ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置兜底值，避免配置文件缺少字段时直接变成零值。
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "ez-admin")
	v.SetDefault("app.env", "dev")
	v.SetDefault("server.addr", ":8080")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "ez_admin")
	v.SetDefault("database.password", "ez_admin_123456")
	v.SetDefault("database.name", "ez_admin")
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
}

// bindEnvs 让环境变量能稳定参与结构体解析。
func bindEnvs(v *viper.Viper) {
	keys := []string{
		"app.name",
		"app.env",
		"server.addr",
		"database.host",
		"database.port",
		"database.user",
		"database.password",
		"database.name",
		"redis.host",
		"redis.port",
		"redis.password",
		"redis.db",
	}

	for _, key := range keys {
		// BindEnv 返回错误通常来自 key 本身，这里的 key 是固定列表。
		_ = v.BindEnv(key)
	}
}
```

::: details 为什么要写 `bindEnvs`
环境变量覆盖要能稳定参与结构体解析。这里把配置项逐个绑定到环境变量，避免出现“环境变量设置了，但结构体里没读到”的情况。
:::

## 🛠️ 使用配置启动服务

修改 `server/main.go`。这一处重点看三个变化：

- 引入 `internal/config` 包。
- 启动地址从 `cfg.Server.Addr` 读取。
- 健康检查接口返回当前运行环境。

```go
package main

import (
	"log"

	"ez-admin-gin/server/internal/config" // [!code ++]

	"github.com/gin-gonic/gin"
)

func main() {
	// 启动服务前先加载配置。
	cfg, err := config.Load() // [!code ++]
	if err != nil {
		log.Fatalf("load config: %v", err) // [!code ++]
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			// 返回 env，方便验证配置文件已经被读取。
			"env":    cfg.App.Env, // [!code ++]
		})
	})

	// 服务端口不再写死，改为读取配置文件。
	if err := r.Run(cfg.Server.Addr); err != nil { // [!code ++]
		log.Fatalf("run server: %v", err) // [!code ++]
	}
}
```

这里先只把 `server.addr` 和 `app.env` 接进启动流程，数据库和 Redis 配置会在后续章节使用。

## 配置优先级

本节采用下面的覆盖顺序：

| 优先级 | 来源 | 示例 |
| --- | --- | --- |
| 低 | 代码默认值 | `server.addr = :8080` |
| 中 | `configs/config.yaml` | `server.addr: ":8080"` |
| 高 | 环境变量 | `EZ_SERVER_ADDR=:18080` |

环境变量名称规则：

| 配置项 | 环境变量 |
| --- | --- |
| `app.env` | `EZ_APP_ENV` |
| `server.addr` | `EZ_SERVER_ADDR` |
| `database.host` | `EZ_DATABASE_HOST` |
| `redis.port` | `EZ_REDIS_PORT` |

## ✅ 验证默认配置

整理依赖：

```bash
# 整理新增依赖
go mod tidy
```

启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

访问健康检查接口：

::: code-group

```powershell [Windows PowerShell]
# 访问默认端口的健康检查接口
Invoke-RestMethod http://localhost:8080/health
```

```bash [macOS / Linux]
# 访问默认端口的健康检查接口
curl http://localhost:8080/health
```

:::

应该能看到 `status` 为 `ok`，`env` 为 `dev`。

验证完成后，回到运行 `go run .` 的终端，按 `Ctrl + C` 停止服务。

## ✅ 验证环境变量覆盖

把服务端口临时改成 `18080`：

::: code-group

```powershell [Windows PowerShell]
# 临时把服务端口覆盖为 18080
$env:EZ_SERVER_ADDR = ":18080"
go run .
```

```bash [macOS / Linux]
# 临时把服务端口覆盖为 18080
EZ_SERVER_ADDR=:18080 go run .
```

:::

再访问：

::: code-group

```powershell [Windows PowerShell]
# 访问被环境变量覆盖后的端口
Invoke-RestMethod http://localhost:18080/health
```

```bash [macOS / Linux]
# 访问被环境变量覆盖后的端口
curl http://localhost:18080/health
```

:::

如果能正常返回，说明环境变量已经覆盖了配置文件里的端口。

::: warning ⚠️ PowerShell 环境变量只影响当前窗口
`$env:EZ_SERVER_ADDR = ":18080"` 只在当前 PowerShell 窗口生效。验证完成后可以执行：

```powershell
# 清理临时环境变量
Remove-Item Env:EZ_SERVER_ADDR
```
:::

## ✅ 确认 Git 状态

回到项目根目录：

::: code-group

```powershell [Windows PowerShell]
# 回到项目根目录，查看本节改动
Set-Location ..
git status
```

```bash [macOS / Linux]
# 回到项目根目录，查看本节改动
cd ..
git status
```

:::

应该能看到：

```text
server/configs/config.yaml
server/internal/config/config.go
server/main.go
server/go.mod
server/go.sum
```

下一节开始接入日志系统：[日志系统](./logging-system)。
