<script setup lang="ts">
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import { NDataTable, NPopconfirm, NSpace, NTag } from 'naive-ui'
import { h } from 'vue'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import EzDataTable from '@/components/ez/EzDataTable.vue'
import { displayText } from '@/utils/format'
import { MenuStatus, MenuType, type AdminMenu } from '@/modules/iam/types/menu'

const props = defineProps<{
  canUse: (code: string) => boolean
  checkedRowKeys: DataTableRowKey[]
  displayMenus: AdminMenu[]
  flatMenuCount: number
  loading: boolean
  selectedCount: number
  stats: { directoryCount: number; menuCount: number; buttonCount: number }
}>()

const expandedRowKeys = defineModel<Array<string | number>>('expandedRowKeys', { required: true })

const emit = defineEmits<{
  checkedRowKeysChange: [keys: DataTableRowKey[]]
  collapseAll: []
  createChild: [row: AdminMenu]
  delete: [row: AdminMenu]
  deleteSelected: []
  edit: [row: AdminMenu]
  expandAll: []
  refresh: []
  toggleStatus: [row: AdminMenu, status: MenuStatus]
}>()

const menuTypeConfig = {
  [MenuType.Directory]: { label: '目录', type: 'info' as const },
  [MenuType.Menu]: { label: '菜单', type: 'success' as const },
  [MenuType.Button]: { label: '按钮', type: 'warning' as const },
}

function formatTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

const columns: DataTableColumns<AdminMenu> = [
  { type: 'selection', width: 44 },
  {
    title: '菜单名称',
    key: 'title',
    minWidth: 180,
    render(row) {
      return h('span', { class: 'font-medium text-[var(--ez-text-main)]' }, displayText(row.title))
    },
  },
  {
    title: '类型',
    key: 'type',
    width: 78,
    align: 'center',
    render(row) {
      const cfg = menuTypeConfig[row.type]

      return h(NTag, { bordered: false, type: cfg.type }, { default: () => cfg.label })
    },
  },
  {
    title: '路由',
    key: 'path',
    minWidth: 128,
    ellipsis: { tooltip: true },
    render(row) {
      return displayText(row.path)
    },
  },
  {
    title: '权限标识',
    key: 'code',
    minWidth: 150,
    ellipsis: { tooltip: true },
    render(row) {
      return displayText(row.code)
    },
  },
  {
    title: '排序',
    key: 'sort',
    width: 72,
    align: 'center',
  },
  {
    title: '状态',
    key: 'status',
    width: 84,
    align: 'center',
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
    title: '更新时间',
    key: 'updated_at',
    width: 150,
    render(row) {
      return formatTime(row.updated_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 164,
    fixed: 'right',
    render(row) {
      const canCreateChild = row.type !== MenuType.Button && props.canUse('system:menu:create')
      const nextStatus =
        row.status === MenuStatus.Enabled ? MenuStatus.Disabled : MenuStatus.Enabled

      return h(
        NSpace,
        { class: 'ez-row-actions', size: 6, align: 'center' },
        {
          default: () =>
            [
              canCreateChild
                ? h(EzActionButton, {
                    iconOnly: true,
                    kind: 'add-child',
                    label: row.type === MenuType.Menu ? '新增按钮' : '新增子级',
                    size: 'tiny',
                    type: 'primary',
                    secondary: true,
                    onClick: () => emit('createChild', row),
                  })
                : null,
              props.canUse('system:menu:update')
                ? h(EzActionButton, {
                    iconOnly: true,
                    kind: 'edit',
                    label: '编辑',
                    size: 'tiny',
                    secondary: true,
                    onClick: () => emit('edit', row),
                  })
                : null,
              props.canUse('system:menu:status')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => emit('toggleStatus', row, nextStatus) },
                    {
                      trigger: () =>
                        h(EzActionButton, {
                          iconOnly: true,
                          kind: nextStatus === MenuStatus.Disabled ? 'disable' : 'enable',
                          label: nextStatus === MenuStatus.Disabled ? '禁用' : '启用',
                          size: 'tiny',
                          secondary: true,
                          tooltip: false,
                          type: nextStatus === MenuStatus.Disabled ? 'error' : 'success',
                        }),
                      default: () =>
                        `确认${nextStatus === MenuStatus.Disabled ? '禁用' : '启用'}该菜单？`,
                    },
                  )
                : null,
              props.canUse('system:menu:delete')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => emit('delete', row) },
                    {
                      trigger: () =>
                        h(EzActionButton, {
                          iconOnly: true,
                          kind: 'delete',
                          label: '删除',
                          size: 'tiny',
                          secondary: true,
                          tooltip: false,
                          type: 'error',
                        }),
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

function rowKey(row: AdminMenu) {
  return row.id
}
</script>

<template>
  <EzDataTable
    :columns="columns"
    :data="displayMenus"
    :loading="loading"
    @refresh="emit('refresh')"
  >
    <template #toolbarSummary>
      <span>
        共 {{ flatMenuCount }} 个节点 · 目录 {{ stats.directoryCount }} · 菜单
        {{ stats.menuCount }} · 按钮 {{ stats.buttonCount }} · 已选 {{ selectedCount }} 项
      </span>
    </template>

    <template #toolbarActions>
      <NSpace :size="12">
        <EzActionButton
          kind="expand"
          label="展开全部"
          quaternary
          size="small"
          @click="emit('expandAll')"
        />
        <EzActionButton
          kind="collapse"
          label="收起全部"
          quaternary
          size="small"
          @click="emit('collapseAll')"
        />
        <NPopconfirm
          v-if="canUse('system:menu:delete')"
          :disabled="selectedCount === 0"
          @positive-click="emit('deleteSelected')"
        >
          <template #trigger>
            <EzActionButton
              kind="delete"
              label="删除选中"
              quaternary
              size="small"
              type="error"
              :disabled="selectedCount === 0"
            />
          </template>
          删除前请确认选中的菜单没有子菜单，也没有分配给任何角色。
        </NPopconfirm>
      </NSpace>
    </template>

    <template #body="{ tableColumns, tableScrollX, tableSize }">
      <NDataTable
        class="ez-table-fill-table menu-table"
        :columns="tableColumns"
        :data="displayMenus"
        :loading="loading"
        :row-key="rowKey"
        :checked-row-keys="checkedRowKeys"
        :expanded-row-keys="expandedRowKeys"
        :pagination="false"
        :scroll-x="tableScrollX"
        :size="tableSize"
        :bordered="false"
        children-key="children"
        flex-height
        @update:checked-row-keys="(keys) => emit('checkedRowKeysChange', keys)"
        @update:expanded-row-keys="expandedRowKeys = $event"
      />
    </template>
  </EzDataTable>
</template>
