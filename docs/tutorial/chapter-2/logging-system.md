---
title: 日志系统
description: "接入 Zap 日志，为请求、错误和关键业务行为提供统一记录。"
---

# 日志系统

上一节已经把运行参数交给配置文件管理。这一节接入结构化日志，让服务启动、请求访问和异常恢复都有统一记录。

::: tip 🎯 本节目标
完成后，服务启动时会输出结构化日志，每次访问接口都会记录请求日志，同时日志会写入 `server/logs/app.log`。
:::

## 日志放在哪里

本节会新增或修改下面这些文件：

```text
server/
├─ configs/
│  └─ config.yaml
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  └─ logger/
│     └─ logger.go
└─ main.go
```

| 位置 | 用途 |
| --- | --- |
| `configs/config.yaml` | 增加日志级别、格式和文件位置 |
| `internal/config/config.go` | 增加日志配置结构 |
| `internal/logger/logger.go` | 初始化 Zap，并提供请求日志和异常恢复中间件 |
| `main.go` | 使用自定义日志替换默认启动方式 |

::: warning ⚠️ 避免重复请求日志
`gin.Default()` 已经自带默认日志和恢复中间件。接入自定义日志后，改用 `gin.New()`，再手动挂载本节创建的日志中间件。
:::

## 🛠️ 安装日志依赖

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

安装 Zap 和日志切割库：

```bash
# 安装结构化日志和日志切割依赖
go get go.uber.org/zap@latest gopkg.in/natefinch/lumberjack.v2@latest
```

这两个依赖的资料入口：

| 依赖 | 用途 | 资料 |
| --- | --- | --- |
| `go.uber.org/zap` | 输出结构化、分级日志 | [Go 包文档](https://pkg.go.dev/go.uber.org/zap) / [项目仓库](https://github.com/uber-go/zap) |
| `gopkg.in/natefinch/lumberjack.v2` | 日志文件切割和保留 | [Go 包文档](https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2) / [项目仓库](https://github.com/natefinch/lumberjack) |

::: details 为什么这里不用标准库日志
标准库日志适合很简单的输出。后台底座后续会记录请求、用户操作、数据库错误和任务执行结果，结构化字段会更方便排查问题。
:::

## 🛠️ 增加日志配置

修改 `server/configs/config.yaml`，在文件末尾增加：

```yaml
log: # [!code ++]
  # 当前日志级别，开发阶段通常使用 info，需要更详细日志时可以临时改成 debug
  level: info # [!code ++]
  # console 更适合本地阅读；生产环境如果要给日志平台采集，可以改成 json
  format: console # [!code ++]
  # 日志文件位置，相对于 server/ 目录
  filename: logs/app.log # [!code ++]
  # 单个日志文件最大大小，单位 MB
  max_size: 100 # [!code ++]
  # 最多保留多少个旧日志文件
  max_backups: 7 # [!code ++]
  # 日志文件最多保留多少天
  max_age: 30 # [!code ++]
  # 是否压缩旧日志文件
  compress: false # [!code ++]
```

字段含义：

| 字段 | 说明 |
| --- | --- |
| `level` | 日志级别，常用值为 `debug`、`info`、`warn`、`error` |
| `format` | 日志格式，开发阶段使用 `console`，生产环境常用 `json` |
| `filename` | 日志文件位置，相对于 `server/` 目录 |
| `max_size` | 单个日志文件最大大小，单位 MB |
| `max_backups` | 最多保留多少个旧日志文件 |
| `max_age` | 日志最多保留多少天 |
| `compress` | 是否压缩旧日志文件 |

::: warning ⚠️ 日志文件不要提交到 Git
项目根目录的 `.gitignore` 已经忽略了 `*.log`。后续生成的 `logs/app.log` 是运行产物，不需要提交。
:::

## 🛠️ 扩展配置结构

修改 `server/internal/config/config.go`。这一处有四个改动：

- 给总配置 `Config` 增加 `Log` 字段。
- 新增 `LogConfig` 结构体。
- 给日志配置补默认值。
- 给日志配置补环境变量绑定。

先给 `Config` 增加 `Log` 字段：

```go{6}
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	// Log 对应 config.yaml 中的 log 配置段。
	Log      LogConfig      `mapstructure:"log"` // [!code ++]
}
```

再新增日志配置结构：

```go
type LogConfig struct {
	// Level 控制输出哪些级别的日志。
	Level string `mapstructure:"level"` // [!code ++]
	// Format 控制日志格式，支持 console 和 json。
	Format string `mapstructure:"format"` // [!code ++]
	// Filename 是日志文件路径。为空时只输出到控制台。
	Filename string `mapstructure:"filename"` // [!code ++]
	// MaxSize 是单个日志文件最大大小，单位 MB。
	MaxSize int `mapstructure:"max_size"` // [!code ++]
	// MaxBackups 是最多保留的旧日志文件数量。
	MaxBackups int `mapstructure:"max_backups"` // [!code ++]
	// MaxAge 是日志文件最多保留天数。
	MaxAge int `mapstructure:"max_age"` // [!code ++]
	// Compress 控制是否压缩旧日志文件。
	Compress bool `mapstructure:"compress"` // [!code ++]
}
```

继续给 `setDefaults` 增加默认值：

```go
// 日志默认值和 config.yaml 保持一致，保证配置文件缺少字段时也能启动。
v.SetDefault("log.level", "info") // [!code ++]
v.SetDefault("log.format", "console") // [!code ++]
v.SetDefault("log.filename", "logs/app.log") // [!code ++]
v.SetDefault("log.max_size", 100) // [!code ++]
v.SetDefault("log.max_backups", 7) // [!code ++]
v.SetDefault("log.max_age", 30) // [!code ++]
v.SetDefault("log.compress", false) // [!code ++]
```

最后给 `bindEnvs` 增加环境变量绑定：

```go
// 绑定环境变量后，可以用 EZ_LOG_LEVEL 这类变量覆盖配置文件。
"log.level", // [!code ++]
"log.format", // [!code ++]
"log.filename", // [!code ++]
"log.max_size", // [!code ++]
"log.max_backups", // [!code ++]
"log.max_age", // [!code ++]
"log.compress", // [!code ++]
```

这样后续可以用 `EZ_LOG_LEVEL=debug` 临时调整日志级别。

## 🛠️ 创建日志包

创建 `server/internal/logger/logger.go`。这是一个新增文件，直接完整写入即可。

```go
package logger

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ez-admin-gin/server/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// New 根据配置创建 Zap Logger。
// 这里同时处理日志级别、日志格式和输出位置。
func New(cfg config.LogConfig) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	// 生产配置的字段更稳定，适合后续接入日志采集平台。
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.MillisDurationEncoder

	var encoder zapcore.Encoder
	if strings.EqualFold(cfg.Format, "json") {
		// json 适合生产环境采集和检索。
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// console 适合本地开发时直接阅读。
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	writeSyncer, err := newWriteSyncer(cfg)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)

	return zap.New(
		core,
		// 记录调用位置，方便定位日志来自哪个文件。
		zap.AddCaller(),
		// error 级别自动带堆栈，方便排查异常。
		zap.AddStacktrace(zapcore.ErrorLevel),
	), nil
}

// GinLogger 记录每一次 HTTP 请求。
func GinLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 先执行后续处理器，才能拿到最终状态码和耗时。
		c.Next()

		if query != "" {
			path = path + "?" + query
		}

		// 用结构化字段记录请求信息，后续按字段过滤会更方便。
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", time.Since(start)),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		if c.Writer.Status() >= http.StatusInternalServerError {
			// 5xx 请求按错误日志记录。
			log.Error("http request", fields...)
			return
		}

		log.Info("http request", fields...)
	}
}

// GinRecovery 捕获 panic，避免服务因为单次请求异常直接退出。
func GinRecovery(log *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error(
			"panic recovered",
			zap.Any("error", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.Stack("stack"),
		)

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// parseLevel 把配置文件中的字符串转成 Zap 识别的日志级别。
func parseLevel(value string) (zapcore.Level, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(value))); err != nil {
		return level, fmt.Errorf("invalid log level %q: %w", value, err)
	}

	return level, nil
}

// newWriteSyncer 决定日志输出到哪里。
func newWriteSyncer(cfg config.LogConfig) (zapcore.WriteSyncer, error) {
	if cfg.Filename == "" {
		// 没有配置文件路径时，只输出到控制台。
		return zapcore.AddSync(os.Stdout), nil
	}

	// 日志目录不存在时自动创建，例如 logs/app.log 会先创建 logs/。
	if err := os.MkdirAll(filepath.Dir(cfg.Filename), 0o755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	// Lumberjack 负责日志切割，避免单个日志文件无限增长。
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// 同时输出到控制台和文件：开发时能直接看，事后也能查文件。
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(fileWriter),
	), nil
}
```

这个包做了三件事：

- 初始化 Zap 日志对象。
- 把日志同时输出到控制台和文件。
- 提供请求日志与异常恢复中间件。

## 🛠️ 在启动入口使用日志

修改 `server/main.go`。这一处重点看四个变化：

- 把 `log` 标准库改名为 `stdlog`，只用于日志初始化前的错误输出。
- 引入 `internal/logger` 和 `go.uber.org/zap`。
- 把 `gin.Default()` 改成 `gin.New()`，避免重复输出默认日志。
- 挂载自定义请求日志和异常恢复中间件。

```go
package main

import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log" // [!code ++]

	"ez-admin-gin/server/internal/config"
	appLogger "ez-admin-gin/server/internal/logger" // [!code ++]

	"github.com/gin-gonic/gin"
	"go.uber.org/zap" // [!code ++]
)

func main() {
	// 先读取配置，日志初始化也需要用到 cfg.Log。
	cfg, err := config.Load()
	if err != nil {
		stdlog.Fatalf("load config: %v", err)
	}

	// 根据配置创建结构化日志对象。
	log, err := appLogger.New(cfg.Log) // [!code ++]
	if err != nil {
		stdlog.Fatalf("create logger: %v", err) // [!code ++]
	}
	defer func() {
		// 退出前刷新缓冲区，避免最后几条日志丢失。
		_ = log.Sync() // [!code ++]
	}()

	// 使用 gin.New()，再手动挂载自定义中间件。
	r := gin.New() // [!code ++]
	r.Use(appLogger.GinLogger(log), appLogger.GinRecovery(log)) // [!code ++]

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    cfg.App.Env,
		})
	})

	// 服务启动日志记录关键运行参数。
	log.Info( // [!code ++]
		"server started",
		zap.String("addr", cfg.Server.Addr),
		zap.String("env", cfg.App.Env),
	)

	if err := r.Run(cfg.Server.Addr); err != nil {
		// Fatal 会记录日志并退出进程。
		log.Fatal("run server", zap.Error(err)) // [!code ++]
	}
}
```

这里保留了标准库 `log`，只用于日志系统还没创建成功之前的启动错误。创建成功后，统一使用 Zap。

## ✅ 整理依赖并启动

整理依赖：

```bash
# 整理新增依赖，更新 go.mod 和 go.sum
go mod tidy
```

启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

启动后，控制台应该能看到类似这样的日志：

```text
2026-04-21T16:40:00.000+0800	INFO	server started	{"addr": ":8080", "env": "dev"}
```

## ✅ 验证请求日志

访问健康检查接口：

::: code-group

```powershell [Windows PowerShell]
# 访问健康检查接口，触发一条请求日志
Invoke-RestMethod http://localhost:8080/health
```

```bash [macOS / Linux]
# 访问健康检查接口，触发一条请求日志
curl http://localhost:8080/health
```

:::

控制台应该新增一条 `http request` 日志，并包含 `status`、`method`、`path`、`latency` 等字段。

## ✅ 验证文件日志

查看日志文件：

::: code-group

```powershell [Windows PowerShell]
# 查看最近 5 行文件日志
Get-Content .\logs\app.log -Tail 5
```

```bash [macOS / Linux]
# 查看最近 5 行文件日志
tail -n 5 logs/app.log
```

:::

如果能看到 `server started` 和 `http request`，说明控制台和文件日志都已经生效。

::: details 如果没有生成 `logs/app.log`
先确认服务是在 `server/` 目录下启动的。如果在项目根目录执行 `go run ./server`，相对路径会变化，日志文件也会生成到不同位置。
:::

## ✅ 验证日志级别覆盖

停止当前服务，临时把日志级别改成 `debug`：

::: code-group

```powershell [Windows PowerShell]
# 临时覆盖日志级别，只影响当前 PowerShell 窗口
$env:EZ_LOG_LEVEL = "debug"
go run .
```

```bash [macOS / Linux]
# 临时覆盖日志级别，只影响当前命令
EZ_LOG_LEVEL=debug go run .
```

:::

这一步主要验证环境变量能够正常覆盖日志配置。验证完成后，PowerShell 可以清理临时环境变量：

```powershell
# 清理临时环境变量
Remove-Item Env:EZ_LOG_LEVEL
```

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
server/configs/config.yaml
server/internal/config/config.go
server/internal/logger/logger.go
server/main.go
server/go.mod
server/go.sum
```

下一节开始连接数据库：[数据库连接](./database-connection)。
