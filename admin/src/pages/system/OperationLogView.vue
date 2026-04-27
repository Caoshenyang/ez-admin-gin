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
  NSpace,
} from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

import { getOperationLogs } from '../../api/operation-log'
import type { OperationLogItem, OperationLogListQuery } from '../../types/operation-log'

type RiskLevel = 'high' | 'medium' | 'low'

const loading = ref(false)
const logs = ref<OperationLogItem[]>([])
const total = ref(0)
const detailVisible = ref(false)
const detailRow = ref<OperationLogItem | null>(null)

const query = reactive<OperationLogListQuery>({
  page: 1,
  page_size: 10,
  username: '',
  method: '',
  path: '',
})

const methodOptions = [
  { label: '方法：全部', value: '' },
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
]

function getRiskLevel(row: OperationLogItem): RiskLevel {
  if (!row.success) return 'high'
  if (row.method === 'POST' || row.method === 'PUT' || row.method === 'DELETE') return 'medium'
  return 'low'
}

const riskMeta: Record<RiskLevel, { label: string; color: string; bg: string }> = {
  high: { label: '高', color: '#D03050', bg: '#FDECEF' },
  medium: { label: '中', color: '#F0A020', bg: '#FFF8E8' },
  low: { label: '低', color: '#18A058', bg: '#EEFBF3' },
}

const methodMeta: Record<string, { color: string; bg: string }> = {
  GET: { color: '#18A058', bg: '#EEFBF3' },
  POST: { color: '#2080F0', bg: '#EAF3FF' },
  PUT: { color: '#F0A020', bg: '#FFF8E8' },
  DELETE: { color: '#D03050', bg: '#FDECEF' },
}

function getModule(path: string): string {
  const segments = path.replace(/^\/api\/v\d+\//, '').split('/')
  const moduleKey = segments[1]
  const moduleMap: Record<string, string> = {
    users: '用户管理',
    roles: '角色权限',
    menus: '菜单管理',
    configs: '配置管理',
    files: '文件管理',
    'operation-logs': '操作日志',
    'login-logs': '登录日志',
    auth: '认证',
  }
  if (!moduleKey) {
    return path
  }
  return moduleMap[moduleKey] ?? moduleKey
}

function getAction(row: OperationLogItem): string {
  if (!row.success && row.error_message) return row.error_message
  const actionMap: Record<string, Record<string, string>> = {
    GET: { default: '查询' },
    POST: { default: '新增/提交' },
    PUT: { default: '更新' },
    DELETE: { default: '删除' },
  }
  return actionMap[row.method]?.default ?? row.method
}

const columns: DataTableColumns<OperationLogItem> = [
  {
    title: '时间',
    key: 'created_at',
    width: 140,
    render(row) {
      return h('span', { class: 'tabular-nums text-[#374151]' }, formatTime(row.created_at))
    },
  },
  {
    title: '操作人',
    key: 'username',
    width: 90,
    render(row) {
      return h('span', { class: 'font-semibold text-[#111827]' }, row.username || '-')
    },
  },
  {
    title: '模块',
    key: 'module',
    width: 90,
    render(row) {
      return h('span', { class: 'text-[#374151]' }, getModule(row.path))
    },
  },
  {
    title: '行为',
    key: 'action',
    width: 130,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[#374151]' }, getAction(row))
    },
  },
  {
    title: 'IP 地址',
    key: 'ip',
    width: 120,
    render(row) {
      return h('span', { class: 'text-[#6B7280]' }, row.ip || '-')
    },
  },
  {
    title: '风险',
    key: 'risk',
    width: 56,
    align: 'center',
    render(row) {
      const risk = getRiskLevel(row)
      const meta = riskMeta[risk]
      return h(
        'span',
        { style: `color:${meta.color};font-weight:600;font-size:13px` },
        meta.label,
      )
    },
  },
  {
    title: '结果',
    key: 'success',
    width: 56,
    align: 'center',
    render(row) {
      return h(
        'span',
        { style: `color:${row.success ? '#18A058' : '#D03050'};font-weight:600;font-size:13px` },
        row.success ? '成功' : '失败',
      )
    },
  },
  {
    title: '',
    key: 'detail',
    width: 48,
    fixed: 'right',
    render(row) {
      return h(
        'span',
        {
          class: 'cursor-pointer text-[#2080F0] hover:underline',
          onClick: () => openDetail(row),
        },
        '详情',
      )
    },
  },
]

function rowProps(row: OperationLogItem) {
  const risk = getRiskLevel(row)
  return { style: `background: ${riskMeta[risk].bg}` }
}

function formatTime(value: string) {
  if (!value) return '-'
  const d = new Date(value)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function formatTimeFull(value: string) {
  if (!value) return '-'
  const d = new Date(value)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function openDetail(row: OperationLogItem) {
  detailRow.value = row
  detailVisible.value = true
}

function handleSearch() {
  query.page = 1
  void loadLogs()
}

function handleReset() {
  query.page = 1
  query.page_size = 10
  query.username = ''
  query.method = ''
  query.path = ''
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
    const data = await getOperationLogs({
      ...query,
      username: query.username?.trim() || undefined,
      method: query.method || undefined,
      path: query.path?.trim() || undefined,
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
        <h1 class="text-[26px] font-bold text-[#111827]">操作日志</h1>
        <p class="mt-1 text-sm text-[#6B7280]">查看系统操作记录，按用户、方法和路径筛选。</p>
      </div>

      <NCard :bordered="false" class="rounded-lg">
        <NSpace align="center" :wrap="true" :size="12">
          <NInput
            v-model:value="query.username"
            clearable
            placeholder="操作人"
            class="w-40"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.method" :options="methodOptions" class="w-36" />
          <NInput
            v-model:value="query.path"
            clearable
            placeholder="路径"
            class="w-52"
            @keyup.enter="handleSearch"
          />
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
          :row-key="(row: OperationLogItem) => row.id"
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

    <NDrawer
      v-model:show="detailVisible"
      :width="400"
      placement="right"
      class="log-drawer"
    >
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
              <span
                v-if="detailRow"
                class="inline-flex h-5 items-center rounded px-2 text-[11px] font-bold"
                :style="{
                  background: riskMeta[getRiskLevel(detailRow)].bg,
                  color: riskMeta[getRiskLevel(detailRow)].color,
                }"
              >
                {{ riskMeta[getRiskLevel(detailRow!)].label }}风险
              </span>
            </div>
            <p v-if="detailRow" class="mt-1 text-xs text-[#64748b]">
              {{ formatTimeFull(detailRow.created_at) }} · {{ detailRow.username }}
            </p>
          </div>
        </template>

        <div v-if="detailRow" class="flex flex-col gap-4">
          <!-- 请求概览 -->
          <div class="detail-section">
            <div class="detail-section__head">请求概览</div>
            <div class="detail-kv">
              <div class="detail-kv__label">请求地址</div>
              <div class="detail-kv__value font-mono text-[13px]">{{ detailRow.method }} {{ detailRow.path }}</div>
            </div>
            <div class="detail-kv">
              <div class="detail-kv__label">路由</div>
              <div class="detail-kv__value">{{ detailRow.route_path || '-' }}</div>
            </div>
            <div class="detail-tags">
              <span
                class="detail-tag"
                :style="{ background: methodMeta[detailRow.method]?.bg, color: methodMeta[detailRow.method]?.color }"
              >{{ detailRow.method }}</span>
              <span class="detail-tag detail-tag--muted">{{ detailRow.latency_ms }} ms</span>
              <span class="detail-tag detail-tag--muted">{{ detailRow.status_code }}</span>
              <span class="detail-tag detail-tag--muted">{{ detailRow.ip || '-' }}</span>
            </div>
          </div>

          <!-- 操作信息 -->
          <div class="detail-section">
            <div class="detail-section__head">操作信息</div>
            <div class="detail-grid">
              <div class="detail-kv">
                <div class="detail-kv__label">操作人</div>
                <div class="detail-kv__value">{{ detailRow.username || '-' }}</div>
              </div>
              <div class="detail-kv">
                <div class="detail-kv__label">所属模块</div>
                <div class="detail-kv__value">{{ getModule(detailRow.path) }}</div>
              </div>
              <div class="detail-kv">
                <div class="detail-kv__label">操作行为</div>
                <div class="detail-kv__value">{{ getAction(detailRow) }}</div>
              </div>
              <div class="detail-kv">
                <div class="detail-kv__label">执行结果</div>
                <div class="detail-kv__value">
                  <span
                    class="inline-flex items-center gap-1"
                    :style="{ color: detailRow.success ? '#18A058' : '#D03050', fontWeight: 600 }"
                  >
                    {{ detailRow.success ? '成功' : '失败' }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 请求参数（终端风格） -->
          <div class="detail-terminal">
            <div class="detail-terminal__head">请求参数</div>
            <div v-if="detailRow.query" class="detail-terminal__line">{{ detailRow.query }}</div>
            <div v-if="!detailRow.query" class="detail-terminal__line detail-terminal__line--muted">无查询参数</div>
            <div v-if="detailRow.user_agent" class="detail-terminal__line detail-terminal__line--dim">
              UA: {{ detailRow.user_agent }}
            </div>
          </div>

          <!-- 失败原因 -->
          <div
            v-if="!detailRow.success"
            class="detail-error"
          >
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

/* 抽屉头部 */
.detail-header {
  padding: 20px 24px 16px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid #e9eff6;
}

/* 分节 */
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
  color: #6B7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

/* KV 对 */
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
  color: #9CA3AF;
}

.detail-kv__value {
  font-size: 13px;
  color: #111827;
  line-height: 1.5;
}

/* 标签行 */
.detail-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 2px;
}

.detail-tag {
  display: inline-flex;
  align-items: center;
  height: 26px;
  padding: 0 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 700;
}

.detail-tag--muted {
  background: #F3F4F6;
  color: #374151;
  font-weight: 600;
}

/* 终端参数块 */
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
  color: #D1D5DB;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-terminal__line {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  line-height: 1.6;
  color: #F9FAFB;
  word-break: break-all;
}

.detail-terminal__line--muted {
  color: #6B7280;
}

.detail-terminal__line--dim {
  color: #9CA3AF;
  font-size: 11px;
  margin-top: 4px;
}

/* 错误块 */
.detail-error {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px 16px;
  border-radius: 10px;
  background: #FDECEF;
}

.detail-error__head {
  font-size: 12px;
  font-weight: 700;
  color: #D03050;
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
  background: #D03050;
  font-size: 11px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.detail-error__msg {
  font-size: 13px;
  color: #111827;
}
</style>
