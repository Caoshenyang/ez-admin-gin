import http from './http'

import type { ApiResponse } from '../types/http'
import type { LoginLogListQuery, LoginLogListResponse } from '../types/login-log'

export async function getLoginLogs(params: LoginLogListQuery) {
  const response = await http.get<ApiResponse<LoginLogListResponse>>('/system/login-logs', { params })
  return response.data.data
}
