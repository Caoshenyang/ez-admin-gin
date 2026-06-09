// ListQuery 列表查询的通用筛选条件（关键字 + 状态）。
export interface ListQuery {
  keyword?: string
  status?: number
}

// PageQuery 分页查询参数，在 ListQuery 基础上增加页码和每页条数。
export interface PageQuery extends ListQuery {
  page: number
  page_size: number
}

// PaginatedResponse 分页响应结构，包含当前页数据列表及分页元信息。
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}
