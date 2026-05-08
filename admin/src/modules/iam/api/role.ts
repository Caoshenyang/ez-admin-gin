import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
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

// getRoles 分页查询角色列表，附带每个角色的权限和菜单信息。
export async function getRoles(params: RoleListQuery) {
  const response = await http.get<ApiResponse<RoleListResponse>>('/system/roles', { params })
  return response.data.data
}

// createRole 创建角色。
export async function createRole(payload: CreateRolePayload) {
  const response = await http.post<ApiResponse<RoleItem>>('/system/roles', payload)
  return response.data.data
}

// updateRole 更新角色基本信息和数据范围。
export async function updateRole(id: number, payload: UpdateRolePayload) {
  const response = await http.post<ApiResponse<RoleItem>>(`/system/roles/${id}/update`, payload)
  return response.data.data
}

// updateRoleStatus 切换角色的启用/禁用状态。
export async function updateRoleStatus(id: number, payload: UpdateRoleStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/roles/${id}/status`,
    payload,
  )
  return response.data.data
}

// updateRolePermissions 更新角色的 Casbin 权限策略。
export async function updateRolePermissions(id: number, payload: UpdateRolePermissionsPayload) {
  const response = await http.post<ApiResponse<{ id: number; permissions: unknown[] }>>(
    `/system/roles/${id}/permissions`,
    payload,
  )
  return response.data.data
}

// updateRoleMenus 更新角色的菜单分配。
export async function updateRoleMenus(id: number, payload: UpdateRoleMenusPayload) {
  const response = await http.post<ApiResponse<{ id: number; menu_ids: number[] }>>(
    `/system/roles/${id}/menus`,
    payload,
  )
  return response.data.data
}
