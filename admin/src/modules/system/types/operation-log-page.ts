import type { OperationLogListQuery } from './operation-log'

// 操作日志页面查询参数，扩展自列表查询参数并增加页面级筛选字段
export interface OperationLogPageQuery extends OperationLogListQuery {
  username: string
  method: string
  path: string
  success: string
}
