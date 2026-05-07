import type { SelectOption } from 'naive-ui'

type CommonStatusValue = 1 | 2

export interface StatusFilterOption extends SelectOption {
  value: 0 | CommonStatusValue
}

export interface StatusFormOption extends SelectOption {
  value: CommonStatusValue
}

export const STATUS_FILTER_OPTIONS: StatusFilterOption[] = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: 1 },
  { label: '禁用', value: 2 },
]

export const STATUS_FORM_OPTIONS: StatusFormOption[] = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 2 },
]
