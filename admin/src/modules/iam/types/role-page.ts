import type { RoleStatus } from './role'

export interface RoleFormModel {
  id: number
  code: string
  name: string
  sort: number
  status: RoleStatus
  remark: string
}

export interface PermissionRow {
  id: number
  path: string
  method: string
}

export type PermissionTab = 'base' | 'menu' | 'button' | 'api' | 'data' | 'users'
