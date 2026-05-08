import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  ConfigItem,
  ConfigListQuery,
  ConfigListResponse,
  CreateConfigPayload,
  UpdateConfigPayload,
  UpdateConfigStatusPayload,
} from '../types/config'

// 获取系统配置列表（分页查询）
export async function getConfigs(params: ConfigListQuery) {
  const response = await http.get<ApiResponse<ConfigListResponse>>('/system/configs', { params })
  return response.data.data
}

// 创建系统配置
export async function createConfig(payload: CreateConfigPayload) {
  const response = await http.post<ApiResponse<ConfigItem>>('/system/configs', payload)
  return response.data.data
}

// 更新系统配置
export async function updateConfig(id: number, payload: UpdateConfigPayload) {
  const response = await http.post<ApiResponse<ConfigItem>>(`/system/configs/${id}/update`, payload)
  return response.data.data
}

// 更新系统配置状态（启用/禁用）
export async function updateConfigStatus(id: number, payload: UpdateConfigStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/configs/${id}/status`,
    payload,
  )
  return response.data.data
}
