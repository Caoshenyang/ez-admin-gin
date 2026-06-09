---
title: 登录与认证流程
description: 前端登录流程、Token 管理、路由守卫
---

# 登录与认证流程

## 登录流程

```
1. 用户在 LoginPage.vue 输入用户名和密码
   ↓
2. useLoginPage composable 调用 auth API
   ↓
3. POST /api/v1/auth/login → 后端验证 → 返回 JWT Token
   ↓
4. Token 存储到 localStorage（记住我）或 sessionStorage
   ↓
5. 路由跳转到 /dashboard
```

## Token 管理

`utils/auth.ts` 负责 Token 的存取：

```typescript
// 存储
setAuthSession(token, remember: boolean)
// remember=true → localStorage
// remember=false → sessionStorage

// 读取
getAccessToken(): string | null

// 清除
clearAuthSession()

// 检查
hasAccessToken(): boolean

// 构造请求头
getAuthorizationHeader(): string
```

## 路由守卫

`router/guard.ts` 实现全局前置守卫：

```typescript
router.beforeEach((to) => {
  if (需要认证的路由 && !hasAccessToken()) {
    return { name: 'login' }
  }
  if (to.path === '/login' && hasAccessToken()) {
    return '/dashboard'
  }
})
```

## API 拦截器

`api/http.ts` 中的拦截器：

**请求拦截器：**

```typescript
instance.interceptors.request.use((config) => {
  const token = getAccessToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
```

**响应拦截器：**

```typescript
instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401 || error.response?.status === 403) {
      clearAuthSession()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)
```

## 认证状态流转

```
未登录
  → 访问任意页面 → 守卫拦截 → 跳转登录页
  → 输入账号密码 → API 登录 → 存储 Token → 跳转 Dashboard
  → 加载动态菜单 → 注册路由 → 渲染侧边栏

已登录
  → 访问页面 → 守卫放行 → 正常加载
  → Token 过期 → API 返回 401 → 清除 Token → 跳转登录页
```
