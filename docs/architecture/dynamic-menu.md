---
title: 动态菜单机制
description: 后端驱动的前端动态菜单注册、路由生成与按钮权限控制
---

# 动态菜单机制

EZ Admin 的菜单系统是后端驱动的：后端定义菜单树和权限，前端根据登录用户的权限动态渲染菜单、注册路由和控制按钮。

## 菜单数据模型

### 三级菜单类型

| 类型 | 值 | 用途 | 前端表现 |
|------|---|------|---------|
| 目录 | 1 | 菜单分组 | 侧边栏可展开的文件夹 |
| 菜单 | 2 | 可访问页面 | 侧边栏菜单项 + Vue Router 路由 |
| 按钮 | 3 | 操作权限 | 页面内按钮的显示/隐藏控制 |

### 菜单表结构

```
sys_menu
├── id              主键
├── parent_id       父菜单 ID（树结构，0 为顶级）
├── type            1=目录, 2=菜单, 3=按钮
├── code            唯一权限编码（如 system:user:list）
├── title           显示名称
├── path            前端路由路径（如 /system/user）
├── component       前端组件路径（如 system/UserView）
├── icon            菜单图标
├── sort            排序值（升序）
├── status          状态（启用/禁用）
├── permission      后端权限标识
├── created_at / updated_at / deleted_at
```

## 菜单加载流程

```
1. 前端登录成功
   ↓
2. 调用 GET /api/v1/auth/menus
   ↓
3. 后端根据用户角色查询关联的菜单 ID 集合
   ├── 过滤 type=1,2 → 构建菜单树返回前端
   └── 过滤 type=3 → 返回按钮权限码列表
   ↓
4. 前端接收数据
   ├── 菜单树 → NMenu 组件渲染侧边栏
   ├── 菜单树 → 动态注册 Vue Router 路由
   └── 按钮权限码 → 存入 usePermission composable
```

## 前端动态路由注册

前端通过 `router/dynamic-menu.ts` 实现动态路由注册：

1. 遍历后端返回的菜单树
2. 根据 `component` 字段匹配前端组件（白名单映射）
3. `type=1`（目录）→ 嵌套路由父级
4. `type=2`（菜单）→ 具体页面路由，使用 `router.addRoute()` 动态添加
5. `type=3`（按钮）→ 提取权限码，不注册路由

### 组件白名单

前端维护一个组件映射白名单，`component` 字段的值对应实际的 Vue 组件：

```typescript
// router/dynamic-menu.ts 中的组件映射
const componentMap: Record<string, Component> = {
  'iam/UserView': () => import('@/modules/iam/pages/UserView.vue'),
  'iam/RoleView': () => import('@/modules/iam/pages/RoleView.vue'),
  // ...
}
```

只有白名单中的组件才会被注册为路由，防止恶意注入。

## 按钮权限控制

### 后端保障

每个按钮权限（`type=3`）对应一个唯一的 `code`，如 `system:user:create`。该 code 同时关联到 Casbin 策略，确保即使前端绕过按钮隐藏，后端仍然会拒绝未授权的 API 调用。

### 前端消费

```typescript
// composables/usePermission.ts
import { usePermission } from '@/composables/usePermission'

const { canUse } = usePermission()

// 模板中使用
<NButton v-if="canUse('system:user:create')">新建</NButton>
<NButton v-if="canUse('system:user:delete')" type="error">删除</NButton>
```

`canUse()` 检查当前用户的按钮权限码列表中是否包含指定 code。

## 菜单管理

系统内置菜单管理页面（`iam/MenuView`），支持：

- 树形展示所有菜单
- 新增/编辑/删除菜单项
- 设置菜单类型、图标、路径、组件、排序
- 配置按钮权限的 code

::: warning
修改菜单后，需要同步更新 Casbin 策略和前端组件白名单，否则新菜单无法正常工作。
:::

## 种子数据

系统迁移中预置了完整的菜单树，涵盖所有系统管理功能：

- 系统管理目录（用户、角色、菜单、部门、岗位）
- 系统工具目录（配置、字典、文件、日志、公告）
- 每个菜单下的标准按钮权限（list/create/update/delete）
