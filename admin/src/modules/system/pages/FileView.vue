<script setup lang="ts">
import { CloudUploadOutline } from '@vicons/ionicons5'
import { NButton, NIcon, NUpload } from 'naive-ui'

import PageHeader from '@/components/PageHeader.vue'
import FileFilterBar from '../components/FileFilterBar.vue'
import FileTable from '../components/FileTable.vue'
import { useFilePage } from '../composables/useFilePage'

const {
  canUse,
  columns,
  extFilterOptions,
  files,
  handlePageChange,
  handlePageSizeChange,
  handleReset,
  handleSearch,
  handleUpload,
  load,
  loading,
  query,
  total,
  uploading,
} = useFilePage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="文件管理" description="上传和管理系统附件，支持图片和常见文档格式。">
        <template #actions>
          <NUpload v-if="canUse('system:file:upload')" :show-file-list="false" :custom-request="handleUpload" :disabled="uploading">
            <NButton type="primary" :loading="uploading">
              <template #icon>
                <NIcon><CloudUploadOutline /></NIcon>
              </template>
              上传文件
            </NButton>
          </NUpload>
        </template>
      </PageHeader>

      <FileFilterBar v-model:query="query" :ext-filter-options="extFilterOptions" @search="handleSearch" @reset="handleReset" />

      <FileTable
        :columns="columns"
        :items="files"
        :loading="loading"
        :query="query"
        :total="total"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        @refresh="load"
      />
    </section>
  </main>
</template>
