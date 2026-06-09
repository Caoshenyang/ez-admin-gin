// 岗位状态常量：启用、禁用
export const PostStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

// 岗位状态联合类型
export type PostStatus = (typeof PostStatus)[keyof typeof PostStatus]

// 岗位列表项数据结构
export interface PostItem {
  id: number
  code: string
  name: string
  sort: number
  status: PostStatus
  remark: string
  created_at: string
  updated_at: string
}

// 岗位列表查询参数
export interface PostListQuery {
  keyword?: string
  status?: PostStatus | 0
}

// 创建岗位的请求载荷
export interface CreatePostPayload {
  code: string
  name: string
  sort: number
  status: PostStatus
  remark: string
}

// 更新岗位的请求载荷（与创建载荷相同）
export type UpdatePostPayload = CreatePostPayload

// 更新岗位状态的请求载荷
export interface UpdatePostStatusPayload {
  status: PostStatus
}
