<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
import { STATUS_FORM_OPTIONS } from '@/constants/status'
import type { AttachmentEditFormModel } from '../types/attachment-page'

defineProps<{
  rules: FormRules
  saving: boolean
  show: boolean
}>()

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()

const formRef = defineModel<FormInst | null>('formRef')
const formModel = defineModel<AttachmentEditFormModel>('model', { required: true })
</script>

<template>
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="ez-modal-width-lg"
    :bordered="false"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        title="编辑附件"
        subtitle="只调整附件的展示信息、分类和状态，不会替换底层文件内容。"
        @close="$emit('update:show', false)"
      />
    </template>

    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="top"
      class="ez-modal-form"
    >
      <NFormItem label="附件名称" path="display_name">
        <NInput v-model:value="formModel.display_name" placeholder="请输入附件名称" />
      </NFormItem>

      <div class="ez-form-grid ez-form-grid--2">
        <NFormItem label="附件分类" path="category">
          <NInput v-model:value="formModel.category" placeholder="附件分类" />
        </NFormItem>
        <NFormItem label="业务类型" path="biz_type">
          <NInput v-model:value="formModel.biz_type" placeholder="业务类型" />
        </NFormItem>
      </div>

      <NFormItem label="状态" path="status">
        <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
      </NFormItem>

      <NFormItem label="备注" path="remark">
        <NInput v-model:value="formModel.remark" type="textarea" :rows="3" placeholder="附件备注" />
      </NFormItem>

      <div class="ez-modal-footer">
        <NButton @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" :loading="saving" @click="$emit('submit')">保存</NButton>
      </div>
    </NForm>
  </NModal>
</template>
