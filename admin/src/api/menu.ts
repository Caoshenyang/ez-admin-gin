import http from './http'

import type {
  AdminMenu,
  AuthMenu,
  CreateMenuPayload,
  UpdateMenuPayload,
  UpdateMenuStatusPayload,
} from '../types/menu'
import type { ApiResponse } from '../types/http'

// getCurrentUserMenus 获取当前登录用户可见的菜单树。
export async function getCurrentUserMenus() {
  const response = await http.get<ApiResponse<AuthMenu[]>>('/auth/menus')
  return response.data.data ?? []
}

export async function getAdminMenus() {
  const response = await http.get<ApiResponse<AdminMenu[]>>('/system/menus')
  return response.data.data ?? []
}

export async function createMenu(payload: CreateMenuPayload) {
  const response = await http.post<ApiResponse<AdminMenu>>('/system/menus', payload)
  return response.data.data
}

export async function updateMenu(id: number, payload: UpdateMenuPayload) {
  const response = await http.post<ApiResponse<AdminMenu>>(`/system/menus/${id}/update`, payload)
  return response.data.data
}

export async function updateMenuStatus(id: number, payload: UpdateMenuStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/menus/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function deleteMenu(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(`/system/menus/${id}/delete`)
  return response.data.data
}
