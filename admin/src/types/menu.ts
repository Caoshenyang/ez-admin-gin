export const MenuType = {
  Directory: 1,
  Menu: 2,
  Button: 3,
} as const

export type MenuType = (typeof MenuType)[keyof typeof MenuType]

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
