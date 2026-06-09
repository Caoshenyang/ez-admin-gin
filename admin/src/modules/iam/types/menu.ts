// 菜单类型常量：目录、菜单、按钮
export const MenuType = {
  Directory: 1,
  Menu: 2,
  Button: 3,
} as const

// 菜单类型联合类型
export type MenuType = (typeof MenuType)[keyof typeof MenuType]

// 菜单状态常量：启用、禁用
export const MenuStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

// 菜单状态联合类型
export type MenuStatus = (typeof MenuStatus)[keyof typeof MenuStatus]

// AuthMenu 对应 /api/v1/auth/menus 返回的菜单节点。
export interface AuthMenu {
  id: number
  parent_id: number
  type: MenuType
  code: string
  title: string
  path: string
  component: string
  icon: string
  sort: number
  children?: AuthMenu[]
}

// AdminMenu 后台管理菜单节点，包含状态、备注和时间信息
export interface AdminMenu {
  id: number
  parent_id: number
  type: MenuType
  code: string
  title: string
  path: string
  component: string
  icon: string
  sort: number
  status: MenuStatus
  remark: string
  children?: AdminMenu[]
  created_at: string
  updated_at: string
}

// 创建菜单的请求载荷
export interface CreateMenuPayload {
  parent_id: number
  type: MenuType
  code: string
  title: string
  path: string
  component: string
  icon: string
  sort: number
  status: MenuStatus
  remark: string
}

// 更新菜单的请求载荷（不含code字段）
export type UpdateMenuPayload = Omit<CreateMenuPayload, 'code'>

// 更新菜单状态的请求载荷
export interface UpdateMenuStatusPayload {
  status: MenuStatus
}
