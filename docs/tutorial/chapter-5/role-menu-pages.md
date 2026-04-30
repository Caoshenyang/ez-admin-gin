---
title: 角色与菜单页面
description: "实现角色管理和菜单管理前端页面。"
---

# 角色与菜单页面

上一节已经把“用户管理”接成了真实页面：可以查询用户、维护状态，并把用户绑定到角色。现在继续补齐权限体系的另一半：角色管理和菜单管理。

完成这一节后，侧边栏里的“角色管理”和“菜单管理”不再停留在占位页，而是可以进入真实管理页面。角色页面负责维护角色、接口权限和菜单权限；菜单页面负责维护目录、页面菜单和按钮权限。

::: tip 🎯 本节目标
这一节会把 `system/RoleView` 和 `system/MenuView` 从占位页换成真实页面，并补齐角色、菜单相关的类型和 API 封装。角色页面采用左侧列表 + 右侧权限面板布局；菜单页面使用 NDataTable 树形展示全宽菜单树，新增和编辑通过弹框表单完成。
:::

## 先看接口边界

角色管理接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/roles` | 角色分页列表 |
| `POST` | `/api/v1/system/roles` | 创建角色 |
| `POST` | `/api/v1/system/roles/:id/update` | 编辑角色基础信息 |
| `POST` | `/api/v1/system/roles/:id/status` | 修改角色状态 |
| `POST` | `/api/v1/system/roles/:id/permissions` | 替换角色接口权限 |
| `POST` | `/api/v1/system/roles/:id/menus` | 替换角色菜单权限 |

菜单管理接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/menus` | 获取完整菜单树 |
| `POST` | `/api/v1/system/menus` | 创建目录、菜单或按钮 |
| `POST` | `/api/v1/system/menus/:id/update` | 编辑菜单 |
| `POST` | `/api/v1/system/menus/:id/status` | 修改菜单状态 |
| `POST` | `/api/v1/system/menus/:id/delete` | 删除菜单 |

::: warning ⚠️ 菜单权限和接口权限是两件事
菜单权限决定“看不看得到入口”，接口权限决定“能不能真的调用接口”。只给角色分配菜单但没有接口权限，页面可能能打开，但请求会被后端拦截；只给接口权限但没有菜单，用户可能有能力访问接口，却没有侧边栏入口。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
└─ src/
   ├─ api/
   │  ├─ menu.ts
   │  └─ role.ts
   ├─ pages/
   │  └─ system/
   │     ├─ RoleView.vue
   │     └─ MenuView.vue
   ├─ router/
   │  └─ dynamic-menu.ts
   └─ types/
      ├─ menu.ts
      └─ role.ts
```

## 开始前先确认

开始之前，先确认下面几件事：

- 已完成上一节 [用户管理页面](./user-pages)。
- 登录后侧边栏能看到“角色管理”和“菜单管理”。
- 当前账号拥有角色与菜单相关按钮权限。
- 后端 `/api/v1/system/roles` 和 `/api/v1/system/menus` 可以正常返回数据。

## 🛠️ 完整代码

下面直接引入本节对应的完整项目文件，默认折叠。需要复制或对照时点击展开即可。

::: details `admin/src/types/role.ts` — 角色类型

```ts
export const RoleStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type RoleStatus = (typeof RoleStatus)[keyof typeof RoleStatus]

export interface RoleItem {
  id: number
  code: string
  name: string
  sort: number
  status: RoleStatus
  remark: string
  menu_ids: number[]
  permissions: Array<{
    path: string
    method: string
  }>
  created_at: string
  updated_at: string
}

export interface RoleListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: RoleStatus | 0
}

export interface RoleListResponse {
  items: RoleItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateRolePayload {
  code: string
  name: string
  sort: number
  status: RoleStatus
  remark: string
}

export interface UpdateRolePayload {
  name: string
  sort: number
  status: RoleStatus
  remark: string
}

export interface UpdateRoleStatusPayload {
  status: RoleStatus
}

export interface RolePermissionItem {
  path: string
  method: string
}

export interface UpdateRolePermissionsPayload {
  permissions: RolePermissionItem[]
}

export interface UpdateRoleMenusPayload {
  menu_ids: number[]
}
```

:::

::: details `admin/src/types/menu.ts` — 菜单类型

```ts
export const MenuType = {
  Directory: 1,
  Menu: 2,
  Button: 3,
} as const

export type MenuType = (typeof MenuType)[keyof typeof MenuType]

export const MenuStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type MenuStatus = (typeof MenuStatus)[keyof typeof MenuStatus]

// AuthMenu 对应 /api/v1/auth/menus 返回的菜单节点。
export interface AuthMenu {
  id: number
  parent_id: number
  type: MenuType
  code: string
  title: string
  path: string
  component: string
  icon: string
  sort: number
  children?: AuthMenu[]
}

export interface AdminMenu {
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
  children?: AdminMenu[]
  created_at: string
  updated_at: string
}

export interface CreateMenuPayload {
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

export type UpdateMenuPayload = Omit<CreateMenuPayload, 'code'>

export interface UpdateMenuStatusPayload {
  status: MenuStatus
}
```

:::

::: details `admin/src/api/role.ts` — 角色接口

```ts
import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  CreateRolePayload,
  RoleItem,
  RoleListQuery,
  RoleListResponse,
  UpdateRoleMenusPayload,
  UpdateRolePayload,
  UpdateRolePermissionsPayload,
  UpdateRoleStatusPayload,
} from '../types/role'

export async function getRoles(params: RoleListQuery) {
  const response = await http.get<ApiResponse<RoleListResponse>>('/system/roles', { params })
  return response.data.data
}

export async function createRole(payload: CreateRolePayload) {
  const response = await http.post<ApiResponse<RoleItem>>('/system/roles', payload)
  return response.data.data
}

export async function updateRole(id: number, payload: UpdateRolePayload) {
  const response = await http.post<ApiResponse<RoleItem>>(`/system/roles/${id}/update`, payload)
  return response.data.data
}

export async function updateRoleStatus(id: number, payload: UpdateRoleStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/roles/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function updateRolePermissions(id: number, payload: UpdateRolePermissionsPayload) {
  const response = await http.post<ApiResponse<{ id: number; permissions: unknown[] }>>(
    `/system/roles/${id}/permissions`,
    payload,
  )
  return response.data.data
}

export async function updateRoleMenus(id: number, payload: UpdateRoleMenusPayload) {
  const response = await http.post<ApiResponse<{ id: number; menu_ids: number[] }>>(
    `/system/roles/${id}/menus`,
    payload,
  )
  return response.data.data
}
```

:::

::: details `admin/src/api/menu.ts` — 菜单接口

```ts
import http from './http'

import type {
  AdminMenu,
  AuthMenu,
  CreateMenuPayload,
  UpdateMenuPayload,
  UpdateMenuStatusPayload,
} from '../types/menu'
import type { ApiResponse } from '../types/http'

// getCurrentUserMenus 获取当前登录用户可见的菜单树。
export async function getCurrentUserMenus() {
  const response = await http.get<ApiResponse<AuthMenu[]>>('/auth/menus')
  return response.data.data ?? []
}

export async function getAdminMenus() {
  const response = await http.get<ApiResponse<AdminMenu[]>>('/system/menus')
  return response.data.data ?? []
}

export async function createMenu(payload: CreateMenuPayload) {
  const response = await http.post<ApiResponse<AdminMenu>>('/system/menus', payload)
  return response.data.data
}

export async function updateMenu(id: number, payload: UpdateMenuPayload) {
  const response = await http.post<ApiResponse<AdminMenu>>(`/system/menus/${id}/update`, payload)
  return response.data.data
}

export async function updateMenuStatus(id: number, payload: UpdateMenuStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/menus/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function deleteMenu(id: number) {
  const response = await http.post<ApiResponse<{ id: number }>>(`/system/menus/${id}/delete`)
  return response.data.data
}
```

:::

::: details `admin/src/pages/system/RoleView.vue` — 角色权限页面

```vue
<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { FormInst, FormRules, SelectOption, TreeOption } from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NCheckbox,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NInputNumber,
  NModal,
  NPopconfirm,
  NSelect,
  NSpace,
  NTabPane,
  NTabs,
  NTag,
  NTree,
  useMessage,
} from 'naive-ui'
import { computed, onMounted, reactive, ref, watch } from 'vue'

import { getAdminMenus } from '../../api/menu'
import {
  createRole,
  getRoles,
  updateRole,
  updateRoleMenus,
  updateRolePermissions,
  updateRoleStatus,
} from '../../api/role'
import { buttonPermissionCodes } from '../../router/dynamic-menu'
import { MenuStatus, MenuType, type AdminMenu } from '../../types/menu'
import {
  RoleStatus,
  type RoleItem,
  type RoleListQuery,
  type RolePermissionItem,
} from '../../types/role'

interface RoleFormModel {
  id: number
  code: string
  name: string
  sort: number
  status: RoleStatus
  remark: string
}

interface PermissionRow {
  id: number
  path: string
  method: string
}

type PermissionTab = 'menu' | 'button' | 'api'

const superAdminRoleCode = 'super_admin'
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

const formRef = ref<FormInst | null>(null)
const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const formModel = reactive<RoleFormModel>({
  id: 0,
  code: '',
  name: '',
  sort: 10,
  status: RoleStatus.Enabled,
  remark: '',
})

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

const selectedRole = computed(() => {
  return roles.value.find((role) => role.id === selectedRoleID.value) ?? null
})

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

const menuTreeOptions = computed<TreeOption[]>(() => {
  return menus.value.map(toTreeOption)
})

const allMenus = computed(() => {
  return flattenMenus(menus.value)
})

const menuIDSet = computed(() => {
  return new Set(
    allMenus.value
      .filter((menu) => menu.type !== MenuType.Button)
      .map((menu) => menu.id),
  )
})

const buttonIDSet = computed(() => {
  return new Set(
    allMenus.value
      .filter((menu) => menu.type === MenuType.Button)
      .map((menu) => menu.id),
  )
})

const checkedMenuCount = computed(() => {
  return checkedMenuIDs.value.filter((id) => menuIDSet.value.has(Number(id))).length
})

const checkedButtonCount = computed(() => {
  return checkedMenuIDs.value.filter((id) => buttonIDSet.value.has(Number(id))).length
})

const checkedTotal = computed(() => checkedMenuIDs.value.length)

const canEditSelectedRole = computed(() => {
  return selectedRole.value !== null && selectedRole.value.code !== superAdminRoleCode
})

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

function formatTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

function statusType(status: RoleStatus) {
  return status === RoleStatus.Enabled ? 'success' : 'error'
}

function resetForm() {
  Object.assign(formModel, {
    id: 0,
    code: '',
    name: '',
    sort: 10,
    status: RoleStatus.Enabled,
    remark: '',
  })
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
  resetForm()
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

async function handleSubmitRole() {
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
  const status =
    role.status === RoleStatus.Enabled ? RoleStatus.Disabled : RoleStatus.Enabled
  await updateRoleStatus(role.id, { status })
  successText.value = `角色已${status === RoleStatus.Enabled ? '启用' : '禁用'}`
  message.success('角色状态已更新')
  await loadRoles()
}

function handleCheckedMenuIDs(keys: Array<string | number>) {
  checkedMenuIDs.value = keys
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
  permissionRows.value.push({
    id: Date.now(),
    path: '',
    method: 'GET',
  })
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

function normalizePermissions(rows: PermissionRow[]): RolePermissionItem[] {
  const seen = new Set<string>()
  const result: RolePermissionItem[] = []

  for (const row of rows) {
    const path = row.path.trim()
    const method = row.method.trim().toUpperCase()

    if (!path || !method) {
      continue
    }

    const key = `${method} ${path}`
    if (seen.has(key)) {
      continue
    }

    seen.add(key)
    result.push({ path, method })
  }

  return result
}

onMounted(async () => {
  await Promise.all([loadMenus(), loadRoles()])
})
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-4 overflow-hidden">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">角色权限</h1>
          <p class="mt-1 text-sm text-[#6B7280]">维护角色本身，以及角色拥有的菜单、按钮和接口权限。</p>
        </div>

        <NSpace>
          <NButton v-if="canUse('system:role:create')" type="primary" @click="openCreate">
            + 新增角色
          </NButton>
          <NButton
            type="primary"
            :loading="saving"
            :disabled="!canEditSelectedRole"
            @click="handleSavePermissions"
          >
            保存权限
          </NButton>
        </NSpace>
      </div>

      <NAlert
        v-if="successText"
        type="success"
        :show-icon="true"
        closable
        class="mx-auto w-full max-w-[520px]"
        @close="successText = ''"
      >
        {{ successText }}
      </NAlert>

      <div class="grid min-h-0 flex-1 grid-cols-[320px_minmax(0,1fr)] gap-4 overflow-hidden">
        <NCard class="rounded-lg" :bordered="false" content-style="height: 100%;">
          <div class="flex h-full flex-col overflow-hidden">
            <div class="mb-4">
              <h2 class="text-lg font-bold text-[#111827]">角色列表</h2>
              <p class="mt-1 text-xs text-[#6B7280]">点击左侧角色后，在右侧维护权限。</p>
            </div>

            <NSpace vertical :size="10" class="mb-4">
              <NInput
                v-model:value="query.keyword"
                clearable
                placeholder="角色编码 / 名称"
                @keyup.enter="handleSearch"
              />
              <div class="grid grid-cols-[1fr_auto] gap-2">
                <NSelect v-model:value="query.status" :options="statusOptions" />
                <NButton @click="handleReset">重置</NButton>
              </div>
            </NSpace>

            <div class="min-h-0 flex-1 space-y-3 overflow-y-auto pr-1">
              <button
                v-for="role in filteredRoles"
                :key="role.id"
                type="button"
                class="role-card"
                :class="{ 'role-card--active': role.id === selectedRoleID }"
                @click="selectRole(role)"
              >
                <span class="flex items-center justify-between gap-2">
                  <span class="min-w-0 truncate text-base font-bold text-[#111827]">
                    {{ role.name }}
                  </span>
                  <NTag :type="statusType(role.status)" :bordered="false" size="small">
                    {{ role.status === RoleStatus.Enabled ? '启用' : '禁用' }}
                  </NTag>
                </span>
                <span class="mt-1 block text-left text-xs text-[#6B7280]">
                  {{ role.code }} · 菜单 {{ role.menu_ids.length }} · 接口 {{ role.permissions.length }}
                </span>
                <span class="mt-2 flex items-center gap-2">
                  <NButton
                    v-if="canUse('system:role:update')"
                    size="tiny"
                    @click.stop="openEdit(role)"
                  >
                    编辑
                  </NButton>
                  <NPopconfirm
                    v-if="canUse('system:role:status') && role.code !== superAdminRoleCode"
                    @positive-click="handleToggleRoleStatus(role)"
                  >
                    <template #trigger>
                      <NButton
                        size="tiny"
                        :type="role.status === RoleStatus.Enabled ? 'error' : 'success'"
                        ghost
                        @click.stop
                      >
                        {{ role.status === RoleStatus.Enabled ? '禁用' : '启用' }}
                      </NButton>
                    </template>
                    确认{{ role.status === RoleStatus.Enabled ? '禁用' : '启用' }}该角色？
                  </NPopconfirm>
                </span>
              </button>
            </div>
          </div>
        </NCard>

        <NCard
          class="min-h-0 rounded-lg"
          :bordered="false"
          content-style="height: 100%; padding: 0;"
        >
          <div class="flex h-full flex-col overflow-hidden">
            <div class="border-b border-[#E5E7EB] px-5 py-5">
              <div class="flex items-start justify-between gap-4">
                <div>
                  <h2 class="text-lg font-bold text-[#111827]">菜单与按钮权限</h2>
                  <p class="mt-2 text-sm text-[#6B7280]">
                    当前角色：
                    <span class="font-semibold text-[#111827]">
                      {{ selectedRole?.name ?? '未选择' }}
                    </span>
                    。半选状态表示部分子权限已授权。
                  </p>
                </div>
                <NTag
                  v-if="selectedRole?.code === superAdminRoleCode"
                  type="warning"
                  :bordered="false"
                >
                  受保护角色
                </NTag>
              </div>
            </div>

            <div class="min-h-0 flex-1 overflow-y-auto px-5 py-4">
              <NTabs v-model:value="activeTab" type="segment" animated>
                <NTabPane name="menu" tab="菜单权限">
                  <div class="permission-toolbar">
                    <NCheckbox :checked="checkedTotal > 0" @update:checked="handleCheckAll">
                      全选
                    </NCheckbox>
                    <NButton text type="primary" @click="handleCheckAll">展开全部</NButton>
                    <NButton text type="primary" @click="handleClearAll">清空全部</NButton>
                  </div>

                  <NTree
                    checkable
                    cascade
                    block-line
                    selectable
                    :data="menuTreeOptions"
                    :checked-keys="checkedMenuIDs"
                    :disabled="!canEditSelectedRole"
                    @update:checked-keys="handleCheckedMenuIDs"
                  />
                </NTabPane>

                <NTabPane name="button" tab="按钮权限">
                  <div class="permission-toolbar">
                    <NButton text type="primary" @click="handleCheckAll">全选可用节点</NButton>
                    <NButton text type="primary" @click="handleClearAll">清空全部</NButton>
                  </div>

                  <NTree
                    checkable
                    cascade
                    block-line
                    selectable
                    :data="menuTreeOptions"
                    :checked-keys="checkedMenuIDs"
                    :disabled="!canEditSelectedRole"
                    @update:checked-keys="handleCheckedMenuIDs"
                  />
                </NTabPane>

                <NTabPane name="api" tab="接口权限">
                  <div class="mb-3 flex items-center justify-between">
                    <p class="text-sm text-[#6B7280]">
                      接口权限按请求路径和方法保存到 Casbin 策略表。
                    </p>
                    <NButton
                      size="small"
                      type="primary"
                      ghost
                      :disabled="!canEditSelectedRole"
                      @click="addPermissionRow"
                    >
                      + 添加接口
                    </NButton>
                  </div>

                  <div class="space-y-3">
                    <div
                      v-for="row in permissionRows"
                      :key="row.id"
                      class="grid grid-cols-[130px_minmax(0,1fr)_80px] items-center gap-3"
                    >
                      <NSelect
                        v-model:value="row.method"
                        :options="methodOptions"
                        :disabled="!canEditSelectedRole"
                      />
                      <NInput
                        v-model:value="row.path"
                        placeholder="/api/v1/system/users"
                        :disabled="!canEditSelectedRole"
                      />
                      <NButton
                        size="small"
                        type="error"
                        ghost
                        :disabled="!canEditSelectedRole"
                        @click="removePermissionRow(row.id)"
                      >
                        删除
                      </NButton>
                    </div>
                  </div>
                </NTabPane>
              </NTabs>
            </div>

            <div class="permission-summary">
              <span>已授权菜单：{{ checkedMenuCount }}</span>
              <span>按钮权限：{{ checkedButtonCount }}</span>
              <span>接口权限：{{ permissionRows.length }}</span>
            </div>
          </div>
        </NCard>
      </div>
    </section>

    <NModal
      v-model:show="formVisible"
      preset="card"
      :closable="false"
      class="role-modal"
      style="width: 560px; max-width: calc(100vw - 32px)"
    >
      <template #header>
        <div class="modal-header">
          <h2>{{ formMode === 'create' ? '新增角色' : '编辑角色' }}</h2>
          <p>
            {{
              formMode === 'create'
                ? '角色编码创建后会成为权限策略主体，建议使用稳定英文标识。'
                : '角色编码保持只读，避免影响已有权限策略。'
            }}
          </p>
          <button type="button" class="modal-close" @click="formVisible = false">
            <NIcon :size="18">
              <CloseOutline />
            </NIcon>
          </button>
        </div>
      </template>

      <NForm
        ref="formRef"
        :model="formModel"
        :rules="rules"
        label-placement="left"
        label-width="80"
      >
        <NFormItem label="角色编码" path="code">
          <NInput
            v-model:value="formModel.code"
            placeholder="demo_operator"
            :disabled="formMode === 'edit'"
          />
        </NFormItem>
        <NFormItem label="角色名称" path="name">
          <NInput v-model:value="formModel.name" placeholder="请输入角色名称" />
        </NFormItem>
        <NFormItem label="排序" path="sort">
          <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
        </NFormItem>
        <NFormItem label="状态" path="status">
          <NSelect v-model:value="formModel.status" :options="statusOptions.slice(1)" />
        </NFormItem>
        <NFormItem label="备注" path="remark">
          <NInput
            v-model:value="formModel.remark"
            type="textarea"
            placeholder="请输入备注"
            :autosize="{ minRows: 3, maxRows: 5 }"
          />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="flex justify-end gap-3">
          <NButton @click="formVisible = false">取消</NButton>
          <NButton type="primary" :loading="saving" @click="handleSubmitRole">
            保存
          </NButton>
        </div>
      </template>
    </NModal>
  </main>
</template>

<style scoped>
.role-card {
  width: 100%;
  border: 1px solid #dbe3ef;
  border-radius: 8px;
  background: #ffffff;
  padding: 14px 12px;
  text-align: left;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    box-shadow 0.2s ease;
}

.role-card:hover {
  border-color: #18a058;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
}

.role-card--active {
  border-color: #d7f2e4;
  background: #e9fbf1;
}

.permission-toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
  padding: 10px 12px;
  border-radius: 6px;
  background: #f7fafc;
}

.permission-summary {
  display: flex;
  gap: 32px;
  margin: 0 20px 20px;
  padding: 16px 18px;
  border-radius: 6px;
  background: #e9fbf1;
  color: #18a058;
  font-weight: 700;
}

.role-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 20px;
  border: 1px solid #dfe9f5;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.role-modal :deep(.n-card-header) {
  padding: 0;
}

.modal-header {
  position: relative;
  padding: 24px 28px;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.modal-header h2 {
  font-size: 19px;
  font-weight: 700;
  color: #111827;
}

.modal-header p {
  margin-top: 8px;
  max-width: 420px;
  font-size: 13px;
  line-height: 1.6;
  color: #64748b;
}

.modal-close {
  position: absolute;
  top: 20px;
  right: 22px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.82);
  color: #64748b;
}
</style>
```

:::

::: details `admin/src/pages/system/MenuView.vue` — 菜单管理页面

```vue
<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules, SelectOption } from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NInputNumber,
  NModal,
  NPopconfirm,
  NSelect,
  NSpace,
  NSwitch,
  NTag,
  NTooltip,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'

import {
  createMenu,
  deleteMenu,
  getAdminMenus,
  updateMenu,
  updateMenuStatus,
} from '../../api/menu'
import { buttonPermissionCodes } from '../../router/dynamic-menu'
import { MenuStatus, MenuType, type AdminMenu } from '../../types/menu'

interface MenuFormModel {
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

interface MenuQuery {
  keyword: string
  type: 0 | MenuType
  status: 0 | MenuStatus
}

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

const formModel = reactive<MenuFormModel>({
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
})

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

const parentOptions = computed<SelectOption[]>(() => {
  const options: SelectOption[] = [{ label: '根节点', value: 0 }]

  for (const menu of flatMenus.value) {
    if (menu.type === MenuType.Button || menu.id === formModel.id) {
      continue
    }

    options.push({
      label: `${'　'.repeat(menuLevel(menu.id))}${menu.title}`,
      value: menu.id,
    })
  }

  return options
})

const displayMenus = computed(() => {
  return filterMenus(menus.value)
})

const directoryCount = computed(() => {
  return flatMenus.value.filter((menu) => menu.type === MenuType.Directory).length
})

const menuCount = computed(() => {
  return flatMenus.value.filter((menu) => menu.type === MenuType.Menu).length
})

const buttonCount = computed(() => {
  return flatMenus.value.filter((menu) => menu.type === MenuType.Button).length
})

const columns: DataTableColumns<AdminMenu> = [
  {
    title: '菜单名称',
    key: 'title',
    minWidth: 240,
    render(row) {
      const typeConfig = {
        [MenuType.Directory]: { label: '目录', type: 'info' as const },
        [MenuType.Menu]: { label: '菜单', type: 'success' as const },
        [MenuType.Button]: { label: '按钮', type: 'warning' as const },
      }
      const cfg = typeConfig[row.type]

      return h('span', { class: 'inline-flex items-center gap-2' }, [
        h('span', { class: 'font-medium text-[#111827]' }, row.title),
        h(
          NTag,
          { size: 'small', bordered: false, round: false, type: cfg.type },
          { default: () => cfg.label },
        ),
      ])
    },
  },
  {
    title: '路由',
    key: 'path',
    minWidth: 130,
    ellipsis: { tooltip: true },
    render(row) {
      return row.path || '-'
    },
  },
  {
    title: '权限标识',
    key: 'code',
    minWidth: 150,
    ellipsis: { tooltip: true },
  },
  {
    title: '排序',
    key: 'sort',
    width: 64,
    align: 'center',
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    align: 'center',
    render(row) {
      return h(
        NTag,
        {
          size: 'small',
          type: row.status === MenuStatus.Enabled ? 'success' : 'error',
          bordered: false,
          round: true,
        },
        { default: () => (row.status === MenuStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    fixed: 'right',
    render(row) {
      const canCreateChild = row.type !== MenuType.Button && canUse('system:menu:create')
      const nextStatus =
        row.status === MenuStatus.Enabled ? MenuStatus.Disabled : MenuStatus.Enabled

      return h(
        NSpace,
        { size: 6, align: 'center' },
        {
          default: () =>
            [
              canCreateChild
                ? h(
                    NButton,
                    {
                      size: 'tiny',
                      type: 'primary',
                      secondary: true,
                      onClick: () => openCreateChild(row),
                    },
                    { default: () => (row.type === MenuType.Menu ? '+ 按钮' : '+ 子级') },
                  )
                : null,
              canUse('system:menu:update')
                ? h(
                    NButton,
                    {
                      size: 'tiny',
                      secondary: true,
                      onClick: () => openEdit(row),
                    },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:menu:status')
                ? h(
                    NTooltip,
                    {},
                    {
                      trigger: () =>
                        h(
                          NPopconfirm,
                          { onPositiveClick: () => handleToggleStatus(row, nextStatus) },
                          {
                            trigger: () =>
                              h(
                                NButton,
                                {
                                  size: 'tiny',
                                  type:
                                    nextStatus === MenuStatus.Disabled ? 'error' : 'success',
                                  secondary: true,
                                },
                                {
                                  default: () =>
                                    nextStatus === MenuStatus.Disabled ? '禁用' : '启用',
                                },
                              ),
                            default: () =>
                              `确认${nextStatus === MenuStatus.Disabled ? '禁用' : '启用'}该菜单？`,
                          },
                        ),
                      default: () => '切换菜单可见状态',
                    },
                  )
                : null,
              canUse('system:menu:delete')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleDelete(row) },
                    {
                      trigger: () =>
                        h(
                          NButton,
                          { size: 'tiny', type: 'error', secondary: true },
                          { default: () => '删除' },
                        ),
                      default: () => '删除前请确认它没有子菜单，也没有分配给任何角色。',
                    },
                  )
                : null,
            ].filter(Boolean),
        },
      )
    },
  },
]

function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}

function rowKey(row: AdminMenu) {
  return row.id
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

function flattenMenus(items: AdminMenu[]): AdminMenu[] {
  const result: AdminMenu[] = []

  for (const item of items) {
    result.push(item)
    result.push(...flattenMenus(item.children ?? []))
  }

  return result
}

function menuLevel(id: number) {
  let level = 0
  let current = flatMenus.value.find((menu) => menu.id === id)

  while (current && current.parent_id !== 0) {
    level += 1
    current = flatMenus.value.find((menu) => menu.id === current?.parent_id)
  }

  return level
}

function filterMenus(items: AdminMenu[]): AdminMenu[] {
  const keyword = query.keyword.trim().toLowerCase()
  const result: AdminMenu[] = []

  for (const item of items) {
    const children = filterMenus(item.children ?? [])
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

function resetForm() {
  Object.assign(formModel, {
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
  })
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
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-4 overflow-hidden">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">菜单管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">维护侧边栏目录、页面菜单和页面内按钮权限。</p>
        </div>

        <NButton v-if="canUse('system:menu:create')" type="primary" @click="openCreateRoot">
          + 新增根目录
        </NButton>
      </div>

      <NAlert
        v-if="successText"
        type="success"
        :show-icon="true"
        closable
        class="mx-auto w-full max-w-[520px]"
        @close="successText = ''"
      >
        {{ successText }}
      </NAlert>

      <NCard :bordered="false" class="rounded-lg">
          <NSpace align="center" :wrap="true">
            <NInput
              v-model:value="query.keyword"
              clearable
              placeholder="菜单名称 / 路由 / 权限标识"
              class="w-56"
            />
            <NSelect v-model:value="query.type" :options="typeOptions" class="w-40" />
            <NSelect v-model:value="query.status" :options="statusOptions" class="w-40" />
            <NButton @click="handleResetQuery">重置</NButton>
          </NSpace>
        </NCard>

      <NCard
        class="min-h-0 flex-1 rounded-lg"
        :bordered="false"
        content-style="height: 100%; padding: 0;"
      >
        <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-2.5">
          <span class="text-xs text-[#6B7280]">
            共 {{ flatMenus.length }} 个节点 · 目录 {{ directoryCount }} · 菜单 {{ menuCount }} · 按钮 {{ buttonCount }}
          </span>
          <NSpace :size="12">
            <NButton text size="small" @click="expandAll">展开全部</NButton>
            <NButton text size="small" @click="collapseAll">收起全部</NButton>
            <NButton text size="small" type="primary" @click="loadMenus">刷新</NButton>
          </NSpace>
        </div>

        <NDataTable
          class="menu-table"
          style="height: calc(100% - 48px)"
          :columns="columns"
          :data="displayMenus"
          :loading="loading"
          :row-key="rowKey"
          :expanded-row-keys="expandedRowKeys"
          :pagination="false"
          :bordered="false"
          children-key="children"
          @update:expanded-row-keys="handleExpandedChange"
          flex-height
        />
      </NCard>
    </section>

    <NModal
      v-model:show="formVisible"
      preset="card"
      :closable="false"
      class="menu-modal"
      style="width: 600px; max-width: calc(100vw - 32px)"
    >
      <template #header>
        <div class="modal-header">
          <h2>{{ formMode === 'create' ? '新增菜单' : '编辑菜单' }}</h2>
          <p>
            {{
              formMode === 'create'
                ? '选择节点类型后填写对应字段。'
                : '权限标识保持只读，避免影响按钮权限判断。'
            }}
          </p>
          <button type="button" class="modal-close" @click="formVisible = false">
            <NIcon :size="18">
              <CloseOutline />
            </NIcon>
          </button>
        </div>
      </template>

      <NForm
        ref="formRef"
        :model="formModel"
        :rules="rules"
        label-placement="left"
        label-width="80"
      >
        <NFormItem label="菜单类型" path="type">
          <div class="type-segment" :class="{ 'is-disabled': formMode === 'edit' }">
            <button
              v-for="opt in formTypeOptions"
              :key="opt.value"
              type="button"
              class="type-segment__btn"
              :class="{ 'type-segment__btn--active': formModel.type === opt.value }"
              :disabled="formMode === 'edit'"
              @click="formModel.type = opt.value as MenuType"
            >
              {{ opt.label }}
            </button>
          </div>
        </NFormItem>

        <NFormItem label="父级节点" path="parent_id">
          <NSelect
            v-model:value="formModel.parent_id"
            filterable
            :options="parentOptions"
          />
        </NFormItem>

        <NFormItem label="菜单名称" path="title">
          <NInput v-model:value="formModel.title" placeholder="请输入菜单名称" />
        </NFormItem>

        <NFormItem label="权限标识" path="code">
          <NInput
            v-model:value="formModel.code"
            placeholder="system:example:list"
            :disabled="formMode === 'edit'"
          />
        </NFormItem>

        <NFormItem v-if="formModel.type !== MenuType.Button" label="路由地址" path="path">
          <NInput v-model:value="formModel.path" placeholder="/system/example" />
        </NFormItem>

        <NFormItem
          v-if="formModel.type === MenuType.Menu"
          label="组件路径"
          path="component"
        >
          <NInput v-model:value="formModel.component" placeholder="system/UserView" />
        </NFormItem>

        <NFormItem label="图标 / 排序">
          <div class="grid w-full grid-cols-[1fr_120px] gap-2">
            <NInput v-model:value="formModel.icon" placeholder="setting / notification / layout-dashboard" />
            <NInputNumber v-model:value="formModel.sort" :min="0" />
          </div>
        </NFormItem>

        <NFormItem label="显示状态">
          <NSwitch
            :value="formModel.status === MenuStatus.Enabled"
            @update:value="
              (checked) => {
                formModel.status = checked ? MenuStatus.Enabled : MenuStatus.Disabled
              }
            "
          />
        </NFormItem>

        <NFormItem label="备注">
          <NInput
            v-model:value="formModel.remark"
            type="textarea"
            placeholder="请输入备注"
            :autosize="{ minRows: 3, maxRows: 5 }"
          />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="flex justify-end gap-3">
          <NButton @click="formVisible = false">取消</NButton>
          <NButton type="primary" :loading="saving" @click="handleSubmit">保存</NButton>
        </div>
      </template>
    </NModal>
  </main>
</template>

<style scoped>
.menu-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #374151;
  background: #fff;
}

.menu-table :deep(.n-data-table-td) {
  color: #374151;
}

.menu-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}

.menu-table :deep(.n-data-table-td .n-data-table-td__content) {
  display: inline-flex;
  align-items: center;
}

.menu-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 20px;
  border: 1px solid #dfe9f5;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.menu-modal :deep(.n-card-header) {
  padding: 0;
}

.modal-header {
  position: relative;
  padding: 24px 28px;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.modal-header h2 {
  font-size: 19px;
  font-weight: 700;
  color: #111827;
}

.modal-header p {
  margin-top: 8px;
  max-width: 420px;
  font-size: 13px;
  line-height: 1.6;
  color: #64748b;
}

.modal-close {
  position: absolute;
  top: 20px;
  right: 22px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.82);
  color: #64748b;
}

.type-segment {
  display: flex;
  gap: 4px;
  padding: 4px;
  border-radius: 6px;
  background: #f3f4f6;
}

.type-segment.is-disabled {
  opacity: 0.6;
  pointer-events: none;
}

.type-segment__btn {
  padding: 4px 20px;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: #6b7280;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}

.type-segment__btn--active {
  background: #fff;
  color: #18a058;
  font-weight: 600;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}
</style>
```

:::

::: details `admin/src/router/dynamic-menu.ts` — 动态路由映射

修改后，`system/RoleView` 和 `system/MenuView` 会从占位页切换为真实页面。

```ts
import {
  AlbumsOutline,
  AppsOutline,
  BeakerOutline,
  BuildOutline,
  CogOutline,
  DocumentTextOutline,
  FolderOpenOutline,
  GridOutline,
  LayersOutline,
  ListOutline,
  NotificationsOutline,
  PeopleOutline,
  PulseOutline,
  ServerOutline,
  SettingsOutline,
  ShieldCheckmarkOutline,
  TimeOutline,
} from '@vicons/ionicons5'
import { NIcon, type MenuOption } from 'naive-ui'
import type { RouteRecordRaw } from 'vue-router'
import { computed, h, shallowRef, type Component } from 'vue'

import { MenuType, type AuthMenu } from '../types/menu'

type RouteComponent = NonNullable<RouteRecordRaw['component']>
type MenuIconComponent = Component

const placeholderPage = () => import('../pages/system/PlaceholderPage.vue')

const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': () => import('../pages/system/HealthView.vue'),
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': () => import('../pages/system/RoleView.vue'),
  'system/MenuView': () => import('../pages/system/MenuView.vue'),
  'system/ConfigView': () => import('../pages/system/ConfigView.vue'),
  'system/FileView': () => import('../pages/system/FileView.vue'),
  'system/OperationLogView': () => import('../pages/system/OperationLogView.vue'),
  'system/LoginLogView': () => import('../pages/system/LoginLogView.vue'),
  'system/NoticeView': () => import('../pages/system/NoticeView.vue'),
}

const defaultMenuIcon = AppsOutline

// 后端 icon 字段只允许命中这份前端白名单，避免把任意字符串直接当组件渲染。
const menuIconMap: Record<string, MenuIconComponent> = {
  albums: AlbumsOutline,
  app: AppsOutline,
  apps: AppsOutline,
  beaker: BeakerOutline,
  blog: DocumentTextOutline,
  build: BuildOutline,
  cog: CogOutline,
  config: BuildOutline,
  dashboard: GridOutline,
  directory: AlbumsOutline,
  document: DocumentTextOutline,
  edit: DocumentTextOutline,
  experiment: BeakerOutline,
  file: FolderOpenOutline,
  files: FolderOpenOutline,
  folder: FolderOpenOutline,
  grid: GridOutline,
  health: PulseOutline,
  history: TimeOutline,
  home: GridOutline,
  layout: GridOutline,
  layoutdashboard: GridOutline,
  layers: LayersOutline,
  list: ListOutline,
  log: ListOutline,
  loginlog: TimeOutline,
  loginlogs: TimeOutline,
  logs: ListOutline,
  menu: LayersOutline,
  menus: LayersOutline,
  monitor: PulseOutline,
  notice: NotificationsOutline,
  notices: NotificationsOutline,
  notification: NotificationsOutline,
  notifications: NotificationsOutline,
  operationlog: ListOutline,
  operationlogs: ListOutline,
  page: DocumentTextOutline,
  people: PeopleOutline,
  person: PeopleOutline,
  role: ShieldCheckmarkOutline,
  roles: ShieldCheckmarkOutline,
  server: ServerOutline,
  setting: SettingsOutline,
  settings: SettingsOutline,
  shield: ShieldCheckmarkOutline,
  system: SettingsOutline,
  time: TimeOutline,
  user: PeopleOutline,
  users: PeopleOutline,
}

const builtinMenuOptions: MenuOption[] = [
  {
    label: '工作台',
    key: '/dashboard',
    icon: renderMenuIcon(GridOutline),
  },
]

export const authMenus = shallowRef<AuthMenu[]>([])

export const sideMenuOptions = computed<MenuOption[]>(() => {
  return [...builtinMenuOptions, ...buildMenuOptions(authMenus.value)]
})

export const buttonPermissionCodes = computed(() => {
  return collectButtonCodes(authMenus.value)
})

export function setAuthMenus(menus: AuthMenu[]) {
  authMenus.value = menus
}

export function clearAuthMenus() {
  authMenus.value = []
}

export function buildDynamicRoutes(menus: AuthMenu[]) {
  return collectPageMenus(menus).map<RouteRecordRaw>((menu) => ({
    path: toChildRoutePath(menu.path),
    name: `menu-${menu.id}`,
    component: resolveRouteComponent(menu.component),
    props: {
      title: menu.title,
      description: `${menu.title} 页面后续会接入真实业务。`,
    },
    meta: {
      title: menu.title,
      menuCode: menu.code,
    },
  }))
}

export function findMenuTitleByPath(path: string) {
  return collectPageMenus(authMenus.value).find((menu) => menu.path === path)?.title
}

function buildMenuOptions(menus: AuthMenu[]) {
  return menus.map(toMenuOption).filter(isMenuOption)
}

function toMenuOption(menu: AuthMenu): MenuOption | null {
  if (menu.type === MenuType.Button) {
    return null
  }

  const children = buildMenuOptions(menu.children ?? [])
  const key = menu.path || menu.code

  return {
    label: menu.title,
    key,
    icon: resolveMenuIcon(menu.icon),
    disabled: menu.type === MenuType.Directory && children.length === 0,
    children: children.length > 0 ? children : undefined,
  }
}

function isMenuOption(option: MenuOption | null): option is MenuOption {
  return option !== null
}

function collectPageMenus(menus: AuthMenu[]) {
  const result: AuthMenu[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Menu && menu.path) {
      result.push(menu)
    }

    result.push(...collectPageMenus(menu.children ?? []))
  }

  return result
}

function collectButtonCodes(menus: AuthMenu[]) {
  const result: string[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Button) {
      result.push(menu.code)
    }

    result.push(...collectButtonCodes(menu.children ?? []))
  }

  return result
}

function resolveRouteComponent(component: string) {
  return routeComponentMap[component] ?? placeholderPage
}

function resolveMenuIcon(icon: string) {
  return renderMenuIcon(menuIconMap[normalizeMenuIcon(icon)] ?? defaultMenuIcon)
}

function renderMenuIcon(icon: MenuIconComponent) {
  return () =>
    h(NIcon, null, {
      default: () => h(icon),
    })
}

function normalizeMenuIcon(icon: string) {
  return icon.trim().toLowerCase().replace(/[^a-z0-9]/g, '')
}

function toChildRoutePath(path: string) {
  return path.replace(/^\/+/, '')
}
```

:::

::: warning ⚠️ 按钮权限的 `code` 要和页面判断一致
例如用户页里判断的是 `system:user:create`、`system:user:update` 这类编码。菜单管理页新增按钮节点时，`code` 必须和页面代码里的 `canUse(code)` 保持一致，否则按钮权限不会生效。
:::

::: details 为什么创建菜单需要 `code`，编辑菜单不允许改 `code`
菜单编码会被按钮权限、角色菜单权限和前端权限判断使用。允许随意修改编码，很容易出现“页面还在，按钮突然不显示”的问题。

后端当前的编辑接口也没有接收 `code` 字段，所以前端编辑表单会把编码作为只读信息处理。
:::

## ✅ 验证结果

先启动后端和前端：

::: code-group

```bash [后端]
cd server
go run .
```

```bash [前端]
cd admin
pnpm dev
```

:::

然后按下面顺序验证：

1. 使用 `admin / EzAdmin@123456` 登录。
2. 点击“系统管理 / 角色管理”，确认角色列表和右侧权限树能正常加载。
3. 新建一个测试角色，例如 `demo_operator`，保存后左侧角色列表中能看到它。
4. 给测试角色分配菜单权限和按钮权限，点击“保存权限”。
5. 切换到“接口权限”，给测试角色分配必要接口权限。
6. 进入”系统管理 / 菜单管理”，点击”+ 新增根目录”打开弹框表单，填写后保存，确认树形表格刷新并展示新节点。
7. 回到“用户管理”，把测试用户绑定到这个角色。
8. 退出登录，再用测试用户登录。
9. 确认侧边栏只显示被授权的菜单，页面按钮也按按钮权限显示。

::: details 如果菜单没有变化，先检查这几件事
- 是否重新登录了。当前菜单在登录后加载，修改角色菜单后建议重新登录验证。
- `sys_role_menu` 是否写入了新的菜单 ID。
- `sys_menu.code` 是否和前端 `canUse(code)` 判断一致。
- 角色是否处于启用状态。
- 用户是否已经绑定到刚刚修改的角色。
:::

## 本节小结

这一节把权限体系最关键的两个页面补齐了：

- 角色页面负责维护角色、接口权限和菜单权限。
- 菜单页面负责维护目录、菜单和按钮权限。
- 动态路由通过 `component` 字段加载真实 Vue 页面。
- 菜单权限控制入口，接口权限控制后端访问，按钮权限控制页面操作体验。

下一节继续补齐系统管理里的剩余页面：[配置与文件页面](./config-file-pages)。
