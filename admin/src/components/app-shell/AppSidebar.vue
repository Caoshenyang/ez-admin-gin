<!-- AppSidebar 渲染侧栏菜单树，高亮当前路由并支持展开/折叠。 -->
<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import { NLayoutSider, NMenu, NScrollbar } from 'naive-ui'

import BrandLogo from '@/components/brand/BrandLogo.vue'

defineProps<{
  activeMenuKey: string
  expandedMenuKeys: string[]
  menuOptions: MenuOption[]
}>()

const emit = defineEmits<{
  navigate: [path: string]
  select: [key: string | number]
  expand: [keys: Array<string | number>]
}>()
</script>

<template>
  <NLayoutSider
    collapse-mode="width"
    bordered
    show-trigger="bar"
    content-class="flex h-full flex-col"
    :collapsed-width="72"
    :width="240"
    :native-scrollbar="false"
    inverted
  >
    <div class="flex items-center px-4 py-4">
      <button
        type="button"
        class="flex items-center border-none bg-transparent px-0 py-0 text-left text-white"
        @click="emit('navigate', '/dashboard')"
      >
        <BrandLogo :width="42" direction="inline" :show-title="true" variant="dark" />
      </button>
    </div>

    <p class="px-4 text-xs font-semibold tracking-wide text-[#6B7280]">主菜单</p>

    <NScrollbar class="mt-3 min-h-0 flex-1 px-2" trigger="none">
      <NMenu
        :value="activeMenuKey"
        :expanded-keys="expandedMenuKeys"
        :options="menuOptions"
        :indent="18"
        :collapsed-icon-size="20"
        inverted
        @update:value="(key) => emit('select', key)"
        @update:expanded-keys="(keys) => emit('expand', keys)"
      />
    </NScrollbar>
  </NLayoutSider>
</template>
