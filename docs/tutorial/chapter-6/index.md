---
title: 第 6 章：核心系统模块
description: "按当前最终版目录结构讲清用户、角色、菜单、配置、文件、日志和公告模块如何落进 module/system 与 module/iam。"
---

# 第 6 章：核心系统模块

前面几章已经把平台底座、认证授权、组织体系和数据权限的骨架定稳了。现在开始进入真正会被后台管理台长期消费的那一层：核心系统模块。

这一章的重点，不再是“功能能不能跑”，而是：

> 这些常用系统能力，如何按当前最终结构稳定落在 `module/*` 下面。

::: tip 本章怎么读
这一章建议带着两个问题往下看：

- 当前一个模块到底应该落在 `auth / setup / system / iam` 的哪一层
- 为什么每个模块要沿着 `dto / repository / service / handler / routes` 这条固定结构落地
:::

## 本章会解决什么

这一章会把下面几类高频后台模块，统一放到当前最终版结构里理解：

- 用户、角色、部门、岗位、菜单这类 IAM 资源
- 配置、文件、登录日志、操作日志、公告这类系统资源
- 数据字典、附件中心、账户中心等常用通用模块
- 模块路由如何从 `bootstrap` 进入，再汇总到 `module/system`
- 模块内部为什么要继续保持职责拆分，而不是重新回到 Handler 直连一切

## 当前真实模块主线

现在仓库里的核心系统模块已经大体收敛到下面这条路径：

```text
main.go (embed)
  ↓
bootstrap/run.go
  ↓
bootstrap/router.go
  ↓
module/auth
module/setup
module/iam/*
module/system/*
```

这意味着第 6 章现在讲的，不是”怎么在某个全局路由文件里多挂一个 handler”，而是：

- 如何把一个模块定义成独立边界
- 如何把它纳入统一的系统路由聚合
- 如何让它继续复用前面几章已经定稳的认证、权限、数据权限和日志链路

## 本章建议顺序

### 1. 先看模块固定结构

先看 [模块固定结构](./module-structure)。

这一页会把当前最终版模块长什么样讲清楚，尤其是：

- `entity / dto / repository / service / handler / routes`
- 哪些模块会额外带 `datascope` 或 `policy`
- 什么情况下复用 `internal/model/*`，什么情况下在模块内单独放 `entity.go`

### 2. 再看内置模块落地流程

再看 [内置模块落地流程](./backend-module-flow)。

这一页会按当前真实结构，把一个系统内置模块从”准备接入”讲到”能被系统路由聚合并对外提供接口”。

### 3. 再看具体核心模块落地

接下来看几个具体的核心模块如何按同一套骨架落地：

- [系统模块示例](./sample-module) — 公告模块
- [数据字典模块落地](./dict-module) — 系统级公共字典
- [附件中心落地](./attachment-center-module) — 复用底层上传能力的业务资源模型
- [账户中心落地](./account-center-module) — 会话内自助能力

### 4. 再看权限、菜单与迁移接入

继续看 [权限、菜单与迁移接入](./permission-menu-migration)。

### 5. 再看前端页面接入流程

再看 [前端页面接入流程](./frontend-page-flow)。

### 6. 最后用验收清单自查

最后用 [模块接入验收清单](./module-integration-checklist) 来判断一个模块是否真正进入后台底座。

## 本章完成后的判断标准

完成这一章后，你至少应该能回答下面几个问题：

1. 为什么当前主线要优先扩 `module/system` 和 `module/iam/*`
2. 一个模块为什么要拆成 `dto / repository / service / handler / routes`
3. 哪些模块需要 `datascope.go`，哪些暂时不需要
4. 为什么系统路由现在由 `bootstrap/router.go` 统一装配
5. 数据字典、附件中心、账户中心分别适合放在哪一层

::: tip 当前主线已完全收敛到 `module/*`
旧的全局 Handler 目录和路由文件已在代码清理中移除。当前所有模块统一沿 `module/<group>/<resource>` 落地，不再存在新旧结构并存的局面。
:::

::: info 边界提醒
本章把后端模块结构、核心模块落地、权限菜单接入、前端页面接入和验收清单统一收口。前端运行时主线放在 [第 7 章](../chapter-7/)，部署、升级与复用放在 [第 8 章](../chapter-8/)。
:::

## 本章小节

- [模块固定结构](./module-structure)
- [内置模块落地流程](./backend-module-flow)
- [系统模块示例](./sample-module)
- [数据字典模块落地](./dict-module)
- [附件中心落地](./attachment-center-module)
- [账户中心落地](./account-center-module)
- [权限、菜单与迁移接入](./permission-menu-migration)
- [前端页面接入流程](./frontend-page-flow)
- [模块接入验收清单](./module-integration-checklist)

## 这一章结束后会走到哪里

当第 6 章把核心系统模块结构收稳后，后面的前端管理台和部署复用就都会有更清晰的后端对照面。

也就是说，这一章真正交付的不是几组 CRUD，而是：

> 一套可以继续承载系统模块和业务模块扩展的后端模块边界。

下一节先看当前主线里的模块固定结构：[模块固定结构](./module-structure)。
