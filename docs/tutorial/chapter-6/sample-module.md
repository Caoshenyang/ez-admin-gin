---
title: 系统模块示例
description: "用当前真实的公告模块走完一条完整接入链路，证明第 6 章定义的系统模块结构已经可以稳定落地。"
---

# 系统模块示例

前面两页已经把“模块应该长什么样”和“后端应该怎么接”讲清楚了。现在还差最后一步：

> 找一个仓库里已经真实存在、前后端都接通了的模块，把整条链路从头到尾串起来看一遍。

这一页继续用公告模块做示例，但不再按旧的扁平结构讲，而是直接对照当前真实代码：

- 后端：`server/internal/module/system/notice`
- 前端：`admin/src/pages/system/NoticeView.vue`
- 菜单入口：`/system/notices`

::: tip 🎯 本节目标
读完这一节，你应该能清楚看到一件事：

当前第 6 章讲的模块结构，不是未来计划，而是已经在仓库里落地并可运行的一套真实形态。
:::

## 为什么继续用公告模块做示例

公告模块很适合当“第一个完整模块”示例，因为它同时具备下面这些特点：

- 结构完整，但业务复杂度不高
- 有分页查询、关键词检索、状态切换、新增和编辑
- 后端已经进入最终版模块结构
- 前端页面、API、动态菜单和按钮权限都已经接通

也就是说，它既足够真实，又不会因为业务本身太复杂，把注意力从模块结构上带偏。

## 先看这条完整链路

当前公告模块已经形成下面这条完整路径：

```text
model/notice.go
  ↓
module/system/notice/
  ├─ dto.go
  ├─ entity.go
  ├─ repository.go
  ├─ service.go
  ├─ handler.go
  ├─ routes.go
  └─ policy.go
  ↓
module/system/routes.go
  ↓
/api/v1/system/notices
  ↓
admin/src/api/notice.ts
  ↓
admin/src/pages/system/NoticeView.vue
  ↓
dynamic-menu.ts
```

这正好把第 6 章前面刚定下来的模块结构，完整证明了一次。

## 后端目录到底长什么样

当前公告模块目录是：

```text
server/internal/module/system/notice/
├─ dto.go
├─ entity.go
├─ repository.go
├─ service.go
├─ handler.go
├─ routes.go
└─ policy.go
```

这里没有 `datascope.go`，原因也很符合当前主线：

- 公告模块是系统级轻量内容资源
- 现阶段没有接组织级数据权限
- 所以它不需要像用户、部门那样再声明资源级查询作用域

这一点本身就很有代表性，因为它说明：

> 模块结构是稳定的，但并不是每个模块都必须机械地拥有完全一样的文件集合。

## `dto.go` 在公告模块里承担了什么

公告模块的 `dto.go` 现在已经集中收了下面这些内容：

- `ListQuery`
- `CreateRequest`
- `UpdateRequest`
- `UpdateStatusRequest`
- `Response`
- `ListResponse`
- `NormalizePage`
- `NormalizeCreateRequest`
- `NormalizeUpdateRequest`
- `NormalizeStatusFilter`
- `ParseNoticeID`

也就是说，公告模块所有“请求长什么样、响应长什么样、参数怎么收口”的问题，都先在这里定清楚了。

这里有一个当前主线里很值得保留的习惯：

> 标题、状态、备注这些字段的合法性检查，优先在 `dto.go` 里通过 `Normalize*` 收口，而不是散落到 `Create`、`Update`、`UpdateStatus` 各个 handler 里分别写一遍。

这样后面 `service.go` 只需要调用统一收口逻辑，不用反复处理同类字段。

## `repository.go` 现在负责哪些事

公告模块的 `repository.go` 目前承担的是典型的“查询和持久化”职责：

- 列表查询
- 按 ID 查公告
- 创建公告
- 更新公告基础字段
- 单独更新公告状态

列表查询这一步已经很有代表性：

```go
items, total, err := s.repo.List(query, page, pageSize, status)
```

而 `repo.List(...)` 内部再继续负责：

- 关键词过滤
- 状态过滤
- 总数统计
- 排序与分页

这说明当前结构里，分页列表的 SQL 形状已经不再塞在 Handler 里，而是明确沉到了 repository。

## `service.go` 在这个模块里为什么仍然有必要

有人会觉得公告模块业务不复杂，是不是可以不拆 service？

当前主线给出的答案是：仍然值得拆。

虽然公告模块没有像配置模块那样要同步 Redis，也没有像文件模块那样要处理落盘，但 `service.go` 仍然承担了两个很重要的职责：

- 调用 `Normalize*` 完成参数归一化
- 明确事务边界

例如创建公告时，当前流程是：

1. `dto.go` 校验并归一化请求
2. `service.go` 组装实体
3. `service.go` 用事务调用 `repo.Create(...)`
4. `service.go` 返回稳定响应对象

这让“业务动作的入口”始终稳定，不会因为当前模块简单，就把协议层和持久化层重新混回去。

## `handler.go` 现在只做什么

公告模块的 Handler 已经很接近当前教程主线想要的样子了：

- 绑定参数
- 解析路径 ID
- 调 service
- 统一写错误响应

比如更新公告时，Handler 只做这些：

1. `ParseNoticeID(c.Param("id"))`
2. `ShouldBindJSON(&req)`
3. `service.Update(noticeID, req)`
4. `response.Success(...)`

这里不会再直接写数据库查询或字段更新逻辑。

这正是当前最终版结构的目标：

> Handler 是协议层，不再是“顺手把业务也一起写掉”的地方。

## `routes.go` 怎么把模块装进系统

公告模块的 `routes.go` 很短，但很关键：

```go
repo := NewRepository(opts.DB)
service := NewService(opts.DB, repo)
handler := NewHandler(service, opts.Log)

group.GET("/notices", handler.List)
group.POST("/notices", handler.Create)
group.POST("/notices/:id/update", handler.Update)
group.POST("/notices/:id/status", handler.UpdateStatus)
```

这个文件解决的是：

- 模块内部依赖怎么构造
- 对外到底暴露哪些接口

上层不需要知道公告模块内部有哪些文件，只需要调用 `notice.RegisterRoutes(...)`。

## `policy.go` 在这个模块里有什么用

公告模块虽然是轻量模块，但它已经把按钮权限常量收在 `policy.go`：

```go
const (
	PermissionList         = "system:notice:list"
	PermissionCreate       = "system:notice:create"
	PermissionUpdate       = "system:notice:update"
	PermissionUpdateStatus = "system:notice:status"
)
```

这一步的价值，不是为了代码行数，而是为了让“这个模块到底有哪些权限点”有一个稳定落点。

后面无论是：

- 写菜单按钮种子
- 写前端 `canUse(...)`
- 写文档

都会更容易保持一致。

## 上层系统路由是怎么把它接进去的

公告模块不是自己去改全局 Gin 引擎，而是由 `module/system/routes.go` 统一聚合进去。

也就是说，真正的链路是：

```text
bootstrap/router.go
  ↓
module/system/routes.go
  ↓
module/system/notice/routes.go
```

这样一来，公告模块天然复用了 `/api/v1/system` 这组统一中间件：

- `middleware.Auth`
- `middleware.LoadActor`
- `middleware.OperationLog`
- `middleware.Permission`

这就是为什么公告接口一接进 `system` 分组，就自动具备了登录校验、接口权限和审计能力。

## 前端这一层是怎么接上的

当前公告模块对应的前端文件主要有三处：

```text
admin/src/api/notice.ts
admin/src/pages/system/NoticeView.vue
admin/src/router/dynamic-menu.ts
```

这三处分别承担：

| 文件 | 职责 |
| --- | --- |
| `api/notice.ts` | 封装 `/system/notices` 相关接口 |
| `pages/system/NoticeView.vue` | 公告管理页面 |
| `dynamic-menu.ts` | 把 `system/NoticeView` 映射成实际页面组件，并收集按钮权限码 |

这里很重要的一点是：前端页面并不是“手写路由进页面”，而是通过后端菜单树动态接进来的。

## 按钮权限在页面里怎么落地

公告页已经真实用了按钮权限控制，例如：

- `system:notice:create`
- `system:notice:update`
- `system:notice:status`

页面里通过：

```ts
buttonPermissionCodes.value.includes(code)
```

判断当前按钮是否显示。

这意味着当前公告模块不只是“接口能调”，而是：

- 后端菜单树返回了按钮节点
- 前端收集了按钮权限码
- 页面层已经在真实使用这些权限点

这正好把第 3 章菜单权限和第 6 章系统模块结构串起来了。

## 为什么公告模块是个很好的模板

公告模块之所以值得当模板，不是因为它最复杂，而是因为它已经覆盖了当前最终版结构里最值得学习的一组最小闭环：

- 有模型复用
- 有 `dto / repository / service / handler / routes`
- 有权限码常量
- 有系统路由聚合
- 有前端 API
- 有前端页面
- 有按钮权限消费

如果后续你要再新增一个轻量系统模块，例如：

- 数据字典
- 帮助中心
- 站内消息模板
- 常用文案管理

基本都可以先按公告模块这条路径来套。

## 怎么验证这一节不是“只会讲结构”

### 1. 接口已经真实可用

当前公告模块对外接口是：

- `GET /api/v1/system/notices`
- `POST /api/v1/system/notices`
- `POST /api/v1/system/notices/:id/update`
- `POST /api/v1/system/notices/:id/status`

这说明它不是概念模块，而是当前系统内置能力。

### 2. 菜单已经真实接入

后端种子里已经存在：

- 菜单：`system:notice`
- 按钮：`system:notice:list`
- 按钮：`system:notice:create`
- 按钮：`system:notice:update`
- 按钮：`system:notice:status`

### 3. 页面已经真实消费权限

前端 `NoticeView.vue` 中已经基于 `canUse(...)` 控制按钮显隐，这说明权限链路已经不是“未来计划”。

## 本节最关键的结论

这一节真正要建立的判断是：

> 当前仓库里的最终版模块结构，已经不是纸面约定，而是已经通过公告模块被完整跑通的一套真实落地方式。

只要你后面新增模块时，沿着公告模块这条路径继续扩，基本不会偏离当前主线。

下一节继续看这套系统模块为什么还必须补权限、菜单和迁移入口：[权限、菜单与迁移接入](./permission-menu-migration)。
