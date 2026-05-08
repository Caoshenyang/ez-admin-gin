import type { DataTableRowKey, TreeSelectOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useSuccessFeedback } from '@/composables/useSuccessFeedback'
import { getDepartments } from '../api/department'
import { getPosts } from '../api/post'
import { getRoles } from '../api/role'
import {
  createUser,
  getUsers,
  updateUser,
  updateUserRoles,
  updateUserStatus,
} from '../api/user'
import type { DepartmentItem } from '../types/department'
import { PostStatus, type PostItem } from '../types/post'
import { RoleStatus, type RoleItem } from '../types/role'
import { UserStatus, type UserItem, type UserListQuery } from '../types/user'
import type { UserFormModel } from '../types/user-page'
import {
  buildDepartmentNameMap,
  buildPostOptions,
  buildRoleFilterOptions,
  buildRoleOptions,
  buildUserCreatePayload,
  buildUserDepartmentTreeOptions,
  buildUserUpdatePayload,
  defaultUserListQuery,
  defaultUserFormModel,
  toUserFormModel,
  normalizeUserListQuery,
  userFormRules,
  userStatusOptions,
} from './user-page.utils'

// 用户管理页面组合式函数，封装用户列表、创建、编辑、角色分配、状态切换等逻辑
export function useUserPage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const { closeSuccess, showSuccess, successText } = useSuccessFeedback()
  const loading = ref(false)
  const roleSaving = ref(false)
  const departments = ref<DepartmentItem[]>([])
  const users = ref<UserItem[]>([])
  const roles = ref<RoleItem[]>([])
  const posts = ref<PostItem[]>([])
  const total = ref(0)
  const checkedRowKeys = ref<DataTableRowKey[]>([])
  const roleVisible = ref(false)
  const roleUser = ref<UserItem | null>(null)
  const selectedRoleIDs = ref<number[]>([])

  // 用户列表查询条件
  const query = reactive<UserListQuery>(defaultUserListQuery())

  const {
    formRef,
    formVisible,
    formMode,
    formModel,
    saving,
    rules,
    openCreate: openCreateBase,
    openEdit: openEditBase,
  } = useModalForm<UserFormModel>(defaultUserFormModel, { rules: userFormRules })

  // 角色ID到名称的映射
  const roleNameMap = computed(() => new Map(roles.value.map((role) => [role.id, role.name])))

  // 部门ID到名称的映射（扁平化后的部门列表）
  const departmentNameMap = computed(() => buildDepartmentNameMap(departments.value))

  // 岗位ID到名称的映射
  const postNameMap = computed(() => new Map(posts.value.map((post) => [post.id, post.name])))

  // 部门树形选择选项
  const departmentTreeOptions = computed<TreeSelectOption[]>(() => buildUserDepartmentTreeOptions(departments.value))

  // 角色下拉选项
  const roleOptions = computed(() => buildRoleOptions(roles.value))

  // 岗位下拉选项
  const postOptions = computed(() => buildPostOptions(posts.value))

  // 角色筛选下拉选项（含"全部"）
  const roleFilterOptions = computed(() => buildRoleFilterOptions(roles.value))

  // 已勾选的行数
  const selectedCount = computed(() => checkedRowKeys.value.length)

  // 根据角色筛选条件过滤后的用户列表
  const displayUsers = computed(() => {
    if (!query.role_id) {
      return users.value
    }

    return users.value.filter((user) => user.role_ids.includes(query.role_id ?? 0))
  })

  // 过滤后的用户总数
  const displayTotal = computed(() => (query.role_id ? displayUsers.value.length : total.value))

  // 处理表格勾选行变化
  function handleCheckedRowKeys(keys: DataTableRowKey[]) {
    checkedRowKeys.value = keys
  }

  // 处理分页页码变化
  function handlePageChange(page: number) {
    query.page = page
    void loadUsers()
  }

  // 处理每页条数变化
  function handlePageSizeChange(pageSize: number) {
    query.page = 1
    query.page_size = pageSize
    void loadUsers()
  }

  // 重置搜索条件并重新加载用户列表
  function handleReset() {
    Object.assign(query, defaultUserListQuery())
    void loadUsers()
  }

  // 搜索用户
  function handleSearch() {
    query.page = 1
    void loadUsers()
  }

  // 打开创建用户的弹窗
  function openCreate() {
    openCreateBase()
  }

  // 打开编辑用户的弹窗，将当前行数据填充到表单
  function openEdit(row: UserItem) {
    openEditBase(toUserFormModel(row))
  }

  // 从服务端加载用户列表
  async function loadUsers() {
    loading.value = true
    try {
      const data = await getUsers(normalizeUserListQuery(query))
      users.value = data.items
      total.value = data.total
      checkedRowKeys.value = []
    } finally {
      loading.value = false
    }
  }

  // 从服务端加载角色列表
  async function loadRoles() {
    const data = await getRoles({
      page: 1,
      page_size: 100,
      status: RoleStatus.Enabled,
    })
    roles.value = data.items
  }

  // 从服务端加载部门树数据
  async function loadDepartments() {
    departments.value = await getDepartments()
  }

  // 从服务端加载岗位列表
  async function loadPosts() {
    posts.value = await getPosts({ status: PostStatus.Enabled })
  }

  // 提交用户表单（新建或更新）
  async function submitForm() {
    await formRef.value?.validate()
    saving.value = true
    try {
      if (formMode.value === 'create') {
        await createUser(buildUserCreatePayload(formModel))
        showSuccess('用户创建成功，临时密码已生成')
        message.success('用户创建成功')
      } else {
        await updateUser(formModel.id, buildUserUpdatePayload(formModel))
        showSuccess('用户信息已更新')
        message.success('用户更新成功')
      }

      formVisible.value = false
      await loadUsers()
    } finally {
      saving.value = false
    }
  }

  // 切换用户的启用/禁用状态
  async function handleToggleStatus(row: UserItem, status: UserStatus) {
    await updateUserStatus(row.id, { status })
    showSuccess(`用户已${status === UserStatus.Enabled ? '启用' : '禁用'}`)
    message.success('用户状态已更新')
    await loadUsers()
  }

  // 打开角色分配弹窗
  function openRole(row: UserItem) {
    roleUser.value = row
    selectedRoleIDs.value = [...row.role_ids]
    roleVisible.value = true
  }

  // 保存用户角色分配
  async function handleSaveRoles() {
    if (!roleUser.value) {
      return
    }

    roleSaving.value = true
    try {
      await updateUserRoles(roleUser.value.id, { role_ids: selectedRoleIDs.value })
      showSuccess('用户角色已更新')
      message.success('用户角色已更新')
      roleVisible.value = false
      await loadUsers()
    } finally {
      roleSaving.value = false
    }
  }

  // 组件挂载时并行加载部门、角色、岗位和用户数据
  onMounted(async () => {
    await Promise.all([loadDepartments(), loadRoles(), loadPosts(), loadUsers()])
  })

  return {
    canUse,
    checkedRowKeys,
    closeSuccess,
    departmentNameMap,
    departmentTreeOptions,
    displayTotal,
    displayUsers,
    formMode,
    formModel,
    formRef,
    formVisible,
    handleCheckedRowKeys,
    handlePageChange,
    handlePageSizeChange,
    handleReset,
    handleSaveRoles,
    handleSearch,
    handleToggleStatus,
    loading,
    openCreate,
    openEdit,
    openRole,
    postNameMap,
    postOptions,
    query,
    roleFilterOptions,
    roleNameMap,
    roleOptions,
    roleSaving,
    roleUser,
    roleVisible,
    rules,
    saving,
    selectedCount,
    selectedRoleIDs,
    statusOptions: userStatusOptions,
    submitForm,
    successText,
  }
}
