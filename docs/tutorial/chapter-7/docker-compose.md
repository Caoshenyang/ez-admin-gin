---
title: Docker Compose 编排
description: "使用 Docker Compose 编排后端、前端、数据库和 Redis 等服务，完成生产环境一键部署。"
---

# Docker Compose 编排

第 1 章用 `compose.local.yml` 启动了 PostgreSQL 和 Redis，只服务本地开发。这一节用一份新的 `compose.prod.yml` 把后端服务、前端 Nginx 也加进来，实现一条命令启动整个后台。

::: tip 🎯 本节目标
理解生产环境 Compose 的完整服务编排，跑通一键启动，并确认所有服务正常联动。
:::

## 从本地到生产：两份 Compose 的区别

| 对比项 | `compose.local.yml` | `compose.prod.yml` |
| --- | --- | --- |
| 定位 | 本地开发，只启动依赖服务 | 生产部署，启动全部服务 |
| 服务数量 | 2 个（PostgreSQL + Redis） | 4 个（PostgreSQL + Redis + Server + Nginx） |
| 数据卷 | 绑定挂载到本机目录 | Docker 命名卷，由 Docker 管理 |
| 端口映射 | PostgreSQL 和 Redis 映射到本机 | 只有 Nginx 暴露端口 |
| 网络隔离 | 默认网络 | 自定义 `backend` 网络 |
| 环境变量 | 硬编码默认值 | 全部支持 `.env` 覆盖 |

核心区别一句话：**本地版只管基础依赖，生产版把整个后台打包编排。**

::: details 为什么本地版不合并进生产版
本地开发时后端和前端通常用热更新方式运行（`go run` / `pnpm dev`），不需要放进 Compose。两份文件各管各的场景，避免一份配置来回切换环境变量。
:::

## 完整配置

<<< ../../../deploy/compose.prod.yml

## 四个服务各自做什么

### PostgreSQL — 业务数据存储

- 使用 `postgres:18-alpine` 镜像，Alpine 体积更小。
- 通过环境变量 `POSTGRES_USER`、`POSTGRES_PASSWORD`、`POSTGRES_DB` 初始化数据库。
- `PGDATA` 指定数据子目录，避免挂载目录权限问题。
- 健康检查用 `pg_isready`，每 10 秒探测一次，最多重试 5 次。

### Redis — 缓存与会话存储

- 使用 `redis:8-alpine` 镜像。
- `--appendonly yes` 开启 AOF 持久化，容器重启后缓存数据不丢失。
- `--requirepass` 从环境变量读取密码，为空时不设密码。
- 健康检查会在设置密码时自动带上认证参数，避免 Redis 已启用密码但探针仍然匿名访问，导致容器一直停在 `unhealthy`。

### Server — 后端 API 服务

- 从 `server/Dockerfile` 构建镜像，包含编译后的 Go 二进制文件。
- 所有配置通过环境变量注入，不依赖配置文件。
- `depends_on` 声明对 PostgreSQL 和 Redis 的依赖，并要求健康检查通过后才启动。
- 上传文件保存到 `uploads_data` 命名卷，数据不会随容器删除而丢失。

::: warning ⚠️ JWT_SECRET 是必填项
`EZ_AUTH_JWT_SECRET` 使用 `${EZ_AUTH_JWT_SECRET:?JWT_SECRET is required}` 语法。如果在 `.env` 文件或环境变量中没有设置这个值，`docker compose up` 会直接报错，不会用空值启动。这是故意的——生产环境不允许 JWT 密钥为空。
:::

### Nginx — 前端静态资源 + 反向代理

- 从 `admin/Dockerfile` 构建镜像，包含打包后的前端静态文件。
- 把 `deploy/nginx/nginx.conf` 以只读方式挂载到容器内，覆盖默认配置。
- 只有这一个服务暴露端口（默认 80），所有外部流量都从 Nginx 进入。
- Nginx 配置的详细说明在[下一节](./nginx-config)。

## 服务依赖与启动顺序

Compose 会按 `depends_on` 的顺序启动服务：

```
postgres (健康检查通过)
    └→ server (等待 postgres + redis 健康)
redis (健康检查通过)   ┘
    └→ nginx (等待 server 启动)
```

启动条件用的是 `condition: service_healthy`，而不是简单的 `service_started`。这意味着：

- PostgreSQL 必须能响应 `pg_isready`，后端才会启动。
- Redis 必须能响应 `ping`，后端才会启动。
- 后端容器启动后，Nginx 才会启动。

这样避免了后端启动时数据库还没准备好的问题。

## 数据持久化

生产环境使用 Docker 命名卷，而不是绑定挂载到本机目录：

| 卷名 | 挂载位置 | 保存什么 |
| --- | --- | --- |
| `postgres_data` | `/var/lib/postgresql/data` | 数据库文件 |
| `redis_data` | `/data` | Redis AOF 持久化文件 |
| `uploads_data` | `/app/uploads` | 用户上传的文件 |

命名卷由 Docker 统一管理，`docker compose down` 停止并删除容器后数据仍然保留。如果需要彻底清空数据，要显式加上 `--volumes` 参数：

```bash
# 停止并删除容器，同时删除命名卷（数据会丢失）
docker compose -f deploy/compose.prod.yml down --volumes
```

::: warning ⚠️ down --volumes 会清空所有数据
执行 `down --volumes` 后，数据库、缓存和上传文件都会被删除。生产环境谨慎操作。
:::

## 网络隔离

所有服务都放在 `backend` 桥接网络中：

- 服务之间可以用服务名互相访问（如 `postgres:5432`、`redis:6379`、`server:8080`）。
- 只有 Nginx 通过 `ports` 暴露了 80 端口到宿主机。
- PostgreSQL 和 Redis 不对外暴露端口，外部无法直接访问数据库和缓存。

这种设计减少了攻击面，数据库和缓存只对后端服务可见。

## 🚀 启动全部服务

在项目根目录执行：

```bash
# 构建镜像并在后台启动全部服务
docker compose -f deploy/compose.prod.yml up -d --build
```

`--build` 会让 Docker 重新构建 `server` 和 `nginx` 的镜像。如果镜像已经构建过且代码没变，可以去掉 `--build` 加快启动。

查看服务状态：

```bash
docker compose -f deploy/compose.prod.yml ps
```

正常情况下应该看到四个服务都是 `running` / `healthy` 状态。

查看后端日志：

```bash
docker compose -f deploy/compose.prod.yml logs -f server
```

::: warning ⚠️ 首次启动前需要准备 .env 文件
`compose.prod.yml` 中至少需要设置 `EZ_AUTH_JWT_SECRET`。在 `deploy/` 目录下创建 `.env` 文件：

```bash
# 在 deploy/ 目录下创建 .env
# 下面是示例值，实际部署请替换为自己的随机密钥
EZ_AUTH_JWT_SECRET=your-random-secret-key-here
```

环境变量的完整说明在[环境变量与初始化数据](./env-and-init-data)一节。
:::

## 常用运维命令

```bash
# 查看所有服务状态
docker compose -f deploy/compose.prod.yml ps

# 查看全部服务日志
docker compose -f deploy/compose.prod.yml logs -f

# 只看某个服务的日志
docker compose -f deploy/compose.prod.yml logs -f server
docker compose -f deploy/compose.prod.yml logs -f nginx

# 重启某个服务
docker compose -f deploy/compose.prod.yml restart server

# 停止全部服务（保留数据）
docker compose -f deploy/compose.prod.yml stop

# 重新启动已停止的服务
docker compose -f deploy/compose.prod.yml start
```

## 小结

这一节把四个服务编排进了同一个 Compose 文件：

- **PostgreSQL + Redis** 提供数据存储和缓存，通过健康检查确保可用后才让后端启动。
- **Server** 作为后端 API 服务，依赖数据库和缓存就绪。
- **Nginx** 作为唯一入口，统一处理前端静态资源和 API 反向代理。
- 数据使用命名卷持久化，网络使用桥接隔离。

接下来看 Nginx 配置的具体细节：[Nginx 配置](./nginx-config)。
