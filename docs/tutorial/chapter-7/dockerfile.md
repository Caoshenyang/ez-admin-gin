---
title: 后端与前端 Dockerfile
description: "为后端 Go 服务和前端 Vue 项目分别编写多阶段构建 Dockerfile，配合 .dockerignore 控制镜像体积和构建速度。"
---

# 后端与前端 Dockerfile

部署的第一步，是让后端和前端各自能被打包成独立、可运行的容器镜像。这一页会为 `server` 和 `admin` 分别编写 Dockerfile，并配上 `.dockerignore` 来排除不需要的文件。

完成后你会得到两个精简镜像：后端是静态编译的 Go 二进制 + Alpine，前端是 Nginx 托管的静态资源。

::: tip 🎯 本节目标
- 为 `server` 和 `admin` 各创建一份 Dockerfile 和 `.dockerignore`
- 理解多阶段构建如何减小镜像体积
- 本地构建并验证镜像可以正常启动
:::

## 后端 Dockerfile

后端采用 Go 多阶段构建：第一阶段编译出静态二进制，第二阶段只拷贝二进制和配置到精简的 Alpine 镜像。

::: details `server/Dockerfile` — Go 多阶段构建
<<< ../../../server/Dockerfile
:::

关键设计点：

| 技术点 | 作用 |
| --- | --- |
| `FROM ... AS builder` | 多阶段构建，编译环境和运行环境隔离 |
| 先 `COPY go.mod go.sum` 再 `COPY .` | 利用 Docker 层缓存，依赖不变时不重复下载 |
| `CGO_ENABLED=0 GOOS=linux` | 静态编译，不依赖 C 库，可以直接跑在 Alpine |
| `-ldflags="-s -w"` | 去掉调试信息，减小二进制体积 |
| 第二阶段 `FROM alpine:3.21` | 运行镜像只有 ~8MB，不含编译工具链 |
| 安装 `ca-certificates tzdata` | 让 HTTPS 请求和时区设置正常工作 |

## 前端 Dockerfile

前端同样是两阶段：第一阶段用 Node 编译出 `dist` 静态资源，第二阶段用 Nginx 来托管。

::: details `admin/Dockerfile` — Node 构建 + Nginx 托管
<<< ../../../admin/Dockerfile
:::

关键设计点：

| 技术点 | 作用 |
| --- | --- |
| `FROM node:22-alpine AS builder` | 使用 Alpine 版 Node，编译阶段尽量轻 |
| 先 `COPY package.json pnpm-lock.yaml` 再 `COPY .` | 同样利用层缓存，依赖不变时跳过 `pnpm install` |
| `corepack enable` | 启用 Corepack，让 `pnpm` 可以直接使用 |
| `--frozen-lockfile` | 严格按 lockfile 安装，保证可复现 |
| 第二阶段 `FROM nginx:1.27-alpine` | 只保留编译产物，不包含 Node 和源码 |
| 删除默认配置 | 后续章节会用自定义 Nginx 配置替代 |

## .dockerignore

`.dockerignore` 的作用和 `.gitignore` 类似：排除不需要进入构建上下文的文件，加快 `docker build` 速度并减小上下文体积。

::: details `server/.dockerignore`
<<< ../../../server/.dockerignore
:::

::: details `admin/.dockerignore`
<<< ../../../admin/.dockerignore
:::

::: warning ⚠️ 别忘了 .dockerignore
如果不排除 `logs/`、`uploads/`、`node_modules/` 这些目录，它们会被完整发送到 Docker 守护进程作为构建上下文，导致构建变慢甚至失败。尤其是 `node_modules`，动辄几百 MB。
:::

## 本地构建与验证

两个 Dockerfile 都写好后，可以在项目根目录分别构建镜像。

构建后端镜像（在 `server/` 目录下执行）：

```bash
cd server
docker build -t ez-admin-server:latest .
```

预期输出（末尾几行）：

```
 => exporting to image
 => => naming to docker.io/library/ez-admin-server:latest
```

构建前端镜像（在 `admin/` 目录下执行）：

```bash
cd admin
docker build -t ez-admin-admin:latest .
```

预期输出（末尾几行）：

```
 => exporting to image
 => => naming to docker.io/library/ez-admin-admin:latest
```

构建完成后，验证镜像是否存在：

```bash
docker images | grep ez-admin
```

预期输出类似：

```
ez-admin-server   latest   ...   ~20MB
ez-admin-admin    latest   ...   ~50MB
```

::: details 为什么后端镜像这么大？
如果后端镜像明显偏大，检查是否忘了 `.dockerignore`，或者 `COPY` 时误带了 `logs/`、`uploads/` 等运行时目录。多阶段构建的正常结果应该在 20MB 左右。
:::

## 小结

- 后端使用 Go 多阶段构建，运行镜像只包含静态二进制和配置文件。
- 前端使用 Node 编译 + Nginx 托管，运行镜像只包含静态资源。
- `.dockerignore` 是必备文件，能显著加快构建速度并避免意外打包运行时文件。
- 两个镜像都可以本地构建并验证。

下一页会把这两个容器组合起来，用 Docker Compose 一键启动完整服务：[Docker Compose 编排](./docker-compose)。
