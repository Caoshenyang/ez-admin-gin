export const LoginLogStatus = {
  Success: 1,
  Failed: 2,
} as const

export type LoginLogStatus = (typeof LoginLogStatus)[keyof typeof LoginLogStatus]

export interface LoginLogItem {
  id: number
  user_id: number
  username: string
  status: LoginLogStatus
  message: string
  ip: string
  user_agent: string
  created_at: string
}

export interface LoginLogListQuery {
  page: number
  page_size: number
  username?: string
  ip?: string
  status?: LoginLogStatus | 0
}

export interface LoginLogListResponse {
  items: LoginLogItem[]
  total: number
  page: number
  page_size: number
}
