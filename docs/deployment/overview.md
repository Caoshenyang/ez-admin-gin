---
title: 部署概览
description: 五种部署变体、Docker Compose 结构、生产部署要点
---

# 部署概览

## 部署变体

EZ Admin 提供五种部署方案，覆盖从开发到生产的全部场景：

| 变体 | 配置文件 | 适用场景 |
|------|---------|---------|
| 本地开发 | `deploy/compose.local.yml` | macOS/Linux 开发环境 |
| Windows 开发 | `deploy/compose.local.win.yml` | Windows 开发环境 |
| 服务器部署 | `deploy/compose.server.yml` | 后端二进制运行在宿主机 |
| 云端部署 | `deploy/compose.deploy.yml` | 全容器化，Docker Hub 镜像 |
| 生产部署 | `deploy/compose.prod.yml` | 生产环境 |

## 服务组成

每个部署变体包含以下服务：

| 服务 | 用途 | 端口 |
|------|------|------|
| PostgreSQL / MySQL | 主数据库 | 5432 / 3306 |
| Redis | 缓存 + 限流 | 6379 |
| 后端 (Go) | API 服务 | 8080 |
| 前端 (Nginx) | 静态资源 + 反向代理 | 80 / 443 |

## 部署文件结构

```
deploy/
├── compose.local.yml          本地开发
├── compose.local.win.yml      Windows 开发
├── compose.server.yml         服务器（二进制 + Docker 基础设施）
├── compose.deploy.yml         云端（全容器化）
├── compose.prod.yml           生产环境
├── .env.example               环境变量模板
├── nginx.conf                 HTTP 反向代理
├── nginx-ssl.conf             HTTPS 配置
├── nginx-native.conf          宿主机 Nginx（HTTP）
├── nginx-native-ssl.conf      宿主机 Nginx（HTTPS）
└── ez-admin.service           Systemd 服务配置
```

## 部署脚本

| 脚本 | 用途 |
|------|------|
| `scripts/setup-server.sh` | 服务器初始化（安装 Docker、创建目录等） |
| `scripts/deploy.sh` | 一键部署到远程服务器 |
| `scripts/update-server.sh` | 更新已部署的服务 |
| `scripts/pack.sh` | 构建并打包（Linux/macOS） |
| `scripts/pack.ps1` | 构建并打包（Windows） |

## Nginx 配置

所有 Nginx 配置包含：

- SPA fallback（`try_files $uri $uri/ /index.html`）
- `/api` 反向代理到后端
- 静态资源 1 年缓存
- Gzip 压缩
- 安全响应头（X-Frame-Options、X-Content-Type-Options 等）
- HTTPS 版本支持 SSL/TLS

## 环境变量

关键环境变量（完整列表见 `deploy/.env.example`）：

| 变量 | 说明 | 必须修改 |
|------|------|---------|
| `EZ_APP_ENV` | 运行环境 | 生产环境设为 `prod` |
| `EZ_DATABASE_HOST` | 数据库地址 | ✅ |
| `EZ_DATABASE_PASSWORD` | 数据库密码 | ✅ |
| `EZ_AUTH_JWT_SECRET` | JWT 密钥 | ✅ 生产必须修改 |
| `EZ_REDIS_HOST` | Redis 地址 | ✅ |
| `EZ_CORS_ALLOWED_ORIGINS` | CORS 白名单 | ✅ 生产必须配置 |

::: warning
生产环境必须修改 JWT 密钥。如果 `EZ_APP_ENV=prod` 且密钥包含 `change-me`，服务将拒绝启动。
:::

## 快速部署流程

### 服务器部署（二进制 + Docker 基础设施）

```bash
# 1. 初始化服务器
./scripts/setup-server.sh user@your-server

# 2. 打包
./scripts/pack.sh

# 3. 部署
./scripts/deploy.sh user@your-server

# 4. 初始化管理员账号
curl -X POST http://your-server:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

### 全容器化部署

```bash
# 1. 配置环境变量
cp deploy/.env.example deploy/.env
# 编辑 .env 填入实际配置

# 2. 启动所有服务
docker compose -f deploy/compose.deploy.yml --env-file deploy/.env up -d

# 3. 初始化管理员账号（同上）
```

## 生产环境检查清单

- [ ] JWT 密钥已修改（`openssl rand -hex 32`）
- [ ] 数据库密码已修改
- [ ] CORS 白名单已配置为实际域名
- [ ] HTTPS 证书已配置
- [ ] Swagger 已关闭（`EZ_SWAGGER_ENABLED=false`）
- [ ] 日志级别调整为 `warn` 或 `info`
- [ ] 文件上传目录权限正确
- [ ] 数据库备份策略已配置
