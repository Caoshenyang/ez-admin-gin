<script setup lang="ts">
import { NAlert, NButton } from 'naive-ui'

import PostFilterBar from '../components/PostFilterBar.vue'
import PostFormModal from '../components/PostFormModal.vue'
import PostTable from '../components/PostTable.vue'
import { usePostPage } from '../composables/usePostPage'

const {
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
  loading,
  openCreate,
  posts,
  query,
  rules,
  saving,
  successText,
} = usePostPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <div class="flex items-center justify-between">
                  <div class="ez-page-header">
            <h1>岗位管理</h1>
          <p>收口岗位基础信息，给用户归属、协作流程和扩展模块提供统一的岗位字典。</p>
        </div>

        <NButton v-if="canUse('system:post:create')" type="primary" @click="openCreate">
          + 新增岗位
        </NButton>
      </div>

      <NAlert v-if="successText" type="success" :show-icon="true" closable class="mx-auto w-full max-w-[520px]" @close="closeSuccess">
        {{ successText }}
      </NAlert>

      <PostFilterBar v-model:query="query" @search="handleSearch" @reset="handleReset" />

      <PostTable :columns="columns" :items="posts" :loading="loading" @refresh="handleSearch" />
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
