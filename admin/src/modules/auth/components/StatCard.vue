<script setup lang="ts">
import { NIcon } from 'naive-ui'
import type { Component } from 'vue'

defineProps<{
  label: string
  value: string
  hint: string
  iconClass: string
  icon: Component
}>()
</script>

<template>
  <div class="stat-card">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-[var(--ez-text-sm)] font-medium text-[var(--ez-text-light)]">{{ label }}</p>
        <p
          class="mt-2 text-[var(--ez-text-hero)] font-bold tracking-tight text-[var(--ez-text-main)]"
        >
          {{ value }}
        </p>
      </div>
      <div class="stat-card__icon" :class="iconClass">
        <NIcon :component="icon" :size="22" />
      </div>
    </div>
    <p class="mt-3 text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">{{ hint }}</p>
  </div>
</template>

<style scoped>
.stat-card {
  background: var(--ez-card-bg);
  border-radius: var(--ez-radius-xl);
  padding: 24px;
  border: 1px solid var(--ez-border);
  box-shadow: var(--ez-shadow-md);
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.stat-card:hover {
  border-color: var(--ez-brand-border);
  box-shadow: var(--ez-shadow-lg);
  transform: translateY(-2px);
}

.stat-card__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: var(--ez-radius-lg);
  flex-shrink: 0;
}
</style>
