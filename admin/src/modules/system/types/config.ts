export const ConfigStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

// ConfigStatus 类型定义。
export type ConfigStatus = (typeof ConfigStatus)[keyof typeof ConfigStatus]

// ConfigItem 类型定义。
export interface ConfigItem {
  id: number
  group_code: string
  key: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
  created_at: string
  updated_at: string
}

// ConfigListQuery 类型定义。
export interface ConfigListQuery {
  page: number
  page_size: number
  keyword?: string
  group_code?: string
  status?: ConfigStatus | 0
}

// ConfigListResponse 类型定义。
export interface ConfigListResponse {
  items: ConfigItem[]
  total: number
  page: number
  page_size: number
}

// CreateConfigPayload 类型定义。
export interface CreateConfigPayload {
  group_code: string
  key: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
}

// UpdateConfigPayload 类型定义。
export interface UpdateConfigPayload {
  group_code: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
}

// UpdateConfigStatusPayload 类型定义。
export interface UpdateConfigStatusPayload {
  status: ConfigStatus
}
