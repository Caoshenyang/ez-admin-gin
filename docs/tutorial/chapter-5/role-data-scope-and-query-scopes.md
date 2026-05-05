---
title: 角色数据范围与查询作用域
description: "按真实代码链路讲清 data_scope、Actor、多角色并集和 gorm.Scopes 如何把数据权限真正压进查询。"
---

# 角色数据范围与查询作用域

上一节先把组织模型的边界定稳了。这一节继续往前走，把数据权限真正从“表结构设计”推进到“请求执行规则”。

换句话说，这一节要回答的是：

> 角色上配好的 `data_scope`，到底是怎么在一次真实请求里，最终变成 SQL 过滤条件的。

::: tip 🎯 本节目标
读完这一节，你应该能顺着当前代码完整走通下面这条链路：

- 角色把数据范围配在 `sys_role` 和 `sys_role_data_scope`
- 请求经过 `LoadActor` 时把这些规则装进 `Actor`
- 平台层用 `datascope.Merge(...)` 合并多角色范围
- 模块通过 `gorm.Scopes(...)` 应用统一过滤
- 用户资源和部门资源各自使用不同的查询作用域
:::

::: info 这一页只负责“把总链路立起来”
这一页重点回答的是：

- `data_scope` 怎样一路走到查询过滤
- 第 5 章这条主线一共分成了哪几层

如果你已经理解总链路，接下来更值得分开读的是：

- [Actor 上下文与多角色并集](./actor-and-grant-merge)：专门看 `LoadActor`、`Grant`、`Summary` 和 `datascope.Merge(...)`
- [共享数据权限接入规范](./shared-datascope-integration-conventions)：专门看新模块该先选哪种 Scope
- [datascope.go 与 Repository 边界](./datascope-and-repository-boundary)：专门看模块内代码该落哪层
- [一次完整请求的权限过滤走读](./request-flow-walkthrough)：专门看真实请求怎样一路走到 Repository
:::

## 先看整条执行链路

当前主线里的数据权限，不是“某个接口里临时多拼一个 where”，而是下面这条完整流程：

```text
角色配置 data_scope
  ↓
sys_role_data_scope（仅 custom_dept 需要）
  ↓
middleware.LoadActor
  ↓
datascope.Actor
  ↓
datascope.Merge(...)
  ↓
UserQueryScope / DepartmentQueryScope
  ↓
module/iam/*/datascope.go
  ↓
Repository 查询结果被裁剪
```

这条链路里每一层都在解决不同问题：

| 层级 | 解决什么问题 |
| --- | --- |
| 角色配置层 | 角色本身允许看什么范围 |
| 请求上下文层 | 当前登录人此刻综合拥有哪些范围 |
| 平台规则层 | 多角色如何合并、超级管理员如何绕过 |
| 模块查询层 | 某类资源应该按哪个字段过滤 |

## 第一步：角色把数据范围配在哪里

当前项目里，角色数据范围不是额外挂在某个临时配置表上的，而是直接进入角色模型：

```go
type Role struct {
	ID        uint
	Code      string
	Name      string
	Sort      int
	DataScope datascope.Scope
	Status    RoleStatus
	Remark    string
}
```

也就是说，`sys_role.data_scope` 本身就是角色模型的一部分。

当前固定支持 5 档：

| 值 | 含义 |
| --- | --- |
| `all` | 全部数据 |
| `dept` | 本部门数据 |
| `dept_and_children` | 本部门及子部门数据 |
| `self` | 仅本人数据 |
| `custom_dept` | 自定义授权部门数据 |

当 `data_scope = custom_dept` 时，角色还会继续通过 `sys_role_data_scope` 绑定一组部门：

```text
sys_role
  ↓
data_scope = custom_dept
  ↓
sys_role_data_scope
  ↓
department_id 列表
```

这就是当前系统里“角色配置数据范围”的完整落点。

## 第二步：请求期为什么一定要先构建 `Actor`

如果只有角色配置，没有请求期装载，那么每个接口都得自己重新理解当前用户的角色、范围、自定义部门和超级管理员状态，这条链很快就会散。

所以当前主线在认证通过后，会先执行 `middleware.LoadActor`，把当前登录人的组织与数据权限上下文一次装进 Gin Context。

当前 `Actor` 的核心结构是：

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

这一层最想解决的是：后面所有模块都不要再自己重新查一次“当前人到底能看什么”。

`Actor` 里最关键的几类信息分别服务不同链路：

| 字段 | 给谁用 |
| --- | --- |
| `UserID` | 本人数据过滤、审计、业务归属 |
| `DepartmentID` | 本部门、本部门及子部门过滤 |
| `RoleCodes` | 接口权限链路继续复用 |
| `Grants` | 数据权限链路使用 |
| `IsSuperAdmin` | 平台层快速绕过数据权限 |

更完整的字段职责和多角色结构，下一页会单独展开：

- [Actor 上下文与多角色并集](./actor-and-grant-merge)

## 第三步：`LoadActor` 实际做了什么

当前 `middleware/actor.go` 里，`LoadActor` 的核心工作可以概括成 4 步：

1. 根据 `CurrentUserID` 查出当前用户基础信息。
2. 查出当前用户拥有的启用角色编码和 `data_scope`。
3. 如果某个角色是 `custom_dept`，再补查它被授予的部门 ID 列表。
4. 组装成 `datascope.Actor` 放进上下文。

这意味着当前主线里，数据权限装载不是“等某个业务模块需要时再查”，而是在统一中间件阶段就先准备好。

## 第四步：多角色为什么统一交给 `datascope.Merge(...)`

单角色很好理解，但企业后台里真正麻烦的是多角色。

比如一个用户可能同时拥有：

- 角色 A：`dept`
- 角色 B：`custom_dept`
- 角色 C：`self`

这时系统必须先回答一个问题：

> 最终应该按交集还是并集处理？

当前主线已经把这个答案固定为并集，并统一写进 `server/internal/platform/datascope/datascope.go`。

总规则可以先压成下面这张表：

| 情况 | 最终结果 |
| --- | --- |
| 任一角色是 `all` | 直接允许全部数据 |
| 同时有 `dept` 和 `custom_dept` | 两者取并集 |
| 有 `dept_and_children` | 允许本部门及整棵子树 |
| 有 `self` | 允许本人数据 |
| 当前用户是 `super_admin` | 直接绕过数据权限 |

这一步在总览页里只需要先记住一个结论：

- 平台层统一解释多角色
- 模块层不再自己发明并集或交集语义

更完整的 `Grant`、`Summary` 和 `Merge(...)` 细节，单独看：

- [Actor 上下文与多角色并集](./actor-and-grant-merge)

::: warning ⚠️ 不要把多角色合并规则散写到各模块
如果用户模块按并集，部门模块按交集，真实业务模块又自己定义第三种规则，后面几乎不可能排查“为什么这个人这里能看到，那里又看不到”。

当前这套主线的核心价值之一，就是把多角色语义提前锁死在平台层。
:::

## 第五步：为什么最终要落到 `gorm.Scopes(...)`

数据权限如果只停留在 `Actor` 里，还是不够。因为真正决定用户能看到什么的，不是内存结构，而是查询本身。

当前主线把数据权限真正压进查询的方式，是 `gorm.Scopes(...)`。

用户资源的接法现在是：

```go
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.UserQueryScope(db, actor, "department_id", "id"))
}
```

部门资源的接法现在是：

```go
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.DepartmentQueryScope(db, actor, "id"))
}
```

这样做的真正价值不是“语法更优雅”，而是三点：

- 过滤规则不会散落在 Handler 和 Service 里
- 同一类资源的列表、详情、更新前查询能复用同一规则
- 模块只需要声明“这个资源按哪个字段过滤”，不需要重写整套并集逻辑

至于模块内 `datascope.go` 和 Repository 为什么要这样配合，后面会单独展开：

- [datascope.go 与 Repository 边界](./datascope-and-repository-boundary)

## 用户资源为什么用 `UserQueryScope`

用户资源不是部门节点本身，而是“归属于某个部门，也可能退化为本人”的资源。

所以用户资源现在传给平台层的是两个关键字段：

| 参数 | 当前值 | 含义 |
| --- | --- | --- |
| `departmentColumn` | `"department_id"` | 用户归属部门字段 |
| `ownerColumn` | `"id"` | 用户本人标识字段 |

也就是说，用户资源当前采用的是：

- 能按部门范围看用户
- 也能在 `self` 情况下退化成只看自己

这就是为什么 `user/repository.go` 里的列表查询一开始就先走：

```go
queryDB := applyDataScope(r.db.Model(&model.User{}), actor)
```

后面再叠加关键词、状态、分页，只是在“已经被权限裁剪过的查询”上继续加工。

## 部门资源为什么要单独用 `DepartmentQueryScope`

部门资源和用户资源不一样。  
部门不是“属于部门的资源”，它自己就是部门节点。

所以部门模块当前的接法是：

```go
return db.Scopes(datascope.DepartmentQueryScope(db, actor, "id"))
```

这背后的判断是：

- 用户资源按“用户属于哪个部门”过滤
- 部门资源按“这个部门节点本身是否可见”过滤

当前实现里还有一个很实用的细节：

> 当角色范围是 `self` 时，部门资源会退化成“当前用户所在部门”。

这样至少还能看到自己的组织归属，不会出现“我能登录系统，但部门树整个空白”的体验。

## 当前用户模块和部门模块是怎么接进来的

### 用户模块

当前用户模块里，下面这些查询前置动作已经会复用统一数据范围：

- 用户列表
- 按范围查用户详情
- 更新用户前先查用户
- 修改用户状态前先查用户
- 绑定用户角色前先查用户

这意味着用户模块已经不是“只有列表做了权限过滤”，而是把“先查目标资源再执行后续动作”这一段整体收进了统一链路。

### 部门模块

部门模块现在同样会在查询和更新前先走范围约束：

- 查部门树
- 查当前可操作部门
- 修改部门前先确认节点在当前范围内
- 切换部门状态前先确认节点在当前范围内

这说明当前主线已经不只是在“属于部门的资源”上落数据权限，也开始把它应用到组织资源自身。

## 继续往下读时，最值得单独拆开的四层

如果你读到这里已经理解了总链路，接下来更顺的阅读方式通常是：

| 页面 | 主要解决什么问题 |
| --- | --- |
| [Actor 上下文与多角色并集](./actor-and-grant-merge) | `LoadActor` 到底装了什么，多角色为什么固定按并集 |
| [共享数据权限接入规范](./shared-datascope-integration-conventions) | 一个新模块更接近 `user`、`department` 还是 `post` |
| [datascope.go 与 Repository 边界](./datascope-and-repository-boundary) | 模块里的数据权限代码到底该落哪层 |
| [一次完整请求的权限过滤走读](./request-flow-walkthrough) | 真实请求怎样一路经过中间件、Service 和 Repository |

这样读的好处是：

- 先建立地图
- 再分层拆开
- 最后再回到真实请求

## 怎么验证这一节已经真正成立

### 1. `/api/v1/auth/me` 已经能返回数据范围摘要

登录后请求：

```text
GET /api/v1/auth/me
```

至少应该能看到：

- `department_id`
- `role_codes`
- `is_super_admin`
- `data_scope`

这说明 `LoadActor` 和 `BuildMeResponse(...)` 已经把当前登录人的数据范围压成了稳定响应。

### 2. 角色列表接口已经返回 `data_scope`

请求：

```text
GET /api/v1/system/roles
```

响应项里应该能看到：

- `data_scope`
- `custom_department_ids`

这说明角色配置层已经具备数据范围表达能力。

### 3. 用户列表已经先按范围裁剪，再做分页和筛选

请求：

```text
GET /api/v1/system/users
```

现在应该先经过统一数据范围过滤，再继续叠加：

- `keyword`
- `status`
- `page`
- `page_size`

这个顺序很重要，因为它说明权限不是结果集出来后再补判断，而是从查询源头就已经收紧。

### 4. 部门树也会被当前范围裁剪

请求：

```text
GET /api/v1/system/departments
```

如果当前登录人不是 `all` 范围，那么返回的就不应该是全量部门树，而是当前角色允许可见的那一部分。

如果你想把第 3、4 步再走得更细，最值得继续看的不是在这页里反复读摘要，而是直接进入：

- [一次完整请求的权限过滤走读](./request-flow-walkthrough)

## 本节最关键的结论

这一节真正要建立的判断是：

> 数据权限不是“在 SQL 后面追加一点条件”，而是一条从角色配置、到请求上下文、到平台合并规则、再到查询作用域的完整执行链路。

只要这条链路立住，后续继续扩部门、岗位、通讯录、审批单、业务单据时，就不需要每个模块重新发明一套“当前人能看什么”的实现。

下一步更推荐先看 [Actor 上下文与多角色并集](./actor-and-grant-merge)，把平台层语义收稳后，再回头看模块接法会更顺。
