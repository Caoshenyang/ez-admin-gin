export const ConfigStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type ConfigStatus = (typeof ConfigStatus)[keyof typeof ConfigStatus]

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

export interface ConfigListQuery {
  page: number
  page_size: number
  keyword?: string
  group_code?: string
  status?: ConfigStatus | 0
}

export interface ConfigListResponse {
  items: ConfigItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateConfigPayload {
  group_code: string
  key: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
}

export interface UpdateConfigPayload {
  group_code: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
}

export interface UpdateConfigStatusPayload {
  status: ConfigStatus
}
