<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NButton, NInput, NSelect } from 'naive-ui'

import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'

defineProps<{
  bizType: string
  category: string
  ext: string
  extOptions: SelectOption[]
  keyword: string
  status: number
  statusOptions: SelectOption[]
}>()

defineEmits<{
  'update:bizType': [value: string]
  'update:category': [value: string]
  'update:ext': [value: string]
  'update:keyword': [value: string]
  'update:status': [value: number]
  search: []
  reset: []
}>()
</script>

<template>
  <EzSearchPanel>
    <NInput :value="keyword" clearable placeholder="附件名称 / 原始文件名" class="ez-search-field ez-search-field--primary" @update:value="(value) => $emit('update:keyword', value)" @keyup.enter="$emit('search')" />
    <NInput :value="category" clearable placeholder="附件分类" class="ez-search-field ez-search-field--sm" @update:value="(value) => $emit('update:category', value)" @keyup.enter="$emit('search')" />
    <NInput :value="bizType" clearable placeholder="业务类型" class="ez-search-field ez-search-field--sm" @update:value="(value) => $emit('update:bizType', value)" @keyup.enter="$emit('search')" />
    <NSelect :value="ext" :options="extOptions" class="ez-search-field ez-search-field--xs" @update:value="(value) => $emit('update:ext', String(value ?? ''))" />
    <NSelect :value="status" :options="statusOptions" class="ez-search-field ez-search-field--xs" @update:value="(value) => $emit('update:status', Number(value ?? 0))" />

    <template #actions>
      <NButton type="primary" @click="$emit('search')">查询</NButton>
      <NButton @click="$emit('reset')">重置</NButton>
    </template>
  </EzSearchPanel>
</template>
