---
title: Vue 3 管理台初始化
description: "整理 Vue 3 管理台基础工程，接入 Naive UI、Tailwind CSS 4 和最小可运行页面结构。"
---

# Vue 3 管理台初始化

第一章已经创建了 `admin/` 子项目。本节不再重新执行脚手架初始化，而是在现有 Vue 工程上完成基础整理：安装管理台常用依赖、清理默认示例、建立后续页面目录，并确认开发服务和构建都能通过。

::: tip 🎯 本节目标
完成后，`admin` 会从 Vue 默认示例页变成一个干净的管理台起步工程：页面入口可访问，Naive UI 可用，Tailwind CSS 4 可用，后续可以直接开始实现登录页。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
├─ package.json
├─ vite.config.ts
└─ src/
   ├─ App.vue
   ├─ main.ts
   ├─ router/
   │  └─ index.ts
   ├─ pages/
   │  ├─ auth/
   │  │  └─ LoginPage.vue
   │  └─ dashboard/
   │     └─ DashboardHome.vue
   └─ styles/
      └─ main.css
```

同时会删除 Vue 脚手架生成的示例组件、示例页面和默认样式文件。

::: info 本节只做前端基础整理
这一节不对接登录接口，也不保存 Token。真实登录流程从下一节 [登录页](./login-page) 开始。
:::

## 🛠️ 安装管理台基础依赖

进入 `admin/` 目录：

```bash
# 在项目根目录执行
cd admin
```

安装管理台会用到的基础依赖：

```bash
# 在 admin/ 目录执行
pnpm add naive-ui @vicons/ionicons5 axios
pnpm add -D tailwindcss @tailwindcss/vite
```

| 依赖 | 用途 | 资料 |
| --- | --- | --- |
| `naive-ui` | 管理台 UI 组件库 | [官方文档](https://www.naiveui.com/) |
| `@vicons/ionicons5` | 图标组件 | [项目仓库](https://github.com/07akioni/xicons) |
| `axios` | 后续封装接口请求 | [官方文档](https://axios-http.com/) |
| `tailwindcss` | 页面骨架、间距、背景和响应式布局 | [官方文档](https://tailwindcss.com/docs/installation/using-vite) |
| `@tailwindcss/vite` | Tailwind CSS 4 的 Vite 插件 | [官方文档](https://tailwindcss.com/docs/installation/using-vite) |

::: details 为什么这里用 Naive UI + Tailwind CSS 4
Naive UI 对 Vue 3 和 TypeScript 友好，表单、表格、弹窗、菜单、布局组件都比较完整，适合快速搭建后台管理台；Tailwind CSS 4 更适合负责页面骨架、栅格、间距、背景和响应式布局。

这一章统一采用混合方案：

- 组件层继续以 Naive UI 为主。
- 页面结构和视觉组织优先交给 Tailwind 4。

这样后面写登录页、后台布局和空状态页时，会比“全组件库”更灵活，也比“全 utility class”更稳。
:::

::: info Tailwind CSS 4 这里走 CSS-first 方案
这一版优先使用 Tailwind CSS 4 官方推荐的 Vite 插件和 `@import "tailwindcss"` 方式接入，不额外创建 `tailwind.config.js`。
:::

## 🛠️ 清理脚手架示例文件

Vue 脚手架默认会生成欢迎页、示例组件和示例 store。本项目后续会重新组织管理台页面，所以先删除这些示例文件。

::: code-group

```powershell [Windows PowerShell]
# 在 admin/ 目录执行，删除脚手架示例文件
Remove-Item .\src\components -Recurse -Force
Remove-Item .\src\views -Recurse -Force
Remove-Item .\src\stores\counter.ts -Force
Remove-Item .\src\assets\base.css -Force
Remove-Item .\src\assets\main.css -Force
Remove-Item .\src\assets\logo.svg -Force
```

```bash [macOS / Linux]
# 在 admin/ 目录执行，删除脚手架示例文件
rm -rf src/components src/views
rm -f src/stores/counter.ts
rm -f src/assets/base.css src/assets/main.css src/assets/logo.svg
```

:::

::: warning ⚠️ 只删除脚手架示例文件
不要删除 `src/router/`、`src/main.ts`、`src/App.vue`。这些文件还会继续作为管理台入口使用。
:::

## 🛠️ 创建管理台目录

创建后续会用到的目录：

::: code-group

```powershell [Windows PowerShell]
# 在 admin/ 目录执行，创建管理台常用目录
New-Item -ItemType Directory -Force .\src\api | Out-Null
New-Item -ItemType Directory -Force .\src\layouts | Out-Null
New-Item -ItemType Directory -Force .\src\pages\auth | Out-Null
New-Item -ItemType Directory -Force .\src\pages\dashboard | Out-Null
New-Item -ItemType Directory -Force .\src\styles | Out-Null
New-Item -ItemType Directory -Force .\src\types | Out-Null
New-Item -ItemType Directory -Force .\src\utils | Out-Null
```

```bash [macOS / Linux]
# 在 admin/ 目录执行，创建管理台常用目录
mkdir -p src/api src/layouts src/pages/auth src/pages/dashboard src/styles src/types src/utils
```

:::

目录职责先简单约定如下：

| 目录 | 用途 |
| --- | --- |
| `src/api` | 接口请求封装 |
| `src/layouts` | 后台布局组件 |
| `src/pages` | 页面组件 |
| `src/router` | 路由配置 |
| `src/stores` | Pinia 状态 |
| `src/styles` | 全局样式 |
| `src/types` | 共享类型 |
| `src/utils` | 工具函数 |

## 🛠️ 接入 Tailwind CSS 4

修改 `admin/vite.config.ts`。本次新增 `tailwindcss()` 插件，让 Vite 可以直接处理 Tailwind CSS 4。

```ts
import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite' // [!code ++]
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    tailwindcss(), // [!code ++]
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
})
```

## 🛠️ 更新入口文件

修改 `src/main.ts`。这一处把样式入口改成 `styles/main.css`，并继续保留 Pinia 和 Vue Router。

```ts
import './styles/main.css' // [!code ++]

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
```

## 🛠️ 创建全局样式

创建 `src/styles/main.css`。这里除了保留全局基础样式，还顺手放一组章节内可复用的主题变量。

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
```

::: details 为什么这里没有 `tailwind.config.js`
Tailwind CSS 4 默认更偏向 CSS-first 配置方式，很多主题变量和样式基线都可以直接放在样式入口里维护。

这对当前教程更友好，因为读者只需要先理解两件事：

- `vite.config.ts` 里接入 `@tailwindcss/vite`
- `src/styles/main.css` 里 `@import "tailwindcss"`

后面如果项目真的需要更复杂的主题、插件或扫描策略，再单独扩展也不晚。
:::

## 🛠️ 更新应用入口组件

修改 `src/App.vue`。这一处接入 Naive UI 的全局容器，并保留路由出口。

```vue
<script setup lang="ts">
import {
  NConfigProvider,
  NDialogProvider,
  NLoadingBarProvider,
  NMessageProvider,
  NNotificationProvider,
} from 'naive-ui'
import { RouterView } from 'vue-router'
</script>

<template>
  <n-config-provider>
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <RouterView />
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>
```

::: details 为什么这里用 Naive UI 的 Provider
消息提示、对话框、通知、加载条等组件通常需要全局上下文。先把 Provider 放在 `App.vue`，后续登录页和管理页就可以直接使用这些能力。
:::

## 🛠️ 创建临时页面

创建 `src/pages/auth/LoginPage.vue`。下一节会把它改成真正的登录页，这里先放一个占位页面用于验证路由，并顺手示范“Tailwind 管骨架，Naive UI 管组件”的写法。

```vue
<script setup lang="ts">
import { NButton, NCard, NSpace } from 'naive-ui'
</script>

<template>
  <main class="grid min-h-screen place-items-center bg-slate-100 px-6 py-10">
    <NCard title="EZ Admin" class="w-full max-w-sm rounded-xl shadow-sm" :bordered="false">
      <NSpace vertical :size="16">
        <p class="text-sm text-[#6B7280]">登录页会在下一节接入真实接口。</p>
        <NButton type="primary" block>进入登录页开发</NButton>
      </NSpace>
    </NCard>
  </main>
</template>
```

创建 `src/pages/dashboard/DashboardHome.vue`。这个页面用于验证登录后的默认页面路由，后续会放进后台布局里。

```vue
<script setup lang="ts">
import { NCard } from 'naive-ui'
</script>

<template>
  <main class="bg-slate-100 px-6 py-8">
    <div class="mx-auto max-w-5xl">
      <NCard title="工作台" class="rounded-xl" :bordered="false">
        <p class="text-sm text-[#6B7280]">后台布局会在后续小节补齐。</p>
      </NCard>
    </div>
  </main>
</template>
```

::: tip 先从页面骨架开始使用 Tailwind
这一节先用 Tailwind 处理页面壳子、宽度、留白和背景，让你先熟悉章节里的样式分工。

像 `NCard`、`NButton` 这些稳定组件，仍然继续交给 Naive UI。
:::

## 🛠️ 更新路由

修改 `src/router/index.ts`。先保留最小路由：访问根路径跳转到 `/login`，登录页和工作台各有一个页面。

```ts
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../pages/auth/LoginPage.vue'),
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

::: tip 当前路由只是最小骨架
本节先保证前端项目干净可运行。登录守卫、动态菜单和后台布局会在后面几节逐步补上。
:::

## ✅ 启动验证

在 `admin/` 目录执行：

```bash
# 启动前端开发服务
pnpm dev
```

打开终端输出的地址，通常是：

```text
http://localhost:5173/
```

应该自动进入 `/login`，页面上能看到：

- 浅灰色页面背景
- 中间的 `EZ Admin` 卡片
- Naive UI 按钮和文字组件正常渲染

再访问：

```text
http://localhost:5173/dashboard
```

应该能看到 `工作台` 页面，并且卡片外层留白和宽度已经由 Tailwind utility class 控制。

::: warning ⚠️ 端口以终端输出为准
如果 `5173` 被占用，Vite 会自动换到下一个可用端口。浏览器访问终端实际输出的地址即可。
:::

## ✅ 构建和检查

继续在 `admin/` 目录执行：

```bash
# 类型检查 + 生产构建
pnpm build
```

再执行：

```bash
# 代码检查
pnpm lint
```

两条命令都通过后，说明管理台基础工程整理完成。

下一节开始实现真正的登录表单和接口请求：[登录页](./login-page)。
