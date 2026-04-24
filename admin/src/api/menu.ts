import http from './http'

import type { AuthMenu } from '../types/menu'
import type { ApiResponse } from '../types/http'

// getCurrentUserMenus 获取当前登录用户可见的菜单树。
export async function getCurrentUserMenus() {
  const response = await http.get<ApiResponse<AuthMenu[]>>('/auth/menus')
  return response.data.data ?? []
}
