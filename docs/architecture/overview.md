---
title: 系统架构概览
description: "按当前代码说明 EZ Admin Gin 的单仓库结构、请求链路、模块边界和扩展原则。"
---

# 系统架构概览

EZ Admin Gin 是一个前后端分离的单仓库后台底座：后端负责接口、权限、数据和部署约束，前端负责管理台交互，文档站负责把可复用路径讲清楚。

::: tip 🎯 先抓住主线
新增能力时，先判断它属于 Auth、IAM、System 还是独立业务分组；再分别接后端模块、前端页面、菜单权限、初始化数据和文档说明。
:::

## 一张图看整体

```text
Browser
  ↓
admin/ Vue 管理台
  ↓ /api/v1
server/ Gin API
  ↓
middleware: CORS / Security / RequestID / Metrics / Logger / Recovery
  ↓
authn + actor + permission + operation log
  ↓
modules: auth / iam / system / setup
  ↓
platform: config / database / redis / authz / datascope / migrate
  ↓
PostgreSQL or MySQL + Redis + uploads
```

## 仓库边界

| 区域 | 目录 | 职责 |
| --- | --- | --- |
| 后端服务 | `server/` | HTTP API、认证授权、业务模块、迁移、Swagger |
| 前端管理台 | `admin/` | 登录、后台布局、动态菜单、业务页面、API 类型消费 |
| 文档站 | `docs/` | 入门、架构、后端、前端、部署和参考手册 |
| 部署资产 | `deploy/` | Compose、Nginx、systemd、环境变量模板 |
| 自动化脚本 | `scripts/` | 打包、部署、更新、SQL 生成、Swagger 生成 |

## 后端分层

后端的入口装配集中在 `server/internal/bootstrap/`，业务能力集中在 `server/internal/modules/`，跨模块基础设施集中在 `server/internal/platform/`。

| 层 | 常见目录 | 职责 |
| --- | --- | --- |
| API | `api/` | 参数绑定、响应序列化、HTTP 语义 |
| Application | `application/` | 用例编排、业务规则、跨仓储协调 |
| Domain | `domain/` | 类型、枚举、领域常量 |
| Infra | `infra/` | GORM 查询、数据权限过滤、外部依赖访问 |
| Platform | `platform/` | 配置、数据库、Redis、认证、授权、中间件、迁移 |

这套分层不是为了追求目录数量，而是为了让变化有稳定落点：HTTP 变化不向下污染业务逻辑，数据库查询不向上泄露到页面接口。

## 当前模块边界

| 模块 | 后端目录 | 前端目录 | 当前能力 |
| --- | --- | --- | --- |
| Auth | `server/internal/modules/auth` | `admin/src/modules/auth` | 登录、刷新、退出、当前用户、账户中心、菜单、Dashboard |
| IAM | `server/internal/modules/iam` | `admin/src/modules/iam` | 用户、角色、菜单、部门、岗位、接口资源 |
| System | `server/internal/modules/system` | `admin/src/modules/system` | 配置、字典、文件、附件、日志、公告、消息、邮件、通知、健康检查 |
| Setup | `server/internal/modules/setup` | 无独立页面 | 首次初始化管理员 |

::: warning 避免文档漂移
模块是否存在，以 `server/internal/modules/*`、`admin/src/modules/*` 和路由注册文件为准。不要只根据旧教程里的列表判断当前能力。
:::

## 请求链路

```text
HTTP Request
  → gin.Engine
  → CORS / SecurityHeaders / RequestID / Metrics / Logger / Recovery
  → 静态上传目录、Swagger、健康检查或 /api/v1 路由
  → Auth 中间件校验 JWT
  → LoadActor 装配当前用户、角色、权限上下文
  → OperationLog 记录后台操作
  → Permission 使用 Casbin 做接口授权
  → Handler 绑定参数
  → Application Service 编排业务
  → Infra Repository 查询数据库并注入数据权限
  → httpx 统一响应
```

健康检查和监控有独立入口：

| 路径 | 说明 |
| --- | --- |
| `/healthz` | 存活检查 |
| `/readyz` | 就绪检查 |
| `/health` | 综合健康检查 |
| `/metrics` | Prometheus 指标 |

## 统一响应

API 统一返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

错误码和 HTTP 状态的对应关系集中在 [错误码参考](/reference/error-code-reference)，不要在每个模块文档里重复维护一套。

## 前端运行模型

前端的固定路由负责登录、布局、Dashboard 和账户中心；后台菜单页面由后端授权菜单生成。

```text
登录成功
  → 保存 token
  → 拉取 /api/v1/auth/me 与 /api/v1/auth/menus
  → route-components.ts 解析 component 字段
  → dynamic-menu.ts 生成侧边栏、搜索索引、按钮权限集合
  → 页面 composable 调用模块 API
```

菜单组件键来自 `admin/src/modules/**/pages/*View.vue`。如果后端菜单配置了不存在的组件，前端会回退到 `PlaceholderPage.vue`，避免整站白屏。

## 扩展原则

- 后端新增资源：优先放入 `server/internal/modules/<group>/<name>`。
- 前端新增页面：优先放入 `admin/src/modules/<group>/pages/*View.vue`。
- 权限新增：同时补菜单、按钮权限码、Casbin 策略和初始化数据。
- 数据权限新增：把过滤逻辑放在仓储层或对应的 datascope 辅助中。
- 文档新增：先更新架构或参考入口，再写具体教程，避免同一事实散落多处。

## 下一步

- [后端概览](/backend/overview) — 看后端目录、路由和配置
- [前端概览](/frontend/overview) — 看前端页面组织和动态菜单
- [当前系统地图](/reference/current-system-map) — 查当前模块、页面、路由和事实来源
