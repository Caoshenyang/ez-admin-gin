export interface ListQuery {
  keyword?: string
  status?: number
}

export interface PageQuery extends ListQuery {
  page: number
  page_size: number
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}
