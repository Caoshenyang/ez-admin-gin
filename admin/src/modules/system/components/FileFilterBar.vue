<script setup lang="ts">
import { NButton, NCard, NInput, NSelect, NSpace } from 'naive-ui'

import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import type { FileListQuery } from '../types/file'

defineProps<{
  extFilterOptions: { label: string; value: string }[]
}>()

defineEmits<{
  reset: []
  search: []
}>()

const query = defineModel<FileListQuery>('query', { required: true })
</script>

<template>
  <NCard :bordered="false" class="rounded-lg">
    <NSpace align="center" :wrap="true">
      <NInput v-model:value="query.keyword" clearable placeholder="文件名" class="w-56" @keyup.enter="$emit('search')" />
      <NSelect v-model:value="query.ext" :options="extFilterOptions" class="w-36" />
      <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" class="w-36" />
      <NButton type="primary" @click="$emit('search')">查询</NButton>
      <NButton @click="$emit('reset')">重置</NButton>
    </NSpace>
  </NCard>
</template>
