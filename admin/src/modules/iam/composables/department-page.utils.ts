import type { SelectOption, TreeSelectOption } from 'naive-ui'

import { DepartmentStatus, type DepartmentItem } from '@/modules/iam/types/department'

import type { DepartmentFormModel, DepartmentPageQuery } from '../types/department-page'

export const departmentStatusOptions: SelectOption[] = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: DepartmentStatus.Enabled },
  { label: '禁用', value: DepartmentStatus.Disabled },
] 

export const departmentFormStatusOptions: SelectOption[] = departmentStatusOptions.slice(1)

export function defaultDepartmentQuery(): DepartmentPageQuery {
  return {
    keyword: '',
    status: 0,
  }
}

// 生成部门表单的默认值
export function defaultDepartmentFormModel(): DepartmentFormModel {
  return {
    id: 0,
    parent_id: 0,
    name: '',
    code: '',
    leader_user_id: 0,
    sort: 0,
    status: DepartmentStatus.Enabled,
    remark: '',
  }
}

export function toDepartmentFormModel(department: DepartmentItem): DepartmentFormModel {
  return {
    id: department.id,
    parent_id: department.parent_id,
    name: department.name,
    code: department.code,
    leader_user_id: department.leader_user_id,
    sort: department.sort,
    status: department.status,
    remark: department.remark,
  }
}

export function normalizeDepartmentQuery(query: DepartmentPageQuery) {
  return {
    keyword: query.keyword.trim() || undefined,
    status: query.status === 0 ? undefined : query.status,
  }
}

export function buildDepartmentPayload(formModel: DepartmentFormModel) {
  return {
    parent_id: formModel.parent_id,
    name: formModel.name,
    code: formModel.code,
    leader_user_id: formModel.leader_user_id,
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark,
  }
}

export function flattenDepartments(items: DepartmentItem[]) {
  const result: DepartmentItem[] = []

  for (const item of items) {
    result.push(item)

    if (item.children?.length) {
      result.push(...flattenDepartments(item.children))
    }
  }

  return result
}

// 将部门树转换为树形选择控件选项
export function buildDepartmentTreeOptions(items: DepartmentItem[]): TreeSelectOption[] {
  return items.map((item) => ({
    label: `${item.name}（${item.code}）`,
    value: item.id,
    children: item.children?.length ? buildDepartmentTreeOptions(item.children) : undefined,
  }))
}
