<script setup lang="ts">
import { AttachOutline } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import { NDataTable, NIcon, NPopconfirm, NSpace, NTag } from 'naive-ui'
import { computed, h } from 'vue'

import EmptyState from '@/components/EmptyState.vue'
import EzActionButton from '@/components/ez/EzActionButton.vue'
import EzDataTable from '@/components/ez/EzDataTable.vue'
import { displayText, formatSize, formatTime } from '@/utils/format'
import { AttachmentStatus, type AttachmentItem } from '../types/attachment'

const props = defineProps<{
  attachments: AttachmentItem[]
  canUse: (code: string) => boolean
  hasRows: boolean
  loading: boolean
  page: number
  pageSize: number
  total: number
}>()

const emit = defineEmits<{
  copy: [url: string]
  edit: [row: AttachmentItem]
  pageChange: [page: number]
  pageSizeChange: [pageSize: number]
  refresh: []
  toggleStatus: [row: AttachmentItem, status: AttachmentStatus]
}>()

const columns = computed<DataTableColumns<AttachmentItem>>(() => [
  {
    title: '附件',
    key: 'display_name',
    minWidth: 220,
    render(row) {
      return h('div', { class: 'flex items-center gap-3' }, [
        h(
          'div',
          {
            class:
              'flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-[var(--ez-radius-xl)] bg-[var(--ez-accent-indigo-soft)] text-[var(--ez-accent-indigo)]',
          },
          [h(NIcon, { size: 18 }, { default: () => h(AttachOutline) })],
        ),
        h('div', { class: 'min-w-0 leading-5' }, [
          h(
            'p',
            { class: 'truncate font-medium text-[var(--ez-text-main)]' },
            displayText(row.display_name),
          ),
          h(
            'p',
            { class: 'truncate text-xs text-[var(--ez-text-sub)]' },
            displayText(row.original_name),
          ),
        ]),
      ])
    },
  },
  {
    title: '分类 / 业务',
    key: 'category',
    width: 150,
    render(row) {
      return h(
        'span',
        { class: 'text-sm text-[var(--ez-text-main)]' },
        [displayText(row.category, '未分类'), displayText(row.biz_type, '通用')].join(' / '),
      )
    },
  },
  {
    title: '类型',
    key: 'ext',
    width: 84,
    render(row) {
      return h(NTag, { size: 'small', bordered: false }, { default: () => row.ext || '-' })
    },
  },
  {
    title: '大小',
    key: 'size',
    width: 92,
    render(row) {
      return formatSize(row.size)
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 78,
    render(row) {
      return h(
        NTag,
        { bordered: false, type: row.status === AttachmentStatus.Enabled ? 'success' : 'error' },
        { default: () => (row.status === AttachmentStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '上传时间',
    key: 'created_at',
    width: 150,
    render(row) {
      return formatTime(row.created_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 128,
    fixed: 'right',
    render(row) {
      const nextStatus =
        row.status === AttachmentStatus.Enabled
          ? AttachmentStatus.Disabled
          : AttachmentStatus.Enabled

      return h(
        NSpace,
        { class: 'ez-row-actions', size: 6, align: 'center' },
        {
          default: () =>
            [
              h(EzActionButton, {
                iconOnly: true,
                kind: 'copy',
                label: '复制链接',
                size: 'tiny',
                secondary: true,
                type: 'info',
                onClick: () => emit('copy', row.url),
              }),
              props.canUse('system:attachment:update')
                ? h(EzActionButton, {
                    iconOnly: true,
                    kind: 'edit',
                    label: '编辑',
                    size: 'tiny',
                    secondary: true,
                    type: 'primary',
                    onClick: () => emit('edit', row),
                  })
                : null,
              props.canUse('system:attachment:status')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => emit('toggleStatus', row, nextStatus) },
                    {
                      trigger: () =>
                        h(EzActionButton, {
                          iconOnly: true,
                          kind: nextStatus === AttachmentStatus.Disabled ? 'disable' : 'enable',
                          label: nextStatus === AttachmentStatus.Disabled ? '禁用' : '启用',
                          size: 'tiny',
                          secondary: true,
                          tooltip: false,
                          type: nextStatus === AttachmentStatus.Disabled ? 'error' : 'success',
                        }),
                      default: () =>
                        `确认${nextStatus === AttachmentStatus.Disabled ? '禁用' : '启用'}该附件？`,
                    },
                  )
                : null,
            ].filter(Boolean),
        },
      )
    },
  },
])
</script>

<template>
  <EzDataTable
    :columns="columns"
    :data="attachments"
    :loading="loading"
    :page="page"
    :page-size="pageSize"
    :summary-text="`共 ${total} 个附件`"
    :total="total"
    @page-change="(page) => emit('pageChange', page)"
    @page-size-change="(pageSize) => emit('pageSizeChange', pageSize)"
    @refresh="emit('refresh')"
  >
    <template #body="{ tableColumns, tableScrollX, tableSize }">
      <NDataTable
        :columns="tableColumns"
        :data="attachments"
        :loading="loading"
        :bordered="false"
        :scroll-x="tableScrollX"
        :size="tableSize"
        flex-height
        class="ez-table-fill-table"
      />

      <div v-if="!loading && !hasRows" class="flex flex-1 items-center justify-center p-4">
        <EmptyState title="当前没有附件记录" description="先上传一份业务附件吧。" kind="create" />
      </div>
    </template>
  </EzDataTable>
</template>
