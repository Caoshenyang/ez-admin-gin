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
