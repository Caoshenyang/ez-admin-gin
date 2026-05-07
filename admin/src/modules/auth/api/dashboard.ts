import http from '@/api/http'

import type { DashboardData } from '../types/dashboard'
import type { ApiResponse } from '@/api/types'

export async function getDashboardSummary() {
  const response = await http.get<ApiResponse<DashboardData>>('/auth/dashboard')
  return response.data.data
}
