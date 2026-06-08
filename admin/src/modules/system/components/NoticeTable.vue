<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'

import EzDataTable from '@/components/ez/EzDataTable.vue'
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
  <EzDataTable
    remote
    :columns="columns"
    :data="items"
    :loading="loading"
    :page="query.page"
    :page-size="query.page_size"
    :row-key="(row: NoticeItem) => row.id"
    :total="total"
    @page-change="$emit('pageChange', $event)"
    @page-size-change="$emit('pageSizeChange', $event)"
    @refresh="$emit('refresh')"
  />
</template>
