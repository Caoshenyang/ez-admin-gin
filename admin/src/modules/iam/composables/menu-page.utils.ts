import type { SelectOption } from 'naive-ui'

import { MenuStatus, MenuType, type AdminMenu, type UpdateMenuPayload } from '@/modules/iam/types/menu'

import type { MenuFormModel, MenuQuery } from '../types/menu-page'

export const menuTypeOptions: SelectOption[] = [
  { label: '类型：全部', value: 0 },
  { label: '目录', value: MenuType.Directory },
  { label: '菜单', value: MenuType.Menu },
  { label: '按钮', value: MenuType.Button },
]

export const menuFormTypeOptions: SelectOption[] = menuTypeOptions.slice(1)

export const menuStatusOptions: SelectOption[] = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: MenuStatus.Enabled },
  { label: '禁用', value: MenuStatus.Disabled },
]

export function defaultMenuQuery(): MenuQuery {
  // 菜单页默认只看启用项，重置筛选和首屏状态共用同一份查询基线。
  return {
    keyword: '',
    type: 0,
    status: MenuStatus.Enabled,
  }
}

// 生成菜单表单的默认值
export function defaultMenuFormModel(): MenuFormModel {
  return {
    id: 0,
    parent_id: 0,
    type: MenuType.Directory,
    code: '',
    title: '',
    path: '',
    component: '',
    icon: '',
    sort: 10,
    status: MenuStatus.Enabled,
    remark: '',
  }
}

export function toMenuFormModel(menu: AdminMenu): MenuFormModel {
  return {
    id: menu.id,
    parent_id: menu.parent_id,
    type: menu.type,
    code: menu.code,
    title: menu.title,
    path: menu.path,
    component: menu.component,
    icon: menu.icon,
    sort: menu.sort,
    status: menu.status,
    remark: menu.remark,
  }
}

export function buildMenuPayload(formModel: MenuFormModel): UpdateMenuPayload {
  const isButton = formModel.type === MenuType.Button

  return {
    parent_id: formModel.parent_id,
    type: formModel.type,
    title: formModel.title.trim(),
    path: isButton ? '' : formModel.path.trim(),
    component: isButton ? '' : formModel.component.trim(),
    icon: formModel.icon.trim(),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}

// 递归扁平化菜单树为一维数组
export function flattenMenus(items: AdminMenu[]): AdminMenu[] {
  const result: AdminMenu[] = []

  for (const item of items) {
    result.push(item)
    result.push(...flattenMenus(item.children ?? []))
  }

  return result
}

// 根据关键词、类型和状态过滤菜单树
export function filterMenus(items: AdminMenu[], query: MenuQuery): AdminMenu[] {
  const keyword = query.keyword.trim().toLowerCase()
  const result: AdminMenu[] = []

  for (const item of items) {
    const children = filterMenus(item.children ?? [], query)
    const matchedKeyword =
      keyword === '' ||
      item.title.toLowerCase().includes(keyword) ||
      item.code.toLowerCase().includes(keyword) ||
      item.path.toLowerCase().includes(keyword)
    const matchedType = query.type === 0 || item.type === query.type
    const matchedStatus = query.status === 0 || item.status === query.status

    if ((matchedKeyword && matchedType && matchedStatus) || children.length > 0) {
      result.push({
        ...item,
        children: children.length > 0 ? children : undefined,
      })
    }
  }

  return result
}

// 构建上级菜单选择选项，排除按钮类型和当前编辑的菜单
export function buildMenuParentOptions(flatMenus: AdminMenu[], currentMenuID: number): SelectOption[] {
  const options: SelectOption[] = [{ label: '根节点', value: 0 }]

  for (const menu of flatMenus) {
    if (menu.type === MenuType.Button || menu.id === currentMenuID) {
      continue
    }

    options.push({
      label: `${'　'.repeat(menuLevel(flatMenus, menu.id))}${menu.title}`,
      value: menu.id,
    })
  }

  return options
}

// 计算菜单在扁平列表中的层级深度
function menuLevel(flatMenus: AdminMenu[], id: number) {
  let level = 0
  let current = flatMenus.find((menu) => menu.id === id)

  while (current && current.parent_id !== 0) {
    level += 1
    current = flatMenus.find((menu) => menu.id === current?.parent_id)
  }

  return level
}
