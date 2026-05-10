---
title: 动态菜单与路由
description: 前端动态路由注册、组件白名单、菜单渲染机制
---

# 动态菜单与路由

## 加载时机

登录成功后，前端请求用户菜单数据：

```
登录成功 → GET /api/v1/auth/menus → 收到菜单树 + 按钮权限码
```

## 菜单数据结构

后端返回的菜单项包含：

```typescript
interface AuthMenu {
  id: number
  type: MenuType          // 1=目录, 2=菜单, 3=按钮
  code: string            // 权限编码
  title: string           // 显示名称
  path: string            // 路由路径
  component: string       // 组件标识
  icon: string            // 图标
  sort: number            // 排序
  children?: AuthMenu[]   // 子菜单
}
```

## 动态路由注册

`router/dynamic-menu.ts` 负责将后端菜单转换为 Vue Router 路由：

```
1. 遍历菜单树
2. type=1（目录）→ 创建嵌套路由父级
3. type=2（菜单）→ 根据 component 查找组件映射 → router.addRoute()
4. type=3（按钮）→ 提取 code → 存入按钮权限列表
```

### 组件白名单

前端维护一个组件映射表，只有白名单中的组件才会被注册：

```typescript
const componentMap: Record<string, () => Promise<Component>> = {
  'dashboard/DashboardHome': () => import('@/modules/auth/pages/DashboardHome.vue'),
  'iam/UserView': () => import('@/modules/iam/pages/UserView.vue'),
  'iam/RoleView': () => import('@/modules/iam/pages/RoleView.vue'),
  'iam/MenuView': () => import('@/modules/iam/pages/MenuView.vue'),
  'iam/DepartmentView': () => import('@/modules/iam/pages/DepartmentView.vue'),
  'iam/PostView': () => import('@/modules/iam/pages/PostView.vue'),
  'system/ConfigView': () => import('@/modules/system/pages/ConfigView.vue'),
  'system/DictView': () => import('@/modules/system/pages/DictView.vue'),
  'system/FileView': () => import('@/modules/system/pages/FileView.vue'),
  'system/OperationLogView': () => import('@/modules/system/pages/OperationLogView.vue'),
  'system/LoginLogView': () => import('@/modules/system/pages/LoginLogView.vue'),
  'system/NoticeView': () => import('@/modules/system/pages/NoticeView.vue'),
}
```

::: warning
新增页面后，必须同时更新后端菜单配置和前端组件白名单，否则路由无法注册。
:::

## 侧边栏渲染

动态菜单数据传递给 `NMenu` 组件渲染侧边栏：

```
AppSidebar
  → NMenu :options="menuOptions"
  → 菜单项点击 → router.push(path)
```

菜单状态由 Pinia store `admin-shell` 管理：

- `activeMenuKey`：当前活跃菜单
- `expandedMenuKeys`：展开的菜单组
- `collapsed`：侧边栏折叠状态

## 工作标签

登录后打开的页面以标签页形式展示在顶部：

```
AppHeader
  → WorkTabs
    → 标签列表（来自 admin-shell store）
    → 点击标签 → 路由跳转
    → 关闭标签 → 移除标签 + 跳转相邻标签
```

## 新增页面步骤

1. 在 `modules/` 下创建页面组件
2. 在 `router/dynamic-menu.ts` 的组件映射中添加条目
3. 在后端菜单管理中添加菜单项，`component` 填写映射 key
4. 将菜单权限分配给角色
5. 对应的 Casbin 策略需同步更新
