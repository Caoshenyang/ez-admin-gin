import http from '@/api/http'

import type {
  AccountProfileResponse,
  LoginRequest,
  LoginResponse,
  UpdateAccountPasswordRequest,
  UpdateAccountProfileRequest,
} from '../types/auth'
import type { ApiResponse } from '@/api/types'

// login 调用后端登录接口。
export async function login(payload: LoginRequest) {
  const response = await http.post<ApiResponse<LoginResponse>>('/auth/login', payload)
  return response.data.data
}

// getAccountProfile 返回当前登录人的账户中心资料。
export async function getAccountProfile() {
  const response = await http.get<ApiResponse<AccountProfileResponse>>('/auth/account')
  return response.data.data
}

// updateAccountProfile 修改当前登录人的基础资料。
export async function updateAccountProfile(payload: UpdateAccountProfileRequest) {
  const response = await http.post<ApiResponse<AccountProfileResponse>>('/auth/account/profile', payload)
  return response.data.data
}

// updateAccountPassword 修改当前登录人的登录密码。
export async function updateAccountPassword(payload: UpdateAccountPasswordRequest) {
  const response = await http.post<ApiResponse<{ updated: boolean }>>('/auth/account/password', payload)
  return response.data.data
}

// logout 调用后端退出登录接口，撤销 refresh token 并黑名单 access token。
export async function logout() {
  await http.post('/auth/logout')
}
