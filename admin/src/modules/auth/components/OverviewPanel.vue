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
        <p class="text-[var(--ez-text-lg)] font-semibold text-[var(--ez-text-main)]">系统概览</p>
        <p class="mt-0.5 text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">{{ statusText }}</p>
      </div>
    </div>

    <!-- info grid -->
    <div class="overview-panel__grid">
      <div v-for="item in healthItems" :key="item.label" class="overview-cell">
        <span class="overview-cell__label">{{ item.label }}</span>
        <NTag :type="getHealthTagType(item.status)" size="small" round :bordered="false">
          {{ item.value }}
        </NTag>
      </div>
      <div v-for="fact in focusFacts" :key="fact.label" class="overview-cell">
        <span class="overview-cell__label">{{ fact.label }}</span>
        <span class="overview-cell__value">{{ fact.value }}</span>
        <span class="overview-cell__hint">{{ fact.hint }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overview-panel {
  padding: 18px;
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-card-bg);
}

.overview-panel__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-surface-subtle);
  color: var(--ez-brand);
  flex-shrink: 0;
}

.overview-panel__grid {
  display: grid;
  margin-top: 16px;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.overview-cell {
  display: flex;
  flex-direction: column;
  gap: 5px;
  border: 1px solid var(--ez-border-light);
  border-radius: var(--ez-radius-sm);
  background: var(--ez-surface-subtle);
  padding: 11px 12px;
}

.overview-cell__label {
  font-size: var(--ez-text-xs);
  color: var(--ez-text-sub);
}

.overview-cell__value {
  font-size: var(--ez-text-md);
  font-weight: 600;
  color: var(--ez-text-main);
}

.overview-cell__hint {
  font-size: var(--ez-text-xs);
  color: var(--ez-text-light);
}

@media (max-width: 768px) {
  .overview-panel__grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
