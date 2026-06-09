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

// AccountDataScopeSummary 用户数据范围权限摘要，描述可查看的数据边界。
export interface AccountDataScopeSummary {
  allow_all: boolean
  require_self: boolean
  include_department: boolean
  include_dept_tree: boolean
  custom_department_ids: number[]
}

// AccountProfileResponse 账户中心资料接口的返回结构。
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

// UpdateAccountProfileRequest 修改账户基础资料的请求体。
export interface UpdateAccountProfileRequest {
  nickname: string
}

// UpdateAccountPasswordRequest 修改登录密码的请求体。
export interface UpdateAccountPasswordRequest {
  old_password: string
  new_password: string
}
