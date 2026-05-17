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

  const query = reactive<MenuQuery>(defaultMenuQuery())

  const formModel = reactive<MenuFormModel>(defaultMenuFormModel())

  const typeOptions: SelectOption[] = menuTypeOptions

  const formTypeOptions: SelectOption[] = menuFormTypeOptions

  const statusOptions: SelectOption[] = menuStatusOptions

  const rules: FormRules = {
    code: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
    title: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  }

  const flatMenus = computed(() => flattenMenus(menus.value))

  const allRowKeys = computed(() => flatMenus.value.map((m) => m.id))

  const displayMenus = computed(() => filterMenus(menus.value, query))

  const directoryCount = computed(() => flatMenus.value.filter((menu) => menu.type === MenuType.Directory).length)

  const menuCount = computed(() => flatMenus.value.filter((menu) => menu.type === MenuType.Menu).length)

  const buttonCount = computed(() => flatMenus.value.filter((menu) => menu.type === MenuType.Button).length)

  const parentOptions = computed<SelectOption[]>(() => {
    return buildMenuParentOptions(flatMenus.value, formModel.id)
  })

  function canUse(code: string) {
    return buttonPermissionCodes.value.includes(code)
  }

  function expandAll() {
    expandedRowKeys.value = allRowKeys.value
  }

  function collapseAll() {
    expandedRowKeys.value = []
  }

  function handleExpandedChange(keys: Array<string | number>) {
    expandedRowKeys.value = keys
  }

  function resetForm() {
    Object.assign(formModel, defaultMenuFormModel())
  }

  async function loadMenus() {
    loading.value = true
    try {
      menus.value = await getAdminMenus()
      expandedRowKeys.value = allRowKeys.value
    } finally {
      loading.value = false
    }
  }

  function openCreateRoot() {
    formMode.value = 'create'
    resetForm()
    formVisible.value = true
  }

  function openCreateChild(row: AdminMenu) {
    formMode.value = 'create'
    resetForm()
    formModel.parent_id = row.id
    formModel.type = row.type === MenuType.Directory ? MenuType.Menu : MenuType.Button
    formModel.sort = row.type === MenuType.Directory ? 1 : 10
    formVisible.value = true
  }

  function openEdit(row: AdminMenu) {
    formMode.value = 'edit'
    Object.assign(formModel, toMenuFormModel(row))
    formVisible.value = true
  }

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

  async function handleToggleStatus(row: AdminMenu, status: MenuStatus) {
    await updateMenuStatus(row.id, { status })
    showSuccess(`菜单已${status === MenuStatus.Enabled ? '启用' : '禁用'}`)
    message.success('菜单状态已更新')
    await loadMenus()
  }

  async function handleDelete(row: AdminMenu) {
    await deleteMenu(row.id)
    showSuccess('菜单已删除')
    message.success('菜单已删除')
    await loadMenus()

    if (formModel.id === row.id) {
      formVisible.value = false
    }
  }

  function handleResetQuery() {
    Object.assign(query, defaultMenuQuery())
  }

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
