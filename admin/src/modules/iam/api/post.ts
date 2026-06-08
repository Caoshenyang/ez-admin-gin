import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  CreatePostPayload,
  PostItem,
  PostListQuery,
  UpdatePostPayload,
  UpdatePostStatusPayload,
} from '../types/post'

// getPosts 查询岗位列表，支持关键字和状态筛选。
export async function getPosts(params: PostListQuery) {
  const response = await http.get<ApiResponse<PostItem[]>>('/system/posts', { params })
  return response.data.data
}

// createPost 创建岗位。
export async function createPost(payload: CreatePostPayload) {
  const response = await http.post<ApiResponse<PostItem>>('/system/posts', payload)
  return response.data.data
}

// updatePost 更新岗位信息。
export async function updatePost(id: number, payload: UpdatePostPayload) {
  const response = await http.post<ApiResponse<PostItem>>(`/system/posts/${id}/update`, payload)
  return response.data.data
}

// updatePostStatus 切换岗位的启用/禁用状态。
export async function updatePostStatus(id: number, payload: UpdatePostStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/posts/${id}/status`,
    payload,
  )
  return response.data.data
}

// deletePost 删除指定岗位。
export async function deletePost(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(`/system/posts/${id}/delete`)
  return response.data.data
}
