<script setup lang="ts">
import {
  ChevronDownOutline,
  ExpandOutline,
  MailOutline,
  MenuOutline,
  MoonOutline,
  NotificationsOutline,
  SearchOutline,
  SunnyOutline,
} from '@vicons/ionicons5'
import type { DropdownOption } from 'naive-ui'
import { NBadge, NButton, NDropdown, NIcon, NInput, NLayoutHeader, NTooltip } from 'naive-ui'
import { computed } from 'vue'
import { useThemeStore } from '../../stores/theme'
import { useNotificationStore } from '../../stores/notification'
import NotificationDrawer from './NotificationDrawer.vue'

defineProps<{
  breadcrumbText: string
  displayName: string
  dropdownOptions: DropdownOption[]
}>()

const emit = defineEmits<{
  toggleSidebar: []
  userAction: [key: string | number]
}>()

const themeStore = useThemeStore()
const notificationStore = useNotificationStore()

const themeIcon = computed(() => (themeStore.isDark ? SunnyOutline : MoonOutline))

const themeTooltip = computed(() => {
  const labels: Record<string, string> = { light: '浅色', dark: '深色', auto: '跟随系统' }
  const next = { light: 'dark', dark: 'auto', auto: 'light' } as const
  const nextLabel = labels[next[themeStore.mode]]
  return `当前: ${labels[themeStore.mode]}（点击切换到${nextLabel}）`
})

function toggleFullscreen() {
  if (document.fullscreenElement) {
    void document.exitFullscreen()
    return
  }
  void document.documentElement.requestFullscreen()
}
</script>

<template>
  <NLayoutHeader bordered class="flex h-14 items-center justify-between bg-[var(--ez-header-bg)] px-4">
    <div class="flex min-w-0 items-center gap-3">
      <NButton quaternary circle class="h-9 w-9 shrink-0 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)]" @click="emit('toggleSidebar')">
        <NIcon :component="MenuOutline" :size="19" />
      </NButton>
      <p class="truncate text-[var(--ez-text-sm)] font-medium text-[var(--ez-text-secondary)]">{{ breadcrumbText }}</p>
    </div>

    <div class="flex items-center gap-1">
      <NInput
        placeholder="搜索菜单 / 页面 / 功能..."
        clearable
        size="small"
        class="w-[320px] max-[1100px]:hidden"
      >
        <template #prefix>
          <NIcon :component="SearchOutline" :size="16" class="text-[var(--ez-text-light)]" />
        </template>
      </NInput>

      <NTooltip trigger="hover">
        <template #trigger>
          <NButton quaternary circle class="h-9 w-9 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)] hover:!text-[var(--ez-text-main)]">
            <NIcon :component="MailOutline" :size="18" />
          </NButton>
        </template>
        消息
      </NTooltip>

      <NBadge :value="notificationStore.unreadCount" :max="99">
        <NButton quaternary circle class="h-9 w-9 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)] hover:!text-[var(--ez-text-main)]" @click="notificationStore.openDrawer()">
          <NIcon :component="NotificationsOutline" :size="18" />
        </NButton>
      </NBadge>

      <NotificationDrawer />

      <NTooltip trigger="hover">
        <template #trigger>
          <NButton quaternary circle class="h-9 w-9 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)] hover:!text-[var(--ez-text-main)]" @click="themeStore.cycleMode()">
            <NIcon :component="themeIcon" :size="18" />
          </NButton>
        </template>
        {{ themeTooltip }}
      </NTooltip>

      <NTooltip trigger="hover">
        <template #trigger>
          <NButton quaternary circle class="h-9 w-9 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)] hover:!text-[var(--ez-text-main)]" @click="toggleFullscreen">
            <NIcon :component="ExpandOutline" :size="18" />
          </NButton>
        </template>
        全屏
      </NTooltip>

      <NDropdown trigger="click" :options="dropdownOptions" @select="(key) => emit('userAction', key)">
        <NButton
          quaternary
          class="ml-1 gap-1.5 rounded-[var(--ez-radius-control)] border border-[var(--ez-border)] px-3 py-1.5 text-[var(--ez-text-main)] transition-colors hover:border-[var(--ez-primary)] hover:bg-[var(--ez-page-bg)]"
        >
          <span class="text-[var(--ez-text-sm)] font-medium">{{ displayName }}</span>
          <NIcon :component="ChevronDownOutline" :size="14" class="text-[var(--ez-text-light)]" />
        </NButton>
      </NDropdown>
    </div>
  </NLayoutHeader>
</template>
