---
title: 后端模块接入流程
description: "按第 8 章的接入主线，说明一个新模块如何从模型、service、handler、routes 一路进入系统聚合路由。"
---

# 后端模块接入流程

前一页已经把模块骨架定下来了。现在继续往前走，把“一个新模块怎么真正接进当前后端”拆成一条能照着做的顺序。

::: tip 🎯 本节目标
读完后，你应该能按当前主线顺序，把一个新模块接到后端：

1. 判断模块分组
2. 定义或复用模型
3. 补 `dto / repository / service / handler / routes`
4. 接进聚合路由
5. 让接口自动复用认证、权限和日志链路
:::

## Step 1：先判断模块应该落在哪个分组

在写任何代码前，先回答：

> 这个模块属于哪个聚合边界？

当前最常见的分组仍然是：

| 分组 | 适合放什么 |
| --- | --- |
| `module/auth` | 登录后身份消费能力 |
| `module/iam/*` | 用户、角色、部门、岗位、菜单 |
| `module/system/*` | 配置、文件、日志、公告等系统支撑资源 |
| 新业务分组 | `crm/*`、`project/*` 这类真实业务模块 |

## Step 2：先定模型，再定模块目录

当前多数模块仍然会复用 `internal/model/*` 作为底层表结构定义。

比较稳的做法是：

1. 在 `internal/model/` 下补或复用模型
2. 在模块内用 `entity.go` 做一层局部别名

这样做的好处是：

- 先保证表结构与接口能力能一起跑通
- 同时给后续模块边界继续收紧留出余地

## Step 3：先补 `dto.go`

`dto.go` 是最适合先开始的地方，因为它会先把模块协议边界定稳。

一个常规模块当前至少建议先补：

- `ListQuery`
- `CreateRequest`
- `UpdateRequest`
- `Response`
- `ListResponse`
- `NormalizePage(...)`
- `NormalizeCreateRequest(...)`
- `NormalizeUpdateRequest(...)`

也就是说，先确定：

- 请求长什么样
- 响应长什么样
- 参数怎么收口

## Step 4：再补 `repository.go`

Repository 这一层优先收下面几类能力：

- 列表查询
- 按 ID 查询
- 唯一性检查
- 创建
- 更新
- 状态修改

如果模块要接数据权限，这一层还会变成：

- 查询作用域真正落地的位置

## Step 5：在 `service.go` 里定业务规则和事务边界

`service.go` 负责的，不是简单转发，而是：

- 调用 `Normalize*`
- 编排校验顺序
- 开事务
- 协调多个 repository 动作
- 处理缓存同步、落盘、关系替换等额外业务动作

可以直接记一句：

> 只要一个动作不再是一条简单 SQL，它就更应该先进入 `service.go`。

## Step 6：`handler.go` 只保留协议层

当前主线里，Handler 最稳的职责仍然是：

1. 绑定参数
2. 从上下文取 `CurrentUserID` 或 `CurrentActor`
3. 调用 service
4. 写统一响应

这一步最值得坚持，因为它直接决定后面模块还能不能稳接：

- 操作日志
- 接口权限
- 数据权限
- 事务边界

## Step 7：在 `routes.go` 里完成模块装配

每个模块都通过自己的 `RegisterRoutes(...)` 暴露装配入口。

这一步通常会：

1. 创建 repository
2. 创建 service
3. 创建 handler
4. 把路由挂到传入的 `*gin.RouterGroup`

这意味着模块真正对外暴露的，不是内部细节，而是：

- 它怎样被挂进系统

## Step 8：把模块接进聚合路由

模块自己的 `routes.go` 写完后，还要接进聚合路由层。

例如系统模块通常接到：

- `module/system/routes.go`

而真实业务模块也应该有自己的上层 `routes.go` 聚合，而不是回退到全局散装注册。

## Step 9：按模块性质决定是否补 `policy.go` 和 `datascope.go`

### `policy.go`

如果模块有稳定权限点，就补。

### `datascope.go`

如果模块的数据会受当前登录人的组织范围约束，就补。

这一步的价值在于：

- 权限和数据范围不会变成“最后才补的边角料”

## Step 10：别忘了它还没真正进入系统

到这里，模块代码本身能跑了，但还不算“真正进入后台系统”。

它后面还必须继续接：

- 菜单
- 按钮
- Casbin 策略
- 前端页面
- 手工验收

这些内容下一页继续展开：[权限、菜单与迁移接入](./permission-menu-integration)。
