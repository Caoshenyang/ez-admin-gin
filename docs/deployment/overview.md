---
title: 部署概览
description: 部署方式总览、适用场景与文件结构
---

# 部署概览

EZ Admin 提供多种部署方式，覆盖本地开发到生产上线。本文帮你快速选择合适的方案，并了解部署目录和脚本的作用。

## 如何选择部署方式

| 你的场景 | 推荐方案 | 配置文件 |
|---------|---------|---------|
| 本地开发（macOS / Linux） | Docker Compose 本地环境 | `compose.local.yml` |
| 本地开发（Windows） | Docker Compose 本地环境 | `compose.local.win.yml` |
| 一台 VPS / 云服务器上线 | 服务器二进制部署（推荐） | `compose.server.yml` |
| 全容器化部署 | Docker 部署 | `compose.deploy.yml` |
| 需要域名和 HTTPS | 域名与 HTTPS | `nginx-native-ssl.conf` |
| 准备正式上线 | 生产环境检查清单 | — |
| 已部署，需要更新 | 更新与回滚 | — |

::: tip 新用户推荐
如果你是第一次把 EZ Admin 部署到服务器，优先阅读 [服务器二进制部署](/deployment/server-binary-deploy)。

它提供了一键打包脚本和自动初始化，流程最短，也最容易排查问题。
:::

## 部署方式对比

| | 服务器二进制部署 | Docker 全容器化部署 |
|---|---|---|
| **后端运行方式** | 宿主机二进制 + systemd | Docker 容器 |
| **前端运行方式** | Nginx 容器挂载静态文件 | Nginx 容器（镜像内置） |
| **数据库 / Redis** | Docker 容器 | Docker 容器 |
| **配置文件** | `compose.server.yml` + `.env` | `compose.deploy.yml` + `.env` |
| **适用场景** | 有 SSH 访问的 VPS、个人服务器 | Docker Hub 镜像分发、CI/CD 流水线 |
| **推荐人群** | 个人项目、小团队 | 需要镜像版本管理的场景 |

## 部署文件结构

```
deploy/
├── compose.local.yml          # 本地开发（macOS / Linux）
├── compose.local.win.yml      # 本地开发（Windows）
├── compose.server.yml         # 服务器部署（二进制 + Docker 基础设施）
├── compose.deploy.yml         # 全容器化部署（Docker Hub 镜像）
├── compose.prod.yml           # 生产环境演示配置
├── .env.example               # 环境变量模板
├── nginx/
│   ├── nginx.conf             # 容器内 Nginx（HTTP）
│   ├── nginx-ssl.conf         # 容器内 Nginx（HTTPS）
│   ├── nginx-native.conf      # 宿主机 Nginx（HTTP）
│   └── nginx-native-ssl.conf  # 宿主机 Nginx（HTTPS）
└── ez-admin.service           # Systemd 服务配置
```

## 部署脚本

| 脚本 | 运行位置 | 用途 |
|------|---------|------|
| `scripts/pack.sh` | 本地（macOS / Linux） | 编译后端、构建前端、打包部署文件 |
| `scripts/pack.ps1` | 本地（Windows） | 同上，输出 `.zip` |
| `scripts/deploy.sh` | 本地 | 一键打包 + 上传 + 远端初始化 |
| `scripts/setup-server.sh` | 服务器 | 首次部署初始化（创建目录、生成密钥、启动服务） |
| `scripts/update-server.sh` | 服务器 | 更新已部署的服务（替换文件、重启后端） |

## 部署流程概览

无论选择哪种部署方式，核心流程一致：

```
本地打包 → 上传到服务器 → 配置环境变量 → 启动服务 → 初始化管理员
```

## 接下来

- [服务器二进制部署](/deployment/server-binary-deploy) — 推荐新用户阅读
- [Docker 部署](/deployment/docker-deploy) — 全容器化方案
- [域名与 HTTPS](/deployment/domain-and-https) — 配置域名和证书
- [更新与回滚](/deployment/update-and-rollback) — 版本更新流程
- [生产环境检查清单](/deployment/production-checklist) — 上线前必查项
