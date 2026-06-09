import { MenuStatus, MenuType } from './menu'

// 菜单表单数据模型
export interface MenuFormModel {
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
}

// 菜单搜索查询参数
export interface MenuQuery {
  keyword: string
  type: 0 | MenuType
  status: 0 | MenuStatus
}
