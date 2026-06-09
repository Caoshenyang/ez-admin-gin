import { createRouter, createWebHistory, type Router, type RouterHistory } from 'vue-router'

import { hasAccessToken } from '../utils/auth'
import { createAuthGuardController } from './guard'

export function createAppRouter(
  history: RouterHistory = createWebHistory(import.meta.env.BASE_URL),
) {
  const router = createRouter({
    history,
    routes: [
      {
        path: '/',
        redirect: () => (hasAccessToken() ? '/dashboard' : '/login'),
      },
      {
        path: '/login',
        name: 'login',
        component: () => import('@/modules/auth/pages/LoginPage.vue'),
      },
      {
        path: '/',
        name: 'admin',
        component: () => import('../layouts/AdminLayout.vue'),
        children: [
          {
            path: 'dashboard',
            name: 'dashboard',
            component: () => import('@/modules/auth/pages/DashboardHome.vue'),
            meta: { title: '工作台' },
          },
          {
            path: 'account/profile',
            name: 'account-profile',
            component: () => import('@/modules/auth/pages/AccountCenterPage.vue'),
            meta: { title: '账户中心' },
          },
          {
            path: ':pathMatch(.*)*',
            name: 'admin-dynamic-fallback',
            component: () => import('@/modules/system/pages/PlaceholderPage.vue'),
            meta: { title: '页面加载中' },
          },
        ],
      },
    ],
  })

  attachAuthGuard(router)
  return router
}

// attachAuthGuard 注册全局前置守卫：未登录跳登录页，登录后首次访问时加载动态路由。
function attachAuthGuard(router: Router) {
  const controller = createAuthGuardController(router)
  router.beforeEach(controller.guard)
  ;(router as Router & { resetDynamicRoutes: () => void }).resetDynamicRoutes =
    controller.resetDynamicRoutes
}

const router = createAppRouter()

// resetDynamicRoutes 用于退出登录或 Token 失效时清理旧账号菜单。
export function resetDynamicRoutes() {
  ;(router as Router & { resetDynamicRoutes: () => void }).resetDynamicRoutes()
}

export default router
