---
title: 教程大纲
description: "EZ Admin Gin 企业级完整版 0-1 教程的大纲，用于确定最终主线章节、小节边界和验证顺序。"
---

# 教程大纲

这份大纲只负责确定企业级完整版教程的主线和小节边界。每一节都围绕最终形态来写，不再把“简化版”作为教程默认路线。

::: tip 当前策略
教程主线固定为 9 章，目标读者是 Java 转 Go 工程师。每一章都要兼顾“能做成什么”“为什么这样设计”“Go vs Java 怎么理解”“执行后该看到什么”。
:::

## 第 1 章：项目定位与仓库初始化

- [章节导读](./chapter-1/)
- [项目仓库初始化](./chapter-1/project-repository-init)
- [Go 后端项目初始化](./chapter-1/backend-init)
- [Vue 管理台项目初始化](./chapter-1/admin-init)
- [VitePress 文档项目初始化](./chapter-1/docs-init)
- [Docker Compose 基础环境](./chapter-1/docker-compose-env)

## 第 2 章：平台基础设施

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
- [RBAC 角色权限模型](./chapter-3/rbac-model)
- [接口级权限控制](./chapter-3/casbin-permission)
- [角色菜单权限](./chapter-3/menu-permission)

## 第 5 章：组织体系与数据权限

- [章节导读](./chapter-5/)
- [组织模型设计](./chapter-5/organization-model-design)
- [角色数据范围与查询作用域](./chapter-5/role-data-scope-and-query-scopes)
- [部门树与部门管理](./chapter-5/department-tree-and-management)
- [岗位管理与用户归属](./chapter-5/post-management-and-user-affiliation)
- 角色数据范围模型
- Actor 上下文
- GORM Scope 数据过滤
- 超级管理员与多角色并集规则

## 第 6 章：核心系统模块

- [章节导读](./chapter-6/)
- 用户管理
- 角色管理
- 菜单管理
- 系统配置
- 文件中心
- 操作日志
- 登录日志
- 公告管理

## 第 7 章：前端企业级管理台

- [章节导读](./chapter-7/)
- [Vue 3 管理台初始化](./chapter-5/vue-project-init)
- [登录页](./chapter-5/login-page)
- [后台布局](./chapter-5/admin-layout)
- [动态菜单](./chapter-5/dynamic-menu)
- [用户管理页面](./chapter-5/user-pages)
- [角色与菜单页面](./chapter-5/role-menu-pages)
- [配置与文件页面](./chapter-5/config-file-pages)
- [日志页面](./chapter-5/log-pages)

## 第 8 章：模块化接入规范

- [章节导读](./chapter-8/)
- [模块固定结构](./chapter-6/module-structure)
- [后端模块接入流程](./chapter-6/backend-module-flow)
- [权限、菜单与迁移接入](./chapter-6/permission-menu-migration)
- [前端页面接入流程](./chapter-6/frontend-page-flow)
- [示例业务模块](./chapter-6/sample-module)

## 第 9 章：部署、升级与复用

- [章节导读](./chapter-9/)
- [后端与前端 Dockerfile](./chapter-7/dockerfile)
- [Docker Compose 编排](./chapter-7/docker-compose)
- [Nginx 配置](./chapter-7/nginx-config)
- [环境变量与初始化数据](./chapter-7/env-and-init-data)
- [部署验证与复用说明](./chapter-7/deployment-and-reuse)
