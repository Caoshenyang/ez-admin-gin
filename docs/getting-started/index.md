---
title: 快速开始
description: 五步启动 EZ Admin 本地开发环境
---

# 快速开始

这一页帮助你在本地跑起来 EZ Admin 的完整开发环境，完成后你能看到一个可登录、可操作的管理后台。

::: tip 🎯 本节目标
启动 PostgreSQL + Redis → 启动后端 → 初始化管理员 → 启动前端 → 浏览器验证。
:::

## 环境要求

| 工具 | 最低版本 |
|------|---------|
| Go | 1.26+ |
| Node.js | 20+ |
| pnpm | 9+ |
| Docker | 24+ |
| Docker Compose | v2+ |
| Git | 2.x |
| make | GNU Make 4+ |

::: details Windows 用户：如何安装 make
Windows 默认不带 make，需要手动安装：

::: code-group

```bash [Chocolatey]
choco install make
```

```bash [Scoop]
scoop install make
```

:::

安装后重启终端，运行 `make --version` 验证。如果暂时不想安装 make，也可以直接查看 Makefile 中对应的原始命令手动执行。
:::

## 1. 克隆仓库

```bash
git clone https://github.com/caoshenyang/ez-admin-gin.git
cd ez-admin-gin
```

## 2. 启动 PostgreSQL 和 Redis

::: code-group

```bash [使用 make]
make docker-up
```

```bash [手动执行]
# macOS / Linux
docker compose -f deploy/compose.local.yml up -d

# Windows
docker compose -f deploy/compose.local.win.yml up -d
```

:::

::: tip
首次启动会拉取 PostgreSQL 和 Redis 镜像，可能需要几分钟。
:::

执行后应看到两个容器 `ez-admin-postgres` 和 `ez-admin-redis` 处于运行状态：

```bash
docker ps --format "table {{.Names}}\t{{.Status}}"
# 期望输出：
# NAMES                STATUS
# ez-admin-postgres    Up ... (healthy)
# ez-admin-redis       Up ... (healthy)
```

## 3. 启动后端

打开一个新终端：

::: code-group

```bash [使用 make]
make server-dev
```

```bash [手动执行]
cd server && go run .
```

:::

后端默认监听 `http://localhost:8080`，首次启动自动执行数据库迁移。

启动成功后会看到类似输出：

```
[INFO] server started on :8080
```

## 4. 初始化管理员账号

```bash
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

执行后应返回包含 `token` 的 JSON。初始化完成后，用返回的账号登录前端。

::: warning ⚠️ 初始化只能执行一次
如果返回 `already initialized` 类型的错误，说明已经初始化过了，直接跳到下一步用已有账号登录即可。
:::

## 5. 启动前端

再打开一个新终端：

::: code-group

```bash [使用 make]
make install && make admin-dev
```

```bash [手动执行]
cd admin && pnpm install && pnpm dev
```

:::

前端默认监听 `http://localhost:5173`，API 请求自动代理到后端。

## 验证

1. 访问 `http://localhost:5173`，应该看到登录页
2. 用初始化的管理员账号登录
3. 登录后应看到 Dashboard 和左侧动态菜单

## 常用命令速查

项目根目录的 Makefile 提供了统一的开发入口，查看所有命令：

```bash
make help
```

| 命令 | 说明 |
|------|------|
| `make server-dev` | 启动后端 |
| `make admin-dev` | 启动前端 |
| `make docs-dev` | 启动文档站 |
| `make test` | 运行后端测试 |
| `make lint` | 运行所有 lint（后端 + 前端） |
| `make admin-check` | 前端类型检查 + lint |
| `make build` | 构建后端二进制 + 前端产物 |
| `make docker-up` | 启动 PostgreSQL + Redis |
| `make docker-down` | 停止 PostgreSQL + Redis |
| `make docker-config` | 验证 Docker Compose 配置 |
| `make docker-build` | 构建 Docker 镜像 |
| `make clean` | 清理构建产物 |

## 下一步

- [系统架构概览](/architecture/overview) — 了解整体设计
- [权限体系](/architecture/rbac) — 理解 RBAC + Casbin + 数据权限
- [后端模块开发](/backend/module-development) — 学习如何添加新模块
- [部署方案](/deployment/overview) — 了解生产环境部署
