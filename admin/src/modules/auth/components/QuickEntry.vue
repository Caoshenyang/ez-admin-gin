<script setup lang="ts">
import { ChevronForwardOutline } from '@vicons/ionicons5'
import { NEmpty, NIcon, NTag } from 'naive-ui'

import type { QuickLink } from '../types/dashboard-page'

defineProps<{
  links: QuickLink[]
}>()

const emit = defineEmits<{
  navigate: [path: string]
}>()
</script>

<template>
  <div>
    <div class="flex items-center justify-between">
      <p class="text-[15px] font-semibold text-[#0F172A]">快捷入口</p>
      <NTag round :bordered="false" type="info" size="small">{{ links.length }} 项</NTag>
    </div>

    <div v-if="links.length > 0" class="mt-4 grid gap-2.5">
      <button
        v-for="item in links"
        :key="item.path"
        type="button"
        class="quick-entry"
        @click="emit('navigate', item.path)"
      >
        <div class="min-w-0">
          <p class="text-[14px] font-medium text-[#0F172A]">{{ item.title }}</p>
          <p class="mt-0.5 truncate text-[13px] text-[#64748B]">{{ item.description }}</p>
        </div>
        <NIcon :component="ChevronForwardOutline" :size="16" class="quick-entry__arrow" />
      </button>
    </div>

    <NEmpty v-else class="mt-4" description="当前角色没有额外页面可跳转" />
  </div>
</template>

<style scoped>
.quick-entry {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  width: 100%;
  border: 1px solid #E5EAF3;
  border-radius: 14px;
  background: #ffffff;
  padding: 12px 14px;
  text-align: left;
  cursor: pointer;
  transition:
    border-color 0.14s ease,
    background-color 0.14s ease,
    box-shadow 0.14s ease;
}

.quick-entry:hover {
  background: #F8FAFC;
  border-color: rgba(37, 99, 235, 0.25);
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.06);
}

.quick-entry__arrow {
  flex-shrink: 0;
  color: #94A3B8;
  transition:
    color 0.14s ease,
    transform 0.14s ease;
}

.quick-entry:hover .quick-entry__arrow {
  color: #2563EB;
  transform: translateX(2px);
}
</style>
