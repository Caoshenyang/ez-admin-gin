<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NInput, NPagination, NSelect } from 'naive-ui'

import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import type { DictTypeItem, DictTypeListQuery } from '../types/dict'

defineProps<{
  canUse: (code: string) => boolean
  columns: DataTableColumns<DictTypeItem>
  items: DictTypeItem[]
  loading: boolean
  rowProps: (row: DictTypeItem) => Record<string, unknown>
  total: number
}>()

defineEmits<{
  create: []
  pageChange: [page: number]
  pageSizeChange: [pageSize: number]
  reset: []
  search: []
}>()

const query = defineModel<DictTypeListQuery>('query', { required: true })
</script>

<template>
  <NCard class="ez-card h-full min-h-0 min-w-0 overflow-hidden rounded-[var(--ez-radius-sm)]" :bordered="false" content-class="ez-card-content-fill">
    <div class="grid h-full min-h-0 grid-rows-[auto_auto_minmax(0,1fr)_auto]">
      <div class="flex items-start justify-between gap-4 border-b border-[var(--ez-border-light)] px-5 pt-[18px] pb-[14px] max-[720px]:px-4">
        <div>
          <p class="text-[var(--ez-text-xs)] font-bold tracking-[0.14em] text-[var(--ez-text-light)] uppercase">Types</p>
          <h2 class="mt-1.5 text-[var(--ez-text-xl)] font-bold text-[var(--ez-text-main)]">字典类型</h2>
        </div>
        <NButton v-if="canUse('system:dict:type:create')" size="small" type="primary" ghost @click="$emit('create')">
          新增
        </NButton>
      </div>

      <div class="grid items-center gap-3 px-5 py-4 min-[1281px]:grid-cols-[minmax(0,1fr)_132px_auto] max-[720px]:px-4">
        <NInput v-model:value="query.keyword" clearable placeholder="编码 / 名称" @keyup.enter="$emit('search')" />
        <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" />
        <div class="ez-filter-actions">
          <NButton type="primary" @click="$emit('search')">查询</NButton>
          <NButton @click="$emit('reset')">重置</NButton>
        </div>
      </div>

      <NDataTable
        remote
        class="dict-table h-full min-h-0"
        :columns="columns"
        :data="items"
        :loading="loading"
        :pagination="false"
        :row-key="(row: DictTypeItem) => row.id"
        :row-props="rowProps"
        :bordered="false"
        :scroll-x="500"
        flex-height
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
    </div>
  </NCard>
</template>

<style scoped>
:deep(.dict-type-row) {
  cursor: pointer;
}

:deep(.dict-type-row--active .n-data-table-td) {
  background: var(--ez-brand-soft);
}
</style>
