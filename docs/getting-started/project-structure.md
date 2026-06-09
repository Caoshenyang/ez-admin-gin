---
title: 项目结构
description: "按当前代码说明 EZ Admin Gin 的仓库目录、模块边界、文档站入口和扩展时应该优先查看的位置。"
---

# 项目结构

这一页帮你把仓库先看成一张地图：哪里是后端，哪里是前端，哪里是文档，新增模块时应该沿哪条主线走。

::: tip 🎯 这页解决什么
读完后，你应该能根据一个需求快速定位到 `server/`、`admin/`、`docs/`、`deploy/` 或 `scripts/` 中对应的文件。
:::

## 顶层目录

```text
ez-admin-gin/
├── server/          # 后端服务：Gin、GORM、Casbin、迁移、Swagger
├── admin/           # 前端管理台：Vue 3、Naive UI、动态菜单
├── docs/            # VitePress 文档站
├── deploy/          # Docker Compose、Nginx、systemd 和环境变量模板
├── scripts/         # 打包、部署、更新、Swagger 与 SQL 生成脚本
├── brand-assets/    # 品牌资源源文件与导出产物
├── Makefile         # 统一开发入口
└── MANUAL_TEST.md   # 发布前人工冒烟检查清单
```

## 当前技术栈

| 区域 | 当前代码依据 | 说明 |
| --- | --- | --- |
| 后端 | `server/go.mod` | Go 1.26.1、Gin、GORM、Casbin、Redis、Prometheus、Swagger |
| 前端 | `admin/package.json` | Vue 3.5、Vue Router 5、Pinia 3、Naive UI、Vite 8、TypeScript 6、Tailwind CSS 4 |
| 文档 | `docs/package.json` | VitePress 2.0 alpha、Vue 3.5 |
| 部署 | `deploy/`、`scripts/` | Docker Compose、Nginx、systemd、二进制部署与全容器部署 |

## 后端目录

```text
server/
├── main.go
├── configs/
│   ├── config.yaml
│   └── rbac_model.conf
├── docs/                  # swag 生成的 OpenAPI / Swagger 文件
├── internal/
│   ├── bootstrap/         # 启动装配、路由聚合、Swagger 注册
│   ├── modules/           # 业务模块
│   ├── platform/          # 跨模块基础设施
│   └── pkg/               # 小型公共工具包
└── migrations/
    ├── mysql/
    └── postgres/
```

后端当前模块边界如下：

| 分组 | 目录 | 当前职责 |
| --- | --- | --- |
| Auth | `server/internal/modules/auth` | 登录、刷新令牌、退出、当前用户、账户中心、授权菜单、Dashboard |
| IAM | `server/internal/modules/iam` | 用户、角色、菜单、部门、岗位、接口资源 |
| System | `server/internal/modules/system` | 配置、字典、文件、附件、操作日志、登录日志、公告、消息、邮件、通知、健康检查 |
| Setup | `server/internal/modules/setup` | 首次初始化管理员 |
| Module Kit | `server/internal/modules/modulekit` | 受保护路由组和公共中间件装配 |

::: warning 注意目录名
当前代码使用的是 `server/internal/modules/`，不是旧文档里曾出现过的 `server/internal/module/`。新增后端模块时，优先沿 `modules/<group>/<name>` 接入。
:::

## 后端模块内部结构

多数资源型模块使用这套结构：

```text
server/internal/modules/iam/user/
├── api/            # handler、dto、路由细节
├── application/    # 业务逻辑和用例编排
├── domain/         # 领域类型、枚举、常量
├── infra/          # GORM 仓储、数据权限过滤
├── routes.go       # 将模块挂到上级路由组
└── services.go     # 装配 service / repository / handler
```

并不是每个模块都必须机械复制全部目录。比如只读资源、聚合模块或基础设施模块可以更轻，但职责边界仍保持一致：API 层不写复杂业务，业务层不直接暴露 HTTP 细节，仓储层负责数据访问。

## 前端目录

```text
admin/src/
├── api/             # Axios 实例、生成类型、契约检查
├── components/      # 全局组件和管理台 Shell 组件
├── composables/     # 全局组合式函数
├── layouts/         # AdminLayout
├── modules/         # auth / iam / system 三个业务分组
├── router/          # 静态路由、守卫、动态菜单和组件解析
├── stores/          # Pinia 状态
├── styles/          # 主题变量和全局样式
├── types/
└── utils/
```

模块页通常按下面方式组织：

```text
admin/src/modules/system/
├── api/             # 与后端接口对应的请求函数
├── components/      # 当前模块页面使用的展示组件
├── composables/     # useXxxPage 页面状态与副作用
├── pages/           # 路由页面入口
└── types/           # 当前模块类型定义
```

动态菜单组件来自 `admin/src/router/route-components.ts` 的 `import.meta.glob('../modules/**/pages/*View.vue')`。菜单里的 `component` 字段通常写成 `system/FileView`、`iam/UserView` 这类格式。

## 文档站目录

```text
docs/
├── .vitepress/
│   ├── config.mts       # 站点配置、导航、侧边栏
│   └── theme/
├── architecture/        # 架构说明
├── backend/             # 后端说明
├── frontend/            # 前端说明
├── deployment/          # 部署说明
├── getting-started/     # 入门入口
├── reference/           # 稳定参考
└── package.json
```

文档站常用命令：

```bash
make docs-dev
make docs-build
```

执行 `make docs-dev` 后，VitePress 默认使用 `docs/.vitepress/config.mts` 中配置的端口 `15174`。

## 怎么继续读

- 想跑起来项目：看 [快速开始](/getting-started/)
- 想理解请求链路：看 [系统架构概览](/architecture/overview)
- 想查当前模块和事实来源：看 [当前系统地图](/reference/current-system-map)
- 想新增业务模块：看 [模块扩展机制](/architecture/module-extension)
