<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NIcon, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

import { STATUS_FORM_OPTIONS } from '@/constants/status'
import type { DictTypeFormModel } from '../composables/useDictPage'

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
const formModel = defineModel<DictTypeFormModel>('model', { required: true })
</script>

<template>
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="compact-dict-modal"
    style="width: 620px; max-width: calc(100vw - 32px)"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <div class="modal-header modal-header--hero">
        <h2 class="modal-header__title">{{ formMode === 'create' ? '新增字典类型' : '编辑字典类型' }}</h2>
        <p class="modal-header__hero-title">
          {{ formMode === 'create' ? '先定义稳定的字典编码，再让字典项围绕它展开。' : '修改字典展示信息，不影响已有字典项的主键归属。' }}
        </p>
        <button type="button" class="modal-close" @click="$emit('update:show', false)">
          <NIcon :size="18">
            <CloseOutline />
          </NIcon>
        </button>
      </div>
    </template>

    <div class="dict-modal-shell">
      <NForm ref="formRef" class="compact-dict-form" :model="formModel" :rules="rules" label-placement="left" label-width="76">
        <section class="form-section form-section--primary">
          <div class="form-section__head">
            <h3>基础信息</h3>
            <p>字典编码建议使用小写字母、数字、冒号、短横线和下划线，便于在前后端直接复用。</p>
          </div>

          <div class="form-section-grid">
            <NFormItem label="字典编码" path="code">
              <NInput v-model:value="formModel.code" placeholder="例如 common:status" :disabled="formMode === 'edit'" />
            </NFormItem>

            <NFormItem label="字典名称" path="name">
              <NInput v-model:value="formModel.name" placeholder="例如 通用状态" />
            </NFormItem>

            <NFormItem label="排序">
              <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
            </NFormItem>

            <NFormItem label="状态">
              <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
            </NFormItem>
          </div>
        </section>

        <section class="form-section form-section--muted">
          <NFormItem label="备注" class="mb-0">
            <NInput v-model:value="formModel.remark" type="textarea" :rows="3" placeholder="补充这个字典的适用场景或业务备注" />
          </NFormItem>
        </section>
      </NForm>
    </div>

    <template #footer>
      <div class="modal-footer-actions">
        <NButton quaternary class="modal-footer-button" @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" class="modal-footer-button modal-footer-button--primary" :loading="saving" @click="$emit('submit')">
          保存
        </NButton>
      </div>
    </template>
  </NModal>
</template>

<style scoped>
.compact-dict-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 32px;
  border: 1px solid #dfe9f5;
  background: #ffffff;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.compact-dict-modal :deep(.n-card-header) {
  padding: 0;
  border-bottom: 1px solid #dfe9f5;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.compact-dict-modal :deep(.n-card__content) {
  padding: 20px 28px 10px;
}

.compact-dict-modal :deep(.n-card__footer) {
  padding: 16px 28px 24px;
  border-top: 1px solid #edf2f7;
  background: rgba(248, 250, 252, 0.85);
}

.compact-dict-form :deep(.n-form-item) {
  margin-bottom: 16px;
}

.compact-dict-form :deep(.n-form-item-label) {
  white-space: nowrap;
  align-items: center;
  padding-right: 14px;
  font-weight: 600;
  color: #374151;
}

.compact-dict-form :deep(.n-form-item-blank) {
  min-height: 40px;
}

.compact-dict-form :deep(.n-input-wrapper),
.compact-dict-form :deep(.n-base-selection) {
  border-radius: 10px;
  background: #fbfcfe;
  box-shadow: none;
}

.dict-modal-shell {
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
  min-height: 124px;
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

  .compact-dict-modal :deep(.n-card__content),
  .compact-dict-modal :deep(.n-card__footer) {
    padding-left: 20px;
    padding-right: 20px;
  }

  .modal-header--hero {
    min-height: 112px;
    padding: 22px 20px 18px;
  }

  .modal-close {
    top: 18px;
    right: 18px;
  }
}
</style>
