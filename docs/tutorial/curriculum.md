---
title: 教程大纲
description: "EZ Admin Gin 企业级完整版 0-1 教程的大纲，用于确定最终主线章节、小节边界和验证顺序。"
---

# 教程大纲

这份大纲只负责确定企业级完整版教程的主线和小节边界。每一节都围绕最终形态来写，不再把“简化版”作为教程默认路线。

::: tip 当前策略
教程主线固定为 9 章，目标读者是 Java 转 Go 工程师。每一章都要兼顾“能做成什么”“为什么这样设计”“Go vs Java 怎么理解”“执行后该看到什么”。
:::

## 第 1 章：项目初始化

- [章节导读](./chapter-1/)
- [项目仓库初始化](./chapter-1/project-repository-init)
- [Go 后端项目初始化](./chapter-1/backend-init)
- [Vue 管理台项目初始化](./chapter-1/admin-init)
- [VitePress 文档项目初始化](./chapter-1/docs-init)
- [Docker Compose 基础环境](./chapter-1/docker-compose-env)

## 第 2 章：后端基础设施

- [章节导读](./chapter-2/)
- [配置管理](./chapter-2/config-management)
- [日志系统](./chapter-2/logging-system)
- [数据库连接](./chapter-2/database-connection)
- [Redis 连接](./chapter-2/redis-connection)
- [统一响应与错误处理](./chapter-2/response-and-errors)
- [路由分组与健康检查](./chapter-2/routing-and-health)

## 第 3 章：认证与登录态

- [章节导读](./chapter-3/)
- [用户模型与登录](./chapter-3/user-model-and-login)
- [Token 签发与解析](./chapter-3/jwt-auth)
- [登录校验中间件](./chapter-3/auth-middleware)

## 第 4 章：接口权限体系

- [章节导读](./chapter-4/)
- [RBAC 角色权限模型](./chapter-4/rbac-model)
- [接口级权限控制](./chapter-4/casbin-permission)
- [角色菜单权限](./chapter-4/menu-permission)

## 第 5 章：组织体系与数据权限

- [章节导读](./chapter-5/)
- [组织模型设计](./chapter-5/organization-model-design)
- [角色数据范围与查询作用域](./chapter-5/role-data-scope-and-query-scopes)
- [Actor 上下文与多角色并集](./chapter-5/actor-and-grant-merge)
- [资源级数据权限接入模式](./chapter-5/module-datascope-patterns)
- [共享数据权限接入规范](./chapter-5/shared-datascope-integration-conventions)
- [datascope.go 与 Repository 边界](./chapter-5/datascope-and-repository-boundary)
- [一次完整请求的权限过滤走读](./chapter-5/request-flow-walkthrough)
- [数据权限落地检查清单](./chapter-5/data-scope-implementation-checklist)
- [部门树与部门管理](./chapter-5/department-tree-and-management)
- [岗位管理与用户归属](./chapter-5/post-management-and-user-affiliation)
- [真实业务模块的数据权限边界](./chapter-5/business-module-datascope-boundaries)
- [岗位资源的数据权限收紧时机](./chapter-5/post-datascope-tightening)

## 第 6 章：核心系统模块

- [章节导读](./chapter-6/)
- [模块固定结构](./chapter-6/module-structure)
- [后端模块接入流程](./chapter-6/backend-module-flow)
- [系统模块示例](./chapter-6/sample-module)
- [数据字典模块落地](./chapter-6/dict-module)
- [附件中心落地](./chapter-6/attachment-center-module)
- [账户中心落地](./chapter-6/account-center-module)
- [权限、菜单与迁移接入](./chapter-6/permission-menu-migration)
- [前端页面接入流程](./chapter-6/frontend-page-flow)
- [模块接入验收清单](./chapter-6/module-integration-checklist)

## 第 7 章：前端企业级管理台

- [章节导读](./chapter-7/)
- [前端运行时结构](./chapter-7/frontend-runtime-structure)
- [管理台工程起步结构](./chapter-7/admin-project-bootstrap)
- [登录态与会话流转](./chapter-7/login-and-session-flow)
- [登录页实现细节](./chapter-7/login-page-implementation)
- [后台壳子、动态菜单与按钮权限](./chapter-7/admin-shell-and-dynamic-menu)
- [后台布局与工作标签](./chapter-7/admin-layout-and-worktabs)
- [动态菜单注册与按钮权限](./chapter-7/dynamic-route-registration)
- [系统模块页面模式](./chapter-7/system-module-pages)
- [用户管理页实现要点](./chapter-7/user-management-page-detail)
- [角色与菜单页实现要点](./chapter-7/role-and-menu-page-detail)
- [配置与文件页实现要点](./chapter-7/config-and-file-page-detail)
- [日志查询页实现要点](./chapter-7/audit-log-pages)

## 第 8 章：部署、升级与复用

- [章节导读](./chapter-8/)
- [环境变量与初始化数据](./chapter-8/env-and-init-data)
- [部署验证与复用说明](./chapter-8/deployment-and-reuse)
- [Compose 与服务运行结构](./chapter-8/compose-and-service-layout)
- [Nginx 与 HTTPS 入口层](./chapter-8/nginx-and-https)
- [部署变体说明](./chapter-8/deployment-variants)
- [更新与回滚策略](./chapter-8/update-and-rollback)
- [回滚分级策略](./chapter-8/rollback-strategy-levels)
- [部署排障 FAQ](./chapter-8/deployment-troubleshooting-faq)
- [长期运维 FAQ](./chapter-8/operations-maintenance-faq)
- [新项目复用清单](./chapter-8/project-reuse-checklist)
