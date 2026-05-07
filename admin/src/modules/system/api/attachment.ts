import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  AttachmentItem,
  AttachmentListQuery,
  AttachmentListResponse,
  CreateAttachmentPayload,
  UpdateAttachmentPayload,
} from '../types/attachment'

export async function getAttachments(params: AttachmentListQuery) {
  const response = await http.get<ApiResponse<AttachmentListResponse>>('/system/attachments', { params })
  return response.data.data
}

export async function createAttachment(file: File, payload: CreateAttachmentPayload) {
  const formData = new FormData()
  formData.append('file', file)
  if (payload.display_name) {
    formData.append('display_name', payload.display_name)
  }
  if (payload.category) {
    formData.append('category', payload.category)
  }
  if (payload.biz_type) {
    formData.append('biz_type', payload.biz_type)
  }
  if (payload.status) {
    formData.append('status', String(payload.status))
  }
  if (payload.remark) {
    formData.append('remark', payload.remark)
  }

  const response = await http.post<ApiResponse<AttachmentItem>>('/system/attachments', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return response.data.data
}

export async function updateAttachment(id: number, payload: UpdateAttachmentPayload) {
  const response = await http.post<ApiResponse<AttachmentItem>>(`/system/attachments/${id}/update`, payload)
  return response.data.data
}

export async function updateAttachmentStatus(id: number, status: number) {
  const response = await http.post<ApiResponse<{ id: number, status: number }>>(
    `/system/attachments/${id}/status`,
    { status },
  )
  return response.data.data
}
