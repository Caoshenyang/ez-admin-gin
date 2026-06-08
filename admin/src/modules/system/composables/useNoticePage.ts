import type { DataTableColumns, FormRules } from 'naive-ui'
import { NButton, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { h } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { displayText, formatTime } from '@/utils/format'
import { createNotice, deleteNotice, getNotices, updateNotice, updateNoticeStatus } from '../api/notice'
import { NoticeStatus, type NoticeItem, type NoticeListQuery } from '../types/notice'
import {
  buildNoticePayload,
  defaultNoticeFormModel,
  defaultNoticeListQuery,
  toNoticeFormModel,
  type NoticeFormModel,
} from './notice-page.utils'

export function useNoticePage() {
  const message = useMessage()
  const { canUse } = usePermission()

  const {
    items: notices,
    total,
    loading,
    query,
    load,
    handleSearch,
    handleReset,
    handlePageChange,
    handlePageSizeChange,
  } = useRemotePagination<NoticeItem, NoticeListQuery>(getNotices, {
    ...defaultNoticeListQuery(),
  })

  const { formRef, formVisible, formMode, formModel, saving, rules, openCreate, openEdit, handleSubmit } =
    useModalForm<NoticeFormModel>(defaultNoticeFormModel, {
      rules: { title: [{ required: true, message: '请输入公告标题', trigger: 'blur' }] } as FormRules,
    })

  const openEditForm = (row: NoticeItem) => {
    openEdit(toNoticeFormModel(row))
  }

  const { handleToggleStatus } = useStatusToggle(updateNoticeStatus, {
    onSuccess: async () => {
      await load()
    },
  })

  const columns: DataTableColumns<NoticeItem> = [
    {
      title: '标题',
      key: 'title',
      width: 180,
      ellipsis: { tooltip: true },
      render(row) {
        return h('span', { class: 'font-semibold text-[var(--ez-text-heading)]' }, displayText(row.title))
      },
    },
    {
      title: '内容',
      key: 'content',
      minWidth: 200,
      ellipsis: { tooltip: true },
      render(row) {
        return h('span', { class: 'text-[var(--ez-text-body)]' }, displayText(row.content))
      },
    },
    {
      title: '排序',
      key: 'sort',
      width: 66,
      align: 'center',
    },
    {
      title: '状态',
      key: 'status',
      width: 78,
      align: 'center',
      render(row) {
        return h(
          NTag,
          { type: row.status === NoticeStatus.Enabled ? 'success' : 'error', bordered: false },
          { default: () => (row.status === NoticeStatus.Enabled ? '启用' : '禁用') },
        )
      },
    },
    {
      title: '更新时间',
      key: 'updated_at',
      width: 150,
      render(row) {
        return h('span', { class: 'tabular-nums text-[var(--ez-text-muted)]' }, formatTime(row.updated_at))
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 172,
      fixed: 'right',
      render(row) {
        const nextStatus = row.status === NoticeStatus.Enabled ? NoticeStatus.Disabled : NoticeStatus.Enabled

        return h(
          NSpace,
          { class: 'ez-row-actions', size: 6, align: 'center' },
          {
            default: () =>
              [
                canUse('system:notice:update')
                  ? h(
                      NButton,
                      { size: 'tiny', secondary: true, type: 'info', onClick: () => openEditForm(row) },
                      { default: () => '编辑' },
                    )
                  : null,
                canUse('system:notice:status')
                  ? h(
                      NPopconfirm,
                      { onPositiveClick: () => handleToggleStatus(row, nextStatus) },
                      {
                        trigger: () =>
                          h(
                            NButton,
                            {
                              size: 'tiny',
                              secondary: true,
                              type: nextStatus === NoticeStatus.Disabled ? 'error' : 'success',
                            },
                            { default: () => (nextStatus === NoticeStatus.Disabled ? '禁用' : '启用') },
                          ),
                        default: () => `确认${nextStatus === NoticeStatus.Disabled ? '禁用' : '启用'}该公告？`,
                      },
                    )
                  : null,
                canUse('system:notice:delete')
                  ? h(
                      NPopconfirm,
                      { onPositiveClick: () => handleDelete(row) },
                      {
                        trigger: () =>
                          h(
                            NButton,
                            { size: 'tiny', secondary: true, type: 'error' },
                            { default: () => '删除' },
                          ),
                        default: () => '确认删除该公告？删除后将不再出现在列表里。',
                      },
                    )
                  : null,
              ].filter(Boolean),
          },
        )
      },
    },
  ]

  async function submitForm() {
    if (formMode.value === 'create') {
      await createNotice(buildNoticePayload(formModel))
      message.success('公告创建成功')
    } else {
      await updateNotice(formModel.id, buildNoticePayload(formModel))
      message.success('公告更新成功')
    }
    await load()
  }

  async function handleDelete(row: NoticeItem) {
    await deleteNotice(row.id)
    message.success('公告已删除')
    await load()
  }

  return {
    canUse,
    columns,
    formMode,
    formModel,
    formRef,
    formVisible,
    handlePageChange,
    handlePageSizeChange,
    handleReset,
    handleSearch,
    handleSubmit,
    load,
    loading,
    notices,
    openCreate,
    openEdit: openEditForm,
    query,
    rules,
    saving,
    submitForm,
    total,
  }
}
