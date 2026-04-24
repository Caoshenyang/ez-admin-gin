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
