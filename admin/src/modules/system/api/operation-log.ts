import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  OperationLogListQuery,
  OperationLogListResponse,
} from '../types/operation-log'

// 获取操作日志列表（分页查询）
export async function getOperationLogs(params: OperationLogListQuery) {
  const response = await http.get<ApiResponse<OperationLogListResponse>>('/system/operation-logs', { params })
  return response.data.data
}
