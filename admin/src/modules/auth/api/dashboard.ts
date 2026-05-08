import http from '@/api/http'

import type { DashboardData } from '../types/dashboard'
import type { ApiResponse } from '@/api/types'

// getDashboardSummary 获取工作台概览数据（当前用户、健康状态、指标统计、最近操作、登录记录和公告）。
export async function getDashboardSummary() {
  const response = await http.get<ApiResponse<DashboardData>>('/auth/dashboard')
  return response.data.data
}
