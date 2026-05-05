---
title: Actor 上下文与多角色并集
description: "围绕 middleware.LoadActor 和 datascope.Merge(...)，讲清当前登录人上下文怎样装载，以及多角色数据范围为什么统一按并集处理。"
---

# Actor 上下文与多角色并集

第 5 章真正让数据权限站稳的，不只是 `sys_role.data_scope` 这个字段，而是：

> 登录后的当前用户上下文里，是否已经装好了后续查询真正要用的数据范围信息。

如果这一层不先定稳，后面的每个模块都要自己重新查角色、重新解释范围、重新判断超级管理员，最后几乎一定会散。

::: tip 🎯 本节目标
读完后，你应该能顺着当前真实代码回答三件事：

1. `LoadActor` 在请求进入业务模块前到底装了什么
2. 多角色数据范围为什么统一按并集解释
3. 为什么 `RoleCodes` 和 `Grants` 必须同时存在
:::

## 先看当前真实落点

当前这条链路主要对应两个地方：

```text
server/internal/middleware/actor.go
server/internal/platform/datascope/datascope.go
```

职责分工很清楚：

| 位置 | 作用 |
| --- | --- |
| `middleware/actor.go` | 从数据库装载当前登录人的组织与数据范围上下文 |
| `platform/datascope/datascope.go` | 定义数据范围枚举、`Actor`、`Grant` 和并集合并规则 |

## 为什么请求期一定要先有 `Actor`

如果数据权限只停留在数据库表里，那么每个模块在真正查数据前，都得自己回答下面这些问题：

- 当前登录人是谁
- 当前登录人属于哪个部门
- 当前登录人有哪些角色
- 每个角色的数据范围是什么
- 有没有自定义部门授权
- 当前是不是超级管理员

这些判断如果分散在 Handler、Service、Repository 里，后面排查“为什么这个人这里能看、那里不能看”会非常痛苦。

所以当前主线先在认证通过后，把这批信息一次收进统一上下文：

```go
type Actor struct {
	UserID       uint
	Username     string
	DepartmentID uint
	RoleCodes    []string
	Grants       []Grant
	IsSuperAdmin bool
}
```

## `Actor` 里的字段分别给谁用

这几个字段虽然都来自当前登录人，但服务的链路并不一样：

| 字段 | 主要用途 |
| --- | --- |
| `UserID` | `self` 范围、审计、资源归属判断 |
| `Username` | 日志、响应摘要、调试排查 |
| `DepartmentID` | `dept` / `dept_and_children` 范围 |
| `RoleCodes` | Casbin 和接口权限链路继续使用 |
| `Grants` | 数据权限链路使用 |
| `IsSuperAdmin` | 快速绕过数据权限限制 |

::: info 为什么这里不只保留 `user_id`
企业后台里，“当前是谁”通常只是最基础的信息。真正影响查询结果的，是：

- 他属于哪个部门
- 他拥有哪些角色
- 这些角色各自授予了什么数据范围

所以 `Actor` 必须比“一个用户 ID”更完整。
:::

## `Grant` 为什么不直接复用角色结构

当前平台层没有直接把整个角色模型塞进 `Actor`，而是压成更轻的一条授权结构：

```go
type Grant struct {
	Scope         Scope
	DepartmentIDs []uint
}
```

这样做的好处是：

- 数据权限只保留自己真正关心的信息
- 不把角色名称、排序、备注之类无关字段带进查询链路
- 后面做合并时，输入结构更稳定

换句话说，`Grant` 不是角色实体，而是：

> 一条已经准备好给数据权限系统消费的授权事实。

## `LoadActor` 当前到底做了什么

当前 `middleware.LoadActor` 的真实执行顺序，可以直接概括成四步：

1. 根据 `CurrentUserID` 查出用户基础信息
2. 查出当前用户绑定的启用角色
3. 把角色的 `data_scope` 压成 `Grants`
4. 如果角色是 `custom_dept`，再补查它被授权的部门 ID

实际代码里，最关键的落点就是：

```go
roleCodes = append(roleCodes, row.Code)
if row.Code == "super_admin" {
	isSuperAdmin = true
}
grants = append(grants, datascope.Grant{
	Scope: datascope.Scope(row.DataScope),
})
```

以及后面的 `attachCustomDeptGrants(...)`：

- 它会额外查询 `sys_role_data_scope`
- 再把自定义部门列表挂回对应的 `Grant`

这说明当前主线不是“查到角色就结束”，而是已经把角色压成了真正可计算的数据权限输入。

## 为什么 `RoleCodes` 和 `Grants` 要并存

这是第 5 章里一个很关键的分层判断。

虽然它们都来源于角色，但服务的是两条不同链路：

| 数据 | 服务哪条链路 |
| --- | --- |
| `RoleCodes` | 接口权限、菜单权限、Casbin |
| `Grants` | 数据权限、查询过滤 |

如果把这两者混成一个模糊结构，后面就很容易出现：

- 菜单权限和数据权限互相影响
- 模块读到角色后，不知道该按哪套语义解释

所以当前结构选择了“来源一致、职责分离”。

## 多角色为什么固定按并集

单角色很好理解，真正麻烦的是一个用户同时拥有多个角色。

例如：

- 角色 A：`dept`
- 角色 B：`custom_dept`
- 角色 C：`self`

这时系统必须先回答一个问题：

> 最终到底按交集还是并集处理？

当前平台层已经把答案固定为：

- 并集

并且统一收在：

```text
server/internal/platform/datascope/datascope.go
```

## `datascope.Merge(...)` 当前怎么合并

当前并集规则可以直接总结成下面这张表：

| 情况 | 最终摘要 |
| --- | --- |
| 任一角色是 `all` | 直接允许全部数据 |
| 有 `dept` | 允许本部门数据 |
| 有 `dept_and_children` | 允许本部门及整棵子树 |
| 有 `custom_dept` | 允许自定义授权部门集合 |
| 有 `self` | 允许本人数据 |
| 当前用户是 `super_admin` | 直接绕过数据权限 |

平台层把这些角色授权最后压成一份 `Summary`：

```go
type Summary struct {
	AllowAll            bool
	RequireSelf         bool
	IncludeDepartment   bool
	IncludeDeptTree     bool
	CustomDepartmentIDs []uint
}
```

这个结构最大的价值是：

- 所有模块看到的都是统一语义
- 超级管理员绕过逻辑只有一份
- 自定义部门集合只需要去重一次

## 为什么这一步必须放在平台层

如果把并集规则交给各模块自己解释，很快就会出现这种失控情况：

- 用户模块按并集
- 部门模块按交集
- 真实业务模块又自己补了第三种解释

当前这套结构刻意把这件事提前锁死在平台层，就是为了避免：

> 同一个用户在不同模块里，数据权限语义不一致。

::: warning ⚠️ 不要让模块自己定义多角色语义
模块层最应该声明的是“这个资源按哪个字段过滤”，而不是重新定义“多角色到底怎么合并”。
当前主线的可复用性，很大一部分就来自这条边界足够清楚。
:::

## `/auth/me` 为什么也要复用这条链路

当前 `/auth/me` 返回的数据范围摘要，并不是另写一套逻辑，而是同样复用了：

- `Actor`
- `datascope.Merge(...)`

这意味着：

- 前端看到的当前登录人数据范围摘要
- Repository 真正执行查询时使用的合并语义

是同一套规则。

这件事很重要，因为它避免了“前端显示一种范围，后端实际执行另一种范围”。

## 这套写法和 Java 常见方案怎么对照

如果你主要做 Java，可以这样理解：

| 当前 Go 写法 | 常见 Java 对照思路 |
| --- | --- |
| `Actor` | `LoginUser` / `UserContext` |
| `Grant` | 已压平的数据权限授权项 |
| `datascope.Merge(...)` | `DataPermissionContext` 聚合规则 |
| `Summary` | AOP / 数据权限插件最终生成的可执行摘要 |

最大的差别不在思想，而在落点：

- Java 项目更常见注解、AOP、SQL 拦截
- 当前这套 Go 主线更偏显式装载、显式合并、显式应用

这会让排查路径更直接。

## 怎么验证这一节已经成立

### 1. `/api/v1/auth/me` 已经返回数据范围摘要

登录后请求：

```text
GET /api/v1/auth/me
```

至少应该能看到：

- `department_id`
- `role_codes`
- `is_super_admin`
- `data_scope`

这说明 `Actor` 已经被成功装载，并且聚合规则已经进了对外响应。

### 2. 多角色用户的 `role_codes` 和 `data_scope` 不再只代表单一角色

如果你给同一个用户绑定多个角色，再请求 `/auth/me`，应该能看到：

- `role_codes` 是角色集合
- `data_scope` 体现的是合并后的摘要，而不是某一个角色原值

### 3. `custom_dept` 角色已经能带出授权部门

给某个角色配置：

- `data_scope = custom_dept`

并绑定若干部门后，请求 `/auth/me`` 或相关资源列表时，应该能看到这些部门授权已经进入当前登录人的范围摘要或实际查询结果。

## 下一步

- 想继续看这些平台规则怎样真正压进资源查询，下一页读 [资源级数据权限接入模式](./module-datascope-patterns)
- 想先回到总览页重新串主线，读 [角色数据范围与查询作用域](./role-data-scope-and-query-scopes)
