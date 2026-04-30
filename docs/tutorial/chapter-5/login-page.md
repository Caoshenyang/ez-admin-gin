---
title: 登录页
description: "参考原型实现后台登录页，优化一屏布局，并打通登录接口、记住登录和当前路由跳转。"
---

# 登录页

这一节把登录入口接到真实 `/api/v1/auth/login`，并把原型里的左右双栏、验证码区、记住登录、默认账号提示和页脚一起落到页面中。

::: tip 🎯 本节目标
完成后，登录页在桌面端会接近原型图；页面高度稳定控制在一屏内；用户名和密码默认填充便于联调；验证码区域保留 UI 位置但不参与当前登录校验；勾选“记住登录”时把 Token 存到 `localStorage`，不勾选时只保存在 `sessionStorage`；当前路由守卫会负责登录页回跳和未登录拦截。
:::

![登录页原型](/prototypes/exports/myUgG.png)

::: info 这一节参考的原型
- 原型文件：`docs/public/prototypes/ui.pen`
- 节点名称：`01 登录页 / Naive UI Form`
- 导出图片：`docs/public/prototypes/exports/myUgG.png`
:::

::: warning ⚠️ 原型要参考，但接口能力要以教程现状为准
原型图里有验证码、记住登录和忘记密码入口，但第 3 章后端登录接口目前仍然只接收 `username`、`password` 两个字段。

所以这一节的处理方式是：

- 页面结构和视觉层级按原型实现。
- 验证码区域先保留，但当前不参与登录校验。
- “忘记密码”按钮当前只提示未接入找回密码流程。
- `LoginPage.vue` 当前默认填充值和提示文案是 `admin / Admin@123456`；如果你的后端仍沿用前文初始化示例 `EzAdmin@123456`，以你实际初始化时设置的密码为准。
:::

## 先看这一页要做成什么

为了后面代码更好跟，我们先把原型拆成三块：

| 区域 | 原型里有什么 | 本节怎么处理 |
| --- | --- | --- |
| 左侧品牌区 | 深色品牌面板、标题、副标题、四条能力说明 | 用 Tailwind CSS 4 组织页面骨架和背景层次 |
| 右侧登录卡片 | 用户名、密码、验证码、记住登录、忘记密码、登录按钮 | 登录接口只提交用户名和密码；验证码只保留 UI；忘记密码按钮当前只提示未接入 |
| 卡片下方补充信息 | 默认账号提示、页脚版权 | 一起实现，方便你验收页面完整度 |

## 本节样式分工

这一节继续遵守第五章的前端约束：

- Tailwind CSS 4 负责页面骨架、栅格、间距、背景和响应式布局。
- `NForm`、`NInput`、`NCheckbox`、`NButton`、`NCard`、`NAlert` 这些组件继续交给 Naive UI。
- 需要原型感的部分，优先通过 Tailwind 调整页面层级，不把组件重新拆成原生标签手写。

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
├─ vite.config.ts
└─ src/
   ├─ api/
   │  ├─ auth.ts
   │  └─ http.ts
   ├─ types/
   │  ├─ auth.ts
   │  └─ http.ts
   ├─ utils/
   │  └─ auth.ts
   ├─ pages/
   │  └─ auth/
   │     └─ LoginPage.vue
   └─ router/
      └─ index.ts
```

| 位置 | 用途 |
| --- | --- |
| `vite.config.ts` | 配置开发代理，让浏览器里的 `/api` 请求转发到后端 |
| `src/types/http.ts` | 定义统一响应结构类型 |
| `src/types/auth.ts` | 定义登录请求和登录响应类型 |
| `src/api/http.ts` | 创建 Axios 实例，并统一挂载认证请求头 |
| `src/api/auth.ts` | 封装登录接口 |
| `src/utils/auth.ts` | 按“记住登录”决定写入 `localStorage` 或 `sessionStorage` |
| `src/pages/auth/LoginPage.vue` | 用 Tailwind 4 组织登录页骨架，并用 Naive UI 承担表单和卡片交互 |
| `src/router/index.ts` | 让根路径和登录页根据本地 Token 做最小跳转 |

## 开始前先确认

开始前先确认两件事：

- 后端服务已经启动，访问 [http://localhost:8080/health](http://localhost:8080/health) 能返回健康检查结果。
- 当前已经完成上一节 [Vue 3 管理台初始化](./vue-project-init)。

::: warning ⚠️ 浏览器里直接请求 `http://localhost:8080` 会遇到跨域限制
前端开发服务默认运行在 `http://localhost:5173`，后端默认运行在 `http://localhost:8080`。如果不处理，浏览器会拦截跨域请求。

这一节使用 Vite 开发代理解决本地联调问题，所以一定要先完成下面的 `vite.config.ts` 修改，再重新启动 `pnpm dev`。
:::

## 🛠️ 配置开发代理

先修改 `admin/vite.config.ts`。本次只加一处：把 `/api` 请求代理到本地后端。

```ts
import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueDevTools(), tailwindcss()],
  server: { // [!code ++]
    proxy: { // [!code ++]
      '/api': { // [!code ++]
        target: 'http://localhost:8080', // [!code ++]
        changeOrigin: true, // [!code ++]
      }, // [!code ++]
    }, // [!code ++]
  }, // [!code ++]
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
```

这样前端里只需要请求 `/api/v1/...`，开发环境下 Vite 会自动转发到 `http://localhost:8080/api/v1/...`。

## 🛠️ 定义接口类型

创建 `admin/src/types/http.ts`。这个文件先定义后端统一响应结构，后面所有接口都可以复用。

```ts
// ApiResponse 对应后端统一响应结构。
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}
```

创建 `admin/src/types/auth.ts`。这里仍然只保留真实登录接口需要的字段，和后端保持一致。

```ts
// LoginRequest 对应登录接口请求体。
export interface LoginRequest {
  username: string
  password: string
}

// LoginResponse 对应登录接口 data 字段。
export interface LoginResponse {
  user_id: number
  username: string
  nickname: string
  access_token: string
  token_type: string
  expires_at: string
}
```

::: details 为什么验证码没有写进 `LoginRequest`
因为当前后端接口还没有验证码参数。如果现在把 `captcha` 也发给 `/api/v1/auth/login`，既没有实际收益，还会让前后端契约变得不一致。

所以当前代码只保留验证码的界面位置和刷新交互，不把它纳入请求体和前端校验。
:::

## 🛠️ 封装本地登录态存储

创建 `admin/src/utils/auth.ts`。这一层负责两件事：

- 统一管理 Token、Token 类型和用户信息。
- 根据“记住登录”决定把数据写进 `localStorage` 还是 `sessionStorage`。

::: details `admin/src/utils/auth.ts` — 本地登录态存储

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

::: details 为什么要同时支持 `localStorage` 和 `sessionStorage`
原型里有“记住登录”复选框，如果无论是否勾选都固定写入 `localStorage`，这个交互就只是一个摆设。

这里的约定是：

- 勾选“记住登录”：刷新浏览器、关闭后重新打开，登录态仍然保留。
- 不勾选“记住登录”：只在当前浏览器会话内保留，关闭标签页或浏览器后失效。
:::

## 🛠️ 创建 Axios 实例

创建 `admin/src/api/http.ts`。这一层负责两件事：

- 统一请求前缀和超时时间。
- 如果本地已经有 Token，就自动带上 `Authorization` 请求头。

::: details `admin/src/utils/request.ts` — Axios 封装

```ts
import axios from 'axios'

import { clearAuthSession, getAuthorizationHeader } from '../utils/auth'

const http = axios.create({
  // 通过 Vite 代理转发到本地后端。
  baseURL: '/api/v1',
  timeout: 10000,
})

http.interceptors.request.use((config) => {
  const authorization = getAuthorizationHeader()

  if (authorization) {
    config.headers.Authorization = authorization
  }

  return config
})

http.interceptors.response.use(
  (response) => response,
  (error) => {
    // 401 时清掉本地旧 Token，避免失效登录态残留。
    if (error.response?.status === 401) {
      clearAuthSession()
    }

    return Promise.reject(error)
  },
)

export default http
```

:::

## 🛠️ 封装登录接口

创建 `admin/src/api/auth.ts`。登录接口只做一件事：发送用户名和密码，并把后端的 `data` 取出来返回给页面。

```ts
import http from './http'

import type { LoginRequest, LoginResponse } from '../types/auth'
import type { ApiResponse } from '../types/http'

// login 调用后端登录接口。
export async function login(payload: LoginRequest) {
  const response = await http.post<ApiResponse<LoginResponse>>('/auth/login', payload)
  return response.data.data
}
```

## 🛠️ 实现登录页

修改 `admin/src/pages/auth/LoginPage.vue`。这次不是只做一个简单表单，而是要把原型里的主要元素都补齐。

本次重点看 6 个点：

- 页面壳子和响应式布局用 Tailwind 4 写。
- 页面整体高度限制在一屏内，不出现浏览器默认滚动条。
- 登录卡片和表单节奏整体收紧，避免左右两侧视觉比例失衡。
- 左侧是品牌介绍区，右侧是登录卡片。
- 表单、按钮、提示卡片继续使用 Naive UI 组件。
- 用户名和密码默认填充，方便当前阶段联调。
- 验证码先保留展示和刷新，不参与当前登录校验。
- “记住登录”决定 `setAuthSession` 写入哪种存储，“忘记密码”先保留入口。

这里直接整体替换成下面内容：

::: details `admin/src/pages/auth/LoginPage.vue` — 登录页

```vue
<script setup lang="ts">
import axios from 'axios'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NCheckbox,
  NForm,
  NFormItem,
  NInput,
  useMessage,
} from 'naive-ui'
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { login } from '../../api/auth'
import { hasAccessToken, setAuthSession } from '../../utils/auth'

const router = useRouter()
const message = useMessage()

const formRef = ref<FormInst | null>(null)
const submitting = ref(false)

const productFeatures = [
  '权限模型：用户 / 角色 / 菜单 / 按钮',
  '工作标签：多页面切换、刷新、关闭其他',
  '审计能力：登录日志、操作日志、风险等级',
  '工程友好：Gin API + Vue 页面快速扩展',
]

function createCaptcha() {
  const alphabet = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
  return Array.from({ length: 4 }, () => {
    const index = Math.floor(Math.random() * alphabet.length)
    return alphabet[index]
  }).join('')
}

const captchaText = ref(createCaptcha())

// 登录表单模型。用户名和密码先默认填充，方便当前阶段联调。
const formModel = reactive({
  username: 'admin',
  password: 'Admin@123456',
  captcha: '',
  rememberLogin: true,
})

const rules: FormRules = {
  username: [
    {
      required: true,
      message: '请输入用户名',
      trigger: ['blur', 'input'],
    },
  ],
  password: [
    {
      required: true,
      message: '请输入密码',
      trigger: ['blur', 'input'],
    },
  ],
}

const footerText = computed(() => {
  return `© ${new Date().getFullYear()} EZ Admin Gin · Naive UI Admin Template`
})

function refreshCaptcha() {
  captchaText.value = createCaptcha()
  formModel.captcha = ''
}

function handleForgotPassword() {
  message.info('当前版本暂未接入找回密码流程')
}

// 如果本地已经有 Token，就直接跳到工作台。
if (hasAccessToken()) {
  void router.replace('/dashboard')
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  submitting.value = true

  try {
    const result = await login({
      username: formModel.username.trim(),
      password: formModel.password,
    })

    setAuthSession(result, formModel.rememberLogin)
    message.success('登录成功')
    await router.push('/dashboard')
  } catch (error) {
    const errorMessage = axios.isAxiosError<{ message?: string }>(error)
      ? error.response?.data?.message ?? '登录失败，请稍后重试'
      : '登录失败，请稍后重试'

    message.error(errorMessage)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <main class="h-screen overflow-hidden bg-[#F5F7FA] px-4 py-4 md:px-5 md:py-5">
    <section
      class="mx-auto grid h-full max-w-[1180px] items-center gap-6 xl:grid-cols-[minmax(0,560px)_400px] xl:justify-between xl:gap-8"
    >
      <section
        class="flex max-h-[720px] min-h-0 flex-col justify-between overflow-hidden rounded-[20px] bg-[#111827] px-7 py-7 md:px-9 md:py-8 xl:px-10 xl:py-9"
      >
        <div>
          <div class="h-14 w-14 rounded-[14px] bg-[#18A058]" />
          <h1 class="mt-6 text-[38px] leading-[1.06] font-bold tracking-tight text-white md:text-[48px]">
          EZ Admin Gin
          </h1>
          <p class="mt-4 text-[15px] leading-7 text-[#D1D5DB] md:text-[17px]">
            面向工程团队的 Naive UI 后台框架
          </p>
        </div>

        <div class="mt-6 rounded-2xl bg-[#1F2937] p-5 md:p-6">
          <ul class="grid list-none gap-4 p-0">
            <li
              v-for="feature in productFeatures"
              :key="feature"
              class="text-[14px] leading-7 text-[#F9FAFB] md:text-[15px]"
            >
              {{ feature }}
            </li>
          </ul>
        </div>
      </section>

      <section class="flex min-h-0 flex-col justify-center gap-2">
        <NCard
          class="rounded-2xl shadow-[0_20px_60px_rgba(15,23,42,0.08)]"
          :bordered="false"
          content-style="padding: 20px;"
        >
          <div class="mb-2.5">
            <h2 class="mb-1 text-[23px] font-bold text-[#111827]">登录控制台</h2>
            <p class="text-sm text-[#6B7280]">请使用管理员账号继续</p>
          </div>

          <NForm
            ref="formRef"
            :model="formModel"
            :rules="rules"
            class="login-form"
            label-placement="top"
            size="medium"
            @submit.prevent="handleSubmit"
          >
            <NFormItem label="用户名" path="username">
              <NInput
                v-model:value="formModel.username"
                class="compact-input"
                placeholder="请输入用户名"
                autocomplete="username"
              />
            </NFormItem>

            <NFormItem label="密码" path="password" class="password-item">
              <NInput
                v-model:value="formModel.password"
                class="compact-input"
                type="password"
                show-password-on="click"
                placeholder="请输入密码"
                autocomplete="current-password"
              />
            </NFormItem>

            <NFormItem class="captcha-item mb-0">
              <div class="grid w-full gap-3 sm:grid-cols-[minmax(0,1fr)_120px]">
                <NInput
                  v-model:value="formModel.captcha"
                  class="compact-input"
                  placeholder="验证码"
                  maxlength="4"
                />

                <button
                  type="button"
                  class="h-8.5 cursor-pointer rounded-lg border border-[#A7F3D0] bg-[#ECFDF5] text-lg font-bold tracking-[0.08em] text-[#18A058]"
                  @click="refreshCaptcha"
                >
                  {{ captchaText }}
                </button>
              </div>
            </NFormItem>

            <div class="my-2.5 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
              <NCheckbox v-model:checked="formModel.rememberLogin">
                记住登录
              </NCheckbox>

              <button
                type="button"
                class="cursor-pointer border-none bg-transparent p-0 text-sm text-[#2080F0]"
                @click="handleForgotPassword"
              >
                忘记密码？
              </button>
            </div>

            <NButton
              attr-type="submit"
              type="primary"
              size="medium"
              block
              color="#18A058"
              :loading="submitting"
              class="login-submit"
            >
              登录
            </NButton>
          </NForm>

          <NAlert
            type="info"
            :show-icon="false"
            class="mt-2.5 compact-alert"
            title="默认账号：admin / Admin@123456"
          >
            验证码当前仅做界面占位。
          </NAlert>
        </NCard>

        <p class="px-1 text-[12px] text-[#9CA3AF]">{{ footerText }}</p>
      </section>
    </section>
  </main>
</template>

<style scoped>
.login-form {
  --n-feedback-height: 8px;
  --n-feedback-padding: 1px 0 0;
  --n-label-height: 18px;
  --n-label-padding: 0 0 3px;
}

.login-form :deep(.n-form-item) {
  margin-bottom: 4px;
}

.login-form :deep(.password-item) {
  margin-bottom: 0;
}

.login-form :deep(.password-item .n-form-item-feedback-wrapper) {
  min-height: 2px;
}

.login-form :deep(.captcha-item) {
  margin-top: -6px;
}

.login-form :deep(.n-form-item:last-child) {
  margin-bottom: 0;
}

.compact-input {
  --n-border-radius: 8px;
  --n-font-size: 14px;
  --n-height: 34px;
  --n-padding-left: 11px;
  --n-padding-right: 11px;
}

.login-submit {
  --n-border-radius: 8px;
  --n-font-size: 14px;
  --n-height: 36px;
}

.compact-alert {
  --n-border-radius: 8px;
  --n-font-size: 13px;
  --n-padding: 8px 10px;
}
</style>
```

:::

::: warning ⚠️ 这一节先不接复杂登录能力
这一节已经把原型图里的关键界面都放出来了，但真正接入的后端能力仍然只有登录接口本身。

所以当前行为边界是：

- `/login` 页可以登录、跳转、保存 Token。
- 验证码区域当前只做界面占位，不参与后端鉴权和前端提交拦截。
- “忘记密码”按钮当前只提示未接入找回密码流程。
- 如果你在未登录状态下手动输入 `/dashboard` 或其他后台页，当前全局守卫会重定向到 `/login`，并带上 `redirect` 查询参数。
:::

## 🛠️ 对齐当前路由跳转

修改 `admin/src/router/index.ts`。按当前仓库代码，登录相关跳转已经统一收敛到全局守卫：

- 访问根路径时，根据本地 Token 决定跳到 `/login` 还是 `/dashboard`。
- 如果已经登录，再访问 `/login` 时直接回到 `/dashboard`。
- 如果未登录访问后台页，统一跳回 `/login`，并保留原始目标地址。

::: details `admin/src/router/index.ts` — 路由守卫

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

先确认后端服务已经运行，再进入 `admin/` 目录执行：

```bash
# 启动管理台开发服务
pnpm dev
```

如果你之前已经开着开发服务，记得先停掉，再重新启动一次，确保新的 `vite.config.ts` 代理配置生效。

浏览器打开终端输出的地址，通常是：

```text
http://localhost:5173/
```

### 1. 验证页面结构是否对齐原型

打开登录页后，先不着急点登录，先看页面结构。

你应该能看到下面这些内容：

- 左侧有深色品牌面板，包含 `EZ Admin Gin` 标题、副标题和四条能力说明。
- 右侧有白色登录卡片，包含“登录控制台”标题。
- 卡片内有用户名、密码、验证码、记住登录、忘记密码和登录按钮。
- 登录按钮下方有默认账号提示，底部有版权文案。

如果你把浏览器宽度缩到移动端，页面应该自动折成上下结构，而不是左右挤压变形。

### 2. 验证默认联调值和验证码占位

打开页面后，应该能直接看到：

- 用户名默认是 `admin`
- 密码输入框当前默认填充值是 `Admin@123456`
- 验证码区域有展示值，也可以点击刷新
- 页面提示会说明“验证码区域当前仅做界面占位”

这一步的重点是确认当前阶段联调更顺手，而不是先把验证码校验做重。

### 3. 验证登录成功

输入管理员账号：

- 用户名：`admin`
- 密码：以你当前初始化结果为准。当前前端默认填充值是 `Admin@123456`，如果你仍沿用前文初始化示例，则可能还是 `EzAdmin@123456`。

勾选“记住登录”后点击“登录”，应该看到：

- 页面提示“登录成功”。
- 页面跳转到 `/dashboard`。
- `Network` 面板中出现：
  - 请求地址：`/api/v1/auth/login`
  - 请求方法：`POST`
  - 响应状态：`200`
  - 响应体中包含 `access_token`
- 浏览器 `Application -> Local Storage` 中出现：
  - `ez-admin-access-token`
  - `ez-admin-token-type`
  - `ez-admin-user-info`

### 4. 验证“记住登录”是否生效

清掉浏览器存储后，再做一次登录，但这次**取消勾选**“记住登录”。

这时应该看到：

- 登录成功后仍然会跳转到 `/dashboard`。
- 数据写进的是 `Application -> Session Storage`，而不是 `Local Storage`。

如果你关闭当前标签页或浏览器后重新打开，登录态应该失效；这说明“记住登录”不是摆设，而是真的影响了存储策略。

### 5. 验证后端错误透传

把密码故意改错，再登录一次。

这时应该看到：

- 页面停留在登录页，不会跳转。
- 页面提示后端返回的错误信息，例如“用户名或密码错误”。

这一步很重要，它可以确认前端没有把后端错误吞掉。

## 常见问题

::: details 浏览器提示跨域错误
先确认 `vite.config.ts` 已经增加了 `/api` 代理配置。

如果文件已经改了，但浏览器仍然报跨域，通常是因为 `pnpm dev` 还是旧进程。把开发服务停掉后重新启动，再刷新浏览器。
:::

::: details 为什么输入错误验证码也还能发起登录
这是当前阶段的预期行为。

这一节先把验证码区域保留下来，对齐原型图，但前端不会先用它拦登录。
:::

::: details 勾选“记住登录”后，Token 还是没有出现在 `Local Storage`
先看 `setAuthSession` 是否接收了第二个参数 `rememberLogin`，并且页面提交时传入的是 `formModel.rememberLogin`。

如果 `setAuthSession(result)` 仍然只有一个参数，就说明存储策略还没有真正接上。
:::

::: details 点击“忘记密码”没有跳转，是不是写漏了
不是。这一节只保留原型里的入口位置，让页面完整度先和原型对齐。

当前点击后只提示未接入找回密码流程，这是本节的预期行为。
:::

下一节开始搭建后台整体骨架：[后台布局](./admin-layout)。
