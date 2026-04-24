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
