import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  CreateDepartmentPayload,
  DepartmentItem,
  DepartmentListQuery,
  UpdateDepartmentPayload,
  UpdateDepartmentStatusPayload,
} from '../types/department'

// getDepartments 查询部门列表，支持关键字和状态筛选。
export async function getDepartments(params: DepartmentListQuery = {}) {
  const response = await http.get<ApiResponse<DepartmentItem[]>>('/system/departments', { params })
  return response.data.data ?? []
}

// createDepartment 创建部门。
export async function createDepartment(payload: CreateDepartmentPayload) {
  const response = await http.post<ApiResponse<DepartmentItem>>('/system/departments', payload)
  return response.data.data
}

// updateDepartment 更新部门信息。
export async function updateDepartment(id: number, payload: UpdateDepartmentPayload) {
  const response = await http.post<ApiResponse<DepartmentItem>>(`/system/departments/${id}/update`, payload)
  return response.data.data
}

// updateDepartmentStatus 切换部门的启用/禁用状态。
export async function updateDepartmentStatus(id: number, payload: UpdateDepartmentStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/departments/${id}/status`,
    payload,
  )
  return response.data.data
}
