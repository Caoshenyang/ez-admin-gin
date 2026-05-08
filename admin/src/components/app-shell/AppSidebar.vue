<!-- AppSidebar 渲染侧栏菜单树，高亮当前路由并支持展开/折叠。 -->
<script setup lang="ts">
import { ChevronBackOutline, ChevronForwardOutline } from '@vicons/ionicons5'
import type { MenuOption } from 'naive-ui'
import { NIcon, NLayoutSider, NMenu, NScrollbar, NTooltip } from 'naive-ui'

import BrandLogo from '@/components/brand/BrandLogo.vue'

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
</script>

<template>
  <NLayoutSider
    collapse-mode="width"
    bordered
    content-class="flex h-full flex-col app-sidebar"
    :collapsed="collapsed"
    :collapsed-width="72"
    :width="240"
    :native-scrollbar="false"
    :show-trigger="false"
    inverted
  >
    <div class="flex items-center px-4 py-4">
      <button
        type="button"
        class="flex items-center border-none bg-transparent px-0 py-0 text-left text-white"
        @click="emit('navigate', '/dashboard')"
      >
        <BrandLogo :width="42" direction="inline" :show-title="!collapsed" variant="dark" />
      </button>
    </div>

    <p v-if="!collapsed" class="px-4 text-[11px] font-medium tracking-[0.06em] text-white/30">主菜单</p>

    <NScrollbar class="mt-2 min-h-0 flex-1 px-2" trigger="none">
      <NMenu
        :value="activeMenuKey"
        :expanded-keys="expandedMenuKeys"
        :options="menuOptions"
        :indent="18"
        :collapsed="collapsed"
        :collapsed-icon-size="20"
        inverted
        @update:value="(key) => emit('select', key)"
        @update:expanded-keys="(keys) => emit('expand', keys)"
      />
    </NScrollbar>

    <div class="shrink-0 border-t border-white/[0.06] px-3 pb-3 pt-2">
      <NTooltip v-if="collapsed" placement="right" :delay="300">
        <template #trigger>
          <button type="button" class="sidebar-toggle" @click="emit('toggle')">
            <NIcon :component="ChevronForwardOutline" :size="18" />
          </button>
        </template>
        展开菜单
      </NTooltip>
      <button v-else type="button" class="sidebar-toggle" @click="emit('toggle')">
        <NIcon :component="ChevronBackOutline" :size="18" />
        <span class="text-[12px] font-medium">收起菜单</span>
      </button>
    </div>
  </NLayoutSider>
</template>

<style scoped>
.sidebar-toggle {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: none;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.06);
  padding: 8px 0;
  color: rgba(255, 255, 255, 0.45);
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    color 0.2s ease;
}

.sidebar-toggle:hover {
  background: rgba(255, 255, 255, 0.12);
  color: rgba(255, 255, 255, 0.85);
}
</style>

<style>
/* 全局：深藏蓝侧边栏背景 */
.app-sidebar {
  background-color: #071A2F !important;
}
</style>
