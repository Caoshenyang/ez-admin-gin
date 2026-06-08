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
    <div class="dict-type-panel-body grid h-full min-h-0 min-w-0 grid-rows-[auto_auto_minmax(0,1fr)_auto]">
      <div class="flex items-start justify-between gap-4 border-b border-[var(--ez-border-light)] px-5 pt-[18px] pb-[14px] max-[720px]:px-4">
        <div>
          <p class="text-[var(--ez-text-xs)] font-bold tracking-[0.14em] text-[var(--ez-text-light)] uppercase">Types</p>
          <h2 class="mt-1.5 text-[var(--ez-text-xl)] font-bold text-[var(--ez-text-main)]">字典类型</h2>
        </div>
        <NButton v-if="canUse('system:dict:type:create')" size="small" type="primary" ghost @click="$emit('create')">
          新增
        </NButton>
      </div>

      <div class="dict-type-filter-grid px-5 py-4 max-[720px]:px-4">
        <NInput v-model:value="query.keyword" clearable placeholder="编码 / 名称" @keyup.enter="$emit('search')" />
        <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" />
        <div class="dict-type-filter-actions">
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
        flex-height
      />

      <div class="ez-table-footer">
        <span>共 {{ total }} 条</span>
        <NPagination
          :page="query.page"
          :page-size="query.page_size"
          :item-count="total"
          :page-slot="3"
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

.dict-type-filter-grid {
  display: grid;
  width: 100%;
  max-width: 100%;
  min-width: 0;
  grid-template-columns: minmax(0, 1fr);
  gap: 10px;
}

.dict-type-filter-actions {
  display: flex;
  grid-column: 1 / -1;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: flex-start;
}

.dict-type-filter-actions :deep(.n-button) {
  min-width: 64px;
}

:deep(.n-data-table-th--fixed-right),
:deep(.n-data-table-td--fixed-right) {
  box-shadow: -8px 0 12px rgba(15, 23, 42, 0.04);
}

.dict-type-panel-body {
  container: dict-type-panel / inline-size;
}

@container dict-type-panel (min-width: 460px) {
  .dict-type-filter-grid {
    grid-template-columns: minmax(180px, 1fr) 118px max-content;
  }

  .dict-type-filter-actions {
    grid-column: auto;
    justify-content: flex-end;
  }
}
</style>
