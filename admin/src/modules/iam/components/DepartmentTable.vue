<script setup lang="ts">
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import { NButton, NDataTable, NPopconfirm, NSpace, NTag } from 'naive-ui'
import { computed, h } from 'vue'

import EzDataTable from '@/components/ez/EzDataTable.vue'
import { displayText } from '@/utils/format'
import { DepartmentStatus, type DepartmentItem } from '@/modules/iam/types/department'

const props = defineProps<{
  canUse: (code: string) => boolean
  checkedRowKeys: DataTableRowKey[]
  departments: DepartmentItem[]
  expandedRowKeys: DataTableRowKey[]
  leaderNameMap: Map<number, string>
  loading: boolean
  selectedCount: number
}>()

const emit = defineEmits<{
  checkedRowKeysChange: [keys: DataTableRowKey[]]
  collapseAll: []
  createChild: [row: DepartmentItem]
  deleteSelected: []
  edit: [row: DepartmentItem]
  expandedRowKeysChange: [keys: DataTableRowKey[]]
  expandAll: []
  refresh: []
  toggleStatus: [row: DepartmentItem, status: DepartmentStatus]
}>()

function formatTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

function countDepartments(items: DepartmentItem[]): number {
  return items.reduce((total, item) => total + 1 + countDepartments(item.children ?? []), 0)
}

const departmentCount = computed(() => countDepartments(props.departments))

const columns = computed<DataTableColumns<DepartmentItem>>(() => [
  { type: 'selection', width: 44 },
  {
    title: '部门名称',
    key: 'name',
    minWidth: 190,
    render(row) {
      return h('span', { class: 'font-medium text-[var(--ez-text-main)]' }, displayText(row.name))
    },
  },
  {
    title: '部门编码',
    key: 'code',
    minWidth: 136,
    ellipsis: { tooltip: true },
    render(row) {
      return displayText(row.code)
    },
  },
  {
    title: '负责人',
    key: 'leader_user_id',
    width: 120,
    render(row) {
      if (!row.leader_user_id) {
        return h('span', { class: 'text-[var(--ez-text-light)]' }, '未设置')
      }

      return props.leaderNameMap.get(row.leader_user_id) ?? `用户 ${row.leader_user_id}`
    },
  },
  {
    title: '排序',
    key: 'sort',
    width: 72,
    align: 'center',
  },
  {
    title: '状态',
    key: 'status',
    width: 84,
    align: 'center',
    render(row) {
      return h(
        NTag,
        {
          bordered: false,
          type: row.status === DepartmentStatus.Enabled ? 'success' : 'error',
        },
        { default: () => (row.status === DepartmentStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 150,
    render(row) {
      return formatTime(row.updated_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 196,
    fixed: 'right',
    render(row) {
      const nextStatus =
        row.status === DepartmentStatus.Enabled
          ? DepartmentStatus.Disabled
          : DepartmentStatus.Enabled

      return h(
        NSpace,
        { size: 8 },
        {
          default: () => [
            props.canUse('system:department:create')
              ? h(
                  NButton,
                  {
                    size: 'tiny',
                    secondary: true,
                    type: 'primary',
                    onClick: () => emit('createChild', row),
                  },
                  { default: () => '新增子部门' },
                )
              : null,
            props.canUse('system:department:update')
              ? h(
                  NButton,
                  {
                    size: 'tiny',
                    secondary: true,
                    type: 'info',
                    onClick: () => emit('edit', row),
                  },
                  { default: () => '编辑' },
                )
              : null,
            props.canUse('system:department:status')
              ? h(
                  NPopconfirm,
                  { onPositiveClick: () => emit('toggleStatus', row, nextStatus) },
                  {
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: 'tiny',
                          secondary: true,
                          type: nextStatus === DepartmentStatus.Disabled ? 'error' : 'success',
                        },
                        { default: () => (nextStatus === DepartmentStatus.Disabled ? '禁用' : '启用') },
                      ),
                    default: () =>
                      `确认${nextStatus === DepartmentStatus.Disabled ? '禁用' : '启用'}该部门？`,
                  },
                )
              : null,
          ],
        },
      )
    },
  },
])
</script>

<template>
  <EzDataTable :columns="columns" :data="departments" :loading="loading" @refresh="emit('refresh')">
    <template #toolbarSummary>
      <span>共 {{ departmentCount }} 个部门节点，已选 {{ selectedCount }} 项</span>
    </template>

    <template #toolbarActions>
      <NSpace :size="12">
        <NButton quaternary size="small" @click="emit('expandAll')">展开全部</NButton>
        <NButton quaternary size="small" @click="emit('collapseAll')">收起全部</NButton>
        <NPopconfirm
          v-if="canUse('system:department:delete')"
          :disabled="selectedCount === 0"
          @positive-click="emit('deleteSelected')"
        >
          <template #trigger>
            <NButton quaternary size="small" type="error" :disabled="selectedCount === 0">删除选中</NButton>
          </template>
          删除前请确认选中的部门没有关联用户、角色数据范围，且子部门也已一并选中。
        </NPopconfirm>
      </NSpace>
    </template>

    <template #body="{ tableColumns, tableSize }">
      <NDataTable
        class="ez-table-fill-table department-table"
        :columns="tableColumns"
        :data="departments"
        :loading="loading"
        :row-key="(row: DepartmentItem) => row.id"
        :checked-row-keys="checkedRowKeys"
        :expanded-row-keys="expandedRowKeys"
        :pagination="false"
        :size="tableSize"
        :bordered="false"
        children-key="children"
        flex-height
        @update:checked-row-keys="(keys) => emit('checkedRowKeysChange', keys)"
        @update:expanded-row-keys="(keys) => emit('expandedRowKeysChange', keys)"
      />
    </template>
  </EzDataTable>
</template>
