import type { DataTableColumns, DataTableRowKey, FormRules } from 'naive-ui'
import { NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { computed, h, ref } from 'vue'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import { useListLoader } from '@/composables/useListLoader'
import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { displayText, formatTime } from '@/utils/format'
import { createPost, deletePost, getPosts, updatePost, updatePostStatus } from '../api/post'
import { PostStatus, type PostItem } from '../types/post'
import type { PostFormModel, PostPageQuery } from '../types/post-page'
import {
  buildPostPayload,
  defaultPostFormModel,
  defaultPostQuery,
  toPostFormModel,
} from './post-page.utils'

const postFormRules: FormRules = {
  code: [{ required: true, message: '请输入岗位编码', trigger: ['blur', 'input'] }],
  name: [{ required: true, message: '请输入岗位名称', trigger: ['blur', 'input'] }],
}

export function usePostPage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const checkedRowKeys = ref<DataTableRowKey[]>([])

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
      await load()
    },
  })

  const columns: DataTableColumns<PostItem> = [
    { type: 'selection', width: 48 },
    {
      title: '岗位名称',
      key: 'name',
      minWidth: 180,
      render(row) {
        return h(
          'span',
          { class: 'font-semibold text-[var(--ez-text-heading)]' },
          displayText(row.name),
        )
      },
    },
    {
      title: '岗位编码',
      key: 'code',
      minWidth: 150,
      ellipsis: { tooltip: true },
      render(row) {
        return h(
          'span',
          { class: 'font-mono text-[var(--ez-text-regular)]' },
          displayText(row.code),
        )
      },
    },
    {
      title: '排序',
      key: 'sort',
      width: 72,
    },
    {
      title: '状态',
      key: 'status',
      width: 84,
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
      width: 150,
      render(row) {
        return formatTime(row.updated_at)
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 128,
      fixed: 'right',
      render(row) {
        const nextStatus =
          row.status === PostStatus.Enabled ? PostStatus.Disabled : PostStatus.Enabled

        return h(
          NSpace,
          { class: 'ez-row-actions', size: 6, align: 'center' },
          {
            default: () => [
              canUse('system:post:update')
                ? h(EzActionButton, {
                    iconOnly: true,
                    kind: 'edit',
                    label: '编辑',
                    size: 'tiny',
                    secondary: true,
                    type: 'info',
                    onClick: () => openEdit(row),
                  })
                : null,
              canUse('system:post:status')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleToggleStatus(row, nextStatus) },
                    {
                      trigger: () =>
                        h(EzActionButton, {
                          iconOnly: true,
                          kind: nextStatus === PostStatus.Disabled ? 'disable' : 'enable',
                          label: nextStatus === PostStatus.Disabled ? '禁用' : '启用',
                          size: 'tiny',
                          secondary: true,
                          tooltip: false,
                          type: nextStatus === PostStatus.Disabled ? 'error' : 'success',
                        }),
                      default: () =>
                        `确认${nextStatus === PostStatus.Disabled ? '禁用' : '启用'}该岗位？`,
                    },
                  )
                : null,
              canUse('system:post:delete')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleDelete(row) },
                    {
                      trigger: () =>
                        h(EzActionButton, {
                          iconOnly: true,
                          kind: 'delete',
                          label: '删除',
                          size: 'tiny',
                          secondary: true,
                          tooltip: false,
                          type: 'error',
                        }),
                      default: () => '删除前请确认该岗位没有绑定到任何用户。',
                    },
                  )
                : null,
            ],
          },
        )
      },
    },
  ]

  const selectedCount = computed(() => checkedRowKeys.value.length)

  function openCreate() {
    openCreateBase()
  }

  function openEdit(row: PostItem) {
    openEditBase(toPostFormModel(row))
  }

  async function submitForm() {
    const payload = buildPostPayload(formModel)
    let createdPost: PostItem | null = null

    if (formMode.value === 'create') {
      createdPost = await createPost(payload)
      message.success('岗位创建成功')
      Object.assign(query, defaultPostQuery())
    } else {
      await updatePost(formModel.id, payload)
      message.success('岗位更新成功')
    }

    await load()
    if (createdPost && !posts.value.some((post) => post.id === createdPost.id)) {
      posts.value = [createdPost, ...posts.value].sort(
        (left, right) => left.sort - right.sort || left.id - right.id,
      )
    }
  }

  async function handleSubmit() {
    await handleSubmitBase(submitForm)
  }

  async function handleToggleStatus(row: PostItem, status: PostStatus) {
    await handleToggleStatusBase(row, status)
  }

  async function handleDelete(row: PostItem) {
    await deletePost(row.id)
    checkedRowKeys.value = checkedRowKeys.value.filter((key) => key !== row.id)
    message.success('岗位已删除')
    await load()
  }

  async function handleDeleteSelected() {
    const ids = checkedRowKeys.value.filter((key): key is number => typeof key === 'number')
    if (ids.length === 0) {
      return
    }

    await Promise.all(ids.map((id) => deletePost(id)))
    checkedRowKeys.value = []
    message.success(`已删除 ${ids.length} 个岗位`)
    await load()
  }

  function handleCheckedRowKeysChange(keys: DataTableRowKey[]) {
    checkedRowKeys.value = keys
  }

  return {
    canUse,
    checkedRowKeys,
    columns,
    formMode,
    formModel,
    formRef,
    formVisible,
    handleReset,
    handleSearch,
    handleSubmit,
    handleCheckedRowKeysChange,
    handleDeleteSelected,
    handleToggleStatus,
    loading,
    openCreate,
    openEdit,
    posts,
    query,
    rules,
    saving,
    selectedCount,
  }
}
