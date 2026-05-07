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
  <NLayoutHeader bordered class="flex h-14 items-center justify-between bg-white px-6">
    <p class="text-sm text-[#374151]">{{ breadcrumbText }}</p>

    <div class="flex items-center gap-2.5">
      <NInput placeholder="搜索菜单 / 页面" clearable class="w-46">
        <template #prefix>
          <NIcon :component="SearchOutline" />
        </template>
      </NInput>

      <NButton quaternary circle>
        <template #icon>
          <NIcon :component="NotificationsOutline" />
        </template>
      </NButton>

      <NButton quaternary circle>
        <template #icon>
          <NIcon :component="ExpandOutline" />
        </template>
      </NButton>

      <NButton quaternary circle>
        <template #icon>
          <NIcon :component="MoonOutline" />
        </template>
      </NButton>

      <NDropdown trigger="click" :options="dropdownOptions" @select="(key) => emit('userAction', key)">
        <NButton secondary>
          <template #icon>
            <NIcon :component="ChevronDownOutline" />
          </template>
          {{ displayName }}
        </NButton>
      </NDropdown>
    </div>
  </NLayoutHeader>
</template>
