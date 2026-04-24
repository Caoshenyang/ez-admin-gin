import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  CreateUserPayload,
  UpdateUserPayload,
  UpdateUserRolesPayload,
  UpdateUserStatusPayload,
  UserItem,
  UserListQuery,
  UserListResponse,
} from '../types/user'

export async function getUsers(params: UserListQuery) {
  const response = await http.get<ApiResponse<UserListResponse>>('/system/users', { params })
  return response.data.data
}

export async function createUser(payload: CreateUserPayload) {
  const response = await http.post<ApiResponse<UserItem>>('/system/users', payload)
  return response.data.data
}

export async function updateUser(id: number, payload: UpdateUserPayload) {
  const response = await http.post<ApiResponse<UserItem>>(`/system/users/${id}/update`, payload)
  return response.data.data
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
