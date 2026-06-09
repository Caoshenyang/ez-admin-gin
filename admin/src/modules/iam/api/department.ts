import http from '@/api/http'

import type { ApiResponse } from '@/api/types'
import type {
  CreateDepartmentPayload,
  DepartmentItem,
  DepartmentListQuery,
  UpdateDepartmentPayload,
  UpdateDepartmentStatusPayload,
} from '../types/department'

// normalizeDepartmentItem 兼容旧数据：负责人未设置时统一使用 null，不能把 0 暴露给页面。
function normalizeDepartmentItem(item: DepartmentItem): DepartmentItem {
  return {
    ...item,
    leader_user_id: item.leader_user_id && item.leader_user_id > 0 ? item.leader_user_id : null,
    children: Array.isArray(item.children) ? item.children.map(normalizeDepartmentItem) : undefined,
  }
}

// getDepartments 查询部门列表，支持关键字和状态筛选。
export async function getDepartments(params: DepartmentListQuery = {}) {
  const response = await http.get<ApiResponse<DepartmentItem[]>>('/system/departments', { params })
  return (response.data.data ?? []).map(normalizeDepartmentItem)
}

// createDepartment 创建部门。
export async function createDepartment(payload: CreateDepartmentPayload) {
  const response = await http.post<ApiResponse<DepartmentItem>>('/system/departments', payload)
  return normalizeDepartmentItem(response.data.data)
}

// updateDepartment 更新部门信息。
export async function updateDepartment(id: number, payload: UpdateDepartmentPayload) {
  const response = await http.post<ApiResponse<DepartmentItem>>(
    `/system/departments/${id}/update`,
    payload,
  )
  return normalizeDepartmentItem(response.data.data)
}

// updateDepartmentStatus 切换部门的启用/禁用状态。
export async function updateDepartmentStatus(id: number, payload: UpdateDepartmentStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/departments/${id}/status`,
    payload,
  )
  return response.data.data
}

// deleteDepartment 删除部门。
export async function deleteDepartment(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(`/system/departments/${id}/delete`)
  return response.data.data
}
