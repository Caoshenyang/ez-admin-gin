export const DepartmentStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type DepartmentStatus = (typeof DepartmentStatus)[keyof typeof DepartmentStatus]

export interface DepartmentItem {
  id: number
  parent_id: number
  ancestors: string
  name: string
  code: string
  leader_user_id: number
  sort: number
  status: DepartmentStatus
  remark: string
  created_at: string
  updated_at: string
  children?: DepartmentItem[]
}

export interface DepartmentListQuery {
  keyword?: string
  status?: DepartmentStatus | 0
}

export interface CreateDepartmentPayload {
  parent_id: number
  name: string
  code: string
  leader_user_id: number
  sort: number
  status: DepartmentStatus
  remark: string
}

export type UpdateDepartmentPayload = CreateDepartmentPayload

export interface UpdateDepartmentStatusPayload {
  status: DepartmentStatus
}
