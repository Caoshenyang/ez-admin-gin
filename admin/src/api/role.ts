import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  CreateRolePayload,
  RoleItem,
  RoleListQuery,
  RoleListResponse,
  UpdateRoleMenusPayload,
  UpdateRolePayload,
  UpdateRolePermissionsPayload,
  UpdateRoleStatusPayload,
} from '../types/role'

export async function getRoles(params: RoleListQuery) {
  const response = await http.get<ApiResponse<RoleListResponse>>('/system/roles', { params })
  return response.data.data
}

export async function createRole(payload: CreateRolePayload) {
  const response = await http.post<ApiResponse<RoleItem>>('/system/roles', payload)
  return response.data.data
}

export async function updateRole(id: number, payload: UpdateRolePayload) {
  const response = await http.post<ApiResponse<RoleItem>>(`/system/roles/${id}/update`, payload)
  return response.data.data
}

export async function updateRoleStatus(id: number, payload: UpdateRoleStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/roles/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function updateRolePermissions(id: number, payload: UpdateRolePermissionsPayload) {
  const response = await http.post<ApiResponse<{ id: number; permissions: unknown[] }>>(
    `/system/roles/${id}/permissions`,
    payload,
  )
  return response.data.data
}

export async function updateRoleMenus(id: number, payload: UpdateRoleMenusPayload) {
  const response = await http.post<ApiResponse<{ id: number; menu_ids: number[] }>>(
    `/system/roles/${id}/menus`,
    payload,
  )
  return response.data.data
}
