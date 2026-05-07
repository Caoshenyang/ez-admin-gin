import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  OperationLogListQuery,
  OperationLogListResponse,
} from '../types/operation-log'

export async function getOperationLogs(params: OperationLogListQuery) {
  const response = await http.get<ApiResponse<OperationLogListResponse>>('/system/operation-logs', { params })
  return response.data.data
}
