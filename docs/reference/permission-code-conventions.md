---
title: 权限码约定
description: "集中说明当前底座里 policy.go、菜单 code、component 编码和按钮权限码之间的关系，以及命名时最值得保持稳定的约定。"
---

# 权限码约定

这页只负责回答一个问题：

> 当前底座里，接口权限、菜单编码、前端按钮权限码到底各自扮演什么角色，命名时应该怎么保持稳定。

::: tip 🎯 这页解决什么
快速确认下面这些问题：

- `policy.go` 里的权限码是给谁用的
- `sys_menu.code` 和按钮显隐是什么关系
- `component` 编码为什么必须稳定
:::

## 先看当前三类“权限相关编码”

当前底座里，最容易混淆的是下面三类内容：

| 类型 | 当前落点 | 主要用途 |
| --- | --- | --- |
| 接口权限码 | `module/*/policy.go` | 给 Casbin 和后端接口鉴权用 |
| 菜单 / 按钮编码 | `sys_menu.code` | 给前端菜单树和按钮显隐用 |
| 页面组件编码 | `sys_menu.component` | 给前端动态路由映射真实页面组件 |

这三者经常会互相关联，但不应该被当成同一个东西。

## `policy.go` 里的权限码是做什么的

当前每个模块都会把稳定权限点收在自己的 `policy.go` 里，例如：

```go
const (
	PermissionList   = "system:file:list"
	PermissionUpload = "system:file:upload"
)
```

它们的主要作用是：

- 作为接口权限的稳定命名
- 让 Casbin 策略和业务模块共用同一套标识
- 避免权限码散落在 Handler、Service、前端页面里

当前已经落地的典型命名大多是：

```text
<group>:<resource>:<action>
```

例如：

- `system:user:list`
- `system:user:create`
- `system:role:update_permissions`
- `system:menu:delete`
- `system:operation-log:list`

## 为什么 `policy.go` 要尽量稳定

因为它一旦变化，会同时影响：

- Casbin 策略种子
- 接口权限判断
- 角色权限配置
- 排障时对“这个接口到底该有什么权限”的理解

所以当前主线里，`policy.go` 的职责不是“凑几个常量”，而是：

> 为模块固定一套长期可复用的权限码命名。

## `sys_menu.code` 和 `policy.go` 是什么关系

当前菜单树里，每个节点也有一个稳定 `code`。

它和 `policy.go` 的关系通常是：

- 页面菜单节点：表达页面入口
- 按钮节点：表达页面内具体操作点

前端页面里的典型消费方式是：

```ts
buttonPermissionCodes.value.includes("system:notice:create")
```

也就是说，按钮显隐依赖的是：

- 当前登录用户从 `/auth/menus` 拿到的菜单树
- 菜单树里按钮节点的 `code`

## 按钮权限码为什么通常和接口权限码保持一致

虽然它们语义上不是同一层，但当前主线里更推荐：

- 按钮权限码尽量和对应接口权限码同名

例如：

- 新增按钮：`system:user:create`
- 编辑按钮：`system:user:update`
- 启停按钮：`system:user:update_status`

这样做的好处是：

- 前后端沟通成本低
- 角色配置和页面显隐更容易对照
- 排查“为什么按钮不显示、接口也不通”时更直接

::: warning ⚠️ 不要把按钮码和接口码随意命成两套
如果前端页面判断的是 `system:user:assign-role`，而后端接口策略写的是另一套完全不同的命名，后面查问题会非常痛苦。
当前更稳的做法是：能复用同一命名就尽量复用。
:::

## `component` 编码为什么也属于这条链路的一部分

菜单节点里还有一个很关键的字段：

- `component`

它本身不是权限码，但它决定了：

- `/auth/menus` 返回这个菜单节点后
- 前端到底把它映射到哪个真实页面组件

当前前端会在 `dynamic-menu.ts` 里维护一张白名单映射：

- `system/NoticeView`
- `system/UserView`
- `system/RoleView`

这意味着：

- `code` 决定入口和按钮语义
- `component` 决定页面落点

两者都必须稳定。

## 当前权限码最常见的动作命名

当前模块里最常见的动作大致可以归成下面几类：

| 动作 | 常见命名 |
| --- | --- |
| 查询列表 | `list` |
| 创建 | `create` |
| 编辑 | `update` |
| 状态切换 | `update_status` / `status` |
| 删除 | `delete` |
| 上传 | `upload` |
| 改角色 | `update_roles` |
| 改接口权限 | `update_permissions` |
| 改菜单授权 | `update_menus` |

当前代码里有个小的历史差异：

- 有的模块用 `update_status`
- 有的模块用 `status`

这说明命名正在逐步收口，但还没有完全统一。

## 当前已经落地的典型权限码

可以直接用下面这些作为现有主线里的代表：

| 模块 | 典型权限码 |
| --- | --- |
| 用户 | `system:user:list`、`system:user:create`、`system:user:update` |
| 角色 | `system:role:update_permissions`、`system:role:update_menus` |
| 菜单 | `system:menu:create`、`system:menu:delete` |
| 部门 | `system:department:update_status` |
| 岗位 | `system:post:update_status` |
| 文件 | `system:file:upload` |
| 配置 | `system:config:value` |
| 日志 | `system:operation-log:list`、`system:login-log:list` |

## 新模块命名时最值得遵守的规则

如果你要给新模块补 `policy.go`、菜单和按钮节点，当前最稳的规则是：

1. 先固定资源分组，例如 `system`, `project`, `sales`
2. 再固定资源名，例如 `customer`, `task`, `invoice`
3. 最后只给动作命名，不把上下文塞进动作名里

也就是优先写成：

```text
project:customer:list
project:customer:create
project:customer:update
```

而不是：

```text
customer:list:project
customer:create:button
```

## 什么最容易导致权限链路失效

当前最常见的几种失效方式是：

1. `policy.go` 改了，但 Casbin 策略种子没同步
2. 菜单按钮节点的 `code` 和页面里的 `canUse(code)` 对不上
3. `component` 编码没有命中前端白名单
4. 角色授权改完后，没有意识到接口权限当前不会自动热刷新

## 相关代码和教程页

如果你要继续查真实实现，优先看：

- `server/internal/module/*/policy.go`
- `server/internal/model/menu.go`
- `admin/src/router/dynamic-menu.ts`

对应教程页：

- [数据库迁移](/backend/migration)
- [路由与菜单](/frontend/route-and-menu)
- [前端概览](/frontend/overview)
