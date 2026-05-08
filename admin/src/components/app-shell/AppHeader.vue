<script setup lang="ts">
import {
  ChevronDownOutline,
  ExpandOutline,
  MoonOutline,
  NotificationsOutline,
  SearchOutline,
} from '@vicons/ionicons5'
import type { DropdownOption } from 'naive-ui'
import { NButton, NDropdown, NIcon, NInput, NLayoutHeader } from 'naive-ui'

defineProps<{
  breadcrumbText: string
  displayName: string
  dropdownOptions: DropdownOption[]
}>()

const emit = defineEmits<{
  userAction: [key: string | number]
}>()
</script>

<template>
  <NLayoutHeader bordered class="app-header">
    <p class="text-[13px] text-[#475569]">{{ breadcrumbText }}</p>

    <div class="flex items-center gap-1.5">
      <NInput
        placeholder="搜索菜单 / 页面"
        clearable
        size="small"
        class="w-44"
      >
        <template #prefix>
          <NIcon :component="SearchOutline" :size="16" />
        </template>
      </NInput>

      <button type="button" class="header-icon-btn">
        <NIcon :component="NotificationsOutline" :size="18" />
      </button>

      <button type="button" class="header-icon-btn">
        <NIcon :component="ExpandOutline" :size="18" />
      </button>

      <button type="button" class="header-icon-btn">
        <NIcon :component="MoonOutline" :size="18" />
      </button>

      <NDropdown trigger="click" :options="dropdownOptions" @select="(key) => emit('userAction', key)">
        <button type="button" class="header-user-btn">
          <span class="text-[13px] font-medium text-[#0F172A]">{{ displayName }}</span>
          <NIcon :component="ChevronDownOutline" :size="14" class="text-[#94A3B8]" />
        </button>
      </NDropdown>
    </div>
  </NLayoutHeader>
</template>

<style scoped>
.app-header {
  display: flex;
  height: 56px;
  align-items: center;
  justify-content: space-between;
  background: #ffffff;
  padding: 0 20px;
}

.header-icon-btn {
  display: inline-flex;
  height: 36px;
  width: 36px;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: #475569;
  cursor: pointer;
  transition:
    background-color 0.15s ease,
    color 0.15s ease;
}

.header-icon-btn:hover {
  background: #F6F8FB;
  color: #0F172A;
}

.header-user-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 1px solid #E5EAF3;
  border-radius: 10px;
  background: #ffffff;
  padding: 6px 12px;
  cursor: pointer;
  transition:
    border-color 0.15s ease,
    background-color 0.15s ease;
}

.header-user-btn:hover {
  border-color: #CBD5E1;
  background: #F6F8FB;
}
</style>
