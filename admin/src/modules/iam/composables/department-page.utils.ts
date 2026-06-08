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
    leader_user_id: null,
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
    leader_user_id: department.leader_user_id ?? null,
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
    name: formModel.name.trim(),
    code: formModel.code.trim(),
    leader_user_id: formModel.leader_user_id,
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark.trim(),
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

export function collectDepartmentSubtreeIDs(items: DepartmentItem[], targetID: number) {
  const result = new Set<number>()

  function collect(nodes: DepartmentItem[]) {
    for (const node of nodes) {
      result.add(node.id)
      collect(node.children ?? [])
    }
  }

  function find(nodes: DepartmentItem[]): boolean {
    for (const node of nodes) {
      if (node.id === targetID) {
        collect([node])
        return true
      }

      if (find(node.children ?? [])) {
        return true
      }
    }

    return false
  }

  if (targetID !== 0) {
    find(items)
  }

  return result
}

export function collectDepartmentRowKeys(items: DepartmentItem[]): number[] {
  return flattenDepartments(items).map((item) => item.id)
}

export function sortDepartmentIDsForDelete(items: DepartmentItem[], ids: Array<string | number>) {
  const selectedIDs = new Set(ids.map(Number))
  const result: Array<{ id: number; depth: number }> = []

  function walk(nodes: DepartmentItem[], depth: number) {
    for (const node of nodes) {
      if (selectedIDs.has(node.id)) {
        result.push({ id: node.id, depth })
      }

      walk(node.children ?? [], depth + 1)
    }
  }

  walk(items, 0)
  return result.sort((a, b) => b.depth - a.depth).map((item) => item.id)
}

// 将部门树转换为树形选择控件选项
export function buildDepartmentTreeOptions(
  items: DepartmentItem[],
  excludedIDs: Set<number> = new Set(),
): TreeSelectOption[] {
  return items
    .filter((item) => !excludedIDs.has(item.id))
    .map((item) => ({
      key: item.id,
      label: `${item.name}（${item.code}）`,
      value: item.id,
      children: item.children?.length ? buildDepartmentTreeOptions(item.children, excludedIDs) : undefined,
    }))
}
