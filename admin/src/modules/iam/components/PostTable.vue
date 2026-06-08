<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'

import EzDataTable from '@/components/ez/EzDataTable.vue'
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
  <EzDataTable :columns="columns" :data="items" :loading="loading" @refresh="$emit('refresh')">
    <template #toolbarSummary>
      <span />
    </template>
  </EzDataTable>
</template>
