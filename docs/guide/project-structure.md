---
title: 项目结构
description: "EZ Admin Gin 的技术栈组成、当前目录结构，以及 v2 企业级底座阶段的目标结构说明。"
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
│   ├── cmd/         # v2 启动入口
│   ├── configs/     # 配置文件（config.yaml）
│   ├── internal/    # 启动装配、平台能力和业务模块
│   └── migrations/  # 数据库迁移
├── admin/           # Vue 3 前端
│   └── src/         # 页面、组件、路由、状态管理
├── docs/            # VitePress 文档站
├── deploy/          # Docker Compose 和 Nginx 配置
│   └── nginx/
└── .agents/         # 开发辅助 skill
```

各目录职责：

- **server/** — Go 后端，当前同时保留兼容入口 `main.go` 和 v2 入口 `cmd/server/`
- **admin/** — Vue 3 前端管理台，页面、组件、路由和状态管理都在 `src/` 下
- **docs/** — VitePress 文档站，就是你现在在读的站点
- **deploy/** — Docker Compose 文件和 Nginx 反向代理配置，分为本地开发环境和生产环境
- **.agents/** — 开发辅助工具配置，正常使用不需要关注

## v2 阶段的目标结构

`v1` 阶段的重点是把后台系统完整跑通，`v2` 阶段的重点是把它升级成可长期扩展的企业级单体底座。后端会逐步收敛到下面这个结构：

```text
server/
├── cmd/
│   └── server/
├── internal/
│   ├── bootstrap/
│   ├── platform/
│   │   ├── config/
│   │   ├── database/
│   │   ├── logger/
│   │   ├── redis/
│   │   ├── authn/
│   │   ├── authz/
│   │   └── datascope/
│   └── module/
│       ├── auth/
│       ├── iam/
│       ├── org/
│       ├── system/
│       ├── dict/
│       └── account/
```

这样调整的核心原因只有一个：让组织体系、数据权限和后续模块扩展有明确落点。

## 当前第一阶段已经落地了什么

这一轮升级先做“结构先行、行为尽量不变”的改造：

- 新增 `bootstrap`，把启动装配从 `main.go` 中抽出来
- 新增 `platform` 命名空间，承接配置、数据库、日志、Redis、认证、鉴权和数据权限基础设施
- 新增 `module/auth`、`module/setup`、`module/system` 路由聚合入口
- 新增部门、岗位、用户岗位、角色数据范围的模型与迁移

这一步还不代表全部业务逻辑都已经迁到 `service / repository`，但后续的结构落点已经建立好了。

## 怎么继续读

- 想理解为什么一定要做这次结构升级：看 [企业级架构升级](/guide/enterprise-architecture)
- 想从 Java 视角理解这套结构：看 [Go vs Java 工程结构](/guide/java-to-go-structure)
