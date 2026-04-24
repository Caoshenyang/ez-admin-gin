export const UserStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type UserStatus = (typeof UserStatus)[keyof typeof UserStatus]

export interface UserItem {
  id: number
  username: string
  nickname: string
  status: UserStatus
  role_ids: number[]
  created_at: string
  updated_at: string
}

export interface UserListQuery {
  page: number
  page_size: number
  keyword?: string
  role_id?: number
  status?: UserStatus
}

export interface UserListResponse {
  items: UserItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateUserPayload {
  username: string
  password: string
  nickname: string
  status: UserStatus
  role_ids: number[]
}

export interface UpdateUserPayload {
  nickname: string
  status: UserStatus
}

export interface UpdateUserStatusPayload {
  status: UserStatus
}

export interface UpdateUserRolesPayload {
  role_ids: number[]
}
