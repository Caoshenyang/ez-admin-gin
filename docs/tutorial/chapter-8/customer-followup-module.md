---
title: CRM 客户跟进模块落地
description: "作为第 8 章扩展阅读，围绕 crm/followup 这条真实业务模块说明第二层业务动作如何继续沿既有骨架稳定扩展。"
---

# CRM 客户跟进模块落地

::: info 扩展阅读定位
这页保留为第 8 章的扩展实现参考，用来说明“基于 `crm/customer` 的第二层业务动作”怎样继续沿既有骨架落地。当前 canonical 主线仍然以 [CRM 客户模块示例](./business-module-example) 为唯一完整案例。
:::

前面的 `crm/customer` 已经证明“一个真实业务资源”可以接进底座。  
这一页继续往前走一步，补的不是第二个无关模块，而是**客户档案之上的业务动作层**：`crm/followup`。

这一步的意义很直接：

> 当前骨架不只适合做一张业务主表，也适合继续长出和它相关的第二层业务模块。

## 为什么继续补 `crm/followup`

`crm/customer` 解决的是“客户档案在不在系统里”。  
`crm/followup` 解决的是“客户进入系统之后，销售推进动作怎么留下业务记录”。

它很适合当下一批可复用模块的第一个样本，因为它同时具备：

- 真实业务关系：直接挂在客户档案之上
- 清晰的数据权限语义：继续服从部门与负责人范围
- 独立菜单和按钮：不是客户页里的临时弹窗逻辑
- 可单独成页：能证明前端页面、后端模块和种子可以继续按同一条主线扩

也就是说，这一页证明的是：

- `crm/customer` 不是孤立示例
- 当前底座已经能承接一个**业务模块家族**

## 这次接入链长什么样

```text
model/customer_followup.go
  ↓
module/crm/followup/
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
/api/v1/crm/followups
  ↓
000010 / 000011 迁移与种子
  ↓
admin/src/types/followup.ts
  ↓
admin/src/api/followup.ts
  ↓
admin/src/pages/crm/CustomerFollowUpView.vue
  ↓
dynamic-menu.ts
```

和 `crm/customer` 对照后，会发现它没有引入第二套结构，而是继续沿既有主线扩下去。

## 后端这一层证明了什么

`crm/followup` 当前继续保持：

- `dto / entity / repository / service / handler / routes / policy / datascope`
- CRM 聚合路由统一注册
- 菜单、按钮、Casbin 与迁移种子一起进入仓库

它新增的核心模型是 `sys_customer_followup`，字段里最关键的是：

- `customer_id`
- `department_id`
- `owner_user_id`

这里不是让跟进记录自己发明一套权限语义，而是**继承客户归属部门和负责人**，继续复用平台已有的数据权限规则。

## 数据权限为什么还能继续成立

客户跟进页继续沿用了 `datascope.UserQueryScope(...)`，只是资源从客户表换成了跟进表：

- 部门字段：`department_id`
- 负责人字段：`owner_user_id`

这样一来：

- `all` 角色可以看全部跟进
- `dept` 角色只能看本部门跟进
- `self` 角色只能看自己负责客户的跟进

这说明一个很关键的结论：

> 平台级数据权限语义，不只适用于“资源主表”，也适用于围绕主资源继续长出来的动作层记录。

## 为什么补了 `customer-options` 接口

这次实现里有一个很值得注意的细节：  
跟进页没有直接复用 `/crm/customers` 列表接口来给表单选客户，而是单独补了：

- `GET /api/v1/crm/followups/customer-options`

这样做是为了守住模块边界：

- 跟进页需要“可选客户”
- 但不应该因此偷偷依赖“客户列表权限”

也就是说，页面需要的辅助数据，仍然由**当前模块自己负责提供**，而不是把权限边界绕回别的模块。

## 前端这一层证明了什么

当前前端接入至少包括：

- `types/followup.ts`
- `api/followup.ts`
- `pages/crm/CustomerFollowUpView.vue`
- `dynamic-menu.ts`

页面里真实消费了：

- `crm:followup:create`
- `crm:followup:update`
- `crm:followup:status`

这说明它不是“能打开的占位页”，而是已经具备：

- 列表筛选
- 客户选择
- 新建
- 编辑
- 状态切换

## 这一页最关键的结论

`crm/followup` 补完之后，第 8 章现在已经不只是证明“能接一个业务模块”，而是进一步证明了：

1. 一个业务主资源可以稳定接进底座
2. 主资源之上的动作层模块也能继续沿同一骨架扩展
3. 数据权限、菜单、按钮、前端页面不会因为进入第二层业务关系就失效

也就是说，当前模块化接入规范已经开始具备“模块家族持续扩展”的能力，而不是只能展示一个孤立样本。

下一页回到整章验收视角，继续看 [模块接入验收清单](./module-integration-checklist)。
