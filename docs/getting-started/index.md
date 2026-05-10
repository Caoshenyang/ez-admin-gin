---
title: 快速开始
description: 五步启动 EZ Admin 本地开发环境
---

# 快速开始

## 环境要求

| 工具 | 最低版本 |
|------|---------|
| Go | 1.26+ |
| Node.js | 20+ |
| pnpm | 9+ |
| Docker | 24+ |
| Docker Compose | v2+ |
| Git | 2.x |

## 1. 克隆仓库

```bash
git clone https://github.com/caoshenyang/ez-admin-gin.git
cd ez-admin-gin
```

## 2. 启动 PostgreSQL 和 Redis

```bash
# macOS / Linux
docker compose -f deploy/compose.local.yml up -d

# Windows
docker compose -f deploy/compose.local.win.yml up -d
```

::: tip
首次启动会拉取 PostgreSQL 和 Redis 镜像，可能需要几分钟。
:::

## 3. 启动后端

```bash
cd server
go run .
```

后端默认监听 `http://localhost:8080`，首次启动自动执行数据库迁移。

## 4. 初始化管理员账号

```bash
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

初始化完成后，用返回的账号登录前端。

## 5. 启动前端

```bash
cd admin
pnpm install
pnpm dev
```

前端默认监听 `http://localhost:5173`，API 请求自动代理到后端。

## 验证

1. 访问 `http://localhost:5173`，应该看到登录页
2. 用初始化的管理员账号登录
3. 登录后应看到 Dashboard 和左侧动态菜单

## 下一步

- [系统架构概览](/architecture/overview) — 了解整体设计
- [权限体系](/architecture/rbac) — 理解 RBAC + Casbin + 数据权限
- [后端模块开发](/backend/module-development) — 学习如何添加新模块
- [部署方案](/deployment/overview) — 了解生产环境部署
