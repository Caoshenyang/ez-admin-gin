import type { DataTableRowKey, FormInst, FormRules, SelectOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import { buttonPermissionCodes } from '@/router/dynamic-menu'
import { createMenu, deleteMenu, getAdminMenus, updateMenu, updateMenuStatus } from '../api/menu'
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
  const loading = ref(false)
  const saving = ref(false)
  const menus = ref<AdminMenu[]>([])
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')
  const formRef = ref<FormInst | null>(null)
  const checkedRowKeys = ref<DataTableRowKey[]>([])
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

  const directoryCount = computed(
    () => flatMenus.value.filter((menu) => menu.type === MenuType.Directory).length,
  )

  const menuCount = computed(
    () => flatMenus.value.filter((menu) => menu.type === MenuType.Menu).length,
  )

  const buttonCount = computed(
    () => flatMenus.value.filter((menu) => menu.type === MenuType.Button).length,
  )

  const selectedCount = computed(() => checkedRowKeys.value.length)

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

  function handleCheckedRowKeys(keys: DataTableRowKey[]) {
    checkedRowKeys.value = keys
  }

  function resetForm() {
    Object.assign(formModel, defaultMenuFormModel())
  }

  async function loadMenus() {
    loading.value = true
    try {
      menus.value = await getAdminMenus()
      expandedRowKeys.value = allRowKeys.value
      checkedRowKeys.value = checkedRowKeys.value.filter((key) =>
        allRowKeys.value.includes(Number(key)),
      )
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
        message.success('菜单创建成功')
      } else {
        await updateMenu(formModel.id, payload)
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
    message.success('菜单状态已更新')
    await loadMenus()
  }

  async function handleDelete(row: AdminMenu) {
    await deleteMenu(row.id)
    message.success('菜单已删除')
    await loadMenus()
    checkedRowKeys.value = checkedRowKeys.value.filter((key) => Number(key) !== row.id)

    if (formModel.id === row.id) {
      formVisible.value = false
    }
  }

  async function handleDeleteSelected() {
    const selectedIds = checkedRowKeys.value.map(Number)

    for (const id of selectedIds) {
      await deleteMenu(id)
    }

    message.success('选中菜单已删除')
    checkedRowKeys.value = []
    await loadMenus()

    if (selectedIds.includes(formModel.id)) {
      formVisible.value = false
    }
  }

  function handleResetQuery() {
    Object.assign(query, defaultMenuQuery())
  }

  function handleSearch() {
    expandedRowKeys.value = allRowKeys.value
  }

  onMounted(loadMenus)

  return {
    buttonCount,
    canUse,
    checkedRowKeys,
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
    handleCheckedRowKeys,
    handleDelete,
    handleDeleteSelected,
    handleExpandedChange,
    handleResetQuery,
    handleSearch,
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
    selectedCount,
    statusOptions,
    typeOptions,
  }
}
