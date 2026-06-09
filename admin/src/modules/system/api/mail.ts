import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  CreateMailAccountPayload,
  CreateMailTemplatePayload,
  MailAccountItem,
  MailAccountListQuery,
  MailAccountListResponse,
  MailLogListQuery,
  MailLogListResponse,
  MailTemplateItem,
  MailTemplateListQuery,
  MailTemplateListResponse,
  RenderMailTemplatePayload,
  RenderMailTemplateResponse,
  SendMailPayload,
  SendMailResponse,
  TestMailAccountPayload,
  UpdateMailAccountPayload,
  UpdateMailStatusPayload,
  UpdateMailTemplatePayload,
} from '../types/mail'

export async function getMailAccounts(params: MailAccountListQuery) {
  const response = await http.get<ApiResponse<MailAccountListResponse>>('/system/mail/accounts', {
    params,
  })
  return response.data.data
}

export async function createMailAccount(payload: CreateMailAccountPayload) {
  const response = await http.post<ApiResponse<MailAccountItem>>('/system/mail/accounts', payload)
  return response.data.data
}

export async function updateMailAccount(id: number, payload: UpdateMailAccountPayload) {
  const response = await http.post<ApiResponse<MailAccountItem>>(
    `/system/mail/accounts/${id}/update`,
    payload,
  )
  return response.data.data
}

export async function updateMailAccountStatus(id: number, payload: UpdateMailStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/mail/accounts/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function deleteMailAccount(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(
    `/system/mail/accounts/${id}/delete`,
  )
  return response.data.data
}

export async function testMailAccount(id: number, payload: TestMailAccountPayload) {
  const response = await http.post<ApiResponse<SendMailResponse>>(
    `/system/mail/accounts/${id}/test`,
    payload,
  )
  return response.data.data
}

export async function getMailTemplates(params: MailTemplateListQuery) {
  const response = await http.get<ApiResponse<MailTemplateListResponse>>('/system/mail/templates', {
    params,
  })
  return response.data.data
}

export async function createMailTemplate(payload: CreateMailTemplatePayload) {
  const response = await http.post<ApiResponse<MailTemplateItem>>(
    '/system/mail/templates',
    payload,
  )
  return response.data.data
}

export async function updateMailTemplate(id: number, payload: UpdateMailTemplatePayload) {
  const response = await http.post<ApiResponse<MailTemplateItem>>(
    `/system/mail/templates/${id}/update`,
    payload,
  )
  return response.data.data
}

export async function updateMailTemplateStatus(id: number, payload: UpdateMailStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/mail/templates/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function deleteMailTemplate(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(
    `/system/mail/templates/${id}/delete`,
  )
  return response.data.data
}

export async function renderMailTemplate(id: number, payload: RenderMailTemplatePayload) {
  const response = await http.post<ApiResponse<RenderMailTemplateResponse>>(
    `/system/mail/templates/${id}/render`,
    payload,
  )
  return response.data.data
}

export async function sendMail(payload: SendMailPayload) {
  const response = await http.post<ApiResponse<SendMailResponse>>('/system/mail/send', payload)
  return response.data.data
}

export async function getMailLogs(params: MailLogListQuery) {
  const response = await http.get<ApiResponse<MailLogListResponse>>('/system/mail/logs', {
    params,
  })
  return response.data.data
}
