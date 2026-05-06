---
title: CRM 客户模块示例
description: "用 `crm/customer` 这条真实业务模块走完一条完整接入链路，证明第 8 章这条模块接入主线已经可以稳定落地。"
---

# CRM 客户模块示例

前面几页已经把结构、后端接入、权限菜单和前端页面接入顺序讲清楚了。现在还差最后一步：

> 找一个仓库里真实存在、前后端都已经落地、并且已经接进测试与迁移链路的模块，把整条接入链从头到尾串起来看一遍。

这一页不再使用仓库里并不存在的通用示例，而是直接回到当前真实主案例：`crm/customer`。

::: info 这页的定位
`crm/customer` 在这里承担的是第 8 章的 canonical 主案例角色。

- 它是仓库里真实存在的业务模块
- 它已经接通迁移、菜单、Casbin、前端页面和数据权限
- `crm/followup` 会继续以它为基础，证明“模块家族可以继续扩展”
:::

::: tip 🎯 本节目标
读完这一节，你应该能清楚看到：

当前第 8 章讲的接入规范，不是抽象建议，而是一条已经在仓库里真实落地的完整路径。
:::

## 为什么选 `crm/customer`

`crm/customer` 很适合当第 8 章的主案例，因为它同时具备：

- 后端结构完整
- 真正挂在 `crm/*` 分组，语义上已经脱离内置系统资源
- 有列表、创建、编辑、状态切换
- 有 `policy.go`
- 有菜单、按钮和 Casbin 种子
- 有真实前端页面与动态菜单映射
- 有数据权限过滤和自动化测试覆盖

也就是说，它既足够真实，又能够完整证明“新业务模块怎样接进后台底座”。

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

当前仓库里的客户模块目录已经稳定落在：

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

- 即使进入真实业务域，模块仍然沿 `dto / entity / datascope / repository / service / handler / routes / policy` 这套固定结构稳定落地

也就是说，第 8 章讲的“先进入固定结构”，不是只给 `system` 模块准备的。

## 权限与菜单这层是怎样成立的

`crm/customer` 当前固定了下面这组稳定权限点：

- `crm:customer:list`
- `crm:customer:create`
- `crm:customer:update`
- `crm:customer:status`

接入系统后，除了接口本身，还同时补齐了：

- `casbin_rule`
- `sys_menu`
- `sys_role_menu`

否则就会出现：

- 接口能调，但角色无权访问
- 菜单不出现
- 页面能进，但按钮都不显示

这也正是第 8 章为什么一直强调“模块接入不能只写代码目录”的原因。

## 前端这一层是怎样成立的

客户模块前端至少对应三处：

- `admin/src/api/customer.ts`
- `admin/src/pages/crm/CustomerView.vue`
- `admin/src/router/dynamic-menu.ts`

这里最关键的连接点是：

- 后端菜单节点里的 `component = crm/CustomerView`
- 前端 `routeComponentMap['crm/CustomerView']`

只要这两边一致，菜单树就能真正把页面打开。

## 数据权限在这个例子里怎么体现

`crm/customer` 比普通 CRUD 更适合当主案例，还有一个关键原因：

> 它已经真实接进了第 5 章的数据权限链路。

当前模块围绕：

- `department_id`
- `owner_user_id`

继续复用平台级 `Actor` 与查询作用域，让不同角色看到不同客户数据。

这意味着第 8 章这里证明的不只是：

- 模块能接进系统

还证明了：

- 一个真实业务模块接进来之后，权限模型不会断层

## 按钮权限在示例里怎么体现

客户页当前真实消费的按钮权限通常是：

- `crm:customer:create`
- `crm:customer:update`
- `crm:customer:status`

这说明一个模块是否“接完”，不能只看页面能不能打开，还要看：

- 页面内的真实操作入口是否也跟着权限树走通

## 为什么这个例子足够说明第 8 章主线已经成立

因为 `crm/customer` 同时证明了：

1. 模块可以先按固定后端骨架稳定落地
2. 权限点可以集中收在 `policy.go`
3. 菜单、按钮和默认授权能进入种子
4. 前端能通过动态菜单真正开页
5. 页面动作显隐能继续服从按钮权限
6. 数据权限可以继续沿平台规则接进真实业务资源
7. 自动化测试已经覆盖成功路径和范围过滤场景

也就是说，第 8 章不再只是“接入建议”，而是一条已经在仓库里被真实证明的接入主线。

## 本节最关键的结论

如果一个新模块要说自己“真正进入后台系统”，至少要像 `crm/customer` 这样同时具备：

- 后端结构
- 权限点
- 菜单种子
- 默认授权
- 前端页面
- 动态菜单映射
- 数据权限接入
- 基本验证链路

做到这里，它才不是一个只在代码层存在的模块，而是真正进入了当前后台底座。

## 后续最自然的扩展方向

当前仓库已经沿这条主线往前走了一步，最自然的下一个扩展方向就是：

- 客户跟进 `crm/followup`

也就是说，第 8 章现在不只证明“能接一个业务模块”，还开始证明：

- 一个业务主资源接进来后，可以继续长出第二层业务动作模块

最终回到整章验收视角，继续看 [模块接入验收清单](./module-integration-checklist)。
