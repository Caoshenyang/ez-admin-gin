---
title: 第 5 章：组织体系与数据权限
description: "把部门、岗位、角色数据范围、Actor 上下文和 gorm.Scopes 查询过滤串成一条企业级后台可复用主线。"
---

# 第 5 章：组织体系与数据权限

前面第 3、4 章已经把“谁能登录、谁能访问哪些接口”这条线收起来了。接下来要补的，是企业后台里另一个决定系统复杂度上限的能力：组织体系和数据权限。

如果认证和接口权限解决的是“你是谁、你能不能进”，那么这一章解决的就是：

> 你进来之后，到底能看到哪些数据。

::: tip 本章怎么读
这一章不要把它当成几张表的补充说明，而要把它当成一条完整执行链路来看：

- 角色里先配置 `data_scope`
- 登录后把角色范围装进 `Actor`
- Repository 再通过 `gorm.Scopes(...)` 把规则真正压进查询
:::

::: info 当前这一章已经闭环
第 5 章现在已经形成了一条完整主线，不再只是“还在补组织附录”：

- 先把组织模型和角色范围定稳
- 再把 `Actor`、`Grant` 和多角色并集讲清楚
- 接着把资源模式与模块接法落到 `datascope.go + Repository`
- 最后用真实请求走读和验收清单收回到代码执行
:::

## 本章会解决什么

这一章会把下面几个容易在企业后台里越做越乱的问题，一次性讲清楚：

- 为什么企业后台一定要先有部门、岗位和用户组织归属
- 为什么角色数据范围不能和接口权限混成一层
- 当前登录人上下文里为什么不能只放 `user_id`
- 多角色并集应该在哪里统一定义
- 为什么数据权限更适合通过 `gorm.Scopes(...)` 落在查询链路里
- 为什么“部门资源本身”和“属于部门的资源”需要不同作用域

## 这一章真正的主线

当前仓库里的数据权限主线，已经不是停留在概念图上的设计，而是下面这条真实链路：

```text
sys_role.data_scope / sys_role_data_scope
  ↓
middleware.LoadActor
  ↓
datascope.Actor + datascope.Merge(...)
  ↓
UserQueryScope / DepartmentQueryScope
  ↓
module/*/datascope.go
  ↓
Repository 查询结果被自动裁剪
```

只要这条主线稳定，后面无论扩用户、部门、岗位，还是扩真实业务模块，数据权限都能沿着同一套模式复用。

## 当前代码已经落到哪里

这章现在主要对应下面这些位置：

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
│  ├─ module/
│  │  └─ iam/
│  │     ├─ department/
│  │     │  └─ datascope.go
│  │     └─ user/
│  │        └─ datascope.go
│  └─ platform/
│     └─ datascope/
│        └─ datascope.go
└─ migrations/
   ├─ mysql/
   └─ postgres/
```

这说明第 5 章现在不再只是“给第 3 章补一个组织附录”，而是已经直接对应真实代码结构。

## 本章建议顺序

### 1. 先把组织模型定稳

先看 [组织模型设计](./organization-model-design)。

这一节回答的是：为什么企业后台的数据权限必须依赖部门、岗位、用户归属和角色范围一起设计。

### 2. 再看数据权限如何真正落到查询

再看 [角色数据范围与查询作用域](./role-data-scope-and-query-scopes)。

这一节是本章主线核心，会把下面几件事一次串起来：

- `data_scope`
- `Actor`
- 多角色并集
- `gorm.Scopes(...)`
- 用户资源与部门资源的两种过滤方式

### 3. 把平台规则和资源接法拆开读

如果你想把这一章真正读透，接着建议继续看：

- [Actor 上下文与多角色并集](./actor-and-grant-merge)
- [资源级数据权限接入模式](./module-datascope-patterns)
- [共享数据权限接入规范](./shared-datascope-integration-conventions)

这样更容易把：

- 请求期上下文装载
- 平台层并集合并
- 现成资源模式
- 新模块选模式

这三层分别看清楚，而不是把它们混在一页里死记。

### 4. 最后看组织资源如何变成真实模块

收完前两节后，再看：

- [部门树与部门管理](./department-tree-and-management)
- [岗位管理与用户归属](./post-management-and-user-affiliation)
- [真实业务模块的数据权限边界](./business-module-datascope-boundaries)
- [岗位资源的数据权限收紧时机](./post-datascope-tightening)

这样看会更顺，因为你会先知道“为什么要这样过滤”，再去看“模块里具体怎么落”。

## 本章最短主线

如果你不是第一次看第 5 章，而是想快速把“数据权限怎样真正落成代码”重新过一遍，当前最短主线建议只读下面 6 页：

| 顺序 | 页面 | 读完能得到什么 |
| --- | --- | --- |
| 1 | [组织模型设计](./organization-model-design) | 先确认组织、岗位、角色范围为什么必须一起设计 |
| 2 | [角色数据范围与查询作用域](./role-data-scope-and-query-scopes) | 建立整条数据权限执行链 |
| 3 | [Actor 上下文与多角色并集](./actor-and-grant-merge) | 看清请求期上下文和平台层并集语义 |
| 4 | [共享数据权限接入规范](./shared-datascope-integration-conventions) | 判断新模块该先选哪种 Scope |
| 5 | [datascope.go 与 Repository 边界](./datascope-and-repository-boundary) | 固定模块内真正落代码的边界 |
| 6 | [一次完整请求的权限过滤走读](./request-flow-walkthrough) | 把前面几页分散的点串成一次真实请求 |

如果你已经在接某个模块，最后再用：

- [数据权限落地检查清单](./data-scope-implementation-checklist)

做一轮验收会更稳。

## 本章各页分工

为了避免把这章读成“很多页都在重复讲 Scope”，可以直接按下面这张表记忆：

| 页面类型 | 当前代表页 | 主要回答什么 |
| --- | --- | --- |
| 总览页 | [角色数据范围与查询作用域](./role-data-scope-and-query-scopes) | 整条链路长什么样 |
| 平台规则页 | [Actor 上下文与多角色并集](./actor-and-grant-merge) | `Actor` 和 `Merge(...)` 到底做了什么 |
| 资源模式页 | [资源级数据权限接入模式](./module-datascope-patterns) | 当前已经存在的三类资源模式各在表达什么 |
| 模块落地页 | [共享数据权限接入规范](./shared-datascope-integration-conventions)、[datascope.go 与 Repository 边界](./datascope-and-repository-boundary) | 新模块该怎么选模式、代码该落哪层 |
| 验证页 | [一次完整请求的权限过滤走读](./request-flow-walkthrough)、[数据权限落地检查清单](./data-scope-implementation-checklist) | 请求如何真实执行、模块怎样验收 |

如果你是带着具体问题回来查，也可以直接按下面这张表回看：

| 现在卡在哪 | 优先回看哪页 |
| --- | --- |
| 不确定 `Actor` 里到底装了什么 | [Actor 上下文与多角色并集](./actor-and-grant-merge) |
| 不确定资源该选哪种 Scope | [资源级数据权限接入模式](./module-datascope-patterns)、[共享数据权限接入规范](./shared-datascope-integration-conventions) |
| 不确定代码该落在哪层 | [datascope.go 与 Repository 边界](./datascope-and-repository-boundary) |
| 不确定真实请求是怎么一路裁剪的 | [一次完整请求的权限过滤走读](./request-flow-walkthrough) |
| 不确定模块是不是已经接完 | [数据权限落地检查清单](./data-scope-implementation-checklist) |

## 本章完成后的判断标准

完成这一章后，你至少应该能回答下面五个问题：

1. 当前用户为什么不只需要 `user_id`，还需要 `department_id`、`role_codes` 和 `grants`
2. 为什么 `data_scope` 必须是角色模型的一部分
3. 多角色并集为什么要统一放在平台层，而不是交给每个模块自己解释
4. 为什么 `gorm.Scopes(...)` 比在 Handler 里散写过滤条件更稳
5. 为什么用户资源、部门资源和岗位资源不能共用完全一样的过滤逻辑

## 本章小节

- [组织模型设计](./organization-model-design)
- [角色数据范围与查询作用域](./role-data-scope-and-query-scopes)
- [Actor 上下文与多角色并集](./actor-and-grant-merge)
- [资源级数据权限接入模式](./module-datascope-patterns)
- [共享数据权限接入规范](./shared-datascope-integration-conventions)
- [datascope.go 与 Repository 边界](./datascope-and-repository-boundary)
- [一次完整请求的权限过滤走读](./request-flow-walkthrough)
- [数据权限落地检查清单](./data-scope-implementation-checklist)
- [部门树与部门管理](./department-tree-and-management)
- [岗位管理与用户归属](./post-management-and-user-affiliation)
- [真实业务模块的数据权限边界](./business-module-datascope-boundaries)
- [岗位资源的数据权限收紧时机](./post-datascope-tightening)

## 这一章结束后会走到哪里

当组织体系和数据权限主线定稳后，后面的系统模块和业务模块就不需要再重新讨论“数据权限放哪层”“当前用户上下文怎么取”这些底层问题了。

也就是说，第 5 章真正交付的不是几张组织表，而是：

> 一套后续模块可以不断复用的数据权限接入范式。

下一节先从组织模型开始：[组织模型设计](./organization-model-design)。
