---
title: 配置管理
description: "设计后端配置文件、环境变量和当前启动入口的配置加载方式。"
---

# 配置管理

第一章已经让后端服务能跑起来。这一节把运行参数从代码里抽出来，让端口、环境、数据库、Redis、日志等配置都从 `configs/config.yaml` 和环境变量读取。

::: tip 先记主线
当前代码里，配置不是在业务 Handler 中零散读取，也不是在旧版 `cmd/server/main.go` 里手动解析，而是由 `server/internal/platform/config` 统一加载，再由 `server/internal/bootstrap/run.go` 分发给日志、数据库、Redis、鉴权和路由层。
:::

## 配置放在哪里

本节重点关注这几个位置：

```text
server/
├─ configs/
│  └─ config.yaml
├─ main.go
└─ internal/
   ├─ bootstrap/
   │  └─ run.go
   └─ platform/
      └─ config/
         └─ config.go
```

| 位置 | 用途 |
| --- | --- |
| `configs/config.yaml` | 本地开发默认配置 |
| `internal/platform/config/config.go` | 统一读取配置、设置默认值、绑定环境变量 |
| `internal/bootstrap/run.go` | 启动时读取配置并把它分发给各平台能力 |
| `server/main.go` | 只保留 `bootstrap.MustRun(...)`，不再内联配置加载流程 |

::: warning 不要再把配置主线理解成 `internal/config -> main.go -> gin.Run`
仓库里仍然保留了 `server/internal/config` 这一层兼容落点，但当前启动主线已经切到 `internal/platform/config` 和 `internal/bootstrap/run.go`。阅读和继续扩展代码时，应优先以这条主线为准。
:::

## 配置文件长什么样

当前服务端默认读取 `server/configs/config.yaml`。基础字段至少包括：

```yaml
app:
  name: ez-admin
  env: dev

server:
  addr: ":8080"

database:
  driver: postgres
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

后续章节会继续补充 `log`、`auth`、`upload` 等配置段，但入口规则不变：都由同一个配置加载器处理。

## 当前项目如何加载配置

当前启动不是在 `main.go` 里直接 `config.Load()`，而是先进入 `bootstrap.MustRun(...)`，再在 `run.go` 里统一加载：

```go
func MustRun(migrationsFS fs.FS, rbacModelPath string) {
	cfg, err := platformConfig.Load()
	if err != nil {
		stdlog.Fatalf("load config: %v", err)
	}

	log, err := platformLogger.New(cfg.Log)
	if err != nil {
		stdlog.Fatalf("create logger: %v", err)
	}

	// 后续数据库、Redis、鉴权、路由都继续复用 cfg
}
```

这层分工的意思是：

- `main.go` 只负责进入启动流程
- `platform/config` 负责把配置读出来
- `bootstrap/run.go` 负责决定配置在启动链上的使用顺序

## `platform/config` 负责什么

当前 `server/internal/platform/config/config.go` 提供的是统一配置入口。它做的事情包括：

- 从 `./configs/config.yaml` 读取 YAML
- 为关键字段设置默认值
- 绑定 `EZ_` 前缀环境变量
- 把结果解析成 `Config` 结构体

典型代码结构可以概括成这样：

```go
v.SetConfigName("config")
v.SetConfigType("yaml")
v.AddConfigPath("./configs")

setDefaults(v)
bindEnvs(v)

v.SetEnvPrefix("EZ")
v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
v.AutomaticEnv()
```

## 配置优先级

当前采用的覆盖顺序是：

| 优先级 | 来源 | 示例 |
| --- | --- | --- |
| 低 | 代码默认值 | `server.addr = :8080` |
| 中 | `configs/config.yaml` | `server.addr: ":8080"` |
| 高 | 环境变量 | `EZ_SERVER_ADDR=:18080` |

常见环境变量映射：

| 配置项 | 环境变量 |
| --- | --- |
| `app.env` | `EZ_APP_ENV` |
| `server.addr` | `EZ_SERVER_ADDR` |
| `database.host` | `EZ_DATABASE_HOST` |
| `redis.port` | `EZ_REDIS_PORT` |

## ✅ 验证默认配置

在 `server/` 目录启动服务：

```bash
go run .
```

访问健康检查：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod http://localhost:8080/health
```

```bash [macOS / Linux]
curl http://localhost:8080/health
```

:::

当前返回已经是统一响应格式，而不是早期章节里的裸 JSON：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "database": "ok",
    "env": "dev",
    "redis": "ok"
  }
}
```

## ✅ 验证环境变量覆盖

临时把端口改成 `18080`：

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

如果能正常返回，说明环境变量覆盖已经生效。

::: warning PowerShell 环境变量只影响当前窗口
验证完成后可以执行：

```powershell
Remove-Item Env:EZ_SERVER_ADDR
```
:::

## 这一页需要同步的认知

1. 当前入口是 `server/main.go`，不是 `cmd/server/main.go`。
2. 当前配置加载入口是 `internal/platform/config`，不是把读取逻辑散在旧目录里。
3. 当前配置的消费者是 `bootstrap/run.go`，不是在 `main.go` 里手写一整套初始化。
4. 当前健康检查返回结构已经是统一响应格式，不再只是 `{"status":"ok"}`。

下一节继续接入日志系统：[日志系统](./logging-system)。
