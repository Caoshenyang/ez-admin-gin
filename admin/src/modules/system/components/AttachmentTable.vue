<script setup lang="ts">
import { AttachOutline, CopyOutline } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NIcon,
  NPagination,
  NPopconfirm,
  NSpace,
  NTag,
  NTooltip,
} from 'naive-ui'
import { computed, h } from 'vue'

import { formatSize, formatTime } from '@/utils/format'
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
    minWidth: 260,
    render(row) {
      return h('div', { class: 'flex items-center gap-3' }, [
        h(
          'div',
          { class: 'flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-[#EEF2FF] text-[#4F46E5]' },
          [h(NIcon, { size: 18 }, { default: () => h(AttachOutline) })],
        ),
        h('div', { class: 'min-w-0 leading-5' }, [
          h('p', { class: 'truncate font-medium text-[#111827]' }, row.display_name),
          h('p', { class: 'truncate text-xs text-[#6B7280]' }, row.original_name),
        ]),
      ])
    },
  },
  {
    title: '分类 / 业务',
    key: 'category',
    width: 180,
    render(row) {
      return h('span', { class: 'text-sm text-[#334155]' }, [row.category || '未分类', row.biz_type || '通用'].join(' / '))
    },
  },
  {
    title: '类型',
    key: 'ext',
    width: 100,
    render(row) {
      return h(NTag, { size: 'small', bordered: false }, { default: () => row.ext || '-' })
    },
  },
  {
    title: '大小',
    key: 'size',
    width: 110,
    render(row) {
      return formatSize(row.size)
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
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
    width: 180,
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
      const nextStatus = row.status === AttachmentStatus.Enabled ? AttachmentStatus.Disabled : AttachmentStatus.Enabled

      return h(NSpace, { size: 8 }, {
        default: () => [
          h(
            NTooltip,
            {},
            {
              trigger: () =>
                h(
                  NButton,
                  { size: 'small', ghost: true, type: 'info', onClick: () => emit('copy', row.url) },
                  { icon: () => h(NIcon, null, { default: () => h(CopyOutline) }) },
                ),
              default: () => '复制链接',
            },
          ),
          props.canUse('system:attachment:update')
            ? h(
                NButton,
                { size: 'small', ghost: true, type: 'primary', onClick: () => emit('edit', row) },
                { default: () => '编辑' },
              )
            : null,
          props.canUse('system:attachment:status')
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
                        type: nextStatus === AttachmentStatus.Disabled ? 'error' : 'success',
                      },
                      { default: () => (nextStatus === AttachmentStatus.Disabled ? '禁用' : '启用') },
                    ),
                  default: () => `确认${nextStatus === AttachmentStatus.Disabled ? '禁用' : '启用'}该附件？`,
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
  <NCard :bordered="false" class="flex-1 overflow-hidden rounded-lg">
    <div class="flex h-full flex-col">
      <NDataTable :columns="columns" :data="attachments" :loading="loading" :bordered="false" flex-height class="flex-1" />

      <div
        v-if="!loading && !hasRows"
        class="flex flex-1 items-center justify-center rounded-2xl border border-dashed border-[#D9DEE8] bg-[#FAFBFC]"
      >
        <NEmpty description="当前没有附件记录，先上传一份业务附件吧。" />
      </div>

      <div class="mt-4 flex justify-end">
        <NPagination
          :page="page"
          :page-size="pageSize"
          :item-count="total"
          show-size-picker
          :page-sizes="[10, 20, 50]"
          @update:page="(page) => emit('pageChange', page)"
          @update:page-size="(pageSize) => emit('pageSizeChange', pageSize)"
        />
      </div>
    </div>
  </NCard>
</template>
