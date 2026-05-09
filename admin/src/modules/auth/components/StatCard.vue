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
        <p class="text-[13px] font-medium text-[#94A3B8]">{{ label }}</p>
        <p class="mt-2 text-[32px] font-bold tracking-tight text-[#0F172A]">
          {{ value }}
        </p>
      </div>
      <div class="stat-card__icon" :class="iconClass">
        <NIcon :component="icon" :size="22" />
      </div>
    </div>
    <p class="mt-3 text-[12px] text-[#94A3B8]">{{ hint }}</p>
  </div>
</template>

<style scoped>
.stat-card {
  background: #ffffff;
  border-radius: 18px;
  padding: 24px;
  border: 1px solid #E6ECF3;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.06);
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.stat-card:hover {
  border-color: rgba(20, 184, 166, 0.22);
  box-shadow: 0 12px 32px rgba(15, 23, 42, 0.10);
  transform: translateY(-2px);
}

.stat-card__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 14px;
  flex-shrink: 0;
}
</style>
