import type { RoleDataScope, RoleStatus } from './role'

export interface RoleFormModel {
  id: number
  code: string
  name: string
  sort: number
  data_scope: RoleDataScope
  custom_department_ids: number[]
  status: RoleStatus
  remark: string
}

export type PermissionTab = 'feature' | 'api'
