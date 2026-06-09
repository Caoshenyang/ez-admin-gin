---
title: 系统架构概览
description: EZ Admin 整体架构设计、模块划分与技术选型
---

# 系统架构概览

## 整体架构

EZ Admin 采用前后端分离的单仓库结构：

```
ez-admin-gin/
├── server/           Go + Gin 后端
│   ├── configs/      配置文件
│   ├── internal/     业务代码（模块化分层）
│   └── migrations/   数据库迁移（MySQL + PostgreSQL）
├── admin/            Vue 3 前端
│   └── src/
│       ├── modules/  业务模块（auth / iam / system）
│       ├── layouts/  布局
│       └── router/   路由（含动态菜单注册）
├── docs/             VitePress 文档站
├── deploy/           Docker Compose、Nginx、部署配置
└── scripts/          部署与打包脚本
```

## 后端分层

后端采用四层架构，每层职责明确：

| 层 | 目录 | 职责 |
|----|------|------|
| Handler | `api/` | HTTP 请求处理、参数绑定、响应序列化 |
| Service | `application/` | 业务逻辑编排、权限校验、跨模块协调 |
| Repository | `infra/` | 数据访问、GORM 查询、数据权限注入 |
| Domain | `domain/` | 领域模型、常量、业务规则 |

平台层（`platform/`）提供跨模块基础设施：认证、授权、配置、数据库、日志、中间件、数据权限、迁移。

## 前端分层

前端采用三层模块化结构：

| 层 | 目录 | 职责 |
|----|------|------|
| Pages | `pages/` | 编排层，拼装 composable + component |
| Composables | `composables/` | 状态管理 + 副作用逻辑 |
| Components | `components/` | 展示组件，通过 props/events 通信 |

全局共享层：`api/`（HTTP 客户端）、`router/`（路由+动态菜单）、`stores/`（Pinia）、`composables/`（组合函数）、`utils/`（工具）。

## 请求链路

```
HTTP Request
  → gin.Engine
  → CORS / RequestID / Logger / Recovery（全局中间件）
  → /api/v1 路由组
  → Auth 中间件（JWT 验证）
  → LoadActor 中间件（加载用户上下文 + 角色 + 权限）
  → Permission 中间件（Casbin 策略匹配）
  → OperationLog 中间件（操作日志记录）
  → Handler（参数绑定）
  → Service（业务逻辑）
  → Repository（数据访问 + 数据权限过滤）
  → 统一响应格式返回
```

## 统一响应格式

所有 API 返回统一的 JSON 结构：

```json
{
  "code": 0,
  "message": "ok",
  "data": { ... }
}
```

错误码体系：`0` 成功、`40000` 请求错误、`40100` 未认证、`40300` 无权限、`40400` 未找到、`50000` 内部错误。

## 模块划分

### 后端模块

| 模块 | 路径 | 职责 |
|------|------|------|
| auth | `modules/auth/` | 登录认证、当前用户、菜单、Dashboard、账户中心 |
| iam | `modules/iam/` | 用户、角色、菜单、部门管理 |
| system | `modules/system/` | 配置、字典、文件、附件、日志、公告 |
| setup | `modules/setup/` | 系统初始化 |

### 前端模块

| 模块 | 路径 | 职责 |
|------|------|------|
| auth | `modules/auth/` | 登录页、Dashboard、账户中心 |
| iam | `modules/iam/` | 用户、角色、菜单、部门、岗位管理页 |
| system | `modules/system/` | 配置、字典、文件、日志、公告管理页 |

## 技术选型

| 层 | 技术 | 选型理由 |
|----|------|---------|
| 后端框架 | Gin | 成熟高性能，中间件生态完善 |
| ORM | GORM | Go 生态主流，支持多数据库 |
| 权限引擎 | Casbin | 灵活的 RBAC/ABAC 策略引擎 |
| 缓存 | Redis | 会话、限流、缓存 |
| 前端框架 | Vue 3 | Composition API + TypeScript |
| UI 组件库 | Naive UI | 后台场景组件齐全，主题可定制 |
| 状态管理 | Pinia | Vue 3 官方推荐 |
| CSS | Tailwind CSS 4 | 实用优先，与 Naive UI 互补 |
| 构建工具 | Vite | 快速开发体验 |
| 文档 | VitePress | Vue 生态文档站 |

## 设计原则

- **模块自治**：每个模块有独立的 api/service/repository/domain 四层
- **平台层共享**：认证、授权、数据权限等横切关注点集中在 platform 层
- **配置外置**：YAML 配置 + 环境变量覆盖，不硬编码
- **数据权限在 Repository 层注入**：Service 不关心数据过滤逻辑
- **前端不直接调用 API**：组件通过 composable 获取数据，页面只做编排
