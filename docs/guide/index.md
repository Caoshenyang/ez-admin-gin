---
title: 快速启动
description: "用最短路径跑起 EZ Admin Gin，并理解它现在已经收敛成什么样的后台底座。"
---

# 快速启动

EZ Admin Gin 当前已经收敛成一条稳定主线：登录、权限、菜单、组织体系、数据权限、前端管理台、部署复用都已经在同一套单仓库里打通。

当前稳定版本：**v1.1.0**

::: tip 🎯 这页解决什么
先帮你判断这个项目是不是你要的，再用最短路径把它跑起来。
:::

## 先理解这套仓库的定位

这不是单独的后端模板，也不是只有页面壳子的前端 Demo，而是一套统一收在一个仓库里的后台底座：

- `server/` 提供 Go + Gin 后端能力
- `admin/` 提供 Vue 3 管理台
- `docs/` 提供教程、指南和参考手册
- `deploy/` 与 `scripts/` 提供部署和打包链路

如果你是第一次进这个仓库，建议固定按这条路径阅读：

1. 先跑起项目
2. 再看 [项目结构](/guide/project-structure)
3. 然后根据需要进入 [从零搭建教程](/tutorial/) 或 [参考手册](/reference/)

## 📋 环境要求

| 依赖 | 版本要求 | 用途 |
| --- | --- | --- |
| Go | >= 1.26 | 后端运行 |
| Node.js | >= 20.19 | 前端运行 |
| pnpm | 最新稳定版 | 前端包管理 |
| Docker & Docker Compose | 最新稳定版 | 本地 PostgreSQL 和 Redis |

## 🚀 五步跑起来

### 第 1 步：启动 PostgreSQL 和 Redis

::: code-group

```bash [macOS / Linux]
docker compose -f deploy/compose.local.yml up -d
```

```bash [Windows]
docker compose -f deploy/compose.local.win.yml up -d
```

:::

### 第 2 步：启动后端

```bash
cd server
go run .
```

后端默认监听 `http://localhost:8080`。首次启动会自动执行迁移并创建系统表、菜单和权限种子数据。

### 第 3 步：初始化管理员账号

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/setup/init -ContentType "application/json" -Body '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

```bash [macOS / Linux]
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

:::

### 第 4 步：启动前端

```bash
cd admin
pnpm install
pnpm dev
```

### 第 5 步：启动文档站

```bash
cd docs
pnpm install
pnpm docs:dev
```

## 跑起来之后怎么继续读

- 想先看仓库怎么组织：读 [项目结构](/guide/project-structure)
- 想理解为什么它这样分层：读 [企业级架构升级](/guide/enterprise-architecture)
- 想把 Go 工程结构映射回熟悉的 Java 经验：读 [Go vs Java 工程结构](/guide/java-to-go-structure)
- 想从空仓库一步步搭出来：读 [从零搭建教程](/tutorial/)
- 想查环境变量、建表语句、模块约定：读 [参考手册](/reference/)
