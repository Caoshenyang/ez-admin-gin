---
title: 动态菜单组件白名单
description: "集中记录 EZ Admin Gin 当前动态菜单里的 component 白名单、icon 白名单、占位页回退规则，以及新增页面时前后端需要同时对齐的字段。"
---

# 动态菜单组件白名单

这页专门用来查前端动态菜单最容易忘的一件事：

> 后端菜单里的 `component` 和 `icon` 不是随便写字符串，而是必须命中前端白名单。

## 当前动态菜单的真实来源

当前前端不会把菜单写死在路由表里，而是按下面这条链路工作：

```text
/api/v1/auth/menus
  ↓
authMenus
  ↓
dynamic-menu.ts
  ↓
侧边菜单 + 动态路由 + 按钮权限集合
```

相关文件：

- `admin/src/router/dynamic-menu.ts`
- `admin/src/router/index.ts`

## `dynamic-menu.ts` 当前负责三件事

这份文件同时承担了三层职责：

| 能力 | 作用 |
| --- | --- |
| `routeComponentMap` | 把后端 `component` 字符串映射成真实 Vue 页面 |
| `menuIconMap` | 把后端 `icon` 字符串映射成 Naive UI 菜单图标 |
| `collectButtonCodes(...)` | 从按钮节点收集按钮权限码，给页面显隐使用 |

所以它不是一个“普通工具函数文件”，而是：

- 前端动态菜单运行时的白名单中心

## 当前 `component` 白名单

当前真正已经注册的页面组件如下：

| 后端 `component` 值 | 前端实际页面 |
| --- | --- |
| `system/HealthView` | `pages/system/HealthView.vue` |
| `system/UserView` | `pages/system/UserView.vue` |
| `system/RoleView` | `pages/system/RoleView.vue` |
| `system/MenuView` | `pages/system/MenuView.vue` |
| `system/ConfigView` | `pages/system/ConfigView.vue` |
| `system/FileView` | `pages/system/FileView.vue` |
| `system/OperationLogView` | `pages/system/OperationLogView.vue` |
| `system/LoginLogView` | `pages/system/LoginLogView.vue` |
| `system/NoticeView` | `pages/system/NoticeView.vue` |

这意味着当前后端在 `sys_menu.component` 里写的值，必须精确命中这张表。

## 如果 `component` 没命中会发生什么

当前前端不会直接报错崩掉，而是回退到：

- `pages/system/PlaceholderPage.vue`

对应逻辑是：

```ts
function resolveRouteComponent(component: string) {
  return routeComponentMap[component] ?? placeholderPage
}
```

这条回退规则的意义是：

- 动态菜单还能注册成功
- 页面不会空白崩溃
- 但会明确暴露“这个菜单的页面还没有真正接好”

::: warning 菜单能显示，不代表页面已经接通
如果你在后台新增了一个菜单节点，但没有同步更新 `routeComponentMap`，用户仍然可能看见菜单入口，只是点进去会落到占位页。
:::

## 当前 `icon` 白名单

后端 `sys_menu.icon` 字段也不是直接拿来当组件名渲染，而是要先经过前端白名单映射。

当前白名单大致覆盖这些别名：

- `dashboard`
- `analytics` / `audit`
- `health`
- `user` / `users`
- `role` / `roles`
- `menu` / `menus`
- `file` / `files` / `folder`
- `notice` / `notification`
- `log` / `logs`
- `operationlog`
- `loginlog`
- `setting` / `settings`
- `server`

如果后端传了一个前端不认识的值，当前会回退到默认图标：

- `AppsOutline`

另外前端还会先做一次规范化：

```ts
icon.trim().toLowerCase().replace(/[^a-z0-9]/g, '')
```

也就是说：

- 大小写差异会被抹平
- `user-log`、`user_log` 这类符号会被归一化

## 菜单类型怎么影响动态路由和按钮权限

当前菜单模型里有三类节点：

| 类型 | 当前作用 |
| --- | --- |
| 目录 | 用于组织侧边栏层级 |
| 菜单 | 会参与动态路由注册 |
| 按钮 | 不参与路由，但会进入按钮权限集合 |

前端当前是这样处理的：

- `MenuType.Menu` 且有 `path` 的节点，会进入 `buildDynamicRoutes(...)`
- `MenuType.Button` 的节点，会进入 `buttonPermissionCodes`

所以当前按钮显隐的真实来源是：

- `/auth/menus` 返回的按钮节点 `code`

## 当前内置固定菜单

除了后端返回的菜单树，前端还会固定插入一个本地内置菜单：

- `工作台`

它来自：

- `builtinMenuOptions`

对应路由是：

- `/dashboard`

这也是为什么当前即使菜单树里没有专门下发工作台节点，登录后仍然会看到工作台入口。

## 动态路由是怎么注册的

当前路由启动逻辑在：

- `admin/src/router/index.ts`

大致流程是：

1. 用户持有 Token 访问受保护页面
2. `beforeEach` 首次请求 `/auth/menus`
3. `setAuthMenus(menus)`
4. `buildDynamicRoutes(menus)` 生成路由记录
5. `router.addRoute('admin', route)` 注册到后台壳子下
6. 重新匹配当前目标地址

这说明当前菜单、路由和按钮权限是一次请求同时驱动的三件事。

## 新增页面时前后端要一起对齐什么

如果你要给一个新模块补页面，当前最稳的对齐清单是：

1. 前端先有真实页面文件，例如 `pages/system/ExampleView.vue`
2. `dynamic-menu.ts` 里补一条 `routeComponentMap`
3. 后端 `sys_menu.component` 写成完全一致的编码，例如 `system/ExampleView`
4. 后端 `sys_menu.path` 写成真实菜单路径
5. 后端 `sys_menu.code` 补页面和按钮节点
6. 页面里按需要消费 `buttonPermissionCodes`

## 最常见的三类问题

### 1. 菜单显示了，但点开是占位页

通常说明：

- `sys_menu.component` 已经存在
- 但前端 `routeComponentMap` 没补

### 2. 页面能打开，但侧边栏图标不对

通常说明：

- `sys_menu.icon` 没命中前端 `menuIconMap`
- 最终回退到了默认图标

### 3. 页面功能按钮一直不显示

通常说明：

- 菜单树里没有对应按钮节点
- 或按钮节点 `code` 和页面判断用的权限码不一致

## 当前这页最适合拿来做什么

你可以把它当成三种排障清单：

1. 查一个 `component` 字符串到底该写什么
2. 查一个 `icon` 字符串当前前端认不认
3. 查一个按钮权限为什么没落到页面上

## 相关教程与参考页

- [权限码约定](./permission-code-conventions)
- [初始化数据参考](./init-data-reference)
- [前端概览](/frontend/overview)
- [模块开发](/backend/module-development)
