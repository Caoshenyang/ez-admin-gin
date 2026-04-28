---
title: Docker Compose 基础环境
description: "准备 PostgreSQL、Redis 等本地依赖，为后续后端开发提供基础环境。"
---

# Docker Compose 基础环境

后端马上会接入数据库和缓存。先用 Docker Compose 准备一套本地基础环境，让后续章节不用反复手动安装 PostgreSQL 和 Redis。

::: tip 🎯 本节目标
创建一份本地 Compose 配置，启动 PostgreSQL 和 Redis，并确认两个服务都能正常响应。
:::

## 环境清单

这一节只需要准备 Docker 相关工具：

| 工具 | 要求 | 用途 |
| --- | --- | --- |
| Docker Desktop | 最新稳定版 | 在本机运行基础服务容器 |
| Docker Compose | Docker Desktop 自带即可 | 用一个配置文件管理多个服务 |

## 🛠️ 环境要求

Windows 和 macOS 推荐安装 Docker Desktop：

- Docker Desktop：[https://docs.docker.com/desktop/](https://docs.docker.com/desktop/)
- Docker Compose 安装说明：[https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)

安装完成后，确认 Docker 和 Compose 都可用：

```bash
# 确认 Docker Engine 和 Compose 插件可用
docker version
docker compose version
```

::: warning ⚠️ 先启动 Docker Desktop
如果命令提示无法连接 Docker daemon，通常是 Docker Desktop 还没有启动，或者 Windows 下 WSL / 虚拟化环境没有准备好。
:::

## 镜像下载慢怎么办

第一次启动时，Docker 需要下载 `postgres:18-alpine` 和 `redis:8-alpine`。如果下载很慢，可以先单独拉取镜像：

```bash
# 单独拉取镜像，便于观察下载进度和失败原因
docker pull postgres:18-alpine
docker pull redis:8-alpine
```

如果经常拉取 Docker Hub 镜像较慢，可以配置 registry mirror。Docker Desktop 可以在 `Settings` -> `Docker Engine` 中加入：

```json
{
  "registry-mirrors": ["https://<你的镜像加速地址>"]
}
```

配置后点击 `Apply & Restart`，再重新执行 `docker pull` 或 `docker compose up -d`。

::: warning ⚠️ 不要随意复制陌生镜像源
镜像源会影响你拉取到的镜像来源。优先使用自己云厂商账号提供的镜像加速地址，或团队内部维护的可信镜像源。
:::

Docker 官方说明：[Docker Hub mirror](https://docs.docker.com/docker-hub/image-library/mirror/)

## 为什么先用 PostgreSQL

本教程先使用 PostgreSQL 和 Redis：

| 服务 | 用途 |
| --- | --- |
| PostgreSQL | 保存用户、角色、菜单、配置、日志等业务数据 |
| Redis | 保存缓存、临时状态、登录相关数据 |

::: details PostgreSQL 和 MySQL 怎么选
两者都可以做后台底座数据库。本教程先选 PostgreSQL，是为了让后续权限、日志、复杂查询和扩展能力有更稳定的默认选择。

如果你更熟悉 MySQL，后续用 GORM 接入时也可以替换；但为了教程主线清晰，前面章节只保留 PostgreSQL 一条路径。
:::

## 创建 Compose 文件

先创建本地数据目录：

::: code-group

```powershell [Windows PowerShell]
# 创建 PostgreSQL 和 Redis 的本地数据目录
New-Item -ItemType Directory -Path D:\ez-admin-gin-data\postgres, D:\ez-admin-gin-data\redis -Force
```

```bash [macOS / Linux]
# 创建 PostgreSQL 和 Redis 的本地数据目录
mkdir -p ~/ez-admin-gin-data/postgres ~/ez-admin-gin-data/redis
```

:::

项目提供了两份本地 Compose 文件，按你的操作系统选择：

| 文件 | 适用平台 | 数据目录 |
| --- | --- | --- |
| `compose.local.yml` | macOS / Linux | `${HOME}/ez-admin-gin-data` |
| `compose.local.win.yml` | Windows | `D:/ez-admin-gin-data` |

::: details `deploy/compose.local.yml` — Docker Compose 本地开发配置

::: code-group

```yaml [macOS / Linux (compose.local.yml)]
name: ez-admin-gin

services:
  postgres:
    # 使用 PostgreSQL 官方 Alpine 镜像，体积更小。
    image: postgres:18-alpine
    container_name: ez-admin-postgres
    restart: unless-stopped
    environment:
      # 这三项会创建本地开发使用的默认用户和数据库。
      POSTGRES_USER: ez_admin
      POSTGRES_PASSWORD: ez_admin_123456
      POSTGRES_DB: ez_admin
      # 把数据库真实数据放到 pgdata 子目录，避免挂载目录权限问题。
      PGDATA: /var/lib/postgresql/data/pgdata
      TZ: Asia/Shanghai
    ports:
      # 左侧是本机端口，右侧是容器端口。
      - "5432:5432"
    volumes:
      # 绑定挂载到本机目录，删除容器后数据仍然保留。
      - ${HOME}/ez-admin-gin-data/postgres:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ez_admin -d ez_admin"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    # 使用 Redis 官方 Alpine 镜像。
    image: redis:8-alpine
    container_name: ez-admin-redis
    restart: unless-stopped
    # 开启 AOF，让本地 Redis 数据可以持久化。
    command: ["redis-server", "--appendonly", "yes"]
    ports:
      - "6379:6379"
    volumes:
      # Redis 数据保存到本机目录。
      - ${HOME}/ez-admin-gin-data/redis:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
```

```yaml [Windows (compose.local.win.yml)]
name: ez-admin-gin

services:
  postgres:
    # 使用 PostgreSQL 官方 Alpine 镜像，体积更小。
    image: postgres:18-alpine
    container_name: ez-admin-postgres
    restart: unless-stopped
    environment:
      # 这三项会创建本地开发使用的默认用户和数据库。
      POSTGRES_USER: ez_admin
      POSTGRES_PASSWORD: ez_admin_123456
      POSTGRES_DB: ez_admin
      # 把数据库真实数据放到 pgdata 子目录，避免挂载目录权限问题。
      PGDATA: /var/lib/postgresql/data/pgdata
      TZ: Asia/Shanghai
    ports:
      # 左侧是本机端口，右侧是容器端口。
      - "5432:5432"
    volumes:
      # 绑定挂载到本机目录，删除容器后数据仍然保留。
      - D:/ez-admin-gin-data/postgres:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ez_admin -d ez_admin"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    # 使用 Redis 官方 Alpine 镜像。
    image: redis:8-alpine
    container_name: ez-admin-redis
    restart: unless-stopped
    # 开启 AOF，让本地 Redis 数据可以持久化。
    command: ["redis-server", "--appendonly", "yes"]
    ports:
      - "6379:6379"
    volumes:
      # Redis 数据保存到本机目录。
      - D:/ez-admin-gin-data/redis:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
```

:::

::: info 镜像版本
截止 2026-04-21，Docker Hub 官方镜像已提供 PostgreSQL 18 和 Redis 8。本教程使用 `postgres:18-alpine`、`redis:8-alpine`，让补丁版本跟随官方镜像更新。
:::

官方镜像资料：

| 镜像 | 资料 |
| --- | --- |
| `postgres` | [Docker Hub](https://hub.docker.com/_/postgres) |
| `redis` | [Docker Hub](https://hub.docker.com/_/redis) |

## 🛠️ 启动基础服务

::: code-group

```bash [macOS / Linux]
# 使用 macOS/Linux 版 Compose 文件启动服务
docker compose -f deploy/compose.local.yml up -d
```

```powershell [Windows PowerShell]
# 使用 Windows 版 Compose 文件启动服务
docker compose -f deploy/compose.local.win.yml up -d
```

:::

查看运行状态：

::: code-group

```bash [macOS / Linux]
docker compose -f deploy/compose.local.yml ps
```

```powershell [Windows PowerShell]
docker compose -f deploy/compose.local.win.yml ps
```

:::

看到 `postgres` 和 `redis` 都处于 running / healthy 状态，就说明基础服务已经启动。

在 Docker Desktop 里，也会看到一个名为 `ez-admin-gin` 的 Stack，里面包含：

- `ez-admin-postgres`
- `ez-admin-redis`

::: warning ⚠️ 端口被占用
如果本机已经安装过 PostgreSQL 或 Redis，`5432`、`6379` 可能被占用。可以先停止本机已有服务，或者把 Compose 文件里的左侧端口改成其他端口。
:::

## ✅ 验证 PostgreSQL

::: code-group

```bash [macOS / Linux]
# 在 postgres 容器中执行一条简单 SQL
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select 1;"
```

```powershell [Windows PowerShell]
# 在 postgres 容器中执行一条简单 SQL
docker compose -f deploy/compose.local.win.yml exec postgres psql -U ez_admin -d ez_admin -c "select 1;"
```

:::

能看到查询结果 `1`，说明数据库可以连接。

## ✅ 验证 Redis

::: code-group

```bash [macOS / Linux]
# 在 redis 容器中执行 ping
docker compose -f deploy/compose.local.yml exec redis redis-cli ping
```

```powershell [Windows PowerShell]
# 在 redis 容器中执行 ping
docker compose -f deploy/compose.local.win.yml exec redis redis-cli ping
```

:::

能看到：

```text
PONG
```

说明 Redis 可以连接。

## 常用管理命令

后续开发时，经常会用到这些命令（以 macOS/Linux 为例，Windows 用户将文件名替换为 `compose.local.win.yml`）：

```bash
# 查看服务状态
docker compose -f deploy/compose.local.yml ps

# 查看所有服务日志
docker compose -f deploy/compose.local.yml logs -f

# 查看单个服务日志
docker compose -f deploy/compose.local.yml logs -f postgres
docker compose -f deploy/compose.local.yml logs -f redis

# 停止服务，但保留容器和数据
docker compose -f deploy/compose.local.yml stop

# 重新启动服务
docker compose -f deploy/compose.local.yml start

# 重启服务
docker compose -f deploy/compose.local.yml restart
```

## 数据保存在哪里

数据目录取决于你使用的 Compose 文件：

| 平台 | Compose 文件 | 数据目录 |
| --- | --- | --- |
| macOS / Linux | `compose.local.yml` | `~/ez-admin-gin-data/postgres/pgdata`、`~/ez-admin-gin-data/redis` |
| Windows | `compose.local.win.yml` | `D:\ez-admin-gin-data\postgres\pgdata`、`D:\ez-admin-gin-data\redis` |

停止并删除容器不会删除这些本地数据：

```bash
# 停止并删除容器，但保留绑定挂载到本机的数据目录
docker compose -f deploy/compose.local.yml down
```

如果确实想清空本地数据，先执行 `down`，再删除对应的数据目录。

::: warning ⚠️ 删除数据目录会清空本地数据
删除 `ez-admin-gin-data` 下的目录后，PostgreSQL 和 Redis 的本地数据会丢失。执行前确认这些数据不再需要。
:::

::: tip 当前只准备环境
这一节只启动本地依赖服务。后端如何读取配置、连接 PostgreSQL 和 Redis，会在后面的基础设施章节里完成。
:::

下一章开始进入后端基础设施：[配置管理](../chapter-2/config-management)。
