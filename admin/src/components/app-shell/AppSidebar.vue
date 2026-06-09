<!-- AppSidebar 渲染侧栏菜单树，高亮当前路由并支持展开/折叠。 -->
<script setup lang="ts">
import { ChevronBackOutline, ChevronForwardOutline } from '@vicons/ionicons5'
import type { ButtonProps, MenuNodeProps, MenuOption, MenuProps, ScrollbarProps } from 'naive-ui'
import { NButton, NIcon, NLayoutSider, NMenu, NScrollbar, NTooltip } from 'naive-ui'
import { nextTick, ref } from 'vue'

import brandLogoMarkLightUrl from '@/assets/brand-logo-mark-light.svg'

const props = defineProps<{
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

const sidebarMenuThemeOverrides: NonNullable<MenuProps['themeOverrides']> = {
  itemColorActiveInverted: 'rgba(37, 99, 255, 0.72)',
  itemColorActiveHoverInverted: 'rgba(37, 99, 255, 0.82)',
  itemColorActiveCollapsedInverted: 'rgba(37, 99, 255, 0.72)',
  itemTextColorActiveInverted: '#ffffff',
  itemTextColorActiveHoverInverted: '#ffffff',
  itemTextColorChildActiveInverted: '#dbeafe',
  itemTextColorChildActiveHoverInverted: '#ffffff',
  itemIconColorActiveInverted: '#ffffff',
  itemIconColorActiveHoverInverted: '#ffffff',
  itemIconColorChildActiveInverted: '#bfdbfe',
  itemIconColorChildActiveHoverInverted: '#ffffff',
  arrowColorChildActiveInverted: '#bfdbfe',
  arrowColorChildActiveHoverInverted: '#ffffff',
}

const sidebarCollapseButtonThemeOverrides: NonNullable<ButtonProps['themeOverrides']> = {
  colorQuaternary:
    'linear-gradient(180deg, rgba(15, 23, 42, 0.5) 0%, rgba(15, 23, 42, 0.3) 100%)',
  colorQuaternaryHover:
    'linear-gradient(180deg, rgba(37, 99, 255, 0.2) 0%, rgba(15, 23, 42, 0.36) 100%)',
  colorQuaternaryPressed:
    'linear-gradient(180deg, rgba(37, 99, 255, 0.28) 0%, rgba(15, 23, 42, 0.42) 100%)',
  border: '1px solid rgba(148, 163, 184, 0.16)',
  borderHover: '1px solid rgba(56, 208, 248, 0.34)',
  borderPressed: '1px solid rgba(56, 208, 248, 0.42)',
  borderFocus: '1px solid rgba(125, 211, 252, 0.44)',
  textColor: 'rgba(255, 255, 255, 0.7)',
  textColorHover: '#ffffff',
  textColorPressed: '#ffffff',
  textColorFocus: '#ffffff',
}

const menuContentRef = ref<HTMLElement | null>(null)

const resolveMenuNodeProps: MenuNodeProps = (option) => ({
  'data-menu-key': option.key === undefined ? undefined : String(option.key),
})

function findMenuItemElement(key: string) {
  const menuItems = menuContentRef.value?.querySelectorAll<HTMLElement>('[data-menu-key]')
  return Array.from(menuItems ?? []).find((item) => item.dataset.menuKey === key) ?? null
}

function scrollExpandedMenuIntoView(key: string) {
  const menuItem = findMenuItemElement(key)
  const submenu = menuItem?.closest<HTMLElement>('.n-submenu')
  const firstChildItem = submenu?.querySelector<HTMLElement>('.n-submenu-children .n-menu-item')

  ;(firstChildItem ?? menuItem)?.scrollIntoView({
    block: 'nearest',
    behavior: 'smooth',
  })
}

function scheduleExpandedMenuScroll(key: string) {
  void nextTick(() => {
    window.requestAnimationFrame(() => {
      scrollExpandedMenuIntoView(key)
      window.setTimeout(() => scrollExpandedMenuIntoView(key), 180)
    })
  })
}

function handleExpandedKeysUpdate(keys: Array<string | number>) {
  const normalizedKeys = keys.map(String)
  const expandedKey = normalizedKeys.find((key) => !props.expandedMenuKeys.includes(key))

  emit('expand', keys)

  if (!props.collapsed && expandedKey) {
    scheduleExpandedMenuScroll(expandedKey)
  }
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
      <NButton
        text
        class="!h-auto !p-0 !text-white hover:!bg-transparent"
        @click="emit('navigate', '/dashboard')"
      >
        <div
          class="flex items-center"
          :class="collapsed ? 'justify-center' : 'w-full justify-start gap-3'"
        >
          <img
            :src="brandLogoMarkLightUrl"
            :width="collapsed ? 28 : 32"
            alt="EZ Admin Gin 品牌图形"
            class="block h-auto shrink-0"
          />
          <div v-if="!collapsed" class="min-w-0 text-left leading-none">
            <p class="truncate text-[15px] font-semibold text-white">EZ Admin Gin</p>
          </div>
        </div>
      </NButton>
    </div>

    <p
      v-if="!collapsed"
      class="px-4 text-[11px] font-semibold tracking-[0.08em] text-[var(--ez-sidebar-text)] uppercase"
    >
      工作台
    </p>

    <NScrollbar
      class="sidebar-menu-scrollbar mt-2 min-h-0 flex-1"
      :class="collapsed ? 'px-0' : 'px-3'"
      trigger="hover"
      :theme-overrides="sidebarScrollbarThemeOverrides"
    >
      <div ref="menuContentRef" class="sidebar-menu-content">
        <NMenu
          :value="activeMenuKey"
          :expanded-keys="expandedMenuKeys"
          :options="menuOptions"
          :indent="20"
          :root-indent="12"
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :accordion="true"
          :node-props="resolveMenuNodeProps"
          :theme-overrides="sidebarMenuThemeOverrides"
          inverted
          @update:value="(key) => emit('select', key)"
          @update:expanded-keys="handleExpandedKeysUpdate"
        />
      </div>
    </NScrollbar>

    <!-- 底部收起/展开按钮 -->
    <div class="sidebar-collapse-zone" :class="{ 'sidebar-collapse-zone--collapsed': collapsed }">
      <NTooltip v-if="collapsed" placement="right" :delay="300">
        <template #trigger>
          <NButton
            quaternary
            aria-label="展开菜单"
            class="sidebar-collapse-button sidebar-collapse-button--icon"
            :theme-overrides="sidebarCollapseButtonThemeOverrides"
            @click="emit('toggle')"
          >
            <span class="sidebar-collapse-button__icon">
              <NIcon :component="ChevronForwardOutline" :size="18" />
            </span>
          </NButton>
        </template>
        展开菜单
      </NTooltip>
      <NButton
        v-else
        quaternary
        aria-label="收起菜单"
        class="sidebar-collapse-button sidebar-collapse-button--full"
        :theme-overrides="sidebarCollapseButtonThemeOverrides"
        @click="emit('toggle')"
      >
        <span class="sidebar-collapse-button__icon">
          <NIcon :component="ChevronBackOutline" :size="18" />
        </span>
        <span class="sidebar-collapse-button__text">收起菜单</span>
      </NButton>
    </div>
  </NLayoutSider>
</template>

<style scoped>
.sidebar-collapse-zone {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
  background: linear-gradient(180deg, rgba(7, 17, 31, 0) 0%, rgba(2, 6, 23, 0.28) 100%);
}

.sidebar-collapse-zone::before {
  position: absolute;
  top: 0;
  right: 12px;
  left: 12px;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(148, 163, 184, 0.22), transparent);
  content: '';
}

.sidebar-collapse-zone--collapsed {
  padding-right: 0;
  padding-left: 0;
}

.sidebar-collapse-zone--collapsed::before {
  right: 10px;
  left: 10px;
}

.sidebar-collapse-button {
  height: 40px;
  border-radius: var(--ez-radius-control);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.07),
    0 8px 18px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.18s ease;
}

.sidebar-collapse-button:hover {
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 10px 22px rgba(2, 6, 23, 0.18);
}

.sidebar-collapse-button:focus-visible {
  outline: 2px solid rgba(125, 211, 252, 0.5);
  outline-offset: 2px;
}

.sidebar-collapse-button--full {
  width: 100%;
  padding: 0 9px;
}

.sidebar-collapse-button--icon {
  width: 40px;
  min-width: 40px;
  padding: 0;
}

.sidebar-collapse-button :deep(.n-button__content) {
  width: 100%;
  min-width: 0;
}

.sidebar-collapse-button--full :deep(.n-button__content) {
  justify-content: flex-start;
  gap: 10px;
}

.sidebar-collapse-button--icon :deep(.n-button__content) {
  justify-content: center;
}

.sidebar-collapse-button__icon {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.82);
  transition:
    background 0.18s ease,
    color 0.18s ease;
}

.sidebar-collapse-button:hover .sidebar-collapse-button__icon {
  background: rgba(56, 208, 248, 0.16);
  color: #7dd3fc;
}

.sidebar-collapse-button__text {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  font-size: var(--ez-text-xs);
  font-weight: 600;
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
