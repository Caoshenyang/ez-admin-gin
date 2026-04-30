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

::: warning ⚠️ 这一节只聚焦后台壳子本身
当前仓库代码已经把菜单接到动态路由链路，但这一节关注的重点仍然是一屏布局、顶部栏、标签栏和工作台骨架。

结合当前实现，可以先记住三点：

- `AdminLayout.vue` 左侧菜单已经改为读取 `sideMenuOptions`，不再维护本地静态菜单数组。
- 工作标签仍然按当前路由维护，这是后台壳子的职责。
- 登录拦截、动态路由注册和菜单清理已经放进 `router/index.ts`，下一节再专门展开这条链路。
:::

## 先明确这一节的布局约束

后台布局和登录页不一样，它会承载后续所有系统页面，所以先把壳子的规则定死，后面会省很多返工。

本节统一遵守下面几条约束：

- 页面整体高度限制在一屏内，根容器使用 `h-screen`。
- 不出现浏览器默认滚动条，`html`、`body`、`#app` 和布局根节点都要关闭浏览器级滚动。
- 页面超出一屏时，优先让内容区内部滚动，而不是让整个浏览器页面滚动。
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
   ├─ router/
   │  ├─ dynamic-menu.ts
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
| `src/layouts/AdminLayout.vue` | 实现侧边栏、顶部栏、工作标签和内容区，并消费动态菜单状态 |
| `src/pages/dashboard/DashboardHome.vue` | 把原来的演示工作台升级成真实项目首页 |
| `src/router/dynamic-menu.ts` | 提供侧边栏菜单状态、页面标题查找和按钮权限集合 |
| `src/router/index.ts` | 把后台布局、登录守卫和动态路由挂载放到同一条路由链路里 |

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

页面里仍然会保留少量 Tailwind 类，用来控制高度、留白、颜色和一屏约束。按当前仓库代码，`AdminLayout.vue` 主要承接四块：

- 左侧菜单区：读取 `sideMenuOptions`
- 顶部栏：面包屑、搜索框、快捷按钮和用户区
- 工作标签：按当前路由维护打开页签
- 路由内容区：承接工作台和动态注册页面

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
        :options="sideMenuOptions"
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

```ts
export interface DashboardCurrentUser {
  user_id: number
  username: string
  nickname: string
}

export interface DashboardHealth {
  env: string
  database: string
  redis: string
}

export interface DashboardMetrics {
  user_total: number
  enabled_user_total: number
  enabled_role_total: number
  config_total: number
  notice_total: number
  file_total: number
  today_operation_total: number
  today_risk_operation_total: number
  today_login_failed_total: number
}

export interface DashboardOperationItem {
  id: number
  username: string
  method: string
  path: string
  status_code: number
  success: boolean
  latency_ms: number
  created_at: string
}

export const DashboardLoginStatus = {
  Success: 1,
  Failed: 2,
} as const

export type DashboardLoginStatus =
  (typeof DashboardLoginStatus)[keyof typeof DashboardLoginStatus]

export interface DashboardLoginItem {
  id: number
  username: string
  status: DashboardLoginStatus
  message: string
  ip: string
  created_at: string
}

export interface DashboardNoticeItem {
  id: number
  title: string
  status: number
  updated_at: string
}

export interface DashboardData {
  current_user: DashboardCurrentUser
  health: DashboardHealth
  metrics: DashboardMetrics
  recent_operations: DashboardOperationItem[]
  recent_logins: DashboardLoginItem[]
  latest_notices: DashboardNoticeItem[]
}
```

:::

::: details `admin/src/api/dashboard.ts` — 工作台接口封装

```ts
import http from './http'

import type { DashboardData } from '../types/dashboard'
import type { ApiResponse } from '../types/http'

export async function getDashboardSummary() {
  const response = await http.get<ApiResponse<DashboardData>>('/auth/dashboard')
  return response.data.data
}
```

:::

::: details `admin/src/pages/dashboard/DashboardHome.vue` — 真实工作台页面

```vue
<script setup lang="ts">
import {
  LayersOutline,
  PulseOutline,
  ShieldCheckmarkOutline,
  TimeOutline,
} from '@vicons/ionicons5'
import { NAlert, NButton, NCard, NEmpty, NIcon, NTag } from 'naive-ui'
import type { Component } from 'vue'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { getDashboardSummary } from '../../api/dashboard'
import { authMenus } from '../../router/dynamic-menu'
import { DashboardLoginStatus, type DashboardData } from '../../types/dashboard'
import { MenuType, type AuthMenu } from '../../types/menu'

interface MetricCard {
  label: string
  value: string
  hint: string
  accent: string
  panelClass: string
  icon: Component
}

interface QuickLink {
  title: string
  path: string
  description: string
}

const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')
const dashboard = ref<DashboardData | null>(null)
const refreshedAt = ref('')

const currentDateLabel = computed(() => {
  return new Intl.DateTimeFormat('zh-CN', {
    month: 'long',
    day: 'numeric',
    weekday: 'long',
  }).format(new Date())
})

const currentUserLabel = computed(() => {
  const user = dashboard.value?.current_user
  if (!user) {
    return '管理员'
  }

  return user.nickname || user.username
})

const visiblePageMenus = computed(() => {
  return flattenPageMenus(authMenus.value)
})

const visiblePageTotal = computed(() => {
  return visiblePageMenus.value.length + 1
})

const healthPath = computed(() => findMenuPathByTitle(visiblePageMenus.value, '系统状态'))
const userManagePath = computed(() => findMenuPathByTitle(visiblePageMenus.value, '用户管理'))

const isHealthy = computed(() => {
  const health = dashboard.value?.health
  if (!health) {
    return false
  }

  return health.database === 'ok' && health.redis === 'ok'
})

const heroStatusText = computed(() => {
  const health = dashboard.value?.health
  if (!health) {
    return loading.value ? '正在拉取当前项目的实时概览...' : '等待首次同步项目运行数据'
  }

  if (health.database === 'ok' && health.redis === 'ok') {
    return '环境、数据库和缓存依赖都在线，可以直接作为管理员首页使用。'
  }

  if (health.database === 'ok') {
    return '数据库正常，Redis 有异常信号，建议优先检查缓存和登录态链路。'
  }

  return '核心依赖出现异常，请先处理环境问题，再继续做业务操作。'
})

const metricCards = computed<MetricCard[]>(() => {
  const metrics = dashboard.value?.metrics

  return [
    {
      label: '启用账号',
      value: formatMetricValue(metrics?.enabled_user_total),
      hint: `总账号 ${formatMetricValue(metrics?.user_total)}`,
      accent: '#1677ff',
      panelClass: 'bg-[#eff6ff]',
      icon: ShieldCheckmarkOutline,
    },
    {
      label: '启用角色',
      value: formatMetricValue(metrics?.enabled_role_total),
      hint: `可访问页面 ${formatMetricValue(visiblePageTotal.value)}`,
      accent: '#0f766e',
      panelClass: 'bg-[#ecfeff]',
      icon: LayersOutline,
    },
    {
      label: '今日操作',
      value: formatMetricValue(metrics?.today_operation_total),
      hint: `失败操作 ${formatMetricValue(metrics?.today_risk_operation_total)}`,
      accent: '#b45309',
      panelClass: 'bg-[#fff7ed]',
      icon: TimeOutline,
    },
    {
      label: '文件沉淀',
      value: formatMetricValue(metrics?.file_total),
      hint: `公告 ${formatMetricValue(metrics?.notice_total)} / 配置 ${formatMetricValue(metrics?.config_total)}`,
      accent: '#047857',
      panelClass: 'bg-[#ecfdf5]',
      icon: PulseOutline,
    },
  ]
})

const healthItems = computed(() => {
  const health = dashboard.value?.health

  return [
    {
      label: '运行环境',
      value: health?.env || 'unknown',
      status: health ? 'ok' : 'pending',
    },
    {
      label: '数据库',
      value: health?.database || 'pending',
      status: health?.database || 'pending',
    },
    {
      label: 'Redis',
      value: health?.redis || 'pending',
      status: health?.redis || 'pending',
    },
  ]
})

const quickLinks = computed<QuickLink[]>(() => {
  return visiblePageMenus.value
    .filter((menu) => menu.path && menu.path !== '/dashboard')
    .slice(0, 6)
    .map((menu) => ({
      title: menu.title,
      path: menu.path,
      description: getQuickLinkDescription(menu.title),
    }))
})

const recentOperations = computed(() => {
  return dashboard.value?.recent_operations ?? []
})

const recentLogins = computed(() => {
  return dashboard.value?.recent_logins ?? []
})

const latestNotices = computed(() => {
  return dashboard.value?.latest_notices ?? []
})

const refreshedLabel = computed(() => {
  return refreshedAt.value ? formatDateTime(refreshedAt.value) : '尚未同步'
})

const focusFacts = computed(() => {
  return [
    {
      label: '可访问页面',
      value: formatMetricValue(visiblePageTotal.value),
      hint: '按当前角色实时计算',
    },
    {
      label: '失败登录',
      value: formatMetricValue(dashboard.value?.metrics.today_login_failed_total),
      hint: '今日累计失败次数',
    },
    {
      label: '最近刷新',
      value: refreshedLabel.value,
      hint: '已同步当前项目快照',
    },
  ]
})

function flattenPageMenus(menus: AuthMenu[]) {
  const result: AuthMenu[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Menu && menu.path) {
      result.push(menu)
    }

    result.push(...flattenPageMenus(menu.children ?? []))
  }

  return result
}

function findMenuPathByTitle(menus: AuthMenu[], title: string) {
  return menus.find((menu) => menu.title === title)?.path || ''
}

function getQuickLinkDescription(title: string) {
  const descriptionMap: Record<string, string> = {
    系统状态: '检查当前环境、数据库和 Redis 状态',
    用户管理: '维护后台账号、状态和角色绑定',
    角色管理: '调整角色权限和菜单分配',
    菜单管理: '维护动态菜单、按钮权限和路由出口',
    系统配置: '查看系统参数和运行期配置项',
    文件管理: '检查上传文件和资源沉淀',
    操作日志: '回看后台操作链路和失败记录',
    登录日志: '检查最近登录结果和来源 IP',
    公告管理: '维护首页公告和系统通知内容',
  }

  return descriptionMap[title] || '进入对应系统页面继续处理业务'
}

function formatMetricValue(value?: number) {
  return typeof value === 'number' ? new Intl.NumberFormat('zh-CN').format(value) : '--'
}

function formatDateTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

function formatRoutePath(path: string) {
  return path.replace(/^\/api\/v1/, '') || path
}

function getHealthTagType(value: string) {
  return value === 'ok' ? 'success' : value === 'pending' ? 'default' : 'error'
}

function getStatusTagType(success: boolean) {
  return success ? 'success' : 'error'
}

function getLoginStatusTagType(status: number) {
  return status === DashboardLoginStatus.Success ? 'success' : 'error'
}

function getLoginStatusLabel(status: number) {
  return status === DashboardLoginStatus.Success ? '成功' : '失败'
}

function getErrorMessage(error: unknown) {
  if (typeof error === 'object' && error !== null) {
    const response = (error as { response?: { data?: { message?: string } } }).response
    if (typeof response?.data?.message === 'string' && response.data.message) {
      return response.data.message
    }
  }

  return '工作台数据获取失败，请稍后重试。'
}

async function loadDashboard() {
  loading.value = true
  errorMessage.value = ''

  try {
    dashboard.value = await getDashboardSummary()
    refreshedAt.value = new Date().toISOString()
  } catch (error) {
    errorMessage.value = getErrorMessage(error)
  } finally {
    loading.value = false
  }
}

function navigateTo(path: string) {
  if (!path) {
    return
  }

  void router.push(path)
}

onMounted(() => {
  void loadDashboard()
})
</script>

<template>
  <main class="flex h-full flex-col gap-4 overflow-auto pr-1">
    <NAlert
      v-if="errorMessage"
      type="error"
      title="工作台同步失败"
      class="rounded-2xl"
      :bordered="false"
    >
      {{ errorMessage }}
    </NAlert>

    <section class="flex flex-wrap items-start justify-between gap-4">
      <div class="max-w-[760px]">
        <p class="text-xs font-semibold uppercase tracking-[0.28em] text-[#94a3b8]">
          Project Workbench
        </p>
        <h1 class="mt-2 text-[30px] font-semibold tracking-[-0.03em] text-[#0f172a]">
          {{ currentUserLabel }}，今天是 {{ currentDateLabel }}
        </h1>
        <p class="mt-3 text-sm leading-7 text-[#475569]">
          先看关键指标，再看异常和待办，这是后台首页最顺手的阅读顺序。
          当前首页直接读取真实项目数据，适合登录后第一眼判断系统是否正常可用。
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <NTag
          :type="isHealthy ? 'success' : 'warning'"
          size="large"
          round
          :bordered="false"
        >
          {{ isHealthy ? '运行稳定' : '存在待检查项' }}
        </NTag>
        <NButton
          v-if="userManagePath"
          type="primary"
          color="#1677ff"
          @click="navigateTo(userManagePath)"
        >
          进入用户管理
        </NButton>
        <NButton v-else-if="healthPath" secondary @click="navigateTo(healthPath)">
          查看系统状态
        </NButton>
        <NButton :loading="loading" @click="void loadDashboard()">刷新工作台</NButton>
      </div>
    </section>

    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <NCard
        v-for="item in metricCards"
        :key="item.label"
        class="rounded-[24px]"
        :bordered="false"
        content-style="padding: 0;"
      >
        <div class="metric-card px-5 py-5" :class="item.panelClass">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-sm font-medium text-[#475569]">{{ item.label }}</p>
              <p class="mt-3 text-[32px] font-semibold tracking-[-0.03em] text-[#0f172a]">
                {{ item.value }}
              </p>
              <p class="mt-2 text-sm" :style="{ color: item.accent }">{{ item.hint }}</p>
            </div>
            <div
              class="flex h-11 w-11 items-center justify-center rounded-2xl bg-white/88 shadow-[0_10px_24px_rgba(15,23,42,0.06)]"
              :style="{ color: item.accent }"
            >
              <NIcon :component="item.icon" :size="20" />
            </div>
          </div>
        </div>
      </NCard>
    </section>

    <section class="grid gap-4 xl:grid-cols-[minmax(0,1.5fr)_360px]">
      <NCard class="rounded-[28px]" :bordered="false" content-style="padding: 0;">
        <div class="overview-panel px-6 py-6">
          <div class="flex flex-wrap items-start justify-between gap-5">
            <div class="max-w-[680px]">
              <div class="flex items-center gap-3">
                <div
                  class="flex h-12 w-12 items-center justify-center rounded-2xl bg-[#0f172a] text-white"
                >
                  <NIcon :component="PulseOutline" :size="22" />
                </div>
                <div>
                  <p class="text-sm font-semibold uppercase tracking-[0.24em] text-[#94a3b8]">
                    Overview
                  </p>
                  <p class="mt-1 text-[24px] font-semibold tracking-[-0.03em] text-[#0f172a]">
                    {{ heroStatusText }}
                  </p>
                </div>
              </div>

              <div class="mt-5 flex flex-wrap gap-2">
                <div
                  v-for="item in healthItems"
                  :key="item.label"
                  class="rounded-2xl border border-white/70 bg-white/82 px-4 py-3"
                >
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-semibold uppercase tracking-[0.18em] text-[#94a3b8]">
                      {{ item.label }}
                    </span>
                    <NTag
                      :type="getHealthTagType(item.status)"
                      size="small"
                      round
                      :bordered="false"
                    >
                      {{ item.value }}
                    </NTag>
                  </div>
                </div>
              </div>

              <div class="mt-5 flex flex-wrap items-center gap-2 text-sm text-[#64748b]">
                <span class="rounded-full bg-white/82 px-3 py-1.5">
                  当前身份 {{ dashboard?.current_user.username || 'waiting' }}
                </span>
                <span class="rounded-full bg-white/82 px-3 py-1.5">
                  健康页 {{ healthPath ? '已接入' : '未开放' }}
                </span>
              </div>
            </div>

            <div class="grid min-w-[220px] gap-3">
              <div
                v-for="item in focusFacts"
                :key="item.label"
                class="rounded-2xl bg-white/82 px-4 py-4 shadow-[0_10px_24px_rgba(15,23,42,0.05)]"
              >
                <p class="text-xs font-semibold uppercase tracking-[0.18em] text-[#94a3b8]">
                  {{ item.label }}
                </p>
                <p class="mt-2 text-lg font-semibold text-[#0f172a]">{{ item.value }}</p>
                <p class="mt-1 text-sm leading-6 text-[#64748b]">{{ item.hint }}</p>
              </div>
            </div>
          </div>
        </div>
      </NCard>

      <section class="grid gap-4">
        <NCard class="rounded-[24px]" :bordered="false" content-style="padding: 22px;">
          <div class="flex items-center justify-between gap-4">
            <div>
              <p class="text-sm font-semibold text-[#111827]">快捷入口</p>
              <p class="mt-1 text-sm text-[#6B7280]">只保留当前角色最常用的几个落点。</p>
            </div>
            <NTag round :bordered="false" type="info">{{ quickLinks.length }} 项</NTag>
          </div>

          <div v-if="quickLinks.length > 0" class="mt-4 grid gap-3">
            <button
              v-for="item in quickLinks"
              :key="item.path"
              type="button"
              class="quick-link-button"
              @click="navigateTo(item.path)"
            >
              <div class="min-w-0">
                <p class="text-sm font-semibold text-[#111827]">{{ item.title }}</p>
                <p class="mt-1 truncate text-sm text-[#64748B]">{{ item.description }}</p>
              </div>
              <span class="shrink-0 text-xs text-[#94A3B8]">进入</span>
            </button>
          </div>

          <NEmpty v-else class="mt-4" description="当前角色没有额外页面可跳转" />
        </NCard>
      </section>
    </section>

    <section class="grid gap-4 xl:grid-cols-[minmax(0,1.3fr)_minmax(320px,0.95fr)]">
      <NCard class="rounded-[24px]" :bordered="false" content-style="padding: 22px;">
        <div class="flex items-center justify-between gap-4">
          <div>
            <p class="text-sm font-semibold text-[#111827]">最近操作</p>
            <p class="mt-1 text-sm text-[#6B7280]">直接读取 `sys_operation_log`，用于判断后台最近发生了什么。</p>
          </div>
          <NTag round :bordered="false" type="info">{{ recentOperations.length }} 条</NTag>
        </div>

        <div v-if="recentOperations.length > 0" class="mt-4 space-y-3">
          <article
            v-for="item in recentOperations"
            :key="item.id"
            class="rounded-2xl border border-[#e5e7eb] px-4 py-3"
          >
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div class="flex min-w-0 items-center gap-2">
                <NTag
                  :type="getStatusTagType(item.success)"
                  size="small"
                  round
                  :bordered="false"
                >
                  {{ item.success ? '成功' : '失败' }}
                </NTag>
                <NTag size="small" round :bordered="false">{{ item.method }}</NTag>
                <span class="truncate text-sm font-medium text-[#111827]">
                  {{ item.username || '系统' }} · {{ formatRoutePath(item.path) }}
                </span>
              </div>
              <span class="text-sm text-[#64748B]">{{ formatDateTime(item.created_at) }}</span>
            </div>

            <div class="mt-3 flex flex-wrap items-center gap-4 text-sm text-[#64748B]">
              <span>状态码 {{ item.status_code }}</span>
              <span>耗时 {{ item.latency_ms }} ms</span>
            </div>
          </article>
        </div>

        <NEmpty v-else class="mt-4" description="还没有操作日志" />
      </NCard>

      <section class="grid gap-4">
        <NCard class="rounded-[24px]" :bordered="false" content-style="padding: 22px;">
          <div class="flex items-center justify-between gap-4">
            <div>
              <p class="text-sm font-semibold text-[#111827]">最近登录</p>
              <p class="mt-1 text-sm text-[#6B7280]">帮助你快速判断是否存在连续失败登录。</p>
            </div>
            <NTag round :bordered="false" type="warning">
              失败 {{ formatMetricValue(dashboard?.metrics.today_login_failed_total) }}
            </NTag>
          </div>

          <div v-if="recentLogins.length > 0" class="mt-4 space-y-3">
            <article
              v-for="item in recentLogins"
              :key="item.id"
              class="rounded-2xl bg-[#f8fafc] px-4 py-3"
            >
              <div class="flex items-center justify-between gap-3">
                <div class="min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-semibold text-[#111827]">{{ item.username }}</span>
                    <NTag
                      :type="getLoginStatusTagType(item.status)"
                      size="small"
                      round
                      :bordered="false"
                    >
                      {{ getLoginStatusLabel(item.status) }}
                    </NTag>
                  </div>
                  <p class="mt-1 truncate text-sm text-[#64748B]">
                    {{ item.message || '登录状态已记录' }}
                  </p>
                </div>
                <span class="text-xs text-[#94A3B8]">{{ item.ip || '-' }}</span>
              </div>
              <p class="mt-3 text-xs text-[#94A3B8]">{{ formatDateTime(item.created_at) }}</p>
            </article>
          </div>

          <NEmpty v-else class="mt-4" description="还没有登录记录" />
        </NCard>

        <NCard class="rounded-[24px]" :bordered="false" content-style="padding: 22px;">
          <div class="flex items-center justify-between gap-4">
            <div>
              <p class="text-sm font-semibold text-[#111827]">最新公告</p>
              <p class="mt-1 text-sm text-[#6B7280]">展示启用中的公告，便于确认系统对外通知是否最新。</p>
            </div>
            <NTag round :bordered="false" type="success">
              {{ formatMetricValue(dashboard?.metrics.notice_total) }} 条
            </NTag>
          </div>

          <div v-if="latestNotices.length > 0" class="mt-4 space-y-3">
            <article
              v-for="item in latestNotices"
              :key="item.id"
              class="rounded-2xl border border-[#e5e7eb] px-4 py-3"
            >
              <p class="text-sm font-semibold text-[#111827]">{{ item.title }}</p>
              <p class="mt-2 text-xs uppercase tracking-[0.18em] text-[#94A3B8]">
                updated {{ formatDateTime(item.updated_at) }}
              </p>
            </article>
          </div>

          <NEmpty v-else class="mt-4" description="当前没有启用中的公告" />
        </NCard>
      </section>
    </section>
  </main>
</template>

<style scoped>
.overview-panel {
  background:
    radial-gradient(circle at top left, rgba(255, 255, 255, 0.82) 0%, rgba(255, 255, 255, 0) 25%),
    radial-gradient(circle at right center, rgba(22, 119, 255, 0.1) 0%, rgba(22, 119, 255, 0) 34%),
    linear-gradient(135deg, #f3f8ff 0%, #f8fbff 52%, #edf7f0 100%);
}

.quick-link-button {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  width: 100%;
  border: 1px solid #e2e8f0;
  border-radius: 18px;
  background: #f8fafc;
  padding: 14px 16px;
  text-align: left;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.quick-link-button:hover {
  border-color: #bfdbfe;
  background: #ffffff;
  box-shadow: 0 14px 28px rgba(15, 23, 42, 0.08);
  transform: translateY(-1px);
}

.metric-card {
  min-height: 152px;
}
</style>
```

:::

## 🛠️ 路由里挂上后台壳子

当前仓库里，后台布局并不是“静态路由 + 占位页”那一版，而是已经并入统一的全局守卫链路：

- `AdminLayout.vue` 固定挂在 `name: 'admin'` 的父路由下。
- `dashboard` 作为内置子路由常驻存在。
- 其余页面在首次进入后台时，通过 `/api/v1/auth/menus` 动态注册。
- `PlaceholderPage.vue` 仍然保留，但只作为未知 `component` 编码的兜底页，不再是本章默认路由出口。

::: details `admin/src/router/index.ts` — 后台布局接进路由

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

### 3. 验证侧边栏、标签和路由出口

登录后，左侧应该能看到 `工作台` 和当前账号有权访问的菜单项。依次点击其中几个页面，重点观察：

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

下一节继续展开当前这套路由链路里的动态菜单部分，包括 `/api/v1/auth/menus`、运行时路由注册和按钮权限集合：[动态菜单](./dynamic-menu)。
