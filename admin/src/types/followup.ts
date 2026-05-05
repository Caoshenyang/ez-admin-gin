export const CustomerFollowUpStatus = {
  Pending: 1,
  Completed: 2,
  Closed: 3,
} as const

export type CustomerFollowUpStatus = (typeof CustomerFollowUpStatus)[keyof typeof CustomerFollowUpStatus]

export interface CustomerFollowUpItem {
  id: number
  customer_id: number
  customer_name: string
  department_id: number
  department_name: string
  owner_user_id: number
  owner_username: string
  owner_nickname: string
  follow_type: string
  subject: string
  content: string
  result: string
  next_follow_at: string | null
  status: CustomerFollowUpStatus
  created_at: string
  updated_at: string
}

export interface CustomerFollowUpListQuery {
  page: number
  page_size: number
  keyword?: string
  follow_type?: string
  customer_id?: number
  status?: CustomerFollowUpStatus | 0
}

export interface CustomerFollowUpListResponse {
  items: CustomerFollowUpItem[]
  total: number
  page: number
  page_size: number
}

export interface CustomerFollowUpCustomerOption {
  id: number
  name: string
  department_id: number
  department_name: string
  owner_user_id: number
  owner_username: string
  owner_nickname: string
}

export interface CustomerFollowUpCustomerOptionQuery {
  keyword?: string
  limit?: number
}

export interface CreateCustomerFollowUpPayload {
  customer_id: number
  follow_type: string
  subject: string
  content: string
  result: string
  next_follow_at: string | null
  status: CustomerFollowUpStatus
}

export interface UpdateCustomerFollowUpPayload {
  follow_type: string
  subject: string
  content: string
  result: string
  next_follow_at: string | null
  status: CustomerFollowUpStatus
}
