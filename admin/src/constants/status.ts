export const STATUS_FILTER_OPTIONS = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: 1 },
  { label: '禁用', value: 2 },
] as const

export const STATUS_FORM_OPTIONS = STATUS_FILTER_OPTIONS.slice(1)
