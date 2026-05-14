<!-- DictItemPanel 展示字典类型下的字典项列表，支持增删改。 -->
<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NEmpty, NInput, NPagination, NSelect, NTag } from 'naive-ui'

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
  <NCard class="ez-card min-h-0 rounded-lg" :bordered="false" content-class="ez-card-content-fill">
    <div class="grid h-full min-h-0 grid-rows-[auto_minmax(0,1fr)]">
      <div class="flex items-start justify-between gap-4 border-b border-[var(--ez-border-light)] px-5 pt-[18px] pb-[14px] max-[720px]:px-4">
        <div>
          <p class="text-[11px] font-bold tracking-[0.14em] text-slate-400 uppercase">Items</p>
          <div class="flex items-center gap-2">
            <h2 class="mt-1.5 text-[18px] font-bold text-[var(--ez-text-main)]">字典项</h2>
            <NTag v-if="selectedType" size="small" type="info" :bordered="false">
              {{ selectedType.name }}
            </NTag>
          </div>
          <p v-if="selectedType" class="mt-1 text-xs text-[var(--ez-text-sub)]">
            {{ selectedType.code }} · 最近更新 {{ formatTime(selectedType.updated_at) }}
          </p>
        </div>
        <NButton
          v-if="canUse('system:dict:item:create')"
          size="small"
          type="primary"
          ghost
          :disabled="!selectedType"
          @click="$emit('create')"
        >
          新增
        </NButton>
      </div>

      <template v-if="selectedType">
        <div class="grid gap-3 px-5 py-4 min-[1281px]:grid-cols-[minmax(0,1fr)_140px_auto] max-[720px]:px-4">
          <NInput v-model:value="query.keyword" clearable placeholder="编码 / 名称 / 值" @keyup.enter="$emit('search')" />
          <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" />
          <div class="flex gap-2.5 min-[1281px]:justify-start max-[1280px]:justify-end">
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
          :row-key="(row: DictItem) => row.id"
          :bordered="false"
          flex-height
        />

        <div class="flex items-center justify-between gap-4 border-t border-[var(--ez-border-light)] px-5 py-3.5 text-[13px] text-[var(--ez-text-sub)] max-[720px]:flex-col max-[720px]:items-stretch max-[720px]:px-4">
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
      </template>

      <div
        v-else
        class="m-4 flex flex-1 items-center justify-center rounded-2xl border border-dashed border-[var(--ez-border)] bg-[var(--ez-page-bg)]"
      >
        <NEmpty description="先从左侧选择一个字典类型，再维护它的字典项。" />
      </div>
    </div>
  </NCard>
</template>

<style scoped>
</style>
