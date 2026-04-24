export const RoleStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type RoleStatus = (typeof RoleStatus)[keyof typeof RoleStatus]

export interface RoleItem {
  id: number
  code: string
  name: string
  sort: number
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

export interface RoleListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: RoleStatus | 0
}

export interface RoleListResponse {
  items: RoleItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateRolePayload {
  code: string
  name: string
  sort: number
  status: RoleStatus
  remark: string
}

export interface UpdateRolePayload {
  name: string
  sort: number
  status: RoleStatus
  remark: string
}

export interface UpdateRoleStatusPayload {
  status: RoleStatus
}

export interface RolePermissionItem {
  path: string
  method: string
}

export interface UpdateRolePermissionsPayload {
  permissions: RolePermissionItem[]
}

export interface UpdateRoleMenusPayload {
  menu_ids: number[]
}
