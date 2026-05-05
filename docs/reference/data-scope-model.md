---
title: 数据权限模型
description: "集中说明当前底座的数据权限枚举、Actor 上下文、多角色并集规则，以及 user / department / post 三类资源的接入边界。"
---

# 数据权限模型

这页不讲完整教程步骤，只负责把当前底座的数据权限约定集中收起来，方便你在实现、排障或扩模块时快速查阅。

::: tip 🎯 这页解决什么
如果你只想快速确认下面这些问题，就直接查这一页：

- 当前固定支持哪几档 `data_scope`
- `Actor` 里到底装了什么
- 多角色为什么按并集
- 用户、部门、岗位三类资源分别怎么接入
:::

## 当前固定支持的 5 档范围

当前底座把角色数据范围固定成后台系统里最常见的五档：

| 值 | 含义 |
| --- | --- |
| `all` | 全部数据 |
| `dept` | 本部门数据 |
| `dept_and_children` | 本部门及子部门数据 |
| `self` | 仅本人数据 |
| `custom_dept` | 自定义授权部门数据 |

这组枚举定义在：

- `server/internal/platform/datascope/datascope.go`

## 数据权限为什么一定依赖组织体系

当前这套模型不是“角色多一个字段”就能成立，而是和组织模型一起配套：

| 结构 | 作用 |
| --- | --- |
| `sys_department` | 提供部门树与组织节点 |
| `sys_post` | 提供岗位字典 |
| `sys_user.department_id` | 表达用户主归属部门 |
| `sys_user_post` | 表达用户和岗位的多对多关系 |
| `sys_role.data_scope` | 表达角色的数据权限范围 |
| `sys_role_data_scope` | 表达角色的自定义部门授权 |

所以当前数据权限的真实前提是：

- 有组织节点
- 有用户归属
- 有角色范围
- 有自定义部门绑定

## 当前请求期上下文长什么样

认证通过后，系统会先装载当前登录人的数据权限上下文，而不是让每个模块自己临时去查角色。

当前 `Actor` 结构是：

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

字段分工可以直接这样理解：

| 字段 | 用途 |
| --- | --- |
| `UserID` | `self` 范围、本人与资源归属判断 |
| `Username` | 日志、响应摘要、排查 |
| `DepartmentID` | 本部门和子部门范围 |
| `RoleCodes` | 给 Casbin / 菜单 / 接口权限继续使用 |
| `Grants` | 给数据权限系统使用 |
| `IsSuperAdmin` | 快速绕过数据权限 |

## `Grant` 和 `Summary` 分别代表什么

平台层没有把完整角色结构塞进数据权限链路，而是先压成一条更轻的授权事实：

```go
type Grant struct {
	Scope         Scope
	DepartmentIDs []uint
}
```

多角色再继续被合并成一份摘要：

```go
type Summary struct {
	AllowAll            bool
	RequireSelf         bool
	IncludeDepartment   bool
	IncludeDeptTree     bool
	CustomDepartmentIDs []uint
}
```

可以把它们理解成：

| 结构 | 角色 |
| --- | --- |
| `Grant` | 单个角色授予给当前用户的一条数据权限 |
| `Summary` | 多角色合并后的最终可执行摘要 |

## 多角色为什么固定按并集

当前底座把多角色规则固定为：

- 并集

也就是说：

| 情况 | 最终结果 |
| --- | --- |
| 任一角色是 `all` | 直接允许全部数据 |
| 有 `dept` | 允许本部门数据 |
| 有 `dept_and_children` | 允许本部门整棵子树 |
| 有 `custom_dept` | 允许授权部门集合 |
| 有 `self` | 允许本人数据 |
| 当前用户是 `super_admin` | 直接绕过数据权限 |

这条规则统一收在：

- `datascope.Merge(...)`

这样做的核心价值是：

- 所有模块语义一致
- 超级管理员逻辑只有一份
- 自定义部门集合只需要去重一次

## 当前真实执行链路

当前主线里的数据权限，不是某个接口里临时拼一段 `WHERE`，而是下面这条完整链路：

```text
sys_role.data_scope / sys_role_data_scope
  ↓
middleware.LoadActor
  ↓
datascope.Actor + datascope.Merge(...)
  ↓
module/*/datascope.go
  ↓
gorm.Scopes(...)
  ↓
Repository 查询结果被裁剪
```

## 当前三类资源接法

当前底座已经给出了三种典型资源接法：

| 资源类型 | 当前接法 | 代表模块 |
| --- | --- | --- |
| 属于部门，也可能退化成本人 | `UserQueryScope(...)` | `user` |
| 资源本身就是范围节点 | `DepartmentQueryScope(...)` | `department` |
| 当前先显式放开 | `return db` | `post` |

### 用户资源

用户资源当前接法是：

```go
return db.Scopes(datascope.UserQueryScope(db, actor, "department_id", "id"))
```

适合：

- 用户
- 客户
- 工单
- 审批单
- 其他既有部门归属、又可能退化成“仅本人可见”的资源

### 部门资源

部门资源当前接法是：

```go
return db.Scopes(datascope.DepartmentQueryScope(db, actor, "id"))
```

适合：

- 部门树
- 区域树
- 分类树
- 其他本身就是范围节点的树形资源

### 岗位资源

岗位模块当前是显式放开：

```go
func applyDataScope(db *gorm.DB) *gorm.DB {
	return db
}
```

这表示当前岗位更接近：

- 系统级组织字典

而不是：

- 带稳定部门归属的资源

## 当前哪些模块已经接入，哪些还没有

当前已经明确接入数据权限主线的典型资源有：

- 用户管理
- 部门管理

当前显式保留落点、但暂时不按部门裁剪的资源有：

- 岗位管理

而像下面这些系统级资源，当前并不属于首批数据权限目标：

- 角色管理
- 菜单管理
- 系统配置
- 登录日志
- 操作日志

## 什么时候不要急着接数据权限

如果一个资源满足下面这些特征，当前更适合先显式放开，而不是仓促接 Scope：

- 它不是组织节点
- 它也没有稳定的部门归属字段
- 它更像系统级公共字典
- 它的组织边界还没定稳

这正是岗位模块当前的处理方式。

## 关键代码位置

如果你要从代码里快速定位当前数据权限实现，优先看下面这些文件：

- `server/internal/platform/datascope/datascope.go`
- `server/internal/middleware/actor.go`
- `server/internal/module/iam/user/datascope.go`
- `server/internal/module/iam/department/datascope.go`
- `server/internal/module/iam/post/datascope.go`

## 相关教程页

如果你要看完整讲解，而不是快速查阅，直接读第 5 章这几页：

- [角色数据范围与查询作用域](/tutorial/chapter-5/role-data-scope-and-query-scopes)
- [Actor 上下文与多角色并集](/tutorial/chapter-5/actor-and-grant-merge)
- [资源级数据权限接入模式](/tutorial/chapter-5/module-datascope-patterns)
- [真实业务模块的数据权限边界](/tutorial/chapter-5/business-module-datascope-boundaries)
- [岗位资源的数据权限收紧时机](/tutorial/chapter-5/post-datascope-tightening)
