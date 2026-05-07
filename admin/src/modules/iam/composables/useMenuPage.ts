import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import { buttonPermissionCodes } from '@/router/dynamic-menu'
import {
  createMenu,
  deleteMenu,
  getAdminMenus,
  updateMenu,
  updateMenuStatus,
} from '../api/menu'
import { MenuStatus, MenuType, type AdminMenu } from '@/modules/iam/types/menu'

export interface MenuFormModel {
  id: number
  parent_id: number
  type: MenuType
  code: string
  title: string
  path: string
  component: string
  icon: string
  sort: number
  status: MenuStatus
  remark: string
}

export interface MenuQuery {
  keyword: string
  type: 0 | MenuType
  status: 0 | MenuStatus
}

function defaultFormModel(): MenuFormModel {
  return {
    id: 0,
    parent_id: 0,
    type: MenuType.Directory,
    code: '',
    title: '',
    path: '',
    component: '',
    icon: '',
    sort: 10,
    status: MenuStatus.Enabled,
    remark: '',
  }
}

function flattenMenus(items: AdminMenu[]): AdminMenu[] {
  const result: AdminMenu[] = []

  for (const item of items) {
    result.push(item)
    result.push(...flattenMenus(item.children ?? []))
  }

  return result
}

function menuLevel(flatMenus: AdminMenu[], id: number) {
  let level = 0
  let current = flatMenus.find((menu) => menu.id === id)

  while (current && current.parent_id !== 0) {
    level += 1
    current = flatMenus.find((menu) => menu.id === current?.parent_id)
  }

  return level
}

function filterMenus(items: AdminMenu[], query: MenuQuery): AdminMenu[] {
  const keyword = query.keyword.trim().toLowerCase()
  const result: AdminMenu[] = []

  for (const item of items) {
    const children = filterMenus(item.children ?? [], query)
    const matchedKeyword =
      keyword === '' ||
      item.title.toLowerCase().includes(keyword) ||
      item.code.toLowerCase().includes(keyword) ||
      item.path.toLowerCase().includes(keyword)
    const matchedType = query.type === 0 || item.type === query.type
    const matchedStatus = query.status === 0 || item.status === query.status

    if ((matchedKeyword && matchedType && matchedStatus) || children.length > 0) {
      result.push({
        ...item,
        children: children.length > 0 ? children : undefined,
      })
    }
  }

  return result
}

export function useMenuPage() {
  const message = useMessage()
  const loading = ref(false)
  const saving = ref(false)
  const menus = ref<AdminMenu[]>([])
  const successText = ref('')
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')
  const formRef = ref<FormInst | null>(null)
  const expandedRowKeys = ref<Array<string | number>>([])

  const query = reactive<MenuQuery>({
    keyword: '',
    type: 0,
    status: MenuStatus.Enabled,
  })

  const formModel = reactive<MenuFormModel>(defaultFormModel())

  const typeOptions: SelectOption[] = [
    { label: '类型：全部', value: 0 },
    { label: '目录', value: MenuType.Directory },
    { label: '菜单', value: MenuType.Menu },
    { label: '按钮', value: MenuType.Button },
  ]

  const formTypeOptions: SelectOption[] = [
    { label: '目录', value: MenuType.Directory },
    { label: '菜单', value: MenuType.Menu },
    { label: '按钮', value: MenuType.Button },
  ]

  const statusOptions: SelectOption[] = [
    { label: '状态：全部', value: 0 },
    { label: '启用', value: MenuStatus.Enabled },
    { label: '禁用', value: MenuStatus.Disabled },
  ]

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
    const options: SelectOption[] = [{ label: '根节点', value: 0 }]

    for (const menu of flatMenus.value) {
      if (menu.type === MenuType.Button || menu.id === formModel.id) {
        continue
      }

      options.push({
        label: `${'　'.repeat(menuLevel(flatMenus.value, menu.id))}${menu.title}`,
        value: menu.id,
      })
    }

    return options
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
    Object.assign(formModel, defaultFormModel())
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
    Object.assign(formModel, {
      id: row.id,
      parent_id: row.parent_id,
      type: row.type,
      code: row.code,
      title: row.title,
      path: row.path,
      component: row.component,
      icon: row.icon,
      sort: row.sort,
      status: row.status,
      remark: row.remark,
    })
    formVisible.value = true
  }

  function normalizedPayload() {
    const isButton = formModel.type === MenuType.Button

    return {
      parent_id: formModel.parent_id,
      type: formModel.type,
      title: formModel.title.trim(),
      path: isButton ? '' : formModel.path.trim(),
      component: isButton ? '' : formModel.component.trim(),
      icon: formModel.icon.trim(),
      sort: formModel.sort,
      status: formModel.status,
      remark: formModel.remark.trim(),
    }
  }

  async function handleSubmit() {
    await formRef.value?.validate()
    saving.value = true
    try {
      const payload = normalizedPayload()

      if (formMode.value === 'create') {
        await createMenu({
          ...payload,
          code: formModel.code.trim(),
        })
        successText.value = '菜单创建成功'
        message.success('菜单创建成功')
      } else {
        await updateMenu(formModel.id, payload)
        successText.value = '菜单信息已更新'
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
    successText.value = `菜单已${status === MenuStatus.Enabled ? '启用' : '禁用'}`
    message.success('菜单状态已更新')
    await loadMenus()
  }

  async function handleDelete(row: AdminMenu) {
    await deleteMenu(row.id)
    successText.value = '菜单已删除'
    message.success('菜单已删除')
    await loadMenus()

    if (formModel.id === row.id) {
      formVisible.value = false
    }
  }

  function handleResetQuery() {
    query.keyword = ''
    query.type = 0
    query.status = MenuStatus.Enabled
  }

  onMounted(loadMenus)

  return {
    buttonCount,
    canUse,
    collapseAll,
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
