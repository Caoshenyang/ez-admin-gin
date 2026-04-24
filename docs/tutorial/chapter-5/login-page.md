---
title: 登录页
description: "实现后台登录页，并打通登录接口和 Token 保存。"
---

# 登录页

这一节把上一节的占位登录页替换成真正可用的登录入口：接入开发代理、调用后端登录接口、保存 Token，并在登录成功后跳转到工作台页面。

::: tip 🎯 本节目标
完成后，浏览器中可以直接输入管理员账号密码完成登录；登录成功后，前端会把 Token 保存到本地，并跳转到 `/dashboard`。
:::

::: info 本节先完成登录主链路
这一节重点是“能登录进去”。完整的后台布局、登录态守卫和动态菜单，会在后面几节继续补齐。
:::

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
| `src/utils/auth.ts` | 读写本地 Token 和登录用户信息 |
| `src/pages/auth/LoginPage.vue` | 实现登录表单和提交逻辑 |
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

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueDevTools()],
  server: {
    proxy: {
      // 浏览器访问 /api 时，由 Vite 转发到本地后端。
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
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

创建 `admin/src/types/auth.ts`。这里放登录请求和登录响应类型。

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

## 🛠️ 封装本地登录态存储

创建 `admin/src/utils/auth.ts`。这一层先统一管理 Token 和用户信息，后面做路由守卫、顶部用户信息、退出登录时都能直接复用。

```ts
import type { LoginResponse } from '../types/auth'

const ACCESS_TOKEN_KEY = 'ez-admin-access-token'
const TOKEN_TYPE_KEY = 'ez-admin-token-type'
const USER_INFO_KEY = 'ez-admin-user-info'

export interface AuthUserInfo {
  userId: number
  username: string
  nickname: string
  expiresAt: string
}

// setAuthSession 在登录成功后保存本地登录态。
export function setAuthSession(payload: LoginResponse) {
  localStorage.setItem(ACCESS_TOKEN_KEY, payload.access_token)
  localStorage.setItem(TOKEN_TYPE_KEY, payload.token_type)
  localStorage.setItem(
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
}

export function getAccessToken() {
  return localStorage.getItem(ACCESS_TOKEN_KEY) ?? ''
}

export function getTokenType() {
  return localStorage.getItem(TOKEN_TYPE_KEY) ?? 'Bearer'
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

::: details 为什么现在就单独抽出 `auth.ts`
如果把 `localStorage` 读写直接散落到页面组件里，后面做路由守卫、退出登录、自动带 Token 请求头时就会开始重复。现在先把读写入口统一起来，后面扩展会轻松很多。
:::

## 🛠️ 创建 Axios 实例

创建 `admin/src/api/http.ts`。这一层负责两件事：

- 统一请求前缀和超时时间。
- 如果本地已经有 Token，就自动带上 `Authorization` 请求头。

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
    // 后面做完整登录态守卫前，先在 401 时清掉本地旧 Token。
    if (error.response?.status === 401) {
      clearAuthSession()
    }

    return Promise.reject(error)
  },
)

export default http
```

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

修改 `admin/src/pages/auth/LoginPage.vue`。这里直接整体替换成下面内容。

```vue
<script setup lang="ts">
import axios from 'axios'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NButton,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NSpace,
  NText,
  useMessage,
} from 'naive-ui'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { login } from '../../api/auth'
import { hasAccessToken, setAuthSession } from '../../utils/auth'

const router = useRouter()
const message = useMessage()

const formRef = ref<FormInst | null>(null)
const submitting = ref(false)

// 登录表单模型。
const formModel = reactive({
  username: '',
  password: '',
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

    setAuthSession(result)
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
  <main class="login-page">
    <section class="login-panel">
      <div class="login-copy">
        <p class="login-eyebrow">EZ Admin</p>
        <h1 class="login-title">通用后台管理系统底座</h1>
        <p class="login-description">
          从这一节开始，前端正式接入后端认证接口。先把登录链路打通，后面的布局、菜单和系统页面才有落脚点。
        </p>
      </div>

      <NCard class="login-card" :bordered="false">
        <div class="login-card-header">
          <h2>账号登录</h2>
          <NText depth="3">请输入管理员账号和密码。</NText>
        </div>

        <NForm
          ref="formRef"
          :model="formModel"
          :rules="rules"
          label-placement="top"
          size="large"
          @submit.prevent="handleSubmit"
        >
          <NFormItem label="用户名" path="username">
            <NInput
              v-model:value="formModel.username"
              placeholder="请输入用户名"
              autocomplete="username"
            />
          </NFormItem>

          <NFormItem label="密码" path="password">
            <NInput
              v-model:value="formModel.password"
              type="password"
              show-password-on="click"
              placeholder="请输入密码"
              autocomplete="current-password"
            />
          </NFormItem>

          <NSpace vertical :size="16">
            <NButton
              attr-type="submit"
              type="primary"
              size="large"
              block
              :loading="submitting"
            >
              登录
            </NButton>

            <NText depth="3" class="login-tip">
              默认管理员账号：`admin`，密码：`EzAdmin@123456`
            </NText>
          </NSpace>
        </NForm>
      </NCard>
    </section>
  </main>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  padding: 24px;
  background:
    radial-gradient(circle at top left, rgba(24, 160, 88, 0.12), transparent 36%),
    linear-gradient(180deg, #f6fbf8 0%, #eef3f9 100%);
}

.login-panel {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(320px, 420px);
  align-items: center;
  gap: 48px;
  min-height: 100vh;
  max-width: 1120px;
  margin: 0 auto;
}

.login-copy {
  max-width: 560px;
}

.login-eyebrow {
  margin: 0 0 16px;
  font-size: 18px;
  font-weight: 700;
  color: #18a058;
}

.login-title {
  margin: 0;
  font-size: clamp(40px, 5vw, 60px);
  line-height: 1.08;
  color: #1f2937;
}

.login-description {
  margin: 24px 0 0;
  font-size: 18px;
  line-height: 1.8;
  color: #4b5563;
}

.login-card {
  width: 100%;
  border-radius: 8px;
  box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
}

.login-card-header {
  margin-bottom: 24px;
}

.login-card-header h2 {
  margin: 0 0 8px;
  font-size: 28px;
  color: #1f2937;
}

.login-tip {
  display: block;
  line-height: 1.7;
}

@media (max-width: 900px) {
  .login-panel {
    grid-template-columns: 1fr;
    gap: 32px;
    padding: 48px 0;
  }
}
</style>
```

::: warning ⚠️ 这里先不在页面里处理复杂登录状态
这一节只在“登录成功后跳转”和“本地已有 Token 时跳过登录页”这两个点上做最小处理。

如果你在登录前直接手输 `/dashboard`，目前仍然能看到占位工作台页面，这不是本节遗漏，而是完整路由守卫还没接入。后面做后台布局时会一起补上。
:::

## 🛠️ 调整路由最小跳转

修改 `admin/src/router/index.ts`。本次主要加两个点：

- 访问根路径时，根据本地 Token 决定跳到 `/login` 还是 `/dashboard`
- 如果已经登录，再访问 `/login` 时直接回到 `/dashboard`

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
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../pages/dashboard/DashboardHome.vue'),
    },
  ],
})

export default router
```

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

现在按下面顺序验证：

1. 打开登录页，确认页面能正常显示用户名和密码输入框。
2. 输入管理员账号 `admin` 和密码 `EzAdmin@123456`。
3. 点击“登录”后，页面应该提示“登录成功”，并跳转到 `/dashboard`。
4. 打开浏览器开发者工具，在 `Application -> Local Storage` 中，应该能看到：
   - `ez-admin-access-token`
   - `ez-admin-token-type`
   - `ez-admin-user-info`
5. 打开 `Network` 面板，应该能看到：
   - 请求地址：`/api/v1/auth/login`
   - 请求方法：`POST`
   - 响应状态：`200`
   - 响应体中包含 `access_token`

## ✅ 失败验证

把密码故意改错，再登录一次。

这时应该看到两件事：

- 页面停留在登录页，不会跳转。
- 页面提示后端返回的错误信息，例如“用户名或密码错误”。

这一步很重要，它可以确认前端没有把后端错误吞掉。

## 常见问题

::: details 浏览器提示跨域错误
先确认 `vite.config.ts` 已经增加了 `/api` 代理配置。

如果文件已经改了，但浏览器仍然报跨域，通常是因为 `pnpm dev` 还是旧进程。把开发服务停掉后重新启动，再刷新浏览器。
:::

::: details 登录按钮一直转圈
先看浏览器开发者工具的 `Network` 面板。

- 如果请求根本没发出去，先看控制台是否有前端报错。
- 如果请求一直挂起，先确认后端服务是否真的启动成功。
- 如果返回了 `500`，优先看后端终端日志。
:::

::: details 登录成功了，但刷新后又回到登录页
先看 `Application -> Local Storage` 里是否真的写入了 `ez-admin-access-token`。

如果没有，通常是 `setAuthSession` 没有执行，或者浏览器当前处于无痕模式 / 存储受限环境。
:::

下一节开始搭建后台整体骨架：[后台布局](./admin-layout)。
