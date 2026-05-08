import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import { useSuccessFeedback } from '@/composables/useSuccessFeedback'
import { buttonPermissionCodes } from '@/router/dynamic-menu'
import {
  createMenu,
  deleteMenu,
  getAdminMenus,
  updateMenu,
  updateMenuStatus,
} from '../api/menu'
import { MenuStatus, MenuType, type AdminMenu } from '@/modules/iam/types/menu'
import { routeComponentOptions } from '@/router/route-components'
import type { MenuFormModel, MenuQuery } from '../types/menu-page'
import {
  buildMenuPayload,
  buildMenuParentOptions,
  defaultMenuQuery,
  defaultMenuFormModel,
  filterMenus,
  flattenMenus,
  menuFormTypeOptions,
  menuStatusOptions,
  menuTypeOptions,
  toMenuFormModel,
} from './menu-page.utils'

// 菜单管理页面组合式函数，封装菜单树加载、创建、编辑、删除、状态切换等逻辑
export function useMenuPage() {
  const message = useMessage()
  const { closeSuccess, showSuccess, successText } = useSuccessFeedback()
  const loading = ref(false)
  const saving = ref(false)
  const menus = ref<AdminMenu[]>([])
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')
  const formRef = ref<FormInst | null>(null)
  const expandedRowKeys = ref<Array<string | number>>([])

  // 搜索查询条件
  const query = reactive<MenuQuery>(defaultMenuQuery())

  // 菜单表单数据
  const formModel = reactive<MenuFormModel>(defaultMenuFormModel())

  // 菜单类型筛选选项
  const typeOptions: SelectOption[] = menuTypeOptions

  // 表单中的菜单类型选项（不含"全部"）
  const formTypeOptions: SelectOption[] = menuFormTypeOptions

  // 状态筛选选项
  const statusOptions: SelectOption[] = menuStatusOptions

  // 表单校验规则
  const rules: FormRules = {
    code: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
    title: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  }

  // 扁平化后的菜单列表（用于统计和搜索）
  const flatMenus = computed(() => flattenMenus(menus.value))

  // 所有菜单行的 key 集合（用于展开全部）
  const allRowKeys = computed(() => flatMenus.value.map((m) => m.id))

  // 根据搜索条件过滤后的菜单树
  const displayMenus = computed(() => filterMenus(menus.value, query))

  // 目录数量统计
  const directoryCount = computed(() => flatMenus.value.filter((menu) => menu.type === MenuType.Directory).length)

  // 菜单数量统计
  const menuCount = computed(() => flatMenus.value.filter((menu) => menu.type === MenuType.Menu).length)

  // 按钮数量统计
  const buttonCount = computed(() => flatMenus.value.filter((menu) => menu.type === MenuType.Button).length)

  // 上级菜单选择选项，排除按钮类型和当前编辑的菜单
  const parentOptions = computed<SelectOption[]>(() => {
    return buildMenuParentOptions(flatMenus.value, formModel.id)
  })

  // 判断当前用户是否拥有指定按钮权限码
  function canUse(code: string) {
    return buttonPermissionCodes.value.includes(code)
  }

  // 展开所有菜单行
  function expandAll() {
    expandedRowKeys.value = allRowKeys.value
  }

  // 折叠所有菜单行
  function collapseAll() {
    expandedRowKeys.value = []
  }

  // 处理展开/折叠行变化
  function handleExpandedChange(keys: Array<string | number>) {
    expandedRowKeys.value = keys
  }

  // 重置菜单表单为默认值
  function resetForm() {
    Object.assign(formModel, defaultMenuFormModel())
  }

  // 从服务端加载菜单树数据
  async function loadMenus() {
    loading.value = true
    try {
      menus.value = await getAdminMenus()
      expandedRowKeys.value = allRowKeys.value
    } finally {
      loading.value = false
    }
  }

  // 打开创建根级菜单的弹窗
  function openCreateRoot() {
    formMode.value = 'create'
    resetForm()
    formVisible.value = true
  }

  // 打开创建子菜单的弹窗，自动设置上级菜单和默认类型
  function openCreateChild(row: AdminMenu) {
    formMode.value = 'create'
    resetForm()
    formModel.parent_id = row.id
    formModel.type = row.type === MenuType.Directory ? MenuType.Menu : MenuType.Button
    formModel.sort = row.type === MenuType.Directory ? 1 : 10
    formVisible.value = true
  }

  // 打开编辑菜单的弹窗，将当前行数据填充到表单
  function openEdit(row: AdminMenu) {
    formMode.value = 'edit'
    Object.assign(formModel, toMenuFormModel(row))
    formVisible.value = true
  }

  // 提交菜单表单（新建或更新）
  async function handleSubmit() {
    await formRef.value?.validate()
    saving.value = true
    try {
      const payload = buildMenuPayload(formModel)

      if (formMode.value === 'create') {
        await createMenu({
          ...payload,
          code: formModel.code.trim(),
        })
        showSuccess('菜单创建成功')
        message.success('菜单创建成功')
      } else {
        await updateMenu(formModel.id, payload)
        showSuccess('菜单信息已更新')
        message.success('菜单信息已更新')
      }

      await loadMenus()
      formVisible.value = false
    } finally {
      saving.value = false
    }
  }

  // 切换菜单的启用/禁用状态
  async function handleToggleStatus(row: AdminMenu, status: MenuStatus) {
    await updateMenuStatus(row.id, { status })
    showSuccess(`菜单已${status === MenuStatus.Enabled ? '启用' : '禁用'}`)
    message.success('菜单状态已更新')
    await loadMenus()
  }

  // 删除指定菜单
  async function handleDelete(row: AdminMenu) {
    await deleteMenu(row.id)
    showSuccess('菜单已删除')
    message.success('菜单已删除')
    await loadMenus()

    if (formModel.id === row.id) {
      formVisible.value = false
    }
  }

  // 重置搜索条件
  function handleResetQuery() {
    Object.assign(query, defaultMenuQuery())
  }

  // 组件挂载时自动加载菜单列表
  onMounted(loadMenus)

  return {
    buttonCount,
    canUse,
    closeSuccess,
    collapseAll,
    componentOptions: routeComponentOptions,
    directoryCount,
    displayMenus,
    expandAll,
    expandedRowKeys,
    flatMenus,
    formMode,
    formModel,
    formRef,
    formTypeOptions,
    formVisible,
    handleDelete,
    handleExpandedChange,
    handleResetQuery,
    handleSubmit,
    handleToggleStatus,
    loadMenus,
    loading,
    menuCount,
    openCreateChild,
    openCreateRoot,
    openEdit,
    parentOptions,
    query,
    rules,
    saving,
    statusOptions,
    successText,
    typeOptions,
  }
}
