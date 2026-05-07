import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  CreatePostPayload,
  PostItem,
  PostListQuery,
  UpdatePostPayload,
  UpdatePostStatusPayload,
} from '../types/post'

export async function getPosts(params: PostListQuery) {
  const response = await http.get<ApiResponse<PostItem[]>>('/system/posts', { params })
  return response.data.data
}

export async function createPost(payload: CreatePostPayload) {
  const response = await http.post<ApiResponse<PostItem>>('/system/posts', payload)
  return response.data.data
}

export async function updatePost(id: number, payload: UpdatePostPayload) {
  const response = await http.post<ApiResponse<PostItem>>(`/system/posts/${id}/update`, payload)
  return response.data.data
}

export async function updatePostStatus(id: number, payload: UpdatePostStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/posts/${id}/status`,
    payload,
  )
  return response.data.data
}
