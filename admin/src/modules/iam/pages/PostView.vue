<script setup lang="ts">
import type { DataTableColumns, FormRules } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  useMessage,
} from 'naive-ui'
import { h } from 'vue'

import { useListLoader } from '@/composables/useListLoader'
import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { STATUS_FILTER_OPTIONS, STATUS_FORM_OPTIONS } from '@/constants/status'
import { displayText, formatTime } from '@/utils/format'
import { createPost, getPosts, updatePost, updatePostStatus } from '../api/post'
import { PostStatus, type CreatePostPayload, type PostItem, type PostListQuery } from '../types/post'

interface PostFormModel extends CreatePostPayload {
  id: number
}

const message = useMessage()
const { canUse } = usePermission()

const { items: posts, loading, query, load, handleSearch, handleReset } = useListLoader<PostItem, PostListQuery>(
  getPosts,
  { keyword: '', status: 0 },
)

// defaultFormModel 函数。
function defaultFormModel(): PostFormModel {
  return { id: 0, code: '', name: '', sort: 0, status: PostStatus.Enabled, remark: '' }
}

const { formRef, formVisible, formMode, formModel, saving, rules, openCreate, openEdit, handleSubmit } =
  useModalForm<PostFormModel>(defaultFormModel, {
    rules: {
      code: [{ required: true, message: '请输入岗位编码', trigger: ['blur', 'input'] }],
      name: [{ required: true, message: '请输入岗位名称', trigger: ['blur', 'input'] }],
    } as FormRules,
  })

const { handleToggleStatus } = useStatusToggle(updatePostStatus, { onSuccess: load })

const columns: DataTableColumns<PostItem> = [
  {
    title: '岗位',
    key: 'name',
    minWidth: 240,
    render(row) {
      return h('div', { class: 'leading-6' }, [
        h('p', { class: 'font-semibold text-[#111827]' }, displayText(row.name)),
        h('p', { class: 'text-xs text-[#6B7280]' }, displayText(row.code)),
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

// onSubmit 函数。
async function onSubmit() {
  const payload = {
    code: formModel.code,
    name: formModel.name,
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark,
  }
  if (formMode.value === 'create') {
    await createPost(payload)
    message.success('岗位创建成功')
  } else {
    await updatePost(formModel.id, payload)
    message.success('岗位更新成功')
  }
  await load()
}
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">岗位管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">收口岗位基础信息，给用户归属、协作流程和扩展模块提供统一的岗位字典。</p>
        </div>

        <NButton v-if="canUse('system:post:create')" type="primary" @click="openCreate">
          + 新增岗位
        </NButton>
      </div>

      <NCard :bordered="false" class="rounded-lg">
        <NSpace align="center" :wrap="true">
          <NInput
            v-model:value="query.keyword"
            clearable
            placeholder="搜索岗位名称 / 编码"
            class="w-64"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" class="w-36" />
          <NButton type="primary" @click="handleSearch">查询</NButton>
          <NButton @click="handleReset">重置</NButton>
        </NSpace>
      </NCard>

      <NCard class="min-h-0 flex-1 rounded-lg" :bordered="false" content-style="display: flex; height: 100%; min-height: 0; flex-direction: column; padding: 0;">
        <NDataTable
          class="h-full"
          style="height: 100%"
          :columns="columns"
          :data="posts"
          :loading="loading"
          :pagination="false"
          :bordered="false"
          flex-height
        />
      </NCard>
    </section>

    <NModal v-model:show="formVisible" preset="card" class="w-[640px]" :title="formMode === 'create' ? '新增岗位' : '编辑岗位'">
      <NForm ref="formRef" :model="formModel" :rules="rules" label-placement="top" class="grid gap-4 md:grid-cols-2">
        <NFormItem label="岗位编码" path="code">
          <NInput v-model:value="formModel.code" placeholder="请输入岗位编码" />
        </NFormItem>

        <NFormItem label="岗位名称" path="name">
          <NInput v-model:value="formModel.name" placeholder="请输入岗位名称" />
        </NFormItem>

        <NFormItem label="排序" path="sort">
          <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
        </NFormItem>

        <NFormItem label="状态" path="status">
          <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
        </NFormItem>

        <NFormItem label="备注" path="remark" class="md:col-span-2">
          <NInput v-model:value="formModel.remark" type="textarea" :rows="4" placeholder="补充记录岗位职责、适用部门或协作边界" />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="flex justify-end gap-3">
          <NButton quaternary @click="formVisible = false">取消</NButton>
          <NButton type="primary" :loading="saving" @click="handleSubmit(onSubmit)">保存</NButton>
        </div>
      </template>
    </NModal>
  </main>
</template>
