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

function normalizeUserItem(item: UserItem): UserItem {
  return {
    ...item,
    role_ids: Array.isArray(item.role_ids) ? item.role_ids : [],
    post_ids: Array.isArray(item.post_ids) ? item.post_ids : [],
  }
}

export async function getUsers(params: UserListQuery) {
  const response = await http.get<ApiResponse<UserListResponse>>('/system/users', { params })
  const data = response.data.data
  return {
    ...data,
    items: data.items.map(normalizeUserItem),
  }
}

export async function createUser(payload: CreateUserPayload) {
  const response = await http.post<ApiResponse<UserItem>>('/system/users', payload)
  return normalizeUserItem(response.data.data)
}

export async function updateUser(id: number, payload: UpdateUserPayload) {
  const response = await http.post<ApiResponse<UserItem>>(`/system/users/${id}/update`, payload)
  return normalizeUserItem(response.data.data)
}

export async function updateUserStatus(id: number, payload: UpdateUserStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/users/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function updateUserRoles(id: number, payload: UpdateUserRolesPayload) {
  const response = await http.post<ApiResponse<{ id: number; role_ids: number[] }>>(
    `/system/users/${id}/roles`,
    payload,
  )
  return response.data.data
}
