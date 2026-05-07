import type { DataTableColumns, FormRules } from 'naive-ui'
import { NButton, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { h, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { formatTime } from '@/utils/format'
import { createNotice, getNotices, updateNotice, updateNoticeStatus } from '../api/notice'
import { NoticeStatus, type NoticeItem, type NoticeListQuery } from '../types/notice'

export interface NoticeFormModel {
  id: number
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

function defaultFormModel(): NoticeFormModel {
  return { id: 0, title: '', content: '', sort: 0, status: NoticeStatus.Enabled, remark: '' }
}

export function useNoticePage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const successText = ref('')

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
    page: 1,
    page_size: 10,
    keyword: '',
    status: 0,
  })

  const { formRef, formVisible, formMode, formModel, saving, rules, openCreate, openEdit, handleSubmit } =
    useModalForm<NoticeFormModel>(defaultFormModel, {
      rules: { title: [{ required: true, message: '请输入公告标题', trigger: 'blur' }] } as FormRules,
    })

  const { handleToggleStatus } = useStatusToggle(updateNoticeStatus, {
    onSuccess: async () => {
      successText.value = '公告状态已更新'
      await load()
    },
  })

  const columns: DataTableColumns<NoticeItem> = [
    {
      title: '标题',
      key: 'title',
      width: 220,
      ellipsis: { tooltip: true },
      render(row) {
        return h('span', { class: 'font-semibold text-[#111827]' }, row.title)
      },
    },
    {
      title: '内容',
      key: 'content',
      minWidth: 240,
      ellipsis: { tooltip: true },
      render(row) {
        return h('span', { class: 'text-[#374151]' }, row.content || '-')
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
                      { size: 'small', ghost: true, type: 'info', onClick: () => openEdit(row) },
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

  async function submitForm() {
    if (formMode.value === 'create') {
      await createNotice({
        title: formModel.title,
        content: formModel.content,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      successText.value = '公告创建成功'
      message.success('公告创建成功')
    } else {
      await updateNotice(formModel.id, {
        title: formModel.title,
        content: formModel.content,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      successText.value = '公告已更新'
      message.success('公告更新成功')
    }
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
    query,
    rules,
    saving,
    submitForm,
    successText,
    total,
  }
}
