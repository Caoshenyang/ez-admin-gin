import type { FormRules, SelectOption, TreeOption } from 'naive-ui'

import { ApiStatus, type AdminAPI } from '@/modules/iam/types/api-resource'
import { MenuStatus, type AdminMenu } from '@/modules/iam/types/menu'
import {
  RoleDataScope,
  RoleStatus,
  type RoleItem,
  type RoleListQuery,
} from '@/modules/iam/types/role'

import type { RoleFormModel } from '../types/role-page'

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

export function flattenRoleAPIs(items: AdminAPI[]) {
  return items
}

export function toAPIPermissionTreeOptions(apis: AdminAPI[]): TreeOption[] {
  const groups = new Map<string, AdminAPI[]>()
  for (const api of apis) {
    const module = api.module || 'other'
    groups.set(module, [...(groups.get(module) ?? []), api])
  }

  return [...groups.entries()]
    .sort(([left], [right]) => moduleOrder(left) - moduleOrder(right) || left.localeCompare(right))
    .map(([module, items]) => ({
      key: `api-module:${module}`,
      label: moduleLabel(module),
      children: items
        .slice()
        .sort((left, right) => left.sort - right.sort || left.id - right.id)
        .map(toAPIPermissionTreeOption),
    }))
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

function toAPIPermissionTreeOption(api: AdminAPI): TreeOption {
  const statusText = api.status === ApiStatus.Enabled ? '' : ' · 禁用'
  return {
    key: api.id,
    label: `${api.name} · ${api.method} ${api.path}${statusText}`,
    disabled: api.status !== ApiStatus.Enabled,
  }
}

function moduleOrder(module: string) {
  const orders = new Map([
    ['iam', 10],
    ['system', 20],
    ['audit', 30],
  ])
  return orders.get(module) ?? 999
}

function moduleLabel(module: string) {
  const labels = new Map([
    ['iam', '权限管理接口'],
    ['system', '系统设置接口'],
    ['audit', '审计监控接口'],
  ])
  return labels.get(module) ?? `${module} 接口`
}

function normalizeCustomDepartmentIDs(formModel: RoleFormModel) {
  if (formModel.data_scope !== RoleDataScope.CustomDept) {
    return []
  }

  return [
    ...new Set(
      formModel.custom_department_ids.map(Number).filter((id) => Number.isFinite(id) && id > 0),
    ),
  ]
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
