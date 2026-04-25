export const NoticeStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type NoticeStatus = (typeof NoticeStatus)[keyof typeof NoticeStatus]

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

export interface NoticeListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: NoticeStatus | 0
}

export interface NoticeListResponse {
  items: NoticeItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateNoticePayload {
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

export interface UpdateNoticePayload {
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

export interface UpdateNoticeStatusPayload {
  status: NoticeStatus
}
