<script setup lang="ts">
import { CloseOutline, EllipsisHorizontal } from '@vicons/ionicons5'
import { NButton, NIcon, NScrollbar } from 'naive-ui'

import type { WorkTab } from '@/stores/admin-shell'

defineProps<{
  activePath: string
  tabs: WorkTab[]
}>()

const emit = defineEmits<{
  navigate: [path: string]
  closeTab: [path: string]
  refresh: []
  closeCurrent: []
  closeOthers: []
  closeAll: []
}>()
</script>

<template>
  <div class="admin-tabs-bar">
    <NScrollbar x-scrollable trigger="none" class="min-w-0 flex-1">
      <div class="admin-tabs-track">
        <button
          v-for="tab in tabs"
          :key="tab.to"
          type="button"
          class="admin-tab-chip"
          :class="{ 'admin-tab-chip--active': activePath === tab.to }"
          @click="emit('navigate', tab.to)"
        >
          <span class="truncate">{{ tab.title }}</span>
          <span
            v-if="tab.closable"
            class="admin-tab-chip__close"
            @click.stop="emit('closeTab', tab.to)"
          >
            <NIcon :component="CloseOutline" :size="14" />
          </span>
        </button>
      </div>
    </NScrollbar>

    <div class="admin-tabs-actions">
      <NButton quaternary size="small" @click="emit('refresh')">刷新</NButton>
      <NButton quaternary size="small" @click="emit('closeCurrent')">关闭当前</NButton>
      <NButton quaternary size="small" @click="emit('closeOthers')">关闭其他</NButton>
      <NButton quaternary size="small" @click="emit('closeAll')">关闭全部</NButton>
      <NButton quaternary circle size="small">
        <template #icon>
          <NIcon :component="EllipsisHorizontal" />
        </template>
      </NButton>
    </div>
  </div>
</template>

<style scoped>
.admin-tabs-bar {
  display: flex;
  min-height: 42px;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #e5e7eb;
  background: #ffffff;
  padding: 0 16px;
}

.admin-tabs-track {
  display: inline-flex;
  min-width: 100%;
  align-items: center;
  gap: 8px;
  padding: 7px 0;
}

.admin-tab-chip {
  display: inline-flex;
  min-width: 0;
  max-width: 220px;
  align-items: center;
  gap: 8px;
  border: 1px solid #d9dee8;
  border-radius: 999px;
  background: #f9fafb;
  padding: 0 12px;
  height: 28px;
  color: #374151;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    color 0.2s ease;
}

.admin-tab-chip--active {
  border-color: #18a058;
  background: #18a058;
  color: #ffffff;
  font-weight: 600;
}

.admin-tab-chip__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 999px;
}

.admin-tab-chip__close:hover {
  background: rgba(255, 255, 255, 0.18);
}

.admin-tabs-actions {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 4px;
}
</style>
