<script setup lang="ts">
import { NIcon, NTag } from 'naive-ui'
import type { Component } from 'vue'

import type { FocusFact, HealthSummaryItem } from '../types/dashboard-page'

defineProps<{
  heroIcon: Component
  healthItems: HealthSummaryItem[]
  focusFacts: FocusFact[]
  statusText: string
  getHealthTagType: (status?: string) => 'success' | 'default' | 'error'
}>()
</script>

<template>
  <div class="overview-panel">
    <!-- header -->
    <div class="flex items-center gap-3">
      <div class="overview-panel__icon">
        <NIcon :component="heroIcon" :size="18" />
      </div>
      <div>
        <p class="text-[16px] font-semibold text-[#0F172A]">系统概览</p>
        <p class="mt-0.5 text-[12px] text-[#94A3B8]">{{ statusText }}</p>
      </div>
    </div>

    <!-- info grid -->
    <div class="overview-panel__grid">
      <div
        v-for="item in healthItems"
        :key="item.label"
        class="overview-cell"
      >
        <span class="overview-cell__label">{{ item.label }}</span>
        <NTag
          :type="getHealthTagType(item.status)"
          size="small"
          round
          :bordered="false"
        >
          {{ item.value }}
        </NTag>
      </div>
      <div
        v-for="fact in focusFacts"
        :key="fact.label"
        class="overview-cell"
      >
        <span class="overview-cell__label">{{ fact.label }}</span>
        <span class="overview-cell__value">{{ fact.value }}</span>
        <span class="overview-cell__hint">{{ fact.hint }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overview-panel {
  padding: 28px;
  background: linear-gradient(135deg, #F4FDFA 0%, #F8FBFF 100%);
  border-radius: 18px;
}

.overview-panel__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border-radius: 12px;
  background: linear-gradient(135deg, #14B8A6 0%, #0D8B82 100%);
  color: #ffffff;
  flex-shrink: 0;
}

.overview-panel__grid {
  display: grid;
  margin-top: 22px;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;
}

.overview-cell {
  background: #ffffff;
  border-radius: 14px;
  padding: 16px 18px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
  display: flex;
  flex-direction: column;
  gap: 6px;
  transition: box-shadow 0.2s ease;
}

.overview-cell:hover {
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
}

.overview-cell__label {
  font-size: 12px;
  color: #64748B;
}

.overview-cell__value {
  font-size: 15px;
  font-weight: 600;
  color: #0F172A;
}

.overview-cell__hint {
  font-size: 11px;
  color: #94A3B8;
}

@media (max-width: 768px) {
  .overview-panel__grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
