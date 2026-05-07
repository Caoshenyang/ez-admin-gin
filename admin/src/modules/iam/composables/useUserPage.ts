import type {
  DataTableRowKey,
  FormRules,
  SelectOption,
  TreeSelectOption,
} from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
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

export interface UserFormModel {
  id: number
  username: string
  password: string
  nickname: string
  department_id: number
  status: UserStatus
  role_ids: number[]
  post_ids: number[]
}

const userStatusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: UserStatus.Enabled },
  { label: '禁用', value: UserStatus.Disabled },
]

const userFormRules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  department_id: [{ required: true, type: 'number', message: '请选择部门', trigger: ['change', 'blur'] }],
}

function defaultUserFormModel(): UserFormModel {
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

export function useUserPage() {
  const { canUse } = usePermission()
  const loading = ref(false)
  const roleSaving = ref(false)
  const departments = ref<DepartmentItem[]>([])
  const users = ref<UserItem[]>([])
  const roles = ref<RoleItem[]>([])
  const posts = ref<PostItem[]>([])
  const total = ref(0)
  const checkedRowKeys = ref<DataTableRowKey[]>([])
  const successText = ref('')
  const roleVisible = ref(false)
  const roleUser = ref<UserItem | null>(null)
  const selectedRoleIDs = ref<number[]>([])

  const query = reactive<UserListQuery>({
    page: 1,
    page_size: 10,
    keyword: '',
    role_id: 0,
    status: 0,
  })

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
  const departmentNameMap = computed(() => {
    return new Map(flattenDepartments(departments.value).map((department) => [department.id, department.name]))
  })
  const postNameMap = computed(() => new Map(posts.value.map((post) => [post.id, post.name])))

  const departmentTreeOptions = computed<TreeSelectOption[]>(() => {
    return [{ label: '未分配部门', value: 0 }, ...buildDepartmentTreeOptions(departments.value)]
  })

  const roleOptions = computed<SelectOption[]>(() => {
    return roles.value.map((role) => ({
      label: `${role.name}（${role.code}）`,
      value: role.id,
    }))
  })

  const postOptions = computed<SelectOption[]>(() => {
    return posts.value.map((post) => ({
      label: `${post.name}（${post.code}）`,
      value: post.id,
    }))
  })

  const roleFilterOptions = computed<SelectOption[]>(() => {
    return [
      { label: '角色：全部', value: 0 },
      ...roles.value.map((role) => ({
        label: role.name,
        value: role.id,
      })),
    ]
  })

  const selectedCount = computed(() => checkedRowKeys.value.length)

  const displayUsers = computed(() => {
    if (!query.role_id) {
      return users.value
    }

    return users.value.filter((user) => user.role_ids.includes(query.role_id ?? 0))
  })

  const displayTotal = computed(() => (query.role_id ? displayUsers.value.length : total.value))

  function closeSuccess() {
    successText.value = ''
  }

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
    query.page = 1
    query.page_size = 10
    query.keyword = ''
    query.role_id = 0
    query.status = 0
    void loadUsers()
  }

  function handleSearch() {
    query.page = 1
    void loadUsers()
  }

  function openCreate() {
    openCreateBase()
  }

  function openEdit(row: UserItem) {
    openEditBase({
      id: row.id,
      username: row.username,
      password: '',
      nickname: row.nickname,
      department_id: row.department_id,
      status: row.status,
      role_ids: row.role_ids,
      post_ids: row.post_ids,
    })
  }

  async function loadUsers() {
    loading.value = true
    try {
      const data = await getUsers({
        ...query,
        keyword: query.keyword?.trim() || undefined,
        role_id: query.role_id === 0 ? undefined : query.role_id,
        status: query.status === 0 ? undefined : query.status,
      })
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
        await createUser({
          username: formModel.username,
          password: formModel.password,
          nickname: formModel.nickname,
          department_id: formModel.department_id,
          status: formModel.status,
          role_ids: formModel.role_ids,
          post_ids: formModel.post_ids,
        })
        successText.value = '用户创建成功，临时密码已生成'
      } else {
        await updateUser(formModel.id, {
          nickname: formModel.nickname,
          department_id: formModel.department_id,
          status: formModel.status,
          post_ids: formModel.post_ids,
        })
        successText.value = '用户信息已更新'
      }

      formVisible.value = false
      await loadUsers()
    } finally {
      saving.value = false
    }
  }

  async function handleToggleStatus(row: UserItem, status: UserStatus) {
    await updateUserStatus(row.id, { status })
    successText.value = `用户已${status === UserStatus.Enabled ? '启用' : '禁用'}`
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
      successText.value = '用户角色已更新'
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

function flattenDepartments(items: DepartmentItem[]) {
  const result: DepartmentItem[] = []

  for (const item of items) {
    result.push(item)

    if (item.children?.length) {
      result.push(...flattenDepartments(item.children))
    }
  }

  return result
}

function buildDepartmentTreeOptions(items: DepartmentItem[]): TreeSelectOption[] {
  return items.map((item) => ({
    label: `${item.name}（${item.code}）`,
    value: item.id,
    children: item.children?.length ? buildDepartmentTreeOptions(item.children) : undefined,
  }))
}
