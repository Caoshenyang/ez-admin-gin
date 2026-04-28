---
title: 第 7 章：部署与复用
description: "编译后端二进制、构建前端静态文件，部署到腾讯云轻量服务器，通过 Cloudflare 配置 HTTPS 域名访问。"
---

# 第 7 章：部署与复用

前六章完成了后台底座的核心功能。这一章把它部署到公网，通过 HTTPS 域名对外提供服务。

::: tip 🎯 本章怎么验证
完成后你能通过 `https://你的域名` 访问后台、登录并执行 CRUD。
:::

## 部署思路

部署方式很直接：

- Docker Compose 只负责 PostgreSQL 和 Redis（基础环境）
- 后端编译成 Linux 二进制，直接在服务器上运行
- 前端构建成静态文件，由 Nginx 托管

不需要 Docker Hub，不需要构建镜像，不需要在服务器上安装 Go 或 Node.js。

## 开始前准备

- 本机有 Go 1.22+ 和 Node.js 22+
- 有一个域名
- 本机已安装 Docker 和 Docker Compose（V2）

## 本章会完成什么

| 小节 | 定位 | 完成什么 |
| --- | --- | --- |
| [部署验证与复用说明](./deployment-and-reuse) | 主线实操 | 本地构建 → 服务器部署 → Cloudflare HTTPS → 更新与复用 |
| [环境变量与初始化数据](./env-and-init-data) | 参考手册 | 环境变量完整清单、覆盖机制、迁移文件结构 |

## 本章完成后的状态

```text
本地编译后端二进制 + 构建前端静态文件
  ↓
上传到腾讯云轻量服务器
  ↓
Cloudflare 托管域名 + HTTPS 上线
  ↓
后续更新：重新编译 → 上传 → 重启
```

开始部署：[部署验证与复用说明](./deployment-and-reuse)。
