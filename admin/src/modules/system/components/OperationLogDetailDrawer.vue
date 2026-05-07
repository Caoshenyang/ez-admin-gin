<script setup lang="ts">
import { NButton, NDrawer, NDrawerContent, NTag } from 'naive-ui'

import type { OperationLogItem } from '../types/operation-log'
import {
  formatTimeFull,
  getAction,
  getModule,
  getRiskLevel,
  riskMeta,
} from '../composables/useOperationLogPage'

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
          <NButton @click="$emit('update:show', false)">关闭</NButton>
        </div>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped>
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
