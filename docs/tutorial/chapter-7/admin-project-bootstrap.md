---
title: 管理台工程起步结构
description: "围绕当前 admin/ 真实工程，讲清依赖、入口、全局 Provider、样式入口和页面目录为什么这样组织。"
---

# 管理台工程起步结构

第 7 章后面的登录页、后台壳子、动态菜单和系统页面，都是建立在同一个前端起步工程上的。

所以在继续往下看具体页面前，先把这一层站稳很重要：

> 当前 `admin/` 这个管理台项目，到底是怎样被组织起来的。

::: tip 🎯 本节目标
读完后，你应该能快速看懂当前前端工程的基础骨架，知道入口文件、全局样式、路由、页面目录和通用依赖分别落在哪里。
:::

## 先看当前真实骨架

当前前端项目的主骨架已经很稳定：

```text
admin/
├─ package.json
├─ vite.config.ts
└─ src/
   ├─ App.vue
   ├─ main.ts
   ├─ api/
   ├─ components/
   ├─ layouts/
   ├─ pages/
   ├─ router/
   ├─ styles/
   ├─ types/
   └─ utils/
```

这里最值得先记住的不是每个文件名，而是职责边界：

| 位置 | 当前主要职责 |
| --- | --- |
| `api/` | 资源级请求封装 |
| `layouts/` | 后台壳子组件 |
| `pages/` | 页面组件本体 |
| `router/` | 路由和动态菜单注册 |
| `styles/` | 全局样式与主题变量 |
| `types/` | 前后端共享的数据结构类型 |
| `utils/` | Token、本地状态和轻工具 |

## 当前依赖栈为什么是这套组合

从 `admin/package.json` 看，当前项目的核心组合是：

- `vue`
- `vue-router`
- `pinia`
- `naive-ui`
- `axios`
- `tailwindcss`

这套组合的分工很明确：

| 依赖 | 当前作用 |
| --- | --- |
| `vue` | 页面和响应式基础 |
| `vue-router` | 登录页、后台壳子和动态路由 |
| `pinia` | 预留状态管理能力 |
| `naive-ui` | 表单、表格、布局、弹窗、菜单等后台组件 |
| `axios` | 请求层封装 |
| `tailwindcss` | 页面骨架、间距、背景和响应式布局 |

这说明当前前端不是“全靠组件库”，也不是“全靠原子类”，而是有意识地把两者拆开使用。

## `main.ts` 现在只做一件事

当前 `admin/src/main.ts` 很干净，它只负责：

1. 引入全局样式
2. 创建 Vue 应用
3. 安装 Pinia
4. 安装 Router
5. 挂载到 `#app`

这种入口越薄越好，因为它意味着：

- 页面初始化逻辑不会堆在根入口里
- 后续新增能力可以更容易判断该落在哪一层

## `App.vue` 为什么只保留全局 Provider

当前 `App.vue` 的角色也非常清楚，它几乎不承担业务页面逻辑，只负责挂一层全局环境：

- `NConfigProvider`
- `NLoadingBarProvider`
- `NDialogProvider`
- `NNotificationProvider`
- `NMessageProvider`
- `RouterView`

也就是说，根组件当前的主要任务是：

> 给整套后台页面提供统一的 Naive UI 运行环境。

这让后面的登录页、布局页、系统页都能直接用消息、通知、对话框，而不用各自重复搭环境。

## 当前样式入口为什么集中在 `styles/main.css`

当前全局样式不是散落在多个地方，而是统一从：

- `admin/src/styles/main.css`

进入。

这里现在主要做了三件事：

1. `@import "tailwindcss"`
2. 定义主题变量
3. 设置页面级基线样式

### 当前这份样式基线最重要的点是什么

最关键的一条，其实不是颜色变量，而是：

- 关闭浏览器级滚动

因为后面整个后台壳子都依赖“一屏应用”的布局心智。如果这里不先收住，`AdminLayout.vue` 很快就会被浏览器原生滚动条打散。

## `vite.config.ts` 为什么已经包含代理和 Tailwind

当前 `vite.config.ts` 已经承担了两件非常实用的开发期配置：

- 接入 `tailwindcss()` Vite 插件
- 把 `/api` 代理到本地后端 `http://localhost:8080`

它的意义分别是：

- 前端可以直接使用 Tailwind 4 的 CSS-first 方案
- 页面代码里只写 `/api/v1/...`，不需要硬编码完整后端域名

这样第 7 章后面所有接口页都能沿用同一套联调方式。

## 当前页面目录已经怎样分层

现在 `admin/src/pages/` 下已经能看出很清楚的分区：

```text
pages/
├─ auth/
├─ dashboard/
└─ system/
```

这和教程主线是完全对齐的：

- `auth/`：登录入口
- `dashboard/`：后台默认首页
- `system/`：系统模块真实管理页

这说明前端页面目录并不是随意长出来的，而是已经在围绕后台产品结构组织。

## 为什么 `components/` 仍然要单独保留

虽然当前第 7 章重点在页面，但 `components/` 仍然有存在价值，例如：

- `BrandLogo.vue`

这种组件既会在登录页出现，也会在后台壳子里出现。

如果把这类可复用小组件硬塞进页面文件里，后面维护成本会越来越高。所以当前前端仍保留了一个很轻的共享组件层。

## 当前这套起步工程和后面的章节是什么关系

可以把它理解成第 7 章后续所有内容的共同基座：

- [登录页实现细节](./login-page-implementation) 会复用 `App.vue` 和 `api/http.ts`
- [后台布局与工作标签](./admin-layout-and-worktabs) 会复用 `styles/main.css` 和 `router/`
- [动态菜单注册与按钮权限](./dynamic-route-registration) 会继续扩展 `router/`
- [系统模块页面模式](./system-module-pages) 会继续在 `pages/system/*` 上展开

所以这一页虽然不讲具体业务，但它实际上决定了后面每一页“该放哪儿、怎么接”。

## 一份快速自检清单

如果你要确认当前管理台起步工程是不是在正确基线上，可以快速看下面几项：

1. `main.ts` 是否只负责挂载应用与安装插件
2. `App.vue` 是否主要承担全局 Provider，而不是业务逻辑
3. `styles/main.css` 是否已经成为全局样式入口
4. `vite.config.ts` 是否已经包含 Tailwind 插件和 `/api` 代理
5. `pages/` 是否已经按 `auth / dashboard / system` 分层

## 下一步

- 想继续看这套起步工程如何进入真实登录链路，下一页读 [登录页实现细节](./login-page-implementation)
- 想继续看登录后的后台页面容器怎么接起来，读 [后台布局与工作标签](./admin-layout-and-worktabs)
