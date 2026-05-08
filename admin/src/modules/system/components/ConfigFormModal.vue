<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NIcon, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

import { STATUS_FORM_OPTIONS } from '@/constants/status'
import type { ConfigFormModel } from '../composables/config-page.utils'

defineProps<{
  formMode: 'create' | 'edit'
  rules: FormRules
  saving: boolean
  show: boolean
}>()

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()

const formRef = defineModel<FormInst | null>('formRef')
const formModel = defineModel<ConfigFormModel>('model', { required: true })
</script>

<template>
  <NModal :show="show" preset="card" :closable="false" class="compact-config-modal" style="width: 600px; max-width: calc(100vw - 32px)" @update:show="(value) => $emit('update:show', value)">
    <template #header>
      <div class="modal-header modal-header--hero">
        <h2 class="modal-header__title">{{ formMode === 'create' ? '新增配置' : '编辑配置' }}</h2>
        <p class="modal-header__hero-title">{{ formMode === 'create' ? '填写配置分组、键和值，保存后立即生效' : '修改配置名称和值，键创建后不可更改' }}</p>
        <button type="button" class="modal-close" @click="$emit('update:show', false)">
          <NIcon :size="18"><CloseOutline /></NIcon>
        </button>
      </div>
    </template>

    <div class="config-modal-shell">
      <NForm ref="formRef" class="compact-config-form" :model="formModel" :rules="rules" label-placement="left" label-width="76">
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
              <NInput v-model:value="formModel.key" placeholder="例如 site_name" :disabled="formMode === 'edit'" />
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
            <NInput v-model:value="formModel.value" type="textarea" :rows="3" placeholder="请输入配置值" />
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
        <NButton quaternary class="modal-footer-button" @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" class="modal-footer-button modal-footer-button--primary" :loading="saving" @click="$emit('submit')">保存</NButton>
      </div>
    </template>
  </NModal>
</template>

<style scoped>
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

.compact-config-form :deep(.n-input-wrapper),
.compact-config-form :deep(.n-base-selection) {
  border-radius: 10px;
  background: #fbfcfe;
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
}

.form-section {
  border: 1px solid #e9eff6;
  border-radius: 14px;
  background: #ffffff;
  padding: 18px 18px 4px;
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
</style>
