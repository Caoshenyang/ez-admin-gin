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
  <NLayoutHeader bordered class="flex h-16 items-center justify-between bg-white px-6">
    <p class="text-[13px] text-[var(--ez-text-sub)]">{{ breadcrumbText }}</p>

    <div class="flex items-center gap-1.5">
      <NInput
        placeholder="搜索菜单 / 页面"
        clearable
        size="small"
        class="w-[300px]"
      >
        <template #prefix>
          <NIcon :component="SearchOutline" :size="16" class="text-[#94A3B8]" />
        </template>
      </NInput>

      <NButton quaternary circle class="h-9 w-9 rounded-[10px] !text-[var(--ez-text-sub)] hover:!text-[var(--ez-text-main)]">
        <NIcon :component="NotificationsOutline" :size="18" />
      </NButton>

      <NButton quaternary circle class="h-9 w-9 rounded-[10px] !text-[var(--ez-text-sub)] hover:!text-[var(--ez-text-main)]">
        <NIcon :component="ExpandOutline" :size="18" />
      </NButton>

      <NButton quaternary circle class="h-9 w-9 rounded-[10px] !text-[var(--ez-text-sub)] hover:!text-[var(--ez-text-main)]">
        <NIcon :component="MoonOutline" :size="18" />
      </NButton>

      <NDropdown trigger="click" :options="dropdownOptions" @select="(key) => emit('userAction', key)">
        <NButton
          quaternary
          class="gap-1.5 rounded-xl border border-[var(--ez-border)] px-3.5 py-2 transition-colors hover:border-slate-300 hover:bg-slate-50"
        >
          <span class="text-[13px] font-medium text-[#0F172A]">{{ displayName }}</span>
          <NIcon :component="ChevronDownOutline" :size="14" class="text-[#94A3B8]" />
        </NButton>
      </NDropdown>
    </div>
  </NLayoutHeader>
</template>
