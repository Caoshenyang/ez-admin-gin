// 角色状态常量：启用、禁用
export const RoleStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

// 角色状态联合类型
export type RoleStatus = (typeof RoleStatus)[keyof typeof RoleStatus]

// 角色数据权限范围常量
export const RoleDataScope = {
  All: 'all',
  Dept: 'dept',
  DeptAndChildren: 'dept_and_children',
  Self: 'self',
  CustomDept: 'custom_dept',
} as const

// 角色数据权限范围联合类型
export type RoleDataScope = (typeof RoleDataScope)[keyof typeof RoleDataScope]

// 角色列表项数据结构
export interface RoleItem {
  id: number
  code: string
  name: string
  sort: number
  data_scope: RoleDataScope
  custom_department_ids: number[]
  status: RoleStatus
  remark: string
  menu_ids: number[]
  permissions: Array<{
    path: string
    method: string
  }>
  created_at: string
  updated_at: string
}

// 角色列表查询参数
export interface RoleListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: RoleStatus | 0
}

// 角色列表响应数据
export interface RoleListResponse {
  items: RoleItem[]
  total: number
  page: number
  page_size: number
}

// 创建角色的请求载荷
export interface CreateRolePayload {
  code: string
  name: string
  sort: number
  data_scope: RoleDataScope
  custom_department_ids: number[]
  status: RoleStatus
  remark: string
}

// 更新角色的请求载荷（不含code字段）
export interface UpdateRolePayload {
  name: string
  sort: number
  data_scope: RoleDataScope
  custom_department_ids: number[]
  status: RoleStatus
  remark: string
}

// 更新角色状态的请求载荷
export interface UpdateRoleStatusPayload {
  status: RoleStatus
}

// API权限项数据结构
export interface RolePermissionItem {
  path: string
  method: string
}

// 更新角色API权限的请求载荷
export interface UpdateRolePermissionsPayload {
  permissions: RolePermissionItem[]
}

// 更新角色菜单权限的请求载荷
export interface UpdateRoleMenusPayload {
  menu_ids: number[]
}
