import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  CreateDepartmentPayload,
  DepartmentItem,
  DepartmentListQuery,
  UpdateDepartmentPayload,
  UpdateDepartmentStatusPayload,
} from '../types/department'

export async function getDepartments(params: DepartmentListQuery = {}) {
  const response = await http.get<ApiResponse<DepartmentItem[]>>('/system/departments', { params })
  return response.data.data ?? []
}

export async function createDepartment(payload: CreateDepartmentPayload) {
  const response = await http.post<ApiResponse<DepartmentItem>>('/system/departments', payload)
  return response.data.data
}

export async function updateDepartment(id: number, payload: UpdateDepartmentPayload) {
  const response = await http.post<ApiResponse<DepartmentItem>>(`/system/departments/${id}/update`, payload)
  return response.data.data
}

export async function updateDepartmentStatus(id: number, payload: UpdateDepartmentStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/departments/${id}/status`,
    payload,
  )
  return response.data.data
}
