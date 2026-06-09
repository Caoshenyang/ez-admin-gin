<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'

import EzDataTable from '@/components/ez/EzDataTable.vue'
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
  <EzDataTable
    remote
    :columns="columns"
    :data="items"
    :loading="loading"
    :page="query.page"
    :page-size="query.page_size"
    :row-key="(row: ConfigItem) => row.id"
    :total="total"
    @page-change="$emit('pageChange', $event)"
    @page-size-change="$emit('pageSizeChange', $event)"
    @refresh="$emit('refresh')"
  />
</template>
