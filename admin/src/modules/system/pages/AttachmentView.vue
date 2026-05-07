<script setup lang="ts">
import axios from 'axios'
import { AttachOutline, CopyOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules, UploadFileInfo } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NModal,
  NPagination,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  NTooltip,
  NUpload,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'

import {
  createAttachment,
  getAttachments,
  updateAttachment,
  updateAttachmentStatus,
} from '../api/attachment'
import { buttonPermissionCodes } from '@/router/dynamic-menu'
import {
  AttachmentStatus,
  type AttachmentItem,
  type AttachmentListQuery,
  type UpdateAttachmentPayload,
} from '../types/attachment'

interface UploadFormModel {
  display_name: string
  category: string
  biz_type: string
  status: AttachmentStatus
  remark: string
}

interface EditFormModel extends UpdateAttachmentPayload {
  id: number
}

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const attachments = ref<AttachmentItem[]>([])
const total = ref(0)
const uploadModalVisible = ref(false)
const editModalVisible = ref(false)
const selectedUploadFile = ref<File | null>(null)
const uploadFileList = ref<UploadFileInfo[]>([])
const uploadFormRef = ref<FormInst | null>(null)
const editFormRef = ref<FormInst | null>(null)

const query = reactive<AttachmentListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  category: '',
  biz_type: '',
  ext: '',
  status: 0,
})

const uploadFormModel = reactive<UploadFormModel>({
  display_name: '',
  category: '',
  biz_type: '',
  status: AttachmentStatus.Enabled,
  remark: '',
})

const editFormModel = reactive<EditFormModel>({
  id: 0,
  display_name: '',
  category: '',
  biz_type: '',
  status: AttachmentStatus.Enabled,
  remark: '',
})

const statusOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: AttachmentStatus.Enabled },
  { label: '禁用', value: AttachmentStatus.Disabled },
]

const statusFormOptions = [
  { label: '启用', value: AttachmentStatus.Enabled },
  { label: '禁用', value: AttachmentStatus.Disabled },
]

const extFilterOptions = [
  { label: '类型：全部', value: '' },
  { label: '图片', value: '.png' },
  { label: 'PDF', value: '.pdf' },
  { label: 'Excel', value: '.xlsx' },
  { label: 'Word', value: '.docx' },
]

const uploadRules: FormRules = {
  display_name: [{ max: 255, message: '附件名称不能超过 255 个字符', trigger: ['blur', 'input'] }],
  category: [{ max: 64, message: '附件分类不能超过 64 个字符', trigger: ['blur', 'input'] }],
  biz_type: [{ max: 64, message: '业务类型不能超过 64 个字符', trigger: ['blur', 'input'] }],
  remark: [{ max: 255, message: '备注不能超过 255 个字符', trigger: ['blur', 'input'] }],
}

const editRules: FormRules = {
  display_name: [{ required: true, message: '请输入附件名称', trigger: ['blur', 'input'] }],
  category: [{ max: 64, message: '附件分类不能超过 64 个字符', trigger: ['blur', 'input'] }],
  biz_type: [{ max: 64, message: '业务类型不能超过 64 个字符', trigger: ['blur', 'input'] }],
  remark: [{ max: 255, message: '备注不能超过 255 个字符', trigger: ['blur', 'input'] }],
}

const columns: DataTableColumns<AttachmentItem> = [
  {
    title: '附件',
    key: 'display_name',
    minWidth: 260,
    render(row) {
      return h('div', { class: 'flex items-center gap-3' }, [
        h(
          'div',
          {
            class: 'flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-[#EEF2FF] text-[#4F46E5]',
          },
          [h(NIcon, { size: 18 }, { default: () => h(AttachOutline) })],
        ),
        h('div', { class: 'min-w-0 leading-5' }, [
          h('p', { class: 'truncate font-medium text-[#111827]' }, row.display_name),
          h('p', { class: 'truncate text-xs text-[#6B7280]' }, row.original_name),
        ]),
      ])
    },
  },
  {
    title: '分类 / 业务',
    key: 'category',
    width: 180,
    render(row) {
      const text = [row.category || '未分类', row.biz_type || '通用'].join(' / ')
      return h('span', { class: 'text-sm text-[#334155]' }, text)
    },
  },
  {
    title: '类型',
    key: 'ext',
    width: 100,
    render(row) {
      return h(NTag, { size: 'small', bordered: false }, { default: () => row.ext || '-' })
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
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return h(
        NTag,
        { bordered: false, type: row.status === AttachmentStatus.Enabled ? 'success' : 'error' },
        { default: () => (row.status === AttachmentStatus.Enabled ? '启用' : '禁用') },
      )
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
    width: 220,
    fixed: 'right',
    render(row) {
      const nextStatus =
        row.status === AttachmentStatus.Enabled ? AttachmentStatus.Disabled : AttachmentStatus.Enabled

      return h(
        NSpace,
        { size: 8 },
        {
          default: () =>
            [
              h(
                NTooltip,
                {},
                {
                  trigger: () =>
                    h(
                      NButton,
                      { size: 'small', ghost: true, type: 'info', onClick: () => copyURL(row.url) },
                      { icon: () => h(NIcon, null, { default: () => h(CopyOutline) }) },
                    ),
                  default: () => '复制链接',
                },
              ),
              canUse('system:attachment:update')
                ? h(
                    NButton,
                    { size: 'small', ghost: true, type: 'primary', onClick: () => openEditModal(row) },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:attachment:status')
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
                            type: nextStatus === AttachmentStatus.Disabled ? 'error' : 'success',
                          },
                          { default: () => (nextStatus === AttachmentStatus.Disabled ? '禁用' : '启用') },
                        ),
                      default: () => `确认${nextStatus === AttachmentStatus.Disabled ? '禁用' : '启用'}该附件？`,
                    },
                  )
                : null,
            ].filter(Boolean),
        },
      )
    },
  },
]

const hasRows = computed(() => attachments.value.length > 0)

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

function copyURL(url: string) {
  navigator.clipboard.writeText(url).then(
    () => message.success('链接已复制'),
    () => message.error('复制失败'),
  )
}

function openUploadModal() {
  uploadModalVisible.value = true
}

function resetUploadModal() {
  uploadFormModel.display_name = ''
  uploadFormModel.category = ''
  uploadFormModel.biz_type = ''
  uploadFormModel.status = AttachmentStatus.Enabled
  uploadFormModel.remark = ''
  uploadFileList.value = []
  selectedUploadFile.value = null
}

function openEditModal(row: AttachmentItem) {
  editFormModel.id = row.id
  editFormModel.display_name = row.display_name
  editFormModel.category = row.category
  editFormModel.biz_type = row.biz_type
  editFormModel.status = row.status
  editFormModel.remark = row.remark
  editModalVisible.value = true
}

function handleSearch() {
  query.page = 1
  void loadAttachments()
}

function handleReset() {
  query.page = 1
  query.page_size = 10
  query.keyword = ''
  query.category = ''
  query.biz_type = ''
  query.ext = ''
  query.status = 0
  void loadAttachments()
}

function handlePageChange(page: number) {
  query.page = page
  void loadAttachments()
}

function handlePageSizeChange(pageSize: number) {
  query.page = 1
  query.page_size = pageSize
  void loadAttachments()
}

function handleUpdateFileList(fileList: UploadFileInfo[]) {
  uploadFileList.value = fileList.slice(-1)
  selectedUploadFile.value = uploadFileList.value[0]?.file ?? null
}

async function loadAttachments() {
  loading.value = true
  try {
    const data = await getAttachments({
      ...query,
      keyword: query.keyword?.trim() || undefined,
      category: query.category?.trim() || undefined,
      biz_type: query.biz_type?.trim() || undefined,
      ext: query.ext || undefined,
      status: query.status === 0 ? undefined : query.status,
    })
    attachments.value = data.items
    total.value = data.total
  } finally {
    loading.value = false
  }
}

async function handleSubmitUpload() {
  try {
    await uploadFormRef.value?.validate()
  } catch {
    return
  }

  if (!selectedUploadFile.value) {
    message.error('请选择要上传的附件')
    return
  }

  saving.value = true
  try {
    await createAttachment(selectedUploadFile.value, {
      display_name: uploadFormModel.display_name.trim() || undefined,
      category: uploadFormModel.category.trim() || undefined,
      biz_type: uploadFormModel.biz_type.trim() || undefined,
      status: uploadFormModel.status,
      remark: uploadFormModel.remark.trim() || undefined,
    })
    message.success('附件上传成功')
    uploadModalVisible.value = false
    resetUploadModal()
    await loadAttachments()
  } catch (error) {
    const errorMessage = axios.isAxiosError<{ message?: string }>(error)
      ? error.response?.data?.message ?? '附件上传失败'
      : '附件上传失败'
    message.error(errorMessage)
  } finally {
    saving.value = false
  }
}

async function handleSubmitEdit() {
  try {
    await editFormRef.value?.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    await updateAttachment(editFormModel.id, {
      display_name: editFormModel.display_name.trim(),
      category: editFormModel.category.trim(),
      biz_type: editFormModel.biz_type.trim(),
      status: editFormModel.status,
      remark: editFormModel.remark.trim(),
    })
    message.success('附件信息已更新')
    editModalVisible.value = false
    await loadAttachments()
  } catch (error) {
    const errorMessage = axios.isAxiosError<{ message?: string }>(error)
      ? error.response?.data?.message ?? '更新附件失败'
      : '更新附件失败'
    message.error(errorMessage)
  } finally {
    saving.value = false
  }
}

async function handleToggleStatus(row: AttachmentItem, status: AttachmentStatus) {
  try {
    await updateAttachmentStatus(row.id, status)
    message.success(status === AttachmentStatus.Enabled ? '附件已启用' : '附件已禁用')
    await loadAttachments()
  } catch (error) {
    const errorMessage = axios.isAxiosError<{ message?: string }>(error)
      ? error.response?.data?.message ?? '更新附件状态失败'
      : '更新附件状态失败'
    message.error(errorMessage)
  }
}

onMounted(() => {
  void loadAttachments()
})
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">附件中心</h1>
          <p class="mt-1 text-sm text-[#6B7280]">
            复用底层文件上传链路，把附件整理成可分类、可检索、可业务复用的统一资源。
          </p>
        </div>

        <NButton v-if="canUse('system:attachment:upload')" type="primary" @click="openUploadModal">
          上传附件
        </NButton>
      </div>

      <NCard :bordered="false" class="rounded-lg">
        <NSpace align="center" :wrap="true">
          <NInput
            v-model:value="query.keyword"
            clearable
            placeholder="附件名称 / 原始文件名"
            class="w-64"
            @keyup.enter="handleSearch"
          />
          <NInput
            v-model:value="query.category"
            clearable
            placeholder="附件分类"
            class="w-40"
            @keyup.enter="handleSearch"
          />
          <NInput
            v-model:value="query.biz_type"
            clearable
            placeholder="业务类型"
            class="w-40"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.ext" :options="extFilterOptions" class="w-36" />
          <NSelect v-model:value="query.status" :options="statusOptions" class="w-32" />
          <NButton type="primary" @click="handleSearch">查询</NButton>
          <NButton quaternary @click="handleReset">重置</NButton>
        </NSpace>
      </NCard>

      <NCard :bordered="false" class="flex-1 overflow-hidden rounded-lg">
        <div class="flex h-full flex-col">
          <NDataTable
            :columns="columns"
            :data="attachments"
            :loading="loading"
            :bordered="false"
            flex-height
            class="flex-1"
          />

          <div
            v-if="!loading && !hasRows"
            class="flex flex-1 items-center justify-center rounded-2xl border border-dashed border-[#D9DEE8] bg-[#FAFBFC]"
          >
            <NEmpty description="当前没有附件记录，先上传一份业务附件吧。" />
          </div>

          <div class="mt-4 flex justify-end">
            <NPagination
              :page="query.page"
              :page-size="query.page_size"
              :item-count="total"
              show-size-picker
              :page-sizes="[10, 20, 50]"
              @update:page="handlePageChange"
              @update:page-size="handlePageSizeChange"
            />
          </div>
        </div>
      </NCard>
    </section>

    <NModal
      v-model:show="uploadModalVisible"
      preset="card"
      title="上传附件"
      class="max-w-[620px] rounded-2xl"
      :bordered="false"
      @after-leave="resetUploadModal"
    >
      <NForm ref="uploadFormRef" :model="uploadFormModel" :rules="uploadRules" label-placement="top">
        <NFormItem label="附件名称" path="display_name">
          <NInput
            v-model:value="uploadFormModel.display_name"
            placeholder="可留空，默认使用原始文件名"
          />
        </NFormItem>

        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="附件分类" path="category">
            <NInput v-model:value="uploadFormModel.category" placeholder="例如 contract / avatar" />
          </NFormItem>
          <NFormItem label="业务类型" path="biz_type">
            <NInput v-model:value="uploadFormModel.biz_type" placeholder="例如 customer / system-template" />
          </NFormItem>
        </div>

        <NFormItem label="状态" path="status">
          <NSelect v-model:value="uploadFormModel.status" :options="statusFormOptions" />
        </NFormItem>

        <NFormItem label="备注" path="remark">
          <NInput
            v-model:value="uploadFormModel.remark"
            type="textarea"
            :rows="3"
            placeholder="给这份附件补一段业务备注"
          />
        </NFormItem>

        <NFormItem label="上传文件">
          <NUpload
            :default-upload="false"
            :max="1"
            :file-list="uploadFileList"
            @update:file-list="handleUpdateFileList"
          >
            <NButton>选择文件</NButton>
          </NUpload>
        </NFormItem>

        <div class="flex justify-end gap-3">
          <NButton @click="uploadModalVisible = false">取消</NButton>
          <NButton type="primary" :loading="saving" @click="handleSubmitUpload">上传并入库</NButton>
        </div>
      </NForm>
    </NModal>

    <NModal
      v-model:show="editModalVisible"
      preset="card"
      title="编辑附件"
      class="max-w-[620px] rounded-2xl"
      :bordered="false"
    >
      <NForm ref="editFormRef" :model="editFormModel" :rules="editRules" label-placement="top">
        <NFormItem label="附件名称" path="display_name">
          <NInput v-model:value="editFormModel.display_name" placeholder="请输入附件名称" />
        </NFormItem>

        <div class="grid gap-4 md:grid-cols-2">
          <NFormItem label="附件分类" path="category">
            <NInput v-model:value="editFormModel.category" placeholder="附件分类" />
          </NFormItem>
          <NFormItem label="业务类型" path="biz_type">
            <NInput v-model:value="editFormModel.biz_type" placeholder="业务类型" />
          </NFormItem>
        </div>

        <NFormItem label="状态" path="status">
          <NSelect v-model:value="editFormModel.status" :options="statusFormOptions" />
        </NFormItem>

        <NFormItem label="备注" path="remark">
          <NInput v-model:value="editFormModel.remark" type="textarea" :rows="3" placeholder="附件备注" />
        </NFormItem>

        <div class="flex justify-end gap-3">
          <NButton @click="editModalVisible = false">取消</NButton>
          <NButton type="primary" :loading="saving" @click="handleSubmitEdit">保存</NButton>
        </div>
      </NForm>
    </NModal>
  </main>
</template>
