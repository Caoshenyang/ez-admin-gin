import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  ConfigItem,
  ConfigListQuery,
  ConfigListResponse,
  CreateConfigPayload,
  UpdateConfigPayload,
  UpdateConfigStatusPayload,
} from '../types/config'

export async function getConfigs(params: ConfigListQuery) {
  const response = await http.get<ApiResponse<ConfigListResponse>>('/system/configs', { params })
  return response.data.data
}

export async function createConfig(payload: CreateConfigPayload) {
  const response = await http.post<ApiResponse<ConfigItem>>('/system/configs', payload)
  return response.data.data
}

export async function updateConfig(id: number, payload: UpdateConfigPayload) {
  const response = await http.post<ApiResponse<ConfigItem>>(`/system/configs/${id}/update`, payload)
  return response.data.data
}

export async function updateConfigStatus(id: number, payload: UpdateConfigStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/configs/${id}/status`,
    payload,
  )
  return response.data.data
}
