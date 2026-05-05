export const DictStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type DictStatus = (typeof DictStatus)[keyof typeof DictStatus]

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

export interface DictTypeListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: DictStatus | 0
}

export interface DictTypeListResponse {
  items: DictTypeItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateDictTypePayload {
  code: string
  name: string
  sort: number
  status: DictStatus
  remark: string
}

export interface UpdateDictTypePayload {
  name: string
  sort: number
  status: DictStatus
  remark: string
}

export interface UpdateDictTypeStatusPayload {
  status: DictStatus
}

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

export interface DictItemListQuery {
  page: number
  page_size: number
  type_id: number
  keyword?: string
  status?: DictStatus | 0
}

export interface DictItemListResponse {
  items: DictItem[]
  total: number
  page: number
  page_size: number
}

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

export interface UpdateDictItemPayload {
  label: string
  value: string
  tag_type: string
  sort: number
  status: DictStatus
  remark: string
}

export interface UpdateDictItemStatusPayload {
  status: DictStatus
}
