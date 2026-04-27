---
title: 项目结构
description: "EZ Admin Gin 的技术栈组成和目录结构说明。"
---

# 项目结构

::: tip 🎯 这页解决什么
帮你快速了解项目用了哪些技术、各目录负责什么，方便后续定位文件。
:::

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Go 1.26、Gin、GORM、PostgreSQL、Redis、Casbin |
| 前端 | Vue 3.5、TypeScript、Naive UI、TailwindCSS 4、Vite 8 |
| 文档 | VitePress 2.0 |
| 部署 | Docker Compose、Nginx |

## 目录结构

```
ez-admin-gin/
├── server/          # Go 后端
│   ├── configs/     # 配置文件（config.yaml）
│   └── internal/    # 业务代码（model / repository / service / handler / router）
├── admin/           # Vue 3 前端
│   └── src/         # 页面、组件、路由、状态管理
├── docs/            # VitePress 文档站
├── deploy/          # Docker Compose 和 Nginx 配置
│   └── nginx/
└── .agents/         # 开发辅助 skill
```

各目录职责：

- **server/** — Go 后端，入口是 `main.go`，业务代码按 `model → repository → service → handler → router` 分层
- **admin/** — Vue 3 前端管理台，页面、组件、路由和状态管理都在 `src/` 下
- **docs/** — VitePress 文档站，就是你现在在读的站点
- **deploy/** — Docker Compose 文件和 Nginx 反向代理配置，分为本地开发环境和生产环境
- **.agents/** — 开发辅助工具配置，正常使用不需要关注
