import http from './http'

import type { ApiResponse } from '../types/http'
import type { PostItem, PostListQuery } from '../types/post'

export async function getPosts(params: PostListQuery) {
  const response = await http.get<ApiResponse<PostItem[]>>('/system/posts', { params })
  return response.data.data
}
