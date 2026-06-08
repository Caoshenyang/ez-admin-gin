<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NTag } from 'naive-ui'
import { h } from 'vue'

import EzDataTable from '@/components/ez/EzDataTable.vue'
import { displayText, formatTime } from '@/utils/format'
import { LoginLogStatus, type LoginLogItem } from '../types/login-log'

const props = defineProps<{
  loading: boolean
  logs: LoginLogItem[]
  page: number
  pageSize: number
  total: number
}>()

const emit = defineEmits<{
  'page-change': [page: number]
  'page-size-change': [pageSize: number]
  refresh: []
}>()

const columns: DataTableColumns<LoginLogItem> = [
  {
    title: '登录结果',
    key: 'status',
    width: 94,
    render(row) {
      const ok = row.status === LoginLogStatus.Success
      return h(
        NTag,
        { bordered: false, type: ok ? 'success' : 'error' },
        { default: () => (ok ? '成功' : '失败') },
      )
    },
  },
  {
    title: '用户',
    key: 'username',
    width: 112,
    render(row) {
      return h(
        'span',
        { class: 'font-semibold text-[var(--ez-text-main)]' },
        displayText(row.username),
      )
    },
  },
  {
    title: '登录时间',
    key: 'created_at',
    width: 150,
    render(row) {
      return h('span', { class: 'text-[var(--ez-text-main)]' }, formatTime(row.created_at))
    },
  },
  {
    title: '消息',
    key: 'message',
    minWidth: 170,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[var(--ez-text-main)]' }, displayText(row.message))
    },
  },
  {
    title: 'IP 地址',
    key: 'ip',
    width: 128,
    render(row) {
      return h(
        'span',
        { class: 'font-mono text-[var(--ez-text-sm)] text-[var(--ez-text-sub)]' },
        displayText(row.ip),
      )
    },
  },
  {
    title: 'User-Agent',
    key: 'user_agent',
    minWidth: 190,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[var(--ez-text-light)]' }, displayText(row.user_agent))
    },
  },
]

function rowProps(row: LoginLogItem) {
  if (row.status === LoginLogStatus.Failed) {
    return { class: 'log-table-row log-table-row--failed' }
  }

  return { class: 'log-table-row' }
}
</script>

<template>
  <EzDataTable
    remote
    class="log-table"
    :columns="columns"
    :data="props.logs"
    :loading="props.loading"
    :page="props.page"
    :page-size="props.pageSize"
    :row-key="(row: LoginLogItem) => row.id"
    :row-props="rowProps"
    :total="props.total"
    @page-change="emit('page-change', $event)"
    @page-size-change="emit('page-size-change', $event)"
    @refresh="emit('refresh')"
  />
</template>

<style scoped>
.log-table :deep(.log-table-row--failed .n-data-table-td) {
  background: var(--ez-danger-bg);
}
</style>
