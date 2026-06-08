<script setup lang="ts">
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'
import {
  NButton,
  NDataTable,
  NPopconfirm,
  NSpace,
  NTag,
} from 'naive-ui'
import { computed, h } from 'vue'

import EzDataTable from '@/components/ez/EzDataTable.vue'
import { displayText, formatTime } from '@/utils/format'
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
  delete: [row: UserItem]
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
    minWidth: 150,
    render(row) {
      return h('div', { class: 'leading-6' }, [
        h('p', { class: 'font-semibold text-[var(--ez-text-main)]' }, displayText(row.username)),
        h('p', { class: 'text-xs text-[var(--ez-text-sub)]' }, displayText(row.nickname)),
      ])
    },
  },
  {
    title: '部门',
    key: 'department_id',
    minWidth: 120,
    render(row) {
      if (row.department_id === 0) {
        return h('span', { class: 'text-sm text-[var(--ez-text-light)]' }, '未分配')
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
    minWidth: 150,
    render(row) {
      if (row.role_ids.length === 0) {
        return h('span', { class: 'text-sm text-[var(--ez-text-light)]' }, '未分配')
      }

      return h('div', { class: 'flex flex-wrap gap-1' }, row.role_ids.map((roleID) =>
        h(NTag, { size: 'small', bordered: false }, { default: () => props.roleNameMap.get(roleID) ?? `角色 ${roleID}` }),
      ))
    },
  },
  {
    title: '岗位',
    key: 'post_ids',
    minWidth: 150,
    render(row) {
      if (row.post_ids.length === 0) {
        return h('span', { class: 'text-sm text-[var(--ez-text-light)]' }, '未绑定')
      }

      return h('div', { class: 'flex flex-wrap gap-1' }, row.post_ids.map((postID) =>
        h(
          NTag,
          { size: 'small', bordered: false, type: 'info' },
          { default: () => props.postNameMap.get(postID) ?? `岗位 ${postID}` },
        ),
      ))
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 86,
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
    width: 150,
    render(row) {
      return formatTime(row.created_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    minWidth: 190,
    render(row) {
      const nextStatus = row.status === UserStatus.Enabled ? UserStatus.Disabled : UserStatus.Enabled

      return h(NSpace, { size: 8, align: 'center', wrap: true }, {
        default: () => [
          props.canUse('system:user:update')
            ? h(
                NButton,
                {
                  size: 'tiny',
                  secondary: true,
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
                        size: 'tiny',
                        secondary: true,
                        type: nextStatus === UserStatus.Disabled ? 'error' : 'success',
                        class: 'min-w-[48px]',
                      },
                      { default: () => (nextStatus === UserStatus.Disabled ? '禁用' : '启用') },
                    ),
                  default: () => `确认${nextStatus === UserStatus.Disabled ? '禁用' : '启用'}该用户？`,
                },
              )
            : null,
          props.canUse('system:user:delete')
            ? h(
                NPopconfirm,
                { onPositiveClick: () => emit('delete', row) },
                {
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: 'tiny',
                        secondary: true,
                        type: 'error',
                        class: 'min-w-[48px]',
                      },
                      { default: () => '删除' },
                    ),
                  default: () => '删除后用户账号和角色/岗位关联会一起移除，确认继续？',
                },
              )
            : null,
          props.canUse('system:user:assign-role')
            ? h(
                NButton,
                {
                  size: 'tiny',
                  secondary: true,
                  type: 'primary',
                  onClick: () => emit('role', row),
                },
                { default: () => '分配角色' },
              )
            : null,
        ].filter(Boolean),
      })
    },
  },
])
</script>

<template>
  <EzDataTable
    remote
    :columns="columns"
    :data="users"
    :loading="loading"
    :page="page"
    :page-size="pageSize"
    :total="displayTotal"
    @page-change="(page) => emit('pageChange', page)"
    @page-size-change="(pageSize) => emit('pageSizeChange', pageSize)"
    @refresh="emit('refresh')"
  >
    <template #toolbarSummary>
      <span>已选 {{ selectedCount }} 项</span>
    </template>

    <template #body="{ tableColumns, tableSize }">
      <NDataTable
        remote
        class="ez-table-fill-table"
        :columns="tableColumns"
        :data="users"
        :loading="loading"
        :pagination="false"
        :row-key="(row: UserItem) => row.id"
        :checked-row-keys="checkedRowKeys"
        :size="tableSize"
        :bordered="false"
        flex-height
        @update:checked-row-keys="(keys) => emit('checkedRowKeysChange', keys)"
      />
    </template>

    <template #footerSummary>
      <span>共 {{ displayTotal }} 条，已选择 {{ selectedCount }} 条</span>
    </template>
  </EzDataTable>
</template>
