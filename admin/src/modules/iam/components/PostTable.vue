<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable } from 'naive-ui'

import type { PostItem } from '../types/post'

defineProps<{
  columns: DataTableColumns<PostItem>
  items: PostItem[]
  loading: boolean
}>()

defineEmits<{
  refresh: []
}>()
</script>

<template>
  <NCard class="ez-table-card min-h-0 flex-1" :bordered="false" content-class="ez-card-content-fill">
    <div class="flex items-center justify-end border-b border-[var(--ez-border)] px-4 py-3">
      <NButton text type="primary" @click="$emit('refresh')">刷新</NButton>
    </div>

    <NDataTable
      class="h-full"
      :columns="columns"
      :data="items"
      :loading="loading"
      :pagination="false"
      :bordered="false"
      flex-height
    />
  </NCard>
</template>
