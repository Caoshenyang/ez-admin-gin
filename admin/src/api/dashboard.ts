import http from './http'

import type { DashboardData } from '../types/dashboard'
import type { ApiResponse } from '../types/http'

export async function getDashboardSummary() {
  const response = await http.get<ApiResponse<DashboardData>>('/auth/dashboard')
  return response.data.data
}
