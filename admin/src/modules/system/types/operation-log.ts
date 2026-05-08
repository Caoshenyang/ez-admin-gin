// OperationLogItem 类型定义。
export interface OperationLogItem {
  id: number
  user_id: number
  username: string
  method: string
  path: string
  route_path: string
  query: string
  ip: string
  user_agent: string
  status_code: number
  latency_ms: number
  success: boolean
  error_message: string
  created_at: string
}

// OperationLogListQuery 类型定义。
export interface OperationLogListQuery {
  page: number
  page_size: number
  username?: string
  method?: string
  path?: string
  success?: string
}

// OperationLogListResponse 类型定义。
export interface OperationLogListResponse {
  items: OperationLogItem[]
  total: number
  page: number
  page_size: number
}
