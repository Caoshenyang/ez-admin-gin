import http from '@/api/http'
import type { ApiResponse } from '@/api/types'
import type {
  NotificationItem,
  NotificationListQuery,
  NotificationListResponse,
  UnreadCountResponse,
  MarkReadPayload,
} from '@/types/notification'

export async function getNotifications(params: NotificationListQuery) {
  const response = await http.get<ApiResponse<NotificationListResponse>>('/system/notifications', { params })
  return response.data.data
}

export async function getUnreadCount() {
  const response = await http.get<ApiResponse<UnreadCountResponse>>('/system/notifications/unread-count')
  return response.data.data
}

export async function markRead(payload: MarkReadPayload) {
  const response = await http.post<ApiResponse<null>>('/system/notifications/mark-read', payload)
  return response.data
}

export async function markAllRead() {
  const response = await http.post<ApiResponse<null>>('/system/notifications/mark-all-read')
  return response.data
}
