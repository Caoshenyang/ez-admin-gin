<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import {
  NButton,
  NDataTable,
  NPagination,
  NTag,
} from 'naive-ui'
import { h } from 'vue'

import EzTableCard from '@/components/ez/EzTableCard.vue'
import TableStatsBar from '@/components/TableStatsBar.vue'
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
    width: 120,
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
    width: 140,
    render(row) {
      return h('span', { class: 'font-semibold text-[var(--ez-text-main)]' }, displayText(row.username))
    },
  },
  {
    title: '登录时间',
    key: 'created_at',
    width: 180,
    render(row) {
      return h('span', { class: 'text-[var(--ez-text-main)]' }, formatTime(row.created_at))
    },
  },
  {
    title: '消息',
    key: 'message',
    minWidth: 200,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[var(--ez-text-main)]' }, displayText(row.message))
    },
  },
  {
    title: 'IP 地址',
    key: 'ip',
    width: 150,
    render(row) {
      return h('span', { class: 'font-mono text-[var(--ez-text-sm)] text-[var(--ez-text-sub)]' }, displayText(row.ip))
    },
  },
  {
    title: 'User-Agent',
    key: 'user_agent',
    minWidth: 220,
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
  <EzTableCard>
    <TableStatsBar>
      <span>共 {{ total }} 条</span>
      <template #actions>
        <NButton text type="primary" @click="emit('refresh')">刷新</NButton>
      </template>
    </TableStatsBar>

    <NDataTable
      remote
      class="log-table"
      :columns="columns"
      :data="props.logs"
      :loading="props.loading"
      :pagination="false"
      :row-key="(row: LoginLogItem) => row.id"
      :row-props="rowProps"
      :bordered="false"
      :scroll-x="1010"
    />

    <div class="ez-table-footer">
      <span>共 {{ total }} 条</span>
      <NPagination
        :page="props.page"
        :page-size="props.pageSize"
        :item-count="props.total"
        :page-sizes="[10, 20, 50]"
        show-size-picker
        @update:page="emit('page-change', $event)"
        @update:page-size="emit('page-size-change', $event)"
      />
    </div>
  </EzTableCard>
</template>

<style scoped>
.log-table :deep(.log-table-row--failed .n-data-table-td) {
  background: var(--ez-danger-bg);
}
</style>
