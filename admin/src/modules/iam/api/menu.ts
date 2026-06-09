import http from '@/api/http'

import type {
  AdminMenu,
  AuthMenu,
  CreateMenuPayload,
  UpdateMenuPayload,
  UpdateMenuStatusPayload,
} from '../types/menu'
import type { ApiResponse } from '@/api/types'

// getCurrentUserMenus 获取当前登录用户可见的菜单树。
export async function getCurrentUserMenus() {
  const response = await http.get<ApiResponse<AuthMenu[]>>('/auth/menus')
  return response.data.data ?? []
}

// getAdminMenus 查询所有菜单（管理端），返回完整的后台菜单树。
export async function getAdminMenus() {
  const response = await http.get<ApiResponse<AdminMenu[]>>('/system/menus')
  return response.data.data ?? []
}

// createMenu 创建菜单。
export async function createMenu(payload: CreateMenuPayload) {
  const response = await http.post<ApiResponse<AdminMenu>>('/system/menus', payload)
  return response.data.data
}

// updateMenu 更新菜单信息。
export async function updateMenu(id: number, payload: UpdateMenuPayload) {
  const response = await http.post<ApiResponse<AdminMenu>>(`/system/menus/${id}/update`, payload)
  return response.data.data
}

// updateMenuStatus 切换菜单的启用/禁用状态。
export async function updateMenuStatus(id: number, payload: UpdateMenuStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/menus/${id}/status`,
    payload,
  )
  return response.data.data
}

// deleteMenu 删除菜单。
export async function deleteMenu(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(`/system/menus/${id}/delete`)
  return response.data.data
}
