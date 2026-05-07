// ApiResponse 对应后端统一响应结构。
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}
