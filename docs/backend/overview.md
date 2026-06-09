---
title: 后端概览
description: "按当前代码说明后端目录、启动流程、路由分组、中间件栈、配置和验证命令。"
---

# 后端概览

后端位于 `server/`，入口是 `server/main.go`，启动装配由 `server/internal/bootstrap/` 完成。业务模块按 Auth、IAM、System、Setup 分组，公共能力放在 `platform/`。

::: tip 🎯 这页怎么读
想改接口，先看“路由分组”；想改业务，先看“模块结构”；想改配置，先看“配置来源”；改完后用最后的命令验证。
:::

## 目录结构

```text
server/
├── main.go
├── configs/
│   ├── config.yaml             # 默认配置
│   └── rbac_model.conf         # Casbin 模型
├── docs/                       # Swagger / OpenAPI 生成产物
├── internal/
│   ├── bootstrap/              # 启动、路由聚合、Swagger 注册
│   ├── modules/                # auth / iam / system / setup / modulekit
│   ├── platform/               # config / database / redis / authn / authz / middleware 等
│   └── pkg/                    # errorsx / httpx / actorx / paging
└── migrations/
    ├── mysql/
    └── postgres/
```

`server/internal/platform/model/` 当前有 25 个 GORM 模型文件，表结构交付以 `server/migrations/{mysql,postgres}/full_schema_and_seed.sql` 为准。

## 启动流程

```text
main.go
  → bootstrap.MustRun()
  → 加载 configs/config.yaml 与 EZ_* 环境变量
  → 校验生产环境安全配置
  → 连接数据库
  → 执行迁移
  → 连接 Redis
  → 初始化 JWT 管理器与 Casbin
  → NewRouter() 注册中间件、静态上传目录、Swagger、模块路由
  → 启动 HTTP 服务
```

## 路由分组

后端 API 统一挂在 `/api/v1` 下，健康检查和监控是例外。

| 分组 | 路由 | 当前能力 |
| --- | --- | --- |
| Auth | `/api/v1/auth` | 登录、刷新令牌、退出、当前用户、账户资料、改密码、授权菜单、Dashboard |
| Setup | `/api/v1/setup/init` | 首次初始化管理员 |
| IAM | `/api/v1/system/users` | 用户列表、新增、更新、状态、角色分配、删除 |
| IAM | `/api/v1/system/roles` | 角色列表、新增、更新、状态、接口权限、菜单权限、删除 |
| IAM | `/api/v1/system/menus` | 菜单列表、新增、更新、状态、删除 |
| IAM | `/api/v1/system/departments` | 部门列表、新增、更新、状态、删除 |
| IAM | `/api/v1/system/posts` | 岗位列表、新增、更新、状态、删除 |
| IAM | `/api/v1/system/apis` | 接口资源列表 |
| System | `/api/v1/system/configs` | 配置列表、新增、更新、状态、删除、按 key 取值 |
| System | `/api/v1/system/dict-types`、`/dict-items` | 字典类型和字典项管理 |
| System | `/api/v1/system/files`、`/attachments` | 文件上传、文件列表、附件上传与附件状态 |
| System | `/api/v1/system/operation-logs`、`/login-logs` | 操作日志、登录日志 |
| System | `/api/v1/system/notices` | 通知公告管理 |
| System | `/api/v1/system/message-templates`、`/message-reminders` | 消息模板、消息提醒规则 |
| System | `/api/v1/system/mail/*` | 邮箱账号、邮件模板、邮件发送、邮件日志 |
| System | `/api/v1/system/notifications` | 站内通知、未读数、标记已读 |
| System | `/api/v1/system/notifications/ws` | 站内通知 WebSocket |

独立基础路由：

| 路由 | 说明 |
| --- | --- |
| `/healthz` | 存活检查 |
| `/readyz` | 就绪检查 |
| `/health`、`/api/v1/system/health` | 健康检查 |
| `/metrics` | Prometheus 指标 |
| `/uploads/*` | 上传文件公开访问路径，前缀由 `upload.public_path` 控制 |

::: warning 路由以代码为准
如果这里和实际接口有差异，先查看 `server/internal/modules/**/api/routes.go` 与 `server/internal/bootstrap/router.go`，再更新文档。
:::

## 中间件栈

全局中间件在 `bootstrap.NewRouter()` 中注册：

| 中间件 | 职责 |
| --- | --- |
| CORS | 跨域控制，开发环境允许 localhost |
| SecurityHeaders | 生产/开发环境安全响应头 |
| RequestID | 请求追踪 ID |
| Metrics | Prometheus 指标采集 |
| GinLogger / GinRecovery | 请求日志与 panic 恢复 |

后台受保护分组由 `modulekit.NewProtectedSystemGroup()` 装配：

| 中间件 | 职责 |
| --- | --- |
| Auth | 校验访问令牌和登出黑名单 |
| LoadActor | 装配当前用户、角色、权限上下文 |
| OperationLog | 记录后台操作日志 |
| Permission | 使用 Casbin 校验接口权限 |

Auth 模块自己的受保护接口只挂 `Auth` 和 `LoadActor`，用于当前用户信息、账户中心、菜单和 Dashboard。

## 配置来源

配置结构定义在 `server/internal/platform/config/config.go`，默认文件是 `server/configs/config.yaml`，环境变量前缀是 `EZ_`。

| 配置段 | 说明 |
| --- | --- |
| `app` | 应用名称、运行环境 |
| `server` | HTTP 监听地址 |
| `database` | 数据库驱动、连接信息、连接池 |
| `redis` | Redis 连接和连接池 |
| `auth` | JWT 密钥、访问令牌和刷新令牌有效期 |
| `upload` | 上传目录、公开访问路径、大小和扩展名白名单 |
| `log` | 日志级别、格式、轮转 |
| `swagger` | 是否启用 Swagger UI |
| `cors` | 允许跨域来源 |
| `rate_limit` | 登录限流和账号锁定 |

::: warning 生产环境必须改默认密钥
当 `app.env=prod` 时，服务会校验 JWT 密钥、CORS、Swagger 和上传大小。生产环境至少要通过 `EZ_AUTH_JWT_SECRET` 设置安全随机密钥。
:::

## 常用验证命令

在仓库根目录执行：

```bash
make server-vet
make swagger
make docs-build
```

如果改了接口返回结构，还要执行：

```bash
make check-types
```

`make check-types` 会重新生成 Swagger，再检查 `admin/src/api/generated.ts` 是否同步。
