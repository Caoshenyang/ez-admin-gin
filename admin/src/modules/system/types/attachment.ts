export const AttachmentStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

// AttachmentStatus 类型定义。
export type AttachmentStatus = (typeof AttachmentStatus)[keyof typeof AttachmentStatus]

// AttachmentItem 类型定义。
export interface AttachmentItem {
  id: number
  file_id: number
  display_name: string
  category: string
  biz_type: string
  original_name: string
  file_name: string
  ext: string
  mime_type: string
  size: number
  url: string
  uploader_id: number
  status: AttachmentStatus
  remark: string
  created_at: string
  updated_at: string
}

// AttachmentListQuery 类型定义。
export interface AttachmentListQuery {
  page: number
  page_size: number
  keyword?: string
  category?: string
  biz_type?: string
  ext?: string
  status?: AttachmentStatus | 0
}

// AttachmentListResponse 类型定义。
export interface AttachmentListResponse {
  items: AttachmentItem[]
  total: number
  page: number
  page_size: number
}

// CreateAttachmentPayload 类型定义。
export interface CreateAttachmentPayload {
  display_name?: string
  category?: string
  biz_type?: string
  status?: AttachmentStatus
  remark?: string
}

// UpdateAttachmentPayload 类型定义。
export interface UpdateAttachmentPayload {
  display_name: string
  category: string
  biz_type: string
  status: AttachmentStatus
  remark: string
}
