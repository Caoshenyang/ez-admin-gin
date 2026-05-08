// 用户状态常量：启用、禁用
export const UserStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

// 用户状态联合类型
export type UserStatus = (typeof UserStatus)[keyof typeof UserStatus]

// 用户列表项数据结构
export interface UserItem {
  id: number
  username: string
  nickname: string
  department_id: number
  status: UserStatus
  role_ids: number[]
  post_ids: number[]
  created_at: string
  updated_at: string
}

// 用户列表查询参数
export interface UserListQuery {
  page: number
  page_size: number
  keyword?: string
  role_id?: number
  status?: UserStatus | 0
}

// 用户列表响应数据
export interface UserListResponse {
  items: UserItem[]
  total: number
  page: number
  page_size: number
}

// 创建用户的请求载荷
export interface CreateUserPayload {
  username: string
  password: string
  nickname: string
  department_id: number
  status: UserStatus
  role_ids: number[]
  post_ids: number[]
}

// 更新用户的请求载荷（不含username、password、role_ids字段）
export interface UpdateUserPayload {
  nickname: string
  department_id: number
  status: UserStatus
  post_ids: number[]
}

// 更新用户状态的请求载荷
export interface UpdateUserStatusPayload {
  status: UserStatus
}

// 更新用户角色分配的请求载荷
export interface UpdateUserRolesPayload {
  role_ids: number[]
}
