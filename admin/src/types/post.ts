export const PostStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type PostStatus = (typeof PostStatus)[keyof typeof PostStatus]

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

export interface PostListQuery {
  keyword?: string
  status?: PostStatus | 0
}
