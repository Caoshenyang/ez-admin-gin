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
  <NCard class="min-h-0 rounded-lg" :bordered="false" content-style="display: flex; height: 100%; min-height: 0; flex-direction: column; padding: 0;">
    <div class="dict-card-shell">
      <div class="dict-card-shell__header">
        <div>
          <p class="dict-card-shell__eyebrow">Types</p>
          <h2 class="dict-card-shell__title">字典类型</h2>
        </div>
        <NButton v-if="canUse('system:dict:type:create')" size="small" type="primary" ghost @click="$emit('create')">
          新增
        </NButton>
      </div>

      <div class="dict-card-shell__filters">
        <NInput v-model:value="query.keyword" clearable placeholder="编码 / 名称" @keyup.enter="$emit('search')" />
        <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" />
        <div class="dict-filter-actions">
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
        style="height: 100%;"
        flex-height
      />

      <div class="dict-card-shell__footer">
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
.dict-card-shell {
  display: grid;
  min-height: 0;
  height: 100%;
  grid-template-rows: auto auto minmax(0, 1fr) auto;
}

.dict-card-shell__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  border-bottom: 1px solid #e5e7eb;
  padding: 18px 20px 14px;
}

.dict-card-shell__eyebrow {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #94a3b8;
}

.dict-card-shell__title {
  margin-top: 6px;
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.dict-card-shell__filters {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 140px auto;
  gap: 12px;
  padding: 16px 20px;
}

.dict-filter-actions {
  display: flex;
  gap: 10px;
}

.dict-card-shell__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  border-top: 1px solid #e5e7eb;
  padding: 14px 20px;
  font-size: 13px;
  color: #6b7280;
}

.dict-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #374151;
  background: #fff;
}

.dict-table :deep(.n-data-table-td) {
  color: #374151;
}

.dict-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}

:deep(.dict-type-row) {
  cursor: pointer;
}

:deep(.dict-type-row--active .n-data-table-td) {
  background: #eef6ff;
}

@media (max-width: 1280px) {
  .dict-card-shell__filters {
    grid-template-columns: minmax(0, 1fr);
  }

  .dict-filter-actions {
    justify-content: flex-end;
  }
}

@media (max-width: 720px) {
  .dict-card-shell__header,
  .dict-card-shell__filters,
  .dict-card-shell__footer {
    padding-left: 16px;
    padding-right: 16px;
  }

  .dict-card-shell__footer {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
