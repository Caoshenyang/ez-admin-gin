---
title: 文档治理说明与基线清单
description: "记录当前文档主线中的 canonical 页面、归档策略和后续维护规则。"
search: false
---

# 文档治理说明与基线清单

::: warning 维护页说明
这是一张偏维护者视角的治理页，不是新读者的默认入口。第一次进入项目时，请优先从 [快速启动](/guide/) 或 [从零搭建教程](/tutorial/) 开始。
:::

这页只做一件事：固定当前文档主线里的页面边界，避免后续继续一边补正文、一边重新判断“这页到底算不算主线”。

::: tip 🎯 这页怎么用
如果你在继续维护文档，先看这页再决定要不要改某一页：

- `canonical`：主线页，可以继续维护、补充、挂在导航里。
- `残留页`：内容还有查阅价值，但不再作为主线入口。
- `历史页`：已经和当前章节边界冲突，只保留说明价值，不再参与主线导航。
:::

## Canonical：当前主线页

当前主线的判断规则固定为：

- 教程页以 [教程大纲](/tutorial/curriculum) 里出现的页面为准。
- 指南页以 `/guide/` 侧边栏里保留的入口页为准。
- 参考页以 [参考手册首页](/reference/) 当前列出的页面为准。

也就是说，后续只要一页没有进入这些入口，就默认不属于当前主线。

### 当前教程 canonical 边界

- 第 1-2 章：项目初始化、后端基础设施
- 第 3 章：登录、JWT、认证中间件
- 第 4 章：RBAC、Casbin、角色菜单权限
- 第 5 章：组织体系、`Actor`、`gorm.Scopes(...)`、资源模式、模块接入规范、请求走读、验收清单
- 第 6 章：模块固定结构、后端模块接入、示例模块、权限菜单接入
- 第 7 章：前端运行时、登录态、布局、动态菜单、系统页模式和详细案例
- 第 8 章：模块化接入规范
- 第 9 章：部署、升级、回滚、排障、长期运维、项目复用

## 残留页：保留内容价值，但不再作为主线入口

这些页面仍然有查阅价值，但不再代表当前章节主线：

- `tutorial/chapter-4/user-management.md`
- `tutorial/chapter-4/role-management.md`
- `tutorial/chapter-4/menu-management.md`
- `tutorial/chapter-4/system-config.md`
- `tutorial/chapter-4/file-upload.md`
- `tutorial/chapter-4/operation-logs.md`
- `tutorial/chapter-4/login-logs.md`
- `tutorial/chapter-4/notice-management.md`

这些页面主要保留“模块个案说明”价值。当前主线里，系统模块结构以第 6 章为准，前端系统页面消费以第 7 章为准。

## 历史页：退出主线，不再参与章节导航

这些页面和当前章节边界已经冲突，只保留历史说明价值：

### 旧的权限页

- `tutorial/chapter-3/rbac-model.md` → stub
- `tutorial/chapter-3/casbin-permission.md` → stub
- `tutorial/chapter-3/menu-permission.md` → stub
- `tutorial/archive/chapter-3/*` → 归档正文

当前 canonical 位置已经切到第 4 章。

### 旧的前端页

- `tutorial/chapter-5/vue-project-init.md` → stub
- `tutorial/chapter-5/login-page.md` → stub
- `tutorial/chapter-5/admin-layout.md` → stub
- `tutorial/chapter-5/dynamic-menu.md` → stub
- `tutorial/chapter-5/user-pages.md` → stub
- `tutorial/chapter-5/role-menu-pages.md` → stub
- `tutorial/chapter-5/config-file-pages.md` → stub
- `tutorial/chapter-5/log-pages.md` → stub
- `tutorial/archive/chapter-5/*` → 归档正文

当前 canonical 位置已经切到第 7 章。

### 旧的模块接入页

- `tutorial/chapter-6/frontend-page-flow.md` → stub
- `tutorial/archive/chapter-6/frontend-page-flow.md` → 归档正文

当前 canonical 位置已经切到第 8 章。

### 旧的部署页

- `tutorial/chapter-7/env-and-init-data.md` → stub
- `tutorial/chapter-7/deployment-and-reuse.md` → stub
- `tutorial/archive/chapter-7/*` → 归档正文

当前 canonical 位置已经切到第 9 章。

## 后续维护规则

后续继续写文档时，固定遵守下面 4 条规则：

1. 新内容优先进入 canonical 页面，不再平行新增一套旧章节版本。
2. 旧链接可以保留 stub，但旧正文统一进入 `tutorial/archive/`。
3. 任何页面迁移，都要同步检查 `README`、站点首页、VitePress 侧边栏和站内 canonical 链接。
4. 同主题只允许一个 canonical 页面；旧页面不能继续承载完整主线正文。
5. 状态页只在阶段完成或阶段切换时更新一次，不再随着单页改动频繁改写。
