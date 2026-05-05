import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  CreateDictItemPayload,
  CreateDictTypePayload,
  DictItem,
  DictItemListQuery,
  DictItemListResponse,
  DictTypeItem,
  DictTypeListQuery,
  DictTypeListResponse,
  UpdateDictItemPayload,
  UpdateDictItemStatusPayload,
  UpdateDictTypePayload,
  UpdateDictTypeStatusPayload,
} from '../types/dict'

export async function getDictTypes(params: DictTypeListQuery) {
  const response = await http.get<ApiResponse<DictTypeListResponse>>('/system/dict-types', { params })
  return response.data.data
}

export async function createDictType(payload: CreateDictTypePayload) {
  const response = await http.post<ApiResponse<DictTypeItem>>('/system/dict-types', payload)
  return response.data.data
}

export async function updateDictType(id: number, payload: UpdateDictTypePayload) {
  const response = await http.post<ApiResponse<DictTypeItem>>(`/system/dict-types/${id}/update`, payload)
  return response.data.data
}

export async function updateDictTypeStatus(id: number, payload: UpdateDictTypeStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/dict-types/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function getDictItems(params: DictItemListQuery) {
  const response = await http.get<ApiResponse<DictItemListResponse>>('/system/dict-items', { params })
  return response.data.data
}

export async function createDictItem(payload: CreateDictItemPayload) {
  const response = await http.post<ApiResponse<DictItem>>('/system/dict-items', payload)
  return response.data.data
}

export async function updateDictItem(id: number, payload: UpdateDictItemPayload) {
  const response = await http.post<ApiResponse<DictItem>>(`/system/dict-items/${id}/update`, payload)
  return response.data.data
}

export async function updateDictItemStatus(id: number, payload: UpdateDictItemStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/dict-items/${id}/status`,
    payload,
  )
  return response.data.data
}
