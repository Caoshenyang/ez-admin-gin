<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NPagination } from 'naive-ui'

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
  <NCard class="min-h-0 flex-1 rounded-lg" :bordered="false" content-style="padding: 0;">
    <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-3">
      <span class="text-sm text-[#6B7280]">共 {{ total }} 个文件</span>
      <NButton text type="primary" @click="$emit('refresh')">刷新</NButton>
    </div>

    <NDataTable
      remote
      class="file-table"
      :columns="columns"
      :data="items"
      :loading="loading"
      :pagination="false"
      :row-key="(row: FileItem) => row.id"
      :bordered="false"
    />

    <div class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]">
      <span>共 {{ total }} 个文件</span>
      <NPagination
        :page="query.page"
        :page-size="query.page_size"
        :item-count="total"
        :page-sizes="[10, 20, 50]"
        show-size-picker
        @update:page="(page) => $emit('pageChange', page)"
        @update:page-size="(pageSize) => $emit('pageSizeChange', pageSize)"
      />
    </div>
  </NCard>
</template>

<style scoped>
.file-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #374151;
  background: #fff;
}

.file-table :deep(.n-data-table-td) {
  color: #374151;
}

.file-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}
</style>
