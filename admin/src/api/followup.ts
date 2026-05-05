import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  CreateCustomerFollowUpPayload,
  CustomerFollowUpCustomerOption,
  CustomerFollowUpCustomerOptionQuery,
  CustomerFollowUpItem,
  CustomerFollowUpListQuery,
  CustomerFollowUpListResponse,
  UpdateCustomerFollowUpPayload,
} from '../types/followup'

export async function getCustomerFollowUps(params: CustomerFollowUpListQuery) {
  const response = await http.get<ApiResponse<CustomerFollowUpListResponse>>('/crm/followups', { params })
  return response.data.data
}

export async function getCustomerFollowUpCustomerOptions(params: CustomerFollowUpCustomerOptionQuery = {}) {
  const response = await http.get<ApiResponse<CustomerFollowUpCustomerOption[]>>(
    '/crm/followups/customer-options',
    { params },
  )
  return response.data.data
}

export async function createCustomerFollowUp(payload: CreateCustomerFollowUpPayload) {
  const response = await http.post<ApiResponse<CustomerFollowUpItem>>('/crm/followups', payload)
  return response.data.data
}

export async function updateCustomerFollowUp(id: number, payload: UpdateCustomerFollowUpPayload) {
  const response = await http.post<ApiResponse<CustomerFollowUpItem>>(`/crm/followups/${id}/update`, payload)
  return response.data.data
}

export async function updateCustomerFollowUpStatus(id: number, status: number) {
  const response = await http.post<ApiResponse<{ id: number, status: number }>>(
    `/crm/followups/${id}/status`,
    { status },
  )
  return response.data.data
}
