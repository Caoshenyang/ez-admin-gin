import {
  NoticeStatus,
  type CreateNoticePayload,
  type NoticeItem,
  type NoticeListQuery,
  type UpdateNoticePayload,
} from '../types/notice'

export interface NoticeFormModel {
  id: number
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

export function defaultNoticeListQuery(): NoticeListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    status: 0,
  }
}

export function defaultNoticeFormModel(): NoticeFormModel {
  return {
    id: 0,
    title: '',
    content: '',
    sort: 0,
    status: NoticeStatus.Enabled,
    remark: '',
  }
}

export function toNoticeFormModel(notice: NoticeItem): NoticeFormModel {
  return {
    id: notice.id,
    title: notice.title,
    content: notice.content,
    sort: notice.sort,
    status: notice.status,
    remark: notice.remark,
  }
}

export function buildNoticePayload(
  formModel: NoticeFormModel,
): CreateNoticePayload | UpdateNoticePayload {
  return {
    title: formModel.title,
    content: formModel.content,
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark,
  }
}
