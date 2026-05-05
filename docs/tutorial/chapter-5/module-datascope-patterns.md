---
title: 资源级数据权限接入模式
description: "围绕 user、department、post 三类现成资源模式，讲清不同资源为什么要用不同 Scope，以及它们各自代表什么边界。"
---

# 资源级数据权限接入模式

`Actor` 和多角色并集规则定稳后，第 5 章还缺一个很容易被混在一起的判断：

> 当前仓库里已经存在的三类资源模式，到底分别在表达什么边界？

这一页专门回答这个问题。

::: tip 🎯 本节目标
读完后，你应该能判断：

1. 为什么用户资源和部门资源不能共用同一条 Scope
2. 为什么岗位模块当前会显式保留“暂不裁剪”的落点
3. 为什么真实业务模块优先应该先往这三类现成模式里靠
:::

::: info 这一页只负责“先看清现成的三类资源模式”
这一页重点解决的是：

- `user / department / post` 这三类模式各在表达什么
- 它们为什么不能共用一套过滤语义

如果你要解决的是别的问题，优先这样分开读：

- 想看新模块该怎么选模式：读 [共享数据权限接入规范](./shared-datascope-integration-conventions)
- 想看模块内代码到底该落哪层：读 [datascope.go 与 Repository 边界](./datascope-and-repository-boundary)
- 想看真实请求怎样一路跑过这些层：读 [一次完整请求的权限过滤走读](./request-flow-walkthrough)
:::

## 先看当前真实落点

当前第 5 章里，资源级数据权限主要落在这三个文件：

```text
server/internal/module/iam/user/datascope.go
server/internal/module/iam/department/datascope.go
server/internal/module/iam/post/datascope.go
```

它们和平台层的关系可以先压成下面这一层：

```text
platform/datascope/datascope.go
  ↓
module/*/datascope.go
```

也就是说：

- 平台层定义“规则怎么解释”
- 模块层声明“这个资源按哪个字段套规则”

## 模块层真正要声明的是什么

模块层最关键的不是“再发明一套权限系统”，而是回答下面这个问题：

> 这个资源应该按什么维度过滤？

在当前主线里，大致分成三类。

| 资源类型 | 代表模块 | 当前接法 |
| --- | --- | --- |
| 属于某个部门、也可退化到本人 | `user` | `UserQueryScope(...)` |
| 资源本身就是部门节点 | `department` | `DepartmentQueryScope(...)` |
| 当前先不按部门范围裁剪 | `post` | 显式放开 |

如果你只是想先记住结论，可以直接记成一句话：

- `user` 看“属于哪个部门，也能不能退化到本人”
- `department` 看“这个节点本身在不在范围里”
- `post` 看“当前是不是还属于显式放开的系统级字典”

## 第一类：用户资源

用户资源并不是部门节点本身，而是：

- 用户归属于某个部门
- 在 `self` 场景下又可能只允许看自己

所以用户模块当前传给平台层的是两个字段：

| 参数 | 当前值 | 含义 |
| --- | --- | --- |
| `departmentColumn` | `"department_id"` | 用户归属部门字段 |
| `ownerColumn` | `"id"` | 用户本人标识字段 |

对应代码是：

```go
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.UserQueryScope(db, actor, "department_id", "id"))
}
```

### 这条模式实际表达了什么

它本质上在说：

- `dept` 时按用户的 `department_id` 过滤
- `dept_and_children` 时按当前部门子树过滤
- `custom_dept` 时按授权部门集合过滤
- `self` 时退化成只允许 `id = 当前用户`

所以用户资源天然就是：

> 部门范围 + 本人范围的组合型资源。

## 第二类：部门资源

部门资源和用户资源不一样。

部门自己就是部门节点，不是“属于部门的资源”。  
所以部门模块当前直接把部门主键字段传给平台层：

```go
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.DepartmentQueryScope(db, actor, "id"))
}
```

### 这条模式实际表达了什么

它在说：

- 这个部门节点本身是否在当前可见范围内
- `dept_and_children` 时是否属于当前部门子树
- `custom_dept` 时是否命中授权部门集合

并且当前还有一个很重要的退化规则：

> 当范围是 `self` 时，部门资源退化成当前用户所在部门。

这样部门管理页至少还能看到自己的组织归属，不会整页空白。

## 第三类：岗位资源

岗位模块当前也保留了 `datascope.go`，但内容是显式放开：

```go
func applyDataScope(db *gorm.DB) *gorm.DB {
	return db
}
```

这不是缺失，而是当前主线里的一个有意识决定：

- 岗位资源当前先按系统级组织字典处理
- 先不按部门范围裁剪
- 等以后真的引入“岗位归属部门”或更复杂组织结构时，再把这一层收紧

### 为什么显式放开比省略更好

因为显式落点会告诉后续开发者：

- 这里不是忘了接数据权限
- 而是当前阶段就是有意不裁剪

这会让后面扩展岗位模块时边界更清楚。

## 当前三种模块接法应该怎么理解

可以直接用下面这张表记忆：

| 模块 | 当前接法 | 为什么这么做 |
| --- | --- | --- |
| `user` | `UserQueryScope(db, actor, "department_id", "id")` | 用户既属于部门，也可能退化成本人资源 |
| `department` | `DepartmentQueryScope(db, actor, "id")` | 部门本身就是范围节点 |
| `post` | `return db` | 当前阶段先作为系统级组织字典放开 |

这页真正想交付的是一个更稳定的判断：

> 资源模式本身有边界，后续新模块优先应该先往现成模式里靠，而不是急着发明第四种语义。

如果你已经想继续往“新模块该怎么选”走，下一页更合适的是：

- [共享数据权限接入规范](./shared-datascope-integration-conventions)

如果你想继续把这个判断放到新业务模块里用，下一页最值得继续看：

- [真实业务模块的数据权限边界](./business-module-datascope-boundaries)
- [岗位资源的数据权限收紧时机](./post-datascope-tightening)

## 这套写法和 Java 常见方案怎么对照

如果你来自 Java，可以这样类比：

| 当前 Go 写法 | 常见 Java 对照思路 |
| --- | --- |
| `UserQueryScope / DepartmentQueryScope` | 不同资源类型各自对应的数据权限过滤器 |
| 模块内 `datascope.go` | 资源级数据权限声明点 |
| 显式放开 `return db` | 某类资源当前暂不纳入过滤的明确边界 |

最大的差异依然在于：

- Java 项目更常见框架自动织入
- 当前这套 Go 主线更偏显式分类和模块内固定落点

## 怎么验证这一节已经成立

### 1. 用户列表会先过权限，再过筛选

请求：

```text
GET /api/v1/system/users
```

现在应该先按当前 `Actor` 的数据范围裁剪，再继续叠加：

- `keyword`
- `status`
- `page`
- `page_size`

### 2. 部门树不会默认返回全量节点

请求：

```text
GET /api/v1/system/departments
```

如果当前不是 `all` 范围，那么返回结果应该已经是裁剪后的树，而不是系统全量部门。

### 3. 岗位模块当前明确不做部门范围裁剪

请求：

```text
GET /api/v1/system/posts
```

当前岗位列表仍然按系统级字典处理，这说明“暂不裁剪”已经是显式规则，而不是遗漏。

## 下一步

- 想看新模块怎样把这三类现成模式拿来用，读 [共享数据权限接入规范](./shared-datascope-integration-conventions)
- 想看这些模式进入模块后，代码应该落到哪层，读 [datascope.go 与 Repository 边界](./datascope-and-repository-boundary)
- 想继续看新业务模块到底该选哪种数据权限模式，读 [真实业务模块的数据权限边界](./business-module-datascope-boundaries)
- 想继续看岗位什么时候才值得真正收紧范围，读 [岗位资源的数据权限收紧时机](./post-datascope-tightening)
- 想回到平台层规则本身，读 [Actor 上下文与多角色并集](./actor-and-grant-merge)
- 想继续回到总览页看完整执行链路，读 [角色数据范围与查询作用域](./role-data-scope-and-query-scopes)
