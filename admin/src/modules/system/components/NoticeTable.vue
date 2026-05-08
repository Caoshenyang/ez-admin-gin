<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NPagination } from 'naive-ui'

import type { NoticeItem, NoticeListQuery } from '../types/notice'

defineProps<{
  columns: DataTableColumns<NoticeItem>
  items: NoticeItem[]
  loading: boolean
  query: NoticeListQuery
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
      <span class="text-sm text-[#6B7280]">共 {{ total }} 条</span>
      <NButton text type="primary" @click="$emit('refresh')">刷新</NButton>
    </div>

    <NDataTable
      remote
      class="notice-table"
      :columns="columns"
      :data="items"
      :loading="loading"
      :pagination="false"
      :row-key="(row: NoticeItem) => row.id"
      :bordered="false"
    />

    <div class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]">
      <span>共 {{ total }} 条</span>
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
.notice-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #4b5563;
  background: #f9fafb;
  font-size: 13px;
}

.notice-table :deep(.n-data-table-td) {
  color: #374151;
  font-size: 14px;
  padding: 10px 16px;
}

.notice-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: unset !important;
}

.notice-table :deep(.n-data-table-tr:hover) {
  filter: brightness(0.97);
}
</style>
