export interface DashboardCurrentUser {
  user_id: number
  username: string
  nickname: string
}

export interface DashboardHealth {
  env: string
  database: string
  redis: string
}

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

export const DashboardLoginStatus = {
  Success: 1,
  Failed: 2,
} as const

export type DashboardLoginStatus =
  (typeof DashboardLoginStatus)[keyof typeof DashboardLoginStatus]

export interface DashboardLoginItem {
  id: number
  username: string
  status: DashboardLoginStatus
  message: string
  ip: string
  created_at: string
}

export interface DashboardNoticeItem {
  id: number
  title: string
  status: number
  updated_at: string
}

export interface DashboardData {
  current_user: DashboardCurrentUser
  health: DashboardHealth
  metrics: DashboardMetrics
  recent_operations: DashboardOperationItem[]
  recent_logins: DashboardLoginItem[]
  latest_notices: DashboardNoticeItem[]
}
