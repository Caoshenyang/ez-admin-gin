import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type { LoginLogListQuery, LoginLogListResponse } from '../types/login-log'

// 获取登录日志列表（分页查询）
export async function getLoginLogs(params: LoginLogListQuery) {
  const response = await http.get<ApiResponse<LoginLogListResponse>>('/system/login-logs', {
    params,
  })
  return response.data.data
}
