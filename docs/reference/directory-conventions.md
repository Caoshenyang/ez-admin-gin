---
title: 目录约定
description: "集中说明 EZ Admin Gin 当前仓库各主目录、server/internal 分层、前端模块和文档站目录的职责边界。"
---

# 目录约定

这页用于回答一个问题：看到一个文件需求时，它应该放到哪里？

::: tip 先看当前代码
当前后端主线是 `server/internal/bootstrap`、`server/internal/platform`、`server/internal/modules` 和 `server/internal/pkg`。旧路径或历史草图不再作为新增代码依据。
:::

## 顶层目录

| 目录 | 职责 |
| --- | --- |
| `server/` | 后端服务、配置、迁移、Swagger、上传目录 |
| `admin/` | Vue 管理台前端 |
| `docs/` | VitePress 文档站 |
| `deploy/` | Docker Compose、Nginx、systemd、环境变量模板 |
| `scripts/` | 打包、部署、更新、Swagger 和 SQL 生成脚本 |
| `brand-assets/` | 品牌资源源文件和导出产物 |
| `.agents/` | 本仓库辅助技能配置，正常业务开发不需要改 |

## `server/internal`

| 目录 | 职责 |
| --- | --- |
| `bootstrap/` | 启动装配、路由总装、Swagger 注册 |
| `modules/` | 认证、IAM、系统、初始化等业务能力 |
| `platform/` | 跨模块基础设施 |
| `pkg/` | 轻量公共工具包 |

最重要的启动入口：

- `server/internal/bootstrap/run.go`
- `server/internal/bootstrap/router.go`

## `platform/`

`platform/` 放跨多个模块复用的底座能力：

| 目录 | 作用 |
| --- | --- |
| `authn/` | JWT、刷新令牌、登出黑名单 |
| `authz/` | Casbin 权限执行器 |
| `config/` | 配置读取、环境变量覆盖、运行时配置 |
| `database/` | 数据库连接和事务辅助 |
| `datascope/` | Actor、数据范围、查询过滤 |
| `logger/` | 日志初始化和 Gin 日志中间件 |
| `middleware/` | 认证、当前用户、权限、操作日志、CORS、安全头等中间件 |
| `migrate/` | 数据库迁移执行 |
| `model/` | GORM 模型 |
| `redis/` | Redis 连接 |

判断标准很简单：多个模块都要用，才考虑放进 `platform/`；只属于某个业务模块，就放进对应模块。

## `modules/`

`modules/` 承接对外形成接口或后台能力的模块。

| 分组 | 说明 |
| --- | --- |
| `modules/auth` | 登录、刷新、退出、当前用户、账户中心、菜单、Dashboard |
| `modules/iam/*` | 用户、角色、菜单、部门、岗位、接口资源 |
| `modules/system/*` | 配置、字典、文件、附件、日志、公告、消息、邮件、通知、健康检查 |
| `modules/setup` | 首次初始化管理员 |
| `modules/modulekit` | 受保护路由组和公共中间件装配 |

成熟资源模块通常包含：

```text
api/            # handler、dto、子路由
application/    # service、ports
domain/         # 类型、常量、枚举
infra/          # repository、数据权限、外部依赖
routes.go       # 模块挂载
services.go     # 依赖装配
```

模块可以按复杂度裁剪，不必为了统一而强行补齐所有目录。

## `pkg/`

`pkg/` 放不携带业务含义的小工具：

| 目录 | 说明 |
| --- | --- |
| `actorx/` | 当前操作者上下文辅助 |
| `errorsx/` | 业务错误类型和错误码 |
| `httpx/` | 统一响应辅助 |
| `paging/` | 分页参数和分页响应 |

如果工具开始依赖某个业务模块，它就不应该继续放在 `pkg/`。

## `admin/src`

| 目录 | 职责 |
| --- | --- |
| `api/` | HTTP 客户端、生成类型、契约检查 |
| `components/` | 全局组件、Shell 组件、后台基础组件 |
| `composables/` | 全局组合式函数 |
| `layouts/` | 管理台布局 |
| `modules/` | `auth`、`iam`、`system` 页面和模块代码 |
| `router/` | 静态路由、守卫、动态菜单、组件解析 |
| `stores/` | Pinia 状态 |
| `styles/` | 主题变量和全局样式 |
| `types/` | 全局类型 |
| `utils/` | 工具函数 |

模块内页面优先按 `api / components / composables / pages / types` 组织。

## `docs/`

| 目录 | 职责 |
| --- | --- |
| `.vitepress/` | VitePress 配置、导航、主题 |
| `getting-started/` | 入门路径 |
| `architecture/` | 架构说明 |
| `backend/` | 后端说明 |
| `frontend/` | 前端说明 |
| `deployment/` | 部署说明 |
| `reference/` | 可反复查阅的稳定参考 |

文档修改后至少执行：

```bash
make docs-build
```

## 放置判断

1. 启动和总装问题：放 `server/internal/bootstrap`。
2. 跨模块基础设施：放 `server/internal/platform`。
3. 业务模块能力：放 `server/internal/modules/<group>/<name>`。
4. 小型无业务工具：放 `server/internal/pkg`。
5. 前端页面：放 `admin/src/modules/<group>/pages`。
6. 文档入口或稳定约定：放 `docs/reference`，教程放到对应专题目录。

## 相关页面

- [项目结构](/getting-started/project-structure)
- [系统架构概览](/architecture/overview)
- [模块规范](./module-conventions)
- [当前系统地图](./current-system-map)
