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

// normalizeRoleItem 确保角色关联数组始终可直接用于表单和树控件。
function normalizeRoleItem(item: RoleItem): RoleItem {
  return {
    ...item,
    custom_department_ids: Array.isArray(item.custom_department_ids)
      ? item.custom_department_ids
      : [],
    api_ids: Array.isArray(item.api_ids) ? item.api_ids : [],
    menu_ids: Array.isArray(item.menu_ids) ? item.menu_ids : [],
    permissions: Array.isArray(item.permissions) ? item.permissions : [],
  }
}

// getRoles 分页查询角色列表，附带每个角色的权限和菜单信息。
export async function getRoles(params: RoleListQuery) {
  const response = await http.get<ApiResponse<RoleListResponse>>('/system/roles', { params })
  const data = response.data.data
  return {
    ...data,
    items: data.items.map(normalizeRoleItem),
  }
}

// createRole 创建角色。
export async function createRole(payload: CreateRolePayload) {
  const response = await http.post<ApiResponse<RoleItem>>('/system/roles', payload)
  return normalizeRoleItem(response.data.data)
}

// updateRole 更新角色基本信息和数据范围。
export async function updateRole(id: number, payload: UpdateRolePayload) {
  const response = await http.post<ApiResponse<RoleItem>>(`/system/roles/${id}/update`, payload)
  return normalizeRoleItem(response.data.data)
}

// updateRoleStatus 切换角色的启用/禁用状态。
export async function updateRoleStatus(id: number, payload: UpdateRoleStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/roles/${id}/status`,
    payload,
  )
  return response.data.data
}

// updateRolePermissions 更新角色的接口权限关联，后端会同步 Casbin 策略。
export async function updateRolePermissions(id: number, payload: UpdateRolePermissionsPayload) {
  const response = await http.post<ApiResponse<{ id: number; code: string; api_ids: number[] }>>(
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

// deleteRole 删除指定角色。
export async function deleteRole(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(`/system/roles/${id}/delete`)
  return response.data.data
}
