export const CustomerStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type CustomerStatus = (typeof CustomerStatus)[keyof typeof CustomerStatus]

export interface CustomerItem {
  id: number
  name: string
  contact_name: string
  phone: string
  level: string
  source: string
  department_id: number
  department_name: string
  owner_user_id: number
  owner_username: string
  owner_nickname: string
  status: CustomerStatus
  remark: string
  created_at: string
  updated_at: string
}

export interface CustomerListQuery {
  page: number
  page_size: number
  keyword?: string
  level?: string
  source?: string
  status?: CustomerStatus | 0
}

export interface CustomerListResponse {
  items: CustomerItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateCustomerPayload {
  name: string
  contact_name: string
  phone: string
  level: string
  source: string
  status: CustomerStatus
  remark: string
}

export type UpdateCustomerPayload = CreateCustomerPayload
