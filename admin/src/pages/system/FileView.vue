<script setup lang="ts">
import { CloudUploadOutline, CopyOutline, DocumentOutline, ImageOutline } from '@vicons/ionicons5'
import type { DataTableColumns, UploadFileInfo } from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NDataTable,
  NIcon,
  NInput,
  NPagination,
  NSelect,
  NSpace,
  NTag,
  NTooltip,
  NUpload,
  useMessage,
} from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

import { getFiles, uploadFile } from '../../api/file'
import { buttonPermissionCodes } from '../../router/dynamic-menu'
import { FileStatus, type FileItem, type FileListQuery } from '../../types/file'

const message = useMessage()
const loading = ref(false)
const files = ref<FileItem[]>([])
const total = ref(0)
const successText = ref('')
const uploading = ref(false)

const query = reactive<FileListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  ext: '',
  status: 0,
})

const statusFilterOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: FileStatus.Enabled },
  { label: '禁用', value: FileStatus.Disabled },
]

const extFilterOptions = [
  { label: '类型：全部', value: '' },
  { label: '图片', value: '.png' },
  { label: 'JPG', value: '.jpg' },
  { label: 'PDF', value: '.pdf' },
  { label: 'Excel', value: '.xlsx' },
  { label: 'Word', value: '.docx' },
]

const imageExts = ['.png', '.jpg', '.jpeg', '.gif', '.webp', '.svg']

const columns: DataTableColumns<FileItem> = [
  {
    title: '文件',
    key: 'original_name',
    minWidth: 260,
    render(row) {
      const isImage = imageExts.includes(row.ext.toLowerCase())
      return h('div', { class: 'flex items-center gap-3' }, [
        h(
          'div',
          {
            class:
              'flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg',
            style: isImage ? 'background:#f0f7ff;color:#3b82f6' : 'background:#f3f4f6;color:#6b7280',
          },
          [
            h(
              NIcon,
              { size: 18 },
              { default: () => h(isImage ? ImageOutline : DocumentOutline) },
            ),
          ],
        ),
        h('div', { class: 'min-w-0 leading-5' }, [
          h('p', { class: 'truncate font-medium text-[#111827]' }, row.original_name),
          h('p', { class: 'truncate text-xs text-[#6B7280]' }, row.mime_type),
        ]),
      ])
    },
  },
  {
    title: '类型',
    key: 'ext',
    width: 100,
    render(row) {
      return h(NTag, { size: 'small', bordered: false }, { default: () => row.ext })
    },
  },
  {
    title: '大小',
    key: 'size',
    width: 110,
    render(row) {
      return formatSize(row.size)
    },
  },
  {
    title: '上传时间',
    key: 'created_at',
    width: 180,
    render(row) {
      return formatTime(row.created_at)
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    fixed: 'right',
    render(row) {
      return h(
        NTooltip,
        {},
        {
          trigger: () =>
            h(
              NButton,
              { size: 'small', ghost: true, type: 'info', onClick: () => copyURL(row) },
              { icon: () => h(NIcon, null, { default: () => h(CopyOutline) }) },
            ),
          default: () => '复制链接',
        },
      )
    },
  },
]

function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}

function formatTime(value: string) {
  return value ? new Date(value).toLocaleString() : '-'
}

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function copyURL(row: FileItem) {
  navigator.clipboard.writeText(row.url).then(
    () => message.success('链接已复制'),
    () => message.error('复制失败'),
  )
}

function handleSearch() {
  query.page = 1
  void loadFiles()
}

function handleReset() {
  query.page = 1
  query.page_size = 10
  query.keyword = ''
  query.ext = ''
  query.status = 0
  void loadFiles()
}

function handlePageChange(page: number) {
  query.page = page
  void loadFiles()
}

function handlePageSizeChange(pageSize: number) {
  query.page = 1
  query.page_size = pageSize
  void loadFiles()
}

async function loadFiles() {
  loading.value = true
  try {
    const data = await getFiles({
      ...query,
      keyword: query.keyword?.trim() || undefined,
      ext: query.ext || undefined,
      status: query.status === 0 ? undefined : query.status,
    })
    files.value = data.items
    total.value = data.total
  } finally {
    loading.value = false
  }
}

async function handleUpload({ file }: { file: UploadFileInfo }) {
  if (!file.file) return

  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', file.file)
    await uploadFile(formData)
    successText.value = `文件 ${file.name} 上传成功`
    message.success('文件上传成功')
    await loadFiles()
  } catch {
    message.error('文件上传失败')
  } finally {
    uploading.value = false
  }
}

onMounted(() => {
  void loadFiles()
})
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-4 overflow-hidden">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">文件管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">上传和管理系统附件，支持图片和常见文档格式。</p>
        </div>

        <NUpload
          v-if="canUse('system:file:upload')"
          :show-file-list="false"
          :custom-request="handleUpload"
          :disabled="uploading"
        >
          <NButton type="primary" :loading="uploading">
            <template #icon>
              <NIcon><CloudUploadOutline /></NIcon>
            </template>
            上传文件
          </NButton>
        </NUpload>
      </div>

      <NAlert
        v-if="successText"
        type="success"
        :show-icon="true"
        closable
        class="mx-auto w-full max-w-[520px]"
        @close="successText = ''"
      >
        {{ successText }}
      </NAlert>

      <NCard :bordered="false" class="rounded-lg">
        <NSpace align="center" :wrap="true">
          <NInput
            v-model:value="query.keyword"
            clearable
            placeholder="文件名"
            class="w-56"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.ext" :options="extFilterOptions" class="w-36" />
          <NSelect v-model:value="query.status" :options="statusFilterOptions" class="w-36" />
          <NButton type="primary" @click="handleSearch">查询</NButton>
          <NButton @click="handleReset">重置</NButton>
        </NSpace>
      </NCard>

      <NCard
        class="min-h-0 flex-1 rounded-lg"
        :bordered="false"
        content-style="height: 100%; padding: 0;"
      >
        <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-3">
          <span class="text-sm text-[#6B7280]">共 {{ total }} 个文件</span>
          <NButton text type="primary" @click="loadFiles">刷新</NButton>
        </div>

        <NDataTable
          remote
          class="file-table h-full"
          style="height: calc(100% - 105px)"
          :columns="columns"
          :data="files"
          :loading="loading"
          :pagination="false"
          :row-key="(row: FileItem) => row.id"
          :bordered="false"
          flex-height
        />

        <div
          class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]"
        >
          <span>共 {{ total }} 个文件</span>
          <NPagination
            :page="query.page"
            :page-size="query.page_size"
            :item-count="total"
            :page-sizes="[10, 20, 50]"
            show-size-picker
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </NCard>
    </section>
  </main>
</template>

<style scoped>
.file-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #374151;
  background: #fff;
}

.file-table :deep(.n-data-table-td) {
  color: #374151;
}

.file-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}
</style>
