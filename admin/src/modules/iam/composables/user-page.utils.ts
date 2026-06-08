import type { FormRules, SelectOption, TreeSelectOption } from 'naive-ui'

import type { DepartmentItem } from '@/modules/iam/types/department'
import type { PostItem } from '@/modules/iam/types/post'
import type { RoleItem } from '@/modules/iam/types/role'
import { UserStatus, type UserItem, type UserListQuery } from '@/modules/iam/types/user'

import { buildDepartmentTreeOptions, flattenDepartments } from './department-page.utils'
import type { UserFormModel } from '../types/user-page'

export const userStatusOptions: SelectOption[] = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: UserStatus.Enabled },
  { label: '禁用', value: UserStatus.Disabled },
]

export const userFormRules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  department_id: [
    { required: true, type: 'number', message: '请选择部门', trigger: ['change', 'blur'] },
  ],
}

export function defaultUserFormModel(): UserFormModel {
  return {
    id: 0,
    username: '',
    password: '',
    nickname: '',
    department_id: 0,
    status: UserStatus.Enabled,
    role_ids: [],
    post_ids: [],
  }
}

export function defaultUserListQuery(): UserListQuery {
  // 列表默认查询集中在这里，搜索重置和首屏加载共用同一份基线。
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    role_id: 0,
    status: 0,
    department_id: 0,
  }
}

export function normalizeUserListQuery(query: UserListQuery) {
  return {
    ...query,
    keyword: query.keyword?.trim() || undefined,
    role_id: query.role_id === 0 ? undefined : query.role_id,
    status: query.status === 0 ? undefined : query.status,
    department_id: query.department_id === 0 ? undefined : query.department_id,
  }
}

export function toUserFormModel(user: UserItem): UserFormModel {
  return {
    id: user.id,
    username: user.username,
    password: '',
    nickname: user.nickname,
    department_id: user.department_id,
    status: user.status,
    // 表单编辑时复制数组，避免勾选变更直接回写到表格行数据。
    role_ids: [...user.role_ids],
    post_ids: [...user.post_ids],
  }
}

export function buildUserCreatePayload(formModel: UserFormModel) {
  return {
    username: formModel.username,
    password: formModel.password,
    nickname: formModel.nickname,
    department_id: formModel.department_id,
    status: formModel.status,
    role_ids: formModel.role_ids,
    post_ids: formModel.post_ids,
  }
}

export function buildUserUpdatePayload(formModel: UserFormModel) {
  return {
    nickname: formModel.nickname,
    department_id: formModel.department_id,
    status: formModel.status,
    post_ids: formModel.post_ids,
  }
}

export function buildUserDepartmentTreeOptions(departments: DepartmentItem[]): TreeSelectOption[] {
  return [{ label: '未分配部门', value: 0 }, ...buildDepartmentTreeOptions(departments)]
}

export function buildUserDepartmentFilterTreeOptions(
  departments: DepartmentItem[],
): TreeSelectOption[] {
  return buildDepartmentTreeOptions(departments)
}

export function filterDepartmentTreeOptions(
  options: TreeSelectOption[],
  keyword: string,
): TreeSelectOption[] {
  const normalizedKeyword = keyword.trim().toLowerCase()
  if (!normalizedKeyword) {
    return options
  }

  return options.flatMap((option) => {
    const children = Array.isArray(option.children)
      ? filterDepartmentTreeOptions(option.children as TreeSelectOption[], normalizedKeyword)
      : []
    const label = String(option.label ?? '').toLowerCase()
    const value = String(option.value ?? '').toLowerCase()

    if (
      label.includes(normalizedKeyword) ||
      value.includes(normalizedKeyword) ||
      children.length > 0
    ) {
      return [{ ...option, children: children.length > 0 ? children : undefined }]
    }

    return []
  })
}

export function buildRoleOptions(roles: RoleItem[]): SelectOption[] {
  return roles.map((role) => ({
    label: `${role.name}（${role.code}）`,
    value: role.id,
  }))
}

export function buildPostOptions(posts: PostItem[]): SelectOption[] {
  return posts.map((post) => ({
    label: `${post.name}（${post.code}）`,
    value: post.id,
  }))
}

export function buildRoleFilterOptions(roles: RoleItem[]): SelectOption[] {
  return [
    { label: '角色：全部', value: 0 },
    ...roles.map((role) => ({
      label: role.name,
      value: role.id,
    })),
  ]
}

export function buildDepartmentNameMap(departments: DepartmentItem[]) {
  // 用户列表和详情都需要按部门 ID 反查名称，这里统一复用同一份扁平索引。
  return new Map(
    flattenDepartments(departments).map((department) => [department.id, department.name]),
  )
}

export function countDepartments(departments: DepartmentItem[]) {
  return flattenDepartments(departments).length
}
