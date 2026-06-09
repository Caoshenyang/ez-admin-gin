import type {
  MessageReceiverType,
  MessageStatus,
  MessageTemplateType,
} from './message'

export interface MessageTemplateFormModel {
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
}

export interface MessageReminderFormModel {
  id: number
  code: string
  name: string
  trigger_event: string
  template_id: number | null
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
