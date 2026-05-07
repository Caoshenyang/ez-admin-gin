import type { FormRules, SelectOption, TreeOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref, watch } from 'vue'

import { buttonPermissionCodes } from '@/router/dynamic-menu'
import { getAdminMenus } from '../api/menu'
import {
  createRole,
  getRoles,
  updateRole,
  updateRoleMenus,
  updateRolePermissions,
  updateRoleStatus,
} from '../api/role'
import { MenuStatus, MenuType, type AdminMenu } from '@/modules/iam/types/menu'
import {
  RoleStatus,
  type RoleItem,
  type RoleListQuery,
  type RolePermissionItem,
} from '../types/role'

export interface RoleFormModel {
  id: number
  code: string
  name: string
  sort: number
  status: RoleStatus
  remark: string
}

export interface PermissionRow {
  id: number
  path: string
  method: string
}

export type PermissionTab = 'menu' | 'button' | 'api'

const superAdminRoleCode = 'super_admin'

function defaultRoleFormModel(): RoleFormModel {
  return {
    id: 0,
    code: '',
    name: '',
    sort: 10,
    status: RoleStatus.Enabled,
    remark: '',
  }
}

function toTreeOption(menu: AdminMenu): TreeOption {
  const typeText = menu.type === MenuType.Directory ? '目录' : menu.type === MenuType.Menu ? '菜单' : '按钮'
  const statusText = menu.status === MenuStatus.Enabled ? '' : '（禁用）'

  return {
    key: menu.id,
    label: `${menu.title}  ${typeText}  ${menu.code}${statusText}`,
    children: menu.children?.map(toTreeOption),
    disabled: menu.status !== MenuStatus.Enabled,
  }
}

function flattenMenus(items: AdminMenu[]) {
  const result: AdminMenu[] = []

  for (const item of items) {
    result.push(item)
    result.push(...flattenMenus(item.children ?? []))
  }

  return result
}

export function useRolePage() {
  const message = useMessage()
  const loading = ref(false)
  const saving = ref(false)
  const roles = ref<RoleItem[]>([])
  const menus = ref<AdminMenu[]>([])
  const selectedRoleID = ref<number | null>(null)
  const activeTab = ref<PermissionTab>('menu')
  const checkedMenuIDs = ref<Array<string | number>>([])
  const permissionRows = ref<PermissionRow[]>([])
  const successText = ref('')

  const query = reactive<RoleListQuery>({
    page: 1,
    page_size: 100,
    keyword: '',
    status: 0,
  })

  const formRef = ref()
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')
  const formModel = reactive<RoleFormModel>(defaultRoleFormModel())

  const statusOptions: SelectOption[] = [
    { label: '状态：全部', value: 0 },
    { label: '启用', value: RoleStatus.Enabled },
    { label: '禁用', value: RoleStatus.Disabled },
  ]

  const methodOptions: SelectOption[] = [
    { label: 'GET', value: 'GET' },
    { label: 'POST', value: 'POST' },
    { label: 'PUT', value: 'PUT' },
    { label: 'DELETE', value: 'DELETE' },
  ]

  const rules: FormRules = {
    code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }],
    name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  }

  const selectedRole = computed(() => roles.value.find((role) => role.id === selectedRoleID.value) ?? null)

  const filteredRoles = computed(() => {
    const keyword = query.keyword?.trim().toLowerCase() ?? ''

    return roles.value.filter((role) => {
      const matchedKeyword =
        keyword === '' ||
        role.code.toLowerCase().includes(keyword) ||
        role.name.toLowerCase().includes(keyword)
      const matchedStatus = query.status === 0 || role.status === query.status

      return matchedKeyword && matchedStatus
    })
  })

  const menuTreeOptions = computed<TreeOption[]>(() => menus.value.map(toTreeOption))
  const allMenus = computed(() => flattenMenus(menus.value))
  const menuIDSet = computed(() => new Set(allMenus.value.filter((menu) => menu.type !== MenuType.Button).map((menu) => menu.id)))
  const buttonIDSet = computed(() => new Set(allMenus.value.filter((menu) => menu.type === MenuType.Button).map((menu) => menu.id)))
  const checkedMenuCount = computed(() => checkedMenuIDs.value.filter((id) => menuIDSet.value.has(Number(id))).length)
  const checkedButtonCount = computed(() => checkedMenuIDs.value.filter((id) => buttonIDSet.value.has(Number(id))).length)
  const checkedTotal = computed(() => checkedMenuIDs.value.length)
  const canEditSelectedRole = computed(() => selectedRole.value !== null && selectedRole.value.code !== superAdminRoleCode)

  watch(selectedRole, (role) => {
    if (!role) {
      checkedMenuIDs.value = []
      permissionRows.value = []
      return
    }

    checkedMenuIDs.value = [...role.menu_ids]
    permissionRows.value = role.permissions.map((permission, index) => ({
      id: index + 1,
      path: permission.path,
      method: permission.method,
    }))
  })

  function canUse(code: string) {
    return buttonPermissionCodes.value.includes(code)
  }

  function statusType(status: RoleStatus) {
    return status === RoleStatus.Enabled ? 'success' : 'error'
  }

  function selectRole(role: RoleItem) {
    selectedRoleID.value = role.id
  }

  async function loadRoles() {
    loading.value = true
    try {
      const data = await getRoles({
        page: 1,
        page_size: 100,
        keyword: undefined,
        status: 0,
      })
      roles.value = data.items

      if (!selectedRoleID.value && data.items.length > 0) {
        selectedRoleID.value = data.items[0]?.id ?? null
      }
      if (selectedRoleID.value && !data.items.some((role) => role.id === selectedRoleID.value)) {
        selectedRoleID.value = data.items[0]?.id ?? null
      }
    } finally {
      loading.value = false
    }
  }

  async function loadMenus() {
    menus.value = await getAdminMenus()
  }

  async function handleSearch() {
    await loadRoles()
  }

  function handleReset() {
    query.keyword = ''
    query.status = 0
    void loadRoles()
  }

  function openCreate() {
    formMode.value = 'create'
    Object.assign(formModel, defaultRoleFormModel())
    formVisible.value = true
  }

  function openEdit(role: RoleItem) {
    formMode.value = 'edit'
    Object.assign(formModel, {
      id: role.id,
      code: role.code,
      name: role.name,
      sort: role.sort,
      status: role.status,
      remark: role.remark,
    })
    formVisible.value = true
  }

  async function submitRole() {
    await formRef.value?.validate()
    saving.value = true
    try {
      if (formMode.value === 'create') {
        const created = await createRole({
          code: formModel.code.trim(),
          name: formModel.name.trim(),
          sort: formModel.sort,
          status: formModel.status,
          remark: formModel.remark.trim(),
        })
        selectedRoleID.value = created.id
        successText.value = '角色创建成功'
        message.success('角色创建成功')
      } else {
        await updateRole(formModel.id, {
          name: formModel.name.trim(),
          sort: formModel.sort,
          status: formModel.status,
          remark: formModel.remark.trim(),
        })
        successText.value = '角色信息已更新'
        message.success('角色信息已更新')
      }

      formVisible.value = false
      await loadRoles()
    } finally {
      saving.value = false
    }
  }

  async function handleToggleRoleStatus(role: RoleItem) {
    const status = role.status === RoleStatus.Enabled ? RoleStatus.Disabled : RoleStatus.Enabled
    await updateRoleStatus(role.id, { status })
    successText.value = `角色已${status === RoleStatus.Enabled ? '启用' : '禁用'}`
    message.success('角色状态已更新')
    await loadRoles()
  }

  function handleCheckAll() {
    checkedMenuIDs.value = allMenus.value.filter((menu) => menu.status === MenuStatus.Enabled).map((menu) => menu.id)
  }

  function handleClearAll() {
    checkedMenuIDs.value = []
  }

  function addPermissionRow() {
    permissionRows.value.push({
      id: Date.now(),
      path: '',
      method: 'GET',
    })
  }

  function removePermissionRow(id: number) {
    permissionRows.value = permissionRows.value.filter((row) => row.id !== id)
  }

  function normalizePermissions(rows: PermissionRow[]): RolePermissionItem[] {
    const seen = new Set<string>()
    const result: RolePermissionItem[] = []

    for (const row of rows) {
      const path = row.path.trim()
      const method = row.method.trim().toUpperCase()

      if (!path || !method) continue

      const key = `${method} ${path}`
      if (seen.has(key)) continue

      seen.add(key)
      result.push({ path, method })
    }

    return result
  }

  async function handleSavePermissions() {
    if (!selectedRole.value || !canEditSelectedRole.value) {
      return
    }

    saving.value = true
    try {
      if (activeTab.value === 'api') {
        const permissions = normalizePermissions(permissionRows.value)
        await updateRolePermissions(selectedRole.value.id, { permissions })
        successText.value = '接口权限已更新'
        message.success('接口权限已更新')
      } else {
        await updateRoleMenus(selectedRole.value.id, {
          menu_ids: checkedMenuIDs.value.map(Number),
        })
        successText.value = '菜单与按钮权限已更新'
        message.success('菜单与按钮权限已更新')
      }

      await loadRoles()
    } finally {
      saving.value = false
    }
  }

  onMounted(async () => {
    await Promise.all([loadMenus(), loadRoles()])
  })

  return {
    activeTab,
    addPermissionRow,
    canEditSelectedRole,
    canUse,
    checkedButtonCount,
    checkedMenuCount,
    checkedMenuIDs,
    checkedTotal,
    filteredRoles,
    formMode,
    formModel,
    formRef,
    formVisible,
    handleCheckAll,
    handleClearAll,
    handleReset,
    handleSavePermissions,
    handleSearch,
    handleToggleRoleStatus,
    menuTreeOptions,
    methodOptions,
    openCreate,
    openEdit,
    permissionRows,
    query,
    removePermissionRow,
    rules,
    saving,
    selectRole,
    selectedRole,
    selectedRoleID,
    statusOptions,
    statusType,
    submitRole,
    successText,
    superAdminRoleCode,
    loading,
  }
}
