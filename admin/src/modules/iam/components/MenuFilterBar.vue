<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NButton, NInput, NSelect } from 'naive-ui'

import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
import { MenuStatus, MenuType } from '@/modules/iam/types/menu'

const props = defineProps<{
  keyword: string
  status: 0 | MenuStatus
  statusOptions: SelectOption[]
  type: 0 | MenuType
  typeOptions: SelectOption[]
}>()

const emit = defineEmits<{
  reset: []
  search: []
  'update:keyword': [value: string]
  'update:status': [value: 0 | MenuStatus]
  'update:type': [value: 0 | MenuType]
}>()
</script>

<template>
  <EzSearchPanel>
    <NInput
      :value="props.keyword"
      clearable
      placeholder="菜单名称 / 路由 / 权限标识"
      class="ez-search-field ez-search-field--primary"
      @update:value="emit('update:keyword', $event ?? '')"
      @keyup.enter="emit('search')"
    />
    <NSelect :value="props.type" :options="typeOptions" class="ez-search-field ez-search-field--md" @update:value="emit('update:type', $event as 0 | MenuType)" />
    <NSelect :value="props.status" :options="statusOptions" class="ez-search-field ez-search-field--sm" @update:value="emit('update:status', $event as 0 | MenuStatus)" />

    <template #actions>
      <NButton type="primary" @click="emit('search')">查询</NButton>
      <NButton @click="emit('reset')">重置</NButton>
    </template>
  </EzSearchPanel>
</template>
