<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NPopconfirm, NSpace, NTag } from 'naive-ui'
import { h } from 'vue'

import { displayText } from '@/utils/format'
import { DepartmentStatus, type DepartmentItem } from '@/modules/iam/types/department'

const props = defineProps<{
  canUse: (code: string) => boolean
  departments: DepartmentItem[]
  loading: boolean
}>()

const emit = defineEmits<{
  edit: [row: DepartmentItem]
  toggleStatus: [row: DepartmentItem, status: DepartmentStatus]
}>()

// formatTime 函数。
function formatTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

const columns: DataTableColumns<DepartmentItem> = [
  {
    title: '部门',
    key: 'name',
    minWidth: 260,
    render(row) {
      return h('div', { class: 'leading-6' }, [
        h('p', { class: 'font-semibold text-[var(--ez-text-main)]' }, displayText(row.name)),
        h('p', { class: 'text-xs text-[var(--ez-text-sub)]' }, displayText(row.code)),
      ])
    },
  },
  {
    title: '负责人',
    key: 'leader_user_id',
    width: 120,
    render(row) {
      return row.leader_user_id === 0 ? '未设置' : `用户 ${row.leader_user_id}`
    },
  },
  {
    title: '排序',
    key: 'sort',
    width: 90,
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row) {
      return h(
        NTag,
        {
          bordered: false,
          type: row.status === DepartmentStatus.Enabled ? 'success' : 'error',
        },
        { default: () => (row.status === DepartmentStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 180,
    render(row) {
      return formatTime(row.updated_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render(row) {
      const nextStatus =
        row.status === DepartmentStatus.Enabled
          ? DepartmentStatus.Disabled
          : DepartmentStatus.Enabled

      return h(
        NSpace,
        { size: 8 },
        {
          default: () => [
            props.canUse('system:department:update')
              ? h(
                  NButton,
                  {
                    size: 'small',
                    ghost: true,
                    type: 'info',
                    onClick: () => emit('edit', row),
                  },
                  { default: () => '编辑' },
                )
              : null,
            props.canUse('system:department:status')
              ? h(
                  NPopconfirm,
                  { onPositiveClick: () => emit('toggleStatus', row, nextStatus) },
                  {
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: 'small',
                          ghost: true,
                          type: nextStatus === DepartmentStatus.Disabled ? 'error' : 'success',
                        },
                        { default: () => (nextStatus === DepartmentStatus.Disabled ? '禁用' : '启用') },
                      ),
                    default: () =>
                      `确认${nextStatus === DepartmentStatus.Disabled ? '禁用' : '启用'}该部门？`,
                  },
                )
              : null,
          ],
        },
      )
    },
  },
]
</script>

<template>
  <NCard class="ez-table-card min-h-0 flex-1" :bordered="false" content-class="ez-card-content-reset">
    <NDataTable
      class="department-table"
      :columns="columns"
      :data="departments"
      :loading="loading"
      :pagination="false"
      :bordered="false"
      children-key="children"
    />
  </NCard>
</template>
