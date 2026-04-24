---
title: 动态菜单
description: "根据后端菜单权限生成前端路由和侧边菜单。"
---

# 动态菜单

上一节已经把后台壳子搭好了，但侧边栏菜单还是写死在前端。现在把它升级成真正的权限菜单：登录后请求 `/api/v1/auth/menus`，根据后端返回的菜单树生成 `NMenu`，同时把可访问页面注册成前端路由。

::: tip 🎯 本节目标
完成后，侧边栏不再依赖前端静态数组，而是由当前登录用户的角色权限决定。刷新页面后，前端会重新拉取菜单并恢复动态路由；没有权限的菜单不会显示，也不能直接通过地址访问。
:::

## 先明确动态菜单的边界

动态菜单听起来像“后端返回什么，前端就渲染什么”，但真实项目里不能这么粗放。我们要把边界先定清楚：

| 后端字段 | 前端用途 |
| --- | --- |
| `type = 1` | 目录，只出现在侧边栏分组里 |
| `type = 2` | 菜单，会注册成可访问路由 |
| `type = 3` | 按钮权限，本节先收集起来，后续页面里再用于按钮显隐 |
| `path` | 路由地址，也是 `NMenu` 的 `key` |
| `component` | 后端保存的组件编码，前端用白名单映射到真实页面组件 |
| `children` | 递归生成菜单层级 |

::: warning ⚠️ 不要直接把后端 `component` 拼成动态 import
`component` 是数据库配置，不应该直接参与任意路径导入。前端要维护一张明确的组件白名单：后端只负责返回组件编码，前端只允许编码命中白名单后加载对应页面。

这样做虽然多写几行映射，但能避免菜单配置错误导致构建路径不可控，也更适合后续做页面级权限审计。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
└─ src/
   ├─ api/
   │  └─ menu.ts
   ├─ layouts/
   │  └─ AdminLayout.vue
   ├─ router/
   │  ├─ dynamic-menu.ts
   │  └─ index.ts
   └─ types/
      └─ menu.ts
```

| 位置 | 用途 |
| --- | --- |
| `src/types/menu.ts` | 定义后端菜单树类型 |
| `src/api/menu.ts` | 请求当前登录用户可见菜单 |
| `src/router/dynamic-menu.ts` | 把菜单树转换成侧边栏选项和动态路由 |
| `src/router/index.ts` | 登录后加载菜单，并注册动态路由 |
| `src/layouts/AdminLayout.vue` | 从动态菜单状态读取 `NMenu` 数据 |

::: info 这一节继续遵守前端约束
侧边栏仍然使用 Naive UI 的 `NMenu`，不要退回手写菜单列表；Tailwind CSS 4 只负责布局外层、间距和视觉微调。工作标签继续保留上一节贴近原型的轻量自实现方式。
:::

## 开始前先确认

开始之前，先确认下面几件事：

- 已完成上一节 [后台布局](./admin-layout)，并且登录后可以进入 `/dashboard`。
- 后端 `/api/v1/auth/menus` 已经能按当前登录用户返回菜单树。
- `admin/src/api/http.ts` 已经统一注入 `Authorization` 请求头。
- 当前系统页面还没有真实实现，所以本节会先把后端组件编码映射到占位页。

## 🛠️ 定义菜单类型

新增 `admin/src/types/menu.ts`。

```ts
export const MenuType = {
  Directory: 1,
  Menu: 2,
  Button: 3,
} as const

export type MenuType = (typeof MenuType)[keyof typeof MenuType]

// AuthMenu 对应 /api/v1/auth/menus 返回的菜单节点。
export interface AuthMenu {
  id: number
  parent_id: number
  type: MenuType
  code: string
  title: string
  path: string
  component: string
  icon: string
  sort: number
  children?: AuthMenu[]
}
```

这里没有把 `status` 放进前端类型，是因为 `/api/v1/auth/menus` 已经只返回启用状态的菜单。前端只关心“当前用户能看到什么”，不需要再判断菜单是否启用。

## 🛠️ 封装菜单接口

新增 `admin/src/api/menu.ts`。

```ts
import http from './http'

import type { AuthMenu } from '../types/menu'
import type { ApiResponse } from '../types/http'

// getCurrentUserMenus 获取当前登录用户可见的菜单树。
export async function getCurrentUserMenus() {
  const response = await http.get<ApiResponse<AuthMenu[]>>('/auth/menus')
  return response.data.data ?? []
}
```

这个接口会自动带上 `Authorization` 请求头，因为上一节已经在 `src/api/http.ts` 里统一处理了请求拦截器。

## 🛠️ 编写菜单转换工具

新增 `admin/src/router/dynamic-menu.ts`。这个文件做三件事：

- 保存当前用户菜单树，供布局读取。
- 把菜单树转换成 Naive UI 的 `MenuOption[]`。
- 把 `type = 2` 的菜单转换成 Vue Router 动态路由。

```ts
import type { MenuOption } from 'naive-ui'
import type { RouteRecordRaw } from 'vue-router'
import { computed, shallowRef } from 'vue'

import { MenuType, type AuthMenu } from '../types/menu'

type RouteComponent = NonNullable<RouteRecordRaw['component']>

const placeholderPage = () => import('../pages/system/PlaceholderPage.vue')

const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': placeholderPage,
  'system/UserView': placeholderPage,
  'system/RoleView': placeholderPage,
  'system/MenuView': placeholderPage,
  'system/ConfigView': placeholderPage,
  'system/FileView': placeholderPage,
  'system/OperationLogView': placeholderPage,
  'system/LoginLogView': placeholderPage,
}

const builtinMenuOptions: MenuOption[] = [
  {
    label: '工作台',
    key: '/dashboard',
  },
]

export const authMenus = shallowRef<AuthMenu[]>([])

export const sideMenuOptions = computed<MenuOption[]>(() => {
  return [...builtinMenuOptions, ...buildMenuOptions(authMenus.value)]
})

export const buttonPermissionCodes = computed(() => {
  return collectButtonCodes(authMenus.value)
})

export function setAuthMenus(menus: AuthMenu[]) {
  authMenus.value = menus
}

export function clearAuthMenus() {
  authMenus.value = []
}

export function buildDynamicRoutes(menus: AuthMenu[]) {
  return collectPageMenus(menus).map<RouteRecordRaw>((menu) => ({
    path: toChildRoutePath(menu.path),
    name: `menu-${menu.id}`,
    component: resolveRouteComponent(menu.component),
    props: {
      title: menu.title,
      description: `${menu.title} 页面后续会接入真实业务。`,
    },
    meta: {
      title: menu.title,
      menuCode: menu.code,
    },
  }))
}

export function findMenuTitleByPath(path: string) {
  return collectPageMenus(authMenus.value).find((menu) => menu.path === path)?.title
}

function buildMenuOptions(menus: AuthMenu[]) {
  return menus.map(toMenuOption).filter(isMenuOption)
}

function toMenuOption(menu: AuthMenu): MenuOption | null {
  if (menu.type === MenuType.Button) {
    return null
  }

  const children = buildMenuOptions(menu.children ?? [])
  const key = menu.path || menu.code

  return {
    label: menu.title,
    key,
    disabled: menu.type === MenuType.Directory && children.length === 0,
    children: children.length > 0 ? children : undefined,
  }
}

function isMenuOption(option: MenuOption | null): option is MenuOption {
  return option !== null
}

function collectPageMenus(menus: AuthMenu[]) {
  const result: AuthMenu[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Menu && menu.path) {
      result.push(menu)
    }

    result.push(...collectPageMenus(menu.children ?? []))
  }

  return result
}

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

function resolveRouteComponent(component: string) {
  return routeComponentMap[component] ?? placeholderPage
}

function toChildRoutePath(path: string) {
  return path.replace(/^\/+/, '')
}
```

::: details 为什么这里先把系统页面都映射到占位页
后端初始化的菜单已经包含 `system/UserView`、`system/RoleView`、`system/MenuView` 等组件编码，但第五章后面的页面还没正式实现。

本节先把这些编码映射到 `PlaceholderPage.vue`，好处是动态菜单链路能先跑通。等后续写到用户、角色、菜单页面时，只需要把白名单里的某一项替换成真实页面即可。

例如后续用户管理页面完成后，把这一行：

```ts
'system/UserView': placeholderPage,
```

替换成：

```ts
'system/UserView': () => import('../pages/system/UserView.vue'),
```
:::

## 🛠️ 在路由守卫里加载菜单

修改 `admin/src/router/index.ts`。这里要把“是否登录”和“动态菜单是否准备好”放到全局守卫里处理。

这一处建议直接替换完整文件，不要只复制其中几行。原因是动态菜单会同时影响登录拦截、路由注册、刷新恢复和退出清理，拆开改很容易漏掉其中一环。

替换后，这个文件只保留 3 类路由：

- `/login`：登录页。
- `/dashboard`：固定内置的工作台。
- 后端菜单接口返回的页面：运行时挂到 `name: 'admin'` 的后台布局下面。

```ts
import { createRouter, createWebHistory } from 'vue-router'

import { getCurrentUserMenus } from '../api/menu'
import { clearAuthSession, hasAccessToken } from '../utils/auth'
import {
  buildDynamicRoutes,
  clearAuthMenus,
  setAuthMenus,
} from './dynamic-menu'

const removeDynamicRouteCallbacks: Array<() => void> = []
let dynamicRoutesReady = false

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: () => (hasAccessToken() ? '/dashboard' : '/login'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../pages/auth/LoginPage.vue'),
    },
    {
      path: '/',
      name: 'admin',
      component: () => import('../layouts/AdminLayout.vue'),
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('../pages/dashboard/DashboardHome.vue'),
          meta: { title: '工作台' },
        },
      ],
    },
  ],
})

router.beforeEach(async (to) => {
  if (to.path === '/login') {
    return hasAccessToken() ? '/dashboard' : true
  }

  if (!hasAccessToken()) {
    resetDynamicRoutes()
    return {
      path: '/login',
      query: {
        redirect: to.fullPath,
      },
    }
  }

  if (!dynamicRoutesReady) {
    try {
      const menus = await getCurrentUserMenus()
      setAuthMenus(menus)

      for (const route of buildDynamicRoutes(menus)) {
        removeDynamicRouteCallbacks.push(router.addRoute('admin', route))
      }

      // 动态路由刚注册完成，需要重新匹配一次当前目标地址。
      dynamicRoutesReady = true
      return to.fullPath
    } catch {
      clearAuthSession()
      resetDynamicRoutes()
      return '/login'
    }
  }

  return true
})

// resetDynamicRoutes 用于退出登录或 Token 失效时清理旧账号菜单。
export function resetDynamicRoutes() {
  for (const removeRoute of removeDynamicRouteCallbacks) {
    removeRoute()
  }

  removeDynamicRouteCallbacks.length = 0
  dynamicRoutesReady = false
  clearAuthMenus()
}

export default router
```

替换后重点看两处：

- `name: 'admin'` 必须保留。后面 `router.addRoute('admin', route)` 会把动态页面挂到这个布局下面。
- `return to.fullPath` 必须保留。它会让 Vue Router 在动态路由注册完成后重新匹配当前地址。否则刷新 `/system/users` 时，第一次匹配会找不到刚注册的路由。

::: warning ⚠️ 退出登录时要重置动态路由
动态路由是运行时加进去的。如果退出登录后不清掉，换另一个低权限账号登录时，旧账号注册过的路由可能还留在前端路由表里。

所以这里把清理动作封装成 `resetDynamicRoutes()`，退出登录、Token 失效、未登录访问受保护页面时都可以复用。
:::

## 🛠️ 让后台布局使用动态菜单

修改 `admin/src/layouts/AdminLayout.vue`。这里只改脚本里和菜单相关的部分：删除上一节的静态 `menuItems`，改成从 `dynamic-menu.ts` 读取菜单选项。

这里也建议先替换完整 `<script setup>`。这样比“删几行、加几行”更稳，能避免出现下面两类常见错误：

- `Cannot find name 'menuItems'`：旧的 `activeMenuKey` 还在引用已删除的静态菜单数组。
- `'routeTitle' is declared but its value is never read`：新增了 `routeTitle`，但没有让工作标签复用它。

把 `AdminLayout.vue` 顶部的整个 `<script setup lang="ts">...</script>` 替换为下面这段：

```vue
<script setup lang="ts">
import {
  ChevronDownOutline,
  EllipsisHorizontal,
  ExpandOutline,
  LogOutOutline,
  MoonOutline,
  NotificationsOutline,
  SearchOutline,
} from '@vicons/ionicons5'
import type { DropdownOption } from 'naive-ui'
import {
  NButton,
  NDropdown,
  NIcon,
  NInput,
  NLayout,
  NLayoutContent,
  NLayoutHeader,
  NLayoutSider,
  NMenu,
  useMessage,
} from 'naive-ui'
import { computed, h, ref, watch } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'

import { resetDynamicRoutes } from '../router'
import { findMenuTitleByPath, sideMenuOptions } from '../router/dynamic-menu'
import { clearAuthSession, getAuthUserInfo } from '../utils/auth'

interface WorkTab {
  title: string
  to: string
  closable: boolean
}

const route = useRoute()
const router = useRouter()
const message = useMessage()

const openTabs = ref<WorkTab[]>([{ title: '工作台', to: '/dashboard', closable: false }])

const currentUser = computed(() => getAuthUserInfo())
const displayName = computed(() => {
  return currentUser.value?.nickname || currentUser.value?.username || '管理员'
})

const routeTitle = computed(() => {
  return String(route.meta.title ?? findMenuTitleByPath(route.path) ?? '工作台')
})

const breadcrumbText = computed(() => {
  return `首页 / ${routeTitle.value}`
})

const activeMenuKey = computed(() => {
  return route.path
})

const dropdownOptions: DropdownOption[] = [
  {
    label: '退出登录',
    key: 'logout',
    icon: () =>
      h(NIcon, null, {
        default: () => h(LogOutOutline),
      }),
  },
]

function ensureCurrentTab() {
  const title = routeTitle.value
  if (!title || route.path === '/login') {
    return
  }

  const exists = openTabs.value.some((tab) => tab.to === route.path)
  if (exists) {
    return
  }

  openTabs.value.push({
    title,
    to: route.path,
    closable: route.path !== '/dashboard',
  })
}

function navigateTo(path: string) {
  void router.push(path)
}

function handleMenuUpdate(key: string | number) {
  navigateTo(String(key))
}

function handleCloseTab(path: string) {
  const nextTabs = openTabs.value.filter((tab) => tab.to !== path)
  openTabs.value =
    nextTabs.length > 0 ? nextTabs : [{ title: '工作台', to: '/dashboard', closable: false }]

  if (route.path === path) {
    const fallback = openTabs.value[openTabs.value.length - 1]
    if (!fallback) {
      void router.push('/dashboard')
      return
    }

    void router.push(fallback.to)
  }
}

function handleCloseOtherTabs() {
  const current = openTabs.value.find((tab) => tab.to === route.path)

  openTabs.value = [{ title: '工作台', to: '/dashboard', closable: false }]

  if (current && current.to !== '/dashboard') {
    openTabs.value.push(current)
  }
}

function handleRefresh() {
  window.location.reload()
}

function handleUserAction(key: string | number) {
  if (key !== 'logout') {
    return
  }

  clearAuthSession()
  resetDynamicRoutes()
  message.success('已退出登录')
  void router.replace('/login')
}

watch(
  () => route.fullPath,
  () => {
    ensureCurrentTab()
  },
  { immediate: true },
)
</script>
```

替换完脚本后，再改模板里的 `NMenu`。只需要确认 `:options` 使用的是 `sideMenuOptions`：

```vue{4}
<NMenu
  class="mt-3"
  :value="activeMenuKey"
  :options="sideMenuOptions"
  :indent="18"
  inverted
  @update:value="handleMenuUpdate"
/>
```

这一段完成后，可以按下面 5 条快速自查：

- `AdminLayout.vue` 里不再导入 `MenuOption`。
- `AdminLayout.vue` 里不再声明 `MenuItem`、`menuItems`、`menuOptions`。
- `activeMenuKey` 只返回 `route.path`，不再扫描静态菜单数组。
- `routeTitle` 同时被 `breadcrumbText` 和 `ensureCurrentTab` 使用。
- 退出登录时会调用 `resetDynamicRoutes()`。

::: info 为什么这里不把 `NMenu` 拆成一堆按钮
动态菜单是典型的后台基础能力，Naive UI 已经提供了键盘交互、展开层级、激活状态和禁用状态。我们只需要把后端菜单树转换成 `MenuOption[]`，没必要再自己维护一套菜单组件。
:::

## 🧪 验证动态菜单

启动前后端后，用管理员账号登录：

```text
admin / EzAdmin@123456
```

登录成功后重点看这几件事：

| 验证点 | 预期结果 |
| --- | --- |
| Network 中请求 `/api/v1/auth/menus` | 返回当前用户可见菜单树 |
| 左侧菜单 | 出现后端返回的系统管理菜单 |
| 点击用户管理、角色权限等菜单 | 能进入对应 `/system/...` 地址 |
| 工作标签 | 标题跟随当前菜单页面变化 |
| 刷新 `/system/users` | 页面能恢复，菜单不会丢失 |
| 退出登录再登录 | 动态路由重新加载，不沿用旧菜单 |

如果要从命令行验证接口，可以先登录拿到 Token，再请求菜单：

::: code-group

```powershell [PowerShell]
$login = Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/auth/login `
  -ContentType "application/json" `
  -Body '{"username":"admin","password":"EzAdmin@123456"}'

Invoke-RestMethod `
  -Method Get `
  -Uri http://localhost:8080/api/v1/auth/menus `
  -Headers @{ Authorization = "$($login.data.token_type) $($login.data.access_token)" }
```

```bash [curl]
TOKEN=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"EzAdmin@123456"}' \
  | jq -r '.data.access_token')

curl http://localhost:8080/api/v1/auth/menus \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

最后跑一次前端检查：

```bash
cd admin
pnpm exec oxlint .
pnpm exec vue-tsc --noEmit
```

## 常见问题

### 刷新动态页面后变成空白

优先检查 `router.beforeEach` 里动态路由注册完成后，是否返回了 `to.fullPath`。这个返回值会让路由重新匹配一次刚刚注册的动态地址。

### 菜单显示了，但点击没有页面

检查后端 `component` 是否命中了 `routeComponentMap`。如果没有命中，本节的代码会先回退到 `PlaceholderPage.vue`，不会直接报错。后续接真实页面时，再把对应编码映射到真实组件。

### 目录菜单点不了

这是正常的。`type = 1` 的目录只负责分组，不注册页面路由。真正可点击的页面节点应该是 `type = 2`。

### 按钮权限在哪里用

本节先把 `type = 3` 的按钮权限收集到 `buttonPermissionCodes`，后续写用户管理、角色管理页面时，再用它控制“新增、编辑、禁用、删除”等按钮是否显示。

## 本节小结

这一节完成了前端权限菜单的核心链路：

```text
登录成功
  ↓
请求 /api/v1/auth/menus
  ↓
保存当前用户菜单树
  ↓
生成 NMenu 侧边栏
  ↓
注册 type = 2 的动态路由
  ↓
刷新页面后重新加载并恢复访问
```

下一节开始进入真实业务页面，把用户列表、查询和基础操作接到前端页面里：[用户管理页面](./user-pages)。
