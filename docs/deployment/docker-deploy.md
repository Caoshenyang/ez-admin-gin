---
title: Docker 部署
description: 全容器化部署方案，使用 Docker Hub 镜像或本地构建
---

# Docker 部署

本文介绍 EZ Admin 的全容器化部署方式。所有服务（后端、前端、数据库、Redis、Nginx）都运行在 Docker 容器中，适合需要镜像版本管理或 CI/CD 集成的场景。

::: tip 先阅读部署概览
如果你还不确定选择哪种部署方式，请先阅读 [部署概览](/deployment/overview)。
:::

## 适用场景

- 通过 Docker Hub 分发镜像
- 需要精确的版本管理和回滚能力
- 团队有成熟的 Docker 运维体系
- 不希望宿主机直接运行后端二进制

## 与服务器二进制部署的区别

| | 服务器二进制部署 | Docker 全容器化 |
|---|---|---|
| 后端 | 二进制 + systemd | Docker 容器 |
| 前端 | Nginx 容器挂载静态文件 | Nginx 容器（镜像内置） |
| 网络模式 | host 网络 | bridge 网络 |
| 镜像来源 | 本地编译 | Docker Hub 或本地构建 |
| 配置文件 | `compose.server.yml` | `compose.deploy.yml` |
| 后端管理 | `systemctl restart ez-admin` | `docker compose restart server` |

## Compose 文件说明

| 文件 | 用途 | 网络模式 | 后端运行方式 |
|------|------|---------|-------------|
| `compose.local.yml` | 本地开发（macOS / Linux） | bridge | 不包含后端，需本地 `go run` |
| `compose.local.win.yml` | 本地开发（Windows） | bridge | 不包含后端，需本地 `go run` |
| `compose.server.yml` | 服务器二进制部署 | host | 宿主机二进制 |
| `compose.deploy.yml` | 全容器化部署 | bridge | Docker 容器 |
| `compose.prod.yml` | 生产环境演示配置 | bridge | Docker 容器 |

::: details 各文件的定位
- **本地开发**使用 `compose.local.yml`（macOS/Linux）或 `compose.local.win.yml`（Windows），只启动 PostgreSQL 和 Redis。
- **服务器二进制部署**使用 `compose.server.yml`，启动 PostgreSQL、Redis 和 Nginx，后端由 systemd 管理。
- **Docker 全容器化部署**使用 `compose.deploy.yml`，所有服务都运行在容器中。
- **`compose.prod.yml`** 是生产环境演示配置，结构与 `compose.deploy.yml` 类似，使用本地构建而非 Docker Hub 镜像。
:::

## 前置条件

- 已安装 Docker 和 Docker Compose
- 如果使用 Docker Hub 镜像：需要 Docker Hub 账号，且已推送镜像
- 如果使用本地构建：项目源码中需要有 `server/Dockerfile` 和 `admin/Dockerfile`
- 开放 80 / 443 端口

## 使用 Docker Hub 镜像

### 1. 配置环境变量

```bash
mkdir -p /opt/ez-admin && cd /opt/ez-admin

# 创建环境变量文件
cp .env.example .env
```

编辑 `.env`，必须配置以下变量：

```bash
# Docker Hub 用户名（必须，镜像名格式为 <USERNAME>/ez-admin-server:latest）
DOCKERHUB_USERNAME=your-dockerhub-username

# JWT 密钥（必须，生产环境不能使用默认值）
EZ_AUTH_JWT_SECRET=$(openssl rand -hex 32)

# 数据库密码（建议修改）
EZ_DATABASE_PASSWORD=your-strong-password
```

::: warning
`DOCKERHUB_USERNAME` 和 `EZ_AUTH_JWT_SECRET` 是必须设置的变量。如果未设置，`docker compose up` 会报错退出。
:::

### 2. 准备 Nginx 配置

```bash
mkdir -p nginx/ssl
```

将 Nginx 配置文件放入 `nginx/` 目录：

- HTTP 模式：`nginx/nginx.conf`
- HTTPS 模式：`nginx/nginx-ssl.conf`

通过 `EZ_NGINX_CONF` 环境变量选择配置文件（默认 `nginx.conf`）。

### 3. 启动服务

```bash
docker compose -f compose.deploy.yml --env-file .env up -d
```

### 4. 初始化管理员

```bash
curl -X POST http://localhost/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin@123456","nickname":"管理员"}'
```

::: danger
首次登录后请立即修改默认密码。
:::

## 使用本地构建

如果你不想使用 Docker Hub，可以使用 `compose.prod.yml` 从源码构建：

```bash
docker compose -f compose.prod.yml up -d
```

`compose.prod.yml` 会从 `server/Dockerfile` 和 `admin/Dockerfile` 构建镜像，不需要 Docker Hub 账号。

::: details 前提条件
使用本地构建需要：
- 项目源码在服务器上
- 服务器有 Go 和 Node.js 构建环境（Docker 多阶段构建会处理依赖）
- `server/Dockerfile` 和 `admin/Dockerfile` 存在
:::

## 网络架构

`compose.deploy.yml` 使用 bridge 网络模式，所有服务通过容器名通信：

```
客户端
  │
  ▼ :80 / :443
┌──────────────┐
│  nginx 容器   │ ← 端口映射到宿主机
└──────┬───────┘
       │ bridge 网络
  ┌────┴────┐
  ▼         ▼
┌───────┐ ┌─────────┐
│ server │ │ postgres │ ← 容器间通过服务名访问
│ :8080  │ │ :5432    │
└───┬───┘ └─────────┘
    │
    ▼
┌───────┐
│ redis  │
│ :6379  │
└───────┘
```

与服务器二进制部署的关键差异：

- 后端通过容器名 `postgres` 连接数据库，而非 `127.0.0.1`
- Nginx 通过容器名 `server` 代理后端
- 只有 Nginx 暴露端口到宿主机（80 / 443）

## 环境变量

`compose.deploy.yml` 中大部分环境变量通过 `.env` 文件注入：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DOCKERHUB_USERNAME` | Docker Hub 用户名，**必须设置** | 无 |
| `EZ_AUTH_JWT_SECRET` | JWT 密钥，生产必须修改 | 无 |
| `EZ_APP_ENV` | 运行环境 | `prod` |
| `EZ_DATABASE_PASSWORD` | 数据库密码 | `ez_admin_123456` |
| `EZ_NGINX_PORT` | Nginx 监听端口 | `80` |
| `EZ_NGINX_CONF` | Nginx 配置文件名 | `nginx.conf` |

完整变量说明见 [环境变量参考](/reference/environment-variables-reference)。

## 镜像构建

如果你需要自行构建并推送镜像：

```bash
# 构建后端镜像
docker build -t your-username/ez-admin-server:latest ./server

# 构建前端镜像
docker build -t your-username/ez-admin-web:latest ./admin

# 推送到 Docker Hub
docker push your-username/ez-admin-server:latest
docker push your-username/ez-admin-web:latest
```

::: warning
项目目前没有内置自动 CI/CD 镜像发布流程。如果你需要自动化构建和推送镜像，需要自行配置 GitHub Actions 或其他 CI 工具。
:::

## 数据持久化

`compose.deploy.yml` 使用 Docker named volumes 管理数据：

| Volume | 用途 |
|--------|------|
| `postgres_data` | PostgreSQL 数据 |
| `redis_data` | Redis 持久化 |
| `uploads_data` | 上传文件 |

查看卷位置：

```bash
docker volume ls | grep ez-admin
```

::: danger
删除容器不会丢失数据，但 `docker compose down -v` 会删除所有 volumes。备份数据前不要使用 `-v` 选项。
:::

## 常用操作

```bash
# 查看所有容器状态
docker compose -f compose.deploy.yml ps

# 查看后端日志
docker compose -f compose.deploy.yml logs -f server

# 重启后端
docker compose -f compose.deploy.yml restart server

# 停止所有服务
docker compose -f compose.deploy.yml down

# 重新构建并启动
docker compose -f compose.deploy.yml up -d --build
```

## 常见问题

### `DOCKERHUB_USERNAME is required` 报错

`compose.deploy.yml` 要求设置 `DOCKERHUB_USERNAME` 环境变量。如果不想使用 Docker Hub，改用 `compose.prod.yml` 从源码构建。

### 容器启动后 API 502

```bash
# 检查所有容器状态
docker compose -f compose.deploy.yml ps

# 如果 server 容器未运行，查看日志
docker compose -f compose.deploy.yml logs server
```

### Nginx 配置文件找不到

确保 `nginx/` 目录下存在 `EZ_NGINX_CONF` 指定的配置文件。默认查找 `nginx/nginx.conf`。

### 数据库连接失败

在 bridge 网络模式下，后端通过容器名 `postgres` 连接数据库。确认 `compose.deploy.yml` 中后端的 `EZ_DATABASE_HOST` 环境变量是 `postgres`，而不是 `127.0.0.1`。
