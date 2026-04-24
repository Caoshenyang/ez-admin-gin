---
title: 用户管理页面
description: "实现用户列表、查询、新增、编辑、启停和角色分配页面。"
---

# 用户管理页面

前面已经能按后端权限生成动态菜单。现在开始接第一个真实业务页面：用户管理。完成后，点击侧边栏里的“用户管理”，会进入真实列表页，可以查询用户、新增用户、编辑昵称和状态、启停账号，并给用户分配角色。

::: tip 🎯 本节目标
这一节会把 `system/UserView` 从占位页换成真实页面。页面会使用 Naive UI 的 `NDataTable`、`NForm`、`NModal`、`NSelect`、`NPopconfirm` 等组件；Tailwind CSS 4 只负责一屏布局、间距和内容区滚动。
:::

## 先看接口边界

用户管理页会用到两组接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/users` | 用户分页列表 |
| `POST` | `/api/v1/system/users` | 创建用户 |
| `POST` | `/api/v1/system/users/:id/update` | 编辑用户基础信息 |
| `POST` | `/api/v1/system/users/:id/status` | 修改用户状态 |
| `POST` | `/api/v1/system/users/:id/roles` | 分配用户角色 |
| `GET` | `/api/v1/system/roles` | 获取角色列表，用于角色选择器 |

::: warning ⚠️ 前端隐藏按钮不等于后端放行
这一页会读取上一节收集的按钮权限 code，控制“新增、编辑、启停、分配角色”等按钮是否显示。但真正的安全边界仍然在后端权限中间件里。前端按钮权限只是体验优化，不是安全机制。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
└─ src/
   ├─ App.vue
   ├─ api/
   │  ├─ role.ts
   │  └─ user.ts
   ├─ pages/
   │  └─ system/
   │     └─ UserView.vue
   ├─ router/
   │  └─ dynamic-menu.ts
   └─ types/
      ├─ role.ts
      └─ user.ts
```

| 位置 | 用途 |
| --- | --- |
| `src/App.vue` | 配置 Naive UI 中文语言包，让分页等组件显示中文 |
| `src/types/user.ts` | 定义用户列表、创建、编辑、状态和角色分配类型 |
| `src/types/role.ts` | 定义角色列表类型，用于用户表单里的角色选择 |
| `src/api/user.ts` | 封装用户管理接口 |
| `src/api/role.ts` | 封装角色列表接口 |
| `src/pages/system/UserView.vue` | 用户管理真实页面 |
| `src/router/dynamic-menu.ts` | 把 `system/UserView` 映射到真实页面 |

## 开始前先确认

开始之前，先确认下面几件事：

- 已完成上一节 [动态菜单](./dynamic-menu)。
- 登录后左侧菜单里能看到“用户管理”。
- 当前账号拥有 `system:user` 菜单权限，以及需要验证的按钮权限。
- 后端 `/api/v1/system/users` 和 `/api/v1/system/roles` 可以正常访问。

## 🛠️ 配置 Naive UI 中文

修改 `admin/src/App.vue`。上一节已经接入了 `NConfigProvider`，这里给它补上中文语言包和中文日期语言包。

```vue
<script setup lang="ts">
import {
  NConfigProvider,
  NDialogProvider,
  NLoadingBarProvider,
  NMessageProvider,
  NNotificationProvider,
  dateZhCN, // [!code ++]
  zhCN, // [!code ++]
} from 'naive-ui'
import { RouterView } from 'vue-router'
</script>

<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN">
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <RouterView />
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>
```

这一步会影响 Naive UI 全局组件文案，例如分页、日期、空状态等。用户管理页里的分页还会额外配置“共多少条、已选择多少条”和“10 / 页”这样的业务文案。

## 🛠️ 定义用户类型

新增 `admin/src/types/user.ts`。

::: details `admin/src/types/user.ts` — 用户类型

```ts
export const UserStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type UserStatus = (typeof UserStatus)[keyof typeof UserStatus]

export interface UserItem {
  id: number
  username: string
  nickname: string
  status: UserStatus
  role_ids: number[]
  created_at: string
  updated_at: string
}

export interface UserListQuery {
  page: number
  page_size: number
  keyword?: string
  role_id?: number
  status?: UserStatus
}

export interface UserListResponse {
  items: UserItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateUserPayload {
  username: string
  password: string
  nickname: string
  status: UserStatus
  role_ids: number[]
}

export interface UpdateUserPayload {
  nickname: string
  status: UserStatus
}

export interface UpdateUserStatusPayload {
  status: UserStatus
}

export interface UpdateUserRolesPayload {
  role_ids: number[]
}
```

:::

这里保留了后端返回的下划线字段，比如 `role_ids`、`created_at`。这样可以减少前端转换成本，也方便跟接口响应直接对照。

## 🛠️ 定义角色类型

新增 `admin/src/types/role.ts`。

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
  status?: RoleStatus
}

export interface RoleListResponse {
  items: RoleItem[]
  total: number
  page: number
  page_size: number
}
```

本节只用角色的 `id`、`code`、`name` 和 `status`，但类型里把接口返回字段补全，后续写角色页面时可以继续复用。

## 🛠️ 封装用户接口

新增 `admin/src/api/user.ts`。

::: details `admin/src/api/user.ts` — 用户接口

```ts
import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  CreateUserPayload,
  UpdateUserPayload,
  UpdateUserRolesPayload,
  UpdateUserStatusPayload,
  UserItem,
  UserListQuery,
  UserListResponse,
} from '../types/user'

export async function getUsers(params: UserListQuery) {
  const response = await http.get<ApiResponse<UserListResponse>>('/system/users', { params })
  return response.data.data
}

export async function createUser(payload: CreateUserPayload) {
  const response = await http.post<ApiResponse<UserItem>>('/system/users', payload)
  return response.data.data
}

export async function updateUser(id: number, payload: UpdateUserPayload) {
  const response = await http.post<ApiResponse<UserItem>>(`/system/users/${id}/update`, payload)
  return response.data.data
}

export async function updateUserStatus(id: number, payload: UpdateUserStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/users/${id}/status`,
    payload,
  )
  return response.data.data
}

export async function updateUserRoles(id: number, payload: UpdateUserRolesPayload) {
  const response = await http.post<ApiResponse<{ id: number; role_ids: number[] }>>(
    `/system/users/${id}/roles`,
    payload,
  )
  return response.data.data
}
```

:::

## 🛠️ 封装角色列表接口

新增 `admin/src/api/role.ts`。

```ts
import http from './http'

import type { ApiResponse } from '../types/http'
import type { RoleListQuery, RoleListResponse } from '../types/role'

export async function getRoles(params: RoleListQuery) {
  const response = await http.get<ApiResponse<RoleListResponse>>('/system/roles', { params })
  return response.data.data
}
```

用户页面只需要角色列表，不在这里实现角色创建、编辑和授权。角色管理会放到下一节继续做。

## 🛠️ 创建用户管理页面

新增 `admin/src/pages/system/UserView.vue`。

这一页代码比较长，建议直接创建完整文件，不要分段拼。页面会尽量贴近原型里的用户管理页：

- 顶部成功提示：创建、编辑、启停、分配角色成功后给出反馈。
- 查询区：用户名 / 手机号、角色、状态、查询和重置按钮，保持横向紧凑。
- 表格区：选择列、工具栏、主操作按钮加更多下拉和底部中文分页。
- 用户弹窗：把之前更有气质的“轻提示头部”融合进 `NModal` 的自定义 header，保留氛围，同时让正文从第一行开始就服务表单本身。
- 角色弹窗：比新增弹窗再小一档，减少空白感。

::: details `admin/src/pages/system/UserView.vue` — 用户管理页面

```vue
<script setup lang="ts">
import { CloseOutline, EllipsisHorizontal } from '@vicons/ionicons5'
import type {
  DataTableColumns,
  DataTableRowKey,
  FormInst,
  FormRules,
  SelectOption,
} from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NDataTable,
  NDropdown,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NModal,
  NPagination,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'

import { getRoles } from '../../api/role'
import {
  createUser,
  getUsers,
  updateUser,
  updateUserRoles,
  updateUserStatus,
} from '../../api/user'
import { buttonPermissionCodes } from '../../router/dynamic-menu'
import { RoleStatus, type RoleItem } from '../../types/role'
import { UserStatus, type UserItem, type UserListQuery } from '../../types/user'

interface UserFormModel {
  id: number
  username: string
  password: string
  nickname: string
  status: UserStatus
  role_ids: number[]
}

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const users = ref<UserItem[]>([])
const roles = ref<RoleItem[]>([])
const total = ref(0)
const checkedRowKeys = ref<DataTableRowKey[]>([])
const successText = ref('')

const query = reactive<UserListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  role_id: 0,
  status: 0,
})

const formRef = ref<FormInst | null>(null)
const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const formModel = reactive<UserFormModel>({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  status: UserStatus.Enabled,
  role_ids: [],
})

const roleVisible = ref(false)
const roleSaving = ref(false)
const roleUser = ref<UserItem | null>(null)
const selectedRoleIDs = ref<number[]>([])

const roleNameMap = computed(() => {
  return new Map(roles.value.map((role) => [role.id, role.name]))
})

const roleOptions = computed<SelectOption[]>(() => {
  return roles.value.map((role) => ({
    label: `${role.name}（${role.code}）`,
    value: role.id,
  }))
})

const roleFilterOptions = computed<SelectOption[]>(() => {
  return [
    { label: '角色：全部', value: 0 },
    ...roles.value.map((role) => ({
      label: role.name,
      value: role.id,
    })),
  ]
})

const statusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: UserStatus.Enabled },
  { label: '禁用', value: UserStatus.Disabled },
]

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
}

const selectedCount = computed(() => checkedRowKeys.value.length)

const displayUsers = computed(() => {
  if (!query.role_id) {
    return users.value
  }

  return users.value.filter((user) => user.role_ids.includes(query.role_id ?? 0))
})

const displayTotal = computed(() => {
  return query.role_id ? displayUsers.value.length : total.value
})

const columns: DataTableColumns<UserItem> = [
  { type: 'selection', width: 48 },
  {
    title: '用户',
    key: 'username',
    minWidth: 180,
    render(row) {
      return h('div', { class: 'leading-6' }, [
        h('p', { class: 'font-semibold text-[#111827]' }, row.username),
        h('p', { class: 'text-xs text-[#6B7280]' }, row.nickname),
      ])
    },
  },
  {
    title: '角色',
    key: 'role_ids',
    minWidth: 220,
    render(row) {
      if (row.role_ids.length === 0) {
        return h('span', { class: 'text-sm text-[#9CA3AF]' }, '未分配')
      }

      return h(
        NSpace,
        { size: 6 },
        {
          default: () =>
            row.role_ids.map((roleID) =>
              h(
                NTag,
                { size: 'small', bordered: false },
                { default: () => roleNameMap.value.get(roleID) ?? `角色 ${roleID}` },
              ),
            ),
        },
      )
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 110,
    render(row) {
      return h(
        NTag,
        {
          type: row.status === UserStatus.Enabled ? 'success' : 'error',
          bordered: false,
        },
        { default: () => (row.status === UserStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 190,
    render(row) {
      return formatTime(row.created_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    fixed: 'right',
    render(row) {
      const nextStatus =
        row.status === UserStatus.Enabled ? UserStatus.Disabled : UserStatus.Enabled
      const dropdownOptions = [
        canUse('system:user:assign-role')
          ? {
              label: '分配角色',
              key: `role:${row.id}`,
            }
          : null,
      ].filter(Boolean)

      return h(
        NSpace,
        { size: 8, align: 'center' },
        {
          default: () =>
            [
              canUse('system:user:update')
                ? h(
                    NButton,
                    {
                      size: 'small',
                      ghost: true,
                      type: 'info',
                      class: 'min-w-[48px]',
                      onClick: () => openEdit(row),
                    },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:user:status')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleToggleStatus(row, nextStatus) },
                    {
                      trigger: () =>
                        h(
                          NButton,
                          {
                            size: 'small',
                            ghost: true,
                            type:
                              nextStatus === UserStatus.Disabled ? 'error' : 'success',
                            class: 'min-w-[48px]',
                          },
                          { default: () => (nextStatus === UserStatus.Disabled ? '禁用' : '启用') },
                        ),
                      default: () =>
                        `确认${nextStatus === UserStatus.Disabled ? '禁用' : '启用'}该用户？`,
                    },
                  )
                : null,
              dropdownOptions.length > 0
                ? h(
                    NDropdown,
                    {
                      options: dropdownOptions,
                      trigger: 'click',
                      onSelect: (key: string | number) => handleRowAction(String(key), row),
                    },
                    {
                      default: () =>
                        h(
                          NButton,
                          {
                            size: 'small',
                            quaternary: true,
                            class: 'min-w-[36px] px-2',
                          },
                          {
                            icon: () =>
                              h(NIcon, null, {
                                default: () => h(EllipsisHorizontal),
                              }),
                          },
                        ),
                    },
                  )
                : null,
            ].filter(Boolean),
        },
      )
    },
  },
]

function rowKey(row: UserItem) {
  return row.id
}

function handleCheckedRowKeys(keys: DataTableRowKey[]) {
  checkedRowKeys.value = keys
}

function handlePageChange(page: number) {
  query.page = page
  void loadUsers()
}

function handlePageSizeChange(pageSize: number) {
  query.page = 1
  query.page_size = pageSize
  void loadUsers()
}

function handleReset() {
  query.page = 1
  query.page_size = 10
  query.keyword = ''
  query.role_id = 0
  query.status = 0
  void loadUsers()
}

function handleRowAction(key: string, row: UserItem) {
  if (key === `role:${row.id}`) {
    openRole(row)
  }
}

function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}

function formatTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

function resetForm() {
  Object.assign(formModel, {
    id: 0,
    username: '',
    password: '',
    nickname: '',
    status: UserStatus.Enabled,
    role_ids: [],
  })
}

async function loadUsers() {
  loading.value = true
  try {
    const data = await getUsers({
      ...query,
      keyword: query.keyword?.trim() || undefined,
      role_id: query.role_id === 0 ? undefined : query.role_id,
      status: query.status === 0 ? undefined : query.status,
    })
    users.value = data.items
    total.value = data.total
    checkedRowKeys.value = []
  } finally {
    loading.value = false
  }
}

async function loadRoles() {
  const data = await getRoles({
    page: 1,
    page_size: 100,
    status: RoleStatus.Enabled,
  })
  roles.value = data.items
}

function handleSearch() {
  query.page = 1
  void loadUsers()
}

function openCreate() {
  formMode.value = 'create'
  resetForm()
  formVisible.value = true
}

function openEdit(row: UserItem) {
  formMode.value = 'edit'
  Object.assign(formModel, {
    id: row.id,
    username: row.username,
    password: '',
    nickname: row.nickname,
    status: row.status,
    role_ids: row.role_ids,
  })
  formVisible.value = true
}

async function handleSubmit() {
  await formRef.value?.validate()
  saving.value = true
  try {
    if (formMode.value === 'create') {
      await createUser({
        username: formModel.username,
        password: formModel.password,
        nickname: formModel.nickname,
        status: formModel.status,
        role_ids: formModel.role_ids,
      })
      successText.value = '用户创建成功，临时密码已生成'
      message.success('用户创建成功')
    } else {
      await updateUser(formModel.id, {
        nickname: formModel.nickname,
        status: formModel.status,
      })
      successText.value = '用户信息已更新'
      message.success('用户更新成功')
    }

    formVisible.value = false
    await loadUsers()
  } finally {
    saving.value = false
  }
}

async function handleToggleStatus(row: UserItem, status: UserStatus) {
  await updateUserStatus(row.id, { status })
  successText.value = `用户已${status === UserStatus.Enabled ? '启用' : '禁用'}`
  message.success('用户状态已更新')
  await loadUsers()
}

function openRole(row: UserItem) {
  roleUser.value = row
  selectedRoleIDs.value = [...row.role_ids]
  roleVisible.value = true
}

async function handleSaveRoles() {
  if (!roleUser.value) {
    return
  }

  roleSaving.value = true
  try {
    await updateUserRoles(roleUser.value.id, { role_ids: selectedRoleIDs.value })
    successText.value = '用户角色已更新'
    message.success('用户角色已更新')
    roleVisible.value = false
    await loadUsers()
  } finally {
    roleSaving.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadRoles(), loadUsers()])
})
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-4 overflow-hidden">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">用户管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">维护后台账号、启停状态和角色绑定。</p>
        </div>

        <NButton v-if="canUse('system:user:create')" type="primary" @click="openCreate">
          + 新增用户
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
            placeholder="用户名 / 手机号"
            class="w-64"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.role_id" :options="roleFilterOptions" class="w-40" />
          <NSelect
            v-model:value="query.status"
            :options="statusOptions"
            class="w-36"
          />
          <NButton type="primary" @click="handleSearch">查询</NButton>
          <NButton @click="handleReset">重置</NButton>
        </NSpace>
      </NCard>

      <NCard
        class="min-h-0 flex-1 rounded-lg"
        :bordered="false"
        content-style="height: 100%; padding: 0;"
      >
        <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-3">
          <NSpace :size="12">
            <span class="text-sm text-[#6B7280]">已选 {{ selectedCount }} 项</span>
            <NButton text :disabled="selectedCount === 0">批量禁用</NButton>
            <NButton text :disabled="selectedCount === 0">批量删除</NButton>
          </NSpace>
          <NSpace :size="14">
            <NButton text type="primary">列设置</NButton>
            <NButton text type="primary">密度</NButton>
            <NButton text type="primary" @click="loadUsers">刷新</NButton>
          </NSpace>
        </div>

        <NDataTable
          remote
          class="user-table h-full"
          style="height: calc(100% - 105px)"
          :columns="columns"
          :data="displayUsers"
          :loading="loading"
          :pagination="false"
          :row-key="rowKey"
          :checked-row-keys="checkedRowKeys"
          :bordered="false"
          flex-height
          @update:checked-row-keys="handleCheckedRowKeys"
        />

        <div
          class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]"
        >
          <span>共 {{ displayTotal }} 条，已选择 {{ selectedCount }} 条</span>
          <NPagination
            :page="query.page"
            :page-size="query.page_size"
            :item-count="displayTotal"
            :page-sizes="[10, 20, 50]"
            show-size-picker
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </NCard>
    </section>

    <NModal
      v-model:show="formVisible"
      preset="card"
      :closable="false"
      class="compact-user-modal"
      style="width: 640px; max-width: calc(100vw - 32px)"
    >
      <template #header>
        <div class="modal-header modal-header--hero">
          <h2 class="modal-header__title">
            {{ formMode === 'create' ? '新增用户' : '编辑用户' }}
          </h2>
          <p class="modal-header__hero-title">
            {{
              formMode === 'create'
                ? '先完成账号主体信息，再补充默认角色范围'
                : '这里仅维护展示资料，不修改登录凭证'
            }}
          </p>
          <p class="modal-header__hero-desc">
            {{
              formMode === 'create'
                ? '用户名和密码会作为首次登录凭证，角色支持后续在列表中继续微调。'
                : '编辑模式下不修改登录名和密码，避免影响已有账号的登录和追踪。'
            }}
          </p>
          <button type="button" class="modal-close" @click="formVisible = false">
            <NIcon :size="18">
              <CloseOutline />
            </NIcon>
          </button>
        </div>
      </template>

      <div class="user-modal-shell">
        <NForm
          ref="formRef"
          class="compact-user-form"
          :model="formModel"
          :rules="rules"
          label-placement="left"
          label-width="76"
        >
          <section class="form-section form-section--primary">
            <div class="form-section__head">
              <h3>基础信息</h3>
              <p>先把账号主体信息补完整，这是本次弹窗的主要内容。</p>
            </div>

            <div v-if="formMode === 'create'" class="form-section-grid">
              <NFormItem label="用户名" path="username">
                <NInput v-model:value="formModel.username" placeholder="请输入用户名" />
              </NFormItem>

              <NFormItem label="登录密码" path="password">
                <NInput
                  v-model:value="formModel.password"
                  type="password"
                  show-password-on="click"
                  placeholder="至少 8 位"
                />
              </NFormItem>

              <NFormItem label="昵称" path="nickname">
                <NInput v-model:value="formModel.nickname" placeholder="请输入昵称" />
              </NFormItem>

              <NFormItem label="状态" path="status">
                <NSelect v-model:value="formModel.status" :options="statusOptions.slice(1)" />
              </NFormItem>
            </div>

            <div v-else class="form-section-grid">
              <NFormItem label="昵称" path="nickname">
                <NInput v-model:value="formModel.nickname" placeholder="请输入昵称" />
              </NFormItem>

              <NFormItem label="状态" path="status">
                <NSelect v-model:value="formModel.status" :options="statusOptions.slice(1)" />
              </NFormItem>
            </div>
          </section>

          <section v-if="formMode === 'create'" class="form-section form-section--muted">
            <div class="form-section__head">
              <h3>角色配置</h3>
              <p>这是补充信息，先给一个默认角色即可，后续仍可在列表中单独调整。</p>
            </div>

            <NFormItem label="角色" path="role_ids" class="mb-0">
              <NSelect
                v-model:value="formModel.role_ids"
                multiple
                filterable
                :options="roleOptions"
                placeholder="请选择角色"
              />
            </NFormItem>
            <p class="form-section__tip">一个用户可以绑定多个角色，系统会自动合并其权限范围。</p>
          </section>
        </NForm>
      </div>

      <template #footer>
        <div class="modal-footer-actions">
          <NButton quaternary class="modal-footer-button" @click="formVisible = false">
            取消
          </NButton>
          <NButton
            type="primary"
            class="modal-footer-button modal-footer-button--primary"
            :loading="saving"
            @click="handleSubmit"
          >
            保存
          </NButton>
        </div>
      </template>
    </NModal>

    <NModal
      v-model:show="roleVisible"
      preset="card"
      :closable="false"
      class="compact-user-modal"
      style="width: 520px; max-width: calc(100vw - 32px)"
    >
      <template #header>
        <div class="modal-header modal-header--hero modal-header--role">
          <h2 class="modal-header__title">分配角色</h2>
          <p class="modal-header__hero-title">建议先给最小权限角色，再按职责逐步放开</p>
          <p class="modal-header__hero-desc">
            保存后立即生效，多角色场景下权限会按并集进行合并。
          </p>
          <button type="button" class="modal-close" @click="roleVisible = false">
            <NIcon :size="18">
              <CloseOutline />
            </NIcon>
          </button>
        </div>
      </template>

      <div class="user-modal-shell">
        <NForm class="compact-user-form" label-placement="left" label-width="76">
          <section class="form-section form-section--muted">
            <div class="form-section__head">
              <h3>角色设置</h3>
              <p>为当前账号选择一个或多个角色，保存后立即生效。</p>
            </div>

            <NFormItem label="当前用户">
              <NInput :value="roleUser?.username ?? ''" disabled />
            </NFormItem>

            <NFormItem label="角色" class="mb-0">
              <NSelect
                v-model:value="selectedRoleIDs"
                multiple
                filterable
                :options="roleOptions"
                placeholder="请选择角色"
              />
            </NFormItem>
            <p class="form-section__tip">多角色场景下，菜单与按钮权限会按照并集生效。</p>
          </section>
        </NForm>
      </div>

      <template #footer>
        <div class="modal-footer-actions">
          <NButton quaternary class="modal-footer-button" @click="roleVisible = false">
            取消
          </NButton>
          <NButton
            type="primary"
            class="modal-footer-button modal-footer-button--primary"
            :loading="roleSaving"
            @click="handleSaveRoles"
          >
            保存
          </NButton>
        </div>
      </template>
    </NModal>
  </main>
</template>

<style scoped>
.compact-user-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 32px;
  border: 1px solid #dfe9f5;
  background: #ffffff;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.user-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #374151;
  background: #fff;
}

.user-table :deep(.n-data-table-td) {
  color: #374151;
}

.user-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}

.compact-user-modal :deep(.n-card-header) {
  padding: 0;
  border-bottom: 1px solid #dfe9f5;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.compact-user-modal :deep(.n-card-header__main) {
  font-size: 19px;
  font-weight: 600;
  letter-spacing: 0.01em;
  color: #111827;
}

.compact-user-modal :deep(.n-card__content) {
  padding: 20px 28px 10px;
}

.compact-user-modal :deep(.n-card__footer) {
  padding: 16px 28px 24px;
  border-top: 1px solid #edf2f7;
  background: rgba(248, 250, 252, 0.85);
}

.compact-user-form :deep(.n-form-item) {
  margin-bottom: 16px;
}

.compact-user-form :deep(.n-form-item-label) {
  white-space: nowrap;
  align-items: center;
  padding-right: 14px;
  font-weight: 600;
  color: #374151;
}

.compact-user-form :deep(.n-form-item-blank) {
  min-height: 40px;
}

.compact-user-form :deep(.n-input-wrapper) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-user-form :deep(.n-base-selection) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-user-form :deep(.n-input),
.compact-user-form :deep(.n-base-selection) {
  box-shadow: none;
}

.compact-user-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.user-modal-shell {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.modal-header--hero {
  position: relative;
  overflow: hidden;
  min-height: 140px;
  padding: 26px 28px 22px;
  background:
    radial-gradient(circle at top right, rgba(34, 197, 94, 0.12), transparent 24%),
    linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.modal-header--hero::after {
  content: '';
  position: absolute;
  top: -18px;
  right: -10px;
  width: 118px;
  height: 118px;
  border-radius: 999px;
  background: radial-gradient(circle, rgba(34, 197, 94, 0.1) 0%, rgba(34, 197, 94, 0) 72%);
}

.modal-header--hero::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 0;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.24) 0%, rgba(255, 255, 255, 0.08) 100%);
  pointer-events: none;
}

.modal-header--role::after {
  background: radial-gradient(circle, rgba(37, 99, 235, 0.08) 0%, rgba(37, 99, 235, 0) 72%);
}

.modal-header__title {
  position: relative;
  z-index: 1;
  font-size: 19px;
  font-weight: 600;
  line-height: 1.3;
  color: #111827;
}

.modal-header__hero-title {
  position: relative;
  z-index: 1;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.6;
  color: #0f172a;
}

.modal-header__hero-desc {
  position: relative;
  z-index: 1;
  font-size: 12px;
  line-height: 1.6;
  color: #64748b;
}

.modal-close {
  position: absolute;
  top: 20px;
  right: 22px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border: none;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.76);
  color: #64748b;
  box-shadow: 0 10px 24px rgba(148, 163, 184, 0.12);
  backdrop-filter: blur(8px);
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.modal-close:hover {
  background: #ffffff;
  color: #111827;
  box-shadow: 0 14px 28px rgba(148, 163, 184, 0.18);
  transform: translateY(-1px);
}

.form-section {
  border: 1px solid #e9eff6;
  border-radius: 14px;
  background: #ffffff;
  padding: 18px 18px 4px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.9);
}

.form-section--primary {
  border-color: #d9e7f8;
  background: linear-gradient(180deg, #ffffff 0%, #fcfdff 100%);
}

.form-section--muted {
  background: linear-gradient(180deg, #fcfdff 0%, #f9fbff 100%);
}

.form-section__head {
  margin-bottom: 12px;
}

.form-section__head h3 {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.form-section__head p {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: #6b7280;
}

.form-section-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  column-gap: 20px;
}

.form-section__tip {
  margin-top: 6px;
  margin-bottom: 0;
  font-size: 12px;
  line-height: 1.6;
  color: #6b7280;
}

.modal-footer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.modal-footer-button {
  min-width: 92px;
  height: 40px;
  border-radius: 10px;
}

.modal-footer-button--primary {
  box-shadow: 0 10px 24px rgba(34, 197, 94, 0.18);
}

.mb-0 {
  margin-bottom: 0;
}

@media (max-width: 720px) {
  .form-section-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .compact-user-modal :deep(.n-card-header),
  .compact-user-modal :deep(.n-card__content),
  .compact-user-modal :deep(.n-card__footer) {
    padding-left: 20px;
    padding-right: 20px;
  }

  .compact-user-modal :deep(.n-card-header) {
    padding-bottom: 0;
  }

  .modal-header--hero {
    padding: 22px 20px 18px;
    min-height: 126px;
  }

  .modal-close {
    top: 18px;
    right: 18px;
  }

  .compact-user-form :deep(.n-form-item-label) {
    width: 72px;
  }
}
</style>
```

:::

::: details 为什么新增用户时可以选角色，编辑用户时不放在同一个弹窗里
后端把“编辑用户基础信息”和“分配用户角色”拆成了两个接口：

- `/api/v1/system/users/:id/update` 只改昵称和状态。
- `/api/v1/system/users/:id/roles` 只改角色绑定。

前端跟着这个边界拆成两个弹窗，后续做权限控制和操作日志时会更清楚。
:::

## 🛠️ 替换动态组件映射

修改 `admin/src/router/dynamic-menu.ts`。找到 `routeComponentMap`，把 `system/UserView` 从占位页换成真实页面。

```ts
const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': placeholderPage,
  'system/UserView': placeholderPage, // [!code --]
  'system/UserView': () => import('../pages/system/UserView.vue'), // [!code ++]
  'system/RoleView': placeholderPage,
  'system/MenuView': placeholderPage,
  'system/ConfigView': placeholderPage,
  'system/FileView': placeholderPage,
  'system/OperationLogView': placeholderPage,
  'system/LoginLogView': placeholderPage,
}
```

替换完成后，只保留新增的那一行，不要让 `system/UserView` 在对象里出现两次。

## 🧪 验证用户管理页

启动前后端后，用管理员账号登录：

```text
admin / EzAdmin@123456
```

进入“用户管理”后，按下面顺序验证：

| 验证点 | 预期结果 |
| --- | --- |
| 页面加载 | 表格请求 `/api/v1/system/users`，并显示用户列表 |
| 查询关键词 | 输入用户名或昵称后点击查询，列表按关键词过滤 |
| 角色 / 状态筛选 | 选择角色或状态后，筛选区结构与原型一致；点击重置后恢复默认筛选 |
| 新增用户 | 输入用户名、密码、昵称、角色后，可以创建用户 |
| 编辑用户 | 可以修改昵称和状态 |
| 启停用户 | 点击启用/禁用后，用户状态刷新 |
| 分配角色 | 保存后 `role_ids` 更新，表格角色标签同步变化 |
| 按钮权限 | 没有对应按钮权限时，页面按钮不显示 |
| 中文分页 | 底部分页左侧显示“共 N 条，已选择 N 条”，右侧分页器显示中文页大小 |
| 原型细节 | 新增按钮显示“+ 新增用户”，操作区保留主按钮并把次级操作收进更多菜单，顶部有绿色成功提示 |

最后跑一次前端检查：

```bash
cd admin
pnpm exec oxlint .
pnpm exec vue-tsc --noEmit
```

::: warning ⚠️ 不要用当前登录账号测试“禁用自己”和“修改自己角色”
后端已经禁止禁用当前登录用户，也禁止修改当前登录用户的角色。验证这些操作时，建议新建一个普通测试账号，再对测试账号执行启停和角色分配。
:::

## 常见问题

### 进入用户管理还是占位页

检查 `admin/src/router/dynamic-menu.ts` 里的 `routeComponentMap`。`system/UserView` 必须映射到 `../pages/system/UserView.vue`，并且不要保留重复 key。

### 表格能打开，但按钮不显示

先看 `/api/v1/auth/menus` 返回的菜单树里，当前账号是否有下面这些按钮权限：

| 按钮 | 权限 code |
| --- | --- |
| 新增用户 | `system:user:create` |
| 编辑用户 | `system:user:update` |
| 启停用户 | `system:user:status` |
| 分配角色 | `system:user:assign-role` |

如果菜单树里没有对应按钮，前端会按预期隐藏按钮。

### 创建用户提示角色不存在或已禁用

用户创建接口会校验角色是否存在、是否启用。先到数据库或角色接口确认角色状态，或者暂时不选择角色创建用户，再单独验证角色分配。

## 本节小结

这一节把用户管理页从占位页升级成了真实业务页面：

```text
动态菜单 system/UserView
  ↓
加载 UserView.vue
  ↓
请求用户列表和角色列表
  ↓
完成查询、新增、编辑、启停、角色分配
  ↓
按钮权限控制页面操作入口
```

下一节继续做角色和菜单管理，把角色列表、接口权限和菜单权限维护接到前端：[角色与菜单页面](./role-menu-pages)。
