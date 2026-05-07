<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NPagination } from 'naive-ui'

import type { ConfigItem, ConfigListQuery } from '../types/config'

defineProps<{
  columns: DataTableColumns<ConfigItem>
  items: ConfigItem[]
  loading: boolean
  query: ConfigListQuery
  total: number
}>()

defineEmits<{
  pageChange: [page: number]
  pageSizeChange: [pageSize: number]
  refresh: []
}>()
</script>

<template>
  <NCard class="min-h-0 flex-1 rounded-lg" :bordered="false" content-style="height: 100%; padding: 0;">
    <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-3">
      <span class="text-sm text-[#6B7280]">共 {{ total }} 条</span>
      <NButton text type="primary" @click="$emit('refresh')">刷新</NButton>
    </div>

    <NDataTable
      remote
      class="config-table h-full"
      style="height: calc(100% - 105px)"
      :columns="columns"
      :data="items"
      :loading="loading"
      :pagination="false"
      :row-key="(row: ConfigItem) => row.id"
      :bordered="false"
      flex-height
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
.config-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #374151;
  background: #fff;
}

.config-table :deep(.n-data-table-td) {
  color: #374151;
}

.config-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}
</style>
