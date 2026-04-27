---
title: 后台布局
description: "实现后台基础布局，包括侧边栏、顶部栏、工作标签和一屏内容区。"
---

# 后台布局

登录页已经能完成认证链路，下一步要把登录后的默认页面从“单独一个工作台页”升级成真正的后台壳子。这一节会按原型补出后台布局的主骨架：左侧侧边栏、顶部栏、工作标签，以及放置业务页面的内容区。

::: tip 🎯 本节目标
完成后，管理台会进入一个固定一屏高度的后台框架：浏览器默认滚动条不出现，侧边栏和顶部栏固定，内容区通过路由切换不同页面，工作标签能表现基础的打开、关闭和关闭其他交互。
:::

![后台布局原型](/prototypes/exports/mipsi.png)

::: info 这一节参考的原型
- 原型文件：`docs/public/prototypes/ui.pen`
- 节点名称：`02 工作台 / Dashboard`
- 导出图片：`docs/public/prototypes/exports/mipsi.png`
:::

::: warning ⚠️ 这一节先实现静态后台壳子
这一节的重点是把后台框架搭稳，不是把所有系统页面一次性接完。

所以当前边界是：

- 侧边栏菜单先使用前端静态配置。
- 工作标签先基于当前前端路由维护。
- 用户管理、角色权限、菜单管理、操作日志、系统设置先放占位页，用于验证布局和路由出口。
- 下一节 [动态菜单](./dynamic-menu) 再把菜单和路由与后端 `/api/v1/auth/menus` 真正打通。
:::

## 先明确这一节的布局约束

后台布局和登录页不一样，它会承载后续所有系统页面，所以先把壳子的规则定死，后面会省很多返工。

本节统一遵守下面几条约束：

- 页面整体高度限制在一屏内，根容器使用 `h-screen`。
- 不出现浏览器默认滚动条，`html`、`body`、`#app` 和布局根节点都要关闭浏览器级滚动。
- 如果后续有超出一屏的内容，优先让内容区内部滚动，而不是让整个浏览器页面滚动。
- Naive UI 优先负责后台布局、菜单、标签、输入框、按钮、卡片、空状态、下拉菜单这些成熟组件。
- Tailwind CSS 4 负责组件外层的尺寸约束、留白、颜色微调和响应式补充。

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
└─ src/
   ├─ api/
   │  └─ dashboard.ts
   ├─ layouts/
   │  └─ AdminLayout.vue
   ├─ pages/
   │  ├─ dashboard/
   │  │  └─ DashboardHome.vue
   │  └─ system/
   │     └─ PlaceholderPage.vue
   ├─ router/
   │  └─ index.ts
   ├─ styles/
   │  └─ main.css
   ├─ types/
   │  └─ dashboard.ts
   └─ utils/
      └─ auth.ts
```

| 位置 | 用途 |
| --- | --- |
| `src/styles/main.css` | 关闭浏览器级滚动，保证后台壳子固定一屏 |
| `src/utils/auth.ts` | 读取本地登录用户信息，用在顶部栏用户区 |
| `src/types/dashboard.ts` | 约束工作台概览接口返回值，避免页面继续写死假数据 |
| `src/api/dashboard.ts` | 调用 `/api/v1/auth/dashboard`，统一获取工作台真实数据 |
| `src/layouts/AdminLayout.vue` | 实现侧边栏、顶部栏、工作标签和内容区 |
| `src/pages/dashboard/DashboardHome.vue` | 把原来的演示工作台升级成真实项目首页 |
| `src/pages/system/PlaceholderPage.vue` | 给静态菜单提供可复用的占位页 |
| `src/router/index.ts` | 把布局挂到受保护路由上，让工作台和系统页共用后台壳子 |

## 开始前先确认

开始前先确认下面几件事：

- 已完成上一节 [登录页](./login-page)。
- 登录成功后，当前已经能跳转到 `/dashboard`。
- `admin/src/styles/main.css` 已经接入 `@import "tailwindcss"`。
- 当前项目后端已经提供 `GET /api/v1/auth/dashboard`，登录后可返回工作台概览数据。

## 🛠️ 先把浏览器默认滚动关掉

修改 `admin/src/styles/main.css`。本次主要补一件事：把浏览器级滚动收掉，后面的后台布局才能稳定限制在一屏内。

::: details `admin/src/styles/main.css` — 关闭浏览器级滚动

```css
@import "tailwindcss";

@theme {
  --font-sans:
    "Inter", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  --color-brand-500: #18a058;
  --color-brand-600: #169250;
  --color-surface-page: #f5f7fb;
  --color-text-main: #1f2430;
}

:root {
  color-scheme: light;
  font-family: var(--font-sans);
  color: var(--color-text-main);
  background: var(--color-surface-page);
  font-synthesis: none;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

* {
  box-sizing: border-box;
}

html,
body,
#app {
  min-width: 320px;
  min-height: 100vh;
  margin: 0;
  /* 关闭浏览器默认滚动条 */ /* [!code ++] */
  overflow: hidden; /* [!code ++] */
}

body {
  background: var(--color-surface-page);
}

a {
  color: inherit;
  text-decoration: none;
}

button,
input,
textarea,
select {
  font: inherit;
}

button:not(:disabled), /* [!code ++] */
[role='button']:not([aria-disabled='true']), /* [!code ++] */
.n-button:not(.n-button--disabled) { /* [!code ++] */
  cursor: pointer; /* [!code ++] */
} /* [!code ++] */
```

:::

::: details 为什么这里直接关掉浏览器级滚动
后台管理台和内容站不一样，它更像一个应用壳子。只要布局层开始依赖“顶部栏固定、侧边栏固定、工作标签固定”，浏览器默认滚动就会把整套结构拉散。

这一节先统一成：

- 浏览器页面不滚。
- 真正需要滚动的地方，后续放到内部内容区单独处理。

这样顶部栏、标签栏和侧边栏才能稳定贴在原型位置上。
:::

## 🛠️ 补齐登录用户读取能力

修改 `admin/src/utils/auth.ts`。顶部栏需要显示当前登录用户昵称，所以这里补一个读取本地用户信息的函数。

本次要改两点：

- 增加 `getAuthUserInfo`。
- 解析失败时顺手清理异常登录态，避免顶部栏读到坏数据。

::: details `admin/src/utils/auth.ts` — 补齐登录用户读取能力

```ts
import type { LoginResponse } from '../types/auth'

const ACCESS_TOKEN_KEY = 'ez-admin-access-token'
const TOKEN_TYPE_KEY = 'ez-admin-token-type'
const USER_INFO_KEY = 'ez-admin-user-info'

type StorageMode = 'local' | 'session'

export interface AuthUserInfo {
  userId: number
  username: string
  nickname: string
  expiresAt: string
}

function getStorage(mode: StorageMode) {
  return mode === 'local' ? localStorage : sessionStorage
}

function readStorageValue(key: string) {
  return localStorage.getItem(key) ?? sessionStorage.getItem(key) ?? ''
}

// setAuthSession 在登录成功后保存本地登录态。
export function setAuthSession(payload: LoginResponse, rememberLogin: boolean) {
  clearAuthSession()

  const storage = getStorage(rememberLogin ? 'local' : 'session')

  storage.setItem(ACCESS_TOKEN_KEY, payload.access_token)
  storage.setItem(TOKEN_TYPE_KEY, payload.token_type)
  storage.setItem(
    USER_INFO_KEY,
    JSON.stringify({
      userId: payload.user_id,
      username: payload.username,
      nickname: payload.nickname,
      expiresAt: payload.expires_at,
    } satisfies AuthUserInfo),
  )
}

export function clearAuthSession() {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(TOKEN_TYPE_KEY)
  localStorage.removeItem(USER_INFO_KEY)

  sessionStorage.removeItem(ACCESS_TOKEN_KEY)
  sessionStorage.removeItem(TOKEN_TYPE_KEY)
  sessionStorage.removeItem(USER_INFO_KEY)
}

export function getAccessToken() {
  return readStorageValue(ACCESS_TOKEN_KEY)
}

export function getTokenType() {
  return readStorageValue(TOKEN_TYPE_KEY) || 'Bearer'
}

export function hasAccessToken() {
  return getAccessToken() !== ''
}

export function getAuthUserInfo() { // [!code ++]
  const raw = readStorageValue(USER_INFO_KEY) // [!code ++]
  if (!raw) { // [!code ++]
    return null // [!code ++]
  } // [!code ++]

  try { // [!code ++]
    return JSON.parse(raw) as AuthUserInfo // [!code ++]
  } catch { // [!code ++]
    clearAuthSession() // [!code ++]
    return null // [!code ++]
  } // [!code ++]
} // [!code ++]

// getAuthorizationHeader 统一拼接 Authorization 请求头。
export function getAuthorizationHeader() {
  const accessToken = getAccessToken()
  if (!accessToken) {
    return ''
  }

  return `${getTokenType()} ${accessToken}`
}
```

:::

## 🛠️ 创建后台布局组件

创建 `admin/src/layouts/AdminLayout.vue`。这一步是本节的核心，它负责把原型里的四层结构一次性搭起来。

这里优先使用 Naive UI 自带组件承接后台壳子：

- `NLayout` / `NLayoutSider` / `NLayoutHeader` / `NLayoutContent` 负责页面框架。
- `NMenu` 负责侧边栏菜单，不再手写菜单按钮列表。
- 工作标签保留轻量自实现，以贴近原型里的小标签视觉。

页面里仍然会保留少量 Tailwind 类，用来控制高度、留白、颜色和一屏约束。

- 左侧静态菜单
- 顶部栏
- 工作标签
- 路由内容区

这里直接完整写入即可。

::: details `admin/src/layouts/AdminLayout.vue` — 后台布局组件

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
import type { DropdownOption, MenuOption } from 'naive-ui'
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

import { clearAuthSession, getAuthUserInfo } from '../utils/auth'

interface MenuItem {
  title: string
  to: string
}

interface WorkTab {
  title: string
  to: string
  closable: boolean
}

const route = useRoute()
const router = useRouter()
const message = useMessage()

const menuItems: MenuItem[] = [
  { title: '工作台', to: '/dashboard' },
  { title: '用户管理', to: '/users' },
  { title: '角色权限', to: '/roles' },
  { title: '菜单管理', to: '/menus' },
  { title: '操作日志', to: '/logs' },
  { title: '系统设置', to: '/settings' },
]

const menuOptions: MenuOption[] = menuItems.map((item) => ({
  label: item.title,
  key: item.to,
}))

const openTabs = ref<WorkTab[]>([
  { title: '工作台', to: '/dashboard', closable: false },
])

const currentUser = computed(() => getAuthUserInfo())
const displayName = computed(() => {
  return currentUser.value?.nickname || currentUser.value?.username || '管理员'
})

const breadcrumbText = computed(() => {
  const title = String(route.meta.title ?? '工作台')
  return `首页 / ${title}`
})

const activeMenuKey = computed(() => {
  return menuItems.some((item) => item.to === route.path) ? route.path : null
})

const dropdownOptions: DropdownOption[] = [
  {
    label: '退出登录',
    key: 'logout',
    icon: () =>
      h(
        NIcon,
        null,
        {
          default: () => h(LogOutOutline),
        },
      ),
  },
]

function ensureCurrentTab() {
  const title = String(route.meta.title ?? '')
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
  openTabs.value = nextTabs.length > 0
    ? nextTabs
    : [{ title: '工作台', to: '/dashboard', closable: false }]

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

  openTabs.value = [
    { title: '工作台', to: '/dashboard', closable: false },
  ]

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

<template>
  <NLayout class="h-screen overflow-hidden bg-[#F5F7FA]" has-sider :native-scrollbar="false">
    <NLayoutSider
      inverted
      :width="240"
      :native-scrollbar="false"
      content-class="flex h-full flex-col"
      content-style="padding: 18px 16px; background: #111827;"
    >
      <button
        type="button"
        class="flex h-12 items-center gap-2.5 border-none bg-transparent px-0 text-left text-white"
        @click="navigateTo('/dashboard')"
      >
        <span class="h-7 w-7 rounded-[5px] bg-[#18A058]" />
        <span class="text-lg font-bold">EZ Admin</span>
      </button>

      <p class="mt-6 text-xs font-semibold tracking-wide text-[#6B7280]">主菜单</p>

      <NMenu
        class="mt-3"
        :value="activeMenuKey"
        :options="menuOptions"
        :indent="18"
        inverted
        @update:value="handleMenuUpdate"
      />
    </NLayoutSider>

    <NLayout class="h-screen min-w-0 overflow-hidden bg-[#F5F7FA]" :native-scrollbar="false">
      <NLayoutHeader
        class="flex h-14 items-center justify-between border-b border-[#E5E7EB] bg-white px-6"
      >
        <p class="text-sm text-[#374151]">{{ breadcrumbText }}</p>

        <div class="flex items-center gap-2.5">
          <NInput placeholder="搜索菜单 / 页面" clearable class="w-46">
            <template #prefix>
              <NIcon :component="SearchOutline" />
            </template>
          </NInput>

          <NButton quaternary circle>
            <template #icon>
              <NIcon :component="NotificationsOutline" />
            </template>
          </NButton>

          <NButton quaternary circle>
            <template #icon>
              <NIcon :component="ExpandOutline" />
            </template>
          </NButton>

          <NButton quaternary circle>
            <template #icon>
              <NIcon :component="MoonOutline" />
            </template>
          </NButton>

          <NDropdown trigger="click" :options="dropdownOptions" @select="handleUserAction">
            <NButton secondary>
              <template #icon>
                <NIcon :component="ChevronDownOutline" />
              </template>
              {{ displayName }}
            </NButton>
          </NDropdown>
        </div>
      </NLayoutHeader>

      <div class="flex h-10.5 items-center gap-2 border-b border-[#E5E7EB] bg-white px-4">
        <div class="flex min-w-0 flex-1 items-center gap-2 overflow-hidden">
          <button
            v-for="tab in openTabs"
            :key="tab.to"
            type="button"
            class="flex h-7 items-center justify-center rounded border px-4 text-[13px]"
            :class="
              route.path === tab.to
                ? 'border-[#18A058] bg-[#18A058] font-semibold text-white'
                : 'border-[#D9DEE8] bg-[#F9FAFB] text-[#374151]'
            "
            @click="navigateTo(tab.to)"
          >
            <span>{{ tab.title }}</span>
            <span
              v-if="tab.closable"
              class="ml-1 cursor-pointer"
              @click.stop="handleCloseTab(tab.to)"
            >
              ×
            </span>
          </button>
        </div>

        <div class="flex shrink-0 items-center gap-1">
          <NButton quaternary size="small" @click="handleRefresh">刷新</NButton>
          <NButton quaternary size="small" @click="handleCloseOtherTabs">关闭其他</NButton>
          <NButton quaternary circle size="small">
            <template #icon>
              <NIcon :component="EllipsisHorizontal" />
            </template>
          </NButton>
        </div>
      </div>

      <NLayoutContent
        :native-scrollbar="false"
        content-style="height: calc(100vh - 98px); padding: 32px; overflow: hidden; background: #F5F7FA;"
      >
        <RouterView />
      </NLayoutContent>
    </NLayout>
  </NLayout>
</template>
```

:::

::: details 为什么这里把滚动控制放在布局内部
后台壳子一旦采用“左侧固定 + 顶部固定 + 标签固定”的结构，就必须保证中间内容区能独立控制高度。

所以这份布局代码里有两个很关键的约束：

- 最外层 `NLayout` 使用 `h-screen overflow-hidden`，把整个后台限制在一屏。
- `NLayoutContent` 使用固定高度和 `overflow: hidden`，让业务页面在内部区域里渲染，不把浏览器撑出滚动条。

如果后面某个系统页面内容很多，再在 `RouterView` 内部页面自己的内容容器上加 `overflow-y-auto` 就够了。
:::

## 🛠️ 升级工作台页面

修改 `admin/src/pages/dashboard/DashboardHome.vue`。这一步不再继续堆“访问趋势”之类的演示占位，而是让首页直接读取当前项目的真实数据：

- 四个指标卡片直接放在最上面，先满足管理员“扫一眼看状态”的需求。
- 顶部问候区收成更克制的总览面板，展示当前用户、当天日期和系统整体状态，不再单独拆出健康检查说明卡。
- 右侧快捷入口只展示当前角色真正可访问的页面，并压缩成更短的信息块。
- 底部区域回显最近操作、最近登录和最新公告，方便做管理员日常巡检。

为了避免前端再次把接口字段写散，这里顺手把工作台接口类型和调用封装也补上。

::: details `admin/src/types/dashboard.ts` — 工作台接口类型

<<< ../../../admin/src/types/dashboard.ts

:::

::: details `admin/src/api/dashboard.ts` — 工作台接口封装

<<< ../../../admin/src/api/dashboard.ts

:::

::: details `admin/src/pages/dashboard/DashboardHome.vue` — 真实工作台页面

<<< ../../../admin/src/pages/dashboard/DashboardHome.vue

:::

## 🛠️ 创建可复用占位页

创建 `admin/src/pages/system/PlaceholderPage.vue`。这一页用于先把后台布局的路由出口接上，后续做到具体系统页时再逐步替换成真实内容。

::: details `admin/src/pages/system/PlaceholderPage.vue` — 可复用占位页

```vue
<script setup lang="ts">
import { NButton, NCard, NEmpty } from 'naive-ui'

defineProps<{
  title: string
  description: string
}>()
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-6 overflow-hidden">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[28px] font-bold text-[#111827]">{{ title }}</h1>
          <p class="mt-1 text-sm text-[#6B7280]">{{ description }}</p>
        </div>

        <NButton tertiary type="primary">
          后续接入
        </NButton>
      </div>

      <NCard
        class="min-h-0 flex-1 rounded-lg"
        :bordered="false"
        content-style="height: 100%;"
      >
        <div class="flex h-full items-center justify-center">
          <NEmpty description="本页会在后续小节继续补齐">
            <template #extra>
              <p class="text-sm text-[#6B7280]">当前先验证后台布局、路由出口和工作标签。</p>
            </template>
          </NEmpty>
        </div>
      </NCard>
    </section>
  </main>
</template>
```

:::

## 🛠️ 把后台布局接进路由

修改 `admin/src/router/index.ts`。这一处要做的核心变化有三件：

- 登录后不再直接进入单独页面，而是进入统一后台布局。
- 后台布局下挂工作台和静态系统页。
- `meta.title` 先补上，方便顶部栏和工作标签读取页面标题。

::: details `admin/src/router/index.ts` — 后台布局接进路由

```ts
import { createRouter, createWebHistory } from 'vue-router'

import { hasAccessToken } from '../utils/auth'

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
      beforeEnter: () => {
        if (hasAccessToken()) {
          return '/dashboard'
        }

        return true
      },
    },
    {
      path: '/',
      component: () => import('../layouts/AdminLayout.vue'), // [!code ++]
      beforeEnter: () => { // [!code ++]
        if (!hasAccessToken()) { // [!code ++]
          return '/login' // [!code ++]
        } // [!code ++]

        return true // [!code ++]
      }, // [!code ++]
      children: [ // [!code ++]
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('../pages/dashboard/DashboardHome.vue'),
          meta: { title: '工作台' }, // [!code ++]
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '用户管理',
            description: '这一页下一节会开始接入真实用户列表和操作表单。',
          },
          meta: { title: '用户管理' }, // [!code ++]
        },
        {
          path: 'roles',
          name: 'roles',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '角色权限',
            description: '当前先验证后台布局和标签栏，角色页面后续章节继续补齐。',
          },
          meta: { title: '角色权限' }, // [!code ++]
        },
        {
          path: 'menus',
          name: 'menus',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '菜单管理',
            description: '这一页下一节会开始与动态菜单能力衔接。',
          },
          meta: { title: '菜单管理' }, // [!code ++]
        },
        {
          path: 'logs',
          name: 'logs',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '操作日志',
            description: '当前先保留路由出口，后续章节再接真实日志页面。',
          },
          meta: { title: '操作日志' }, // [!code ++]
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '系统设置',
            description: '当前先验证后台布局结构，配置页后续章节继续补齐。',
          },
          meta: { title: '系统设置' }, // [!code ++]
        },
      ],
    },
  ],
})

export default router
```

:::

## ✅ 启动验证

先确认后端服务仍然运行，再进入 `admin/` 目录执行：

```bash
pnpm dev
```

浏览器打开终端输出的地址，然后用管理员账号登录。

### 1. 验证后台壳子是否贴近原型

登录成功后，应该看到：

- 左侧深色侧边栏固定在页面左边。
- 右侧顶部有面包屑、搜索框、图标按钮和用户区。
- 顶部栏下方有工作标签。
- 主内容区先显示顶部统计卡片，再往下展示真实工作台数据，而不是写死的演示数字。

### 2. 验证页面保持一屏，不出现浏览器默认滚动条

打开浏览器后，重点看两件事：

- 页面整体高度应该稳定在一屏内。
- 浏览器窗口右侧**不应该**出现默认滚动条。

如果出现了浏览器滚动条，优先检查：

- `html`、`body`、`#app` 是否真的加了 `overflow: hidden`
- 布局根节点是否用了 `h-screen overflow-hidden`
- 右侧主区域里是否缺少 `min-h-0`

### 3. 验证静态菜单和路由出口

依次点击左侧菜单里的：

- `工作台`
- `用户管理`
- `角色权限`
- `菜单管理`
- `操作日志`
- `系统设置`

这时应该看到：

- 顶部面包屑文字跟着变化。
- 工作标签会增加当前打开的页面。
- 内容区会切到对应页面，而不是整页刷新。

### 4. 验证工作台已经接入真实项目数据

回到 `/dashboard` 后，重点看下面几处：

- 指标卡片数值会跟数据库里当前的用户、角色、文件、公告和日志数据联动。
- 统计卡片会优先出现在页面最上方，便于快速扫描。
- `Overview` 会直接概括环境状态和最近刷新信息，不再单独占一块“健康检查入口”卡片。
- “最近操作”展示的是 `sys_operation_log` 中的最新记录。
- “最近登录”展示的是 `sys_login_log` 中的最新记录。
- “最新公告”展示的是启用中的公告，而不是前端手写数组。
- “快捷入口”只会展示当前登录角色真正有权访问的页面。

### 5. 验证工作标签的基础交互

继续验证标签栏：

- 点击某个标签，应该切回对应页面。
- 点击可关闭标签上的 `×`，应该关闭该标签。
- 点击“关闭其他”后，应该只保留“工作台”和当前页。
- 点击“刷新”会重新加载当前页。

::: warning ⚠️ 这里的“刷新”先用整页刷新实现
这一节先把标签栏交互跑通，所以 `handleRefresh` 直接用了 `window.location.reload()`。

这不是最终形态，但对当前教程阶段足够了。后面如果你想把刷新做成更细粒度的组件重载，再单独抽象也不晚。
:::

## 常见问题

::: details 登录后还是直接看到旧的简单工作台页
先确认 `router/index.ts` 里 `/dashboard` 是否已经改成挂在 `AdminLayout.vue` 的子路由下。

如果仍然是直接：

```ts
{
  path: '/dashboard',
  component: () => import('../pages/dashboard/DashboardHome.vue'),
}
```

那就说明后台壳子还没有真正接入。
:::

::: details 页面虽然进了后台，但浏览器还是出现滚动条
这是这一节最容易漏掉的点，优先按下面顺序检查：

1. `main.css` 的 `html`、`body`、`#app` 是否都加了 `overflow: hidden`
2. `AdminLayout.vue` 最外层是否用了 `h-screen overflow-hidden`
3. 右侧列和内容区是否用了 `min-h-0`

很多时候不是内容太多，而是某一层 flex 子节点没有允许自己收缩。
:::

::: details 顶部栏显示的用户名是“管理员”，不是登录昵称
先看 `getAuthUserInfo` 是否已经从 `USER_INFO_KEY` 里解析了本地用户信息。

如果没有这个函数，或者 JSON 解析失败后没有回退逻辑，顶部栏就拿不到真实昵称。
:::

下一节继续把左侧菜单从“前端写死”升级成“根据后端返回生成”：[动态菜单](./dynamic-menu)。
