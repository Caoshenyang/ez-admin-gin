import type { SelectOption } from 'naive-ui'

// CommonStatusValue 通用状态值：1=启用，2=禁用。
type CommonStatusValue = 1 | 2

// StatusFilterOption 状态筛选下拉选项，支持"全部"(0)、启用(1)、禁用(2)。
export interface StatusFilterOption extends SelectOption {
  value: 0 | CommonStatusValue
}

// StatusFormOption 表单中的状态下拉选项，仅支持启用(1)或禁用(2)。
export interface StatusFormOption extends SelectOption {
  value: CommonStatusValue
}

// STATUS_FILTER_OPTIONS 列表页筛选栏使用的状态选项（含"全部"）。
export const STATUS_FILTER_OPTIONS: StatusFilterOption[] = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: 1 },
  { label: '禁用', value: 2 },
]

// STATUS_FORM_OPTIONS 新增/编辑表单中使用的状态选项（不含"全部"）。
export const STATUS_FORM_OPTIONS: StatusFormOption[] = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 2 },
]
