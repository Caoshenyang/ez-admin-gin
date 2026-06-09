export const NotificationType = {
  System: 1,
  Security: 2,
  Task: 3,
  Message: 4,
} as const

export type NotificationType = (typeof NotificationType)[keyof typeof NotificationType]

export interface NotificationItem {
  id: number
  type: NotificationType
  title: string
  content: string
  extra: Record<string, unknown> | null
  is_read: boolean
  created_at: string
  read_at: string | null
}

export interface NotificationListQuery {
  page: number
  page_size: number
  type?: NotificationType | 0
  is_read?: number // 0=all, 1=unread, 2=read
}

export interface NotificationListResponse {
  items: NotificationItem[]
  total: number
  page: number
  page_size: number
}

export interface UnreadCountResponse {
  count: number
}

export interface MarkReadPayload {
  ids: number[]
}

export interface WSMessage {
  type: 'notification' | 'unread_count' | 'ping'
  data?: unknown
}
