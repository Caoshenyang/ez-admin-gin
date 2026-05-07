import {
  AlbumsOutline,
  AppsOutline,
  BriefcaseOutline,
  BeakerOutline,
  BuildOutline,
  CogOutline,
  DocumentTextOutline,
  FolderOpenOutline,
  GitBranchOutline,
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

import { MenuType, type AuthMenu } from '@/modules/iam/types/menu'

type RouteComponent = NonNullable<RouteRecordRaw['component']>
type MenuIconComponent = Component
export interface AdminMenuOption {
  label: string
  key: string
  menuCode: string
  menuType: MenuType
  routePath?: string
  icon?: MenuOption['icon']
  disabled?: boolean
  children?: AdminMenuOption[]
}

const placeholderPage = () => import('@/modules/system/pages/PlaceholderPage.vue')

const routeComponentMap: Record<string, RouteComponent> = {
  'system/DepartmentView': () => import('@/modules/iam/pages/DepartmentView.vue'),
  'system/HealthView': () => import('@/modules/system/pages/HealthView.vue'),
  'system/AttachmentView': () => import('@/modules/system/pages/AttachmentView.vue'),
  'system/ConfigView': () => import('@/modules/system/pages/ConfigView.vue'),
  'system/DictView': () => import('@/modules/system/pages/DictView.vue'),
  'system/FileView': () => import('@/modules/system/pages/FileView.vue'),
  'system/OperationLogView': () => import('@/modules/system/pages/OperationLogView.vue'),
  'system/LoginLogView': () => import('@/modules/system/pages/LoginLogView.vue'),
  'system/MenuView': () => import('@/modules/iam/pages/MenuView.vue'),
  'system/NoticeView': () => import('@/modules/system/pages/NoticeView.vue'),
  'system/PostView': () => import('@/modules/iam/pages/PostView.vue'),
  'system/RoleView': () => import('@/modules/iam/pages/RoleView.vue'),
  'system/UserView': () => import('@/modules/iam/pages/UserView.vue'),
}

const defaultMenuIcon = AppsOutline

const builtinMenuCodeIconMap: Record<string, MenuIconComponent> = {
  dashboard: GridOutline,
  system: SettingsOutline,
  'system:health': PulseOutline,
  'system:user': PeopleOutline,
  'system:role': ShieldCheckmarkOutline,
  'system:department': GitBranchOutline,
  'system:post': BriefcaseOutline,
  'system:menu': LayersOutline,
  'system:config': CogOutline,
  'system:file': FolderOpenOutline,
  'system:attachment': FolderOpenOutline,
  'system:dict': AlbumsOutline,
  'system:dict-type': AlbumsOutline,
  'system:dict-item': ListOutline,
  'system:operation-log': ListOutline,
  'system:login-log': TimeOutline,
  'system:notice': NotificationsOutline,
}

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
  gitbranch: GitBranchOutline,
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
  briefcase: BriefcaseOutline,
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

const builtinMenuOptions: AdminMenuOption[] = [
  {
    label: '工作台',
    key: 'dashboard',
    menuCode: 'dashboard',
    menuType: MenuType.Menu,
    routePath: '/dashboard',
    icon: renderMenuIcon(GridOutline),
  },
]

export const authMenus = shallowRef<AuthMenu[]>([])

export const sideMenuOptions = computed<AdminMenuOption[]>(() => {
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

export function findMenuCodeByPath(path: string) {
  return collectPageMenus(authMenus.value).find((menu) => menu.path === path)?.code ?? ''
}

export function findMenuOptionByKey(
  key: string,
  options: AdminMenuOption[] = sideMenuOptions.value,
): AdminMenuOption | null {
  for (const option of options) {
    if (option.key === key) {
      return option
    }

    const child = option.children ? findMenuOptionByKey(key, option.children) : null
    if (child) {
      return child
    }
  }

  return null
}

export function collectExpandedMenuKeysByPath(path: string) {
  const chain = findMenuCodeChainByPath(authMenus.value, path)
  return chain.slice(0, -1)
}

function buildMenuOptions(menus: AuthMenu[]): AdminMenuOption[] {
  return menus.map(toMenuOption).filter(isMenuOption)
}

function toMenuOption(menu: AuthMenu): AdminMenuOption | null {
  if (menu.type === MenuType.Button) {
    return null
  }

  const children = buildMenuOptions(menu.children ?? [])

  return {
    label: menu.title,
    key: menu.code,
    menuCode: menu.code,
    menuType: menu.type,
    routePath: menu.type === MenuType.Menu ? menu.path : undefined,
    icon: resolveMenuIcon(menu.code, menu.icon),
    disabled: menu.type === MenuType.Directory && children.length === 0,
    children: children.length > 0 ? children : undefined,
  }
}

function isMenuOption(option: AdminMenuOption | null): option is AdminMenuOption {
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

function resolveMenuIcon(code: string, icon: string) {
  const normalizedCode = normalizeMenuCode(code)
  const builtinIcon = builtinMenuCodeIconMap[normalizedCode]
  if (builtinIcon) {
    return renderMenuIcon(builtinIcon)
  }

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

function normalizeMenuCode(code: string) {
  return code.trim().toLowerCase()
}

function toChildRoutePath(path: string) {
  return path.replace(/^\/+/, '')
}

function findMenuCodeChainByPath(menus: AuthMenu[], path: string, parents: string[] = []): string[] {
  for (const menu of menus) {
    const nextParents = [...parents, menu.code]
    if (menu.type === MenuType.Menu && menu.path === path) {
      return nextParents
    }

    const childChain = findMenuCodeChainByPath(menu.children ?? [], path, nextParents)
    if (childChain.length > 0) {
      return childChain
    }
  }

  return []
}
