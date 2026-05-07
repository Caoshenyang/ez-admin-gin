<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NInput,
  NPagination,
  NSelect,
  NSpace,
  NTag,
} from 'naive-ui'
import { h } from 'vue'

import { useRemotePagination } from '@/composables/useRemotePagination'
import { formatTime } from '@/utils/format'
import { getLoginLogs } from '../api/login-log'
import { LoginLogStatus, type LoginLogItem, type LoginLogListQuery } from '../types/login-log'

const statusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '成功', value: LoginLogStatus.Success },
  { label: '失败', value: LoginLogStatus.Failed },
]

const {
  items: logs,
  total,
  loading,
  query,
  load,
  handleSearch,
  handleReset,
  handlePageChange,
  handlePageSizeChange,
} = useRemotePagination<LoginLogItem, LoginLogListQuery>(
  (params) =>
    getLoginLogs({
      ...params,
      username: params.username?.trim() || undefined,
      ip: params.ip?.trim() || undefined,
      status: params.status === 0 ? undefined : params.status,
    }),
  {
    page: 1,
    page_size: 10,
    username: '',
    ip: '',
    status: 0,
  },
)

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
      return h('span', { class: 'font-semibold text-[#111827]' }, row.username || '-')
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
      return h('span', { class: 'text-[#374151]' }, row.message || '-')
    },
  },
  {
    title: 'IP 地址',
    key: 'ip',
    width: 150,
    render(row) {
      return h('span', { class: 'font-mono text-[13px] text-[#6B7280]' }, row.ip || '-')
    },
  },
  {
    title: 'User-Agent',
    key: 'user_agent',
    minWidth: 220,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[#9CA3AF]' }, row.user_agent || '-')
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
  <main class="admin-page">
    <section class="admin-page-section">
      <div>
        <h1 class="text-[26px] font-bold text-[#111827]">登录日志</h1>
        <p class="mt-1 text-sm text-[#6B7280]">按用户名、IP 和登录状态回看账号登录轨迹，快速识别异常登录或失败重试。</p>
      </div>

      <NCard :bordered="false" class="rounded-lg">
        <NSpace align="center" :wrap="true" :size="12">
          <NInput
            v-model:value="query.username"
            clearable
            placeholder="用户名"
            class="w-40"
            @keyup.enter="handleSearch"
          />
          <NInput
            v-model:value="query.ip"
            clearable
            placeholder="IP 地址"
            class="w-44"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.status" :options="statusOptions" class="w-36" />
          <NButton type="primary" @click="handleSearch">查询</NButton>
          <NButton @click="handleReset">重置</NButton>
        </NSpace>
      </NCard>

      <NCard
        class="min-h-0 flex-1 rounded-lg"
        :bordered="false"
        content-style="height: 100%; padding: 0;"
      >
        <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-3">
          <span class="text-sm text-[#6B7280]">共 {{ total }} 条</span>
          <NButton text type="primary" @click="load">刷新</NButton>
        </div>

        <NDataTable
          remote
          class="log-table h-full"
          style="height: calc(100% - 105px)"
          :columns="columns"
          :data="logs"
          :loading="loading"
          :pagination="false"
          :row-key="(row: LoginLogItem) => row.id"
          :row-props="rowProps"
          :bordered="false"
          flex-height
        />

        <div class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]">
          <span>共 {{ total }} 条</span>
          <NPagination
            :page="query.page"
            :page-size="query.page_size"
            :item-count="total"
            :page-sizes="[10, 20, 50]"
            show-size-picker
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </NCard>
    </section>
  </main>
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
