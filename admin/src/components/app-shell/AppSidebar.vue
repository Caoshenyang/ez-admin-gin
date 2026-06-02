<!-- AppSidebar 渲染侧栏菜单树，高亮当前路由并支持展开/折叠。 -->
<script setup lang="ts">
import { ChevronBackOutline, ChevronForwardOutline } from '@vicons/ionicons5'
import type { MenuOption, ScrollbarProps } from 'naive-ui'
import { NButton, NIcon, NLayoutSider, NMenu, NScrollbar, NTooltip } from 'naive-ui'

import brandLogoMarkLightUrl from '@/assets/brand-logo-mark-light.svg'

defineProps<{
  activeMenuKey: string
  collapsed: boolean
  expandedMenuKeys: string[]
  menuOptions: MenuOption[]
}>()

const emit = defineEmits<{
  navigate: [path: string]
  select: [key: string | number]
  expand: [keys: Array<string | number>]
  toggle: []
}>()

const sidebarScrollbarThemeOverrides: NonNullable<ScrollbarProps['themeOverrides']> = {
  width: '4px',
  height: '4px',
  borderRadius: '999px',
  color: 'var(--ez-sidebar-scrollbar)',
  colorHover: 'var(--ez-sidebar-scrollbar-hover)',
  railColor: 'transparent',
}
</script>

<template>
  <NLayoutSider
    collapse-mode="width"
    bordered
    content-class="flex h-full flex-col bg-[linear-gradient(180deg,var(--ez-sidebar-bg)_0%,var(--ez-sidebar-bg-deep)_100%)]"
    :collapsed="collapsed"
    :collapsed-width="64"
    :width="248"
    :native-scrollbar="false"
    :show-trigger="false"
    inverted
  >
    <div class="flex h-14 items-center justify-center px-3">
      <NButton text class="!h-auto !p-0 !text-white hover:!bg-transparent" @click="emit('navigate', '/dashboard')">
        <div
          class="flex items-center"
          :class="collapsed ? 'justify-center' : 'w-full justify-start gap-3'"
        >
          <img
            :src="brandLogoMarkLightUrl"
            :width="collapsed ? 28 : 32"
            alt="EZ Admin Gin 品牌图形"
            class="block h-auto shrink-0"
          >
          <div v-if="!collapsed" class="min-w-0 text-left leading-none">
            <p class="truncate text-[15px] font-semibold text-white">EZ Admin Gin</p>
          </div>
        </div>
      </NButton>
    </div>

    <p v-if="!collapsed" class="px-4 text-[11px] font-semibold tracking-[0.08em] text-[#64748B] uppercase">
      工作台
    </p>

    <NScrollbar
      class="sidebar-menu-scrollbar mt-2 min-h-0 flex-1"
      :class="collapsed ? 'px-0' : 'px-3'"
      trigger="hover"
      :theme-overrides="sidebarScrollbarThemeOverrides"
    >
      <NMenu
        :value="activeMenuKey"
        :expanded-keys="expandedMenuKeys"
        :options="menuOptions"
        :indent="20"
        :root-indent="12"
        :collapsed="collapsed"
        :collapsed-width="64"
        :collapsed-icon-size="22"
        :accordion="false"
        inverted
        @update:value="(key) => emit('select', key)"
        @update:expanded-keys="(keys) => emit('expand', keys)"
      />
    </NScrollbar>

    <!-- 底部收起/展开按钮 -->
    <div
      class="flex items-center justify-center border-t border-white/6 py-3"
      :class="collapsed ? 'px-0' : 'px-3'"
    >
      <NTooltip v-if="collapsed" placement="right" :delay="300">
        <template #trigger>
          <NButton
            quaternary
            circle
            class="h-9 w-9 shrink-0 rounded-[var(--ez-radius-control)] !bg-white/6 !p-0 !text-white/45 hover:!bg-white/12 hover:!text-white/85"
            @click="emit('toggle')"
          >
            <NIcon :component="ChevronForwardOutline" :size="18" />
          </NButton>
        </template>
        展开菜单
      </NTooltip>
      <NButton
        v-else
        quaternary
        class="flex w-full items-center justify-center gap-2 rounded-[var(--ez-radius-control)] !bg-white/6 px-0 py-2 !text-white/45 hover:!bg-white/12 hover:!text-white/85"
        @click="emit('toggle')"
      >
        <NIcon :component="ChevronBackOutline" :size="18" />
        <span class="text-[var(--ez-text-xs)] font-medium">收起菜单</span>
      </NButton>
    </div>
  </NLayoutSider>
</template>
