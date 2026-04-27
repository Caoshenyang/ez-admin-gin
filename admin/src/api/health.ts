import http from './http'

import type { ApiResponse } from '../types/http'
import type { SystemHealthData } from '../types/health'

export async function getSystemHealth() {
  const response = await http.get<ApiResponse<SystemHealthData>>('/system/health')
  return response.data.data
}
