<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NPagination } from 'naive-ui'

import TableStatsBar from '@/components/TableStatsBar.vue'
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
  <NCard class="ez-table-card min-h-0 flex-1" :bordered="false" content-class="ez-card-content-reset">
    <TableStatsBar>
      <span>共 {{ total }} 条</span>
      <template #actions>
        <NButton text type="primary" @click="$emit('refresh')">刷新</NButton>
      </template>
    </TableStatsBar>

    <NDataTable
      remote
      :columns="columns"
      :data="items"
      :loading="loading"
      :pagination="false"
      :row-key="(row: ConfigItem) => row.id"
      :bordered="false"
    />

    <div class="ez-table-footer">
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
