import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type { SystemHealthData } from '../types/health'

export async function getSystemHealth() {
  const response = await http.get<ApiResponse<SystemHealthData>>('/system/health')
  return response.data.data
}
