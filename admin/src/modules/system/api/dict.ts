import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
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

// 获取字典类型列表（分页查询）
export async function getDictTypes(params: DictTypeListQuery) {
  const response = await http.get<ApiResponse<DictTypeListResponse>>('/system/dict-types', {
    params,
  })
  return response.data.data
}

// 创建字典类型
export async function createDictType(payload: CreateDictTypePayload) {
  const response = await http.post<ApiResponse<DictTypeItem>>('/system/dict-types', payload)
  return response.data.data
}

// 更新字典类型
export async function updateDictType(id: number, payload: UpdateDictTypePayload) {
  const response = await http.post<ApiResponse<DictTypeItem>>(
    `/system/dict-types/${id}/update`,
    payload,
  )
  return response.data.data
}

// 更新字典类型状态（启用/禁用）
export async function updateDictTypeStatus(id: number, payload: UpdateDictTypeStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/dict-types/${id}/status`,
    payload,
  )
  return response.data.data
}

// 删除字典类型
export async function deleteDictType(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(`/system/dict-types/${id}/delete`)
  return response.data.data
}

// 获取字典项列表（分页查询）
export async function getDictItems(params: DictItemListQuery) {
  const response = await http.get<ApiResponse<DictItemListResponse>>('/system/dict-items', {
    params,
  })
  return response.data.data
}

// 创建字典项
export async function createDictItem(payload: CreateDictItemPayload) {
  const response = await http.post<ApiResponse<DictItem>>('/system/dict-items', payload)
  return response.data.data
}

// 更新字典项
export async function updateDictItem(id: number, payload: UpdateDictItemPayload) {
  const response = await http.post<ApiResponse<DictItem>>(
    `/system/dict-items/${id}/update`,
    payload,
  )
  return response.data.data
}

// 更新字典项状态（启用/禁用）
export async function updateDictItemStatus(id: number, payload: UpdateDictItemStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/dict-items/${id}/status`,
    payload,
  )
  return response.data.data
}

// 删除字典项
export async function deleteDictItem(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(`/system/dict-items/${id}/delete`)
  return response.data.data
}
