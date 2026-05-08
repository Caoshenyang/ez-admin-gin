export const NoticeStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

// NoticeStatus 类型定义。
export type NoticeStatus = (typeof NoticeStatus)[keyof typeof NoticeStatus]

// NoticeItem 类型定义。
export interface NoticeItem {
  id: number
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
  created_at: string
  updated_at: string
}

// NoticeListQuery 类型定义。
export interface NoticeListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: NoticeStatus | 0
}

// NoticeListResponse 类型定义。
export interface NoticeListResponse {
  items: NoticeItem[]
  total: number
  page: number
  page_size: number
}

// CreateNoticePayload 类型定义。
export interface CreateNoticePayload {
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

// UpdateNoticePayload 类型定义。
export interface UpdateNoticePayload {
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

// UpdateNoticeStatusPayload 类型定义。
export interface UpdateNoticeStatusPayload {
  status: NoticeStatus
}
