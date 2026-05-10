---
title: 服务器二进制部署
description: 一台 VPS 从零到上线，一键打包 + 自动初始化，最推荐的部署方式
---

# 服务器二进制部署

使用宿主机运行 Go 后端二进制文件，使用 Docker Compose 管理 PostgreSQL、Redis 和 Nginx。这是面向一台 VPS / 云服务器的推荐部署方式，提供一键打包脚本和自动初始化，流程最短，最容易排查问题。

## 适用场景

- 个人项目上线
- 小型后台系统
- 想快速复用这套后台底座
- 希望后端由 systemd 管理，崩溃自动重启
- 希望数据库、Redis、Nginx 由 Docker Compose 管理

## 部署架构

```
Browser
  │
  ▼ :80 / :443
┌──────────────────────────────────────────┐
│  Nginx 容器 (network_mode: host)         │
│  ├── /            → 前端静态文件 (web/)  │
│  ├── /assets/*    → 静态资源（缓存 1 年）│
│  ├── /api/*       → 127.0.0.1:8080       │
│  ├── /uploads/*   → 127.0.0.1:8080       │
│  └── /health      → 127.0.0.1:8080       │
└──────────────────────────────────────────┘
  │
  ▼
┌─────────────────────┐
│  Go 后端 (systemd)  │  ← 二进制直接运行在宿主机 :8080
│  WorkingDirectory    │
│  /opt/ez-admin       │
└────────┬────────────┘
         │            │
         ▼            ▼
┌──────────────┐  ┌──────────┐
│  PostgreSQL   │  │  Redis   │  ← Docker 容器，仅监听 127.0.0.1
│  :5432        │  │  :6379   │
└──────────────┘  └──────────┘
```

关键设计：

- 前端静态资源由 Nginx 容器提供（挂载 `/opt/ez-admin/web`）
- `/api` 请求通过 Nginx 反向代理到宿主机后端
- Go 后端运行在宿主机，由 `systemd` 管理（`ez-admin.service`）
- PostgreSQL 和 Redis 运行在 Docker 中，端口仅绑定 `127.0.0.1`
- `compose.server.yml` 中 Nginx 使用 `network_mode: host`，因此可以通过 `127.0.0.1:8080` 访问宿主机上的后端

## 前置条件

### 本地机器

| 要求 | 说明 |
|------|------|
| Git | 克隆仓库 |
| Go 1.23+ | 编译后端 |
| Node.js 20+ & pnpm | 构建前端 |
| bash | `pack.sh` / `deploy.sh` 需要 bash 环境 |
| SSH 访问 | 连接服务器 |

Windows 用户可使用 `pack.ps1`（PowerShell 脚本）。

### 服务器

| 要求 | 说明 |
|------|------|
| Linux（建议 Ubuntu 22.04+） | 示例以 Ubuntu 为准 |
| Docker + Docker Compose | 运行数据库、Redis、Nginx |
| systemd | 管理后端服务 |
| sudo 权限 | `setup-server.sh` 内部使用 `sudo cp`、`sudo systemctl` 等 |
| 开放 80 / 443 端口 | HTTP / HTTPS 访问 |

::: warning 远程用户需要 sudo 权限
`scripts/deploy.sh` 会远程执行 `setup-server.sh`，而 `setup-server.sh` 内部会使用 `sudo cp`、`sudo systemctl` 等命令。如果服务器禁用了交互式 sudo（如要求输入密码），你需要手动登录服务器执行初始化命令。
:::

::: warning Docker 需要提前安装
`setup-server.sh` 不会自动安装 Docker。如果服务器还没有 Docker 环境，需要你提前准备好。
:::

## 推荐路径：一键部署

在项目根目录执行：

```bash
bash scripts/deploy.sh user@your-server-ip
```

`deploy.sh` 会自动完成以下全部操作：

1. **编译后端** — 交叉编译为 Linux amd64 二进制（`-ldflags="-s -w"` 去除调试信息）
2. **构建前端** — `pnpm build`，产物输出到 `dist/`
3. **打包配置** — 复制 `compose.server.yml`、`.env.example`、`ez-admin.service`、Nginx 配置等
4. **生成压缩包** — `deploy-package.tar.gz`
5. **上传到服务器** — 通过 `scp` 上传到 `/tmp/`
6. **远端初始化** — 解压到 `/opt/ez-admin`，执行 `setup-server.sh`

部署完成后终端会输出：

```
=========================================
✅ 部署完成！

  访问地址：http://your-server-ip
  默认账号：admin / Admin@123456

  查看后端日志：sudo journalctl -u ez-admin -f
  查看容器状态：docker compose -f /opt/ez-admin/compose.server.yml ps
=========================================
```

::: tip 如果你更喜欢手动控制
可以使用 `pack.sh` 单独打包，然后手动上传和初始化。详见下方"手动操作"章节。
:::

## 理解路径：脚本背后做了什么

`setup-server.sh` 在服务器上按顺序执行以下操作：

### 1. 整理文件

- 创建 `/opt/ez-admin/nginx`、`ssl`、`data/postgres`、`data/redis` 目录
- 将 `dist/` 重命名为 `web/`（Nginx 配置中的 `root` 路径）
- 将 `.env.example` 复制为 `.env`（仅首次，不会覆盖已有配置）
- 将 `ez-admin.service` 复制到 `/etc/systemd/system/`
- 赋予后端二进制可执行权限

### 2. 生成 JWT 密钥

如果 `.env` 中的 `EZ_AUTH_JWT_SECRET` 仍为默认值 `change-me-to-a-random-string-at-least-32-chars`，脚本会自动用 `openssl rand -hex 32` 生成随机密钥并替换。

### 3. 启动基础设施

```bash
docker compose -f compose.server.yml up -d
```

启动 PostgreSQL、Redis 和 Nginx 容器。脚本会等待数据库就绪（最多 30 秒）。

### 4. 启动后端

```bash
sudo systemctl enable --now ez-admin
```

通过 systemd 启动后端服务。脚本会等待后端健康检查通过（最多 15 秒）。

### 5. 初始化管理员

脚本尝试调用 `POST /api/v1/setup/init` 创建默认管理员。如果返回 `200` 表示创建成功，`409` 表示管理员已存在则跳过。

## 部署目录结构

部署完成后，服务器上的目录结构如下：

```
/opt/ez-admin/
├── server                    # Go 后端二进制
├── web/                      # 前端构建产物（原 dist/）
│   ├── index.html
│   └── assets/
├── compose.server.yml        # Docker Compose 配置
├── .env                      # 环境变量（从 .env.example 生成）
├── ez-admin.service          # systemd 服务文件（已复制到 /etc/systemd/system/）
├── configs/
│   ├── config.yaml           # 后端配置文件
│   └── rbac_model.conf       # Casbin RBAC 模型
├── nginx/
│   └── nginx-native.conf     # Nginx HTTP 配置
├── ssl/                      # SSL 证书目录（初始为空）
├── data/
│   ├── postgres/             # PostgreSQL 数据卷
│   └── redis/                # Redis 数据卷
├── logs/                     # 后端运行日志（如有配置）
└── uploads/                  # 用户上传文件
```

::: tip 数据持久化
数据库数据存储在 `/opt/ez-admin/data/postgres` 和 `data/redis`，是 Docker 绑定挂载目录。删除容器不会丢失数据，但删除这些目录会。
:::

## 环境变量配置

`.env` 文件由 `setup-server.sh` 从 `.env.example` 自动创建。以下变量来自 `deploy/.env.example`：

### 必须修改

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `EZ_AUTH_JWT_SECRET` | JWT 签名密钥 | 脚本会自动生成；也可手动设置 |
| `EZ_DATABASE_PASSWORD` | 数据库密码 | `ez_admin_123456` |

::: danger 生产环境安全保护
`EZ_APP_ENV=prod` 且 `EZ_AUTH_JWT_SECRET` 包含 `change-me` 时，后端会拒绝启动。这是一项安全保护，不要绕过。
:::

### 建议修改

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `EZ_REDIS_PASSWORD` | Redis 密码 | 空（无密码） |
| `EZ_CORS_ALLOWED_ORIGINS` | CORS 白名单，设为实际域名 | 空 |
| `EZ_LOG_LEVEL` | 日志级别 | `info` |

### 其他变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `EZ_DATABASE_HOST` | 数据库地址 | `127.0.0.1` |
| `EZ_DATABASE_PORT` | 数据库端口 | `5432` |
| `EZ_DATABASE_USER` | 数据库用户 | `ez_admin` |
| `EZ_DATABASE_NAME` | 数据库名 | `ez_admin` |
| `EZ_REDIS_HOST` | Redis 地址 | `127.0.0.1` |
| `EZ_REDIS_PORT` | Redis 端口 | `6379` |
| `EZ_APP_ENV` | 运行环境 | `prod` |
| `EZ_SERVER_ADDR` | 后端监听地址 | `:8080` |
| `EZ_LOG_FORMAT` | 日志格式 | `json` |
| `EZ_RATE_LIMIT_LOGIN_MAX_REQUESTS` | 登录限流最大请求数 | `10` |
| `EZ_RATE_LIMIT_LOGIN_WINDOW_SEC` | 登录限流窗口（秒） | `60` |

完整变量说明见 [环境变量参考](/reference/environment-variables-reference)。

## 默认管理员初始化

`setup-server.sh` 会在部署完成后自动尝试创建默认管理员：

- 用户名：`admin`
- 密码：`Admin@123456`
- 昵称：`管理员`

::: danger 生产部署后必须立即修改默认密码
首次部署脚本会尝试创建默认管理员账号 `admin / Admin@123456`。
部署完成后请立即登录后台修改密码，或在正式开放公网前完成初始化。
不要在公网长期保留默认密码。
:::

如果初始化返回 `409`，表示管理员已存在，脚本会自动跳过。

如果初始化失败（返回非 `200`/`409` 状态码），脚本会打印手动初始化命令：

```bash
curl -X POST http://localhost/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin@123456","nickname":"管理员"}'
```

## 验证部署

部署完成后，逐项检查：

### 服务状态

```bash
# 后端 systemd 服务
sudo systemctl status ez-admin
# 期望输出：Active: active (running)

# Docker 容器
cd /opt/ez-admin
docker compose -f compose.server.yml ps
# 期望输出：三个容器均为 Up / running (healthy)
```

### HTTP 访问

```bash
# 健康检查
curl -sf http://localhost/health
# 期望：返回 200

# API 代理
curl -sf http://localhost/api/v1/health
# 期望：后端正常响应

# 前端页面
curl -sf http://localhost/
# 期望：返回 HTML 页面
```

### 浏览器验证

1. 访问 `http://your-server-ip` — 应看到登录页
2. 使用 `admin / Admin@123456` 登录 — 应进入后台首页
3. 检查左侧菜单是否正常显示
4. 点击菜单项，确认页面正常加载
5. 尝试上传文件，确认上传目录有写入权限

### 日志检查

```bash
# 后端日志
sudo journalctl -u ez-admin -n 50

# Nginx 日志
docker compose -f compose.server.yml logs nginx --tail 20
```

## 手动操作（不使用一键部署）

如果你倾向于手动控制每一步，或服务器禁用了远程 sudo，可以分步操作：

### 1. 本地打包

::: code-group

```bash [macOS / Linux]
bash scripts/pack.sh
```

```powershell [Windows]
.\scripts\pack.ps1
```

:::

### 2. 上传到服务器

```bash
scp deploy-package.tar.gz user@your-server:/tmp/
```

### 3. 服务器初始化

```bash
# SSH 登录服务器
ssh user@your-server

# 创建目录并解压
sudo mkdir -p /opt/ez-admin
cd /opt/ez-admin
sudo tar xzf /tmp/deploy-package.tar.gz

# 执行首次部署脚本
sudo bash setup-server.sh
```

## 常见问题

### 远程 sudo 权限问题

如果 `deploy.sh` 在远端执行 `setup-server.sh` 时因 sudo 权限失败：

```bash
# 方案一：手动登录服务器执行
ssh user@your-server
cd /opt/ez-admin
sudo bash setup-server.sh

# 方案二：配置服务器免密 sudo（仅限受控环境）
# 在服务器上执行：sudo visudo
# 添加：your-user ALL=(ALL) NOPASSWD: ALL
```

### Docker 未安装或未启动

```bash
# 检查 Docker 状态
docker info

# 如果未安装（Ubuntu）
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER
# 重新登录使用户组生效
```

### 80 / 443 端口被占用

```bash
# 检查端口占用
sudo ss -tlnp | grep -E ':80|:443'

# 如果是其他服务占用，停止该服务或修改 Nginx 监听端口
```

### Nginx 能访问但 API 404 / 502

```bash
# 检查后端是否运行
sudo systemctl status ez-admin

# 检查后端健康
curl http://localhost:8080/health

# 如果后端正常但 Nginx 502，检查 Nginx 容器网络模式
docker inspect ez-admin-nginx | grep NetworkMode
# 应该是 "host"
```

### 后端启动失败

```bash
# 查看后端日志
sudo journalctl -u ez-admin -n 50 --no-pager

# 常见原因：
# 1. JWT 密钥未修改 → 编辑 .env 修改 EZ_AUTH_JWT_SECRET
# 2. 数据库未就绪   → 等待几秒后重启：sudo systemctl restart ez-admin
# 3. 端口被占用     → 检查 8080 端口：sudo ss -tlnp | grep 8080
```

### 数据库连接失败

```bash
# 检查 PostgreSQL 容器
docker compose -f compose.server.yml ps postgres

# 检查数据库是否就绪
docker compose -f compose.server.yml exec postgres pg_isready -U ez_admin

# 检查 .env 中的数据库配置
grep EZ_DATABASE /opt/ez-admin/.env
```

### 默认管理员无法初始化

```bash
# 检查后端是否已就绪
curl -sf http://localhost/health

# 手动初始化
curl -X POST http://localhost/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin@123456","nickname":"管理员"}'

# 返回 409 表示管理员已存在，正常
```

### 前端刷新 404

Nginx 配置已包含 SPA fallback（`try_files $uri $uri/ /index.html`）。如果仍然出现 404，检查 Nginx 配置是否正确挂载：

```bash
docker exec ez-admin-nginx cat /etc/nginx/conf.d/default.conf
# 确认包含 try_files 配置
```

### 上传目录权限问题

```bash
# 检查上传目录
ls -la /opt/ez-admin/uploads/

# 确保后端进程有写入权限
sudo chown -R $(whoami) /opt/ez-admin/uploads/
```

## 如何复用为自己的项目

EZ Admin 作为通用后台底座，你可以基于它构建自己的项目：

### 需要修改的部分

| 类别 | 修改位置 | 说明 |
|------|---------|------|
| 项目名称 | 仓库名、README | 全局替换项目名称 |
| Logo / 标题 / 品牌信息 | `admin/` | 前端品牌组件和页面标题 |
| 域名 | Nginx 配置、`.env` | `server_name`、`EZ_CORS_ALLOWED_ORIGINS` |
| 数据库名和密码 | `.env` | `EZ_DATABASE_NAME`、`EZ_DATABASE_PASSWORD` |
| 初始菜单 seed | `server/migrations/` | 种子数据 SQL |
| 初始角色权限 | `server/migrations/` | Casbin 策略种子 |
| 默认管理员策略 | `setup-server.sh` | 初始用户名和密码 |
| CORS 白名单 | `.env` | `EZ_CORS_ALLOWED_ORIGINS` |
| Nginx server_name | `deploy/nginx/` | Nginx 配置文件中的 `server_name` |
| 前端 API 地址 | `admin/.env*` | Vite 代理配置或生产 API 地址 |

### 可以直接保留的部分

- RBAC 权限底座（JWT + Casbin + 数据权限）
- 动态菜单和按钮权限机制
- 组织架构（部门 / 岗位）
- 五级数据权限
- 用户 / 角色 / 部门 / 岗位基础模块
- 后端四层分层结构（handler / service / repository / domain）
- 前端路由和权限结构
- 部署脚本和配置

::: tip 复用流程
Fork 仓库 → 修改品牌和配置 → 在 `server/` 中添加业务模块 → 重新 `pack.sh` 打包部署。部署流程本身不需要修改。
:::
