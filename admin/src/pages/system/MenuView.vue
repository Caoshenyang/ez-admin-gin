<script setup lang="ts">
import type { DataTableColumns, FormInst, FormRules, SelectOption } from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NPopconfirm,
  NSelect,
  NSpace,
  NSwitch,
  NTag,
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
const panelMode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInst | null>(null)

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
    minWidth: 180,
    render(row) {
      return h('div', { class: 'flex items-center gap-2 font-semibold text-[#111827]' }, [
        h('span', typeIcon(row.type)),
        h('span', row.title),
      ])
    },
  },
  {
    title: '路由',
    key: 'path',
    minWidth: 150,
    render(row) {
      return row.path || '-'
    },
  },
  {
    title: '组件',
    key: 'component',
    minWidth: 150,
    render(row) {
      return row.component || '-'
    },
  },
  {
    title: '权限标识',
    key: 'code',
    minWidth: 170,
  },
  {
    title: '序',
    key: 'sort',
    width: 70,
  },
  {
    title: '状态',
    key: 'status',
    width: 96,
    render(row) {
      return h(
        NTag,
        {
          type: row.status === MenuStatus.Enabled ? 'success' : 'error',
          bordered: false,
        },
        { default: () => (row.status === MenuStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 248,
    fixed: 'right',
    render(row) {
      const canCreateChild = row.type !== MenuType.Button && canUse('system:menu:create')
      const nextStatus =
        row.status === MenuStatus.Enabled ? MenuStatus.Disabled : MenuStatus.Enabled

      return h(
        NSpace,
        { size: 8 },
        {
          default: () =>
            [
              canCreateChild
                ? h(
                    NButton,
                    {
                      size: 'small',
                      ghost: true,
                      type: 'success',
                      onClick: () => openCreateChild(row),
                    },
                    { default: () => (row.type === MenuType.Menu ? '按钮' : '新增') },
                  )
                : null,
              canUse('system:menu:update')
                ? h(
                    NButton,
                    {
                      size: 'small',
                      ghost: true,
                      type: 'info',
                      onClick: () => openEdit(row),
                    },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:menu:status')
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
                              nextStatus === MenuStatus.Disabled ? 'error' : 'success',
                          },
                          { default: () => (nextStatus === MenuStatus.Disabled ? '禁用' : '启用') },
                        ),
                      default: () =>
                        `确认${nextStatus === MenuStatus.Disabled ? '禁用' : '启用'}该菜单？`,
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
                          { size: 'small', ghost: true, type: 'error' },
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

function typeIcon(type: MenuType) {
  if (type === MenuType.Directory) {
    return '▾'
  }
  if (type === MenuType.Menu) {
    return '◻'
  }
  return '•'
}

function typeName(type: MenuType) {
  if (type === MenuType.Directory) {
    return '目录'
  }
  if (type === MenuType.Menu) {
    return '菜单'
  }
  return '按钮'
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
  } finally {
    loading.value = false
  }
}

function openCreateRoot() {
  panelMode.value = 'create'
  resetForm()
}

function openCreateChild(row: AdminMenu) {
  panelMode.value = 'create'
  resetForm()
  formModel.parent_id = row.id
  formModel.type = row.type === MenuType.Directory ? MenuType.Menu : MenuType.Button
  formModel.sort = row.type === MenuType.Directory ? 1 : 10
}

function openEdit(row: AdminMenu) {
  panelMode.value = 'edit'
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
}

async function handleSubmit() {
  await formRef.value?.validate()
  saving.value = true
  try {
    const payload = normalizedPayload()

    if (panelMode.value === 'create') {
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
    if (panelMode.value === 'create') {
      resetForm()
    }
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
    openCreateRoot()
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

      <div class="grid min-h-0 flex-1 grid-cols-[minmax(0,1fr)_360px] gap-4 overflow-hidden">
        <section class="flex min-h-0 flex-col gap-4 overflow-hidden">
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
            <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-3">
              <span class="text-sm text-[#374151]">
                菜单树 ｜ 目录 {{ directoryCount }} 个，菜单 {{ menuCount }} 个，按钮权限
                {{ buttonCount }} 个
              </span>
              <NSpace :size="14">
                <NButton text type="primary">展开全部</NButton>
                <NButton text type="primary">收起全部</NButton>
                <NButton text type="primary" @click="loadMenus">刷新</NButton>
              </NSpace>
            </div>

            <NDataTable
              class="menu-table"
              style="height: calc(100% - 103px)"
              :columns="columns"
              :data="displayMenus"
              :loading="loading"
              :row-key="rowKey"
              :pagination="false"
              :bordered="false"
              children-key="children"
              default-expand-all
              flex-height
            />

            <div class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3">
              <span class="text-sm text-[#6B7280]">
                共 {{ flatMenus.length }} 个菜单节点，目录 {{ directoryCount }} 个，按钮权限
                {{ buttonCount }} 个
              </span>
              <NSpace>
                <NButton>导入</NButton>
                <NButton>导出</NButton>
                <NButton>同步路由</NButton>
              </NSpace>
            </div>
          </NCard>
        </section>

        <NCard class="rounded-lg" :bordered="false" content-style="height: 100%;">
          <div class="flex h-full flex-col">
            <div class="mb-4">
              <h2 class="text-xl font-bold text-[#111827]">
                {{ panelMode === 'create' ? '新增菜单' : '编辑菜单' }}
              </h2>
              <p class="mt-1 text-sm text-[#6B7280]">
                {{ panelMode === 'create' ? '选择节点类型后填写对应字段。' : '权限标识保持只读，避免影响按钮权限判断。' }}
              </p>
            </div>

            <NForm
              ref="formRef"
              class="min-h-0 flex-1 overflow-y-auto pr-1"
              :model="formModel"
              :rules="rules"
              label-placement="top"
            >
              <NFormItem label="菜单类型" path="type">
                <NSelect
                  v-model:value="formModel.type"
                  :options="formTypeOptions"
                  :disabled="panelMode === 'edit'"
                />
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
                  :disabled="panelMode === 'edit'"
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
                  <NInput v-model:value="formModel.icon" placeholder="layout-dashboard" />
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

            <div class="mt-4 flex justify-end gap-3 border-t border-[#E5E7EB] pt-4">
              <NButton @click="openCreateRoot">取消</NButton>
              <NButton type="primary" :loading="saving" @click="handleSubmit">确认</NButton>
            </div>
          </div>
        </NCard>
      </div>
    </section>
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
</style>
