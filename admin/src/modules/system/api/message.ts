import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  CreateMessageReminderPayload,
  CreateMessageTemplatePayload,
  MessageReminderItem,
  MessageReminderListQuery,
  MessageReminderListResponse,
  MessageTemplateItem,
  MessageTemplateListQuery,
  MessageTemplateListResponse,
  UpdateMessageReminderPayload,
  UpdateMessageReminderStatusPayload,
  UpdateMessageTemplatePayload,
  UpdateMessageTemplateStatusPayload,
} from '../types/message'

export async function getMessageTemplates(params: MessageTemplateListQuery) {
  const response = await http.get<ApiResponse<MessageTemplateListResponse>>(
    '/system/message-templates',
    { params },
  )
  return response.data.data
}

export async function createMessageTemplate(payload: CreateMessageTemplatePayload) {
  const response = await http.post<ApiResponse<MessageTemplateItem>>(
    '/system/message-templates',
    payload,
  )
  return response.data.data
}

export async function updateMessageTemplate(id: number, payload: UpdateMessageTemplatePayload) {
  const response = await http.post<ApiResponse<MessageTemplateItem>>(
    `/system/message-templates/${id}/update`,
    payload,
  )
  return response.data.data
}

export async function updateMessageTemplateStatus(
  id: number,
  payload: UpdateMessageTemplateStatusPayload,
) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/message-templates/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function getMessageReminders(params: MessageReminderListQuery) {
  const response = await http.get<ApiResponse<MessageReminderListResponse>>(
    '/system/message-reminders',
    { params },
  )
  return response.data.data
}

export async function createMessageReminder(payload: CreateMessageReminderPayload) {
  const response = await http.post<ApiResponse<MessageReminderItem>>(
    '/system/message-reminders',
    payload,
  )
  return response.data.data
}

export async function updateMessageReminder(id: number, payload: UpdateMessageReminderPayload) {
  const response = await http.post<ApiResponse<MessageReminderItem>>(
    `/system/message-reminders/${id}/update`,
    payload,
  )
  return response.data.data
}

export async function updateMessageReminderStatus(
  id: number,
  payload: UpdateMessageReminderStatusPayload,
) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/message-reminders/${id}/status`,
    payload,
  )
  return response.data.data
}
