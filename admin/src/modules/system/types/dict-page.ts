import type {
  CreateDictItemPayload,
  CreateDictTypePayload,
  DictItem,
  DictItemListQuery,
  DictStatus,
  DictTypeItem,
  DictTypeListQuery,
  UpdateDictItemPayload,
  UpdateDictTypePayload,
} from './dict'

export interface DictTypeFormModel {
  id: number
  code: string
  name: string
  sort: number
  status: DictStatus
  remark: string
}

export interface DictItemFormModel {
  id: number
  type_id: number
  item_key: string
  label: string
  value: string
  tag_type: string
  sort: number
  status: DictStatus
  remark: string
}

export type DictTypeCreatePayload = CreateDictTypePayload
export type DictTypeUpdatePayload = UpdateDictTypePayload
export type DictItemCreatePayload = CreateDictItemPayload
export type DictItemUpdatePayload = UpdateDictItemPayload
export type DictTypePageQuery = DictTypeListQuery
export type DictItemPageQuery = DictItemListQuery
export type DictTypePageItem = DictTypeItem
export type DictItemPageItem = DictItem
