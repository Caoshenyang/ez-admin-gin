---
title: 示例业务模块
description: "用当前真实的 crm/customer 模块走完一条完整接入链路，证明第 8 章这条模块接入主线已经可以稳定落地。"
---

# 示例业务模块

前面几页已经把结构、后端接入、权限菜单和前端页面接入顺序讲清楚了。现在还差最后一步：

> 找一个仓库里已经真实存在、前后端都接通了的模块，把整条接入链从头到尾串起来看一遍。

这一页直接用 `crm/customer` 做示例，而且它还是一个真正的非 `system` 分组业务模块。这样第 8 章讲的“模块化接入规范”，就不再只是系统模块内部的重复练习，而是已经可以落到真实业务域。

::: tip 🎯 本节目标
读完这一节，你应该能清楚看到：

当前第 8 章讲的接入规范，不是抽象建议，而是已经能在仓库里落地的一条真实路径。
:::

## 为什么选 `crm/customer`

`crm/customer` 很适合当第 8 章的示例，因为它同时具备：

- 后端结构完整
- 真正挂在 `/api/v1/crm/*`，不再属于 `system`
- 有列表、创建、编辑、状态切换
- 有 `policy.go`
- 有菜单、按钮和 Casbin 种子
- 有真实前端页面与动态菜单映射
- 有数据权限过滤，不只是“接口接通”

也就是说，它既足够真实，又能把“业务模块怎么接进底座”这件事讲清楚。

## 先看完整接入链

```text
model/customer.go
  ↓
module/crm/customer/
  ├─ dto.go
  ├─ entity.go
  ├─ datascope.go
  ├─ repository.go
  ├─ service.go
  ├─ handler.go
  ├─ routes.go
  └─ policy.go
  ↓
module/crm/routes.go
  ↓
bootstrap/router.go
  ↓
/api/v1/crm/customers
  ↓
migration seed data
  ├─ casbin_rule
  ├─ sys_menu
  └─ sys_role_menu
  ↓
admin/src/api/customer.ts
  ↓
admin/src/pages/crm/CustomerView.vue
  ↓
admin/src/router/dynamic-menu.ts
```

这条链几乎把第 8 章前面几页讲过的所有关键点都串起来了。

## 后端这一层是怎样成立的

客户模块目录当前是：

```text
server/internal/module/crm/customer/
├─ dto.go
├─ entity.go
├─ datascope.go
├─ repository.go
├─ service.go
├─ handler.go
├─ routes.go
└─ policy.go
```

它说明了一件事：

- 即使进入真实业务域，模块仍然沿 `dto / entity / repository / service / handler / routes / policy / datascope` 这套固定结构稳定落地

也就是说，第 8 章讲的“先进入固定结构”，不是只给 `system` 模块准备的。

## 权限与菜单这层是怎样成立的

`crm/customer` 当前已经有完整权限点：

- `crm:customer:list`
- `crm:customer:create`
- `crm:customer:update`
- `crm:customer:status`

接入系统后，除了接口本身，还要同时补：

- `casbin_rule`
- `sys_menu`
- `sys_role_menu`

否则就会出现：

- 接口能调，但角色无权访问
- 菜单不出现
- 页面能进，但按钮都不显示

## 前端这一层是怎样成立的

客户模块前端当前至少对应三处：

- `admin/src/api/customer.ts`
- `admin/src/pages/crm/CustomerView.vue`
- `admin/src/router/dynamic-menu.ts`

这里最关键的连接点是：

- 后端菜单节点里的 `component = crm/CustomerView`
- 前端 `routeComponentMap['crm/CustomerView']`

只要这两边一致，菜单树就能真正把页面打开。

## 按钮权限在示例里怎么体现

客户页当前已经真实消费了按钮权限：

- `crm:customer:create`
- `crm:customer:update`
- `crm:customer:status`

这说明一个模块是否“接完”，不能只看页面能不能打开，还要看：

- 页面内的真实操作入口是否也跟着权限树走通

## 数据权限在这个例子里怎么体现

这次的示例更关键的一点，是它还把数据权限真正接进去了。

`crm/customer` 的 `datascope.go` 固定的是一条典型业务资源规则：

- 部门字段用 `department_id`
- 负责人字段用 `owner_user_id`
- 查询时统一走 `datascope.UserQueryScope(...)`

这意味着：

- `all` 角色可以看全部客户
- `dept` 角色只能看本部门客户
- `self` 角色只能看自己负责的客户

也就是说，第 8 章现在不只是证明“模块能接进去”，还证明“模块能按平台统一的数据权限语义接进去”。

## 为什么这个例子足够说明第 8 章主线已经成立

因为 `crm/customer` 同时证明了：

1. 模块可以先按固定后端骨架稳定落地
2. 权限点可以集中收在 `policy.go`
3. 菜单、按钮和默认授权能进入种子
4. 前端能通过动态菜单真正开页
5. 页面动作显隐能继续服从按钮权限
6. 业务资源可以直接接入统一数据权限链路

也就是说，第 8 章不再只是“接入建议”，而是一条已经被真实模块证明过的接入主线。

## 本节最关键的结论

如果一个新模块要说自己“真正进入后台系统”，至少要像 `crm/customer` 这样同时具备：

- 后端结构
- 权限点
- 菜单种子
- 默认授权
- 前端页面
- 动态菜单映射
- 按钮权限消费

做到这里，它才不是一个只在代码层存在的模块，而是真正进入了当前后台底座。

如果你想继续看“客户档案之上的第二层业务动作”怎样接进来，下一页先读 [CRM 客户跟进模块落地](./customer-followup-module)。  
最终再回到整章验收视角：[模块接入验收清单](./module-integration-checklist)。
