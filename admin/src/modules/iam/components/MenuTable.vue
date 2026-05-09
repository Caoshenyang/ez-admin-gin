<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NPopconfirm, NSpace, NTag, NTooltip } from 'naive-ui'
import { h } from 'vue'

import TableStatsBar from '@/components/TableStatsBar.vue'
import { displayText } from '@/utils/format'
import { MenuStatus, MenuType, type AdminMenu } from '@/modules/iam/types/menu'

const props = defineProps<{
  canUse: (code: string) => boolean
  displayMenus: AdminMenu[]
  flatMenuCount: number
  loading: boolean
  stats: { directoryCount: number; menuCount: number; buttonCount: number }
}>()

const expandedRowKeys = defineModel<Array<string | number>>('expandedRowKeys', { required: true })

const emit = defineEmits<{
  collapseAll: []
  createChild: [row: AdminMenu]
  delete: [row: AdminMenu]
  edit: [row: AdminMenu]
  expandAll: []
  refresh: []
  toggleStatus: [row: AdminMenu, status: MenuStatus]
}>()

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
        h('span', { class: 'font-medium text-[#0F172A]' }, displayText(row.title)),
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
      return displayText(row.path)
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
      const canCreateChild = row.type !== MenuType.Button && props.canUse('system:menu:create')
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
                      onClick: () => emit('createChild', row),
                    },
                    { default: () => (row.type === MenuType.Menu ? '+ 按钮' : '+ 子级') },
                  )
                : null,
              props.canUse('system:menu:update')
                ? h(
                    NButton,
                    {
                      size: 'tiny',
                      secondary: true,
                      onClick: () => emit('edit', row),
                    },
                    { default: () => '编辑' },
                  )
                : null,
              props.canUse('system:menu:status')
                ? h(
                    NTooltip,
                    {},
                    {
                      trigger: () =>
                        h(
                          NPopconfirm,
                          { onPositiveClick: () => emit('toggleStatus', row, nextStatus) },
                          {
                            trigger: () =>
                              h(
                                NButton,
                                {
                                  size: 'tiny',
                                  type: nextStatus === MenuStatus.Disabled ? 'error' : 'success',
                                  secondary: true,
                                },
                                { default: () => (nextStatus === MenuStatus.Disabled ? '禁用' : '启用') },
                              ),
                            default: () =>
                              `确认${nextStatus === MenuStatus.Disabled ? '禁用' : '启用'}该菜单？`,
                          },
                        ),
                      default: () => '切换菜单可见状态',
                    },
                  )
                : null,
              props.canUse('system:menu:delete')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => emit('delete', row) },
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

// rowKey 函数。
function rowKey(row: AdminMenu) {
  return row.id
}
</script>

<template>
  <NCard class="ez-table-card min-h-0 flex-1" :bordered="false" content-class="ez-card-content-reset">
    <TableStatsBar>
      <span class="text-xs text-[#64748B]">
        共 {{ flatMenuCount }} 个节点 · 目录 {{ stats.directoryCount }} · 菜单 {{ stats.menuCount }} · 按钮 {{ stats.buttonCount }}
      </span>
      <template #actions>
        <NSpace :size="12">
          <NButton text size="small" @click="emit('expandAll')">展开全部</NButton>
          <NButton text size="small" @click="emit('collapseAll')">收起全部</NButton>
          <NButton text size="small" type="primary" @click="emit('refresh')">刷新</NButton>
        </NSpace>
      </template>
    </TableStatsBar>

    <NDataTable
      class="menu-table"
      :columns="columns"
      :data="displayMenus"
      :loading="loading"
      :row-key="rowKey"
      :expanded-row-keys="expandedRowKeys"
      :pagination="false"
      :bordered="false"
      children-key="children"
      @update:expanded-row-keys="expandedRowKeys = $event"
    />
  </NCard>
</template>
