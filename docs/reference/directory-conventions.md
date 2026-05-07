---
title: 目录约定
description: "集中说明 EZ Admin Gin 当前仓库各主目录职责、server/internal 分层边界，以及哪些目录属于历史兼容区。"
---

# 目录约定

这页的目标很简单：

> 让你看到一个目录时，马上知道它在当前最终结构里应该承担什么职责。

## 仓库顶层目录

当前仓库顶层最重要的目录如下：

| 目录 | 职责 |
| --- | --- |
| `admin/` | Vue 3 + Vite 管理台前端 |
| `server/` | Go 后端服务与迁移、配置、上传目录 |
| `docs/` | VitePress 教程、指南和参考手册 |
| `deploy/` | Compose、Nginx、systemd 等部署工件 |
| `scripts/` | 打包、初始化服务器、更新服务器等脚本 |
| `.github/` | CI、自动化工作流 |

## `admin/` 的职责边界

当前前端工程主要看这些子目录：

| 目录 | 职责 |
| --- | --- |
| `admin/src/api` | 接口封装 |
| `admin/src/components` | 复用组件 |
| `admin/src/layouts` | 后台壳子布局 |
| `admin/src/pages` | 页面级视图 |
| `admin/src/router` | 路由与动态菜单注册 |
| `admin/src/styles` | 全局样式 |
| `admin/src/types` | TS 类型定义 |
| `admin/src/utils` | 鉴权、请求、菜单等工具函数 |

当前页面大致分为：

- `pages/auth`
- `pages/dashboard`
- `pages/system`

这说明现在的管理台不是“组件先行”的结构，而是：

- 页面目录承接业务
- 路由和菜单负责拼装运行时

## `server/` 的职责边界

`server/` 下面除了代码，还有运行时资产：

| 目录 | 职责 |
| --- | --- |
| `server/cmd` | 程序入口 |
| `server/configs` | 本地配置样例 |
| `server/internal` | 核心后端实现 |
| `server/migrations` | 数据库迁移文件 |
| `server/uploads` | 上传文件落盘目录 |
| `server/logs` | 本地日志输出目录 |

::: tip 先区分“源码”和“运行时产物”
`internal`、`cmd`、`migrations` 是源码结构。

`logs`、`uploads` 更接近运行时目录，不应该拿来放业务代码。
:::

## `server/internal` 的主骨架

当前后端最终结构围绕下面几层展开：

| 目录 | 职责 |
| --- | --- |
| `bootstrap` | 启动装配与路由总装 |
| `platform` | 平台级基础能力，如认证、权限、数据权限、配置、数据库、日志、迁移、Redis |
| `module` | 业务模块与系统模块 |
| `middleware` | Gin 中间件 |
| `model` | GORM 模型 |
| `response` | 统一响应输出 |
| `apperror` | 统一业务错误 |

最值得记住的入口是：

- `server/internal/bootstrap/run.go`
- `server/internal/bootstrap/router.go`

也就是当前真实主线已经是：

```text
main.go (embed)
  ↓
bootstrap/run.go
  ↓
bootstrap/router.go
  ↓
module/*
```

## `platform/` 放什么，不放什么

当前 `platform/` 负责“跨模块都会复用的底座能力”，例如：

| 目录 | 作用 |
| --- | --- |
| `platform/authn` | Token/JWT 管理 |
| `platform/authz` | Casbin 权限执行器 |
| `platform/config` | 配置读取 |
| `platform/database` | 数据库连接 |
| `platform/datascope` | Actor、Grant、数据权限合并与查询作用域 |
| `platform/logger` | 日志初始化 |
| `platform/migrate` | 迁移执行 |
| `platform/redis` | Redis 连接 |

更稳的判断标准是：

- 如果一个能力会被多个模块共享，优先考虑 `platform`
- 如果它只属于某个业务模块，不要塞进 `platform`

## `module/` 放什么

`module/` 承接的是“对外能形成一组接口或一段业务能力”的模块。

当前主要分成三组：

| 分组 | 说明 |
| --- | --- |
| `module/auth` | 登录、当前用户、菜单、工作台 |
| `module/iam/*` | 用户、角色、部门、岗位、菜单这类身份与授权模型 |
| `module/system/*` | 配置、文件、公告、日志等系统能力 |

还有一个初始化入口：

- `module/setup`

## 一个模块目录里通常长什么样

当前成熟模块通常会包含这些文件：

| 文件 | 职责 |
| --- | --- |
| `routes.go` | 路由注册与依赖装配 |
| `handler.go` | HTTP 入参与响应 |
| `service.go` | 业务规则、事务边界 |
| `repository.go` | 查询与持久化 |
| `dto.go` | 请求、响应、字段归一化 |
| `entity.go` | 当前模块对 `model` 的类型别名或收口 |
| `policy.go` | 接口权限码常量 |
| `datascope.go` | 模块级数据权限接法 |

但这不是机械强制的“八件套”。

更准确地说：

- 需要数据权限的模块，再补 `datascope.go`
- 需要按钮 / 接口权限稳定命名的模块，再补 `policy.go`
- 极轻模块也可能暂时没有独立 `repository.go`

## `middleware/` 的位置为什么独立

当前中间件不是塞在某个模块里，而是独立为一层，因为它们服务的是整条请求链：

| 文件 | 作用 |
| --- | --- |
| `middleware/auth.go` | 登录态校验 |
| `middleware/actor.go` | 装载当前登录人的组织与数据权限上下文 |
| `middleware/permission.go` | Casbin 接口权限校验 |
| `middleware/operation_log.go` | 记录操作日志 |

这些中间件主要在：

- `module/auth/routes.go`
- `module/system/routes.go`

里被串起来。

## `model/` 和 `entity.go` 的关系

当前数据库模型集中放在：

- `server/internal/model`

例如：

- `user.go`
- `role.go`
- `department.go`
- `menu.go`

模块内的 `entity.go` 更像是：

- 对当前模块所依赖模型的一次本地收口
- 避免在 Handler / Service / Repository 里到处直接引用一大串 `model.Xxx`

所以二者分工是：

| 位置 | 作用 |
| --- | --- |
| `model/` | 全局数据库模型定义 |
| `module/*/entity.go` | 当前模块使用的实体别名与局部语义收口 |

## 当前哪些目录属于历史兼容区

仓库里仍有一些旧目录与当前主线结构并存：

| 目录 | 当前状态 |
| --- | --- |
| `server/internal/config` / `database` / `logger` / `redis` / `token` / `permission` | 与 `platform/*` 对应的旧落点仍在仓库中 |

这类目录的理解方式应该是：

- 当前仓库里仍存在
- 但主线结构已经以 `bootstrap / platform / module` 为准

::: warning 不要再把新能力继续扩到历史目录
如果你要补新模块、新平台能力或新中间件，优先沿当前主线落到：

- `bootstrap`
- `platform`
- `module`
- `app`
- `middleware`

而不是继续扩大旧式扁平目录。
:::

## 新增代码时最实用的放置判断

如果你正在犹豫一个文件该放哪，可以按这个顺序判断：

1. 这是启动装配问题吗？是就放 `bootstrap`
2. 这是跨模块复用的平台能力吗？是就放 `platform`
3. 这是某个模块自己的业务能力吗？是就放 `module/<group>/<name>`
4. 这是请求链统一拦截吗？是就放 `middleware`
5. 这是统一错误或响应协议吗？是就放 `apperror` / `response`

## 相关教程与参考页

- [第 6 章：核心系统模块](../tutorial/chapter-6/)
- [第 7 章：前端企业级管理台](../tutorial/chapter-7/)
- [第 6 章：核心系统模块](../tutorial/chapter-6/)
- [模块规范](./module-conventions)
