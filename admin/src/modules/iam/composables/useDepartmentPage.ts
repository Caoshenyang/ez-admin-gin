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
  // DepartmentItem 类型定义。
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

// 部门管理页面组合式函数，封装部门树加载、创建、编辑、状态切换等逻辑
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
      { label: '作为根部门', value: 0 },
      ...buildDepartmentTreeOptions(departments.value),
    ]
  })

  // 表单校验规则
  const rules: FormRules = {
    name: [{ required: true, message: '请输入部门名称', trigger: ['blur', 'input'] }],
    code: [{ required: true, message: '请输入部门编码', trigger: ['blur', 'input'] }],
  }

  // 判断当前用户是否拥有指定按钮权限码
  function canUse(code: string) {
    return buttonPermissionCodes.value.includes(code)
  }

  // 从服务端加载部门树数据
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

  // 点击搜索按钮，重新加载部门列表
  function handleSearch() {
    void loadDepartments()
  }

  // 重置搜索条件并重新加载部门列表
  function handleReset() {
    Object.assign(query, defaultDepartmentQuery())
    void loadDepartments()
  }

  // 打开新建部门的弹窗
  function openCreate() {
    formMode.value = 'create'
    Object.assign(formModel, defaultDepartmentFormModel())
    formVisible.value = true
  }

  // 打开编辑部门的弹窗，将当前行数据填充到表单
  function openEdit(row: DepartmentItem) {
    formMode.value = 'edit'
    Object.assign(formModel, toDepartmentFormModel(row))
    formVisible.value = true
  }

  // 提交部门表单（新建或更新）
  async function handleSubmit() {
    await formRef.value?.validate()
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

  // 切换部门的启用/禁用状态
  async function handleToggleStatus(row: DepartmentItem, status: DepartmentStatus) {
    try {
      await updateDepartmentStatus(row.id, { status })
      message.success(status === DepartmentStatus.Enabled ? '部门已启用' : '部门已禁用')
      await loadDepartments()
    } catch {
      message.error('部门状态更新失败')
    }
  }

  // 组件挂载时自动加载部门列表
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
