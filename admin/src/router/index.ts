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
      component: () => import('../layouts/AdminLayout.vue'),
      beforeEnter: () => {
        if (!hasAccessToken()) {
          return '/login'
        }

        return true
      },
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('../pages/dashboard/DashboardHome.vue'),
          meta: { title: '工作台' },
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '用户管理',
            description: '这一页下一节会开始接入真实用户列表和操作表单。',
          },
          meta: { title: '用户管理' },
        },
        {
          path: 'roles',
          name: 'roles',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '角色权限',
            description: '当前先验证后台布局和标签栏，角色页面后续章节继续补齐。',
          },
          meta: { title: '角色权限' },
        },
        {
          path: 'menus',
          name: 'menus',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '菜单管理',
            description: '这一页下一节会开始与动态菜单能力衔接。',
          },
          meta: { title: '菜单管理' },
        },
        {
          path: 'logs',
          name: 'logs',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '操作日志',
            description: '当前先保留路由出口，后续章节再接真实日志页面。',
          },
          meta: { title: '操作日志' },
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('../pages/system/PlaceholderPage.vue'),
          props: {
            title: '系统设置',
            description: '当前先验证后台布局结构，配置页后续章节继续补齐。',
          },
          meta: { title: '系统设置' },
        },
      ],
    },
  ],
})

export default router
