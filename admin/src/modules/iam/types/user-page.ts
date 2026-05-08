import type { UserStatus } from './user'

export interface UserFormModel {
  id: number
  username: string
  password: string
  nickname: string
  department_id: number
  status: UserStatus
  role_ids: number[]
  post_ids: number[]
}
