import type { DataTableColumns, FormRules } from 'naive-ui'
import { NButton, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { h } from 'vue'

import { useListLoader } from '@/composables/useListLoader'
import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { useSuccessFeedback } from '@/composables/useSuccessFeedback'
import { displayText, formatTime } from '@/utils/format'
import { createPost, getPosts, updatePost, updatePostStatus } from '../api/post'
import { PostStatus, type PostItem } from '../types/post'
import type { PostFormModel, PostPageQuery } from '../types/post-page'
import { buildPostPayload, defaultPostFormModel, defaultPostQuery, toPostFormModel } from './post-page.utils'

const postFormRules: FormRules = {
  code: [{ required: true, message: '请输入岗位编码', trigger: ['blur', 'input'] }],
  name: [{ required: true, message: '请输入岗位名称', trigger: ['blur', 'input'] }],
}

export function usePostPage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const { closeSuccess, showSuccess, successText } = useSuccessFeedback()

  const {
    items: posts,
    loading,
    query,
    load,
    handleSearch,
    handleReset,
  } = useListLoader<PostItem, PostPageQuery>(getPosts, defaultPostQuery())

  const {
    formRef,
    formVisible,
    formMode,
    formModel,
    saving,
    rules,
    openCreate: openCreateBase,
    openEdit: openEditBase,
    handleSubmit: handleSubmitBase,
  } = useModalForm<PostFormModel>(defaultPostFormModel, { rules: postFormRules })

  const { handleToggleStatus: handleToggleStatusBase } = useStatusToggle(updatePostStatus, {
    onSuccess: async () => {
      showSuccess('岗位状态已更新')
      await load()
    },
  })

  const columns: DataTableColumns<PostItem> = [
    {
      title: '岗位',
      key: 'name',
      minWidth: 240,
      render(row) {
        return h('div', { class: 'leading-6' }, [
          h('p', { class: 'font-semibold text-[var(--ez-text-heading)]' }, displayText(row.name)),
          h('p', { class: 'text-xs text-[var(--ez-text-muted)]' }, displayText(row.code)),
        ])
      },
    },
    {
      title: '排序',
      key: 'sort',
      width: 90,
    },
    {
      title: '状态',
      key: 'status',
      width: 100,
      render(row) {
        return h(
          NTag,
          { bordered: false, type: row.status === PostStatus.Enabled ? 'success' : 'error' },
          { default: () => (row.status === PostStatus.Enabled ? '启用' : '禁用') },
        )
      },
    },
    {
      title: '更新时间',
      key: 'updated_at',
      width: 180,
      render(row) {
        return formatTime(row.updated_at)
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 180,
      render(row) {
        const nextStatus = row.status === PostStatus.Enabled ? PostStatus.Disabled : PostStatus.Enabled

        return h(
          NSpace,
          { size: 8 },
          {
            default: () => [
              canUse('system:post:update')
                ? h(
                    NButton,
                    { size: 'small', ghost: true, type: 'info', onClick: () => openEdit(row) },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:post:status')
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
                            type: nextStatus === PostStatus.Disabled ? 'error' : 'success',
                          },
                          { default: () => (nextStatus === PostStatus.Disabled ? '禁用' : '启用') },
                        ),
                      default: () => `确认${nextStatus === PostStatus.Disabled ? '禁用' : '启用'}该岗位？`,
                    },
                  )
                : null,
            ],
          },
        )
      },
    },
  ]

  function openCreate() {
    openCreateBase()
  }

  function openEdit(row: PostItem) {
    openEditBase(toPostFormModel(row))
  }

  async function submitForm() {
    const payload = buildPostPayload(formModel)

    if (formMode.value === 'create') {
      await createPost(payload)
      showSuccess('岗位创建成功')
      message.success('岗位创建成功')
    } else {
      await updatePost(formModel.id, payload)
      showSuccess('岗位更新成功')
      message.success('岗位更新成功')
    }

    await load()
  }

  async function handleSubmit() {
    await handleSubmitBase(submitForm)
  }

  async function handleToggleStatus(row: PostItem, status: PostStatus) {
    await handleToggleStatusBase(row, status)
  }

  return {
    canUse,
    closeSuccess,
    columns,
    formMode,
    formModel,
    formRef,
    formVisible,
    handleReset,
    handleSearch,
    handleSubmit,
    handleToggleStatus,
    loading,
    openCreate,
    openEdit,
    posts,
    query,
    rules,
    saving,
    successText,
  }
}
