<script setup lang="ts">
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import { NDataTable, NPopconfirm, NSpace } from 'naive-ui'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import EzDataTable from '@/components/ez/EzDataTable.vue'
import type { PostItem } from '../types/post'

defineProps<{
  canUse: (code: string) => boolean
  checkedRowKeys: DataTableRowKey[]
  columns: DataTableColumns<PostItem>
  items: PostItem[]
  loading: boolean
  selectedCount: number
}>()

defineEmits<{
  checkedRowKeysChange: [keys: DataTableRowKey[]]
  deleteSelected: []
  refresh: []
}>()
</script>

<template>
  <EzDataTable :columns="columns" :data="items" :loading="loading" @refresh="$emit('refresh')">
    <template #toolbarSummary>
      <span>已选 {{ selectedCount }} 项</span>
    </template>

    <template #toolbarActions>
      <NSpace :size="8">
        <NPopconfirm
          v-if="canUse('system:post:delete')"
          :disabled="selectedCount === 0"
          @positive-click="$emit('deleteSelected')"
        >
          <template #trigger>
            <EzActionButton
              kind="delete"
              label="删除选中"
              quaternary
              size="small"
              type="error"
              :disabled="selectedCount === 0"
            />
          </template>
          删除前请确认选中的岗位没有绑定到任何用户。
        </NPopconfirm>
      </NSpace>
    </template>

    <template #body="{ tableColumns, tableScrollX, tableSize }">
      <NDataTable
        class="ez-table-fill-table"
        :columns="tableColumns"
        :data="items"
        :loading="loading"
        :pagination="false"
        :row-key="(row: PostItem) => row.id"
        :checked-row-keys="checkedRowKeys"
        :scroll-x="tableScrollX"
        :size="tableSize"
        :bordered="false"
        flex-height
        @update:checked-row-keys="$emit('checkedRowKeysChange', $event)"
      />
    </template>
  </EzDataTable>
</template>
