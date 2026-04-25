import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  NoticeItem,
  NoticeListQuery,
  NoticeListResponse,
  CreateNoticePayload,
  UpdateNoticePayload,
  UpdateNoticeStatusPayload,
} from '../types/notice'

export async function getNotices(params: NoticeListQuery) {
  const response = await http.get<ApiResponse<NoticeListResponse>>('/system/notices', { params })
  return response.data.data
}

export async function createNotice(payload: CreateNoticePayload) {
  const response = await http.post<ApiResponse<NoticeItem>>('/system/notices', payload)
  return response.data.data
}

export async function updateNotice(id: number, payload: UpdateNoticePayload) {
  const response = await http.post<ApiResponse<NoticeItem>>(`/system/notices/${id}/update`, payload)
  return response.data.data
}

export async function updateNoticeStatus(id: number, payload: UpdateNoticeStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/notices/${id}/status`,
    payload,
  )
  return response.data.data
}
