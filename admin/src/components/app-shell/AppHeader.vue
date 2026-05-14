<script setup lang="ts">
import {
  ChevronDownOutline,
  MoonOutline,
  NotificationsOutline,
  SearchOutline,
  SunnyOutline,
} from '@vicons/ionicons5'
import type { DropdownOption } from 'naive-ui'
import { NButton, NDropdown, NIcon, NInput, NLayoutHeader, NTooltip } from 'naive-ui'
import { computed } from 'vue'
import { useThemeStore } from '../../stores/theme'

defineProps<{
  breadcrumbText: string
  displayName: string
  dropdownOptions: DropdownOption[]
}>()

const emit = defineEmits<{
  userAction: [key: string | number]
}>()

const themeStore = useThemeStore()

const themeIcon = computed(() => (themeStore.isDark ? SunnyOutline : MoonOutline))

const themeTooltip = computed(() => {
  const labels: Record<string, string> = { light: '浅色', dark: '深色', auto: '跟随系统' }
  const next = { light: 'dark', dark: 'auto', auto: 'light' } as const
  const nextLabel = labels[next[themeStore.mode]]
  return `当前: ${labels[themeStore.mode]}（点击切换到${nextLabel}）`
})
</script>

<template>
  <NLayoutHeader bordered class="flex h-16 items-center justify-between bg-[var(--ez-card-bg)] px-6">
    <p class="text-[13px] text-[var(--ez-text-sub)]">{{ breadcrumbText }}</p>

    <div class="flex items-center gap-1.5">
      <NInput
        placeholder="搜索菜单 / 页面"
        clearable
        size="small"
        class="w-[300px]"
      >
        <template #prefix>
          <NIcon :component="SearchOutline" :size="16" class="text-[var(--ez-text-light)]" />
        </template>
      </NInput>

      <NButton quaternary circle class="h-9 w-9 rounded-[10px] !text-[var(--ez-text-sub)] hover:!text-[var(--ez-text-main)]">
        <NIcon :component="NotificationsOutline" :size="18" />
      </NButton>

      <NTooltip trigger="hover">
        <template #trigger>
          <NButton quaternary circle class="h-9 w-9 rounded-[10px] !text-[var(--ez-text-sub)] hover:!text-[var(--ez-text-main)]" @click="themeStore.cycleMode()">
            <NIcon :component="themeIcon" :size="18" />
          </NButton>
        </template>
        {{ themeTooltip }}
      </NTooltip>

      <NDropdown trigger="click" :options="dropdownOptions" @select="(key) => emit('userAction', key)">
        <NButton
          quaternary
          class="gap-1.5 rounded-xl border border-[var(--ez-border)] px-3.5 py-2 text-[var(--ez-text-main)] transition-colors hover:border-[var(--ez-brand)] hover:bg-[var(--ez-page-bg)]"
        >
          <span class="text-[13px] font-medium">{{ displayName }}</span>
          <NIcon :component="ChevronDownOutline" :size="14" class="text-[var(--ez-text-light)]" />
        </NButton>
      </NDropdown>
    </div>
  </NLayoutHeader>
</template>
