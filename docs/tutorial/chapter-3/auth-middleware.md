---
title: 登录校验中间件
description: "把 access token 接进请求链路，并在认证通过后加载 Actor 上下文，让后续接口既知道当前是谁，也知道当前人的组织与数据范围。"
---

# 登录校验中间件

上一节已经拿到了 `access_token`。这一节开始把它真正接进请求链路：请求进入受保护接口前，先完成 Token 校验，再把当前用户身份和更进一步的 Actor 上下文写入 Gin 上下文。

::: tip 🎯 本节目标
完成后，`/api/v1/auth/me`、`/api/v1/auth/menus`、`/api/v1/auth/dashboard` 这类受保护接口都必须携带有效 Token；认证通过后，不只可以拿到 `user_id` 和 `username`，还可以继续加载组织和数据范围摘要。
:::

## 本节会改什么

当前主线里，这一节主要对应下面这些位置：

```text
server/
├─ internal/
│  ├─ middleware/
│  │  ├─ auth.go
│  │  └─ actor.go
│  └─ module/
│     └─ auth/
│        ├─ me_handler.go
│        ├─ me_service.go
│        └─ routes.go
```

| 位置 | 用途 |
| --- | --- |
| `middleware/auth.go` | 解析 `Authorization` 请求头，校验 Token，并写入最基础身份信息 |
| `middleware/actor.go` | 在认证通过后继续加载角色、部门和数据范围摘要 |
| `module/auth/me_handler.go` | 对外暴露 `/api/v1/auth/me` |
| `module/auth/routes.go` | 给 `/auth/me`、`/auth/menus`、`/auth/dashboard` 挂载认证链路 |

## 当前请求链路已经不是“只认 Token”

这一节最重要的变化是：当前认证链路已经分成两层，而不是只有一个简单的 JWT 校验中间件。

### 第一层：`Auth`

负责：

- 从 `Authorization` 头中解析 Bearer Token
- 校验签名和过期时间
- 把 `current_user_id` 和 `current_username` 写入上下文

### 第二层：`LoadActor`

负责：

- 根据当前用户继续查角色
- 读取部门归属
- 聚合数据范围授权
- 把 Actor 上下文写入 Gin Context

这一步是企业级后台和普通 Demo 差异很明显的地方，因为后面真正的数据权限和工作台摘要都依赖 Actor，而不是只依赖一个 `user_id`。

::: warning ⚠️ 认证通过不等于 Actor 一定完整
`Auth` 负责“这个 Token 是否有效”，`LoadActor` 负责“这个人当前拥有哪些组织与数据范围信息”。这两个阶段的职责不要混在一起。
:::

## 请求头约定仍然保持简单稳定

后续所有需要登录的接口，都按同一格式传递 Token：

```http
Authorization: Bearer <access_token>
```

当前中间件会严格检查：

- 是否存在 `Authorization`
- 是否是 `Bearer`
- Token 是否能被成功解析

## `/auth/me` 现在已经不只是一个“回显用户 ID”的接口

当前主线里，`/api/v1/auth/me` 的价值已经升级了。  
它现在不只用来验证“我登录了没有”，还承担两件事：

- 给前端返回当前登录用户的基础摘要
- 预览当前用户的组织和数据范围信息

这也是为什么这一节现在会和后面的第 5 章产生自然衔接：  
你在这里先看到 `Actor` 和数据范围摘要，下一章再继续深入它们真正如何参与查询过滤。

## 为什么认证链路要提前挂上 `LoadActor`

很多后台 Demo 会在这里停在“Token 校验通过就结束了”。当前主线没有这么做，是因为后面这些能力都需要 Actor：

- `/auth/me` 数据范围摘要
- 用户列表的数据权限过滤
- 部门和岗位的组织边界
- 前端工作台的当前用户上下文

如果等到第 5 章才临时补 Actor，中间会出现大量重复改路由和 Handler 的返工。  
所以这一步虽然看起来比普通教程重一点，但它恰好是在替后面省成本。

## 怎么验证这一节已经做成

### 1. 不带 Token 会直接失败

请求：

```text
GET /api/v1/auth/me
```

如果没有 `Authorization` 请求头，应该直接返回 `401`。

### 2. 带错误 Token 也会失败

把一个明显错误的 Token 放进：

```http
Authorization: Bearer wrong-token
```

同样应该返回 `401`，说明第一层 `Auth` 已经接管住登录态校验。

### 3. 带有效 Token 时 `/auth/me` 能返回更完整摘要

使用有效 Token 请求：

```text
GET /api/v1/auth/me
```

当前主线下，响应里应该已经能看到这类字段：

- `user_id`
- `username`
- `department_id`
- `role_codes`
- `is_super_admin`
- `data_scope`

这说明第二层 `LoadActor` 也已经接进来了。

## 本节最关键的收获

这一节真正建立的判断标准是：

> 企业级后台里的认证中间件，不只是“看一下 Token 对不对”，而是要把后续权限和数据范围真正会用到的上下文提前加载好。

这一节完成后，请求链路里已经不只是“当前是谁”，而是开始拥有“这个人在系统里处于什么组织和权限位置”的基础上下文。

下一节继续补齐前端真正会消费的菜单权限链路：[角色菜单权限](./menu-permission)。
