// 工作台当前登录用户信息
export interface DashboardCurrentUser {
  user_id: number
  username: string
  nickname: string
}

// 系统健康状态信息
export interface DashboardHealth {
  env: string
  database: string
  redis: string
}

// 工作台统计指标数据
export interface DashboardMetrics {
  user_total: number
  enabled_user_total: number
  enabled_role_total: number
  config_total: number
  notice_total: number
  file_total: number
  today_operation_total: number
  today_risk_operation_total: number
  today_login_failed_total: number
}

// 工作台最近操作记录项
export interface DashboardOperationItem {
  id: number
  username: string
  method: string
  path: string
  status_code: number
  success: boolean
  latency_ms: number
  created_at: string
}

// 登录状态常量：成功、失败
export const DashboardLoginStatus = {
  Success: 1,
  Failed: 2,
} as const

// 登录状态联合类型
export type DashboardLoginStatus =
  (typeof DashboardLoginStatus)[keyof typeof DashboardLoginStatus]

// 工作台最近登录记录项
export interface DashboardLoginItem {
  id: number
  username: string
  status: DashboardLoginStatus
  message: string
  ip: string
  created_at: string
}

// 工作台最新公告项
export interface DashboardNoticeItem {
  id: number
  title: string
  status: number
  updated_at: string
}

// 工作台完整数据结构
export interface DashboardData {
  current_user: DashboardCurrentUser
  health: DashboardHealth
  metrics: DashboardMetrics
  recent_operations: DashboardOperationItem[]
  recent_logins: DashboardLoginItem[]
  latest_notices: DashboardNoticeItem[]
}
