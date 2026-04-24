import http from './http'

import type { LoginRequest, LoginResponse } from '../types/auth'
import type { ApiResponse } from '../types/http'

// login 调用后端登录接口。
export async function login(payload: LoginRequest) {
  const response = await http.post<ApiResponse<LoginResponse>>('/auth/login', payload)
  return response.data.data
}
