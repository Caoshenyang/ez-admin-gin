import type { FormInst, FormRules, TreeSelectOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import { buttonPermissionCodes } from '@/router/dynamic-menu'
import {
  createDepartment,
  getDepartments,
  updateDepartment,
  updateDepartmentStatus,
} from '../api/department'
import {
  DepartmentStatus,
  type CreateDepartmentPayload,
  type DepartmentItem,
  type DepartmentListQuery,
} from '../types/department'

export interface DepartmentFormModel extends CreateDepartmentPayload {
  id: number
}

function resetFormModel(): DepartmentFormModel {
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

function buildDepartmentTreeOptions(items: DepartmentItem[]): TreeSelectOption[] {
  return items.map((item) => ({
    label: `${item.name}（${item.code}）`,
    value: item.id,
    children: item.children?.length ? buildDepartmentTreeOptions(item.children) : undefined,
  }))
}

export function useDepartmentPage() {
  const message = useMessage()
  const loading = ref(false)
  const saving = ref(false)
  const departments = ref<DepartmentItem[]>([])
  const formRef = ref<FormInst | null>(null)
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')

  const query = reactive<DepartmentListQuery>({
    keyword: '',
    status: 0,
  })

  const formModel = reactive<DepartmentFormModel>(resetFormModel())

  const statusOptions = [
    { label: '状态：全部', value: 0 },
    { label: '启用', value: DepartmentStatus.Enabled },
    { label: '禁用', value: DepartmentStatus.Disabled },
  ]

  const formStatusOptions = statusOptions.slice(1)

  const parentOptions = computed<TreeSelectOption[]>(() => {
    return [
      { label: '作为根部门', value: 0 },
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
        keyword: query.keyword?.trim() || undefined,
        status: query.status === 0 ? undefined : query.status,
      })
    } finally {
      loading.value = false
    }
  }

  function handleSearch() {
    void loadDepartments()
  }

  function handleReset() {
    query.keyword = ''
    query.status = 0
    void loadDepartments()
  }

  function openCreate() {
    formMode.value = 'create'
    Object.assign(formModel, resetFormModel())
    formVisible.value = true
  }

  function openEdit(row: DepartmentItem) {
    formMode.value = 'edit'
    Object.assign(formModel, {
      id: row.id,
      parent_id: row.parent_id,
      name: row.name,
      code: row.code,
      leader_user_id: row.leader_user_id,
      sort: row.sort,
      status: row.status,
      remark: row.remark,
    })
    formVisible.value = true
  }

  async function handleSubmit() {
    await formRef.value?.validate()
    saving.value = true

    try {
      if (formMode.value === 'create') {
        await createDepartment({
          parent_id: formModel.parent_id,
          name: formModel.name,
          code: formModel.code,
          leader_user_id: formModel.leader_user_id,
          sort: formModel.sort,
          status: formModel.status,
          remark: formModel.remark,
        })
        message.success('部门创建成功')
      } else {
        await updateDepartment(formModel.id, {
          parent_id: formModel.parent_id,
          name: formModel.name,
          code: formModel.code,
          leader_user_id: formModel.leader_user_id,
          sort: formModel.sort,
          status: formModel.status,
          remark: formModel.remark,
        })
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
