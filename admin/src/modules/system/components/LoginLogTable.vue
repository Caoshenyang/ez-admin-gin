<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NPagination,
  NTag,
} from 'naive-ui'
import { h } from 'vue'

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
      return h('span', { class: 'font-semibold text-[#111827]' }, displayText(row.username))
    },
  },
  {
    title: '登录时间',
    key: 'created_at',
    width: 180,
    render(row) {
      return h('span', { class: 'text-[#374151]' }, formatTime(row.created_at))
    },
  },
  {
    title: '消息',
    key: 'message',
    minWidth: 200,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[#374151]' }, displayText(row.message))
    },
  },
  {
    title: 'IP 地址',
    key: 'ip',
    width: 150,
    render(row) {
      return h('span', { class: 'font-mono text-[13px] text-[#6B7280]' }, displayText(row.ip))
    },
  },
  {
    title: 'User-Agent',
    key: 'user_agent',
    minWidth: 220,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[#9CA3AF]' }, displayText(row.user_agent))
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
  <NCard class="min-h-0 flex-1 rounded-lg" :bordered="false" content-style="padding: 0;">
    <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-3">
      <span class="text-sm text-[#6B7280]">共 {{ total }} 条</span>
      <NButton text type="primary" @click="emit('refresh')">刷新</NButton>
    </div>

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
    />

    <div class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]">
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
  </NCard>
</template>

<style scoped>
.log-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #4b5563;
  background: #f9fafb;
  font-size: 13px;
}

.log-table :deep(.n-data-table-td) {
  color: #374151;
  font-size: 14px;
  padding: 10px 16px;
}

.log-table :deep(.log-table-row--failed .n-data-table-td) {
  background: #fef2f2;
}

.log-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}
</style>
