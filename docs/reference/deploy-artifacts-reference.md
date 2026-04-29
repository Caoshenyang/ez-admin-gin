---
title: Docker 部署文件参考
description: "解析项目中 Dockerfile、Compose 和部署配置的设计：多阶段构建、服务编排、环境变量注入和数据持久化。"
---

# Docker 部署文件参考

这一页解析项目中 Docker 相关文件的设计：Dockerfile 如何构建镜像、Compose 如何编排服务、环境变量如何注入。部署时不需要修改这些文件，遇到问题需要排查或想自定义时，回来查这一页。

## Dockerfile

项目有两份 Dockerfile，都使用多阶段构建——第一阶段编译，第二阶段只保留运行产物。这样运行镜像不包含编译工具链，体积更小、攻击面更小。

### 后端 Dockerfile

::: details `server/Dockerfile`
```dockerfile
# ---- 构建阶段：编译 Go 二进制 ----
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# 先复制依赖文件，利用 Docker 层缓存加速后续构建。
COPY go.mod go.sum ./
RUN go mod download

# 再复制源码并编译。
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# ---- 运行阶段：只保留二进制和配置 ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/server .
COPY configs/ ./configs/

# 上传目录和日志目录需要在运行时存在。
RUN mkdir -p uploads logs

EXPOSE 8080

CMD ["./server"]
```
:::

| 技术点 | 作用 |
| --- | --- |
| `FROM ... AS builder` | 多阶段构建，编译环境和运行环境隔离 |
| 先 `COPY go.mod go.sum` 再 `COPY .` | 利用 Docker 层缓存：依赖不变时跳过 `go mod download`，只重新编译代码 |
| `CGO_ENABLED=0 GOOS=linux` | 静态编译，不依赖 C 库，可以直接跑在 Alpine |
| `-ldflags="-s -w"` | 去掉调试信息和符号表，减小二进制体积 |
| 第二阶段 `FROM alpine:3.21` | 运行镜像约 8MB，不含 Go 编译器 |
| 安装 `ca-certificates tzdata` | 让 HTTPS 请求和时区设置正常工作 |

### 前端 Dockerfile

::: details `admin/Dockerfile`
```dockerfile
# ---- 构建阶段：编译前端资源 ----
FROM node:22-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用 Docker 层缓存。
COPY package.json pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

# 再复制源码并构建。
COPY . .
RUN pnpm build

# ---- 运行阶段：Nginx 托管静态资源 ----
FROM nginx:1.27-alpine

# 删除默认配置，使用项目自定义配置。
RUN rm /etc/nginx/conf.d/default.conf

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```
:::

| 技术点 | 作用 |
| --- | --- |
| 先 `COPY package.json pnpm-lock.yaml` 再 `COPY .` | 利用层缓存：依赖不变时跳过 `pnpm install` |
| `corepack enable` | 启用 Corepack，让 `pnpm` 可以直接使用 |
| `--frozen-lockfile` | 严格按 lockfile 安装，保证可复现 |
| 第二阶段 `FROM nginx:1.27-alpine` | 只保留编译产物，不包含 Node 和源码 |

### .dockerignore

`.dockerignore` 排除不需要进入构建上下文的文件，加快 `docker build` 速度并减小上下文体积。

> [!WARNING]
> 如果不排除 `logs/`、`uploads/`、`node_modules/` 这些目录，它们会被完整发送到 Docker 守护进程作为构建上下文。尤其是 `node_modules`，动辄几百 MB，会导致构建变慢甚至失败。

## Compose 文件

项目有两份 Compose 文件，各管一个场景：

| 文件 | 场景 | 镜像来源 |
| --- | --- | --- |
| `deploy/compose.prod.yml` | 本地一键构建并启动全部服务 | `build:` 从源码构建 |
| `deploy/compose.deploy.yml` | 云服务器部署，从 Docker Hub 拉取 | `image:` 从 Docker Hub 拉取 |

两者的服务定义、网络和卷配置基本一致，区别只在后端和 Nginx 的镜像来源。

### 四个服务

| 服务 | 作用 | 端口 |
| --- | --- | --- |
| PostgreSQL | 业务数据存储 | 仅内部网络（5432） |
| Redis | 缓存与会话存储 | 仅内部网络（6379） |
| Server | 后端 API 服务 | 仅内部网络（8080） |
| Nginx | 前端静态资源 + API 反向代理 | 对外暴露（80、443） |

### 启动顺序

```text
postgres（健康检查通过）
    └→ server（等待 postgres + redis 健康）
redis（健康检查通过）  ┘
    └→ nginx（等待 server 启动）
```

`depends_on` 使用 `condition: service_healthy` 而不是 `service_started`，确保数据库能响应查询后端才启动，避免连接失败。

### 数据持久化

| 卷名 | 挂载位置 | 保存什么 |
| --- | --- | --- |
| `postgres_data` | `/var/lib/postgresql/data` | 数据库文件 |
| `redis_data` | `/data` | Redis AOF 持久化文件 |
| `uploads_data` | `/app/uploads` | 用户上传的文件 |

命名卷由 Docker 统一管理。`docker compose down` 停止并删除容器后数据仍然保留。彻底清空需要加 `--volumes`：

```bash
docker compose -f compose.deploy.yml down --volumes
```

### 网络隔离

所有服务放在 `backend` 桥接网络中：

- 服务之间用服务名互相访问（如 `postgres:5432`、`server:8080`）。
- 只有 Nginx 通过 `ports` 暴露端口到宿主机。
- PostgreSQL 和 Redis 不对外暴露，外部无法直接访问数据库和缓存。

## 环境变量注入

Compose 中所有 `EZ_` 前缀的环境变量都支持通过 `.env` 文件覆盖：

```yaml
# compose.deploy.yml 中的写法
EZ_AUTH_JWT_SECRET: ${EZ_AUTH_JWT_SECRET:?JWT_SECRET is required}
```

语法说明：

| 写法 | 含义 |
| --- | --- |
| `${VAR:-default}` | 变量未设置时使用默认值 |
| `${VAR:?error message}` | 变量未设置时报错退出 |

环境变量的完整清单和覆盖机制见 [环境变量与初始化数据](../tutorial/chapter-7/env-and-init-data)。
