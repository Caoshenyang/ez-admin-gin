import http from './http'

import type { ApiResponse } from '../types/http'
import type { RoleListQuery, RoleListResponse } from '../types/role'

export async function getRoles(params: RoleListQuery) {
  const response = await http.get<ApiResponse<RoleListResponse>>('/system/roles', { params })
  return response.data.data
}
