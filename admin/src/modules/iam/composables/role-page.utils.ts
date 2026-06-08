import type { FormRules, SelectOption, TreeOption } from 'naive-ui'

import { MenuStatus, type AdminMenu } from '@/modules/iam/types/menu'
import {
  RoleDataScope,
  RoleStatus,
  type RoleItem,
  type RoleListQuery,
  type RolePermissionItem,
} from '@/modules/iam/types/role'

import type { PermissionRow, RoleFormModel } from '../types/role-page'

export const superAdminRoleCode = 'super_admin'

export const roleStatusOptions: SelectOption[] = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: RoleStatus.Enabled },
  { label: '禁用', value: RoleStatus.Disabled },
]

export const roleDataScopeOptions: SelectOption[] = [
  { label: '全部数据', value: RoleDataScope.All },
  { label: '本部门数据', value: RoleDataScope.Dept },
  { label: '本部门及下级', value: RoleDataScope.DeptAndChildren },
  { label: '仅本人数据', value: RoleDataScope.Self },
  { label: '自定义部门', value: RoleDataScope.CustomDept },
]

export const roleDataScopeLabels = new Map<RoleDataScope, string>(
  roleDataScopeOptions.map((option) => [option.value as RoleDataScope, String(option.label)]),
)

export const roleDataScopeHelps = new Map<RoleDataScope, string>([
  [RoleDataScope.All, '可查看所有组织与业务数据'],
  [RoleDataScope.Dept, '仅查看当前归属部门数据'],
  [RoleDataScope.DeptAndChildren, '包含本部门和下级部门数据'],
  [RoleDataScope.Self, '仅查看本人创建或归属的数据'],
  [RoleDataScope.CustomDept, '按已选授权部门控制可见数据'],
])

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
    data_scope: RoleDataScope.Self,
    custom_department_ids: [],
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
    data_scope: role.data_scope,
    custom_department_ids: [...role.custom_department_ids],
    status: role.status,
    remark: role.remark,
  }
}

export function toFeaturePermissionTreeOptions(menus: AdminMenu[]): TreeOption[] {
  return menus.map(toFeaturePermissionTreeOption)
}

export function flattenRoleMenus(items: AdminMenu[]) {
  const result: AdminMenu[] = []

  for (const item of items) {
    result.push(item)
    result.push(...flattenRoleMenus(item.children ?? []))
  }

  return result
}

function toFeaturePermissionTreeOption(menu: AdminMenu): TreeOption {
  const children = (menu.children ?? []).map(toFeaturePermissionTreeOption)
  const statusText = menu.status === MenuStatus.Enabled ? '' : ' · 禁用'
  const codeText = menu.code ? ` · ${menu.code}` : ''

  return {
    key: menu.id,
    label: `${menu.title}${codeText}${statusText}`,
    children,
    disabled: menu.status !== MenuStatus.Enabled,
  }
}

export function toPermissionRows(role: RoleItem): PermissionRow[] {
  // 权限面板需要本地行 ID，切换角色时在这里一次性补齐。
  return (role.permissions ?? []).map((permission, index) => ({
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

function normalizeCustomDepartmentIDs(formModel: RoleFormModel) {
  if (formModel.data_scope !== RoleDataScope.CustomDept) {
    return []
  }

  return [...new Set(formModel.custom_department_ids.map(Number).filter((id) => Number.isFinite(id) && id > 0))]
}

export function buildRoleCreatePayload(formModel: RoleFormModel) {
  return {
    code: formModel.code.trim(),
    name: formModel.name.trim(),
    sort: formModel.sort,
    data_scope: formModel.data_scope,
    custom_department_ids: normalizeCustomDepartmentIDs(formModel),
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}

export function buildRoleUpdatePayload(formModel: RoleFormModel) {
  return {
    name: formModel.name.trim(),
    sort: formModel.sort,
    data_scope: formModel.data_scope,
    custom_department_ids: normalizeCustomDepartmentIDs(formModel),
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}
