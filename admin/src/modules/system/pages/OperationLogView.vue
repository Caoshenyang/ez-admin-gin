<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NDrawer,
  NDrawerContent,
  NInput,
  NPagination,
  NSelect,
  NTag,
} from 'naive-ui'
import { h, ref } from 'vue'

import { useRemotePagination } from '@/composables/useRemotePagination'
import { formatTime } from '@/utils/format'
import { getOperationLogs } from '../api/operation-log'
import type { OperationLogItem, OperationLogListQuery } from '../types/operation-log'

type RiskLevel = 'high' | 'medium' | 'low'

const detailVisible = ref(false)
const detailRow = ref<OperationLogItem | null>(null)

const methodOptions = [
  { label: '方法：全部', value: '' },
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
]

const successOptions = [
  { label: '结果：全部', value: '' },
  { label: '成功', value: 'true' },
  { label: '失败', value: 'false' },
]

const {
  items: logs,
  total,
  loading,
  query,
  load,
  handleSearch,
  handleReset: handleResetBase,
  handlePageChange,
  handlePageSizeChange,
} = useRemotePagination<OperationLogItem, OperationLogListQuery>(
  (params) =>
    getOperationLogs({
      ...params,
      username: params.username?.trim() || undefined,
      method: params.method || undefined,
      path: params.path?.trim() || undefined,
      success: params.success || undefined,
    }),
  {
    page: 1,
    page_size: 10,
    username: '',
    method: '',
    path: '',
    success: '',
  },
)

const riskMeta: Record<RiskLevel, { label: string; tagType: 'error' | 'warning' | 'success'; bg: string }> = {
  high: { label: '高风险', tagType: 'error', bg: '#fef2f2' },
  medium: { label: '中风险', tagType: 'warning', bg: '#fff7ed' },
  low: { label: '低风险', tagType: 'success', bg: '#f0fdf4' },
}

function handleReset() {
  handleResetBase()
  detailVisible.value = false
}

function getRiskLevel(row: OperationLogItem): RiskLevel {
  if (!row.success) return 'high'
  if (row.method === 'POST' || row.method === 'PUT' || row.method === 'DELETE') return 'medium'
  return 'low'
}

function getModule(path: string): string {
  const segments = path.replace(/^\/api\/v\d+\//, '').split('/')
  const moduleKey = segments[1]
  const moduleMap: Record<string, string> = {
    users: '用户管理',
    roles: '角色权限',
    menus: '菜单管理',
    configs: '配置管理',
    'dict-types': '字典管理',
    'dict-items': '字典管理',
    files: '文件管理',
    attachments: '附件管理',
    notices: '公告管理',
    'operation-logs': '操作日志',
    'login-logs': '登录日志',
    health: '系统状态',
    auth: '认证',
  }

  return moduleKey ? (moduleMap[moduleKey] ?? moduleKey) : path
}

function getAction(row: OperationLogItem): string {
  if (!row.success && row.error_message) return row.error_message

  const actionMap: Record<string, string> = {
    GET: '查询',
    POST: '提交',
    PUT: '更新',
    DELETE: '删除',
  }

  return actionMap[row.method] ?? row.method
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
        { text: true, type: 'primary', onClick: () => openDetail(row) },
        { default: () => '查看' },
      )
    },
  },
]

function rowProps(row: OperationLogItem) {
  return {
    class: 'operation-table-row',
    style: `background:${riskMeta[getRiskLevel(row)].bg};`,
  }
}

function openDetail(row: OperationLogItem) {
  detailRow.value = row
  detailVisible.value = true
}

function formatTimeFull(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <div>
        <h1 class="text-[26px] font-bold text-[#111827]">操作日志</h1>
        <p class="mt-1 text-sm text-[#6B7280]">追踪后台接口操作行为，快速查看风险等级、执行结果和具体请求细节。</p>
      </div>

      <NCard :bordered="false" class="rounded-lg">
        <div class="grid gap-3 xl:grid-cols-[180px_150px_minmax(0,1fr)_150px_auto]">
          <NInput
            v-model:value="query.username"
            clearable
            placeholder="操作人"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.method" :options="methodOptions" />
          <NInput
            v-model:value="query.path"
            clearable
            placeholder="请求路径"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.success" :options="successOptions" />
          <div class="flex gap-3 xl:justify-end">
            <NButton type="primary" @click="handleSearch">查询</NButton>
            <NButton @click="handleReset">重置</NButton>
          </div>
        </div>
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

    <NDrawer v-model:show="detailVisible" :width="420" placement="right" class="log-drawer">
      <NDrawerContent
        :native-scrollbar="false"
        :body-content-style="{ padding: '20px 24px 24px' }"
        :header-style="{ padding: 0 }"
        :footer-style="{ padding: '16px 24px', borderTop: '1px solid #edf2f7', background: 'rgba(248,250,252,0.85)' }"
      >
        <template #header>
          <div class="detail-header">
            <div class="flex items-center gap-3">
              <span class="text-lg font-bold text-[#111827]">日志详情</span>
              <NTag
                v-if="detailRow"
                :bordered="false"
                :type="riskMeta[getRiskLevel(detailRow)].tagType"
              >
                {{ riskMeta[getRiskLevel(detailRow)].label }}
              </NTag>
            </div>
            <p v-if="detailRow" class="mt-1 text-xs text-[#64748B]">
              {{ formatTimeFull(detailRow.created_at) }} · {{ detailRow.username || '-' }}
            </p>
          </div>
        </template>

        <div v-if="detailRow" class="flex flex-col gap-4">
          <div class="detail-section">
            <div class="detail-section__head">请求概览</div>
            <div class="detail-kv">
              <div class="detail-kv__label">请求地址</div>
              <div class="detail-kv__value font-mono text-[13px]">{{ detailRow.method }} {{ detailRow.path }}</div>
            </div>
            <div class="detail-kv">
              <div class="detail-kv__label">路由模板</div>
              <div class="detail-kv__value">{{ detailRow.route_path || '-' }}</div>
            </div>
            <div class="detail-kv">
              <div class="detail-kv__label">模块 / 行为</div>
              <div class="detail-kv__value">{{ getModule(detailRow.path) }} · {{ getAction(detailRow) }}</div>
            </div>
          </div>

          <div class="detail-section">
            <div class="detail-section__head">执行结果</div>
            <div class="detail-grid">
              <div class="detail-kv">
                <div class="detail-kv__label">状态码</div>
                <div class="detail-kv__value">{{ detailRow.status_code }}</div>
              </div>
              <div class="detail-kv">
                <div class="detail-kv__label">耗时</div>
                <div class="detail-kv__value">{{ detailRow.latency_ms }} ms</div>
              </div>
              <div class="detail-kv">
                <div class="detail-kv__label">IP 地址</div>
                <div class="detail-kv__value">{{ detailRow.ip || '-' }}</div>
              </div>
              <div class="detail-kv">
                <div class="detail-kv__label">执行结果</div>
                <div class="detail-kv__value">
                  <span :style="{ color: detailRow.success ? '#18A058' : '#D03050', fontWeight: 600 }">
                    {{ detailRow.success ? '成功' : '失败' }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div class="detail-terminal">
            <div class="detail-terminal__head">请求上下文</div>
            <div class="detail-terminal__line">{{ detailRow.query || '无查询参数' }}</div>
            <div class="detail-terminal__line detail-terminal__line--dim">UA: {{ detailRow.user_agent || '-' }}</div>
          </div>

          <div v-if="!detailRow.success" class="detail-error">
            <div class="detail-error__head">失败原因</div>
            <div class="detail-error__body">
              <span class="detail-error__code">HTTP {{ detailRow.status_code }}</span>
              <span class="detail-error__msg">{{ detailRow.error_message || '未知错误' }}</span>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end">
            <NButton @click="detailVisible = false">关闭</NButton>
          </div>
        </template>
      </NDrawerContent>
    </NDrawer>
  </main>
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

.detail-header {
  padding: 20px 24px 16px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid #e9eff6;
}

.detail-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 14px 16px;
  border: 1px solid #e9eff6;
  border-radius: 10px;
  background: #fff;
}

.detail-section__head {
  font-size: 12px;
  font-weight: 700;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px 16px;
}

.detail-kv {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.detail-kv__label {
  font-size: 11px;
  font-weight: 600;
  color: #9ca3af;
}

.detail-kv__value {
  font-size: 13px;
  color: #111827;
  line-height: 1.5;
}

.detail-terminal {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 12px 14px;
  border-radius: 8px;
  background: #111827;
}

.detail-terminal__head {
  font-size: 11px;
  font-weight: 700;
  color: #d1d5db;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-terminal__line {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  line-height: 1.6;
  color: #f9fafb;
  word-break: break-all;
}

.detail-terminal__line--dim {
  color: #9ca3af;
}

.detail-error {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px 16px;
  border-radius: 10px;
  background: #fef2f2;
}

.detail-error__head {
  font-size: 12px;
  font-weight: 700;
  color: #d03050;
}

.detail-error__body {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.detail-error__code {
  display: inline-flex;
  align-items: center;
  height: 22px;
  padding: 0 8px;
  border-radius: 4px;
  background: #d03050;
  font-size: 11px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.detail-error__msg {
  font-size: 13px;
  color: #111827;
}

@media (max-width: 1280px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
