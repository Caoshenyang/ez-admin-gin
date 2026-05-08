import { DepartmentStatus, type CreateDepartmentPayload } from './department'

// 部门表单数据模型，继承自创建载荷并增加 id 字段
export interface DepartmentFormModel extends CreateDepartmentPayload {
  id: number
}

// 部门页面查询参数
export interface DepartmentPageQuery {
  keyword: string
  status: 0 | DepartmentStatus
}
