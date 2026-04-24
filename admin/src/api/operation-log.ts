import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  OperationLogListQuery,
  OperationLogListResponse,
} from '../types/operation-log'

export async function getOperationLogs(params: OperationLogListQuery) {
  const response = await http.get<ApiResponse<OperationLogListResponse>>('/system/operation-logs', { params })
  return response.data.data
}
