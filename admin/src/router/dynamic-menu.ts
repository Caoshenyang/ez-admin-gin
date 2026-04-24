import type { MenuOption } from 'naive-ui'
import type { RouteRecordRaw } from 'vue-router'
import { computed, shallowRef } from 'vue'

import { MenuType, type AuthMenu } from '../types/menu'

type RouteComponent = NonNullable<RouteRecordRaw['component']>

const placeholderPage = () => import('../pages/system/PlaceholderPage.vue')

const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': placeholderPage,
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': placeholderPage,
  'system/MenuView': placeholderPage,
  'system/ConfigView': placeholderPage,
  'system/FileView': placeholderPage,
  'system/OperationLogView': placeholderPage,
  'system/LoginLogView': placeholderPage,
}

const builtinMenuOptions: MenuOption[] = [
  {
    label: '工作台',
    key: '/dashboard',
  },
]

export const authMenus = shallowRef<AuthMenu[]>([])

export const sideMenuOptions = computed<MenuOption[]>(() => {
  return [...builtinMenuOptions, ...buildMenuOptions(authMenus.value)]
})

export const buttonPermissionCodes = computed(() => {
  return collectButtonCodes(authMenus.value)
})

export function setAuthMenus(menus: AuthMenu[]) {
  authMenus.value = menus
}

export function clearAuthMenus() {
  authMenus.value = []
}

export function buildDynamicRoutes(menus: AuthMenu[]) {
  return collectPageMenus(menus).map<RouteRecordRaw>((menu) => ({
    path: toChildRoutePath(menu.path),
    name: `menu-${menu.id}`,
    component: resolveRouteComponent(menu.component),
    props: {
      title: menu.title,
      description: `${menu.title} 页面后续会接入真实业务。`,
    },
    meta: {
      title: menu.title,
      menuCode: menu.code,
    },
  }))
}

export function findMenuTitleByPath(path: string) {
  return collectPageMenus(authMenus.value).find((menu) => menu.path === path)?.title
}

function buildMenuOptions(menus: AuthMenu[]) {
  return menus.map(toMenuOption).filter(isMenuOption)
}

function toMenuOption(menu: AuthMenu): MenuOption | null {
  if (menu.type === MenuType.Button) {
    return null
  }

  const children = buildMenuOptions(menu.children ?? [])
  const key = menu.path || menu.code

  return {
    label: menu.title,
    key,
    disabled: menu.type === MenuType.Directory && children.length === 0,
    children: children.length > 0 ? children : undefined,
  }
}

function isMenuOption(option: MenuOption | null): option is MenuOption {
  return option !== null
}

function collectPageMenus(menus: AuthMenu[]) {
  const result: AuthMenu[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Menu && menu.path) {
      result.push(menu)
    }

    result.push(...collectPageMenus(menu.children ?? []))
  }

  return result
}

function collectButtonCodes(menus: AuthMenu[]) {
  const result: string[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Button) {
      result.push(menu.code)
    }

    result.push(...collectButtonCodes(menu.children ?? []))
  }

  return result
}

function resolveRouteComponent(component: string) {
  return routeComponentMap[component] ?? placeholderPage
}

function toChildRoutePath(path: string) {
  return path.replace(/^\/+/, '')
}
