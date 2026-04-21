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
Set-Location .\server
```

```bash [macOS / Linux]
cd server
```

:::

安装 Viper：

```bash
go get github.com/spf13/viper@latest
```

Viper 用来读取 YAML 配置，并支持环境变量覆盖。这里先只用它处理基础配置，不引入更复杂的配置中心。

## 🛠️ 创建配置文件

创建 `server/configs/config.yaml`：

```yaml
app:
  name: ez-admin
  env: dev

server:
  addr: ":8080"

database:
  host: localhost
  port: 5432
  user: ez_admin
  password: ez_admin_123456
  name: ez_admin

redis:
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

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type ServerConfig struct {
	Addr string `mapstructure:"addr"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	setDefaults(v)
	bindEnvs(v)

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
		_ = v.BindEnv(key)
	}
}
```

::: details 为什么要写 `bindEnvs`
环境变量覆盖要能稳定参与结构体解析。这里把配置项逐个绑定到环境变量，避免出现“环境变量设置了，但结构体里没读到”的情况。
:::

## 🛠️ 使用配置启动服务

修改 `server/main.go`：

```go
package main

import (
	"log"

	"ez-admin-gin/server/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    cfg.App.Env,
		})
	})

	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("run server: %v", err)
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
go mod tidy
```

启动服务：

```bash
go run .
```

访问健康检查接口：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod http://localhost:8080/health
```

```bash [macOS / Linux]
curl http://localhost:8080/health
```

:::

应该能看到 `status` 为 `ok`，`env` 为 `dev`。

验证完成后，回到运行 `go run .` 的终端，按 `Ctrl + C` 停止服务。

## ✅ 验证环境变量覆盖

把服务端口临时改成 `18080`：

::: code-group

```powershell [Windows PowerShell]
$env:EZ_SERVER_ADDR = ":18080"
go run .
```

```bash [macOS / Linux]
EZ_SERVER_ADDR=:18080 go run .
```

:::

再访问：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod http://localhost:18080/health
```

```bash [macOS / Linux]
curl http://localhost:18080/health
```

:::

如果能正常返回，说明环境变量已经覆盖了配置文件里的端口。

::: warning ⚠️ PowerShell 环境变量只影响当前窗口
`$env:EZ_SERVER_ADDR = ":18080"` 只在当前 PowerShell 窗口生效。验证完成后可以执行：

```powershell
Remove-Item Env:EZ_SERVER_ADDR
```
:::

## ✅ 确认 Git 状态

回到项目根目录：

::: code-group

```powershell [Windows PowerShell]
Set-Location ..
git status
```

```bash [macOS / Linux]
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
