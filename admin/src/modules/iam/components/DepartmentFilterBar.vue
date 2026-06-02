<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NButton, NInput, NSelect } from 'naive-ui'

import EzSearchPanel from '@/components/ez/EzSearchPanel.vue'
import { DepartmentStatus } from '@/modules/iam/types/department'

const props = defineProps<{
  keyword: string
  status: 0 | DepartmentStatus
  statusOptions: SelectOption[]
}>()

const emit = defineEmits<{
  reset: []
  search: []
  'update:keyword': [value: string]
  'update:status': [value: 0 | DepartmentStatus]
}>()
</script>

<template>
  <EzSearchPanel>
    <NInput
      :value="props.keyword"
      clearable
      placeholder="搜索部门名称 / 编码"
      class="ez-search-field ez-search-field--primary"
      @update:value="emit('update:keyword', $event ?? '')"
      @keyup.enter="emit('search')"
    />
    <NSelect :value="props.status" :options="statusOptions" class="ez-search-field ez-search-field--sm" @update:value="emit('update:status', $event as 0 | DepartmentStatus)" />

    <template #actions>
      <NButton type="primary" @click="emit('search')">查询</NButton>
      <NButton @click="emit('reset')">重置</NButton>
    </template>
  </EzSearchPanel>
</template>
