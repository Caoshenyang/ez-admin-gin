import type { DataTableColumns, FormRules } from 'naive-ui'
import { NButton, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { h } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useSuccessFeedback } from '@/composables/useSuccessFeedback'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { displayText, formatTime } from '@/utils/format'
import { createNotice, getNotices, updateNotice, updateNoticeStatus } from '../api/notice'
import { NoticeStatus, type NoticeItem, type NoticeListQuery } from '../types/notice'
import {
  buildNoticePayload,
  defaultNoticeFormModel,
  defaultNoticeListQuery,
  toNoticeFormModel,
  type NoticeFormModel,
} from './notice-page.utils'

// 公告管理页面组合式函数，封装公告列表、创建、编辑、状态切换等逻辑
export function useNoticePage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const { closeSuccess, showSuccess, successText } = useSuccessFeedback()

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
      showSuccess('公告状态已更新')
      await load()
    },
  })

  // 公告列表表格列定义
  const columns: DataTableColumns<NoticeItem> = [
    {
      title: '标题',
      key: 'title',
      width: 220,
      ellipsis: { tooltip: true },
      render(row) {
        return h('span', { class: 'font-semibold text-[#111827]' }, displayText(row.title))
      },
    },
    {
      title: '内容',
      key: 'content',
      minWidth: 240,
      ellipsis: { tooltip: true },
      render(row) {
        return h('span', { class: 'text-[#374151]' }, displayText(row.content))
      },
    },
    {
      title: '排序',
      key: 'sort',
      width: 80,
      align: 'center',
    },
    {
      title: '状态',
      key: 'status',
      width: 90,
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
      width: 160,
      render(row) {
        return h('span', { class: 'tabular-nums text-[#6B7280]' }, formatTime(row.updated_at))
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 180,
      fixed: 'right',
      render(row) {
        const nextStatus = row.status === NoticeStatus.Enabled ? NoticeStatus.Disabled : NoticeStatus.Enabled

        return h(
          NSpace,
          { size: 8, align: 'center' },
          {
            default: () =>
              [
                canUse('system:notice:update')
                  ? h(
                      NButton,
                      { size: 'small', ghost: true, type: 'info', onClick: () => openEditForm(row) },
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
                              size: 'small',
                              ghost: true,
                              type: nextStatus === NoticeStatus.Disabled ? 'error' : 'success',
                            },
                            { default: () => (nextStatus === NoticeStatus.Disabled ? '禁用' : '启用') },
                          ),
                        default: () => `确认${nextStatus === NoticeStatus.Disabled ? '禁用' : '启用'}该公告？`,
                      },
                    )
                  : null,
              ].filter(Boolean),
          },
        )
      },
    },
  ]

  // 提交公告表单（新建或更新）
  async function submitForm() {
    if (formMode.value === 'create') {
      await createNotice(buildNoticePayload(formModel))
      showSuccess('公告创建成功')
      message.success('公告创建成功')
    } else {
      await updateNotice(formModel.id, buildNoticePayload(formModel))
      showSuccess('公告已更新')
      message.success('公告更新成功')
    }
    await load()
  }

  return {
    canUse,
    closeSuccess,
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
    successText,
    total,
  }
}
