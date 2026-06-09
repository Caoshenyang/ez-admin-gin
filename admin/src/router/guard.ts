import type { RouteLocationNormalized, Router } from 'vue-router'

import { getCurrentUserMenus } from '@/modules/iam/api/menu'
import { clearAuthSession, hasAccessToken } from '@/utils/auth'

import { buildDynamicRoutes, clearAuthMenus, setAuthMenus } from './dynamic-menu'

// AdminRouteRegistrar 仅需 addRoute 方法的路由注册器类型。
type AdminRouteRegistrar = Pick<Router, 'addRoute'>

// createAuthGuardController 创建路由守卫控制器，负责登录态判断、动态路由注册与重置。
export function createAuthGuardController(router: AdminRouteRegistrar) {
  // removeDynamicRouteCallbacks 存储已注册动态路由的移除回调，登出时批量清理。
  const removeDynamicRouteCallbacks: Array<() => void> = []
  // dynamicRoutesReady 标记动态路由是否已完成注册，避免重复请求后端菜单。
  let dynamicRoutesReady = false

  // resetDynamicRoutes 移除所有已注册的动态路由并重置状态。
  function resetDynamicRoutes() {
    for (const removeRoute of removeDynamicRouteCallbacks) {
      removeRoute()
    }

    removeDynamicRouteCallbacks.length = 0
    dynamicRoutesReady = false
    clearAuthMenus()
  }

  // guard 核心路由守卫函数，根据目标路由和登录态决定导航结果（放行 / 重定向 / 登录）。
  async function guard(to: RouteLocationNormalized) {
    // 已登录用户访问 /login 时自动跳转到工作台。
    if (to.path === '/login') {
      return hasAccessToken() ? '/dashboard' : true
    }

    // 未登录时重定向到登录页，并携带原目标路径以便登录后回跳。
    if (!hasAccessToken()) {
      resetDynamicRoutes()
      return {
        path: '/login',
        query: {
          redirect: to.fullPath,
        },
      }
    }

    // 首次进入系统时从后端拉取菜单并注册动态路由，完成后重新导航以确保路由命中。
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

    // 动态路由加载完成后，若导航命中兜底路由则重定向到工作台。
    if (to.name === 'admin-dynamic-fallback') {
      return '/dashboard'
    }

    return true
  }

  return {
    guard,
    resetDynamicRoutes,
  }
}
