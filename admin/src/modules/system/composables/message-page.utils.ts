import type { SelectOption } from 'naive-ui'

import {
  MessageReceiverType,
  MessageStatus,
  MessageTemplateType,
  type CreateMessageReminderPayload,
  type CreateMessageTemplatePayload,
  type MessageReminderItem,
  type MessageReminderListQuery,
  type MessageTemplateItem,
  type MessageTemplateListQuery,
  type UpdateMessageReminderPayload,
  type UpdateMessageTemplatePayload,
} from '../types/message'
import type { MessageReminderFormModel, MessageTemplateFormModel } from '../types/message-page'

export const MESSAGE_TEMPLATE_TYPE_OPTIONS: SelectOption[] = [
  { label: '站内通知', value: MessageTemplateType.Notification },
  { label: '待办提醒', value: MessageTemplateType.Todo },
  { label: '告警提醒', value: MessageTemplateType.Alert },
]

export const MESSAGE_TEMPLATE_TYPE_FILTER_OPTIONS: SelectOption[] = [
  { label: '类型：全部', value: 0 },
  ...MESSAGE_TEMPLATE_TYPE_OPTIONS,
]

export const MESSAGE_RECEIVER_TYPE_OPTIONS: SelectOption[] = [
  { label: '按角色', value: MessageReceiverType.Role },
  { label: '指定用户', value: MessageReceiverType.User },
  { label: '按部门', value: MessageReceiverType.Department },
  { label: '业务发起人', value: MessageReceiverType.Initiator },
  { label: '业务负责人', value: MessageReceiverType.Assignee },
]

export const MESSAGE_RECEIVER_TYPE_FILTER_OPTIONS: SelectOption[] = [
  { label: '接收人：全部', value: 0 },
  ...MESSAGE_RECEIVER_TYPE_OPTIONS,
]

export function defaultMessageTemplateListQuery(): MessageTemplateListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    type: 0,
    status: 0,
  }
}

export function defaultMessageReminderListQuery(): MessageReminderListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    trigger_event: '',
    template_id: undefined,
    receiver_type: 0,
    status: 0,
  }
}

export function defaultMessageTemplateFormModel(): MessageTemplateFormModel {
  return {
    id: 0,
    code: '',
    name: '',
    title: '',
    content: '',
    type: MessageTemplateType.Notification,
    variables: '',
    sort: 0,
    status: MessageStatus.Enabled,
    is_system: true,
    remark: '',
  }
}

export function defaultMessageReminderFormModel(): MessageReminderFormModel {
  return {
    id: 0,
    code: '',
    name: '',
    trigger_event: '',
    template_id: null,
    channels: 'notification',
    receiver_type: MessageReceiverType.Role,
    receiver_values: '',
    advance_minutes: 0,
    link_url: '',
    sort: 0,
    status: MessageStatus.Enabled,
    is_system: true,
    remark: '',
  }
}

export function toMessageTemplateFormModel(item: MessageTemplateItem): MessageTemplateFormModel {
  return {
    id: item.id,
    code: item.code,
    name: item.name,
    title: item.title,
    content: item.content,
    type: item.type,
    variables: item.variables,
    sort: item.sort,
    status: item.status,
    is_system: item.is_system,
    remark: item.remark,
  }
}

export function toMessageReminderFormModel(item: MessageReminderItem): MessageReminderFormModel {
  return {
    id: item.id,
    code: item.code,
    name: item.name,
    trigger_event: item.trigger_event,
    template_id: item.template_id,
    channels: item.channels,
    receiver_type: item.receiver_type,
    receiver_values: item.receiver_values,
    advance_minutes: item.advance_minutes,
    link_url: item.link_url,
    sort: item.sort,
    status: item.status,
    is_system: item.is_system,
    remark: item.remark,
  }
}

export function buildMessageTemplateCreatePayload(
  formModel: MessageTemplateFormModel,
): CreateMessageTemplatePayload {
  return {
    code: formModel.code,
    name: formModel.name,
    title: formModel.title,
    content: formModel.content,
    type: formModel.type,
    variables: formModel.variables,
    sort: formModel.sort,
    status: formModel.status,
    is_system: formModel.is_system,
    remark: formModel.remark,
  }
}

export function buildMessageTemplateUpdatePayload(
  formModel: MessageTemplateFormModel,
): UpdateMessageTemplatePayload {
  return {
    name: formModel.name,
    title: formModel.title,
    content: formModel.content,
    type: formModel.type,
    variables: formModel.variables,
    sort: formModel.sort,
    status: formModel.status,
    is_system: formModel.is_system,
    remark: formModel.remark,
  }
}

export function buildMessageReminderCreatePayload(
  formModel: MessageReminderFormModel,
): CreateMessageReminderPayload {
  return {
    code: formModel.code,
    name: formModel.name,
    trigger_event: formModel.trigger_event,
    template_id: formModel.template_id ?? 0,
    channels: formModel.channels,
    receiver_type: formModel.receiver_type,
    receiver_values: formModel.receiver_values,
    advance_minutes: formModel.advance_minutes,
    link_url: formModel.link_url,
    sort: formModel.sort,
    status: formModel.status,
    is_system: formModel.is_system,
    remark: formModel.remark,
  }
}

export function buildMessageReminderUpdatePayload(
  formModel: MessageReminderFormModel,
): UpdateMessageReminderPayload {
  return {
    name: formModel.name,
    trigger_event: formModel.trigger_event,
    template_id: formModel.template_id ?? 0,
    channels: formModel.channels,
    receiver_type: formModel.receiver_type,
    receiver_values: formModel.receiver_values,
    advance_minutes: formModel.advance_minutes,
    link_url: formModel.link_url,
    sort: formModel.sort,
    status: formModel.status,
    is_system: formModel.is_system,
    remark: formModel.remark,
  }
}

export function templateTypeLabel(value: MessageTemplateType) {
  return findOptionLabel(MESSAGE_TEMPLATE_TYPE_OPTIONS, value)
}

export function receiverTypeLabel(value: MessageReceiverType) {
  return findOptionLabel(MESSAGE_RECEIVER_TYPE_OPTIONS, value)
}

export function receiverValuePlaceholder(value: MessageReceiverType) {
  if (value === MessageReceiverType.Role) return '角色编码，多个用英文逗号分隔'
  if (value === MessageReceiverType.User) return '用户 ID，多个用英文逗号分隔'
  if (value === MessageReceiverType.Department) return '部门 ID，多个用英文逗号分隔'
  return '无需填写，发送时由业务上下文决定'
}

function findOptionLabel(options: SelectOption[], value: number) {
  const matched = options.find((option) => option.value === value)
  return typeof matched?.label === 'string' ? matched.label : '-'
}
