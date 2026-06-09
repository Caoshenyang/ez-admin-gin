<script setup lang="ts">
import type { Component } from 'vue'
import { computed } from 'vue'
import EzIconBox from './EzIconBox.vue'

const props = defineProps<{
  icon?: Component
  iconColor?: string
  title: string
  trend?: string
  trendType?: 'up' | 'down' | 'neutral'
  value: string | number
}>()

const trendClass = computed(() => {
  if (props.trendType === 'down') return 'ez-metric-card__trend--down'
  if (props.trendType === 'up') return 'ez-metric-card__trend--up'
  return 'ez-metric-card__trend--neutral'
})
</script>

<template>
  <article class="ez-metric-card">
    <div class="ez-metric-card__top">
      <div>
        <p class="ez-metric-card__title">{{ title }}</p>
        <strong class="ez-metric-card__value">{{ value }}</strong>
      </div>
      <EzIconBox v-if="icon" :icon="icon" :color="iconColor ?? 'var(--ez-primary)'" />
    </div>
    <p v-if="trend" class="ez-metric-card__trend" :class="trendClass">{{ trend }}</p>
  </article>
</template>

<style scoped>
.ez-metric-card {
  overflow: hidden;
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-card-bg);
  padding: 15px;
  box-shadow: var(--ez-component-shadow);
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.ez-metric-card:hover {
  border-color: var(--ez-brand-border);
  box-shadow: var(--ez-shadow-card);
  transform: translateY(-1px);
}

.ez-metric-card__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
}

.ez-metric-card__title {
  margin: 0;
  color: var(--ez-text-secondary);
  font-size: var(--ez-text-xs);
  font-weight: 700;
}

.ez-metric-card__value {
  display: block;
  margin-top: 6px;
  color: var(--ez-text-main);
  font-size: 24px;
  font-weight: 600;
  line-height: 1.1;
}

.ez-metric-card__trend {
  margin: 10px 0 0;
  font-size: 12px;
  font-weight: 600;
}

.ez-metric-card__trend--up {
  color: var(--ez-success);
}

.ez-metric-card__trend--down {
  color: var(--ez-danger);
}

.ez-metric-card__trend--neutral {
  color: var(--ez-text-secondary);
}
</style>
