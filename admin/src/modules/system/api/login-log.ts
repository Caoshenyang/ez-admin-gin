import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type { LoginLogListQuery, LoginLogListResponse } from '../types/login-log'

export async function getLoginLogs(params: LoginLogListQuery) {
  const response = await http.get<ApiResponse<LoginLogListResponse>>('/system/login-logs', { params })
  return response.data.data
}
