---
title: 组织模型设计
description: "在企业级完整版主线里先把部门、岗位、用户归属和 Actor 上下文定稳，为后续数据权限落地打基础。"
---

# 组织模型设计

这一节先不急着写复杂查询，而是把组织体系的底座定稳。数据权限之所以常常越做越乱，往往不是过滤条件不会写，而是部门、岗位、用户归属和当前登录人上下文一开始就没有设计清楚。

::: tip 🎯 本节目标
完成后，项目会具备一套稳定的组织模型约定：有部门表、有岗位表、有用户组织归属，也有当前登录人 `Actor` 上下文。后续无论接用户管理、部门管理还是业务模块，都能共用同一条数据权限主线。
:::

## 本节会改什么

这一节会把下面这些位置作为最终结构里的固定基础：

```text
server/
├─ internal/
│  ├─ middleware/
│  │  └─ actor.go
│  ├─ model/
│  │  ├─ department.go
│  │  ├─ post.go
│  │  ├─ role.go
│  │  ├─ role_data_scope.go
│  │  ├─ user.go
│  │  └─ user_post.go
│  └─ platform/
│     └─ datascope/
│        └─ datascope.go
└─ migrations/
   ├─ mysql/
   └─ postgres/
```

| 位置 | 作用 |
| --- | --- |
| `internal/model/department.go` | 定义部门实体 |
| `internal/model/post.go` | 定义岗位实体 |
| `internal/model/user.go` | 给用户补 `department_id` |
| `internal/model/role.go` | 给角色补 `data_scope` |
| `internal/model/role_data_scope.go` | 表达角色到自定义部门范围的绑定 |
| `internal/middleware/actor.go` | 认证通过后装载当前登录人的组织与数据权限上下文 |
| `internal/platform/datascope/datascope.go` | 固化数据范围枚举、合并规则和查询作用域 |

## 为什么先做组织体系，再做数据权限

企业后台里的数据权限，真正依赖的不是一段 `WHERE` 条件，而是下面这些主数据关系：

- 用户属于哪个部门
- 用户挂了哪些岗位
- 用户有哪些角色
- 角色的数据范围是什么
- 自定义部门范围到底授给了哪些部门

只要这几个基础关系没有定稳，后面的过滤逻辑就一定会东一块、西一块地散在各个 Handler 里。

::: warning ⚠️ 不要把数据权限理解成“给角色多加一个字段就结束了”
`data_scope` 只是入口，不是全部。真正能支撑企业后台长期扩展的，是“组织模型 + Actor 上下文 + 查询作用域”这三件事一起成立。
:::

## 组织模型的最小闭环

这一轮先把组织体系控制在一个够用、但不会过度设计的范围里：

| 实体 | 关键字段 | 说明 |
| --- | --- | --- |
| `sys_department` | `parent_id`、`ancestors` | 表达部门树 |
| `sys_post` | `code`、`name`、`status` | 表达岗位 |
| `sys_user` | `department_id` | 表达用户主归属部门 |
| `sys_user_post` | `user_id`、`post_id` | 表达用户和岗位的多对多关系 |
| `sys_role` | `data_scope` | 表达角色的数据范围 |
| `sys_role_data_scope` | `role_id`、`department_id` | 表达角色的自定义部门授权 |

这里没有继续往下拆“用户主岗”“兼职岗”“岗位组”等更复杂模型，是因为当前主线先服务企业级后台底座的首版落地。等主线稳定后，再往更复杂的人事场景扩展会更稳。

## 部门树为什么用 `parent_id + ancestors`

部门树首版不依赖数据库专属能力，而是统一采用：

- `parent_id`
- `ancestors`

例如一条部门记录的 `ancestors` 是：

```text
0,1,3
```

这表示它位于根节点下面，祖先链路依次经过 `1`、`3`。

这样做有三个直接好处：

1. PostgreSQL 和 MySQL 都容易落地。
2. “本部门及子部门”查询更容易统一实现。
3. 读者跟着教程排查数据问题时，也更直观。

::: details 为什么这一步不直接上递归 CTE
递归 CTE 在 PostgreSQL 上很好用，但它不是这套教程的默认路径。当前底座要同时兼容 PostgreSQL 和 MySQL，也要尽量让结构本身容易理解、容易排查，所以首版优先选择 `parent_id + ancestors`。
:::

## 角色数据范围固定为 5 档

当前主线把角色数据范围直接固定成后台系统里最常见的 5 档：

| 值 | 含义 |
| --- | --- |
| `all` | 全部数据 |
| `dept` | 本部门数据 |
| `dept_and_children` | 本部门及子部门数据 |
| `self` | 仅本人数据 |
| `custom_dept` | 自定义授权部门数据 |

这一组枚举已经固化在 `internal/platform/datascope/datascope.go`，后续模块不需要再自己发明命名。

## Actor 上下文应该包含什么

认证通过后，系统不应该只知道“当前用户 ID 是多少”，还应该把后续查询真正要用到的信息一次装进上下文。

当前 `Actor` 至少包含：

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

其中最关键的是两部分：

- `RoleCodes`：继续给接口权限链路使用
- `Grants`：给数据权限链路使用

这就是为什么接口权限和数据权限虽然都依赖角色，但不应该混成一层。

## 多角色合并规则先定死

同一个用户可能同时拥有多个角色。首版规则直接固定为并集：

- 只要任一角色是 `all`，结果就是全部数据
- `dept` 和 `custom_dept` 同时存在时，最终结果是两者并集
- `super_admin` 永远绕过数据权限限制

这套规则的好处是：

- 语义稳定
- 容易和前端配置页对齐
- 容易在教程里讲清楚

## 这一节做完后怎么验证

本节完成后，至少应该能验证三件事。

### 1. `/api/v1/auth/me` 已经能返回组织与数据范围摘要

登录后访问：

```text
GET /api/v1/auth/me
```

应该能看到：

- 当前用户基础信息
- `department_id`
- `role_codes`
- `is_super_admin`
- `data_scope` 摘要

### 2. 组织相关模型已经进入迁移脚本

检查下面两个文件里是否已经包含部门、岗位、用户岗位和角色数据范围相关建表语句：

- `/server/migrations/postgres/000003_enterprise_foundation.up.sql`
- `/server/migrations/mysql/000003_enterprise_foundation.up.sql`

### 3. 用户资源已经开始接入数据权限作用域

当前主线已经让用户管理先接入了查询级数据权限。也就是说，后续当角色数据范围配置好后，用户列表和用户编辑类接口会开始按 `Actor` 上下文约束可见范围。

::: info 这一节的重点不是“功能很多”，而是“主线定稳”
组织体系和数据权限一旦底层关系设计清楚，后面的模块接入、按钮权限、查询过滤和真实业务扩展都会顺很多。反过来，如果这里草草带过，后面每加一个模块都要重新解释一次“为什么这样过滤”。
:::

## 下一步

组织模型定稳后，下一节就可以继续进入数据权限真正落查询链路的部分：角色数据范围、`Actor` 装载、`gorm.Scopes(...)` 和资源级过滤约定会一起串起来。
