export const DictStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

// DictStatus 类型定义。
export type DictStatus = (typeof DictStatus)[keyof typeof DictStatus]

// DictTypeItem 类型定义。
export interface DictTypeItem {
  id: number
  code: string
  name: string
  sort: number
  status: DictStatus
  remark: string
  created_at: string
  updated_at: string
}

// DictTypeListQuery 类型定义。
export interface DictTypeListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: DictStatus | 0
}

// DictTypeListResponse 类型定义。
export interface DictTypeListResponse {
  items: DictTypeItem[]
  total: number
  page: number
  page_size: number
}

// CreateDictTypePayload 类型定义。
export interface CreateDictTypePayload {
  code: string
  name: string
  sort: number
  status: DictStatus
  remark: string
}

// UpdateDictTypePayload 类型定义。
export interface UpdateDictTypePayload {
  name: string
  sort: number
  status: DictStatus
  remark: string
}

// UpdateDictTypeStatusPayload 类型定义。
export interface UpdateDictTypeStatusPayload {
  status: DictStatus
}

// DictItem 类型定义。
export interface DictItem {
  id: number
  type_id: number
  item_key: string
  label: string
  value: string
  tag_type: string
  sort: number
  status: DictStatus
  remark: string
  created_at: string
  updated_at: string
}

// DictItemListQuery 类型定义。
export interface DictItemListQuery {
  page: number
  page_size: number
  type_id: number
  keyword?: string
  status?: DictStatus | 0
}

// DictItemListResponse 类型定义。
export interface DictItemListResponse {
  items: DictItem[]
  total: number
  page: number
  page_size: number
}

// CreateDictItemPayload 类型定义。
export interface CreateDictItemPayload {
  type_id: number
  item_key: string
  label: string
  value: string
  tag_type: string
  sort: number
  status: DictStatus
  remark: string
}

// UpdateDictItemPayload 类型定义。
export interface UpdateDictItemPayload {
  label: string
  value: string
  tag_type: string
  sort: number
  status: DictStatus
  remark: string
}

// UpdateDictItemStatusPayload 类型定义。
export interface UpdateDictItemStatusPayload {
  status: DictStatus
}
