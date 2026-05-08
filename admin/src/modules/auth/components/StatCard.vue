<script setup lang="ts">
import { NIcon } from 'naive-ui'
import type { Component } from 'vue'

defineProps<{
  label: string
  value: string
  hint: string
  accent: string
  iconBg: string
  icon: Component
}>()
</script>

<template>
  <div class="stat-card">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-[13px] font-medium text-[#94A3B8]">{{ label }}</p>
        <p class="mt-2 text-[32px] font-semibold tracking-tight text-[#0F172A]">
          {{ value }}
        </p>
      </div>
      <div
        class="stat-card__icon"
        :style="{ background: iconBg, color: accent }"
      >
        <NIcon :component="icon" :size="20" />
      </div>
    </div>
    <p class="mt-3 text-[12px] text-[#94A3B8]">{{ hint }}</p>
  </div>
</template>

<style scoped>
.stat-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 24px;
  border: 1px solid #E5EAF3;
  box-shadow: 0 1px 2px rgba(16, 24, 40, 0.04), 0 1px 3px rgba(16, 24, 40, 0.06);
  transition:
    border-color 0.14s ease,
    box-shadow 0.14s ease,
    transform 0.14s ease;
}

.stat-card:hover {
  border-color: rgba(37, 99, 235, 0.25);
  box-shadow: 0 6px 16px rgba(15, 23, 42, 0.08);
  transform: translateY(-2px);
}

.stat-card__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 14px;
  flex-shrink: 0;
}
</style>
