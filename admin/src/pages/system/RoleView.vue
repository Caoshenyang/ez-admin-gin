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
