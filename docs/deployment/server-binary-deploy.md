---
title: 服务器二进制部署
description: 从本地打包到服务器上线，包含完整步骤、环境变量配置和常见问题
---

# 服务器二进制部署

本文介绍如何将 EZ Admin 部署到一台 Linux 服务器上。后端以二进制方式通过 systemd 管理，前端由 Nginx 容器托管，数据库和 Redis 使用 Docker 运行。

::: tip 读完本文你可以
- 在本地一键打包部署文件
- 在服务器上完成首次部署并访问系统
- 了解每个步骤的命令含义和预期结果
:::

## 部署架构

```
客户端浏览器
    │
    ▼
┌──────────────────────────────────────┐
│  Nginx 容器 (host 网络)              │
│  ├── /          → 前端静态文件 (web/) │
│  ├── /api/*    → 127.0.0.1:8080      │
│  └── /uploads/* → 127.0.0.1:8080     │
└──────────────────────────────────────┘
    │
    ▼
┌──────────────────┐
│  Go 后端 (systemd) │  ← 二进制直接运行，端口 8080
│  监听 :8080       │
└──────────────────┘
    │           │
    ▼           ▼
┌─────────┐  ┌───────┐
│ PostgreSQL │  │ Redis  │  ← Docker 容器，仅监听 127.0.0.1
│ :5432     │  │ :6379  │
└─────────┘  └───────┘
```

核心特点：

- Nginx 使用 `network_mode: host`，共享宿主机网络栈，`proxy_pass` 直接访问 `127.0.0.1:8080`
- PostgreSQL 和 Redis 仅绑定 `127.0.0.1`，外部无法直连
- 后端通过 systemd 管理，崩溃自动重启

## 前置条件

### 本地机器

| 要求 | 说明 |
|------|------|
| Go 1.23+ | 编译后端 |
| Node.js 20+ & pnpm | 构建前端 |
| tar / scp | 打包和上传（macOS / Linux） |
| SSH 访问 | 连接服务器 |

### 服务器

| 要求 | 说明 |
|------|------|
| Ubuntu 22.04+ 或其他主流 Linux 发行版 | 示例以 Ubuntu 为准 |
| Docker + Docker Compose | 运行数据库、Redis、Nginx |
| 至少 1GB 内存 | 数据库 + Redis + 后端 + Nginx |
| 开放 80 / 443 端口 | HTTP / HTTPS 访问 |

::: warning
Docker 需要已安装并启动。如果服务器还没有 Docker，`setup-server.sh` 不会自动安装，需要你提前准备好。
:::

## 本地打包

在项目根目录执行打包脚本：

::: code-group
```bash [macOS / Linux]
bash scripts/pack.sh
```
```powershell [Windows]
.\scripts\pack.ps1
```
:::

脚本会依次完成：

1. **编译后端** — 交叉编译为 Linux amd64 二进制（`-ldflags="-s -w"` 去除调试信息）
2. **构建前端** — `pnpm build`，产物输出到 `dist/`
3. **打包配置** — 复制 `compose.server.yml`、`.env.example`、`ez-admin.service`、Nginx 配置等
4. **生成压缩包** — `deploy-package.tar.gz`（或 Windows 下的 `.zip`）

打包完成后会看到：

```
✅ 打包完成！上传 deploy-package.tar.gz 到服务器即可。
```

### 打包产物结构

```
deploy-package/
├── server                      # Go 后端二进制
├── dist/                       # 前端构建产物
├── compose.server.yml          # Docker Compose 配置
├── .env.example                # 环境变量模板
├── ez-admin.service            # systemd 服务文件
├── setup-server.sh             # 首次部署脚本
├── update-server.sh            # 更新脚本
├── nginx/
│   └── nginx-native.conf       # Nginx 配置
├── configs/
│   ├── config.yaml             # 后端配置文件
│   └── rbac_model.conf         # Casbin RBAC 模型
└── ssl/                        # SSL 证书目录（空）
```

## 上传部署包

::: code-group
```bash [手动上传]
scp deploy-package.tar.gz user@your-server:/tmp/
```
```bash [一键部署]
bash scripts/deploy.sh user@your-server
```
:::

一键部署脚本 `deploy.sh` 会自动完成打包、上传和远端初始化。如果选择手动上传，接下来需要在服务器上操作。

## 服务器初始化

SSH 登录服务器后执行：

```bash
# 创建部署目录并解压
mkdir -p /opt/ez-admin
cd /opt/ez-admin
tar xzf /tmp/deploy-package.tar.gz

# 执行首次部署脚本
sudo bash setup-server.sh
```

`setup-server.sh` 会完成以下操作：

### 1. 整理文件

- 将 `dist/` 重命名为 `web/`（Nginx 配置中的 root 路径）
- 将 `.env.example` 复制为 `.env`（仅在首次，不会覆盖已有配置）
- 注册 systemd 服务
- 赋予后端二进制可执行权限

### 2. 生成 JWT 密钥

如果 `.env` 中的 JWT 密钥仍为默认值 `change-me-to-a-random-string-at-least-32-chars`，脚本会自动用 `openssl rand -hex 32` 生成一个随机密钥并替换。

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

脚本自动创建默认管理员账号：

- 用户名：`admin`
- 密码：`Admin@123456`

::: danger
首次登录后请立即修改默认密码。
:::

## 环境变量配置

部署脚本会自动从 `.env.example` 创建 `.env`。如果需要自定义配置：

```bash
cd /opt/ez-admin
cp .env.example .env   # 仅首次，已存在则跳过
vim .env               # 编辑配置
```

### 必须修改的变量

| 变量 | 说明 |
|------|------|
| `EZ_AUTH_JWT_SECRET` | JWT 签名密钥，脚本会自动生成，也可手动设置 |
| `EZ_DATABASE_PASSWORD` | 数据库密码，生产环境应使用强密码 |

### 建议修改的变量

| 变量 | 说明 |
|------|------|
| `EZ_REDIS_PASSWORD` | Redis 密码 |
| `EZ_CORS_ALLOWED_ORIGINS` | CORS 白名单，设为实际域名 |

::: warning
`EZ_APP_ENV=prod` 且 JWT 密钥包含 `change-me` 时，后端会拒绝启动。这是一项安全保护。
:::

完整变量说明见 [环境变量参考](/reference/environment-variables-reference)。

## 手动操作（不使用脚本）

如果你倾向于手动控制每一步，可以替代 `setup-server.sh`：

```bash
# 1. 创建目录
cd /opt/ez-admin
mkdir -p data/postgres data/redis nginx web ssl

# 2. 整理文件
mv dist web

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env，至少修改 JWT 密钥和数据库密码

# 4. 启动数据库和 Redis
docker compose -f compose.server.yml up -d

# 5. 等待数据库就绪
docker compose -f compose.server.yml exec -T postgres pg_isready -U ez_admin

# 6. 注册并启动后端
sudo cp ez-admin.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now ez-admin

# 7. 初始化管理员
curl -X POST http://localhost/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

## 访问验证

部署完成后，在浏览器中访问：

```
http://your-server-ip
```

使用 `admin / Admin@123456` 登录。

### 验证服务状态

```bash
# 检查后端服务
sudo systemctl status ez-admin

# 检查 Docker 容器
cd /opt/ez-admin
docker compose -f compose.server.yml ps

# 查看后端日志
sudo journalctl -u ez-admin -f
```

所有服务正常时，`systemctl status` 应显示 `active (running)`，三个容器均应 `Up`。

## 常见问题

### 后端启动失败：JWT 密钥未修改

```
Error: JWT secret must be changed in production environment
```

编辑 `/opt/ez-admin/.env`，将 `EZ_AUTH_JWT_SECRET` 替换为随机字符串：

```bash
# 生成新密钥
openssl rand -hex 32
# 更新 .env
sed -i 's/change-me-to-a-random-string-at-least-32-chars/你生成的密钥/' .env
# 重启后端
sudo systemctl restart ez-admin
```

### 数据库连接失败

确认 Docker 容器正在运行且端口绑定正确：

```bash
docker compose -f compose.server.yml ps
# PostgreSQL 应绑定在 127.0.0.1:5432
```

### Nginx 无法访问后端

检查后端是否在运行：

```bash
curl http://localhost:8080/health
# 应返回 200
```

Nginx 使用 host 网络，`proxy_pass` 直连 `127.0.0.1:8080`。如果后端正常但 Nginx 502，检查 Nginx 容器是否使用了 `network_mode: host`。

### 管理员初始化返回 409

表示管理员账号已存在，这是正常的。只有在首次部署时才会返回 200。

## 复用为自己的项目

EZ Admin 作为通用后台底座，你可以基于它构建自己的项目：

1. **Fork 仓库**，修改项目名称和品牌信息
2. **修改前端** — 品牌名、Logo、主题色在 `admin/` 中配置
3. **扩展后端** — 在 `server/` 中添加业务模块（参见 [模块开发](/backend/module-development)）
4. **重新打包部署** — `pack.sh` 会自动包含你的改动

部署流程不需要任何修改，脚本和配置文件直接复用。需要自定义的只有 `.env` 中的环境变量。
