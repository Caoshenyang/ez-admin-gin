<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

import { STATUS_FORM_OPTIONS } from '@/constants/status'
import FormModalHeader from '@/components/FormModalHeader.vue'
import type { DictTypeFormModel } from '../types/dict-page'

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
    class="ez-modal-width-md"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增字典类型' : '编辑字典类型'"
        :subtitle="
          formMode === 'create'
            ? '先定义稳定的字典编码，再让字典项围绕它展开。'
            : '修改字典展示信息，不影响已有字典项的主键归属。'
        "
        @close="$emit('update:show', false)"
      />
    </template>

    <div class="ez-modal-shell">
      <NForm
        ref="formRef"
        class="ez-modal-form"
        :model="formModel"
        :rules="rules"
        label-placement="left"
        label-width="76"
      >
        <section class="ez-modal-section ez-modal-section--soft">
          <div class="ez-modal-section__head">
            <h3>基础信息</h3>
            <p>字典编码建议使用小写字母、数字、冒号、短横线和下划线，便于在前后端直接复用。</p>
          </div>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="字典编码" path="code">
              <NInput
                v-model:value="formModel.code"
                placeholder="例如 common:status"
                :disabled="formMode === 'edit'"
              />
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

        <section class="ez-modal-section">
          <NFormItem label="备注">
            <NInput
              v-model:value="formModel.remark"
              type="textarea"
              :rows="3"
              placeholder="补充这个字典的适用场景或业务备注"
            />
          </NFormItem>
        </section>
      </NForm>
    </div>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton quaternary class="min-w-[92px]" @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" class="min-w-[92px]" :loading="saving" @click="$emit('submit')">
          保存
        </NButton>
      </div>
    </template>
  </NModal>
</template>
