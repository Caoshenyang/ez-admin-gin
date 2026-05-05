---
title: datascope.go 与 Repository 边界
description: "围绕 applyDataScope、List 和 FindByIDInScope 讲清模块内数据权限为什么要固定落在 datascope.go 与 Repository 入口，而不是散进 Handler 或 Service。"
---

# datascope.go 与 Repository 边界

第 5 章如果只讲到 `Actor`、`Merge(...)` 和资源类型判断，其实还差最后一个最容易被做散的地方：

> 模块里的数据权限到底应该落在哪一层，才能保证列表、详情、更新前检查都沿着同一条边界执行。

这一页专门把这个边界收清楚。

::: tip 🎯 本节目标
读完后，你应该能稳定回答三件事：

1. 为什么 `applyDataScope(...)` 应该固定收在模块自己的 `datascope.go`
2. 为什么 Repository 入口必须先套权限边界，再叠业务查询
3. 为什么 `FindByIDInScope(...)` 这类方法对更新、状态切换同样重要
:::

::: info 这一页默认你已经选好了资源模式
这一页不再展开“这个模块到底更像 `user`、`department` 还是 `post`”。

它默认你已经完成了那一步判断，只继续回答：

- 模块内代码该落哪层
- Repository 为什么必须成为真正的数据权限执行边界

如果你还停在“模式该怎么选”，先回看：

- [共享数据权限接入规范](./shared-datascope-integration-conventions)
:::

## 先看当前真实落点

当前第 5 章最值得对照的真实代码，不在 Handler，而在这几个文件：

```text
server/internal/module/iam/user/datascope.go
server/internal/module/iam/user/repository.go
server/internal/module/iam/department/datascope.go
server/internal/module/iam/department/repository.go
server/internal/module/iam/post/datascope.go
server/internal/module/iam/post/repository.go
```

这几个文件一起表达的是当前主线里的固定顺序：

```text
模块判断资源类型
  ↓
datascope.go 固定 applyDataScope(...)
  ↓
Repository 入口先套权限边界
  ↓
再叠 keyword / status / 排序 / 分页
  ↓
详情、更新前查询也走同一条边界
```

## 为什么 `datascope.go` 不该散进 Handler 或 Service

当前模块内单独保留 `datascope.go`，不是为了多拆一个文件，而是为了把“这个资源到底怎么接数据权限”固定成一个明确落点。

例如用户模块：

```go
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.UserQueryScope(db, actor, "department_id", "id"))
}
```

部门模块：

```go
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.DepartmentQueryScope(db, actor, "id"))
}
```

岗位模块：

```go
func applyDataScope(db *gorm.DB) *gorm.DB {
	return db
}
```

这说明 `datascope.go` 当前承担的是：

- 把模块和平台层规则接起来
- 用最短的形式声明“这个资源按什么维度过滤”

它不负责：

- 查询细节
- 业务规则
- 路径参数处理

## 为什么 Repository 必须先套权限边界

当前 Repository 最稳的写法，不是先查数据再补判断，而是一开始就先拿到已经被权限裁剪过的查询对象。

用户模块的列表入口就是：

```go
queryDB := applyDataScope(r.db.Model(&model.User{}), actor)
```

部门模块的列表入口也是同样顺序：

```go
queryDB := applyDataScope(r.db.Model(&model.Department{}), actor)
```

然后才继续往后追加：

- `keyword`
- `status`
- `Order`
- `Offset / Limit`

这条顺序的意义非常直接：

- 权限边界是查询源头的一部分
- 不是查询完成后的补丁

## 为什么这一步不能推迟到 Service

如果把数据权限推迟到 Service 层，最容易出现的三类问题是：

1. 列表过滤了，但详情没过滤
2. 更新动作只按 ID 查到了资源，却没检查当前人有没有资格看到它
3. 同一个模块里，不同 Repository 方法各自拼一套近似但不完全一致的规则

而当前这种“Repository 开头先统一 `applyDataScope(...)`”的写法，会让同类查询天然站在同一条边界上。

## `List(...)` 为什么只是第一层验证

很多人第一次接数据权限时，会默认把重点全放在列表上。

但企业后台里真正更危险的往往不是：

- 你能不能在列表里看到它

而是：

- 你能不能拿着一个 ID 继续更新它、禁用它、改它的角色关系

所以当前主线里，列表只是第一层：

- `List(...)` 负责把可见范围裁进列表

它还不够。

## 为什么还需要 `FindByIDInScope(...)`

用户模块当前已经专门提供了：

```go
func (r *Repository) FindByIDInScope(db *gorm.DB, actor datascope.Actor, userID uint) (model.User, error)
```

部门模块也有对应做法：

```go
func (r *Repository) FindByIDInScope(db *gorm.DB, actor datascope.Actor, departmentID uint) (model.Department, error)
```

这类方法的真正价值是：

- 让“按 ID 查单条资源”也站在同一条权限边界内

这样后续这些动作才能真正安全：

- 更新基础信息
- 修改状态
- 修改角色
- 调整组织关系

## 当前这条边界在用户模块里怎么成立

用户模块现在已经体现出一条比较完整的模式：

### 列表阶段

```go
queryDB := applyDataScope(r.db.Model(&model.User{}), actor)
```

先裁边界，再加：

- `username / nickname` 的关键字搜索
- 状态筛选
- 分页排序

### 单条资源阶段

```go
err := applyDataScope(db, actor).First(&user, userID).Error
```

这意味着：

- 当前人看不到的用户
- 在更新、改状态、改角色这些动作里同样查不出来

也就是说，当前模块已经不是“列表有权限，动作靠自觉”。

## 部门模块为什么同样要用这条边界

部门模块虽然是树资源，但原理一样。

它的 `FindByIDInScope(...)` 同样是在说：

- 这个部门节点如果不在当前人的可见范围里
- 那么后续编辑和状态切换也不应该继续进行

这让“列表能不能看见”和“能不能继续操作”统一到了一起。

## 岗位模块为什么现在还没有 `FindByIDInScope(...)`

岗位模块当前保留的是：

- 显式放开 `applyDataScope(db) -> db`

所以它的 Repository 里暂时仍然是：

- 普通 `FindByID(...)`

这不是少了一层设计，而是因为当前阶段岗位资源本身还没有组织范围边界可套。

这也再次说明：

- Repository 边界必须跟资源类型一致
- 不是每个模块都机械复制一份 `FindByIDInScope(...)`

## 当前模块内最稳的职责分工

如果把第 5 章这条边界压成最实用的分工表，可以直接记成下面这样：

| 位置 | 负责什么 |
| --- | --- |
| `platform/datascope` | 解释多角色并集和具体 Scope 语义 |
| `module/*/datascope.go` | 声明这个资源按哪个字段套 Scope |
| `repository.go` | 一开始就应用 `applyDataScope(...)`，并把列表与单条查询都收进同一边界 |
| `service.go` | 在已经完成资源可见性检查之后，再处理业务规则和事务 |
| `handler.go` | 处理 HTTP 入参和统一响应 |

这张表最想强调的是：

- `datascope.go` 不是规则定义层
- `repository.go` 才是规则真正落到查询的位置

## 这套边界对真实业务模块意味着什么

如果你后面要接一个真实业务模块，当前更稳的落法不是：

- 先写个列表能用就算了

而是一次性把下面两类入口都想清楚：

1. 列表查询怎么套 `applyDataScope(...)`
2. 按 ID 查单条资源时，是否也要补 `FindByIDInScope(...)`

只做第一条，不做第二条，后面通常会在更新动作里留下边界漏洞。

## 怎么判断这页已经真的收稳

读完这一页后，你至少应该能稳定回答下面四个问题：

1. 为什么 `datascope.go` 是模块自己的声明层，而不是平台层
2. 为什么 Repository 必须在开头就先套权限边界
3. 为什么 `List(...)` 不足以代表这个模块已经接好了数据权限
4. 为什么 `FindByIDInScope(...)` 对更新和状态切换同样重要

如果这四个判断已经稳住，那么当前模块内“代码该落哪层”的问题就已经定住了；接下来再回到请求主线，会更容易看清每一层为什么这样分工。

只要这四个判断已经稳住，第 5 章这条“资源级过滤规则”主线就算真正落地了。

## 下一步

- 想继续看模块该先选哪一种 Scope，读 [共享数据权限接入规范](./shared-datascope-integration-conventions)
- 想继续看资源类型本身怎么判断，读 [真实业务模块的数据权限边界](./business-module-datascope-boundaries)
- 想回到平台层规则，读 [Actor 上下文与多角色并集](./actor-and-grant-merge)
