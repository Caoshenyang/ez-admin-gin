---
title: 第 7 章：前端企业级管理台
description: "围绕当前真实前端结构，串起登录态、后台布局、动态菜单、按钮权限和系统模块页面对接。"
---

# 第 7 章：前端企业级管理台

后端的基础能力、权限体系、组织体系和系统模块都已经具备后，这一章开始把它们接成一个真正可操作的企业级管理台。

::: tip 本章怎么读
这一章建议按“运行时结构 → 登录态 → 后台壳子 → 动态菜单 → 系统页面模式 → 详细案例”的顺序读。这样更容易把前端代码的启动链路、权限链路和页面落点一次串起来。
:::

## 本章会解决什么

- 登录成功后，前端如何保存登录态并拉起初始化链路
- 后台布局如何承接动态菜单、工作台和页面容器
- `/api/v1/auth/menus` 返回的菜单树如何变成前端路由和侧边栏
- 按钮权限码如何从菜单树里提取并落到页面操作显隐
- 系统模块页面如何围绕当前 `admin/src/api/*`、`types/*`、`pages/system/*` 继续扩展

::: info 边界提醒
第 7 章只负责前端运行时和页面承接，不再继续承载部署说明。部署、环境变量、升级和回滚已经统一放到 [第 8 章](../chapter-8/)。
:::

## 当前真实前端结构

现在仓库里的前端主线主要落在：

```text
admin/src/
├─ api/
├─ types/
├─ pages/
│  ├─ auth/
│  ├─ dashboard/
│  └─ system/
├─ layouts/
├─ router/
└─ utils/
```

第 7 章真正要讲清楚的，就是这几层是怎么一起工作的：

```text
LoginPage
  ↓
auth API / token 持久化
  ↓
router/index.ts
  ↓
dynamic-menu.ts
  ↓
AdminLayout
  ↓
pages/system/*
```

## 推荐阅读顺序

这一章自己的核心正文，建议先按下面顺序读：

1. [前端运行时结构](./frontend-runtime-structure)
2. [管理台工程起步结构](./admin-project-bootstrap)
3. [登录态与会话流转](./login-and-session-flow)
4. [登录页实现细节](./login-page-implementation)
5. [后台壳子、动态菜单与按钮权限](./admin-shell-and-dynamic-menu)
6. [后台布局与工作标签](./admin-layout-and-worktabs)
7. [动态菜单注册与按钮权限](./dynamic-route-registration)
8. [系统模块页面模式](./system-module-pages)
9. [用户管理页实现要点](./user-management-page-detail)
10. [角色与菜单页实现要点](./role-and-menu-page-detail)
11. [配置与文件页实现要点](./config-and-file-page-detail)
12. [日志查询页实现要点](./audit-log-pages)

## 本章覆盖范围

第 7 章围绕当前管理台的完整前端链路展开，重点覆盖：

- `main.ts`、`App.vue`、`router/index.ts` 组成的运行时入口
- `utils/auth.ts`、`api/http.ts` 组成的登录态与请求鉴权链路
- `dynamic-menu.ts`、`AdminLayout.vue` 组成的菜单、路由和工作区骨架
- `pages/system/*` 中各类系统页面的共性模式与具体案例

## 本章完成后的判断标准

这一章完全收稳后，应该至少能回答下面几个问题：

1. 登录态在前端是怎么保存和失效处理的
2. 动态菜单为什么要围绕 `/auth/menus` 展开，而不是手写前端静态路由表
3. 按钮权限码为什么能直接从菜单树递归提取
4. 系统页面为什么统一落在 `pages/system/*`
5. 前端页面如何与第 6 章的系统模块结构一一对应

## 本章小节

- [前端运行时结构](./frontend-runtime-structure)
- [管理台工程起步结构](./admin-project-bootstrap)
- [登录态与会话流转](./login-and-session-flow)
- [登录页实现细节](./login-page-implementation)
- [后台壳子、动态菜单与按钮权限](./admin-shell-and-dynamic-menu)
- [后台布局与工作标签](./admin-layout-and-worktabs)
- [动态菜单注册与按钮权限](./dynamic-route-registration)
- [系统模块页面模式](./system-module-pages)
- [用户管理页实现要点](./user-management-page-detail)
- [角色与菜单页实现要点](./role-and-menu-page-detail)
- [配置与文件页实现要点](./config-and-file-page-detail)
- [日志查询页实现要点](./audit-log-pages)

## 下一步

前端管理台这条线继续收稳后，再进入 [第 8 章：部署、升级与复用](../chapter-8/)。
