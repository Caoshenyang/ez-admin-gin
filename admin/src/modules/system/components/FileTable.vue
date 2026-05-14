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
  <NCard class="ez-table-card min-h-0 flex-1" :bordered="false" content-class="ez-card-content-reset">
    <div class="flex items-center justify-between border-b border-[var(--ez-border)] px-4 py-3">
      <span class="text-sm text-[var(--ez-text-sub)]">共 {{ total }} 个文件</span>
      <NButton text type="primary" @click="$emit('refresh')">刷新</NButton>
    </div>

    <NDataTable
      remote
      :columns="columns"
      :data="items"
      :loading="loading"
      :pagination="false"
      :row-key="(row: FileItem) => row.id"
      :bordered="false"
    />

    <div class="flex items-center justify-between border-t border-[var(--ez-border)] px-4 py-3 text-sm text-[var(--ez-text-sub)]">
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
