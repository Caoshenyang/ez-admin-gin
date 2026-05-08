import type { FormRules, SelectOption, TreeOption } from 'naive-ui'

import { MenuStatus, MenuType, type AdminMenu } from '@/modules/iam/types/menu'
import { RoleStatus, type RoleItem, type RoleListQuery, type RolePermissionItem } from '@/modules/iam/types/role'

import type { PermissionRow, RoleFormModel } from '../types/role-page'

export const superAdminRoleCode = 'super_admin'

export const roleStatusOptions: SelectOption[] = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: RoleStatus.Enabled },
  { label: '禁用', value: RoleStatus.Disabled },
]

export const permissionMethodOptions: SelectOption[] = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
]

export const roleFormRules: FormRules = {
  code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
}

export function defaultRoleFormModel(): RoleFormModel {
  return {
    id: 0,
    code: '',
    name: '',
    sort: 10,
    status: RoleStatus.Enabled,
    remark: '',
  }
}

export function defaultRoleListQuery(): RoleListQuery {
  // 角色页的筛选项比较稳定，重置时直接回到这份默认查询即可。
  return {
    page: 1,
    page_size: 100,
    keyword: '',
    status: 0,
  }
}

export function getRoleStatusTagType(status: RoleStatus) {
  return status === RoleStatus.Enabled ? 'success' : 'error'
}

export function toRoleFormModel(role: RoleItem): RoleFormModel {
  return {
    id: role.id,
    code: role.code,
    name: role.name,
    sort: role.sort,
    status: role.status,
    remark: role.remark,
  }
}

export function toRoleMenuTreeOption(menu: AdminMenu): TreeOption {
  const typeText = menu.type === MenuType.Directory ? '目录' : menu.type === MenuType.Menu ? '菜单' : '按钮'
  const statusText = menu.status === MenuStatus.Enabled ? '' : '（禁用）'

  return {
    key: menu.id,
    label: `${menu.title}  ${typeText}  ${menu.code}${statusText}`,
    children: menu.children?.map(toRoleMenuTreeOption),
    disabled: menu.status !== MenuStatus.Enabled,
  }
}

export function flattenRoleMenus(items: AdminMenu[]) {
  const result: AdminMenu[] = []

  for (const item of items) {
    result.push(item)
    result.push(...flattenRoleMenus(item.children ?? []))
  }

  return result
}

export function toPermissionRows(role: RoleItem): PermissionRow[] {
  // 权限面板需要本地行 ID，切换角色时在这里一次性补齐。
  return role.permissions.map((permission, index) => ({
    id: index + 1,
    path: permission.path,
    method: permission.method,
  }))
}

export function defaultPermissionRow(): PermissionRow {
  return {
    id: Date.now(),
    path: '',
    method: 'GET',
  }
}

export function normalizeRolePermissions(rows: PermissionRow[]): RolePermissionItem[] {
  const seen = new Set<string>()
  const result: RolePermissionItem[] = []

  for (const row of rows) {
    const path = row.path.trim()
    const method = row.method.trim().toUpperCase()

    if (!path || !method) continue

    const key = `${method} ${path}`
    if (seen.has(key)) continue

    // API 权限按 “METHOD + PATH” 去重，避免同一行被重复保存到后端。
    seen.add(key)
    result.push({ path, method })
  }

  return result
}

export function buildRoleCreatePayload(formModel: RoleFormModel) {
  return {
    code: formModel.code.trim(),
    name: formModel.name.trim(),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}

export function buildRoleUpdatePayload(formModel: RoleFormModel) {
  return {
    name: formModel.name.trim(),
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}
