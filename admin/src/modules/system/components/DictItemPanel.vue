<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NCard, NDataTable, NInput, NPagination, NSelect, NTag } from 'naive-ui'

import EmptyState from '@/components/EmptyState.vue'
import EzActionButton from '@/components/ez/EzActionButton.vue'
import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import { formatTime } from '@/utils/format'
import type { DictItem, DictItemListQuery, DictTypeItem } from '../types/dict'

defineProps<{
  canUse: (code: string) => boolean
  columns: DataTableColumns<DictItem>
  items: DictItem[]
  loading: boolean
  selectedType: DictTypeItem | null
  total: number
}>()

defineEmits<{
  create: []
  pageChange: [page: number]
  pageSizeChange: [pageSize: number]
  reset: []
  search: []
}>()

const query = defineModel<DictItemListQuery>('query', { required: true })
</script>

<template>
  <NCard
    class="ez-card h-full min-h-0 min-w-0 overflow-hidden rounded-[var(--ez-radius-sm)]"
    :bordered="false"
    content-class="ez-card-content-fill"
  >
    <div class="dict-item-panel-body grid h-full min-h-0 min-w-0 grid-rows-[auto_minmax(0,1fr)]">
      <div
        class="flex items-start justify-between gap-4 border-b border-[var(--ez-border-light)] px-5 pt-[18px] pb-[14px] max-[720px]:px-4"
      >
        <div>
          <p
            class="text-[var(--ez-text-xs)] font-bold tracking-[0.14em] text-[var(--ez-text-light)] uppercase"
          >
            Items
          </p>
          <div class="flex items-center gap-2">
            <h2 class="mt-1.5 text-[var(--ez-text-xl)] font-bold text-[var(--ez-text-main)]">
              字典项
            </h2>
            <NTag v-if="selectedType" size="small" type="info" :bordered="false">
              {{ selectedType.name }}
            </NTag>
          </div>
          <p v-if="selectedType" class="mt-1 text-xs text-[var(--ez-text-sub)]">
            {{ selectedType.code }} · 最近更新 {{ formatTime(selectedType.updated_at) }}
          </p>
        </div>
        <EzActionButton
          v-if="canUse('system:dict:item:create')"
          kind="add"
          label="新增"
          size="small"
          type="primary"
          ghost
          :disabled="!selectedType"
          @click="$emit('create')"
        />
      </div>

      <div
        v-if="selectedType"
        class="grid min-h-0 min-w-0 grid-rows-[auto_minmax(0,1fr)_auto] overflow-hidden"
      >
        <div class="dict-item-filter-grid px-5 py-4 max-[720px]:px-4">
          <NInput
            v-model:value="query.keyword"
            clearable
            placeholder="编码 / 名称 / 值"
            @keyup.enter="$emit('search')"
          />
          <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" />
          <div class="dict-item-filter-actions ez-filter-actions">
            <EzActionButton kind="search" label="查询" type="primary" @click="$emit('search')" />
            <EzActionButton kind="reset" label="重置" @click="$emit('reset')" />
          </div>
        </div>

        <NDataTable
          remote
          class="dict-table h-full min-h-0"
          :columns="columns"
          :data="items"
          :loading="loading"
          :pagination="false"
          :row-key="(row: DictItem) => row.id"
          :bordered="false"
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

      <div v-else class="m-4 flex flex-1 items-center justify-center">
        <EmptyState
          title="先选择字典类型"
          description="从左侧选择一个字典类型后，再维护它的字典项。"
        />
      </div>
    </div>
  </NCard>
</template>

<style scoped>
.dict-item-panel-body {
  container: dict-item-panel / inline-size;
}

.dict-item-filter-grid {
  display: grid;
  width: 100%;
  max-width: 100%;
  min-width: 0;
  align-items: center;
  grid-template-columns: minmax(0, 1fr);
  gap: 12px;
}

.dict-item-filter-actions {
  justify-content: flex-start;
}

@container dict-item-panel (min-width: 560px) {
  .dict-item-filter-grid {
    grid-template-columns: minmax(220px, 1fr) 128px max-content;
  }

  .dict-item-filter-actions {
    grid-column: auto;
    justify-content: flex-end;
  }
}
</style>
