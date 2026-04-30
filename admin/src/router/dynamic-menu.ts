import {
  AlbumsOutline,
  AppsOutline,
  BeakerOutline,
  BuildOutline,
  CogOutline,
  DocumentTextOutline,
  FolderOpenOutline,
  GridOutline,
  LayersOutline,
  ListOutline,
  NotificationsOutline,
  PeopleOutline,
  PulseOutline,
  ServerOutline,
  SettingsOutline,
  ShieldCheckmarkOutline,
  TimeOutline,
} from '@vicons/ionicons5'
import { NIcon, type MenuOption } from 'naive-ui'
import type { RouteRecordRaw } from 'vue-router'
import { computed, h, shallowRef, type Component } from 'vue'

import { MenuType, type AuthMenu } from '../types/menu'

type RouteComponent = NonNullable<RouteRecordRaw['component']>
type MenuIconComponent = Component

const placeholderPage = () => import('../pages/system/PlaceholderPage.vue')

const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': () => import('../pages/system/HealthView.vue'),
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': () => import('../pages/system/RoleView.vue'),
  'system/MenuView': () => import('../pages/system/MenuView.vue'),
  'system/ConfigView': () => import('../pages/system/ConfigView.vue'),
  'system/FileView': () => import('../pages/system/FileView.vue'),
  'system/OperationLogView': () => import('../pages/system/OperationLogView.vue'),
  'system/LoginLogView': () => import('../pages/system/LoginLogView.vue'),
  'system/NoticeView': () => import('../pages/system/NoticeView.vue'),
}

const defaultMenuIcon = AppsOutline

// 后端 icon 字段只允许命中这份前端白名单，避免把任意字符串直接当组件渲染。
const menuIconMap: Record<string, MenuIconComponent> = {
  albums: AlbumsOutline,
  app: AppsOutline,
  apps: AppsOutline,
  beaker: BeakerOutline,
  blog: DocumentTextOutline,
  build: BuildOutline,
  cog: CogOutline,
  config: BuildOutline,
  dashboard: GridOutline,
  directory: AlbumsOutline,
  document: DocumentTextOutline,
  edit: DocumentTextOutline,
  experiment: BeakerOutline,
  file: FolderOpenOutline,
  files: FolderOpenOutline,
  folder: FolderOpenOutline,
  grid: GridOutline,
  health: PulseOutline,
  history: TimeOutline,
  home: GridOutline,
  layout: GridOutline,
  layoutdashboard: GridOutline,
  layers: LayersOutline,
  list: ListOutline,
  log: ListOutline,
  loginlog: TimeOutline,
  loginlogs: TimeOutline,
  logs: ListOutline,
  menu: LayersOutline,
  menus: LayersOutline,
  monitor: PulseOutline,
  notice: NotificationsOutline,
  notices: NotificationsOutline,
  notification: NotificationsOutline,
  notifications: NotificationsOutline,
  operationlog: ListOutline,
  operationlogs: ListOutline,
  page: DocumentTextOutline,
  people: PeopleOutline,
  person: PeopleOutline,
  role: ShieldCheckmarkOutline,
  roles: ShieldCheckmarkOutline,
  server: ServerOutline,
  setting: SettingsOutline,
  settings: SettingsOutline,
  shield: ShieldCheckmarkOutline,
  system: SettingsOutline,
  time: TimeOutline,
  user: PeopleOutline,
  users: PeopleOutline,
}

const builtinMenuOptions: MenuOption[] = [
  {
    label: '工作台',
    key: '/dashboard',
    icon: renderMenuIcon(GridOutline),
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
    icon: resolveMenuIcon(menu.icon),
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

function resolveMenuIcon(icon: string) {
  return renderMenuIcon(menuIconMap[normalizeMenuIcon(icon)] ?? defaultMenuIcon)
}

function renderMenuIcon(icon: MenuIconComponent) {
  return () =>
    h(NIcon, null, {
      default: () => h(icon),
    })
}

function normalizeMenuIcon(icon: string) {
  return icon.trim().toLowerCase().replace(/[^a-z0-9]/g, '')
}

function toChildRoutePath(path: string) {
  return path.replace(/^\/+/, '')
}
