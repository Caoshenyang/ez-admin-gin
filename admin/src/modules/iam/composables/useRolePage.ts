import type { FormRules, SelectOption, TreeOption } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref, watch } from 'vue'

import { buttonPermissionCodes } from '@/router/dynamic-menu'
import { getDepartments } from '../api/department'
import { getAdminMenus } from '../api/menu'
import {
  createRole,
  deleteRole,
  getRoles,
  updateRole,
  updateRoleMenus,
  updateRolePermissions,
  updateRoleStatus,
} from '../api/role'
import { getUsers } from '../api/user'
import { buildDepartmentTreeOptions, flattenDepartments } from './department-page.utils'
import type { DepartmentItem } from '../types/department'
import { MenuStatus, type AdminMenu } from '@/modules/iam/types/menu'
import {
  RoleStatus,
  type RoleItem,
  type RoleListQuery,
} from '../types/role'
import type { UserItem } from '../types/user'
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
  roleDataScopeHelps,
  roleDataScopeLabels,
  roleDataScopeOptions,
  permissionMethodOptions,
  roleFormRules,
  roleStatusOptions,
  superAdminRoleCode,
  toPermissionRows,
  toFeaturePermissionTreeOptions,
  toRoleFormModel,
} from './role-page.utils'

export function useRolePage() {
  const message = useMessage()
  const loading = ref(false)
  const saving = ref(false)
  const relatedUsersLoading = ref(false)
  const roles = ref<RoleItem[]>([])
  const menus = ref<AdminMenu[]>([])
  const departments = ref<DepartmentItem[]>([])
  const relatedUsers = ref<UserItem[]>([])
  const relatedUsersTotal = ref(0)
  const selectedRoleID = ref<number | null>(null)
  const activeTab = ref<PermissionTab>('feature')
  const checkedMenuIDs = ref<Array<string | number>>([])
  const permissionRows = ref<PermissionRow[]>([])

  const query = reactive<RoleListQuery>(defaultRoleListQuery())

  const formRef = ref()
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')
  const formModel = reactive<RoleFormModel>(defaultRoleFormModel())

  const statusOptions: SelectOption[] = roleStatusOptions

  const methodOptions: SelectOption[] = permissionMethodOptions

  const dataScopeOptions: SelectOption[] = roleDataScopeOptions

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

  const menuTreeOptions = computed<TreeOption[]>(() => toFeaturePermissionTreeOptions(menus.value))

  const departmentTreeOptions = computed(() => buildDepartmentTreeOptions(departments.value))

  const departmentNameMap = computed(() =>
    new Map(flattenDepartments(departments.value).map((department) => [department.id, department.name])),
  )

  const dataScopeLabel = computed(() => {
    if (!selectedRole.value) {
      return ''
    }

    return roleDataScopeLabels.get(selectedRole.value.data_scope) ?? selectedRole.value.data_scope
  })

  const dataScopeHelp = computed(() => {
    if (!selectedRole.value) {
      return ''
    }

    return roleDataScopeHelps.get(selectedRole.value.data_scope) ?? ''
  })

  const allMenus = computed(() => flattenRoleMenus(menus.value))

  const checkedFeatureCount = computed(() => checkedMenuIDs.value.length)

  const canEditSelectedRole = computed(() => selectedRole.value !== null && selectedRole.value.code !== superAdminRoleCode)

  const canEditBaseRole = computed(() => canEditSelectedRole.value && canUse('system:role:update'))

  const canToggleSelectedRoleStatus = computed(() => canEditSelectedRole.value && canUse('system:role:status'))

  const canDeleteSelectedRole = computed(() => canEditSelectedRole.value && canUse('system:role:delete'))

  const canEditFeaturePermissions = computed(() => canEditSelectedRole.value && canUse('system:role:menu'))

  const canEditApiPermissions = computed(() => canEditSelectedRole.value && canUse('system:role:permission'))

  const canEditPermissionTab = computed(() =>
    activeTab.value === 'api' ? canEditApiPermissions.value : canEditFeaturePermissions.value,
  )

  const canSavePermissionTab = computed(() => canEditPermissionTab.value)

  const permissionSaveLabel = computed(() =>
    activeTab.value === 'api' ? '保存接口权限' : '保存功能权限',
  )

  // 角色切换后要同步整块权限面板，避免旧角色的勾选和接口行残留。
  watch(selectedRole, (role) => {
    if (!role) {
      checkedMenuIDs.value = []
      permissionRows.value = []
      relatedUsers.value = []
      relatedUsersTotal.value = 0
      return
    }

    checkedMenuIDs.value = [...(role.menu_ids ?? [])]
    permissionRows.value = toPermissionRows(role)
    void loadRelatedUsers()
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

  async function loadDepartments() {
    departments.value = await getDepartments()
  }

  async function loadRelatedUsers() {
    if (!selectedRoleID.value) {
      relatedUsers.value = []
      relatedUsersTotal.value = 0
      return
    }

    relatedUsersLoading.value = true
    try {
      const data = await getUsers({
        page: 1,
        page_size: 20,
        role_id: selectedRoleID.value,
        status: 0,
      })
      relatedUsers.value = data.items
      relatedUsersTotal.value = data.total
    } finally {
      relatedUsersLoading.value = false
    }
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
        message.success('角色创建成功')
      } else {
        await updateRole(formModel.id, buildRoleUpdatePayload(formModel))
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
    message.success('角色状态已更新')
    await loadRoles()
  }

  async function handleDeleteRole(role: RoleItem) {
    await deleteRole(role.id)
    message.success('角色已删除')
    await loadRoles()
  }

  function handleCheckAll() {
    checkedMenuIDs.value = allMenus.value
      .filter((menu) => menu.status === MenuStatus.Enabled)
      .map((menu) => menu.id)
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
    if (!selectedRole.value || !canSavePermissionTab.value) {
      return
    }

    saving.value = true
    try {
      if (activeTab.value === 'api') {
        const permissions = normalizeRolePermissions(permissionRows.value)
        await updateRolePermissions(selectedRole.value.id, { permissions })
        message.success('接口权限已更新')
      } else {
        await updateRoleMenus(selectedRole.value.id, {
          menu_ids: checkedMenuIDs.value.map(Number),
        })
        message.success('功能权限已更新')
      }

      await loadRoles()
    } finally {
      saving.value = false
    }
  }

  onMounted(async () => {
    await Promise.all([loadMenus(), loadDepartments(), loadRoles()])
  })

  return {
    activeTab,
    addPermissionRow,
    canEditBaseRole,
    canDeleteSelectedRole,
    canEditPermissionTab,
    canEditSelectedRole,
    canSavePermissionTab,
    canToggleSelectedRoleStatus,
    canUse,
    checkedFeatureCount,
    checkedMenuIDs,
    dataScopeHelp,
    dataScopeLabel,
    dataScopeOptions,
    departmentNameMap,
    departmentTreeOptions,
    filteredRoles,
    formMode,
    formModel,
    formRef,
    formVisible,
    handleCheckAll,
    handleClearAll,
    handleDeleteRole,
    handleReset,
    handleSavePermissions,
    handleSearch,
    handleToggleRoleStatus,
    loadRelatedUsers,
    menuTreeOptions,
    methodOptions,
    openCreate,
    openEdit,
    permissionSaveLabel,
    permissionRows,
    query,
    relatedUsers,
    relatedUsersLoading,
    relatedUsersTotal,
    removePermissionRow,
    rules,
    saving,
    selectRole,
    selectedRole,
    selectedRoleID,
    statusOptions,
    statusType,
    submitRole,
    superAdminRoleCode,
    loading,
  }
}
