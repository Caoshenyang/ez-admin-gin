---
title: 模块规范
description: "集中说明 EZ Admin Gin 当前后端模块的固定骨架、依赖装配方式、可选文件，以及新增模块时最稳的落地顺序。"
---

# 模块规范

这页不讲某一个具体模块，而是把当前后端模块的固定做法收成一份可复用规范。

如果你准备新增一个 `crm/customer`、`project/task` 之类的新模块，先看这页，能少走很多弯路。

## 当前模块的统一目标

当前模块规范想解决的不是“文件名漂不漂亮”，而是这三件事：

1. 路由装配有固定入口
2. 业务规则、查询逻辑、权限命名不再散落
3. 新模块可以按同一骨架落地，而不是每次临时商量

## 一个成熟模块的典型结构

当前已经比较完整的模块，通常会长成这样：

```text
server/internal/modules/<group>/<name>/
  api/
    dto.go
    handler.go
    routes.go
  application/
    ports.go
    service.go
  domain/
    types.go
  infra/
    repository.go
    datascope.go
  routes.go
  services.go
```

当前真实例子可以直接看：

- `server/internal/modules/iam/user`
- `server/internal/modules/iam/department`
- `server/internal/modules/system/config`
- `server/internal/modules/system/notice`

## 各文件职责

### `routes.go`

职责只有两个：

- 接收 `RouteOptions`
- 组装 `repo -> service -> handler`

典型形态：

```go
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/users", handler.List)
	group.POST("/users", handler.Create)
}
```

这说明当前模块路由层不是：

- 直接写业务逻辑
- 或把依赖构造散到 `bootstrap`

而是由模块自己完成最后一层装配。

### `api/handler.go`

`handler.go` 只负责 HTTP 边界：

- 解析路径参数
- 绑定请求体和查询参数
- 从中间件拿当前用户 / Actor
- 调用 Service
- 统一回写 `httpx.Success` / `httpx.Error`

更稳的原则是：

- Handler 不直接拼复杂 SQL
- Handler 不承担事务编排
- Handler 不保存业务规则

### `application/service.go`

`service.go` 是模块里的业务主心骨，当前通常负责：

- DTO 归一化后的业务校验
- 事务边界
- 聚合多个 Repository 操作
- 缓存同步
- 业务级别的错误抛出

例如 `server/internal/modules/system/config/application/service.go` 里就收了：

- 唯一键检查
- 创建 / 更新事务
- Redis 缓存同步

### `infra/repository.go`

`repository.go` 负责：

- 查询条件拼装
- 列表、详情、存在性查询
- 更新数据库记录
- 把 `gorm.ErrRecordNotFound` 转成明确业务错误

当前推荐把“数据怎么查”放在 Repository，把“为什么这么查”留在 Service。

### `api/dto.go`

`dto.go` 负责：

- 请求结构
- 响应结构
- 字段归一化
- 基础参数校验

当前模块里很常见的模式是：

- `NormalizeCreateRequest(...)`
- `NormalizeUpdateRequest(...)`
- `NormalizePage(...)`

也就是先把不干净的 HTTP 输入压成稳定的业务输入，再进入 Service。

### `domain/types.go`

`domain/types.go` 主要作用是把模块内的领域类型、枚举、常量和局部语义收起来。

如果模块只是在管理一张系统表，也可以在这里定义状态、菜单类型、排序规则等稳定领域概念，而不是把这些常量散落在 Handler 或 Repository 中。

### 权限常量

权限常量用来固定接口权限码，例如：

```go
const (
	PermissionList   = "system:user:list"
	PermissionCreate = "system:user:create"
)
```

它的职责是：

- 给 Casbin 提供稳定权限点
- 给角色授权和排障提供统一命名
- 避免权限码散落在 Handler、前端页面、初始化数据里

### `infra/datascope.go`

只有需要数据权限的模块，才应该补这个文件。

它的职责是：

- 把模块的资源过滤规则固定在一个地方
- 避免把 `gorm.Scopes(...)` 散落到多个查询函数里

例如用户模块当前就是：

```go
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.UserQueryScope(db, actor, "department_id", "id"))
}
```

## `RouteOptions` 是怎么判断该放哪些依赖的

当前每个模块都定义自己的 `RouteOptions`，而不是共用一个超级大对象。

这意味着：

- 用户模块只要 `DB + Log`
- 文件模块需要 `DB + Upload + Log`
- 配置模块需要 `DB + Redis + Log`
- 认证模块需要 `Config + DB + Redis + Token + Log`

更实用的判断标准是：

- 模块真正要用到什么，就在 `RouteOptions` 里声明什么
- 不要为了“统一”把没用到的依赖都塞进去

## 模块应该怎么接到总路由

当前总路由装配在：

- `server/internal/bootstrap/router.go`
- `server/internal/modules/iam/routes.go`
- `server/internal/modules/system/routes.go`
- `server/internal/modules/auth/routes.go`

真实顺序是：

```text
bootstrap/router.go
  ↓
modules/<group>/routes.go
  ↓
具体子模块 RegisterRoutes(...)
```

所以一个新模块要对外可用，不能只写完目录，还要把它注册进对应聚合路由。

## 当前模块命名上的固定约定

### 分组命名

当前后端分组大致遵循：

- `auth`：登录态与当前用户上下文
- `iam`：身份、角色、部门、岗位、菜单
- `system`：配置、文件、日志、公告这类系统能力

新业务模块如果不属于系统底座，更推荐独立成自己的业务分组，而不是继续塞进 `system`。

### 资源命名

模块目录名优先使用：

- 小写英文单词
- 必要时使用简短复合名

例如：

- `user`
- `role`
- `department`
- `operationlog`

### 接口路径命名

当前接口路径整体更接近“资源 + 动作后缀”的管理后台风格，例如：

- `GET /users`
- `POST /users`
- `POST /users/:id/update`
- `POST /users/:id/status`

这说明当前仓库并没有强推纯 RESTful，而是优先：

- 后台系统可读性
- 动作语义直观
- 前后端统一成本低

## 什么情况下可以不补全整套文件

当前规范不是为了追求“每个模块八个文件”，所以可以按实际复杂度裁剪：

| 情况 | 处理方式 |
| --- | --- |
| 纯查询只读模块 | 可以保持较轻的 `application`，但仍建议保留业务层入口 |
| 没有数据权限需求 | 可以没有 `infra/datascope.go` |
| 没有独立权限点 | 可以不单独维护权限常量文件 |
| 初始化或聚合模块 | 可以按实际职责裁剪 `api/application/infra` |

但裁剪的底线是：

- 路由、处理、业务规则三层不要重新糊在一起

## 新增模块最稳的落地顺序

如果你要新增一个正式模块，当前更稳的顺序是：

1. 先确定模型和表结构
2. 再补 `domain/types.go`
3. 再补 `api/dto.go`
4. 再补 `infra/repository.go`
5. 再补 `application/service.go`
6. 再补 `api/handler.go` 和 `api/routes.go`
7. 再补模块根部 `routes.go` 与 `services.go`
8. 需要权限时同步权限码、菜单按钮、接口资源和 Casbin 策略
9. 需要数据权限时补 `infra/datascope.go`
10. 最后注册到聚合路由、菜单、Casbin 和前端页面

## 一个模块什么时候算“真的接好了”

当前判断标准至少包括：

1. 路由已经注册
2. Handler / Service / Repository 已经形成清晰边界
3. 统一响应和错误码已经对齐
4. 权限码、菜单按钮、接口资源和 Casbin 策略已经同步
5. 需要数据权限的资源已经通过仓储层或 `infra/datascope.go` 固定查询作用域
6. 前端页面、菜单、按钮权限已经接通

## 相关教程与参考页

- [模块开发](/backend/module-development)
- [权限码约定](./permission-code-conventions)
- [数据权限模型](./data-scope-model)
- [目录约定](./directory-conventions)
