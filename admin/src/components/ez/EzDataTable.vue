<script setup lang="ts" generic="T extends object">
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import { NButton, NDataTable, NPagination } from 'naive-ui'
import { computed } from 'vue'

import EzTableCard from './EzTableCard.vue'
import TableStatsBar from '../TableStatsBar.vue'

const props = withDefaults(
  defineProps<{
    columns: DataTableColumns<T>
    data: T[]
    loading: boolean
    rowKey?: (row: T) => DataTableRowKey
    rowProps?: (row: T) => Record<string, unknown>
    remote?: boolean
    summaryText?: string
    page?: number
    pageSize?: number
    total?: number
    pageSizes?: number[]
    showSizePicker?: boolean
  }>(),
  {
    remote: false,
    rowKey: undefined,
    rowProps: undefined,
    summaryText: '',
    page: undefined,
    pageSize: undefined,
    total: undefined,
    pageSizes: () => [10, 20, 50],
    showSizePicker: true,
  },
)

const emit = defineEmits<{
  pageChange: [page: number]
  pageSizeChange: [pageSize: number]
  refresh: []
}>()

defineSlots<{
  actions?: () => unknown
  body?: () => unknown
  footerSummary?: () => unknown
  toolbarSummary?: () => unknown
}>()

const hasPagination = computed(() => {
  return props.page !== undefined && props.pageSize !== undefined && props.total !== undefined
})

const displaySummary = computed(() => {
  return props.summaryText || `共 ${props.total ?? props.data.length} 条`
})
</script>

<template>
  <EzTableCard>
    <TableStatsBar>
      <slot name="toolbarSummary">
        <span>{{ displaySummary }}</span>
      </slot>
      <template #actions>
        <slot name="actions">
          <NButton text type="primary" @click="emit('refresh')">刷新</NButton>
        </slot>
      </template>
    </TableStatsBar>

    <div class="ez-table-fill">
      <!--
        普通列表只传 columns/data/rowKey 即可；需要勾选、行样式或额外事件时，
        用 body 插槽接管 NDataTable，仍复用统一工具条和分页外壳。
      -->
      <slot name="body">
        <NDataTable
          :remote="remote"
          class="ez-table-fill-table"
          :columns="columns"
          :data="data"
          :loading="loading"
          :pagination="false"
          :row-key="rowKey"
          :row-props="rowProps"
          :bordered="false"
          flex-height
        />
      </slot>

      <div v-if="hasPagination" class="ez-table-footer">
        <slot name="footerSummary">
          <span>{{ displaySummary }}</span>
        </slot>
        <NPagination
          :page="page"
          :page-size="pageSize"
          :item-count="total"
          :page-sizes="pageSizes"
          :show-size-picker="showSizePicker"
          @update:page="emit('pageChange', $event)"
          @update:page-size="emit('pageSizeChange', $event)"
        />
      </div>
    </div>
  </EzTableCard>
</template>
