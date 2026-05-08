import type { FormRules } from 'naive-ui'

import { AttachmentStatus, type AttachmentItem, type AttachmentListQuery } from '../types/attachment'
import type { AttachmentEditFormModel, AttachmentUploadFormModel } from '../types/attachment-page'

export const attachmentUploadRules: FormRules = {
  display_name: [{ max: 255, message: '附件名称不能超过 255 个字符', trigger: ['blur', 'input'] }],
  category: [{ max: 64, message: '附件分类不能超过 64 个字符', trigger: ['blur', 'input'] }],
  biz_type: [{ max: 64, message: '业务类型不能超过 64 个字符', trigger: ['blur', 'input'] }],
  remark: [{ max: 255, message: '备注不能超过 255 个字符', trigger: ['blur', 'input'] }],
}

export const attachmentEditRules: FormRules = {
  display_name: [{ required: true, message: '请输入附件名称', trigger: ['blur', 'input'] }],
  category: [{ max: 64, message: '附件分类不能超过 64 个字符', trigger: ['blur', 'input'] }],
  biz_type: [{ max: 64, message: '业务类型不能超过 64 个字符', trigger: ['blur', 'input'] }],
  remark: [{ max: 255, message: '备注不能超过 255 个字符', trigger: ['blur', 'input'] }],
}

export const attachmentExtFilterOptions = [
  { label: '类型：全部', value: '' },
  { label: '图片', value: '.png' },
  { label: 'PDF', value: '.pdf' },
  { label: 'Excel', value: '.xlsx' },
  { label: 'Word', value: '.docx' },
]

export function defaultAttachmentQuery(): AttachmentListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    category: '',
    biz_type: '',
    ext: '',
    status: 0,
  }
}

export function defaultAttachmentUploadForm(): AttachmentUploadFormModel {
  return {
    display_name: '',
    category: '',
    biz_type: '',
    status: AttachmentStatus.Enabled,
    remark: '',
  }
}

export function defaultAttachmentEditForm(): AttachmentEditFormModel {
  return {
    id: 0,
    display_name: '',
    category: '',
    biz_type: '',
    status: AttachmentStatus.Enabled,
    remark: '',
  }
}

export function toAttachmentEditFormModel(item: AttachmentItem): AttachmentEditFormModel {
  return {
    id: item.id,
    display_name: item.display_name,
    category: item.category,
    biz_type: item.biz_type,
    status: item.status,
    remark: item.remark,
  }
}

export function buildAttachmentUploadPayload(formModel: AttachmentUploadFormModel) {
  return {
    display_name: formModel.display_name.trim() || undefined,
    category: formModel.category.trim() || undefined,
    biz_type: formModel.biz_type.trim() || undefined,
    status: formModel.status,
    remark: formModel.remark.trim() || undefined,
  }
}

export function buildAttachmentEditPayload(formModel: AttachmentEditFormModel) {
  return {
    display_name: formModel.display_name.trim(),
    category: formModel.category.trim(),
    biz_type: formModel.biz_type.trim(),
    status: formModel.status,
    remark: formModel.remark.trim(),
  }
}
