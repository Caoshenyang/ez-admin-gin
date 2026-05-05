---
title: 一次完整请求的权限过滤走读
description: "以 GET /api/v1/system/users 为例，把 Auth、LoadActor、Permission、Handler、Service、Repository 和 datascope.go 串成一次完整的数据权限执行过程。"
---

# 一次完整请求的权限过滤走读

第 5 章前面几页已经把“规则是什么”“模块该怎么接”“边界该落哪层”拆开讲了。  
这一页不再补新概念，只做一件事：

> 把一次真实请求从进路由到出结果完整走一遍，让前面那些分散的点真正连成一条链。

::: tip 🎯 本节目标
读完后，你应该能顺着当前真实代码走通下面这条请求：

- `GET /api/v1/system/users?page=1&page_size=10`

并看清：

1. 登录态校验在哪一层完成
2. `Actor` 在哪一层装载
3. 接口权限和数据权限分别在哪一层生效
4. 为什么最后真正裁剪结果的是 Repository，而不是 Handler
:::

::: info 这一页只负责“把前面几层重新串起来”
如果你已经知道：

- 模块该选哪种 Scope
- `datascope.go` 和 Repository 该怎么分工

那么这一页的作用就只剩一个：

- 把前面拆开的几层重新还原成一次真实请求

所以读到这里时，不要把它当成新规则页，而要把它当成第 5 章的主线串讲页。
:::

::: tip 这页最适合放在最后读
如果你前面已经大致知道：

- `Actor` 里装了什么
- 模块该选哪种 Scope
- `datascope.go` 和 Repository 怎么分工

那这一页就是第 5 章最适合拿来“收口”的一页。它的作用不是补新知识，而是帮你确认前面的分层已经连成一条真实请求。
:::

## 先看这次走读选哪条请求

这一页用用户列表接口做样本，因为它同时具备：

- 路由保护
- `Actor` 装载
- Casbin 接口权限
- `UserQueryScope(...)`
- 列表分页和筛选

对应接口是：

```text
GET /api/v1/system/users
```

## 第 1 步：请求先进入系统模块总路由

当前系统模块总入口在：

- `server/internal/module/system/routes.go`

用户列表接口不会直接裸挂在 Gin 根路由下，而是先进入：

```text
/api/v1/system
```

这组路由会统一串上 4 个中间件：

1. `middleware.Auth`
2. `middleware.LoadActor`
3. `middleware.OperationLog`
4. `middleware.Permission`

也就是说，用户模块还没开始处理业务之前，登录态、当前用户上下文、操作审计和接口权限就已经先被框进来了。

## 第 2 步：`Auth` 先把“当前是谁”放进上下文

当前第一层是：

- `server/internal/middleware/auth.go`

它负责：

1. 读取 `Authorization: Bearer <token>`
2. 校验 Access Token
3. 把 `current_user_id` 和 `current_username` 写进 Gin Context

经过这一层后，系统已经知道：

- 你是不是已登录
- 当前用户 ID 是谁
- 当前用户名是谁

但这时还不知道：

- 你属于哪个部门
- 你有哪些角色
- 你的数据范围是什么

这些要留给下一层。

## 第 3 步：`LoadActor` 再把“当前人能看什么”装进去

第二层是：

- `server/internal/middleware/actor.go`

它会基于 `CurrentUserID` 继续去查数据库，组装出完整的：

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

这一层真正做成的是：

1. 查当前用户基础信息
2. 查当前用户拥有的启用角色
3. 把 `data_scope` 压成 `Grant`
4. 如果角色是 `custom_dept`，补出授权部门列表
5. 如果角色编码里出现 `super_admin`，标记超级管理员

到这一步，请求上下文里才第一次真正具备：

- 接口权限要用的 `RoleCodes`
- 数据权限要用的 `Grants`

## 第 4 步：`Permission` 先拦接口，再放业务继续走

第四层是：

- `server/internal/middleware/permission.go`

它会优先尝试从 `CurrentActor(c)` 里读出 `RoleCodes`，然后拿：

- 当前请求路径
- 当前 HTTP 方法

去做 Casbin 判断。

这一步要特别分清：

| 权限类型 | 当前看什么 |
| --- | --- |
| 接口权限 | `RoleCodes + path + method` |
| 数据权限 | `Actor.Grants + datascope.Merge(...) + Scope` |

也就是说，当前请求到了这里，只能说明：

- 这个人有权访问“用户列表接口”

还不能说明：

- 这个人有权看见“所有用户数据”

真正的数据裁剪还没开始。

## 第 5 步：用户模块 Handler 只负责把请求转进去

真正进入用户模块后，先到：

- `server/internal/module/iam/user/handler.go`

`List(c)` 的工作很克制，主要只有三件事：

1. 从上下文里取 `Actor`
2. 绑定查询参数 `ListQuery`
3. 调用 `service.List(actor, query)`

它并不在 Handler 里做这些事：

- 不自己查角色
- 不自己算并集
- 不自己拼 `WHERE department_id ...`

这正是当前第 5 章想稳住的边界。

## 第 6 步：Service 负责业务整形，不负责解释 Scope

进入：

- `server/internal/module/iam/user/service.go`

后，`List(actor, query)` 会先做：

```go
page, pageSize := NormalizePage(query.Page, query.PageSize)
```

然后把已经准备好的：

- `actor`
- `query`
- `page`
- `pageSize`

一起交给 Repository：

```go
users, total, err := s.repo.List(actor, query, page, pageSize)
```

也就是说，这一层负责的是：

- 分页参数归一化
- 结果聚合
- 角色 / 岗位信息的二次拼装

它也不负责重新定义：

- 多角色按并集还是交集
- 用户资源该按哪个字段过滤

这些都应该已经在更下层稳定下来。

## 第 7 步：真正裁剪结果的是 Repository 开头那一行

真正让结果变小的关键动作，发生在：

- `server/internal/module/iam/user/repository.go`

入口第一行就是：

```go
queryDB := applyDataScope(r.db.Model(&model.User{}), actor)
```

这一步的意义是：

- 从查询源头开始，先拿到“已经被权限裁剪过”的查询对象

后面再叠加：

- `keyword`
- `status`
- `Count`
- `Order`
- `Offset / Limit`

换句话说，当前顺序是：

```text
先按数据权限裁边界
  ↓
再按业务条件缩结果
  ↓
最后做分页
```

这和“先查出一堆数据，再在内存里挑一部分返回”完全不是一回事。

## 第 8 步：`applyDataScope(...)` 其实只是模块自己的声明层

用户 Repository 里那一行之所以能工作，是因为模块本身已经在：

- `server/internal/module/iam/user/datascope.go`

固定好了：

```go
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.UserQueryScope(db, actor, "department_id", "id"))
}
```

这一层只声明两件事：

| 参数 | 当前值 | 意思 |
| --- | --- | --- |
| `departmentColumn` | `department_id` | 用户归属部门字段 |
| `ownerColumn` | `id` | 用户本人标识字段 |

它不是在重新发明规则，而是在告诉平台层：

- 用户资源该按哪两个字段解释 `dept` / `self`

这一层如果你还想看得更细，应该回到：

- [datascope.go 与 Repository 边界](./datascope-and-repository-boundary)

## 第 9 步：平台层才真正解释并集和 Scope 语义

`applyDataScope(...)` 往下再走，最终会落到：

- `server/internal/platform/datascope/datascope.go`

这里才是当前第 5 章真正统一规则的地方：

1. `datascope.Merge(...)` 先把多角色范围压成 `Summary`
2. `UserQueryScope(...)` 再根据：
   - `AllowAll`
   - `IncludeDepartment`
   - `IncludeDeptTree`
   - `CustomDepartmentIDs`
   - `RequireSelf`
   生成真正的 GORM 过滤条件

所以在这条请求里：

- `Actor` 决定当前有哪些授权事实
- `Merge(...)` 决定这些授权怎么按并集汇总
- `UserQueryScope(...)` 决定这些汇总规则最终怎么变成查询过滤

如果你在这一步卡住，优先回看的是：

- [Actor 上下文与多角色并集](./actor-and-grant-merge)

## 第 10 步：为什么更新动作也必须走同一条链

如果第 5 章只盯着列表，其实还不够。

用户模块的更新、状态切换、角色修改，当前都会先调用：

- `FindByIDInScope(...)`

也就是：

```go
err := applyDataScope(db, actor).First(&user, userID).Error
```

这意味着：

- 当前人列表里看不到的用户
- 在更新和改状态时同样查不出来

这才算真正把“能不能看到”和“能不能继续操作”压到了同一条边界上。

## 这条请求走读真正想说明什么

把整条链压缩起来，当前用户列表请求的执行顺序就是：

```text
Auth
  ↓
LoadActor
  ↓
Permission
  ↓
user.Handler.List
  ↓
user.Service.List
  ↓
user.Repository.List
  ↓
user.applyDataScope(...)
  ↓
datascope.UserQueryScope(...)
  ↓
数据库返回已经被裁剪后的用户列表
```

这里最值得记住的不是文件名，而是职责边界：

| 层 | 负责什么 |
| --- | --- |
| `Auth` | 你是谁 |
| `LoadActor` | 你拥有哪些范围事实 |
| `Permission` | 你能不能访问这个接口 |
| `datascope.go` | 这个资源按哪个字段接数据权限 |
| `Repository` | 把权限边界压进查询 |

## 读完这页后你应该能稳定判断的事

到这里，如果第 5 章主线已经跟上，你应该能自然回答下面四个问题：

1. 为什么接口权限和数据权限虽然都依赖角色，但一定要分层
2. 为什么数据权限真正生效的位置在 Repository，而不是 Handler
3. 为什么 `datascope.go` 只是资源声明层，不是规则定义层
4. 为什么更新动作也必须复用同一条 `applyDataScope(...)` 边界

只要这四点已经稳住，第 5 章这条执行主线就基本成形了。

## 这页和前面几页怎么分工

如果你读到这里开始觉得“前面几页是不是都在讲同一件事”，可以直接用下面这张表区分：

| 页面 | 更适合解决什么问题 |
| --- | --- |
| [角色数据范围与查询作用域](./role-data-scope-and-query-scopes) | 整条链路总览 |
| [Actor 上下文与多角色并集](./actor-and-grant-merge) | `LoadActor` 和 `Merge(...)` 为什么这样设计 |
| [共享数据权限接入规范](./shared-datascope-integration-conventions) | 新模块该先选哪种 Scope |
| [datascope.go 与 Repository 边界](./datascope-and-repository-boundary) | 模块内代码边界该怎样固定 |
| 当前这页 | 把上面几层真正串成一次请求执行过程 |

这样读会更顺，因为：

- 前几页在拆层解释
- 这一页在重新把它们拼回一条链

## 如果这条请求某一步看不懂，最该回看哪页

可以直接按下面的回看顺序定位：

| 卡住的位置 | 优先回看 |
| --- | --- |
| 为什么要先装 `Actor` | [Actor 上下文与多角色并集](./actor-and-grant-merge) |
| 为什么用户和部门不是同一条 Scope | [资源级数据权限接入模式](./module-datascope-patterns) |
| 为什么代码要落在 `datascope.go` 和 Repository | [datascope.go 与 Repository 边界](./datascope-and-repository-boundary) |
| 为什么真实业务模块不能机械照抄用户模块 | [共享数据权限接入规范](./shared-datascope-integration-conventions) |

## 下一步

- 想继续回到模块内固定边界，读 [datascope.go 与 Repository 边界](./datascope-and-repository-boundary)
- 想继续回到资源类型判断，读 [共享数据权限接入规范](./shared-datascope-integration-conventions)
- 想回到平台层上下文和并集合并，读 [Actor 上下文与多角色并集](./actor-and-grant-merge)
