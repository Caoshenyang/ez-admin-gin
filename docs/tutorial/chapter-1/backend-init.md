---
title: Go 后端项目初始化
description: "初始化 server 子项目，为 Gin 后端服务准备最小可运行基础。"
---

# Go 后端项目初始化

上一节已经创建了 `server/` 目录。现在进入这个目录，初始化后端 module，并写一个最小可运行的健康检查接口。

::: tip 🎯 本节目标
让 `server/` 从空目录变成一个能启动的最小后端服务，并通过 `/health` 接口确认服务可用。
:::

完成后，`server/` 目录会变成这样：

```text
server/
├─ go.mod
├─ go.sum
└─ main.go
```

## 🛠️ 环境要求

本教程使用当前最新稳定版 Go。截止 2026-04-21，Go 官方下载页的最新稳定版是 `go1.26.2`。

先确认本机版本：

```bash
# 确认当前 Go 版本
go version
```

如果还没有安装 Go，先到官方页面下载安装：

- Go 下载页：[https://go.dev/dl/](https://go.dev/dl/)
- Go 安装说明：[https://go.dev/doc/install](https://go.dev/doc/install)

安装完成后重新打开终端，再执行 `go version`。能看到 `go1.26.2` 或更新的稳定版即可继续。

## 🛠️ 初始化 module

进入 `server/` 目录：

::: code-group

```powershell [Windows PowerShell]
# 进入后端项目目录
Set-Location .\server
```

```bash [macOS / Linux]
# 进入后端项目目录
cd server
```

:::

初始化 module：

```bash
# 初始化后端 module，module 路径和仓库目录保持一致
go mod init ez-admin-gin/server
```

然后把 `go.mod` 明确更新到本教程使用的 Go 版本：

```bash
# 明确使用本教程对应的 Go 版本
go get go@1.26.2
```

::: details 关于 `go.mod` 里的版本
从 Go 1.26 开始，`go mod init` 会默认写入低一档的 `go` 版本，例如使用 `go1.26.2` 初始化时，可能先生成 `go 1.25.0`。

这是官方为了让新 module 默认兼容当前仍被支持的 Go 版本。这里执行 `go get go@1.26.2`，是为了让本教程的 `go.mod` 明确使用当前版本。
:::

## 🛠️ 添加最小入口

创建 `main.go`：

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	// gin.Default() 会附带默认日志和 panic 恢复中间件。
	r := gin.Default()

	// 健康检查接口，用来确认服务已经正常启动。
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 启动 HTTP 服务，默认监听本机 8080 端口。
	r.Run(":8080")
}
```

安装依赖并整理 module：

```bash
# 安装 Gin，并整理 go.mod / go.sum
go get github.com/gin-gonic/gin@latest
go mod tidy
```

依赖资料入口：

| 依赖 | 用途 | 资料 |
| --- | --- | --- |
| `github.com/gin-gonic/gin` | 创建 HTTP 服务和路由 | [Go 包文档](https://pkg.go.dev/github.com/gin-gonic/gin) / [项目仓库](https://github.com/gin-gonic/gin) |

如果上一节留下了 `.gitkeep`，现在可以删掉：

::: code-group

```powershell [Windows PowerShell]
# server/ 已经有真实文件，可以删除占位文件
Remove-Item .gitkeep -ErrorAction SilentlyContinue
```

```bash [macOS / Linux]
# server/ 已经有真实文件，可以删除占位文件
rm -f .gitkeep
```

:::

## ✅ 启动并验证

在 `server/` 目录下启动服务：

```bash
# 在 server/ 目录启动当前 module
go run .
```

看到类似下面的输出，就说明服务已经启动：

```text
Listening and serving HTTP on :8080
```

::: warning ⚠️ 端口被占用
如果提示 `8080` 端口被占用，先停止占用该端口的进程，或临时把 `main.go` 里的 `:8080` 改成其他端口再验证。
:::

保持服务运行，打开另一个终端访问健康检查接口：

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
  "status": "ok"
}
```

验证完成后，回到运行 `go run .` 的终端，按 `Ctrl + C` 停止服务。

下一节开始初始化管理台项目：[Vue 管理台项目初始化](./admin-init)。
