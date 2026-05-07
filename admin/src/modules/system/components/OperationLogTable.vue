<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NPagination, NTag } from 'naive-ui'
import { h } from 'vue'

import TableStatsBar from '@/components/TableStatsBar.vue'
import { formatTime } from '@/utils/format'
import type { OperationLogItem } from '../types/operation-log'
import {
  getAction,
  getModule,
  getRiskLevel,
  riskMeta,
} from '../composables/useOperationLogPage'

defineProps<{
  loading: boolean
  logs: OperationLogItem[]
  page: number
  pageSize: number
  total: number
}>()

const emit = defineEmits<{
  detail: [row: OperationLogItem]
  pageChange: [page: number]
  pageSizeChange: [pageSize: number]
  refresh: []
}>()

function rowProps(row: OperationLogItem) {
  return {
    class: 'operation-table-row',
    style: `background:${riskMeta[getRiskLevel(row)].bg};`,
  }
}

const columns: DataTableColumns<OperationLogItem> = [
  {
    title: '操作时间',
    key: 'created_at',
    width: 180,
    render(row) {
      return h('span', { class: 'text-[#374151]' }, formatTime(row.created_at))
    },
  },
  {
    title: '操作人',
    key: 'username',
    width: 120,
    render(row) {
      return h('span', { class: 'font-semibold text-[#111827]' }, row.username || '-')
    },
  },
  {
    title: '模块',
    key: 'path',
    width: 120,
    render(row) {
      return h('span', { class: 'text-[#374151]' }, getModule(row.path))
    },
  },
  {
    title: '方法',
    key: 'method',
    width: 100,
    render(row) {
      return h(NTag, { bordered: false, type: row.method === 'GET' ? 'success' : 'info' }, { default: () => row.method })
    },
  },
  {
    title: '行为',
    key: 'action',
    minWidth: 180,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[#374151]' }, getAction(row))
    },
  },
  {
    title: '风险',
    key: 'risk',
    width: 110,
    render(row) {
      const risk = riskMeta[getRiskLevel(row)]
      return h(NTag, { bordered: false, type: risk.tagType }, { default: () => risk.label })
    },
  },
  {
    title: '结果',
    key: 'success',
    width: 90,
    render(row) {
      return h(
        NTag,
        { bordered: false, type: row.success ? 'success' : 'error' },
        { default: () => (row.success ? '成功' : '失败') },
      )
    },
  },
  {
    title: '详情',
    key: 'detail',
    width: 80,
    fixed: 'right',
    render(row) {
      return h(
        NButton,
        { text: true, type: 'primary', onClick: () => emit('detail', row) },
        { default: () => '查看' },
      )
    },
  },
]
</script>

<template>
  <NCard
    class="min-h-0 flex-1 rounded-lg"
    :bordered="false"
    content-style="height: 100%; padding: 0;"
  >
    <TableStatsBar>
      <span class="text-sm text-[#6B7280]">共 {{ total }} 条</span>
      <template #actions>
        <NButton text type="primary" @click="emit('refresh')">刷新</NButton>
      </template>
    </TableStatsBar>

    <NDataTable
      remote
      class="operation-table h-full"
      style="height: calc(100% - 105px)"
      :columns="columns"
      :data="logs"
      :loading="loading"
      :pagination="false"
      :row-key="(row: OperationLogItem) => row.id"
      :row-props="rowProps"
      :bordered="false"
      flex-height
    />

    <div class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]">
      <span>共 {{ total }} 条</span>
      <NPagination
        :page="page"
        :page-size="pageSize"
        :item-count="total"
        :page-sizes="[10, 20, 50]"
        show-size-picker
        @update:page="emit('pageChange', $event)"
        @update:page-size="emit('pageSizeChange', $event)"
      />
    </div>
  </NCard>
</template>

<style scoped>
.operation-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #4b5563;
  background: #f9fafb;
  font-size: 13px;
}

.operation-table :deep(.n-data-table-td) {
  color: #374151;
  font-size: 14px;
  padding: 10px 16px;
}

.operation-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  filter: brightness(0.985);
}
</style>
