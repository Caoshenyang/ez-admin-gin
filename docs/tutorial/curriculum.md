---
title: 教程大纲
description: "EZ Admin Gin 从零搭建教程的小节划分，用于确定后续一章一章完善和验证的顺序。"
---

# 教程大纲

这份大纲只负责确定教程主线和小节边界。每一节后续完善时，都要补齐前置条件、操作步骤、期望结果和失败排查。

::: tip 当前策略
先把 `tutorial/` 从零搭建主线跑通。等全部完成并经过验证后，再统一整理“开始这里”、参考手册和路线图。
:::

## 第一阶段：项目初始化

- [章节导读](./chapter-1/)
- [仓库创建与 Git 初始化](./chapter-1/repository-and-sequence)
- [单仓库目录结构](./chapter-1/directory-structure)
- [Go 后端项目初始化](./chapter-1/backend-init)
- [Vue 前端项目初始化](./chapter-1/frontend-init)
- [VitePress 文档项目初始化](./chapter-1/docs-init)
- [Docker Compose 基础环境](./chapter-1/docker-compose-env)

## 第二阶段：后端基础设施

- [章节导读](./chapter-2/)
- [配置管理](./chapter-2/config-management)
- [日志系统](./chapter-2/logging-system)
- [数据库连接](./chapter-2/database-connection)
- [Redis 连接](./chapter-2/redis-connection)
- [统一响应与错误处理](./chapter-2/response-and-errors)
- [路由分组与健康检查](./chapter-2/routing-and-health)

## 第三阶段：认证与权限

- [章节导读](./chapter-3/)
- [用户模型与登录接口](./chapter-3/user-model-and-login)
- [JWT 认证](./chapter-3/jwt-auth)
- [认证中间件](./chapter-3/auth-middleware)
- [角色与权限模型](./chapter-3/rbac-model)
- [Casbin 权限控制](./chapter-3/casbin-permission)
- [菜单权限设计](./chapter-3/menu-permission)

## 第四阶段：通用系统模块

- [章节导读](./chapter-4/)
- [用户管理](./chapter-4/user-management)
- [角色管理](./chapter-4/role-management)
- [菜单管理](./chapter-4/menu-management)
- [系统配置](./chapter-4/system-config)
- [文件上传](./chapter-4/file-upload)
- [操作日志](./chapter-4/operation-logs)
- [登录日志](./chapter-4/login-logs)

## 第五阶段：前端管理台

- [章节导读](./chapter-5/)
- [Vue 3 管理台初始化](./chapter-5/vue-project-init)
- [登录页](./chapter-5/login-page)
- [后台布局](./chapter-5/admin-layout)
- [动态菜单](./chapter-5/dynamic-menu)
- [用户管理页面](./chapter-5/user-pages)
- [角色与菜单页面](./chapter-5/role-menu-pages)
- [配置与文件页面](./chapter-5/config-file-pages)

## 第六阶段：业务模块接入规范

- [章节导读](./chapter-6/)
- [模块固定结构](./chapter-6/module-structure)
- [后端模块接入流程](./chapter-6/backend-module-flow)
- [权限、菜单与迁移接入](./chapter-6/permission-menu-migration)
- [前端页面接入流程](./chapter-6/frontend-page-flow)
- [示例业务模块](./chapter-6/sample-module)

## 第七阶段：部署与复用

- [章节导读](./chapter-7/)
- [后端与前端 Dockerfile](./chapter-7/dockerfile)
- [Docker Compose 编排](./chapter-7/docker-compose)
- [Nginx 配置](./chapter-7/nginx-config)
- [环境变量与初始化数据](./chapter-7/env-and-init-data)
- [部署验证与复用说明](./chapter-7/deployment-and-reuse)
