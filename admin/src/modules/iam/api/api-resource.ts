import http from '@/api/http'
import type { ApiResponse } from '@/api/types'

import type { AdminAPI } from '../types/api-resource'

// getAdminAPIs 查询所有接口权限元数据。
export async function getAdminAPIs() {
  const response = await http.get<ApiResponse<AdminAPI[]>>('/system/apis')
  return response.data.data ?? []
}
