<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NInputNumber,
  NModal,
  NPagination,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  useMessage,
} from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

import { createNotice, getNotices, updateNotice, updateNoticeStatus } from '../../api/notice'
import { buttonPermissionCodes } from '../../router/dynamic-menu'
import {
  NoticeStatus,
  type NoticeItem,
  type NoticeListQuery,
} from '../../types/notice'

interface NoticeFormModel {
  id: number
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const notices = ref<NoticeItem[]>([])
const total = ref(0)

const query = reactive<NoticeListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  status: 0,
})

const formRef = ref<FormInst | null>(null)
const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const formModel = reactive<NoticeFormModel>({
  id: 0,
  title: '',
  content: '',
  sort: 0,
  status: NoticeStatus.Enabled,
  remark: '',
})

const statusFilterOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: NoticeStatus.Enabled },
  { label: '禁用', value: NoticeStatus.Disabled },
]

const statusFormOptions = [
  { label: '启用', value: NoticeStatus.Enabled },
  { label: '禁用', value: NoticeStatus.Disabled },
]

const rules: FormRules = {
  title: [{ required: true, message: '请输入公告标题', trigger: 'blur' }],
}

const columns: DataTableColumns<NoticeItem> = [
  {
    title: '标题',
    key: 'title',
    width: 220,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'font-semibold text-[#111827]' }, row.title)
    },
  },
  {
    title: '内容',
    key: 'content',
    minWidth: 240,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[#374151]' }, row.content || '-')
    },
  },
  {
    title: '排序',
    key: 'sort',
    width: 80,
    align: 'center',
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: row.status === NoticeStatus.Enabled ? 'success' : 'error', bordered: false },
        { default: () => (row.status === NoticeStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 160,
    render(row) {
      return h('span', { class: 'tabular-nums text-[#6B7280]' }, formatTime(row.updated_at))
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render(row) {
      const nextStatus = row.status === NoticeStatus.Enabled ? NoticeStatus.Disabled : NoticeStatus.Enabled

      return h(
        NSpace,
        { size: 8, align: 'center' },
        {
          default: () =>
            [
              canUse('system:notice:update')
                ? h(
                    NButton,
                    { size: 'small', ghost: true, type: 'info', onClick: () => openEdit(row) },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:notice:status')
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
                            type: nextStatus === NoticeStatus.Disabled ? 'error' : 'success',
                          },
                          { default: () => (nextStatus === NoticeStatus.Disabled ? '禁用' : '启用') },
                        ),
                      default: () => `确认${nextStatus === NoticeStatus.Disabled ? '禁用' : '启用'}该公告？`,
                    },
                  )
                : null,
            ].filter(Boolean),
        },
      )
    },
  },
]

function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}

function formatTime(value: string) {
  if (!value) return '-'
  const d = new Date(value)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function resetForm() {
  Object.assign(formModel, {
    id: 0,
    title: '',
    content: '',
    sort: 0,
    status: NoticeStatus.Enabled,
    remark: '',
  })
}

function handleSearch() {
  query.page = 1
  void loadNotices()
}

function handleReset() {
  query.page = 1
  query.page_size = 10
  query.keyword = ''
  query.status = 0
  void loadNotices()
}

function handlePageChange(page: number) {
  query.page = page
  void loadNotices()
}

function handlePageSizeChange(pageSize: number) {
  query.page = 1
  query.page_size = pageSize
  void loadNotices()
}

function openCreate() {
  formMode.value = 'create'
  resetForm()
  formVisible.value = true
}

function openEdit(row: NoticeItem) {
  formMode.value = 'edit'
  Object.assign(formModel, {
    id: row.id,
    title: row.title,
    content: row.content,
    sort: row.sort,
    status: row.status,
    remark: row.remark,
  })
  formVisible.value = true
}

async function loadNotices() {
  loading.value = true
  try {
    const data = await getNotices({
      ...query,
      keyword: query.keyword?.trim() || undefined,
      status: query.status === 0 ? undefined : query.status,
    })
    notices.value = data.items
    total.value = data.total
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  await formRef.value?.validate()
  saving.value = true
  try {
    if (formMode.value === 'create') {
      await createNotice({
        title: formModel.title,
        content: formModel.content,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      message.success('公告创建成功')
    } else {
      await updateNotice(formModel.id, {
        title: formModel.title,
        content: formModel.content,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      message.success('公告更新成功')
    }

    formVisible.value = false
    await loadNotices()
  } finally {
    saving.value = false
  }
}

async function handleToggleStatus(row: NoticeItem, status: NoticeStatus) {
  await updateNoticeStatus(row.id, { status })
  message.success('公告状态已更新')
  await loadNotices()
}

onMounted(() => {
  void loadNotices()
})
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-4 overflow-hidden">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">公告管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">管理系统公告，支持按标题搜索和状态筛选。</p>
        </div>

        <NButton v-if="canUse('system:notice:create')" type="primary" @click="openCreate">
          + 新增公告
        </NButton>
      </div>

      <NCard :bordered="false" class="rounded-lg">
        <NSpace align="center" :wrap="true">
          <NInput
            v-model:value="query.keyword"
            clearable
            placeholder="公告标题"
            class="w-56"
            @keyup.enter="handleSearch"
          />
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
          <span class="text-sm text-[#6B7280]">共 {{ total }} 条</span>
          <NButton text type="primary" @click="loadNotices">刷新</NButton>
        </div>

        <NDataTable
          remote
          class="notice-table h-full"
          style="height: calc(100% - 105px)"
          :columns="columns"
          :data="notices"
          :loading="loading"
          :pagination="false"
          :row-key="(row: NoticeItem) => row.id"
          :bordered="false"
          flex-height
        />

        <div
          class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]"
        >
          <span>共 {{ total }} 条</span>
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

    <NModal
      v-model:show="formVisible"
      preset="card"
      :closable="false"
      class="compact-notice-modal"
      style="width: 600px; max-width: calc(100vw - 32px)"
    >
      <template #header>
        <div class="modal-header modal-header--hero">
          <h2 class="modal-header__title">
            {{ formMode === 'create' ? '新增公告' : '编辑公告' }}
          </h2>
          <p class="modal-header__hero-title">
            {{
              formMode === 'create'
                ? '填写公告标题和内容，保存后可立即展示'
                : '修改公告标题和内容，状态变更即时生效'
            }}
          </p>
          <button type="button" class="modal-close" @click="formVisible = false">
            <NIcon :size="18">
              <CloseOutline />
            </NIcon>
          </button>
        </div>
      </template>

      <div class="notice-modal-shell">
        <NForm
          ref="formRef"
          class="compact-notice-form"
          :model="formModel"
          :rules="rules"
          label-placement="left"
          label-width="76"
        >
          <section class="form-section form-section--primary">
            <div class="form-section__head">
              <h3>公告信息</h3>
              <p>标题不超过 128 个字符，内容支持任意文本。</p>
            </div>

            <div class="form-section-grid">
              <NFormItem label="标题" path="title">
                <NInput v-model:value="formModel.title" placeholder="公告标题" />
              </NFormItem>

              <NFormItem label="排序">
                <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
              </NFormItem>
            </div>
          </section>

          <section class="form-section form-section--muted">
            <div class="form-section__head">
              <h3>公告内容</h3>
            </div>

            <NFormItem label="内容" class="mb-0">
              <NInput
                v-model:value="formModel.content"
                type="textarea"
                :rows="4"
                placeholder="请输入公告内容"
              />
            </NFormItem>
          </section>

          <section class="form-section form-section--muted">
            <div class="form-section-grid">
              <NFormItem label="状态">
                <NSelect v-model:value="formModel.status" :options="statusFormOptions" />
              </NFormItem>

              <NFormItem label="备注">
                <NInput v-model:value="formModel.remark" placeholder="可选" />
              </NFormItem>
            </div>
          </section>
        </NForm>
      </div>

      <template #footer>
        <div class="modal-footer-actions">
          <NButton quaternary class="modal-footer-button" @click="formVisible = false">
            取消
          </NButton>
          <NButton
            type="primary"
            class="modal-footer-button modal-footer-button--primary"
            :loading="saving"
            @click="handleSubmit"
          >
            保存
          </NButton>
        </div>
      </template>
    </NModal>
  </main>
</template>

<style scoped>
.notice-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #4B5563;
  background: #F9FAFB;
  font-size: 13px;
}

.notice-table :deep(.n-data-table-td) {
  color: #374151;
  font-size: 14px;
  padding: 10px 16px;
}

.notice-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: unset !important;
}

.notice-table :deep(.n-data-table-tr) {
  transition: none;
}

.notice-table :deep(.n-data-table-tr:hover) {
  filter: brightness(0.97);
}

.compact-notice-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 32px;
  border: 1px solid #dfe9f5;
  background: #ffffff;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.compact-notice-modal :deep(.n-card-header) {
  padding: 0;
  border-bottom: 1px solid #dfe9f5;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.compact-notice-modal :deep(.n-card-header__main) {
  font-size: 19px;
  font-weight: 600;
  letter-spacing: 0.01em;
  color: #111827;
}

.compact-notice-modal :deep(.n-card__content) {
  padding: 20px 28px 10px;
}

.compact-notice-modal :deep(.n-card__footer) {
  padding: 16px 28px 24px;
  border-top: 1px solid #edf2f7;
  background: rgba(248, 250, 252, 0.85);
}

.compact-notice-form :deep(.n-form-item) {
  margin-bottom: 16px;
}

.compact-notice-form :deep(.n-form-item-label) {
  white-space: nowrap;
  align-items: center;
  padding-right: 14px;
  font-weight: 600;
  color: #374151;
}

.compact-notice-form :deep(.n-form-item-blank) {
  min-height: 40px;
}

.compact-notice-form :deep(.n-input-wrapper) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-notice-form :deep(.n-base-selection) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-notice-form :deep(.n-input),
.compact-notice-form :deep(.n-base-selection) {
  box-shadow: none;
}

.compact-notice-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.notice-modal-shell {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.modal-header--hero {
  position: relative;
  overflow: hidden;
  min-height: 120px;
  padding: 26px 28px 22px;
  background:
    radial-gradient(circle at top right, rgba(34, 197, 94, 0.12), transparent 24%),
    linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.modal-header__title {
  position: relative;
  z-index: 1;
  font-size: 19px;
  font-weight: 600;
  line-height: 1.3;
  color: #111827;
}

.modal-header__hero-title {
  position: relative;
  z-index: 1;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.6;
  color: #0f172a;
}

.modal-close {
  position: absolute;
  top: 20px;
  right: 22px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border: none;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.76);
  color: #64748b;
  box-shadow: 0 10px 24px rgba(148, 163, 184, 0.12);
  backdrop-filter: blur(8px);
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.modal-close:hover {
  background: #ffffff;
  color: #111827;
  box-shadow: 0 14px 28px rgba(148, 163, 184, 0.18);
  transform: translateY(-1px);
}

.form-section {
  border: 1px solid #e9eff6;
  border-radius: 14px;
  background: #ffffff;
  padding: 18px 18px 4px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.9);
}

.form-section--primary {
  border-color: #d9e7f8;
  background: linear-gradient(180deg, #ffffff 0%, #fcfdff 100%);
}

.form-section--muted {
  background: linear-gradient(180deg, #fcfdff 0%, #f9fbff 100%);
}

.form-section__head {
  margin-bottom: 12px;
}

.form-section__head h3 {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.form-section__head p {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: #6b7280;
}

.form-section-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  column-gap: 20px;
}

.modal-footer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.modal-footer-button {
  min-width: 92px;
  height: 40px;
  border-radius: 10px;
}

.modal-footer-button--primary {
  box-shadow: 0 10px 24px rgba(34, 197, 94, 0.18);
}

.mb-0 {
  margin-bottom: 0;
}

@media (max-width: 720px) {
  .form-section-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .compact-notice-modal :deep(.n-card-header),
  .compact-notice-modal :deep(.n-card__content),
  .compact-notice-modal :deep(.n-card__footer) {
    padding-left: 20px;
    padding-right: 20px;
  }

  .compact-notice-modal :deep(.n-card-header) {
    padding-bottom: 0;
  }

  .modal-header--hero {
    padding: 22px 20px 18px;
    min-height: 110px;
  }

  .modal-close {
    top: 18px;
    right: 18px;
  }
}
</style>
