import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  CreateCustomerPayload,
  CustomerItem,
  CustomerListQuery,
  CustomerListResponse,
  UpdateCustomerPayload,
} from '../types/customer'

export async function getCustomers(params: CustomerListQuery) {
  const response = await http.get<ApiResponse<CustomerListResponse>>('/crm/customers', { params })
  return response.data.data
}

export async function createCustomer(payload: CreateCustomerPayload) {
  const response = await http.post<ApiResponse<CustomerItem>>('/crm/customers', payload)
  return response.data.data
}

export async function updateCustomer(id: number, payload: UpdateCustomerPayload) {
  const response = await http.post<ApiResponse<CustomerItem>>(`/crm/customers/${id}/update`, payload)
  return response.data.data
}

export async function updateCustomerStatus(id: number, status: number) {
  const response = await http.post<ApiResponse<{ id: number, status: number }>>(
    `/crm/customers/${id}/status`,
    { status },
  )
  return response.data.data
}
