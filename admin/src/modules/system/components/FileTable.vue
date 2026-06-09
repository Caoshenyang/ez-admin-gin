<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'

import EzDataTable from '@/components/ez/EzDataTable.vue'
import type { FileItem, FileListQuery } from '../types/file'

defineProps<{
  columns: DataTableColumns<FileItem>
  items: FileItem[]
  loading: boolean
  query: FileListQuery
  total: number
}>()

defineEmits<{
  pageChange: [page: number]
  pageSizeChange: [pageSize: number]
  refresh: []
}>()
</script>

<template>
  <EzDataTable
    remote
    :columns="columns"
    :data="items"
    :loading="loading"
    :page="query.page"
    :page-size="query.page_size"
    :row-key="(row: FileItem) => row.id"
    :summary-text="`共 ${total} 个文件`"
    :total="total"
    @page-change="$emit('pageChange', $event)"
    @page-size-change="$emit('pageSizeChange', $event)"
    @refresh="$emit('refresh')"
  />
</template>
