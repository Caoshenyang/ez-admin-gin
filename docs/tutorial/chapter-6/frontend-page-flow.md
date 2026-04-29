---
title: 前端页面接入流程
description: "说明一个新业务模块如何接入前端类型、API、页面和路由。"
---

# 前端页面接入流程

上一节已经把后端接口、权限和菜单补齐了。这一节把前端部分接上，让模块能从侧边栏打开并调用真实接口。

前端接入遵循固定的 4 步流程：类型 → API → 页面 → 路由。每一步都可以参照第5章已经写好的系统页面。

::: tip 🎯 本节目标
读完这一节后，你应该能按照固定顺序把任意一个后端模块接成可用的前端页面：定义类型、封装接口、创建 Vue 页面、注册路由映射。
:::

## 接入顺序和对应文件

| 步骤 | 文件 | 做什么 |
| --- | --- | --- |
| 1. 定义类型 | `admin/src/types/{module}.ts` | 状态常量、数据接口、查询参数、请求载荷 |
| 2. 封装 API | `admin/src/api/{module}.ts` | 对应后端每个接口的前端请求函数 |
| 3. 创建页面 | `admin/src/pages/{group}/{Module}View.vue` | 搜索区 + 数据表 + 弹框表单 |
| 4. 注册路由 | `admin/src/router/dynamic-menu.ts` | 把 component 映射到真实页面组件 |

## 第一步：定义类型

类型文件负责三件事：状态枚举、后端返回的数据接口、请求参数接口。

参照 `admin/src/types/user.ts` 的模式：

```text
types/user.ts
├── UserStatus        — 状态常量（Enabled / Disabled）
├── UserItem          — 单条数据接口（对应后端 JSON 字段）
├── UserListQuery     — 列表查询参数
├── UserListResponse  — 列表响应（items + total + page + page_size）
├── CreateUserPayload — 创建请求载荷
├── UpdateUserPayload — 编辑请求载荷
└── ...               — 其他操作载荷
```

::: details 类型文件的固定模式

```typescript
// 1. 状态常量
export const XxxStatus = {
  Enabled: 1,
  Disabled: 2,
} as const
export type XxxStatus = (typeof XxxStatus)[keyof typeof XxxStatus]

// 2. 数据项（字段名和后端 JSON 保持一致）
export interface XxxItem {
  id: number
  name: string
  status: XxxStatus
  created_at: string
  updated_at: string
}

// 3. 查询参数
export interface XxxListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: XxxStatus | 0
}

// 4. 列表响应
export interface XxxListResponse {
  items: XxxItem[]
  total: number
  page: number
  page_size: number
}

// 5. 创建/编辑载荷
export interface CreateXxxPayload {
  name: string
  status: XxxStatus
}
export interface UpdateXxxPayload {
  name: string
  status: XxxStatus
}
```

:::

关键约定：

- 字段名使用 `snake_case`，和后端 JSON 保持一致。
- 查询参数中 `status` 的类型写成 `XxxStatus | 0`，`0` 代表"全部"。
- 状态使用 `as const` + 类型推导，避免魔法数字。

## 第二步：封装 API

API 文件负责把每个后端接口封装成异步函数。

参照 `admin/src/api/user.ts` 的模式：

```typescript
import http from './http'
import type { ApiResponse } from '../types/http'
import type { XxxItem, XxxListQuery, XxxListResponse } from '../types/xxx'

// 列表查询
export async function getXxxs(params: XxxListQuery) {
  const response = await http.get<ApiResponse<XxxListResponse>>('/system/xxxs', { params })
  return response.data.data
}

// 创建
export async function createXxx(payload: CreateXxxPayload) {
  const response = await http.post<ApiResponse<XxxItem>>('/system/xxxs', payload)
  return response.data.data
}

// 编辑（注意后端使用 POST /:id/update 而不是 PUT）
export async function updateXxx(id: number, payload: UpdateXxxPayload) {
  const response = await http.post<ApiResponse<XxxItem>>(`/system/xxxs/${id}/update`, payload)
  return response.data.data
}
```

::: warning ⚠️ 后端编辑接口不是 RESTful PUT
后端统一使用 `POST /resources/:id/update` 而不是 `PUT /resources/:id`。前端封装时要注意路径匹配，否则会 404。
:::

## 第三步：创建页面

页面组件遵循统一的布局模式。可以参照第5章任意一个管理页面（用户、角色、菜单、配置）。

每个页面通常包含以下部分：

```text
ModuleView.vue
├── <script setup>
│   ├── 导入（Naive UI 组件、API、类型、权限）
│   ├── 状态（loading、列表、总数、查询参数）
│   ├── 表格列定义（columns）
│   ├── 弹框表单（formVisible、formModel、formMode）
│   ├── 数据加载函数（loadXxx）
│   ├── 操作函数（openCreate、openEdit、handleSubmit）
│   └── 权限函数（canUse）
├── <template>
│   ├── 标题 + 操作按钮
│   ├── 搜索卡片
│   ├── 数据表格卡片（表头 + NDataTable + 分页）
│   └── 新增/编辑弹框（NModal）
└── <style scoped>
    ├── 表格样式（表头加粗、hover 高亮）
    └── 弹框样式（圆角、渐变头部）
```

按钮权限用 `canUse` 函数控制：

```typescript
import { buttonPermissionCodes } from '../../router/dynamic-menu'

function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}
```

模板中配合 `v-if` 使用：

```vue
<NButton v-if="canUse('module:xxx:create')" type="primary" @click="openCreate">
  + 新增
</NButton>
```

::: details 按钮权限编码必须和菜单管理中的 code 一致
`canUse` 检查的是后端 `sys_menu` 表中 `type=button` 的记录的 `code` 字段。如果页面里写的是 `module:xxx:create`，菜单管理中新增按钮节点时，`code` 也必须填 `module:xxx:create`。
:::

## 第四步：注册路由映射

最后一步是在 `admin/src/router/dynamic-menu.ts` 中把 component 映射到真实页面。

```typescript
const routeComponentMap: Record<string, RouteComponent> = {
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': () => import('../pages/system/RoleView.vue'),
  // 新增模块映射
  'system/XxxView': () => import('../pages/system/XxxView.vue'),
}
```

映射关系的来源：

1. 后端迁移文件中种子菜单的 `Component` 字段（如 `"system/XxxView"`）。
2. 前端 `routeComponentMap` 的 key 必须和种子数据完全一致。
3. 如果 key 不在 map 中，会 fallback 到占位页。

::: warning ⚠️ 四个位置要一致
一个模块要正常工作，下面四处必须对齐：
1. 后端迁移文件中种子菜单的 `Component` 字段（如 `system/XxxView`）
2. 前端 `dynamic-menu.ts` 的 `routeComponentMap` key
3. 前端页面文件的实际路径（如 `pages/system/XxxView.vue`）
4. 菜单管理界面中菜单节点的 `component` 值（如果通过管理界面修改过）
:::

## 验证清单

完成以上 4 步后，按下面顺序验证：

1. 前端开发服务正常启动，无编译错误。
2. 用管理员登录，侧边栏能看到新模块的菜单入口。
3. 点击菜单，页面正常加载（不是占位页）。
4. 列表能正常请求后端接口并展示数据。
5. 新增/编辑弹框能正常提交并刷新列表。
6. 按钮权限在非管理员角色下正常隐藏。

## 本节小结

这一节把前端接入流程拆成了固定的 4 步：

- 类型定义保证字段名和后端一致。
- API 封装保证请求路径和方法正确。
- 页面组件遵循统一的搜索 + 表格 + 弹框模式。
- 路由映射把 component 字段关联到真实 Vue 文件。

下一节用一个完整的示例模块把前面所有步骤串起来：[示例业务模块](./sample-module)。
