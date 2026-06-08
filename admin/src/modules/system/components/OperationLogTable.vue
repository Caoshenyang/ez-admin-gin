<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NTag } from 'naive-ui'
import { h } from 'vue'

import EzDataTable from '@/components/ez/EzDataTable.vue'
import { displayText, formatTime } from '@/utils/format'
import type { OperationLogItem } from '../types/operation-log'
import {
  getAction,
  getModule,
  getRiskLevel,
  riskMeta,
} from '../composables/operation-log-page.utils'

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
    width: 150,
    render(row) {
      return h('span', { class: 'text-[var(--ez-text-main)]' }, formatTime(row.created_at))
    },
  },
  {
    title: '操作人',
    key: 'username',
    width: 104,
    render(row) {
      return h('span', { class: 'font-semibold text-[var(--ez-text-main)]' }, displayText(row.username))
    },
  },
  {
    title: '模块',
    key: 'path',
    width: 104,
    render(row) {
      return h('span', { class: 'text-[var(--ez-text-main)]' }, getModule(row.path))
    },
  },
  {
    title: '方法',
    key: 'method',
    width: 84,
    render(row) {
      return h(NTag, { bordered: false, type: row.method === 'GET' ? 'success' : 'info' }, { default: () => displayText(row.method) })
    },
  },
  {
    title: '行为',
    key: 'action',
    minWidth: 160,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[var(--ez-text-main)]' }, getAction(row))
    },
  },
  {
    title: '风险',
    key: 'risk',
    width: 92,
    render(row) {
      const risk = riskMeta[getRiskLevel(row)]
      return h(NTag, { bordered: false, type: risk.tagType }, { default: () => risk.label })
    },
  },
  {
    title: '结果',
    key: 'success',
    width: 78,
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
    width: 70,
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
  <EzDataTable
    remote
    class="operation-table"
    :columns="columns"
    :data="logs"
    :loading="loading"
    :page="page"
    :page-size="pageSize"
    :row-key="(row: OperationLogItem) => row.id"
    :row-props="rowProps"
    :total="total"
    @page-change="emit('pageChange', $event)"
    @page-size-change="emit('pageSizeChange', $event)"
    @refresh="emit('refresh')"
  />
</template>
