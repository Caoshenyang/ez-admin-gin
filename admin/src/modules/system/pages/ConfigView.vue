<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormRules } from 'naive-ui'
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
import { h, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { STATUS_FILTER_OPTIONS, STATUS_FORM_OPTIONS } from '@/constants/status'
import { createConfig, getConfigs, updateConfig, updateConfigStatus } from '../api/config'
import { ConfigStatus, type ConfigItem, type ConfigListQuery } from '../types/config'

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
const { canUse } = usePermission()
const successText = ref('')

const {
  items: configs,
  total,
  loading,
  query,
  load,
  handleSearch,
  handleReset,
  handlePageChange,
  handlePageSizeChange,
} = useRemotePagination<ConfigItem, ConfigListQuery>(getConfigs, {
  page: 1,
  page_size: 10,
  keyword: '',
  group_code: '',
  status: 0,
})

function defaultFormModel(): ConfigFormModel {
  return { id: 0, group_code: '', key: '', name: '', value: '', sort: 0, status: ConfigStatus.Enabled, remark: '' }
}

const { formRef, formVisible, formMode, formModel, saving, rules, openCreate, openEdit, handleSubmit } =
  useModalForm<ConfigFormModel>(defaultFormModel, {
    rules: {
      group_code: [{ required: true, message: '请输入配置分组', trigger: 'blur' }],
      key: [{ required: true, message: '请输入配置键', trigger: 'blur' }],
      name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
      value: [{ required: true, message: '请输入配置值', trigger: 'blur' }],
    } as FormRules,
  })

const { handleToggleStatus } = useStatusToggle<ConfigItem>(updateConfigStatus, {
  onSuccess: async () => {
    successText.value = '配置状态已更新'
    await load()
  },
})

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

async function onSubmit() {
  const payload = {
    group_code: formModel.group_code,
    name: formModel.name,
    value: formModel.value,
    sort: formModel.sort,
    status: formModel.status,
    remark: formModel.remark,
  }
  if (formMode.value === 'create') {
    await createConfig({ ...payload, key: formModel.key })
    successText.value = '配置创建成功'
    message.success('配置创建成功')
  } else {
    await updateConfig(formModel.id, payload)
    successText.value = '配置已更新'
    message.success('配置更新成功')
  }
  await load()
}
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
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
          <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" class="w-36" />
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
          <NButton text type="primary" @click="load">刷新</NButton>
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
                <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
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
            @click="handleSubmit(onSubmit)"
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
