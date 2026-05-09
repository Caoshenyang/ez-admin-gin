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
import { resolveRouteComponent } from './route-components'

// MenuIconComponent 菜单图标对应的 Vue 组件类型。
type MenuIconComponent = Component

// AdminMenuOption 侧栏菜单项结构，用于渲染 Naive UI 菜单。
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

// defaultMenuIcon 菜单图标的默认兜底图标。
const defaultMenuIcon = AppsOutline

// builtinMenuCodeIconMap 内置菜单编码到图标的映射，优先级最高。
const builtinMenuCodeIconMap: Record<string, MenuIconComponent> = {
  dashboard: AppsOutline,
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
  dashboard: AppsOutline,
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
  home: AppsOutline,
  layout: AppsOutline,
  layoutdashboard: AppsOutline,
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

// builtinMenuOptions 不依赖后端菜单的内置菜单项（如工作台）。
const builtinMenuOptions: AdminMenuOption[] = [
  {
    label: '工作台',
    key: 'dashboard',
    menuCode: 'dashboard',
    menuType: MenuType.Menu,
    routePath: '/dashboard',
    icon: renderMenuIcon(AppsOutline),
  },
]

// authMenus 存储后端返回的当前用户授权菜单树。
export const authMenus = shallowRef<AuthMenu[]>([])

// sideMenuOptions 计算属性：合并内置菜单与后端授权菜单，生成侧栏菜单选项列表。
export const sideMenuOptions = computed<AdminMenuOption[]>(() => {
  return [...builtinMenuOptions, ...buildMenuOptions(authMenus.value)]
})

// buttonPermissionCodes 计算属性：从授权菜单树中提取所有按钮权限编码，用于权限指令判断。
export const buttonPermissionCodes = computed(() => {
  return collectButtonCodes(authMenus.value)
})

// setAuthMenus 保存后端返回的授权菜单树到响应式变量。
export function setAuthMenus(menus: AuthMenu[]) {
  authMenus.value = menus
}

// clearAuthMenus 清空授权菜单，通常在登出时调用。
export function clearAuthMenus() {
  authMenus.value = []
}

// buildDynamicRoutes 将后端菜单树扁平化为 vue-router 路由配置，仅取 type=menu 的节点。
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

// findMenuTitleByPath 根据路由路径查找对应的菜单标题。
export function findMenuTitleByPath(path: string) {
  return collectPageMenus(authMenus.value).find((menu) => menu.path === path)?.title
}

// findMenuCodeByPath 根据路由路径查找对应的菜单编码，找不到则返回空字符串。
export function findMenuCodeByPath(path: string) {
  return collectPageMenus(authMenus.value).find((menu) => menu.path === path)?.code ?? ''
}

// findMenuOptionByKey 在菜单选项树中递归查找指定 key 的菜单项。
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

// collectExpandedMenuKeysByPath 返回路径对应的菜单链（不含自身），用于自动展开侧栏。
export function collectExpandedMenuKeysByPath(path: string) {
  const chain = findMenuCodeChainByPath(authMenus.value, path)
  return chain.slice(0, -1)
}

// buildMenuOptions 将后端菜单树转换为前端菜单选项数组，过滤掉按钮类型的菜单。
function buildMenuOptions(menus: AuthMenu[]): AdminMenuOption[] {
  return menus.map(toMenuOption).filter(isMenuOption)
}

// toMenuOption 将单个后端菜单节点转换为前端菜单选项，按钮类型返回 null。
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
    // 空子目录在前端禁用点击，避免导航到无内容的目录节点。
    disabled: menu.type === MenuType.Directory && children.length === 0,
    children: children.length > 0 ? children : undefined,
  }
}

// isMenuOption 类型守卫：判断 toMenuOption 的返回值是否为有效菜单选项。
function isMenuOption(option: AdminMenuOption | null): option is AdminMenuOption {
  return option !== null
}

// collectPageMenus 递归收集菜单树中所有类型为"页面菜单"且有路径的节点。
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

// collectButtonCodes 递归收集菜单树中所有按钮类型的权限编码。
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

// resolveMenuIcon 根据菜单编码和后端 icon 字段解析菜单图标，优先使用内置映射。
function resolveMenuIcon(code: string, icon: string) {
  const normalizedCode = normalizeMenuCode(code)
  const builtinIcon = builtinMenuCodeIconMap[normalizedCode]
  if (builtinIcon) {
    return renderMenuIcon(builtinIcon)
  }

  return renderMenuIcon(menuIconMap[normalizeMenuIcon(icon)] ?? defaultMenuIcon)
}

// renderMenuIcon 将图标组件包装为 Naive UI 菜单所需的渲染函数。
function renderMenuIcon(icon: MenuIconComponent) {
  return () =>
    h(
      NIcon,
      null,
      {
        default: () => h(icon),
      },
    )
}

// normalizeMenuIcon 将后端 icon 字段标准化：去空格、转小写、移除非字母数字字符。
function normalizeMenuIcon(icon: string) {
  return icon.trim().toLowerCase().replace(/[^a-z0-9]/g, '')
}

// normalizeMenuCode 将菜单编码标准化：去空格、转小写。
function normalizeMenuCode(code: string) {
  return code.trim().toLowerCase()
}

// toChildRoutePath 移除路由路径的前导斜杠，使其成为嵌套子路由路径。
function toChildRoutePath(path: string) {
  return path.replace(/^\/+/, '')
}

// findMenuCodeChainByPath 递归查找目标路径在菜单树中的编码链（含自身），用于侧栏自动展开。
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
