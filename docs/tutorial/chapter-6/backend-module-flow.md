---
title: 内置模块落地流程
description: "按当前最终版结构说明一个系统内置模块如何落地到 module/system 与 module/iam，从模型复用到 service、handler、routes 与系统聚合。"
---

# 内置模块落地流程

前一页已经把模块骨架定下来了。这一页继续往前走，把“一个新模块怎么真正接进当前后端”这件事拆成一条可以照着走的路径。

这条路径的重点不是“把文件建出来”，而是：

> 让新模块从一开始就沿当前最终版结构接入，而不是先写进旧目录，后面再返工迁移。

::: tip 🎯 本节目标
读完这一节，你应该能按当前主线顺序，把一个新模块接到后端：

1. 定义或复用底层模型
2. 在 `module/<group>/<resource>` 下补齐 `dto / entity / repository / service / handler / routes`
3. 通过 `module/system/routes.go` 或其他上层聚合模块把它接进系统
4. 让新接口自动复用认证、日志、接口权限和数据权限链路
:::

## 先看当前真实接入路径

现在一个系统模块真正进入后端，不再是”改某个全局路由文件再 new 一个 handler”，而是下面这条路径：

```text
internal/model/*
  ↓
module/<group>/<resource>/*
  ↓
module/<group>/routes.go
  ↓
bootstrap/router.go
  ↓
Gin Engine
```

以系统配置模块为例，实际是这样接起来的：

```text
model/system_config.go
  ↓
module/system/config/
  ├─ dto.go
  ├─ entity.go
  ├─ repository.go
  ├─ service.go
  ├─ handler.go
  └─ routes.go
  ↓
module/system/routes.go
  ↓
bootstrap/router.go
```

## Step 1：先判断模块应该落在哪个分组

在开始写代码之前，先回答一个问题：

> 这个模块属于哪个聚合边界？

当前主线里最常见的是三类：

| 分组 | 适合放什么 |
| --- | --- |
| `module/auth` | 登录后身份消费能力，如 `me / menus / dashboard` |
| `module/iam/*` | 用户、角色、部门、岗位、菜单这类身份与权限基础资源 |
| `module/system/*` | 配置、文件、日志、公告等系统支撑资源 |

如果你后面接一个真实业务模块，也建议优先照这个方式继续扩，比如：

```text
module/business/order/
module/business/customer/
```

而不是再塞回一个越来越大的 `system` 里。

## Step 2：先定义或复用底层模型

当前大多数模块仍然复用 `internal/model/*` 作为底层表结构定义。

这一步通常包含两件事：

1. 在 `internal/model/` 下新增或补齐模型。
2. 在模块内通过 `entity.go` 收一层别名。

例如用户模块：

```go
type Entity = model.User
type UserRoleEntity = model.UserRole
```

这样做的好处是：

- 先继续复用全局模型，避免一开始过度抽象
- 同时给后续模块边界继续收紧留出口

::: info 当前教程里的推荐顺序
首版先让表结构和模块能力跑通，再考虑是否要把领域实体和持久化模型进一步分开。对当前仓库阶段来说，`internal/model + module/*/entity.go` 是更稳的过渡方案。
:::

## Step 3：在模块目录里先补 `dto.go`

`dto.go` 是最适合先开始的地方，因为它会帮你先把模块协议边界定清楚。

当前主线里，`dto.go` 通常会包含：

- 列表查询参数
- 创建请求
- 更新请求
- 状态修改请求
- 对外响应结构
- `Normalize*` 函数
- `Valid*` 函数

一个很重要的判断标准是：

> 参数归一化、状态合法性和分页收口，优先放在 `dto.go`，不要散写到 Handler 各个方法里。

这样后面 `service.go` 只需要调用统一入口，而不用每个方法自己重新修参数。

## Step 4：在 `repository.go` 里收查询和持久化

接下来再补 `repository.go`。

这一层优先放下面几类能力：

- 列表查询
- 按 ID 查询
- 唯一性检查
- 可用性检查
- 创建、更新、状态修改
- 关系替换

如果模块接了数据权限，也是在这里真正把作用域应用到查询上。

例如用户模块当前就是从这里开始先收紧查询：

```go
queryDB := applyDataScope(r.db.Model(&model.User{}), actor)
```

这意味着权限过滤不是结果出来后再补判断，而是从查询起点就已经开始生效。

## Step 5：在 `service.go` 里定业务规则和事务边界

`service.go` 是当前最终版结构最关键的一层之一。

它通常负责：

- 调用 `Normalize*`
- 编排校验顺序
- 开事务
- 调多个 repository 方法
- 做缓存同步、落盘、聚合构造等额外业务动作

几个当前仓库里的典型例子：

- 配置模块：更新数据库后同步 Redis
- 文件模块：先校验上传文件，再落盘，再写库
- 角色模块：统一更新接口权限、菜单权限、自定义部门范围
- 部门模块：维护部门树变更和祖先路径更新

也就是说：

> 如果某个动作不只是“一条简单查询”，优先让它在 `service.go` 里成为明确的业务编排，而不是堆进 handler。

## Step 6：在 `handler.go` 里只保留协议层职责

当前主线里的 Handler，已经不再推荐直接承接复杂业务逻辑。

它更适合只做下面几件事：

1. `ShouldBindQuery` / `ShouldBindJSON`
2. 从上下文取 `CurrentUserID` 或 `CurrentActor`
3. 调 service
4. 用统一响应写回结果

像用户模块和部门模块现在都会在 Handler 层先取 `Actor`：

```go
actor, ok := middleware.CurrentActor(c)
```

但真正的数据权限应用，并不在 Handler 里写 SQL，而是交给 service 和 repository 继续往下走。

这一步的目标很明确：

> Handler 负责把 HTTP 请求翻译成模块调用，而不是成为“什么都塞”的总入口。

## Step 7：在 `routes.go` 里完成模块装配

当前每个模块都通过自己的 `RegisterRoutes(...)` 暴露装配入口。

这一步通常会：

1. 创建 repository
2. 创建 service
3. 创建 handler
4. 把路由挂到传入的 `*gin.RouterGroup`

例如系统配置模块当前就是：

```go
repo := NewRepository(opts.DB)
service := NewService(opts.DB, repo, opts.Redis, opts.Log)
handler := NewHandler(service, opts.Log)

group.GET("/configs", handler.List)
group.POST("/configs", handler.Create)
group.POST("/configs/:id/update", handler.Update)
group.POST("/configs/:id/status", handler.UpdateStatus)
group.GET("/configs/value/:key", handler.Value)
```

这说明一个模块真正对外暴露的，是“如何被挂进路由”，而不是要求上层自己知道它内部该怎么 new。

## Step 8：把模块接进聚合路由

写完模块自己的 `routes.go` 后，还差最后一步：把它接进聚合路由层。

对于系统模块，通常会接到：

- `module/system/routes.go`

再由它统一挂到 `/api/v1/system` 路由组。

这一步很关键，因为 `/api/v1/system` 当前已经统一挂着：

- `middleware.Auth`
- `middleware.LoadActor`
- `middleware.OperationLog`
- `middleware.Permission`

所以只要你把新模块正确接进这层，它就会自动复用认证、接口权限、操作日志和当前登录人上下文这条链路。

## Step 9：按模块性质决定是否补 `policy.go` 和 `datascope.go`

当基础 CRUD 跑通后，再判断这个模块是否需要补两个可选文件。

### `policy.go`

如果这个模块会稳定拥有独立权限点，就把权限码常量集中放在这里。

适合的模块包括：

- 用户
- 角色
- 菜单
- 部门
- 岗位
- 真实业务资源

### `datascope.go`

如果这个模块的数据天然会受当前登录人范围约束，就补 `datascope.go`，把接入方式显式写出来。

适合的模块包括：

- 用户
- 部门
- 以后真实业务里的“属于某部门 / 某负责人”的资源

这一步能让“这个资源到底怎么接数据权限”成为模块内一个很清楚的声明点。

## Step 10：别忘了迁移、权限和菜单入口

到这里模块代码本身已经能跑，但对后台系统来说，还不算真正“接完”。

一个完整模块通常还要继续补：

- 表结构迁移
- Casbin 接口权限种子
- 菜单与按钮种子
- 超级管理员默认绑定

这些内容会在本章后面的 [权限、菜单与迁移接入](./permission-menu-migration) 里继续展开。

## 当前主线下的最小接入清单

如果你现在要新增一个常规系统模块，推荐最小顺序就是：

1. 先定落点：`iam` 还是 `system`，还是新的业务分组
2. 定义 `internal/model/*`
3. 新建模块目录
4. 补 `dto.go`
5. 补 `entity.go`
6. 补 `repository.go`
7. 补 `service.go`
8. 补 `handler.go`
9. 补 `routes.go`
10. 接进聚合路由
11. 视情况补 `policy.go` / `datascope.go`
12. 再补迁移、权限和菜单种子

## 当前不该再走的旧路径

这一页最值得明确的一点是：

> 当前教程主线已经不再推荐”Model 写在 `internal/model`，逻辑全塞进一个巨型 handler 文件，最后去改一个全局大路由”这套方式。

原因不是它完全不能工作，而是它已经不适合当前仓库这条企业级完整版主线继续扩展。

## 本页最关键的结论

这一页真正要建立的判断是：

> 新模块接入后端时，最重要的不是先把某个接口写通，而是先让它正确进入当前最终版模块边界和系统聚合链路。

只要这一步做对，后面的权限、菜单、前端页面和真实业务扩展都会顺很多。

下一节继续看如何把这套结构放进一个完整系统模块示例里：[系统模块示例](./sample-module)。
