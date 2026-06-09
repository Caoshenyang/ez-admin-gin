<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'

import EzDataTable from '@/components/ez/EzDataTable.vue'
import type { MessageTemplateItem, MessageTemplateListQuery } from '../types/message'

defineProps<{
  columns: DataTableColumns<MessageTemplateItem>
  items: MessageTemplateItem[]
  loading: boolean
  query: MessageTemplateListQuery
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
    :row-key="(row: MessageTemplateItem) => row.id"
    :total="total"
    @page-change="$emit('pageChange', $event)"
    @page-size-change="$emit('pageSizeChange', $event)"
    @refresh="$emit('refresh')"
  />
</template>
