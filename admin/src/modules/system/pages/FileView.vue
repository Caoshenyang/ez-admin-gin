<script setup lang="ts">
import { CloudUploadOutline } from '@vicons/ionicons5'
import { NAlert, NButton, NIcon, NUpload } from 'naive-ui'

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
  successText,
  total,
  uploading,
} = useFilePage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">文件管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">上传和管理系统附件，支持图片和常见文档格式。</p>
        </div>

        <NUpload v-if="canUse('system:file:upload')" :show-file-list="false" :custom-request="handleUpload" :disabled="uploading">
          <NButton type="primary" :loading="uploading">
            <template #icon>
              <NIcon><CloudUploadOutline /></NIcon>
            </template>
            上传文件
          </NButton>
        </NUpload>
      </div>

      <NAlert v-if="successText" type="success" :show-icon="true" closable class="mx-auto w-full max-w-[520px]" @close="successText = ''">
        {{ successText }}
      </NAlert>

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
