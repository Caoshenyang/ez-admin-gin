---
title: 模块初始化模板
description: "把一个新模块从目录骨架、路由装配、权限常量、菜单种子到前端页面接入，压成一份可直接复用的初始化模板。"
---

# 模块初始化模板

这页专门服务一个很实际的需求：

> 如果现在要新开一个模块，我最少应该先补哪些文件、常量和接入点，才能不偏离当前最终结构？

::: tip 这页怎么用
这不是某个具体业务模块的教程，而是一份“新模块开工模板”。

适合在你准备新增：

- `project/customer`
- `project/task`
- `system/dict`

这类模块时直接对照执行。
:::

## 先记住当前最小主线

一个新模块从“目录存在”到“前后端真正可用”，当前最小主线是：

```text
建模块目录
  ↓
补 dto / repository / service / handler / routes
  ↓
补 policy.go
  ↓
注册到对应聚合路由
  ↓
补菜单 / 按钮 / Casbin 种子
  ↓
前端补 api / 页面 / dynamic-menu 映射
```

## 后端目录最小模板

当前一个成熟模块最常见的目录骨架是：

```text
server/internal/module/<group>/<name>/
├─ dto.go
├─ entity.go
├─ repository.go
├─ service.go
├─ handler.go
├─ routes.go
├─ policy.go
└─ datascope.go
```

但不是每个模块都必须一次补满全部文件。

### 默认必备

- `dto.go`
- `repository.go`
- `service.go`
- `handler.go`
- `routes.go`
- `policy.go`

### 视情况补充

- `entity.go`
  适合：想把模块主实体局部别名化或后续可能替换包装
- `datascope.go`
  适合：这个模块需要进入第 5 章的数据权限主线

## `routes.go` 最小模板

当前模块路由装配的稳定模板就是：

```go
func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
	repo := NewRepository(opts.DB)
	service := NewService(opts.DB, repo)
	handler := NewHandler(service, opts.Log)

	group.GET("/examples", handler.List)
	group.POST("/examples", handler.Create)
}
```

你真正要改的通常只有三类东西：

1. 模块依赖
2. 资源路径
3. 暴露哪些接口

## `policy.go` 最小模板

新模块一开始就应该把权限点固定下来，不要等页面写完再回头补。

当前最稳的模板是：

```go
const (
	PermissionList   = "project:customer:list"
	PermissionCreate = "project:customer:create"
	PermissionUpdate = "project:customer:update"
	PermissionDelete = "project:customer:delete"
)
```

命名优先保持：

```text
<group>:<resource>:<action>
```

## `dto.go` 最小模板

当前最值得先固定的是这几类结构：

- `ListQuery`
- `CreateRequest`
- `UpdateRequest`
- `Response`
- `ListResponse`
- `NormalizePage(...)`
- `NormalizeCreateRequest(...)`
- `NormalizeUpdateRequest(...)`

也就是说，一开始就把：

- 请求长什么样
- 响应长什么样
- 参数怎么归一化

先收稳，不要让 Handler 自己慢慢长出参数清洗逻辑。

## `repository.go` 最小模板

即使模块很简单，也建议先把这 4 个动作分出来：

- `List(...)`
- `FindByID(...)` 或 `FindByIDInScope(...)`
- `Create(...)`
- `Update(...)`

如果模块属于第 5 章主线，还要进一步判断：

- 是不是需要 `datascope.go`
- 单条查询是不是也应该走 `FindByIDInScope(...)`

## `datascope.go` 什么时候补

不是每个模块都需要一上来补 `datascope.go`。

当前最稳的判断顺序是：

1. 先看这个模块是不是要进入组织体系与数据权限主线
2. 再看它更接近 `user`、`department` 还是显式放开
3. 如果要接数据权限，再补 `datascope.go`

相关页：

- [资源级数据权限接入模式](./module-datascope-patterns)
- [共享数据权限接入规范](./shared-datascope-integration-conventions)

## 总路由接入点别忘了

模块目录写完后，还要把它接进对应聚合路由。

当前主要聚合点通常是：

- `server/internal/module/auth/routes.go`
- `server/internal/module/system/routes.go`
- 未来新增业务分组自己的 `routes.go`

否则模块虽然存在，但不会真正对外可用。

## 初始化数据最少要补什么

如果模块需要真正进入后台界面，通常至少还要补三类初始化数据：

1. Casbin 策略
2. 菜单节点
3. 按钮节点

最小思路是：

```text
policy.go
  ↓
casbin_rule
  ↓
sys_menu（目录 / 菜单 / 按钮）
  ↓
sys_role_menu
```

这一步没补齐，前端常见表现就是：

- 菜单不出现
- 页面能开但按钮都没权限
- 接口打过去直接 403

## 前端最小接入模板

当前前端最小需要这几步：

1. `admin/src/api/<resource>.ts`
2. `admin/src/types/<resource>.ts`
3. `admin/src/pages/.../<View>.vue`
4. `admin/src/router/dynamic-menu.ts` 里补 `routeComponentMap`

如果页面有按钮动作，还要继续对齐：

5. 页面里补 `canUse(code)`
6. 后端 `sys_menu` 里补按钮节点 `code`

## 一个最小检查表

新模块开工后，最值得顺着这张清单过一遍：

| 项 | 最低标准 |
| --- | --- |
| 后端目录 | `dto / repository / service / handler / routes / policy` 已有 |
| 路由 | 已注册到对应聚合路由 |
| 权限常量 | `policy.go` 已固定命名 |
| 初始化数据 | Casbin、菜单、按钮节点已补 |
| 前端页面 | 页面文件和 API 已接通 |
| 动态路由 | `component` 已命中前端白名单 |
| 按钮权限 | 页面 `canUse(...)` 与按钮节点 `code` 对齐 |
| 数据权限 | 若属于第 5 章主线，`datascope.go` 已判断清楚 |

## 当前最适合直接参考哪个真实模块

如果你想看一套现成样本，当前最适合直接对照：

- 公告模块：结构完整、复杂度适中
- 用户模块：结构完整，且带数据权限边界
- 文件模块：额外带上传配置和文件落盘

对应教程页：

- [第 6 章：示例业务模块](/tutorial/chapter-6/sample-module)
- [第 6 章：模块接入验收清单](/tutorial/chapter-6/module-integration-checklist)

## 和哪些参考页一起看最顺

- [模块规范](./module-conventions)
- [权限码约定](./permission-code-conventions)
- [动态菜单组件白名单](./dynamic-menu-component-reference)
- [初始化数据参考](./init-data-reference)
- [查询与分页约定](./query-and-pagination-conventions)
