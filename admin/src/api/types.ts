// ApiResponse 对应后端统一响应结构，包含业务状态码、提示信息和实际数据。
export interface ApiResponse<T> {
  // code 业务状态码，0 通常表示成功。
  code: number
  // message 提示信息，用于向用户展示操作结果。
  message: string
  // data 实际返回的业务数据，泛型 T 由调用方指定。
  data: T
}
