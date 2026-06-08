<script setup lang="ts">
import { NButton } from 'naive-ui'

import PageHeader from '@/components/PageHeader.vue'
import PostFilterBar from '../components/PostFilterBar.vue'
import PostFormModal from '../components/PostFormModal.vue'
import PostTable from '../components/PostTable.vue'
import { usePostPage } from '../composables/usePostPage'

const {
  canUse,
  checkedRowKeys,
  columns,
  formMode,
  formModel,
  formRef,
  formVisible,
  handleCheckedRowKeysChange,
  handleDeleteSelected,
  handleReset,
  handleSearch,
  handleSubmit,
  loading,
  openCreate,
  posts,
  query,
  rules,
  saving,
  selectedCount,
} = usePostPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="岗位管理" description="收口岗位基础信息，给用户归属、协作流程和扩展模块提供统一的岗位字典。">
        <template #actions>
          <NButton v-if="canUse('system:post:create')" type="primary" @click="openCreate">
            + 新增岗位
          </NButton>
        </template>
      </PageHeader>

      <PostFilterBar v-model:query="query" @search="handleSearch" @reset="handleReset" />

      <PostTable
        :can-use="canUse"
        :checked-row-keys="checkedRowKeys"
        :columns="columns"
        :items="posts"
        :loading="loading"
        :selected-count="selectedCount"
        @checked-row-keys-change="handleCheckedRowKeysChange"
        @delete-selected="handleDeleteSelected"
        @refresh="handleSearch"
      />
    </section>

    <PostFormModal
      v-model:show="formVisible"
      v-model:form-ref="formRef"
      v-model:model="formModel"
      :form-mode="formMode"
      :rules="rules"
      :saving="saving"
      @submit="handleSubmit"
    />
  </main>
</template>
