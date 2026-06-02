import type { FormInst, FormRules, TreeSelectOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, nextTick, onMounted, reactive, ref } from 'vue'

import { buttonPermissionCodes } from '@/router/dynamic-menu'
import {
  createDepartment,
  getDepartments,
  updateDepartment,
  updateDepartmentStatus,
} from '../api/department'
import {
  DepartmentStatus,
  type DepartmentItem,
} from '../types/department'
import type { DepartmentFormModel, DepartmentPageQuery } from '../types/department-page'
import {
  buildDepartmentPayload,
  buildDepartmentTreeOptions,
  departmentFormStatusOptions,
  departmentStatusOptions,
  defaultDepartmentQuery,
  defaultDepartmentFormModel,
  normalizeDepartmentQuery,
  toDepartmentFormModel,
} from './department-page.utils'

export function useDepartmentPage() {
  const message = useMessage()
  const loading = ref(false)
  const saving = ref(false)
  const departments = ref<DepartmentItem[]>([])
  const formRef = ref<FormInst | null>(null)
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')

  const query = reactive<DepartmentPageQuery>(defaultDepartmentQuery())

  const formModel = reactive<DepartmentFormModel>(defaultDepartmentFormModel())

  const statusOptions = departmentStatusOptions
  const formStatusOptions = departmentFormStatusOptions

  // 上级部门树形选择选项，包含"作为根部门"选项
  const parentOptions = computed<TreeSelectOption[]>(() => {
    return [
      { key: 0, label: '作为根部门', value: 0 },
      ...buildDepartmentTreeOptions(departments.value),
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
    } finally {
      loading.value = false
    }
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

  onMounted(() => {
    void loadDepartments()
  })

  return {
    canUse,
    departments,
    formMode,
    formModel,
    formRef,
    formStatusOptions,
    formVisible,
    handleReset,
    handleSearch,
    handleSubmit,
    handleToggleStatus,
    loading,
    openCreate,
    openEdit,
    parentOptions,
    query,
    rules,
    saving,
    statusOptions,
  }
}
