<script setup lang="ts">
import EzActionButton from '@/components/ez/EzActionButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import NoticeFilterBar from '../components/NoticeFilterBar.vue'
import NoticeFormModal from '../components/NoticeFormModal.vue'
import NoticeTable from '../components/NoticeTable.vue'
import { useNoticePage } from '../composables/useNoticePage'

const {
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
  total,
} = useNoticePage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="公告管理" description="管理系统公告，支持按标题搜索和状态筛选。">
        <template #actions>
          <EzActionButton
            v-if="canUse('system:notice:create')"
            kind="add"
            label="新增公告"
            type="primary"
            @click="openCreate"
          />
        </template>
      </PageHeader>

      <NoticeFilterBar v-model:query="query" @search="handleSearch" @reset="handleReset" />

      <NoticeTable
        :columns="columns"
        :items="notices"
        :loading="loading"
        :query="query"
        :total="total"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        @refresh="load"
      />
    </section>

    <NoticeFormModal
      v-model:show="formVisible"
      v-model:form-ref="formRef"
      :form-mode="formMode"
      :model="formModel"
      :rules="rules"
      :saving="saving"
      @submit="handleSubmit(submitForm)"
    />
  </main>
</template>
