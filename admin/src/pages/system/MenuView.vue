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
