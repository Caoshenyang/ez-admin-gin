<script setup lang="ts" generic="T extends object">
import {
  EyeOffOutline,
  OptionsOutline,
  RefreshOutline,
  ResizeOutline,
} from '@vicons/ionicons5'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import {
  NButton,
  NCheckbox,
  NCheckboxGroup,
  NDataTable,
  NIcon,
  NPagination,
  NPopover,
  NRadioButton,
  NRadioGroup,
  NTooltip,
} from 'naive-ui'
import { computed, ref, watch } from 'vue'

import EzTableCard from './EzTableCard.vue'
import TableStatsBar from '../TableStatsBar.vue'

type TableSize = 'small' | 'medium' | 'large'
type TableColumn = DataTableColumns<T>[number]

interface ColumnSetting {
  key: string
  label: string
}

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
    showRefresh?: boolean
    showColumnSettings?: boolean
    showDensity?: boolean
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
    showRefresh: true,
    showColumnSettings: true,
    showDensity: true,
  },
)

const emit = defineEmits<{
  pageChange: [page: number]
  pageSizeChange: [pageSize: number]
  refresh: []
}>()

defineSlots<{
  actions?: () => unknown
  body?: (props: {
    tableColumns: DataTableColumns<T>
    tableScrollX: number
    tableSize: TableSize
  }) => unknown
  footerSummary?: () => unknown
  toolbarActions?: () => unknown
  toolbarSummary?: () => unknown
}>()

const tableSize = ref<TableSize>('medium')
const visibleColumnKeys = ref<string[]>([])

const hasPagination = computed(() => {
  return props.page !== undefined && props.pageSize !== undefined && props.total !== undefined
})

const displaySummary = computed(() => {
  return props.summaryText || `共 ${props.total ?? props.data.length} 条`
})

const columnSettings = computed<ColumnSetting[]>(() => {
  return props.columns.flatMap((column) => {
    if (!('key' in column) || typeof column.key !== 'string') {
      return []
    }

    if (column.type === 'selection' || column.key === 'actions') {
      return []
    }

    const title = typeof column.title === 'string' ? column.title : column.key

    return [{ key: column.key, label: title }]
  })
})

const configurableColumnKeys = computed(() => columnSettings.value.map((column) => column.key))

const tableColumns = computed<DataTableColumns<T>>(() => {
  if (!props.showColumnSettings) {
    return props.columns
  }

  const visibleKeys = new Set(visibleColumnKeys.value)

  return props.columns.filter((column) => {
    if (!('key' in column) || typeof column.key !== 'string') {
      return true
    }

    if (!configurableColumnKeys.value.includes(column.key)) {
      return true
    }

    return visibleKeys.has(column.key)
  })
})

const tableScrollX = computed(() => {
  return tableColumns.value.reduce((total, column) => total + resolveColumnWidth(column), 0)
})

const hasHiddenColumns = computed(() => {
  return visibleColumnKeys.value.length < configurableColumnKeys.value.length
})

const hasViewControls = computed(() => {
  return props.showRefresh || props.showDensity || (props.showColumnSettings && columnSettings.value.length > 0)
})

watch(
  configurableColumnKeys,
  (keys, oldKeys = []) => {
    if (oldKeys.length === 0) {
      visibleColumnKeys.value = [...keys]
      return
    }

    const current = new Set(visibleColumnKeys.value)
    const oldKeySet = new Set(oldKeys)
    const mergedKeys = keys.filter((key) => current.has(key) || !oldKeySet.has(key))

    visibleColumnKeys.value = mergedKeys.length > 0 ? mergedKeys : [...keys]
  },
  { immediate: true },
)

function resetColumns() {
  visibleColumnKeys.value = [...configurableColumnKeys.value]
}

function resolveColumnWidth(column: TableColumn): number {
  if (column.type === 'selection') {
    return toWidthNumber(column.width) ?? 48
  }

  if ('children' in column && Array.isArray(column.children)) {
    return column.children.reduce((total, childColumn) => total + resolveColumnWidth(childColumn), 0)
  }

  return toWidthNumber(column.width) ?? toWidthNumber(column.minWidth) ?? 120
}

function toWidthNumber(value: unknown): number | undefined {
  if (typeof value === 'number') {
    return value
  }

  if (typeof value === 'string') {
    const matched = value.match(/^(\d+(?:\.\d+)?)px$/)
    return matched ? Number(matched[1]) : undefined
  }

  return undefined
}
</script>

<template>
  <EzTableCard>
    <TableStatsBar>
      <slot name="toolbarSummary">
        <span>{{ displaySummary }}</span>
      </slot>
      <template #actions>
        <slot name="actions">
          <div class="ez-table-actions">
            <div v-if="$slots.toolbarActions" class="ez-table-action-group ez-table-action-group--business">
              <slot name="toolbarActions" />
            </div>

            <div v-if="hasViewControls" class="ez-table-action-group ez-table-action-group--view">
              <NTooltip v-if="showRefresh" trigger="hover">
                <template #trigger>
                  <NButton
                    quaternary
                    size="small"
                    class="ez-table-view-button"
                    :disabled="loading"
                    :aria-label="loading ? '正在刷新' : '刷新列表'"
                    title="刷新"
                    @click="emit('refresh')"
                  >
                    <template #icon>
                      <NIcon :component="RefreshOutline" />
                    </template>
                  </NButton>
                </template>
                刷新
              </NTooltip>
              <span v-if="showRefresh && showDensity" class="ez-table-view-divider" aria-hidden="true" />
              <NPopover v-if="showDensity" trigger="click" placement="bottom-end">
                <template #trigger>
                  <NButton
                    quaternary
                    size="small"
                    class="ez-table-view-button"
                    aria-label="调整表格密度"
                    title="密度"
                  >
                    <template #icon>
                      <NIcon :component="ResizeOutline" />
                    </template>
                  </NButton>
                </template>
                <div class="ez-table-popover">
                  <div class="ez-table-popover__title">表格大小</div>
                  <NRadioGroup v-model:value="tableSize" size="small">
                    <NRadioButton value="small">紧凑</NRadioButton>
                    <NRadioButton value="medium">默认</NRadioButton>
                    <NRadioButton value="large">宽松</NRadioButton>
                  </NRadioGroup>
                </div>
              </NPopover>
              <span
                v-if="showColumnSettings && columnSettings.length > 0 && (showRefresh || showDensity)"
                class="ez-table-view-divider"
                aria-hidden="true"
              />
              <NPopover
                v-if="showColumnSettings && columnSettings.length > 0"
                trigger="click"
                placement="bottom-end"
              >
                <template #trigger>
                  <NButton
                    quaternary
                    size="small"
                    :class="['ez-table-view-button', { 'ez-table-view-button--active': hasHiddenColumns }]"
                    :type="hasHiddenColumns ? 'primary' : 'default'"
                    aria-label="设置显示列"
                    title="列设置"
                  >
                    <template #icon>
                      <NIcon :component="OptionsOutline" />
                    </template>
                  </NButton>
                </template>
                <div class="ez-table-popover ez-table-column-panel">
                  <div class="ez-table-popover__header">
                    <span class="ez-table-popover__title">显示列</span>
                    <NButton text size="tiny" type="primary" @click="resetColumns">重置</NButton>
                  </div>
                  <NCheckboxGroup v-model:value="visibleColumnKeys">
                    <div class="ez-table-column-panel__checks">
                      <NCheckbox
                        v-for="column in columnSettings"
                        :key="column.key"
                        :value="column.key"
                        size="small"
                      >
                        {{ column.label }}
                      </NCheckbox>
                    </div>
                  </NCheckboxGroup>
                  <div v-if="hasHiddenColumns" class="ez-table-popover__hint">
                    <NIcon :component="EyeOffOutline" :size="14" />
                    已隐藏 {{ configurableColumnKeys.length - visibleColumnKeys.length }} 列
                  </div>
                </div>
              </NPopover>
            </div>
          </div>
        </slot>
      </template>
    </TableStatsBar>

    <div class="ez-table-fill">
      <!--
        普通列表只传 columns/data/rowKey 即可；需要勾选、行样式或额外事件时，
        用 body 插槽接管 NDataTable，仍复用统一工具条和分页外壳。
      -->
      <slot name="body" :table-columns="tableColumns" :table-scroll-x="tableScrollX" :table-size="tableSize">
        <NDataTable
          :remote="remote"
          class="ez-table-fill-table"
          :columns="tableColumns"
          :data="data"
          :loading="loading"
          :pagination="false"
          :scroll-x="tableScrollX"
          :row-key="rowKey"
          :row-props="rowProps"
          :size="tableSize"
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
