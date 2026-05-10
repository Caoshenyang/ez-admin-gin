---
title: 后端概览
description: Go + Gin 后端的架构分层、模块结构、请求处理流程
---

# 后端概览

## 目录结构

```
server/
├── main.go                     入口文件
├── configs/                    配置文件
│   ├── config.yaml             主配置
│   └── rbac_model.conf         Casbin RBAC 模型
├── internal/
│   ├── bootstrap/              启动引导
│   │   ├── run.go              初始化流程
│   │   ├── router.go           路由注册
│   │   └── swagger.go          Swagger 文档
│   ├── modules/                业务模块
│   │   ├── auth/               认证模块
│   │   ├── iam/                身份与访问管理
│   │   ├── system/             系统模块
│   │   ├── setup/              系统初始化
│   │   └── modulekit/          模块工具包
│   ├── platform/               平台层（跨模块基础设施）
│   │   ├── authn/              JWT 认证
│   │   ├── authz/              Casbin 授权
│   │   ├── config/             配置管理
│   │   ├── database/           数据库连接
│   │   ├── datascope/          数据权限
│   │   ├── logger/             日志
│   │   ├── middleware/         HTTP 中间件
│   │   ├── migrate/            迁移管理
│   │   ├── model/              GORM 模型（18 张表）
│   │   └── redis/              Redis 连接
│   └── pkg/                    公共工具包
│       ├── errorsx/            错误处理
│       ├── httpx/              HTTP 工具
│       ├── actorx/             Actor 上下文
│       └── paging/             分页
├── migrations/                 数据库迁移
│   ├── mysql/                  MySQL 方言
│   └── postgres/               PostgreSQL 方言
├── uploads/                    文件上传目录
└── logs/                       日志文件
```

## 四层架构

每个业务模块遵循 handler → service → repository → domain 四层结构：

```
api/handlers.go    接收 HTTP 请求，绑定参数，调用 service，返回响应
application/       业务逻辑，权限校验，跨模块协调
infra/             GORM 数据访问，数据权限注入
domain/types.go    领域类型，常量，枚举
```

## API 路由结构

```
/api/v1/
├── auth/                          认证模块
│   ├── POST /login                登录
│   ├── GET  /me                   当前用户信息
│   ├── GET  /account              账户详情
│   ├── POST /account/profile      更新个人信息
│   ├── POST /account/password     修改密码
│   ├── GET  /menus                当前用户菜单树
│   └── GET  /dashboard            Dashboard 数据
├── iam/                           身份与访问管理
│   ├── /users                     用户管理 CRUD
│   ├── /roles                     角色管理 CRUD
│   ├── /menus                     菜单管理 CRUD
│   └── /departments               部门管理 CRUD
├── system/                        系统模块
│   ├── /configs                   系统配置
│   ├── /dict-types                字典类型
│   ├── /dict-items                字典项
│   ├── /files                     文件上传
│   ├── /attachments               附件管理
│   ├── /operation-logs            操作日志
│   ├── /login-logs                登录日志
│   └── /notices                   通知公告
└── setup/                         系统初始化
    └── POST /init                 初始化管理员账号
```

## 中间件栈

| 中间件 | 位置 | 职责 |
|--------|------|------|
| CORS | 全局 | 跨域资源共享 |
| RequestID | 全局 | 请求追踪 ID |
| Logger | 全局 | 请求日志 |
| Recovery | 全局 | Panic 恢复 |
| Auth | 认证路由组 | JWT 验证 |
| LoadActor | 认证路由组 | 加载用户上下文 |
| Permission | 需鉴权路由 | Casbin 策略匹配 |
| RateLimit | 登录接口 | 登录限流 |
| OperationLog | 需记录路由 | 操作日志 |

## 启动流程

```
main.go
  → bootstrap.MustRun()
    → 加载配置（config.yaml + 环境变量）
    → 连接数据库
    → 执行迁移
    → 连接 Redis
    → 初始化 Casbin
    → 注册路由
    → 启动 HTTP 服务
```

## 配置管理

配置通过 `configs/config.yaml` 管理，支持环境变量覆盖（前缀 `EZ_`）。

| 配置段 | 关键项 | 说明 |
|--------|--------|------|
| app | name, env | 应用名称和环境 |
| server | addr | 监听地址（默认 `:8080`） |
| database | driver, host, port, user, password, name | 数据库连接 |
| redis | host, port, password | Redis 连接 |
| log | level, format, filename | 日志配置 |
| auth | jwt_secret, access_token_ttl | JWT 配置 |
| cors | allowed_origins | CORS 白名单 |
| rate_limit | login_max_requests, login_window_sec | 登录限流 |
| upload | dir, max_size_mb, allowed_exts | 文件上传 |

::: warning
生产环境必须通过环境变量 `EZ_AUTH_JWT_SECRET` 覆盖 JWT 密钥。如果 `env=prod` 且密钥包含 `change-me`，服务将拒绝启动。
:::
