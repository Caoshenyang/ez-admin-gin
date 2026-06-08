<script setup lang="ts">
import { NButton } from 'naive-ui'

import PageHeader from '@/components/PageHeader.vue'
import { STATUS_FILTER_OPTIONS } from '@/constants/status'
import AttachmentEditModal from '../components/AttachmentEditModal.vue'
import AttachmentFilterBar from '../components/AttachmentFilterBar.vue'
import AttachmentTable from '../components/AttachmentTable.vue'
import AttachmentUploadModal from '../components/AttachmentUploadModal.vue'
import { useAttachmentPage } from '../composables/useAttachmentPage'

const {
  attachments,
  canUse,
  editFormModel,
  editFormRef,
  editModalVisible,
  editRules,
  editSaving,
  extFilterOptions,
  handleEditSubmit,
  handlePageChange,
  handlePageSizeChange,
  handleReset,
  handleSearch,
  handleToggleStatus,
  handleUpdateFileList,
  hasRows,
  loading,
  openEditModal,
  openUploadModal,
  query,
  resetUploadModal,
  saving,
  submitEdit,
  submitUpload,
  total,
  uploadFileList,
  uploadFormModel,
  uploadFormRef,
  uploadModalVisible,
  uploadRules,
} = useAttachmentPage()

function copyURL(url: string) {
  void navigator.clipboard.writeText(url)
}

async function handleUploadSubmit() {
  await submitUpload()
}

async function handleEditSubmitAction() {
  await handleEditSubmit(submitEdit)
}

async function handleStatusChange(row: Parameters<typeof handleToggleStatus>[0], status: Parameters<typeof handleToggleStatus>[1]) {
  await handleToggleStatus(row, status)
}
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="附件中心" description="复用底层文件上传链路，把附件整理成可分类、可检索、可业务复用的统一资源。">
        <template #actions>
          <NButton v-if="canUse('system:attachment:upload')" type="primary" @click="openUploadModal">
            上传附件
          </NButton>
        </template>
      </PageHeader>

      <AttachmentFilterBar
        :biz-type="query.biz_type ?? ''"
        :category="query.category ?? ''"
        :ext="query.ext ?? ''"
        :ext-options="extFilterOptions"
        :keyword="query.keyword ?? ''"
        :status="query.status ?? 0"
        :status-options="STATUS_FILTER_OPTIONS"
        @update:biz-type="query.biz_type = $event"
        @update:category="query.category = $event"
        @update:ext="query.ext = $event"
        @update:keyword="query.keyword = $event"
        @update:status="query.status = $event as 0 | 1 | 2"
        @search="handleSearch"
        @reset="handleReset"
      />

      <AttachmentTable
        :attachments="attachments"
        :can-use="canUse"
        :has-rows="hasRows"
        :loading="loading"
        :page="query.page"
        :page-size="query.page_size"
        :total="total"
        @copy="copyURL"
        @edit="openEditModal"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        @refresh="handleSearch"
        @toggle-status="handleStatusChange"
      />
    </section>

    <AttachmentUploadModal
      v-model:show="uploadModalVisible"
      v-model:form-ref="uploadFormRef"
      :file-list="uploadFileList"
      :model="uploadFormModel"
      :rules="uploadRules"
      :saving="saving"
      @update:file-list="handleUpdateFileList"
      @update:show="(value) => { uploadModalVisible = value; if (!value) resetUploadModal() }"
      @submit="handleUploadSubmit"
    />

    <AttachmentEditModal
      v-model:show="editModalVisible"
      v-model:form-ref="editFormRef"
      :model="editFormModel"
      :rules="editRules"
      :saving="editSaving"
      @submit="handleEditSubmitAction"
    />
  </main>
</template>
