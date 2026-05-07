export const MenuType = {
  Directory: 1,
  Menu: 2,
  Button: 3,
} as const

export type MenuType = (typeof MenuType)[keyof typeof MenuType]

export const MenuStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

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

export type UpdateMenuPayload = Omit<CreateMenuPayload, 'code'>

export interface UpdateMenuStatusPayload {
  status: MenuStatus
}
