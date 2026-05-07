<script setup lang="ts">
import { EllipsisHorizontal } from '@vicons/ionicons5'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NDropdown,
  NIcon,
  NPagination,
  NPopconfirm,
  NSpace,
  NTag,
} from 'naive-ui'
import { computed, h } from 'vue'

import { formatTime } from '@/utils/format'
import { UserStatus, type UserItem } from '../types/user'

const props = defineProps<{
  checkedRowKeys: DataTableRowKey[]
  departmentNameMap: Map<number, string>
  displayTotal: number
  loading: boolean
  page: number
  pageSize: number
  postNameMap: Map<number, string>
  roleNameMap: Map<number, string>
  selectedCount: number
  users: UserItem[]
  canUse: (code: string) => boolean
}>()

const emit = defineEmits<{
  checkedRowKeysChange: [keys: DataTableRowKey[]]
  edit: [row: UserItem]
  pageChange: [page: number]
  pageSizeChange: [pageSize: number]
  refresh: []
  role: [row: UserItem]
  toggleStatus: [row: UserItem, status: UserStatus]
}>()

const columns = computed<DataTableColumns<UserItem>>(() => [
  { type: 'selection', width: 48 },
  {
    title: '用户',
    key: 'username',
    minWidth: 180,
    render(row) {
      return h('div', { class: 'leading-6' }, [
        h('p', { class: 'font-semibold text-[#111827]' }, row.username),
        h('p', { class: 'text-xs text-[#6B7280]' }, row.nickname),
      ])
    },
  },
  {
    title: '部门',
    key: 'department_id',
    minWidth: 180,
    render(row) {
      if (row.department_id === 0) {
        return h('span', { class: 'text-sm text-[#9CA3AF]' }, '未分配')
      }

      return h(
        NTag,
        { size: 'small', bordered: false, type: 'warning' },
        { default: () => props.departmentNameMap.get(row.department_id) ?? `部门 ${row.department_id}` },
      )
    },
  },
  {
    title: '角色',
    key: 'role_ids',
    minWidth: 220,
    render(row) {
      if (row.role_ids.length === 0) {
        return h('span', { class: 'text-sm text-[#9CA3AF]' }, '未分配')
      }

      return h(NSpace, { size: 6 }, {
        default: () => row.role_ids.map((roleID) =>
          h(NTag, { size: 'small', bordered: false }, { default: () => props.roleNameMap.get(roleID) ?? `角色 ${roleID}` }),
        ),
      })
    },
  },
  {
    title: '岗位',
    key: 'post_ids',
    minWidth: 220,
    render(row) {
      if (row.post_ids.length === 0) {
        return h('span', { class: 'text-sm text-[#9CA3AF]' }, '未绑定')
      }

      return h(NSpace, { size: 6 }, {
        default: () => row.post_ids.map((postID) =>
          h(
            NTag,
            { size: 'small', bordered: false, type: 'info' },
            { default: () => props.postNameMap.get(postID) ?? `岗位 ${postID}` },
          ),
        ),
      })
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 110,
    render(row) {
      return h(
        NTag,
        { type: row.status === UserStatus.Enabled ? 'success' : 'error', bordered: false },
        { default: () => (row.status === UserStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 190,
    render(row) {
      return formatTime(row.created_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    fixed: 'right',
    render(row) {
      const nextStatus = row.status === UserStatus.Enabled ? UserStatus.Disabled : UserStatus.Enabled
      const dropdownOptions = props.canUse('system:user:assign-role')
        ? [{ label: '分配角色', key: `role:${row.id}` }]
        : []

      return h(NSpace, { size: 8, align: 'center' }, {
        default: () => [
          props.canUse('system:user:update')
            ? h(
                NButton,
                {
                  size: 'small',
                  ghost: true,
                  type: 'info',
                  class: 'min-w-[48px]',
                  onClick: () => emit('edit', row),
                },
                { default: () => '编辑' },
              )
            : null,
          props.canUse('system:user:status')
            ? h(
                NPopconfirm,
                { onPositiveClick: () => emit('toggleStatus', row, nextStatus) },
                {
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: 'small',
                        ghost: true,
                        type: nextStatus === UserStatus.Disabled ? 'error' : 'success',
                        class: 'min-w-[48px]',
                      },
                      { default: () => (nextStatus === UserStatus.Disabled ? '禁用' : '启用') },
                    ),
                  default: () => `确认${nextStatus === UserStatus.Disabled ? '禁用' : '启用'}该用户？`,
                },
              )
            : null,
          dropdownOptions.length > 0
            ? h(
                NDropdown,
                {
                  options: dropdownOptions,
                  trigger: 'click',
                  onSelect: () => emit('role', row),
                },
                {
                  default: () =>
                    h(
                      NButton,
                      { size: 'small', quaternary: true, class: 'min-w-[36px] px-2' },
                      { icon: () => h(NIcon, null, { default: () => h(EllipsisHorizontal) }) },
                    ),
                },
              )
            : null,
        ].filter(Boolean),
      })
    },
  },
])
</script>

<template>
  <NCard class="min-h-0 flex-1 rounded-lg" :bordered="false" content-style="height: 100%; padding: 0;">
    <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-3">
      <NSpace :size="12">
        <span class="text-sm text-[#6B7280]">已选 {{ selectedCount }} 项</span>
        <NButton text :disabled="selectedCount === 0">批量禁用</NButton>
        <NButton text :disabled="selectedCount === 0">批量删除</NButton>
      </NSpace>
      <NSpace :size="14">
        <NButton text type="primary">列设置</NButton>
        <NButton text type="primary">密度</NButton>
        <NButton text type="primary" @click="emit('refresh')">刷新</NButton>
      </NSpace>
    </div>

    <NDataTable
      remote
      class="user-table h-full"
      style="height: calc(100% - 105px)"
      :columns="columns"
      :data="users"
      :loading="loading"
      :pagination="false"
      :row-key="(row: UserItem) => row.id"
      :checked-row-keys="checkedRowKeys"
      :bordered="false"
      flex-height
      @update:checked-row-keys="(keys) => emit('checkedRowKeysChange', keys)"
    />

    <div class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]">
      <span>共 {{ displayTotal }} 条，已选择 {{ selectedCount }} 条</span>
      <NPagination
        :page="page"
        :page-size="pageSize"
        :item-count="displayTotal"
        :page-sizes="[10, 20, 50]"
        show-size-picker
        @update:page="(page) => emit('pageChange', page)"
        @update:page-size="(pageSize) => emit('pageSizeChange', pageSize)"
      />
    </div>
  </NCard>
</template>

<style scoped>
.user-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #374151;
  background: #fff;
}

.user-table :deep(.n-data-table-td) {
  color: #374151;
}

.user-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}
</style>
