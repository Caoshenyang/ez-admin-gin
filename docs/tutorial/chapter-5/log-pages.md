---
title: 日志页面
description: "实现操作日志和登录日志查询页面。"
---

# 日志页面

上一节已经把配置和文件管理接成了真实页面。现在补齐系统管理的最后一组功能：操作日志和登录日志。

完成这一节后，侧边栏里的"操作日志"和"登录日志"不再停留在占位页。操作日志页面展示每次 API 调用的方法、路径、状态码和耗时；登录日志页面展示每次登录尝试的用户、IP、状态和 User-Agent。

::: tip 🎯 本节目标
这一节会把 `system/OperationLogView` 和 `system/LoginLogView` 从占位页换成真实页面，并补齐日志相关的类型和 API 封装。两个页面都是只读查询，不需要新增、编辑或删除操作。
:::

## 先看接口边界

日志接口只有查询，不提供写入：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/operation-logs` | 操作日志分页列表 |
| `GET` | `/api/v1/system/login-logs` | 登录日志分页列表 |

操作日志支持按用户名、请求方法、路径和是否成功筛选。登录日志支持按用户名、IP 和登录状态筛选。

::: details 操作日志是怎么产生的
后端通过 `OperationLog` 中间件自动记录每次 API 调用。请求完成后，中间件会把用户、方法、路径、状态码、耗时和错误信息写入 `sys_operation_log` 表。登录日志则是在登录接口内部写入 `sys_login_log` 表，记录每次登录尝试的结果。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
└─ src/
   ├─ api/
   │  ├─ operation-log.ts
   │  └─ login-log.ts
   ├─ pages/
   │  └─ system/
   │     ├─ OperationLogView.vue
   │     └─ LoginLogView.vue
   ├─ router/
   │  └─ dynamic-menu.ts
   └─ types/
      ├─ operation-log.ts
      └─ login-log.ts
```

## 开始前先确认

开始之前，先确认下面几件事：

- 已完成上一节 [配置与文件页面](./config-file-pages)。
- 登录后侧边栏能看到"操作日志"和"登录日志"。
- 后端 `/api/v1/system/operation-logs` 和 `/api/v1/system/login-logs` 可以正常返回数据。
- 数据库中已经有若干操作日志和登录日志记录（登录和操作几次后自动产生）。

## 🛠️ 完整代码

下面直接引入本节对应的完整项目文件，默认折叠。需要复制或对照时点击展开即可。

::: details `admin/src/types/operation-log.ts` — 操作日志类型

```ts
export interface OperationLogItem {
  id: number
  user_id: number
  username: string
  method: string
  path: string
  route_path: string
  query: string
  ip: string
  user_agent: string
  status_code: number
  latency_ms: number
  success: boolean
  error_message: string
  created_at: string
}

export interface OperationLogListQuery {
  page: number
  page_size: number
  username?: string
  method?: string
  path?: string
  success?: string
}

export interface OperationLogListResponse {
  items: OperationLogItem[]
  total: number
  page: number
  page_size: number
}
```

:::

::: details `admin/src/types/login-log.ts` — 登录日志类型

```ts
export const LoginLogStatus = {
  Success: 1,
  Failed: 2,
} as const

export type LoginLogStatus = (typeof LoginLogStatus)[keyof typeof LoginLogStatus]

export interface LoginLogItem {
  id: number
  user_id: number
  username: string
  status: LoginLogStatus
  message: string
  ip: string
  user_agent: string
  created_at: string
}

export interface LoginLogListQuery {
  page: number
  page_size: number
  username?: string
  ip?: string
  status?: LoginLogStatus | 0
}

export interface LoginLogListResponse {
  items: LoginLogItem[]
  total: number
  page: number
  page_size: number
}
```

:::

::: details `admin/src/api/operation-log.ts` — 操作日志接口

```ts
import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  OperationLogListQuery,
  OperationLogListResponse,
} from '../types/operation-log'

export async function getOperationLogs(params: OperationLogListQuery) {
  const response = await http.get<ApiResponse<OperationLogListResponse>>('/system/operation-logs', { params })
  return response.data.data
}
```

:::

::: details `admin/src/api/login-log.ts` — 登录日志接口

```ts
import http from './http'

import type { ApiResponse } from '../types/http'
import type { LoginLogListQuery, LoginLogListResponse } from '../types/login-log'

export async function getLoginLogs(params: LoginLogListQuery) {
  const response = await http.get<ApiResponse<LoginLogListResponse>>('/system/login-logs', { params })
  return response.data.data
}
```

:::

::: details `admin/src/pages/system/OperationLogView.vue` — 操作日志页面

```vue
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
```

:::

::: details `admin/src/pages/system/LoginLogView.vue` — 登录日志页面

```vue
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
```

:::

::: details `admin/src/router/dynamic-menu.ts` — 动态路由映射

修改后，`system/OperationLogView` 和 `system/LoginLogView` 会从占位页切换为真实页面。

```ts
import type { MenuOption } from 'naive-ui'
import type { RouteRecordRaw } from 'vue-router'
import { computed, shallowRef } from 'vue'

import { MenuType, type AuthMenu } from '../types/menu'

type RouteComponent = NonNullable<RouteRecordRaw['component']>

const placeholderPage = () => import('../pages/system/PlaceholderPage.vue')

const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': () => import('../pages/system/HealthView.vue'),
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': () => import('../pages/system/RoleView.vue'),
  'system/MenuView': () => import('../pages/system/MenuView.vue'),
  'system/ConfigView': () => import('../pages/system/ConfigView.vue'),
  'system/FileView': () => import('../pages/system/FileView.vue'),
  'system/OperationLogView': () => import('../pages/system/OperationLogView.vue'),
  'system/LoginLogView': () => import('../pages/system/LoginLogView.vue'),
  'system/NoticeView': () => import('../pages/system/NoticeView.vue'),
}

const builtinMenuOptions: MenuOption[] = [
  {
    label: '工作台',
    key: '/dashboard',
  },
]

export const authMenus = shallowRef<AuthMenu[]>([])

export const sideMenuOptions = computed<MenuOption[]>(() => {
  return [...builtinMenuOptions, ...buildMenuOptions(authMenus.value)]
})

export const buttonPermissionCodes = computed(() => {
  return collectButtonCodes(authMenus.value)
})

export function setAuthMenus(menus: AuthMenu[]) {
  authMenus.value = menus
}

export function clearAuthMenus() {
  authMenus.value = []
}

export function buildDynamicRoutes(menus: AuthMenu[]) {
  return collectPageMenus(menus).map<RouteRecordRaw>((menu) => ({
    path: toChildRoutePath(menu.path),
    name: `menu-${menu.id}`,
    component: resolveRouteComponent(menu.component),
    props: {
      title: menu.title,
      description: `${menu.title} 页面后续会接入真实业务。`,
    },
    meta: {
      title: menu.title,
      menuCode: menu.code,
    },
  }))
}

export function findMenuTitleByPath(path: string) {
  return collectPageMenus(authMenus.value).find((menu) => menu.path === path)?.title
}

function buildMenuOptions(menus: AuthMenu[]) {
  return menus.map(toMenuOption).filter(isMenuOption)
}

function toMenuOption(menu: AuthMenu): MenuOption | null {
  if (menu.type === MenuType.Button) {
    return null
  }

  const children = buildMenuOptions(menu.children ?? [])
  const key = menu.path || menu.code

  return {
    label: menu.title,
    key,
    disabled: menu.type === MenuType.Directory && children.length === 0,
    children: children.length > 0 ? children : undefined,
  }
}

function isMenuOption(option: MenuOption | null): option is MenuOption {
  return option !== null
}

function collectPageMenus(menus: AuthMenu[]) {
  const result: AuthMenu[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Menu && menu.path) {
      result.push(menu)
    }

    result.push(...collectPageMenus(menu.children ?? []))
  }

  return result
}

function collectButtonCodes(menus: AuthMenu[]) {
  const result: string[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Button) {
      result.push(menu.code)
    }

    result.push(...collectButtonCodes(menu.children ?? []))
  }

  return result
}

function resolveRouteComponent(component: string) {
  return routeComponentMap[component] ?? placeholderPage
}

function toChildRoutePath(path: string) {
  return path.replace(/^\/+/, '')
}
```

:::

## ✅ 验证结果

先启动后端和前端：

::: code-group

```bash [后端]
cd server
go run .
```

```bash [前端]
cd admin
pnpm dev
```

:::

然后按下面顺序验证：

1. 使用 `admin / EzAdmin@123456` 登录。
2. 点击"系统管理 / 操作日志"，确认日志列表能正常加载。
3. 按用户名筛选，确认只显示匹配记录。
4. 按请求方法筛选（如 POST），确认过滤生效。
5. 按路径关键词筛选，确认模糊匹配正常。
6. 点击某条记录的"详情"，确认能看到错误信息、请求参数或 User-Agent。
7. 进入"系统管理 / 登录日志"，确认日志列表能正常加载。
8. 按用户名、IP 或状态筛选，确认过滤生效。

## 本节小结

这一节把系统管理的日志查询页面补齐了：

- 操作日志页面展示 API 调用记录，支持按用户、方法、路径筛选，点击详情可查看错误信息和请求参数。
- 登录日志页面展示登录尝试记录，支持按用户名、IP 和状态筛选。
- 两个页面都是只读查询，没有新增、编辑或删除操作。
- 日志数据由后端中间件和登录接口自动产生，不需要前端手动写入。

到这里，第 5 章前端管理台的所有页面都已完成。
