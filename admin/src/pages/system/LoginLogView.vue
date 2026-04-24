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
} from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

import { getLoginLogs } from '../../api/login-log'
import { LoginLogStatus, type LoginLogItem, type LoginLogListQuery } from '../../types/login-log'

const loading = ref(false)
const logs = ref<LoginLogItem[]>([])
const total = ref(0)

const query = reactive<LoginLogListQuery>({
  page: 1,
  page_size: 10,
  username: '',
  ip: '',
  status: 0,
})

const statusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '成功', value: LoginLogStatus.Success },
  { label: '失败', value: LoginLogStatus.Failed },
]

const columns: DataTableColumns<LoginLogItem> = [
  {
    title: '时间',
    key: 'created_at',
    width: 140,
    render(row) {
      return h('span', { class: 'tabular-nums text-[#374151]' }, formatTime(row.created_at))
    },
  },
  {
    title: '用户',
    key: 'username',
    width: 100,
    render(row) {
      return h('span', { class: 'font-semibold text-[#111827]' }, row.username || '-')
    },
  },
  {
    title: '结果',
    key: 'status',
    width: 56,
    align: 'center',
    render(row) {
      const ok = row.status === LoginLogStatus.Success
      return h(
        'span',
        { style: `color:${ok ? '#18A058' : '#D03050'};font-weight:600;font-size:13px` },
        ok ? '成功' : '失败',
      )
    },
  },
  {
    title: '消息',
    key: 'message',
    width: 180,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[#374151]' }, row.message || '-')
    },
  },
  {
    title: 'IP',
    key: 'ip',
    width: 120,
    render(row) {
      return h('span', { class: 'text-[#6B7280]' }, row.ip || '-')
    },
  },
  {
    title: 'User-Agent',
    key: 'user_agent',
    minWidth: 180,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[#9CA3AF]' }, row.user_agent || '-')
    },
  },
]

function rowProps(row: LoginLogItem) {
  if (row.status === LoginLogStatus.Failed) {
    return { style: 'background: #FDECEF' }
  }
  return {}
}

function formatTime(value: string) {
  if (!value) return '-'
  const d = new Date(value)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function handleSearch() {
  query.page = 1
  void loadLogs()
}

function handleReset() {
  query.page = 1
  query.page_size = 10
  query.username = ''
  query.ip = ''
  query.status = 0
  void loadLogs()
}

function handlePageChange(page: number) {
  query.page = page
  void loadLogs()
}

function handlePageSizeChange(pageSize: number) {
  query.page = 1
  query.page_size = pageSize
  void loadLogs()
}

async function loadLogs() {
  loading.value = true
  try {
    const data = await getLoginLogs({
      ...query,
      username: query.username?.trim() || undefined,
      ip: query.ip?.trim() || undefined,
      status: query.status === 0 ? undefined : query.status,
    })
    logs.value = data.items
    total.value = data.total
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void loadLogs()
})
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-4 overflow-hidden">
      <div>
        <h1 class="text-[26px] font-bold text-[#111827]">登录日志</h1>
        <p class="mt-1 text-sm text-[#6B7280]">查看账号登录记录，按用户名、IP 和状态筛选。</p>
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
          <NButton text type="primary" @click="loadLogs">刷新</NButton>
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
          :bordered="false"
          :row-props="rowProps"
          flex-height
        />

        <div
          class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]"
        >
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
  color: #4B5563;
  background: #F9FAFB;
  font-size: 13px;
}

.log-table :deep(.n-data-table-td) {
  color: #374151;
  font-size: 14px;
  padding: 10px 16px;
}

.log-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: unset !important;
}

.log-table :deep(.n-data-table-tr) {
  transition: none;
}

.log-table :deep(.n-data-table-tr:hover) {
  filter: brightness(0.97);
}
</style>
