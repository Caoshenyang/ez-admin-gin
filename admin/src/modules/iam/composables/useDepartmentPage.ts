import type { DataTableRowKey, FormInst, FormRules, TreeSelectOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, nextTick, onMounted, reactive, ref } from 'vue'

import { buttonPermissionCodes } from '@/router/dynamic-menu'
import {
  createDepartment,
  deleteDepartment,
  getDepartments,
  updateDepartment,
  updateDepartmentStatus,
} from '../api/department'
import { getUsers } from '../api/user'
import {
  DepartmentStatus,
  type DepartmentItem,
} from '../types/department'
import { UserStatus, type UserItem } from '../types/user'
import type { DepartmentFormModel, DepartmentPageQuery } from '../types/department-page'
import {
  buildDepartmentPayload,
  buildDepartmentTreeOptions,
  collectDepartmentRowKeys,
  collectDepartmentSubtreeIDs,
  departmentFormStatusOptions,
  departmentStatusOptions,
  defaultDepartmentQuery,
  defaultDepartmentFormModel,
  normalizeDepartmentQuery,
  sortDepartmentIDsForDelete,
  toDepartmentFormModel,
} from './department-page.utils'

export function useDepartmentPage() {
  const message = useMessage()
  const loading = ref(false)
  const saving = ref(false)
  const departments = ref<DepartmentItem[]>([])
  const users = ref<UserItem[]>([])
  const checkedRowKeys = ref<DataTableRowKey[]>([])
  const expandedRowKeys = ref<DataTableRowKey[]>([])
  const formRef = ref<FormInst | null>(null)
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')

  const query = reactive<DepartmentPageQuery>(defaultDepartmentQuery())

  const formModel = reactive<DepartmentFormModel>(defaultDepartmentFormModel())

  const statusOptions = departmentStatusOptions
  const formStatusOptions = departmentFormStatusOptions

  const flatDepartments = computed(() => collectDepartmentRowKeys(departments.value))

  const selectedCount = computed(() => checkedRowKeys.value.length)

  const leaderNameMap = computed(() => {
    return new Map(users.value.map((user) => [user.id, user.nickname || user.username]))
  })

  // 上级部门树形选择选项，包含"作为根部门"选项
  const parentOptions = computed<TreeSelectOption[]>(() => {
    const excludedIDs = formMode.value === 'edit'
      ? collectDepartmentSubtreeIDs(departments.value, formModel.id)
      : new Set<number>()

    return [
      { key: 0, label: '作为根部门', value: 0 },
      ...buildDepartmentTreeOptions(departments.value, excludedIDs),
    ]
  })

  const rules: FormRules = {
    name: [{ required: true, message: '请输入部门名称', trigger: ['blur', 'input'] }],
    code: [{ required: true, message: '请输入部门编码', trigger: ['blur', 'input'] }],
  }

  function canUse(code: string) {
    return buttonPermissionCodes.value.includes(code)
  }

  async function loadDepartments() {
    loading.value = true
    try {
      departments.value = await getDepartments({
        ...normalizeDepartmentQuery(query),
      })
      checkedRowKeys.value = []
    } finally {
      loading.value = false
    }
  }

  async function loadUsers() {
    const data = await getUsers({
      page: 1,
      page_size: 100,
      status: UserStatus.Enabled,
    })
    users.value = data.items
  }

  function handleSearch() {
    void loadDepartments()
  }

  function handleReset() {
    Object.assign(query, defaultDepartmentQuery())
    void loadDepartments()
  }

  function openCreate() {
    formMode.value = 'create'
    Object.assign(formModel, defaultDepartmentFormModel())
    formVisible.value = true
    void nextTick(() => {
      formRef.value?.restoreValidation()
    })
  }

  function openEdit(row: DepartmentItem) {
    formMode.value = 'edit'
    Object.assign(formModel, toDepartmentFormModel(row))
    formVisible.value = true
    void nextTick(() => {
      formRef.value?.restoreValidation()
    })
  }

  function openCreateChild(row: DepartmentItem) {
    formMode.value = 'create'
    Object.assign(formModel, defaultDepartmentFormModel(), {
      parent_id: row.id,
      sort: 0,
    })
    formVisible.value = true
    void nextTick(() => {
      formRef.value?.restoreValidation()
    })
  }

  function handleCheckedRowKeys(keys: DataTableRowKey[]) {
    checkedRowKeys.value = keys
  }

  function handleExpandedRowKeys(keys: DataTableRowKey[]) {
    expandedRowKeys.value = keys
  }

  function expandAll() {
    expandedRowKeys.value = flatDepartments.value
  }

  function collapseAll() {
    expandedRowKeys.value = []
  }

  async function handleSubmit() {
    try {
      await formRef.value?.validate()
    } catch {
      return
    }

    saving.value = true

    try {
      if (formMode.value === 'create') {
        await createDepartment(buildDepartmentPayload(formModel))
        message.success('部门创建成功')
      } else {
        await updateDepartment(formModel.id, buildDepartmentPayload(formModel))
        message.success('部门更新成功')
      }

      formVisible.value = false
      await loadDepartments()
      await loadUsers()
    } catch {
      message.error(formMode.value === 'create' ? '部门创建失败' : '部门更新失败')
    } finally {
      saving.value = false
    }
  }

  async function handleToggleStatus(row: DepartmentItem, status: DepartmentStatus) {
    try {
      await updateDepartmentStatus(row.id, { status })
      message.success(status === DepartmentStatus.Enabled ? '部门已启用' : '部门已禁用')
      await loadDepartments()
    } catch {
      message.error('部门状态更新失败')
    }
  }

  async function handleDeleteSelected() {
    const deleteIDs = sortDepartmentIDsForDelete(departments.value, checkedRowKeys.value)
    if (deleteIDs.length === 0) {
      return
    }

    loading.value = true
    try {
      for (const id of deleteIDs) {
        await deleteDepartment(id)
      }

      message.success(deleteIDs.length === 1 ? '部门已删除' : `已删除 ${deleteIDs.length} 个部门`)
      checkedRowKeys.value = []
      await loadDepartments()
    } catch {
      message.error('部门删除失败，请确认没有子部门、关联用户或角色数据范围')
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    void Promise.all([loadDepartments(), loadUsers()])
  })

  return {
    canUse,
    checkedRowKeys,
    collapseAll,
    departments,
    expandAll,
    expandedRowKeys,
    formMode,
    formModel,
    formRef,
    formStatusOptions,
    formVisible,
    handleCheckedRowKeys,
    handleDeleteSelected,
    handleExpandedRowKeys,
    handleReset,
    handleSearch,
    handleSubmit,
    handleToggleStatus,
    leaderNameMap,
    loading,
    openCreate,
    openCreateChild,
    openEdit,
    parentOptions,
    query,
    rules,
    saving,
    selectedCount,
    statusOptions,
  }
}
