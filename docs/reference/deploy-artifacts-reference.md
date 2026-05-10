---
title: Docker 与部署产物参考
description: "解析项目中 Dockerfile、Compose、Nginx 和部署脚本的设计，帮助理解当前主线部署与其他部署变体。"
---

# Docker 与部署产物参考

这一页解析项目中 Docker、Compose、Nginx 和部署脚本相关文件的设计。它不是部署步骤页，而是帮助你看懂“当前主线为什么这样拆”和“其他部署变体落在哪些文件里”。

## 先看整体分组

当前仓库里的部署产物可以先按下面四组理解：

| 分组 | 关键文件 | 作用 |
| --- | --- | --- |
| 当前默认主线 | `scripts/pack.sh`、`scripts/setup-server.sh`、`scripts/update-server.sh`、`deploy/compose.server.yml` | 宿主机后端 + 容器基础设施 |
| 全容器部署变体 | `deploy/compose.deploy.yml`、`deploy/compose.prod.yml` | 后端和前端也进入容器 |
| Nginx 入口层 | `deploy/nginx/nginx-native*.conf`、`deploy/nginx/nginx*.conf` | 匹配不同部署形态的入口配置 |
| 本地开发基础环境 | `deploy/compose.local.yml`、`deploy/compose.local.win.yml` | 第 1 章本地开发用的 PostgreSQL / Redis 基础环境 |

如果你当前只关心第 8 章默认部署主线，优先看：

- `scripts/pack.sh`
- `scripts/setup-server.sh`
- `scripts/update-server.sh`
- `deploy/compose.server.yml`
- `deploy/nginx/nginx-native.conf`
- `deploy/nginx/nginx-native-ssl.conf`

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

当前仓库里有五份 Compose 文件，但它们不是同一层面的东西：

| 文件 | 当前角色 |
| --- | --- |
| `deploy/compose.local.yml` | macOS / Linux 本地开发基础环境 |
| `deploy/compose.local.win.yml` | Windows 本地开发基础环境 |
| `deploy/compose.server.yml` | 当前默认服务器主线 |
| `deploy/compose.deploy.yml` | 镜像分发型全容器部署 |
| `deploy/compose.prod.yml` | 源码构建型全容器部署 |

### `compose.server.yml`：当前默认服务器主线

这份文件当前最重要，因为第 8 章默认就是围绕它来讲的。

它表达的是：

- PostgreSQL、Redis、Nginx 进容器
- 后端二进制跑在宿主机
- Nginx 通过 `host` 网络反代 `127.0.0.1:8080`

所以它要和下面这些文件一起看：

- `scripts/setup-server.sh`
- `scripts/update-server.sh`
- `deploy/nginx/nginx-native.conf`
- `deploy/nginx/nginx-native-ssl.conf`

### `compose.deploy.yml`：镜像分发型全容器部署

这份文件适合：

- 服务器只负责拉镜像
- 前后端镜像已经在 Docker Hub 或其他仓库里
- 后端、前端、数据库、Redis、Nginx 都统一放到容器里

它依赖：

- `DOCKERHUB_USERNAME`
- `ez-admin-server:latest`
- `ez-admin-web:latest`

并配套使用：

- `deploy/nginx/nginx.conf`
- `deploy/nginx/nginx-ssl.conf`

### `compose.prod.yml`：源码构建型全容器部署

这份文件和 `compose.deploy.yml` 的整体结构很接近，但后端与前端镜像改成现场 `build`。

它更适合：

- 服务器或 CI 直接持有源码
- 希望在构建机现场产出镜像
- 想验证全容器生产结构

### 四个核心服务

| 服务 | 作用 | 端口 |
| --- | --- | --- |
| PostgreSQL | 业务数据存储 | 仅内部网络（5432） |
| Redis | 缓存与会话存储 | 仅内部网络（6379） |
| Server | 后端 API 服务 | `compose.server.yml` 中由宿主机承担；全容器变体中由容器承担 |
| Nginx | 前端静态资源 + API 反向代理 | 对外暴露（80、443） |

### 启动顺序

```text
postgres（健康检查通过）
    └→ server（等待 postgres + redis 健康）
redis（健康检查通过）  ┘
    └→ nginx（等待 server 启动）
```

在全容器变体里，`depends_on` 使用 `condition: service_healthy` 而不是 `service_started`，确保数据库能响应查询后端才启动，避免连接失败。

### 数据持久化

| 卷名 | 挂载位置 | 保存什么 |
| --- | --- | --- |
| `postgres_data` | `/var/lib/postgresql/data` | 数据库文件 |
| `redis_data` | `/data` | Redis AOF 持久化文件 |
| `uploads_data` | `/app/uploads` | 全容器结构下的用户上传文件 |

命名卷由 Docker 统一管理。`docker compose down` 停止并删除容器后数据仍然保留。彻底清空需要加 `--volumes`：

```bash
docker compose -f compose.deploy.yml down --volumes
```

### 网络隔离

全容器变体中的服务放在 `backend` 桥接网络中：

- 服务之间用服务名互相访问（如 `postgres:5432`、`server:8080`）。
- 只有 Nginx 通过 `ports` 暴露端口到宿主机。
- PostgreSQL 和 Redis 不对外暴露，外部无法直接访问数据库和缓存。

而在 `compose.server.yml` 这条默认主线里，Nginx 改用 `host` 网络，目的是直接访问宿主机上的 `127.0.0.1:8080`。

## 环境变量注入

Compose 中所有 `EZ_` 前缀的环境变量都支持通过 `.env` 文件覆盖：

```yaml
# compose.server.yml / compose.deploy.yml 中都能看到类似写法
EZ_AUTH_JWT_SECRET: ${EZ_AUTH_JWT_SECRET:?JWT_SECRET is required}
```

语法说明：

| 写法 | 含义 |
| --- | --- |
| `${VAR:-default}` | 变量未设置时使用默认值 |
| `${VAR:?error message}` | 变量未设置时报错退出 |

环境变量的完整清单和覆盖机制见 [环境变量参考](./environment-variables-reference)。

## Nginx 配置文件怎么配套看

当前四份 Nginx 配置可以这样记：

| 文件 | 适配哪种结构 |
| --- | --- |
| `deploy/nginx/nginx-native.conf` | 宿主机后端 + 容器基础设施 |
| `deploy/nginx/nginx-native-ssl.conf` | 宿主机后端 + 容器基础设施 + HTTPS |
| `deploy/nginx/nginx.conf` | 全容器部署 |
| `deploy/nginx/nginx-ssl.conf` | 全容器部署 + HTTPS |

它们最核心的差别只有一个：

- Native 版本反代 `127.0.0.1:8080`
- 容器版本反代 `server:8080`

## 部署脚本为什么也要算进部署产物

因为当前默认主线不是“手动敲一串命令”，而是明确依赖三份脚本：

| 文件 | 作用 |
| --- | --- |
| `scripts/pack.sh` / `scripts/pack.ps1` | 本地打包部署包 |
| `scripts/setup-server.sh` | 首次部署初始化 |
| `scripts/update-server.sh` | 日常更新 |

这三份脚本和 `compose.server.yml` 一起，才组成当前完整的服务器交付链路。
