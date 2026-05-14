<script setup lang="ts">
import { NButton, NDrawer, NDrawerContent, NTag } from 'naive-ui'

import { displayText } from '@/utils/format'
import type { OperationLogItem } from '../types/operation-log'
import {
  formatTimeFull,
  getAction,
  getModule,
  getRiskLevel,
  riskMeta,
} from '../composables/operation-log-page.utils'

defineProps<{
  detailRow: OperationLogItem | null
  show: boolean
}>()

defineEmits<{
  'update:show': [value: boolean]
}>()
</script>

<template>
  <NDrawer :show="show" :width="420" placement="right" class="log-drawer" @update:show="(value) => $emit('update:show', value)">
    <NDrawerContent :native-scrollbar="false" body-content-class="p-5 pt-5" header-class="p-0" footer-class="border-t border-[var(--ez-border-light)] bg-slate-50/80 px-6 py-4">
      <template #header>
        <div class="border-b border-[var(--ez-border-light)] bg-[linear-gradient(135deg,#f8fafc_0%,#f1f5f9_100%)] px-6 py-5">
          <div class="flex items-center gap-3">
            <span class="text-lg font-bold text-[var(--ez-text-main)]">日志详情</span>
            <NTag
              v-if="detailRow"
              :bordered="false"
              :type="riskMeta[getRiskLevel(detailRow)].tagType"
            >
              {{ riskMeta[getRiskLevel(detailRow)].label }}
            </NTag>
          </div>
          <p v-if="detailRow" class="mt-1 text-xs text-[var(--ez-text-sub)]">
            {{ formatTimeFull(detailRow.created_at) }} · {{ detailRow.username || '-' }}
          </p>
        </div>
      </template>

      <div v-if="detailRow" class="flex flex-col gap-4">
        <div class="flex flex-col gap-2.5 rounded-[10px] border border-[var(--ez-border-light)] bg-[var(--ez-card-bg)] px-4 py-3.5">
          <div class="text-[12px] font-bold tracking-[0.05em] text-gray-500 uppercase">请求概览</div>
          <div class="flex flex-col gap-0.5">
            <div class="text-[11px] font-semibold text-slate-400">请求地址</div>
            <div class="font-mono text-[13px] leading-6 text-[var(--ez-text-main)]">{{ displayText(detailRow.method) }} {{ displayText(detailRow.path) }}</div>
          </div>
          <div class="flex flex-col gap-0.5">
            <div class="text-[11px] font-semibold text-slate-400">路由模板</div>
            <div class="text-[13px] leading-6 text-[var(--ez-text-main)]">{{ detailRow.route_path || '-' }}</div>
          </div>
          <div class="flex flex-col gap-0.5">
            <div class="text-[11px] font-semibold text-slate-400">模块 / 行为</div>
            <div class="text-[13px] leading-6 text-[var(--ez-text-main)]">{{ getModule(detailRow.path) }} · {{ getAction(detailRow) }}</div>
          </div>
        </div>

        <div class="flex flex-col gap-2.5 rounded-[10px] border border-[var(--ez-border-light)] bg-[var(--ez-card-bg)] px-4 py-3.5">
          <div class="text-[12px] font-bold tracking-[0.05em] text-gray-500 uppercase">执行结果</div>
          <div class="grid gap-x-4 gap-y-2.5 md:grid-cols-2">
            <div class="flex flex-col gap-0.5">
              <div class="text-[11px] font-semibold text-slate-400">状态码</div>
              <div class="text-[13px] leading-6 text-[var(--ez-text-main)]">{{ detailRow.status_code }}</div>
            </div>
            <div class="flex flex-col gap-0.5">
              <div class="text-[11px] font-semibold text-slate-400">耗时</div>
              <div class="text-[13px] leading-6 text-[var(--ez-text-main)]">{{ detailRow.latency_ms }} ms</div>
            </div>
            <div class="flex flex-col gap-0.5">
              <div class="text-[11px] font-semibold text-slate-400">IP 地址</div>
              <div class="text-[13px] leading-6 text-[var(--ez-text-main)]">{{ detailRow.ip || '-' }}</div>
            </div>
            <div class="flex flex-col gap-0.5">
              <div class="text-[11px] font-semibold text-slate-400">执行结果</div>
              <div class="text-[13px] leading-6 text-[var(--ez-text-main)]">
                <span :class="detailRow.success ? 'ez-status-text--success' : 'ez-status-text--danger'">
                  {{ detailRow.success ? '成功' : '失败' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex flex-col gap-1.5 rounded-lg bg-slate-900 px-3.5 py-3">
          <div class="text-[11px] font-bold tracking-[0.05em] text-slate-300 uppercase">请求上下文</div>
          <div class="break-all font-mono text-[12px] leading-6 text-slate-50">{{ displayText(detailRow.query, '无查询参数') }}</div>
          <div class="break-all font-mono text-[12px] leading-6 text-slate-400">UA: {{ displayText(detailRow.user_agent) }}</div>
        </div>

        <div v-if="!detailRow.success" class="flex flex-col gap-2 rounded-[10px] bg-rose-50 px-4 py-3.5">
          <div class="text-[12px] font-bold text-rose-700">失败原因</div>
          <div class="flex items-baseline gap-2">
            <span class="inline-flex h-[22px] shrink-0 items-center rounded px-2 text-[11px] font-bold text-white bg-[#d03050]">
              HTTP {{ detailRow.status_code }}
            </span>
            <span class="text-[13px] text-[var(--ez-text-main)]">{{ displayText(detailRow.error_message, '未知错误') }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end">
          <NButton @click="$emit('update:show', false)">关闭</NButton>
        </div>
    </template>
  </NDrawerContent>
  </NDrawer>
</template>
