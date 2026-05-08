import type { FormRules, UploadFileInfo } from 'naive-ui'

import type { AttachmentItem, AttachmentListQuery, AttachmentStatus, UpdateAttachmentPayload } from './attachment'

export interface AttachmentUploadFormModel {
  display_name: string
  category: string
  biz_type: string
  status: AttachmentStatus
  remark: string
}

export interface AttachmentEditFormModel extends UpdateAttachmentPayload {
  id: number
}

export type AttachmentPageQuery = AttachmentListQuery
export type AttachmentPageItem = AttachmentItem
export type AttachmentUploadFileList = UploadFileInfo[]
export type AttachmentFormRules = FormRules
