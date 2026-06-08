import type { DataTableRowKey, TreeSelectOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { getDepartments } from '../api/department'
import { getPosts } from '../api/post'
import { getRoles } from '../api/role'
import {
  createUser,
  deleteUser,
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
  buildUserDepartmentFilterTreeOptions,
  buildUserDepartmentTreeOptions,
  buildUserUpdatePayload,
  countDepartments,
  defaultUserListQuery,
  defaultUserFormModel,
  filterDepartmentTreeOptions,
  toUserFormModel,
  normalizeUserListQuery,
  userFormRules,
  userStatusOptions,
} from './user-page.utils'

export function useUserPage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const loading = ref(false)
  const roleSaving = ref(false)
  const departments = ref<DepartmentItem[]>([])
  const users = ref<UserItem[]>([])
  const roles = ref<RoleItem[]>([])
  const posts = ref<PostItem[]>([])
  const departmentKeyword = ref('')
  const total = ref(0)
  const checkedRowKeys = ref<DataTableRowKey[]>([])
  const roleVisible = ref(false)
  const roleUser = ref<UserItem | null>(null)
  const selectedRoleIDs = ref<number[]>([])

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

  const roleNameMap = computed(() => new Map(roles.value.map((role) => [role.id, role.name])))

  const departmentNameMap = computed(() => buildDepartmentNameMap(departments.value))

  const postNameMap = computed(() => new Map(posts.value.map((post) => [post.id, post.name])))

  const departmentTreeOptions = computed<TreeSelectOption[]>(() => buildUserDepartmentTreeOptions(departments.value))

  const departmentFilterTreeOptions = computed<TreeSelectOption[]>(() => buildUserDepartmentFilterTreeOptions(departments.value))

  const filteredDepartmentTreeOptions = computed<TreeSelectOption[]>(() =>
    filterDepartmentTreeOptions(departmentFilterTreeOptions.value, departmentKeyword.value),
  )

  const departmentCount = computed(() => countDepartments(departments.value))

  const selectedDepartmentKeys = computed<Array<string | number>>(() => {
    if (!query.department_id) {
      return []
    }

    return [query.department_id]
  })

  const roleOptions = computed(() => buildRoleOptions(roles.value))

  const postOptions = computed(() => buildPostOptions(posts.value))

  const roleFilterOptions = computed(() => buildRoleFilterOptions(roles.value))

  const selectedCount = computed(() => checkedRowKeys.value.length)

  const displayUsers = computed(() => users.value)

  const displayTotal = computed(() => total.value)

  function handleCheckedRowKeys(keys: DataTableRowKey[]) {
    checkedRowKeys.value = keys
  }

  function handlePageChange(page: number) {
    query.page = page
    void loadUsers()
  }

  function handlePageSizeChange(pageSize: number) {
    query.page = 1
    query.page_size = pageSize
    void loadUsers()
  }

  function handleReset() {
    Object.assign(query, defaultUserListQuery())
    departmentKeyword.value = ''
    void loadUsers()
  }

  function handleSearch() {
    query.page = 1
    void loadUsers()
  }

  function handleSelectDepartment(keys: Array<string | number>) {
    query.department_id = Number(keys[0] ?? 0)
    query.page = 1
    void loadUsers()
  }

  function handleClearDepartment() {
    if (!query.department_id) {
      return
    }

    query.department_id = 0
    query.page = 1
    void loadUsers()
  }

  function openCreate() {
    openCreateBase()
  }

  function openEdit(row: UserItem) {
    openEditBase(toUserFormModel(row))
  }

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

  async function loadRoles() {
    const data = await getRoles({
      page: 1,
      page_size: 100,
      status: RoleStatus.Enabled,
    })
    roles.value = data.items
  }

  async function loadDepartments() {
    departments.value = await getDepartments()
  }

  async function loadPosts() {
    posts.value = await getPosts({ status: PostStatus.Enabled })
  }

  async function submitForm() {
    await formRef.value?.validate()
    saving.value = true
    try {
      if (formMode.value === 'create') {
        await createUser(buildUserCreatePayload(formModel))
        message.success('用户创建成功')
      } else {
        await updateUser(formModel.id, buildUserUpdatePayload(formModel))
        message.success('用户更新成功')
      }

      formVisible.value = false
      await loadUsers()
    } finally {
      saving.value = false
    }
  }

  async function handleToggleStatus(row: UserItem, status: UserStatus) {
    await updateUserStatus(row.id, { status })
    message.success('用户状态已更新')
    await loadUsers()
  }

  async function handleDelete(row: UserItem) {
    await deleteUser(row.id)
    message.success('用户已删除')
    await loadUsers()
  }

  function openRole(row: UserItem) {
    roleUser.value = row
    selectedRoleIDs.value = [...row.role_ids]
    roleVisible.value = true
  }

  async function handleSaveRoles() {
    if (!roleUser.value) {
      return
    }

    roleSaving.value = true
    try {
      await updateUserRoles(roleUser.value.id, { role_ids: selectedRoleIDs.value })
      message.success('用户角色已更新')
      roleVisible.value = false
      await loadUsers()
    } finally {
      roleSaving.value = false
    }
  }

  onMounted(async () => {
    await Promise.all([loadDepartments(), loadRoles(), loadPosts(), loadUsers()])
  })

  return {
    canUse,
    checkedRowKeys,
    departmentNameMap,
    departmentCount,
    departmentFilterTreeOptions,
    departmentKeyword,
    departmentTreeOptions,
    displayTotal,
    displayUsers,
    filteredDepartmentTreeOptions,
    formMode,
    formModel,
    formRef,
    formVisible,
    handleCheckedRowKeys,
    handlePageChange,
    handlePageSizeChange,
    handleClearDepartment,
    handleDelete,
    handleReset,
    handleSaveRoles,
    handleSearch,
    handleSelectDepartment,
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
    selectedDepartmentKeys,
    selectedRoleIDs,
    statusOptions: userStatusOptions,
    submitForm,
  }
}
