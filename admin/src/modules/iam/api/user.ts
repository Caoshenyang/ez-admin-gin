import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  CreateUserPayload,
  UpdateUserPayload,
  UpdateUserRolesPayload,
  UpdateUserStatusPayload,
  UserItem,
  UserListQuery,
  UserListResponse,
} from '../types/user'

// normalizeUserItem 确保用户的 role_ids 和 post_ids 始终是数组，防止后端返回 null。
function normalizeUserItem(item: UserItem): UserItem {
  return {
    ...item,
    role_ids: Array.isArray(item.role_ids) ? item.role_ids : [],
    post_ids: Array.isArray(item.post_ids) ? item.post_ids : [],
  }
}

// getUsers 分页查询用户列表，附带角色和岗位信息。
export async function getUsers(params: UserListQuery) {
  const response = await http.get<ApiResponse<UserListResponse>>('/system/users', { params })
  const data = response.data.data
  return {
    ...data,
    items: data.items.map(normalizeUserItem),
  }
}

// createUser 创建用户并关联角色和岗位。
export async function createUser(payload: CreateUserPayload) {
  const response = await http.post<ApiResponse<UserItem>>('/system/users', payload)
  return normalizeUserItem(response.data.data)
}

// updateUser 更新用户基本信息和岗位关联。
export async function updateUser(id: number, payload: UpdateUserPayload) {
  const response = await http.post<ApiResponse<UserItem>>(`/system/users/${id}/update`, payload)
  return normalizeUserItem(response.data.data)
}

// updateUserStatus 切换用户的启用/禁用状态。
export async function updateUserStatus(id: number, payload: UpdateUserStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/users/${id}/status`,
    payload,
  )
  return response.data.data
}

// updateUserRoles 更新用户的角色分配。
export async function updateUserRoles(id: number, payload: UpdateUserRolesPayload) {
  const response = await http.post<ApiResponse<{ id: number; role_ids: number[] }>>(
    `/system/users/${id}/roles`,
    payload,
  )
  return response.data.data
}

// deleteUser 删除指定用户。
export async function deleteUser(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(`/system/users/${id}/delete`)
  return response.data.data
}
