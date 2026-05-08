import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  NoticeItem,
  NoticeListQuery,
  NoticeListResponse,
  CreateNoticePayload,
  UpdateNoticePayload,
  UpdateNoticeStatusPayload,
} from '../types/notice'

// 获取公告列表（分页查询）
export async function getNotices(params: NoticeListQuery) {
  const response = await http.get<ApiResponse<NoticeListResponse>>('/system/notices', { params })
  return response.data.data
}

// 创建公告
export async function createNotice(payload: CreateNoticePayload) {
  const response = await http.post<ApiResponse<NoticeItem>>('/system/notices', payload)
  return response.data.data
}

// 更新公告
export async function updateNotice(id: number, payload: UpdateNoticePayload) {
  const response = await http.post<ApiResponse<NoticeItem>>(`/system/notices/${id}/update`, payload)
  return response.data.data
}

// 更新公告状态（启用/禁用）
export async function updateNoticeStatus(id: number, payload: UpdateNoticeStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/notices/${id}/status`,
    payload,
  )
  return response.data.data
}
