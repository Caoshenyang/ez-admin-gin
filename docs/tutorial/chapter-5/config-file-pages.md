---
title: 配置与文件页面
description: "实现系统配置管理页面和文件上传管理页面。"
---

# 配置与文件页面

上一节已经把角色和菜单接成了真实页面。现在继续补齐系统管理剩下的两个功能页：配置管理和文件管理。

完成这一节后，侧边栏里的"配置管理"和"文件管理"不再停留在占位页。配置页面负责维护系统键值配置，按分组归类管理；文件页面负责上传附件、查看文件列表和复制文件链接。

::: tip 🎯 本节目标
这一节会把 `system/ConfigView` 和 `system/FileView` 从占位页换成真实页面，并补齐配置和文件相关的类型和 API 封装。配置页面采用搜索 + 数据表 + 弹框表单布局；文件页面使用上传按钮 + 数据表布局，支持按文件名和类型筛选。
:::

## 先看接口边界

配置管理接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/configs` | 配置分页列表 |
| `POST` | `/api/v1/system/configs` | 创建配置 |
| `POST` | `/api/v1/system/configs/:id/update` | 编辑配置 |
| `POST` | `/api/v1/system/configs/:id/status` | 修改配置状态 |

文件管理接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/files` | 文件分页列表 |
| `POST` | `/api/v1/system/files` | 上传文件 |

::: warning ⚠️ 配置键创建后不可更改
配置键（`key`）会被后端缓存到 Redis，也被其他模块引用。允许随意修改键名，容易导致缓存失效或引用断裂。所以编辑模式下，配置键是只读字段。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
└─ src/
   ├─ api/
   │  ├─ config.ts
   │  └─ file.ts
   ├─ pages/
   │  └─ system/
   │     ├─ ConfigView.vue
   │     └─ FileView.vue
   ├─ router/
   │  └─ dynamic-menu.ts
   └─ types/
      ├─ config.ts
      └─ file.ts
```

## 开始前先确认

开始之前，先确认下面几件事：

- 已完成上一节 [角色与菜单页面](./role-menu-pages)。
- 登录后侧边栏能看到"配置管理"和"文件管理"。
- 当前账号拥有配置与文件相关按钮权限。
- 后端 `/api/v1/system/configs` 和 `/api/v1/system/files` 可以正常返回数据。

## 🛠️ 完整代码

下面直接引入本节对应的完整项目文件，默认折叠。需要复制或对照时点击展开即可。

::: details `admin/src/types/config.ts` — 配置类型

```ts
export const ConfigStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type ConfigStatus = (typeof ConfigStatus)[keyof typeof ConfigStatus]

export interface ConfigItem {
  id: number
  group_code: string
  key: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
  created_at: string
  updated_at: string
}

export interface ConfigListQuery {
  page: number
  page_size: number
  keyword?: string
  group_code?: string
  status?: ConfigStatus | 0
}

export interface ConfigListResponse {
  items: ConfigItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateConfigPayload {
  group_code: string
  key: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
}

export interface UpdateConfigPayload {
  group_code: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
}

export interface UpdateConfigStatusPayload {
  status: ConfigStatus
}
```

:::

::: details `admin/src/types/file.ts` — 文件类型

```ts
export const FileStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type FileStatus = (typeof FileStatus)[keyof typeof FileStatus]

export interface FileItem {
  id: number
  storage: string
  original_name: string
  file_name: string
  ext: string
  mime_type: string
  size: number
  sha256: string
  path: string
  url: string
  uploader_id: number
  status: FileStatus
  remark: string
  created_at: string
  updated_at: string
}

export interface FileListQuery {
  page: number
  page_size: number
  keyword?: string
  ext?: string
  status?: FileStatus | 0
}

export interface FileListResponse {
  items: FileItem[]
  total: number
  page: number
  page_size: number
}
```

:::

::: details `admin/src/api/config.ts` — 配置接口

```ts
import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  ConfigItem,
  ConfigListQuery,
  ConfigListResponse,
  CreateConfigPayload,
  UpdateConfigPayload,
  UpdateConfigStatusPayload,
} from '../types/config'

export async function getConfigs(params: ConfigListQuery) {
  const response = await http.get<ApiResponse<ConfigListResponse>>('/system/configs', { params })
  return response.data.data
}

export async function createConfig(payload: CreateConfigPayload) {
  const response = await http.post<ApiResponse<ConfigItem>>('/system/configs', payload)
  return response.data.data
}

export async function updateConfig(id: number, payload: UpdateConfigPayload) {
  const response = await http.post<ApiResponse<ConfigItem>>(`/system/configs/${id}/update`, payload)
  return response.data.data
}

export async function updateConfigStatus(id: number, payload: UpdateConfigStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/configs/${id}/status`,
    payload,
  )
  return response.data.data
}
```

:::

::: details `admin/src/api/file.ts` — 文件接口

```ts
import http from './http'

import type { ApiResponse } from '../types/http'
import type { FileItem, FileListQuery, FileListResponse } from '../types/file'

export async function getFiles(params: FileListQuery) {
  const response = await http.get<ApiResponse<FileListResponse>>('/system/files', { params })
  return response.data.data
}

export async function uploadFile(formData: FormData) {
  const response = await http.post<ApiResponse<FileItem>>('/system/files', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return response.data.data
}
```

:::

::: details `admin/src/pages/system/ConfigView.vue` — 配置管理页面

```vue
<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import {
  NAlert,
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

import { createConfig, getConfigs, updateConfig, updateConfigStatus } from '../../api/config'
import { buttonPermissionCodes } from '../../router/dynamic-menu'
import {
  ConfigStatus,
  type ConfigItem,
  type ConfigListQuery,
} from '../../types/config'

interface ConfigFormModel {
  id: number
  group_code: string
  key: string
  name: string
  value: string
  sort: number
  status: ConfigStatus
  remark: string
}

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const configs = ref<ConfigItem[]>([])
const total = ref(0)
const successText = ref('')

const query = reactive<ConfigListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  group_code: '',
  status: 0,
})

const formRef = ref<FormInst | null>(null)
const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const formModel = reactive<ConfigFormModel>({
  id: 0,
  group_code: '',
  key: '',
  name: '',
  value: '',
  sort: 0,
  status: ConfigStatus.Enabled,
  remark: '',
})

const statusFilterOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: ConfigStatus.Enabled },
  { label: '禁用', value: ConfigStatus.Disabled },
]

const statusFormOptions = [
  { label: '启用', value: ConfigStatus.Enabled },
  { label: '禁用', value: ConfigStatus.Disabled },
]

const rules: FormRules = {
  group_code: [{ required: true, message: '请输入配置分组', trigger: 'blur' }],
  key: [{ required: true, message: '请输入配置键', trigger: 'blur' }],
  name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
  value: [{ required: true, message: '请输入配置值', trigger: 'blur' }],
}

const columns: DataTableColumns<ConfigItem> = [
  {
    title: '分组',
    key: 'group_code',
    width: 140,
    render(row) {
      return h(NTag, { size: 'small', bordered: false, type: 'info' }, { default: () => row.group_code })
    },
  },
  {
    title: '键',
    key: 'key',
    width: 200,
    ellipsis: { tooltip: true },
  },
  {
    title: '名称',
    key: 'name',
    width: 160,
  },
  {
    title: '值',
    key: 'value',
    minWidth: 180,
    ellipsis: { tooltip: true },
  },
  {
    title: '排序',
    key: 'sort',
    width: 80,
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return h(
        NTag,
        { type: row.status === ConfigStatus.Enabled ? 'success' : 'error', bordered: false },
        { default: () => (row.status === ConfigStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render(row) {
      const nextStatus = row.status === ConfigStatus.Enabled ? ConfigStatus.Disabled : ConfigStatus.Enabled

      return h(
        NSpace,
        { size: 8, align: 'center' },
        {
          default: () =>
            [
              canUse('system:config:update')
                ? h(
                    NButton,
                    { size: 'small', ghost: true, type: 'info', onClick: () => openEdit(row) },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:config:status')
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
                            type: nextStatus === ConfigStatus.Disabled ? 'error' : 'success',
                          },
                          { default: () => (nextStatus === ConfigStatus.Disabled ? '禁用' : '启用') },
                        ),
                      default: () => `确认${nextStatus === ConfigStatus.Disabled ? '禁用' : '启用'}该配置？`,
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
  return value ? new Date(value).toLocaleString() : '-'
}

function resetForm() {
  Object.assign(formModel, {
    id: 0,
    group_code: '',
    key: '',
    name: '',
    value: '',
    sort: 0,
    status: ConfigStatus.Enabled,
    remark: '',
  })
}

function handleSearch() {
  query.page = 1
  void loadConfigs()
}

function handleReset() {
  query.page = 1
  query.page_size = 10
  query.keyword = ''
  query.group_code = ''
  query.status = 0
  void loadConfigs()
}

function handlePageChange(page: number) {
  query.page = page
  void loadConfigs()
}

function handlePageSizeChange(pageSize: number) {
  query.page = 1
  query.page_size = pageSize
  void loadConfigs()
}

function openCreate() {
  formMode.value = 'create'
  resetForm()
  formVisible.value = true
}

function openEdit(row: ConfigItem) {
  formMode.value = 'edit'
  Object.assign(formModel, {
    id: row.id,
    group_code: row.group_code,
    key: row.key,
    name: row.name,
    value: row.value,
    sort: row.sort,
    status: row.status,
    remark: row.remark,
  })
  formVisible.value = true
}

async function loadConfigs() {
  loading.value = true
  try {
    const data = await getConfigs({
      ...query,
      keyword: query.keyword?.trim() || undefined,
      group_code: query.group_code?.trim() || undefined,
      status: query.status === 0 ? undefined : query.status,
    })
    configs.value = data.items
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
      await createConfig({
        group_code: formModel.group_code,
        key: formModel.key,
        name: formModel.name,
        value: formModel.value,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      successText.value = '配置创建成功'
      message.success('配置创建成功')
    } else {
      await updateConfig(formModel.id, {
        group_code: formModel.group_code,
        name: formModel.name,
        value: formModel.value,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      successText.value = '配置已更新'
      message.success('配置更新成功')
    }

    formVisible.value = false
    await loadConfigs()
  } finally {
    saving.value = false
  }
}

async function handleToggleStatus(row: ConfigItem, status: ConfigStatus) {
  await updateConfigStatus(row.id, { status })
  successText.value = `配置已${status === ConfigStatus.Enabled ? '启用' : '禁用'}`
  message.success('配置状态已更新')
  await loadConfigs()
}

onMounted(() => {
  void loadConfigs()
})
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-4 overflow-hidden">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">配置管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">维护系统键值配置，按分组归类管理。</p>
        </div>

        <NButton v-if="canUse('system:config:create')" type="primary" @click="openCreate">
          + 新增配置
        </NButton>
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
            placeholder="键 / 名称"
            class="w-56"
            @keyup.enter="handleSearch"
          />
          <NInput
            v-model:value="query.group_code"
            clearable
            placeholder="分组"
            class="w-44"
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
          <NButton text type="primary" @click="loadConfigs">刷新</NButton>
        </div>

        <NDataTable
          remote
          class="config-table h-full"
          style="height: calc(100% - 105px)"
          :columns="columns"
          :data="configs"
          :loading="loading"
          :pagination="false"
          :row-key="(row: ConfigItem) => row.id"
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
      class="compact-config-modal"
      style="width: 600px; max-width: calc(100vw - 32px)"
    >
      <template #header>
        <div class="modal-header modal-header--hero">
          <h2 class="modal-header__title">
            {{ formMode === 'create' ? '新增配置' : '编辑配置' }}
          </h2>
          <p class="modal-header__hero-title">
            {{
              formMode === 'create'
                ? '填写配置分组、键和值，保存后立即生效'
                : '修改配置名称和值，键创建后不可更改'
            }}
          </p>
          <button type="button" class="modal-close" @click="formVisible = false">
            <NIcon :size="18">
              <CloseOutline />
            </NIcon>
          </button>
        </div>
      </template>

      <div class="config-modal-shell">
        <NForm
          ref="formRef"
          class="compact-config-form"
          :model="formModel"
          :rules="rules"
          label-placement="left"
          label-width="76"
        >
          <section class="form-section form-section--primary">
            <div class="form-section__head">
              <h3>基础信息</h3>
              <p>配置键只允许小写字母、数字、冒号、短横线和下划线。</p>
            </div>

            <div class="form-section-grid">
              <NFormItem label="分组" path="group_code">
                <NInput v-model:value="formModel.group_code" placeholder="例如 site" />
              </NFormItem>

              <NFormItem label="键" path="key">
                <NInput
                  v-model:value="formModel.key"
                  placeholder="例如 site_name"
                  :disabled="formMode === 'edit'"
                />
              </NFormItem>

              <NFormItem label="名称" path="name">
                <NInput v-model:value="formModel.name" placeholder="例如 站点名称" />
              </NFormItem>

              <NFormItem label="排序">
                <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
              </NFormItem>
            </div>
          </section>

          <section class="form-section form-section--muted">
            <div class="form-section__head">
              <h3>配置值</h3>
              <p>配置值支持任意文本，启用后会被缓存到 Redis。</p>
            </div>

            <NFormItem label="值" path="value" class="mb-0">
              <NInput
                v-model:value="formModel.value"
                type="textarea"
                :rows="3"
                placeholder="请输入配置值"
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
.config-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #374151;
  background: #fff;
}

.config-table :deep(.n-data-table-td) {
  color: #374151;
}

.config-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: #f8fbff;
}

.compact-config-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 32px;
  border: 1px solid #dfe9f5;
  background: #ffffff;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.compact-config-modal :deep(.n-card-header) {
  padding: 0;
  border-bottom: 1px solid #dfe9f5;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.compact-config-modal :deep(.n-card-header__main) {
  font-size: 19px;
  font-weight: 600;
  letter-spacing: 0.01em;
  color: #111827;
}

.compact-config-modal :deep(.n-card__content) {
  padding: 20px 28px 10px;
}

.compact-config-modal :deep(.n-card__footer) {
  padding: 16px 28px 24px;
  border-top: 1px solid #edf2f7;
  background: rgba(248, 250, 252, 0.85);
}

.compact-config-form :deep(.n-form-item) {
  margin-bottom: 16px;
}

.compact-config-form :deep(.n-form-item-label) {
  white-space: nowrap;
  align-items: center;
  padding-right: 14px;
  font-weight: 600;
  color: #374151;
}

.compact-config-form :deep(.n-form-item-blank) {
  min-height: 40px;
}

.compact-config-form :deep(.n-input-wrapper) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-config-form :deep(.n-base-selection) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-config-form :deep(.n-input),
.compact-config-form :deep(.n-base-selection) {
  box-shadow: none;
}

.compact-config-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.config-modal-shell {
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

  .compact-config-modal :deep(.n-card-header),
  .compact-config-modal :deep(.n-card__content),
  .compact-config-modal :deep(.n-card__footer) {
    padding-left: 20px;
    padding-right: 20px;
  }

  .compact-config-modal :deep(.n-card-header) {
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
```

:::

::: details `admin/src/pages/system/FileView.vue` — 文件管理页面

```vue
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
```

:::

::: details `admin/src/router/dynamic-menu.ts` — 动态路由映射

修改后，`system/ConfigView` 和 `system/FileView` 会从占位页切换为真实页面。

```ts
import {
  AlbumsOutline,
  AppsOutline,
  BeakerOutline,
  BuildOutline,
  CogOutline,
  DocumentTextOutline,
  FolderOpenOutline,
  GridOutline,
  LayersOutline,
  ListOutline,
  NotificationsOutline,
  PeopleOutline,
  PulseOutline,
  ServerOutline,
  SettingsOutline,
  ShieldCheckmarkOutline,
  TimeOutline,
} from '@vicons/ionicons5'
import { NIcon, type MenuOption } from 'naive-ui'
import type { RouteRecordRaw } from 'vue-router'
import { computed, h, shallowRef, type Component } from 'vue'

import { MenuType, type AuthMenu } from '../types/menu'

type RouteComponent = NonNullable<RouteRecordRaw['component']>
type MenuIconComponent = Component

const placeholderPage = () => import('../pages/system/PlaceholderPage.vue')

const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': () => import('../pages/system/HealthView.vue'),
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': () => import('../pages/system/RoleView.vue'),
  'system/MenuView': () => import('../pages/system/MenuView.vue'),
  'system/ConfigView': () => import('../pages/system/ConfigView.vue'),
  'system/FileView': () => import('../pages/system/FileView.vue'),
  'system/OperationLogView': () => import('../pages/system/OperationLogView.vue'),
  'system/LoginLogView': () => import('../pages/system/LoginLogView.vue'),
  'system/NoticeView': () => import('../pages/system/NoticeView.vue'),
}

const defaultMenuIcon = AppsOutline

// 后端 icon 字段只允许命中这份前端白名单，避免把任意字符串直接当组件渲染。
const menuIconMap: Record<string, MenuIconComponent> = {
  albums: AlbumsOutline,
  app: AppsOutline,
  apps: AppsOutline,
  beaker: BeakerOutline,
  blog: DocumentTextOutline,
  build: BuildOutline,
  cog: CogOutline,
  config: BuildOutline,
  dashboard: GridOutline,
  directory: AlbumsOutline,
  document: DocumentTextOutline,
  edit: DocumentTextOutline,
  experiment: BeakerOutline,
  file: FolderOpenOutline,
  files: FolderOpenOutline,
  folder: FolderOpenOutline,
  grid: GridOutline,
  health: PulseOutline,
  history: TimeOutline,
  home: GridOutline,
  layout: GridOutline,
  layoutdashboard: GridOutline,
  layers: LayersOutline,
  list: ListOutline,
  log: ListOutline,
  loginlog: TimeOutline,
  loginlogs: TimeOutline,
  logs: ListOutline,
  menu: LayersOutline,
  menus: LayersOutline,
  monitor: PulseOutline,
  notice: NotificationsOutline,
  notices: NotificationsOutline,
  notification: NotificationsOutline,
  notifications: NotificationsOutline,
  operationlog: ListOutline,
  operationlogs: ListOutline,
  page: DocumentTextOutline,
  people: PeopleOutline,
  person: PeopleOutline,
  role: ShieldCheckmarkOutline,
  roles: ShieldCheckmarkOutline,
  server: ServerOutline,
  setting: SettingsOutline,
  settings: SettingsOutline,
  shield: ShieldCheckmarkOutline,
  system: SettingsOutline,
  time: TimeOutline,
  user: PeopleOutline,
  users: PeopleOutline,
}

const builtinMenuOptions: MenuOption[] = [
  {
    label: '工作台',
    key: '/dashboard',
    icon: renderMenuIcon(GridOutline),
  },
]

export const authMenus = shallowRef<AuthMenu[]>([])

export const sideMenuOptions = computed<MenuOption[]>(() => {
  return [...builtinMenuOptions, ...buildMenuOptions(authMenus.value)]
})

export const buttonPermissionCodes = computed(() => {
  return collectButtonCodes(authMenus.value)
})

export function setAuthMenus(menus: AuthMenu[]) {
  authMenus.value = menus
}

export function clearAuthMenus() {
  authMenus.value = []
}

export function buildDynamicRoutes(menus: AuthMenu[]) {
  return collectPageMenus(menus).map<RouteRecordRaw>((menu) => ({
    path: toChildRoutePath(menu.path),
    name: `menu-${menu.id}`,
    component: resolveRouteComponent(menu.component),
    props: {
      title: menu.title,
      description: `${menu.title} 页面后续会接入真实业务。`,
    },
    meta: {
      title: menu.title,
      menuCode: menu.code,
    },
  }))
}

export function findMenuTitleByPath(path: string) {
  return collectPageMenus(authMenus.value).find((menu) => menu.path === path)?.title
}

function buildMenuOptions(menus: AuthMenu[]) {
  return menus.map(toMenuOption).filter(isMenuOption)
}

function toMenuOption(menu: AuthMenu): MenuOption | null {
  if (menu.type === MenuType.Button) {
    return null
  }

  const children = buildMenuOptions(menu.children ?? [])
  const key = menu.path || menu.code

  return {
    label: menu.title,
    key,
    icon: resolveMenuIcon(menu.icon),
    disabled: menu.type === MenuType.Directory && children.length === 0,
    children: children.length > 0 ? children : undefined,
  }
}

function isMenuOption(option: MenuOption | null): option is MenuOption {
  return option !== null
}

function collectPageMenus(menus: AuthMenu[]) {
  const result: AuthMenu[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Menu && menu.path) {
      result.push(menu)
    }

    result.push(...collectPageMenus(menu.children ?? []))
  }

  return result
}

function collectButtonCodes(menus: AuthMenu[]) {
  const result: string[] = []

  for (const menu of menus) {
    if (menu.type === MenuType.Button) {
      result.push(menu.code)
    }

    result.push(...collectButtonCodes(menu.children ?? []))
  }

  return result
}

function resolveRouteComponent(component: string) {
  return routeComponentMap[component] ?? placeholderPage
}

function resolveMenuIcon(icon: string) {
  return renderMenuIcon(menuIconMap[normalizeMenuIcon(icon)] ?? defaultMenuIcon)
}

function renderMenuIcon(icon: MenuIconComponent) {
  return () =>
    h(NIcon, null, {
      default: () => h(icon),
    })
}

function normalizeMenuIcon(icon: string) {
  return icon.trim().toLowerCase().replace(/[^a-z0-9]/g, '')
}

function toChildRoutePath(path: string) {
  return path.replace(/^\/+/, '')
}
```

:::

::: details 为什么配置编辑不允许改键
配置键会被后端缓存到 Redis，其他模块也可能通过 `GET /api/v1/system/configs/value/:key` 按键读取值。如果允许编辑键名，需要同步清理旧缓存、写入新缓存，还要更新所有引用方。后端当前的编辑接口也没有接收 `key` 字段，所以前端编辑表单会把键作为只读信息处理。
:::

## ✅ 验证结果

先启动后端和前端：

::: code-group

```bash [后端]
cd server
go run .
```

```bash [前端]
cd admin
pnpm dev
```

:::

然后按下面顺序验证：

1. 使用 `admin / EzAdmin@123456` 登录。
2. 点击"系统管理 / 配置管理"，确认配置列表能正常加载。
3. 点击"+ 新增配置"，填写分组、键、名称和值，保存后列表中出现新配置。
4. 点击编辑，确认键为只读，修改名称和值后保存成功。
5. 点击"禁用"按钮，确认状态切换成功。
6. 进入"系统管理 / 文件管理"，确认文件列表能正常加载。
7. 点击"上传文件"，选择一张图片或文档，确认上传成功后列表中出现新记录。
8. 点击复制链接按钮，确认剪贴板中包含文件 URL。

::: details 如果上传失败，先检查这几件事
- 文件大小是否超过后端配置的 `max_size_mb`（默认 10 MB）。
- 文件扩展名是否在 `allowed_exts` 列表中。
- `uploads` 目录是否有写入权限。
- 后端控制台是否有 `save file` 相关错误日志。
:::

## 本节小结

这一节把系统管理剩余的两个页面补齐了：

- 配置页面负责维护分组键值配置，支持搜索、筛选、新增、编辑和状态切换。
- 文件页面负责上传附件、查看文件列表和复制文件链接。
- 配置键创建后不可更改，避免缓存和引用断裂。
- 文件上传通过 `multipart/form-data` 提交，后端自动生成文件名并计算 SHA256 校验。

到这里，第 5 章前端管理台的所有基础页面都已完成。下一节继续补齐日志查询页面：[日志页面](./log-pages)。
