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
  // RoleItem 类型定义。
  type RoleItem,
  // RoleListQuery 类型定义。
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

// 角色管理页面组合式函数，封装角色列表、菜单权限、API权限分配等逻辑
export function useRolePage() {
  const message = useMessage()
  const { closeSuccess, showSuccess, successText } = useSuccessFeedback()
  const loading = ref(false)
  const saving = ref(false)
  const roles = ref<RoleItem[]>([])
  const menus = ref<AdminMenu[]>([])
  const selectedRoleID = ref<number | null>(null)
  const activeTab = ref<PermissionTab>('menu')
  const checkedMenuIDs = ref<Array<string | number>>([])
  const permissionRows = ref<PermissionRow[]>([])

  // 角色列表查询条件
  const query = reactive<RoleListQuery>(defaultRoleListQuery())

  const formRef = ref()
  const formVisible = ref(false)
  const formMode = ref<'create' | 'edit'>('create')
  const formModel = reactive<RoleFormModel>(defaultRoleFormModel())

  // 状态筛选选项
  const statusOptions: SelectOption[] = roleStatusOptions

  // HTTP方法选项
  const methodOptions: SelectOption[] = permissionMethodOptions

  // 表单校验规则
  const rules: FormRules = roleFormRules

  // 当前选中的角色对象
  const selectedRole = computed(() => roles.value.find((role) => role.id === selectedRoleID.value) ?? null)

  // 根据关键词和状态过滤后的角色列表
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

  // 菜单树形选项
  const menuTreeOptions = computed<TreeOption[]>(() => menus.value.map(toRoleMenuTreeOption))

  // 所有扁平化后的菜单
  const allMenus = computed(() => flattenRoleMenus(menus.value))

  // 非按钮类型的菜单ID集合（目录 + 菜单）
  const menuIDSet = computed(() => new Set(allMenus.value.filter((menu) => menu.type !== MenuType.Button).map((menu) => menu.id)))

  // 按钮类型的菜单ID集合
  const buttonIDSet = computed(() => new Set(allMenus.value.filter((menu) => menu.type === MenuType.Button).map((menu) => menu.id)))

  // 已勾选的菜单数量
  const checkedMenuCount = computed(() => checkedMenuIDs.value.filter((id) => menuIDSet.value.has(Number(id))).length)

  // 已勾选的按钮数量
  const checkedButtonCount = computed(() => checkedMenuIDs.value.filter((id) => buttonIDSet.value.has(Number(id))).length)

  // 已勾选的总数
  const checkedTotal = computed(() => checkedMenuIDs.value.length)

  // 是否可以编辑选中的角色（超级管理员不可编辑）
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

  // 判断当前用户是否拥有指定按钮权限码
  function canUse(code: string) {
    return buttonPermissionCodes.value.includes(code)
  }

  // 根据角色状态返回标签类型（success / error）
  function statusType(status: RoleStatus) {
    return getRoleStatusTagType(status)
  }

  // 选中指定角色
  function selectRole(role: RoleItem) {
    selectedRoleID.value = role.id
  }

  // 从服务端加载角色列表
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

  // 从服务端加载菜单树数据
  async function loadMenus() {
    menus.value = await getAdminMenus()
  }

  // 搜索角色
  async function handleSearch() {
    await loadRoles()
  }

  // 重置搜索条件并重新加载角色列表
  function handleReset() {
    Object.assign(query, defaultRoleListQuery())
    void loadRoles()
  }

  // 打开创建角色的弹窗
  function openCreate() {
    formMode.value = 'create'
    Object.assign(formModel, defaultRoleFormModel())
    formVisible.value = true
  }

  // 打开编辑角色的弹窗，将当前行数据填充到表单
  function openEdit(role: RoleItem) {
    formMode.value = 'edit'
    Object.assign(formModel, toRoleFormModel(role))
    formVisible.value = true
  }

  // 提交角色表单（新建或更新）
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

  // 切换角色的启用/禁用状态
  async function handleToggleRoleStatus(role: RoleItem) {
    const status = role.status === RoleStatus.Enabled ? RoleStatus.Disabled : RoleStatus.Enabled
    await updateRoleStatus(role.id, { status })
    showSuccess(`角色已${status === RoleStatus.Enabled ? '启用' : '禁用'}`)
    message.success('角色状态已更新')
    await loadRoles()
  }

  // 全选所有启用的菜单和按钮
  function handleCheckAll() {
    checkedMenuIDs.value = allMenus.value.filter((menu) => menu.status === MenuStatus.Enabled).map((menu) => menu.id)
  }

  // 清空所有勾选的菜单和按钮
  function handleClearAll() {
    checkedMenuIDs.value = []
  }

  // 新增一行API权限
  function addPermissionRow() {
    permissionRows.value.push(defaultPermissionRow())
  }

  // 删除指定行的API权限
  function removePermissionRow(id: number) {
    permissionRows.value = permissionRows.value.filter((row) => row.id !== id)
  }

  // 保存权限（菜单权限或API权限）
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

  // 组件挂载时并行加载菜单和角色数据
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
