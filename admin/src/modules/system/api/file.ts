import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type { FileItem, FileListQuery, FileListResponse } from '../types/file'

// 获取文件列表（分页查询）
export async function getFiles(params: FileListQuery) {
  const response = await http.get<ApiResponse<FileListResponse>>('/system/files', { params })
  return response.data.data
}

// 上传文件
export async function uploadFile(formData: FormData) {
  const response = await http.post<ApiResponse<FileItem>>('/system/files', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return response.data.data
}
