---
title: Redis 连接
description: "接入 Redis，为缓存、会话扩展和运行期能力做准备。"
---

# Redis 连接

数据库已经接进后端服务了。现在继续把 Redis 接进来，让后端具备缓存、会话、验证码、限流这类运行期能力的基础入口。

::: tip 🎯 本节目标
完成后，服务启动时会连接 Redis，控制台会输出 Redis 连接日志，`/health` 接口会返回 Redis 状态。
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
│  └─ redis/
│     └─ redis.go
└─ main.go
```

| 位置 | 用途 |
| --- | --- |
| `configs/config.yaml` | 增加 Redis 连接池配置 |
| `internal/config/config.go` | 扩展 Redis 配置结构和环境变量绑定 |
| `internal/redis/redis.go` | 创建 Redis 客户端、提供 Ping 和关闭方法 |
| `main.go` | 启动时连接 Redis，并把 Redis 状态接入 `/health` |

::: warning ⚠️ 先确认 Redis 已启动
继续之前，确认第一章的基础环境已经启动：

```bash
# 在项目根目录查看 PostgreSQL 和 Redis 状态
docker compose -f deploy/compose.local.yml ps
```

`redis` 应该处于 `running` 或 `healthy` 状态。
:::

## 🛠️ 安装 Redis 依赖

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

安装 Redis 客户端：

```bash
# 安装 Redis 客户端依赖
go get github.com/redis/go-redis/v9@latest
```

依赖资料入口：

| 依赖 | 用途 | 资料 |
| --- | --- | --- |
| `github.com/redis/go-redis/v9` | Redis 客户端，负责连接、命令执行和连接池管理 | [Go 包文档](https://pkg.go.dev/github.com/redis/go-redis/v9) / [项目仓库](https://github.com/redis/go-redis) |

::: details 为什么这里直接用 go-redis
Redis 本身就是 Key-Value 服务，不像数据库那样需要 ORM。这里直接使用官方常用客户端就够了，后面做缓存封装时也会更直接。
:::

## 🛠️ 扩展 Redis 配置

修改 `server/configs/config.yaml`，给 `redis` 增加连接池配置：

```yaml
redis:
  # Redis 配置和第一章的 Docker Compose 保持一致。
  host: localhost
  port: 6379
  password: ""
  db: 0
  # 最大重试次数，临时网络抖动时会自动重试。
  max_retries: 3 # [!code ++]
  # 最小空闲连接数，避免每次请求都重新建连。
  min_idle_conns: 5 # [!code ++]
  # 连接池大小，限制最多保留多少个 Redis 连接。
  pool_size: 10 # [!code ++]
```

字段含义：

| 字段 | 说明 |
| --- | --- |
| `max_retries` | 命令失败时最多自动重试多少次 |
| `min_idle_conns` | 最少保留多少个空闲连接 |
| `pool_size` | Redis 连接池大小 |

## 🛠️ 扩展配置结构

修改 `server/internal/config/config.go`。这一处有三个改动：

- 给 `RedisConfig` 增加连接池字段。
- 给 Redis 连接池补默认值。
- 给 Redis 连接池补环境变量绑定。

先修改 `RedisConfig`：

```go
// RedisConfig 保存 Redis 连接配置。
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	// MaxRetries 控制 Redis 命令失败时的最大重试次数。
	MaxRetries int `mapstructure:"max_retries"` // [!code ++]
	// MinIdleConns 控制最少保留多少个空闲连接。
	MinIdleConns int `mapstructure:"min_idle_conns"` // [!code ++]
	// PoolSize 控制连接池大小。
	PoolSize int `mapstructure:"pool_size"` // [!code ++]
}
```

继续在 `setDefaults` 里增加默认值：

```go
// Redis 连接池默认值适合本地开发和小型后台起步。
v.SetDefault("redis.max_retries", 3) // [!code ++]
v.SetDefault("redis.min_idle_conns", 5) // [!code ++]
v.SetDefault("redis.pool_size", 10) // [!code ++]
```

最后在 `bindEnvs` 的 `keys` 列表中增加：

```go
// 允许用 EZ_REDIS_POOL_SIZE 这类环境变量覆盖 Redis 连接池配置。
"redis.max_retries", // [!code ++]
"redis.min_idle_conns", // [!code ++]
"redis.pool_size", // [!code ++]
```

常用环境变量示例：

| 配置项 | 环境变量 |
| --- | --- |
| `redis.host` | `EZ_REDIS_HOST` |
| `redis.port` | `EZ_REDIS_PORT` |
| `redis.db` | `EZ_REDIS_DB` |
| `redis.pool_size` | `EZ_REDIS_POOL_SIZE` |

## 🛠️ 创建 Redis 包

创建 `server/internal/redis/redis.go`。这是新增文件，直接完整写入即可。

```go
package redis

import (
	"context"
	"fmt"
	"time"

	"ez-admin-gin/server/internal/config"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var pingTimeout = 3 * time.Second

// New 创建 Redis 客户端，并在启动时完成连通性检查。
func New(cfg config.RedisConfig, log *zap.Logger) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
	})

	if err := Ping(client); err != nil {
		return nil, err
	}

	log.Info(
		"redis connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	return client, nil
}

// Ping 用于健康检查，确认 Redis 当前仍然可连接。
func Ping(client *goredis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	return nil
}

// Close 关闭 Redis 客户端。
func Close(client *goredis.Client) error {
	if err := client.Close(); err != nil {
		return fmt.Errorf("close redis: %w", err)
	}

	return nil
}
```

::: details 为什么 `Ping` 要加超时
健康检查应该尽快返回。如果 Redis 因为网络问题卡住，超时能避免接口长时间挂起。
:::

## 🛠️ 在启动入口接入 Redis

修改 `server/main.go`。这一处重点看四个变化：

- 引入 `internal/redis` 包。
- 服务启动时创建 Redis 客户端，并在退出时关闭它。
- `/health` 增加 Redis 连通性检查。
- 健康检查结果同时返回数据库和 Redis 状态。

```go
package main

import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"
	"net/http"

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	appRedis "ez-admin-gin/server/internal/redis" // [!code ++]

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
	redisClient, err := appRedis.New(cfg.Redis, log) // [!code ++]
	if err != nil {
		log.Fatal("connect redis", zap.Error(err)) // [!code ++]
	}
	defer func() {
		if err := appRedis.Close(redisClient); err != nil { // [!code ++]
			log.Error("close redis", zap.Error(err)) // [!code ++]
		}
	}()

	// 使用 gin.New()，再手动挂载自定义中间件。
	r := gin.New()
	r.Use(appLogger.GinLogger(log), appLogger.GinRecovery(log))

	r.GET("/health", func(c *gin.Context) {
		if err := database.Ping(db); err != nil {
			log.Error("database health check failed", zap.Error(err))
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error",
				"env":      cfg.App.Env,
				"database": "unavailable",
				"redis":    "unknown", // [!code ++]
			})
			return
		}

		if err := appRedis.Ping(redisClient); err != nil { // [!code ++]
			log.Error("redis health check failed", zap.Error(err)) // [!code ++]
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error", // [!code ++]
				"env":      cfg.App.Env,
				"database": "ok", // [!code ++]
				"redis":    "unavailable", // [!code ++]
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"env":      cfg.App.Env,
			"database": "ok",
			"redis":    "ok", // [!code ++]
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

## ✅ 整理依赖并启动

整理依赖：

```bash
# 整理新增依赖，更新 go.mod 和 go.sum
go mod tidy
```

确认 Redis 正在运行：

```bash
# 在项目根目录执行
docker compose -f deploy/compose.local.yml ps
```

回到 `server/` 目录启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

启动后，控制台应该能看到类似日志：

```text
INFO	database connected	{"host": "localhost", "port": 5432, "database": "ez_admin"}
INFO	redis connected	{"host": "localhost", "port": 6379, "db": 0}
INFO	server started	{"addr": ":8080", "env": "dev"}
```

## ✅ 验证健康检查

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

应该看到：

```json
{
  "database": "ok",
  "env": "dev",
  "redis": "ok",
  "status": "ok"
}
```

这说明后端服务、PostgreSQL 和 Redis 都处于可用状态。

## ✅ 验证 Redis 不可用的情况

停止 Redis：

```bash
# 在项目根目录停止 Redis 容器
docker compose -f deploy/compose.local.yml stop redis
```

再次启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

这时服务应该启动失败，并在日志里看到 `connect redis` 相关错误。验证完成后，重新启动 Redis：

```bash
# 在项目根目录重新启动 Redis 容器
docker compose -f deploy/compose.local.yml start redis
```

::: warning ⚠️ 验证失败场景后记得恢复服务
后续章节会继续依赖 Redis。验证 Redis 不可用的情况后，记得重新执行 `start redis`，并确认容器状态恢复正常。
:::

## 常见问题

::: details 提示 `connect: connection refused`
通常是 Redis 容器没有启动，或者端口不是 `6379`。先执行：

```bash
# 查看 Redis 容器状态和端口
docker compose -f deploy/compose.local.yml ps
```
:::

::: details 提示 `NOAUTH Authentication required`
说明 Redis 已经设置了密码，但 `configs/config.yaml` 里的 `redis.password` 还是空字符串。检查配置是否和实际 Redis 一致。
:::

::: details 修改了 Redis 配置但没有生效
如果使用了环境变量，环境变量优先级高于 `config.yaml`。例如 `EZ_REDIS_HOST` 会覆盖 `redis.host`。
:::

下一节开始统一响应与错误处理：[统一响应与错误处理](./response-and-errors)。
