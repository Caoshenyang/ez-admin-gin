import type { FormRules, SelectOption, TreeOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref, watch } from 'vue'

import { useSuccessFeedback } from '@/composables/useSuccessFeedback'
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
} from '../types/role'
import type { PermissionRow, PermissionTab, RoleFormModel } from '../types/role-page'
import {
  buildRoleCreatePayload,
  buildRoleUpdatePayload,
  defaultRoleListQuery,
  defaultPermissionRow,
  defaultRoleFormModel,
  flattenRoleMenus,
  getRoleStatusTagType,
  normalizeRolePermissions,
  permissionMethodOptions,
  roleFormRules,
  roleStatusOptions,
  superAdminRoleCode,
  toPermissionRows,
  toRoleFormModel,
  toRoleMenuTreeOption,
} from './role-page.utils'

export function useRolePage() {
  const message = useMessage()
  const { closeSuccess, showSuccess, successText } = useSuccessFeedback()
  const loading = ref(false)
  const saving = ref(false)
  const roles = ref<RoleItem[]>([])
  const menus = ref<AdminMenu[]>([])
  const selectedRoleID = ref<number | null>(null)
  const activeTab = ref<PermissionTab>('base')
  const checkedMenuIDs = ref<Array<string | number>>([])
  const permissionRows = ref<PermissionRow[]>([])

  const query = reactive<RoleListQuery>(defaultRoleListQuery())

  const formRef = ref()
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')
  const formModel = reactive<RoleFormModel>(defaultRoleFormModel())

  const statusOptions: SelectOption[] = roleStatusOptions

  const methodOptions: SelectOption[] = permissionMethodOptions

  const rules: FormRules = roleFormRules

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

  const menuTreeOptions = computed<TreeOption[]>(() => menus.value.map(toRoleMenuTreeOption))

  const allMenus = computed(() => flattenRoleMenus(menus.value))

  const menuIDSet = computed(() => new Set(allMenus.value.filter((menu) => menu.type !== MenuType.Button).map((menu) => menu.id)))

  const buttonIDSet = computed(() => new Set(allMenus.value.filter((menu) => menu.type === MenuType.Button).map((menu) => menu.id)))

  const checkedMenuCount = computed(() => checkedMenuIDs.value.filter((id) => menuIDSet.value.has(Number(id))).length)

  const checkedButtonCount = computed(() => checkedMenuIDs.value.filter((id) => buttonIDSet.value.has(Number(id))).length)

  const checkedTotal = computed(() => checkedMenuIDs.value.length)

  const canEditSelectedRole = computed(() => selectedRole.value !== null && selectedRole.value.code !== superAdminRoleCode)

  // 角色切换后要同步整块权限面板，避免旧角色的勾选和接口行残留。
  watch(selectedRole, (role) => {
    if (!role) {
      checkedMenuIDs.value = []
      permissionRows.value = []
      return
    }

    checkedMenuIDs.value = [...(role.menu_ids ?? [])]
    permissionRows.value = toPermissionRows(role)
  })

  function canUse(code: string) {
    return buttonPermissionCodes.value.includes(code)
  }

  function statusType(status: RoleStatus) {
    return getRoleStatusTagType(status)
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
    Object.assign(query, defaultRoleListQuery())
    void loadRoles()
  }

  function openCreate() {
    formMode.value = 'create'
    Object.assign(formModel, defaultRoleFormModel())
    formVisible.value = true
  }

  function openEdit(role: RoleItem) {
    formMode.value = 'edit'
    Object.assign(formModel, toRoleFormModel(role))
    formVisible.value = true
  }

  async function submitRole() {
    await formRef.value?.validate()
    saving.value = true
    try {
      if (formMode.value === 'create') {
        const created = await createRole(buildRoleCreatePayload(formModel))
        selectedRoleID.value = created.id
        showSuccess('角色创建成功')
        message.success('角色创建成功')
      } else {
        await updateRole(formModel.id, buildRoleUpdatePayload(formModel))
        showSuccess('角色信息已更新')
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
    showSuccess(`角色已${status === RoleStatus.Enabled ? '启用' : '禁用'}`)
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
    permissionRows.value.push(defaultPermissionRow())
  }

  function removePermissionRow(id: number) {
    permissionRows.value = permissionRows.value.filter((row) => row.id !== id)
  }

  async function handleSavePermissions() {
    if (!selectedRole.value || !canEditSelectedRole.value) {
      return
    }

    saving.value = true
    try {
      if (activeTab.value === 'api') {
        const permissions = normalizeRolePermissions(permissionRows.value)
        await updateRolePermissions(selectedRole.value.id, { permissions })
        showSuccess('接口权限已更新')
        message.success('接口权限已更新')
      } else {
        await updateRoleMenus(selectedRole.value.id, {
          menu_ids: checkedMenuIDs.value.map(Number),
        })
        showSuccess('菜单与按钮权限已更新')
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
    closeSuccess,
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
