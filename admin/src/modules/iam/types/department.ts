// 部门状态常量：启用、禁用
export const DepartmentStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

// 部门状态联合类型
export type DepartmentStatus = (typeof DepartmentStatus)[keyof typeof DepartmentStatus]

// 部门树节点数据结构
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

// 部门列表查询参数
export interface DepartmentListQuery {
  keyword?: string
  status?: DepartmentStatus | 0
}

// 创建部门的请求载荷
export interface CreateDepartmentPayload {
  parent_id: number
  name: string
  code: string
  leader_user_id: number
  sort: number
  status: DepartmentStatus
  remark: string
}

// 更新部门的请求载荷（与创建载荷相同）
export type UpdateDepartmentPayload = CreateDepartmentPayload

// 更新部门状态的请求载荷
export interface UpdateDepartmentStatusPayload {
  status: DepartmentStatus
}
