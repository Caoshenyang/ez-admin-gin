<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NDataTable, NPagination } from 'naive-ui'

import EzTableCard from '@/components/ez/EzTableCard.vue'
import TableStatsBar from '@/components/TableStatsBar.vue'
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
  <EzTableCard>
    <TableStatsBar>
      <span>共 {{ total }} 个文件</span>
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
      :row-key="(row: FileItem) => row.id"
      :bordered="false"
    />

    <div class="ez-table-footer">
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
  </EzTableCard>
</template>
