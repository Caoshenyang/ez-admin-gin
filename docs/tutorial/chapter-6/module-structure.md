---
title: 模块固定结构
description: "对齐当前最终版后端结构，讲清一个模块为什么由 dto、repository、service、handler、routes 这些文件组成。"
---

# 模块固定结构

当一个后台项目进入“长期维护和持续扩模块”的阶段后，最容易先乱掉的，通常不是某个接口逻辑，而是模块边界本身：

- 文件到底放哪里
- 哪一层负责参数校验
- 哪一层负责事务
- 哪一层负责路由聚合
- 哪些规则应该沉在模块内部，哪些应该放到平台层

这一页只做一件事：把当前仓库已经基本成型的最终版模块骨架固定下来。

::: tip 🎯 这一页的目标
看完后，你应该能直接回答：

- 一个新模块应该建哪些文件
- 每个文件各自承担什么职责
- 哪些公共工具已经沉在 `internal/app/` 供所有模块复用
:::

## 先看当前主线里的真实骨架

现在仓库里的核心模块，已经大体收敛到下面这种形态：

```text
server/internal/module/
├─ auth/
│  ├─ dto.go
│  ├─ repository.go
│  ├─ *_service.go
│  ├─ *_handler.go
│  └─ routes.go
├─ iam/
│  ├─ user/
│  ├─ role/
│  ├─ department/
│  ├─ post/
│  └─ menu/
└─ system/
   ├─ config/
   ├─ file/
   ├─ operationlog/
   ├─ loginlog/
   ├─ notice/
   └─ routes.go
```

其中最常见的模块目录会长成这样：

```text
module/<group>/<resource>/
├─ dto.go
├─ entity.go
├─ repository.go
├─ service.go
├─ handler.go
├─ routes.go
├─ policy.go        # 可选
└─ datascope.go     # 可选
```

## 每个文件到底负责什么

### `dto.go`

负责一切和输入输出结构直接相关的内容，例如：

- 请求体
- 查询参数
- 响应结构
- 基础归一化函数
- 参数合法性校验

当前很多模块里的 `NormalizeCreateRequest`、`NormalizePage`、`ValidStatus` 都放在这里。

它解决的是：

> 这个模块对外暴露什么协议，以及这些协议的最小合法边界是什么。

### `entity.go`

负责定义模块内部当前使用的实体别名，常见写法是直接复用 `internal/model/*`：

```go
type Entity = model.User
type UserRoleEntity = model.UserRole
```

这一步的意义不是“多绕一层”，而是给后续模块边界继续演进留出口。

现阶段很多模块还直接复用全局模型，但通过 `entity.go` 先收一层，后面如果真要把领域实体进一步独立出来，影响面会小很多。

### `repository.go`

负责持久化和查询拼装。

这里通常会放：

- 列表查询
- 按 ID 查资源
- 唯一性检查
- 可用性检查
- 关系替换
- 资源聚合查询

如果模块已经接入数据权限，`repository.go` 也会成为查询作用域真正落地的位置。

它解决的是：

> 这个模块怎么和数据库打交道，以及查询最终长什么样。

### `service.go`

负责业务规则和事务边界。

这一层通常会做：

- 调用 `dto.go` 的归一化逻辑
- 组织多步校验
- 开启事务
- 编排多个 repository 操作
- 补缓存同步、落盘流程或聚合结果

比如：

- 配置模块会在这里同步 Redis 缓存
- 文件模块会在这里处理上传落盘和元数据写库
- 角色模块会在这里统一更新数据范围、接口权限或菜单权限

它解决的是：

> 一次业务动作到底由哪些步骤组成，事务应该包到哪一层。

### `handler.go`

负责 HTTP 协议层。

这一层只做几件事：

- 从 Gin 绑定参数
- 从上下文拿当前用户或 `Actor`
- 调 service
- 把错误转成统一响应

当前主线里，`handler` 已经不再是“直接写所有业务逻辑”的地方，而是把协议层和业务层隔开。

### `routes.go`

负责路由注册和模块装配。

通常会做下面几件事：

1. 创建 repository
2. 创建 service
3. 创建 handler
4. 把 handler 方法挂到传入的路由组上

这意味着模块本身不直接依赖全局路由文件，而是通过 `RegisterRoutes(...)` 把自己作为一个可装配单元暴露出去。

### `policy.go`

可选文件，主要用于集中放权限码常量，例如：

```go
const (
	PermissionList = "system:user:list"
	PermissionCreate = "system:user:create"
)
```

它的价值在于把“这个模块有哪些权限点”明确收口在一个地方，而不是散落在前端、迁移种子和文档描述里。

### `datascope.go`

可选文件，用于声明这个资源如何接入数据权限。

例如：

- 用户资源用 `UserQueryScope`
- 部门资源用 `DepartmentQueryScope`
- 岗位资源当前显式放开，不做范围裁剪

它解决的是：

> 这个资源到底按什么字段接入平台层的数据权限规则。

## 为什么这套结构比旧的 Handler 聚合更稳

旧的扁平结构通常会把下面这些职责压在一个文件里：

- 参数绑定
- 业务判断
- 数据库查询
- 事务控制
- 路由注册

前期很快，但模块一多就会出现几个问题：

- 一个文件越来越长
- 规则越来越难复用
- 事务和查询边界混在一起
- 想接数据权限或缓存同步时很难下手

当前这套结构的核心价值，不是为了“像大项目那样分层”，而是为了把演进成本提前降下来。

## 当前不同类型模块各自长什么样

### IAM 模块

当前 IAM 资源基本都已经进入完整形态，例如：

- `module/iam/user`
- `module/iam/role`
- `module/iam/department`
- `module/iam/post`
- `module/iam/menu`

其中用户、部门、岗位会更容易和数据权限打交道，所以常见会带 `datascope.go`。

### System 模块

当前系统资源也基本都在最终结构里，例如：

- `module/system/config`
- `module/system/file`
- `module/system/operationlog`
- `module/system/loginlog`
- `module/system/notice`

这些模块虽然都在 `system` 分组下，但内部复杂度并不一样：

- 配置模块会处理 Redis 缓存同步
- 文件模块会处理上传校验和本地落盘
- 日志模块更偏查询与审计边界
- 公告模块则更接近轻量内容管理

这也是为什么“同样叫系统模块”，仍然需要每个模块有自己的 `service / repository / handler` 边界。

### Auth 模块

认证模块和普通 CRUD 模块又不完全一样。

当前 `module/auth` 更像一组围绕身份消费组织起来的能力：

- `login`
- `me`
- `menus`
- `dashboard`

所以它的 service / handler 文件是按能力拆开的，而不是一个资源一个 service。

这说明：

> 模块固定结构不是要求所有模块长得一模一样，而是要求职责拆分方式稳定且可预期。

## 路由装配为什么放在 `bootstrap/router.go`

当前主线里，模块并不是自己去修改全局路由，而是由 `bootstrap/router.go` 统一装配：

```text
bootstrap/router.go
  ↓
authModule.RegisterRoutes(...)
setupModule.RegisterRoutes(...)
systemModule.RegisterRoutes(...)
```

再由 `module/system/routes.go` 继续把 system 下的具体模块聚合进去。

这样做有两个直接好处：

- 全局入口清楚，知道系统到底装了哪些模块
- 每个模块保持可独立装配，不必直接依赖某个巨型路由文件

## 模块公共工具：`internal/app/`

当模块越来越多时，一些跨模块的通用逻辑会被提取到共享包里。当前 `internal/app/` 提供了三个这样的工具：

| 函数 | 作用 | 典型调用位置 |
| --- | --- | --- |
| `NormalizePage(page, pageSize)` | 分页参数归一化：页码最小 1，每页最小 10、最大 100 | 各模块 `dto.go` |
| `WriteError(c, err, fallbackMsg, log)` | 统一处理 handler 层错误响应，区分 `apperror` 和未知错误 | 各模块 `handler.go` |
| `CurrentActor(c, log)` | 从 `gin.Context` 提取当前登录用户，未登录时自动返回 401 | 需要 `Actor` 的 `handler.go` |
| `UintIDParam(c, param, label, log)` | 从路径参数中解析 `uint` 类型 ID，失败时返回 400 | 需要路径 ID 的 `handler.go` |

这些工具解决的共性问题是：

> 各模块 handler 里反复出现的"解析参数、取当前用户、写错误响应"样板代码，现在可以统一调用，不再每个模块各自实现一份。

举个例子，之前各模块 handler 里取当前登录用户时需要自己处理"未登录则返回 401"的逻辑，现在可以直接写：

```go
actor, ok := app.CurrentActor(c, log)
if !ok {
    return
}
```

而路径参数解析也可以简化为：

```go
id, ok := app.UintIDParam(c, "id", "配置 ID", log)
if !ok {
    return
}
```

::: tip 为什么放在 `internal/app/` 而不是每个模块各自实现
这些函数不是某个特定模块的业务逻辑，而是所有模块共享的协议层工具。把它们收到 `internal/app/` 里，既避免了每个模块重复实现，也让后续修改（比如统一调整错误格式）只需要改一个地方。
:::

## 一个新模块的最小文件清单

在当前最终版结构下，新增一个常规 CRUD 模块时，最小清单通常是：

- `dto.go`
- `entity.go`
- `repository.go`
- `service.go`
- `handler.go`
- `routes.go`

按需要再补：

- `policy.go`
- `datascope.go`

如果这个模块还依赖缓存、对象存储、外部服务等，再由 `RouteOptions` 或 `Service` 构造函数注入对应依赖。

## 本页最关键的结论

这一页真正要建立的判断是：

> 当前模块结构的核心，不是文件数量，而是把协议层、业务层、查询层、路由装配层和资源级规则层明确拆开。

只要这套边界稳定，后续无论继续补系统模块，还是新增真实业务模块，都能沿着同一条主线扩下去。

下一节继续看这套结构在后端接入时如何一步步落地：[后端模块接入流程](./backend-module-flow)。
