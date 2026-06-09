---
title: 当前系统地图
description: "按当前代码汇总 EZ Admin Gin 的模块、页面、路由、配置和文档事实来源，降低文档与实现漂移的风险。"
---

# 当前系统地图

这页用于快速校准“当前项目到底有什么”。它不替代源码；当页面内容和代码冲突时，优先看这里列出的事实来源。

::: tip 🎯 使用场景
改文档、扩模块、查接口、排查菜单组件不匹配时，先用这页确认当前系统边界，再进入具体教程或参考页。
:::

## 事实来源

| 你要确认什么 | 优先查看 |
| --- | --- |
| 后端模块是否存在 | `server/internal/modules/` |
| 后端路由是否存在 | `server/internal/modules/**/api/routes.go`、`server/internal/bootstrap/router.go` |
| 受保护中间件挂载方式 | `server/internal/modules/modulekit/` |
| 配置项和环境变量 | `server/internal/platform/config/config.go`、`server/configs/config.yaml` |
| 数据表和初始化数据 | `server/migrations/{mysql,postgres}/full_schema_and_seed.sql` |
| GORM 模型 | `server/internal/platform/model/` |
| Swagger / OpenAPI | `server/docs/swagger.json`、`server/docs/swagger.yaml` |
| 前端页面是否存在 | `admin/src/modules/**/pages/*View.vue` |
| 动态菜单组件解析 | `admin/src/router/route-components.ts` |
| 菜单、搜索和按钮权限 | `admin/src/router/dynamic-menu.ts`、`admin/src/composables/usePermission.ts` |
| 前端 API 类型 | `admin/src/api/generated.ts`、`admin/src/api/contract-check.ts` |
| 文档导航 | `docs/.vitepress/config.mts` |

## 后端模块地图

| 分组 | 当前目录 | 当前能力 |
| --- | --- | --- |
| Auth | `server/internal/modules/auth` | 登录、刷新令牌、退出、当前用户、账户中心、授权菜单、Dashboard |
| IAM / User | `server/internal/modules/iam/user` | 用户列表、新增、更新、状态、角色分配、删除 |
| IAM / Role | `server/internal/modules/iam/role` | 角色列表、新增、更新、状态、接口权限、菜单权限、删除 |
| IAM / Menu | `server/internal/modules/iam/menu` | 菜单列表、新增、更新、状态、删除 |
| IAM / Department | `server/internal/modules/iam/department` | 部门列表、新增、更新、状态、删除 |
| IAM / Post | `server/internal/modules/iam/post` | 岗位列表、新增、更新、状态、删除 |
| IAM / API Resource | `server/internal/modules/iam/apiresource` | 接口资源列表 |
| System / Config | `server/internal/modules/system/config` | 配置列表、新增、更新、状态、删除、按 key 取值 |
| System / Dict | `server/internal/modules/system/dict` | 字典类型和字典项管理 |
| System / File | `server/internal/modules/system/file` | 文件列表、文件上传 |
| System / Attachment | `server/internal/modules/system/attachment` | 附件列表、上传、更新、状态 |
| System / Operation Log | `server/internal/modules/system/operationlog` | 操作日志列表 |
| System / Login Log | `server/internal/modules/system/loginlog` | 登录日志列表 |
| System / Notice | `server/internal/modules/system/notice` | 公告列表、新增、更新、状态、删除 |
| System / Message | `server/internal/modules/system/message` | 消息模板、消息提醒规则 |
| System / Mail | `server/internal/modules/system/mail` | 邮箱账号、邮件模板、发送邮件、邮件日志 |
| System / Notification | `server/internal/modules/system/notification` | 站内通知、未读数、标记已读、WebSocket |
| System / Health | `server/internal/modules/system/health_handler.go` | 存活、就绪和综合健康检查 |
| Setup | `server/internal/modules/setup` | 首次初始化管理员 |

## 前端页面地图

| 分组 | 当前页面 |
| --- | --- |
| Auth | `LoginPage.vue`、`DashboardHome.vue`、`AccountCenterPage.vue` |
| IAM | `UserView.vue`、`RoleView.vue`、`MenuView.vue`、`DepartmentView.vue`、`PostView.vue` |
| System | `HealthView.vue`、`ConfigView.vue`、`DictView.vue`、`FileView.vue`、`AttachmentView.vue`、`OperationLogView.vue`、`LoginLogView.vue`、`NoticeView.vue`、`MessageView.vue`、`MailView.vue`、`PlaceholderPage.vue` |

动态路由只会自动收集 `admin/src/modules/**/pages/*View.vue`。`DashboardHome.vue` 和 `AccountCenterPage.vue` 属于静态路由，不走菜单组件扫描。

## 路由前缀

| 类型 | 前缀 |
| --- | --- |
| 认证接口 | `/api/v1/auth` |
| 初始化接口 | `/api/v1/setup` |
| 后台受保护接口 | `/api/v1/system` |
| 健康检查 | `/healthz`、`/readyz`、`/health`、`/api/v1/system/health` |
| 监控指标 | `/metrics` |
| 上传文件 | 默认 `/uploads`，以 `upload.public_path` 为准 |

## 校验命令

在仓库根目录执行：

```bash
make server-vet
make admin-check
make check-types
make docs-build
```

::: warning 文档修改后的最低检查
只改文档时至少执行 `make docs-build`。如果文档同步了接口、模块或类型相关内容，建议同时执行 `make check-types`，确认 Swagger 和前端生成类型没有漂移。
:::

## 更新文档时的顺序

1. 先看代码事实来源，不先改措辞。
2. 再更新总览页或参考页，保证入口能说明当前系统边界。
3. 最后更新教程页或部署页，避免同一个事实被散落维护。
4. 修改完成后运行 `make docs-build`。
