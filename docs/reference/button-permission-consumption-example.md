---
title: 按钮权限消费示例
description: "集中说明按钮权限码怎样从 /auth/menus 进入前端页面，以及 canUse(code) 在系统页里的最小稳定用法。"
---

# 按钮权限消费示例

这页只回答一个高频问题：

> 当前底座里，页面上的“新增 / 编辑 / 启停 / 上传”这些按钮，到底是怎么根据权限码显隐的？

::: tip 快速结论
当前按钮权限消费链路很固定：

`/api/v1/auth/menus` → 按钮节点 `code` → `collectButtonCodes(...)` → `buttonPermissionCodes` → 页面里的 `canUse(code)`
:::

## 先看真实链路

当前前端不是单独再请求一份“按钮权限列表”，而是直接复用菜单树：

```text
/api/v1/auth/menus
  ↓
AuthMenu[]
  ↓
dynamic-menu.ts / collectButtonCodes(...)
  ↓
buttonPermissionCodes
  ↓
页面里的 canUse(code)
  ↓
v-if / 下拉动作 / 批量操作显隐
```

相关文件：

- `admin/src/router/dynamic-menu.ts`
- `admin/src/modules/iam/pages/UserView.vue`
- `admin/src/modules/system/pages/FileView.vue`
- `admin/src/modules/system/pages/NoticeView.vue`

## 第一步：按钮权限码先从菜单树里收集出来

当前 `dynamic-menu.ts` 里，按钮权限来源就是按钮节点：

```ts
function collectButtonCodes(menus: AuthMenu[]) {
  const result: string[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Button) {
      result.push(menu.code)
    }

    result.push(...collectButtonCodes(menu.children ?? []))
  }

  return result
}
```

最终前端统一暴露的是：

```ts
export const buttonPermissionCodes = computed(() => {
  return collectButtonCodes(authMenus.value)
})
```

这意味着页面不需要知道：

- 权限树怎么递归
- 按钮节点在哪一层
- 后端是不是把它挂在目录下还是菜单下

页面只关心自己要不要消费某个 `code`。

## 第二步：页面里统一写一个 `canUse(code)`

系统页现在最稳定的写法就是这一层：

```ts
function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}
```

它的价值不是代码短，而是所有页面都在复用同一套判断语义：

- 当前账号有没有这个按钮权限码
- 有就显示操作入口
- 没有就只隐藏入口，不改变后端安全边界

::: warning 按钮权限只负责前端显隐，不负责真正安全
即使页面按钮隐藏了，真正的安全边界也仍然在后端接口权限和业务校验里。

前端按钮权限的职责只是：

- 降低误操作入口
- 让页面反馈更贴近当前角色
- 避免用户点进来才发现“没有权限”
:::

## 第三步：最常见的 3 种消费位置

### 1. 工具栏按钮

这是最常见的一类，例如用户页和公告页顶部的“新增”：

```vue
<NButton v-if="canUse('system:user:create')" type="primary" @click="openCreate">
  新增用户
</NButton>
```

或者：

```vue
<NButton v-if="canUse('system:notice:create')" type="primary" @click="openCreate">
  新增公告
</NButton>
```

适合：

- 新增
- 导出
- 上传
- 批量操作入口

### 2. 表格行操作

例如列表里“编辑 / 启停 / 分配角色”这类行级动作：

```ts
function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}
```

```ts
canUse('system:user:update')
canUse('system:user:status')
canUse('system:user:assign-role')
canUse('system:user:delete')
```

适合：

- 编辑
- 状态切换
- 角色授权
- 单条删除

### 3. 资源形态更特殊的页面动作

文件页最典型，因为它的主动作是上传：

```vue
<NButton
  v-if="canUse('system:file:upload')"
  type="primary"
  :loading="uploading"
>
  上传文件
</NButton>
```

这说明按钮权限粒度不是每页都一样，而是跟资源动作保持一致。

## 当前真实页面可以直接对照什么

如果你想看现成代码，当前最值得对照这几页：

| 页面 | 典型按钮权限码 |
| --- | --- |
| `UserView.vue` | `system:user:create` / `system:user:update` / `system:user:status` / `system:user:assign-role` / `system:user:delete` |
| `RoleView.vue` | `system:role:create` / `system:role:update` / `system:role:status` / `system:role:permission` / `system:role:menu` / `system:role:delete` |
| `PostView.vue` | `system:post:create` / `system:post:update` / `system:post:status` / `system:post:delete` |
| `ConfigView.vue` | `system:config:create` / `system:config:update` / `system:config:status` / `system:config:delete` |
| `DictView.vue` | `system:dict:type:create` / `system:dict:type:update` / `system:dict:type:delete` / `system:dict:item:create` / `system:dict:item:delete` |
| `NoticeView.vue` | `system:notice:create` / `system:notice:update` / `system:notice:status` / `system:notice:delete` |
| `FileView.vue` | `system:file:upload` |
| `MenuView.vue` | `system:menu:create` / `system:menu:update` / `system:menu:status` / `system:menu:delete` |
| `DepartmentView.vue` | `system:department:create` / `system:department:update` / `system:department:status` / `system:department:delete` |

## 新页面最小可复用模板

如果你在接一个新页面，当前最小稳定模板就是：

```ts
import { buttonPermissionCodes } from '../../router/dynamic-menu'

function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}
```

然后在模板里：

```vue
<NButton v-if="canUse('crm:customer:create')" type="primary">
  新增客户
</NButton>
```

你只需要保证三件事：

1. 后端 `sys_menu` 里已经有对应按钮节点
2. 按钮节点的 `code` 和页面里写的一致
3. 当前角色已经被授权到这个按钮节点

## 最常见的 4 个排查点

### 1. 页面按钮一直不显示

优先检查：

- `/api/v1/auth/menus` 里有没有这个按钮节点
- 按钮节点 `code` 是否和页面里的 `canUse(...)` 完全一致
- 当前角色是否拿到了这个按钮节点授权

### 2. 按钮显示了，但接口仍然报无权限

通常说明：

- 前端按钮权限和后端接口权限都存在
- 但当前角色只拿到了按钮节点，没有拿到接口策略

也就是说，前端显隐和后端接口鉴权不是同一层。

### 3. 菜单能开页，但按钮全没了

通常说明：

- 菜单节点授权了
- 按钮节点没有授权

因为按钮权限码来自菜单树里的 `type = Button` 节点，而不是页面菜单节点本身。

### 4. 改完角色授权后，前端马上没变化

先确认：

- 页面有没有重新请求 `/auth/menus`
- 登录态有没有刷新

前端按钮权限集合依赖当前菜单树，不是直接监听角色后台的改动。

## 和哪些页一起看最顺

- [权限码约定](./permission-code-conventions)
- [动态菜单组件白名单](./dynamic-menu-component-reference)
- [路由与菜单](/frontend/route-and-menu)
- [动态菜单](/architecture/dynamic-menu)
