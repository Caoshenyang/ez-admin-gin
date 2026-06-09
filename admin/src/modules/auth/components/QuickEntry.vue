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
      <div>
        <p class="text-[var(--ez-text-lg)] font-bold text-[var(--ez-text-main)]">常用入口</p>
        <p class="mt-1 text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">
          按当前角色开放的高频处理页面。
        </p>
      </div>
      <NTag round :bordered="false" type="info" size="small">{{ links.length }} 项</NTag>
    </div>

    <div v-if="links.length > 0" class="mt-4 grid gap-2">
      <button
        v-for="item in links"
        :key="item.path"
        type="button"
        class="quick-entry"
        @click="emit('navigate', item.path)"
      >
        <div class="min-w-0">
          <p class="text-[var(--ez-text-sm)] font-bold text-[var(--ez-text-main)]">
            {{ item.title }}
          </p>
          <p class="mt-0.5 truncate text-[13px] text-[var(--ez-text-sub)]">
            {{ item.description }}
          </p>
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
  border: 1px solid var(--ez-border-light);
  border-radius: var(--ez-radius-sm);
  background: var(--ez-surface-subtle);
  padding: 11px 12px;
  text-align: left;
  cursor: pointer;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    box-shadow 0.2s ease;
}

.quick-entry:hover {
  background: var(--ez-brand-soft);
  border-color: var(--ez-brand-border);
}

.quick-entry__arrow {
  flex-shrink: 0;
  color: var(--ez-text-light);
  transition:
    color 0.2s ease,
    transform 0.2s ease;
}

.quick-entry:hover .quick-entry__arrow {
  color: var(--ez-brand);
  transform: translateX(4px);
}
</style>
