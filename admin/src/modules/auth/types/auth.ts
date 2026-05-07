// LoginRequest 对应登录接口请求体。
export interface LoginRequest {
  username: string
  password: string
}

// LoginResponse 对应登录接口 data 字段。
export interface LoginResponse {
  user_id: number
  username: string
  nickname: string
  access_token: string
  token_type: string
  expires_at: string
}

export interface AccountDataScopeSummary {
  allow_all: boolean
  require_self: boolean
  include_department: boolean
  include_dept_tree: boolean
  custom_department_ids: number[]
}

export interface AccountProfileResponse {
  user_id: number
  username: string
  nickname: string
  department_id: number
  department_name: string
  status: number
  role_codes: string[]
  is_super_admin: boolean
  data_scope: AccountDataScopeSummary
  updated_at: string
}

export interface UpdateAccountProfileRequest {
  nickname: string
}

export interface UpdateAccountPasswordRequest {
  old_password: string
  new_password: string
}
