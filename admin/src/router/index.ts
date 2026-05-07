import { createRouter, createWebHistory } from 'vue-router'

import { getCurrentUserMenus } from '@/modules/iam/api/menu'
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

  if (to.name === 'admin-dynamic-fallback') {
    return '/dashboard'
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
