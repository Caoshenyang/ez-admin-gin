export const MessageStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type MessageStatus = (typeof MessageStatus)[keyof typeof MessageStatus]

export const MessageTemplateType = {
  Notification: 1,
  Todo: 2,
  Alert: 3,
} as const

export type MessageTemplateType =
  (typeof MessageTemplateType)[keyof typeof MessageTemplateType]

export const MessageReceiverType = {
  Role: 1,
  User: 2,
  Department: 3,
  Initiator: 4,
  Assignee: 5,
} as const

export type MessageReceiverType =
  (typeof MessageReceiverType)[keyof typeof MessageReceiverType]

export interface MessageTemplateItem {
  id: number
  code: string
  name: string
  title: string
  content: string
  type: MessageTemplateType
  variables: string
  sort: number
  status: MessageStatus
  is_system: boolean
  remark: string
  created_at: string
  updated_at: string
}

export interface MessageTemplateListQuery {
  page: number
  page_size: number
  keyword?: string
  type?: MessageTemplateType | 0
  status?: MessageStatus | 0
}

export interface MessageTemplateListResponse {
  items: MessageTemplateItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateMessageTemplatePayload {
  code: string
  name: string
  title: string
  content: string
  type: MessageTemplateType
  variables: string
  sort: number
  status: MessageStatus
  is_system: boolean
  remark: string
}

export interface UpdateMessageTemplatePayload {
  name: string
  title: string
  content: string
  type: MessageTemplateType
  variables: string
  sort: number
  status: MessageStatus
  is_system: boolean
  remark: string
}

export interface UpdateMessageTemplateStatusPayload {
  status: MessageStatus
}

export interface MessageReminderItem {
  id: number
  code: string
  name: string
  trigger_event: string
  template_id: number
  template_code: string
  template_name: string
  channels: string
  receiver_type: MessageReceiverType
  receiver_values: string
  advance_minutes: number
  link_url: string
  sort: number
  status: MessageStatus
  is_system: boolean
  remark: string
  created_at: string
  updated_at: string
}

export interface MessageReminderListQuery {
  page: number
  page_size: number
  keyword?: string
  trigger_event?: string
  template_id?: number
  receiver_type?: MessageReceiverType | 0
  status?: MessageStatus | 0
}

export interface MessageReminderListResponse {
  items: MessageReminderItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateMessageReminderPayload {
  code: string
  name: string
  trigger_event: string
  template_id: number
  channels: string
  receiver_type: MessageReceiverType
  receiver_values: string
  advance_minutes: number
  link_url: string
  sort: number
  status: MessageStatus
  is_system: boolean
  remark: string
}

export interface UpdateMessageReminderPayload {
  name: string
  trigger_event: string
  template_id: number
  channels: string
  receiver_type: MessageReceiverType
  receiver_values: string
  advance_minutes: number
  link_url: string
  sort: number
  status: MessageStatus
  is_system: boolean
  remark: string
}

export interface UpdateMessageReminderStatusPayload {
  status: MessageStatus
}
