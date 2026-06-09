import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type { SystemHealthData } from '../types/health'

// 获取系统健康状态信息
export async function getSystemHealth() {
  const response = await http.get<ApiResponse<SystemHealthData>>('/system/health')
  return response.data.data
}
