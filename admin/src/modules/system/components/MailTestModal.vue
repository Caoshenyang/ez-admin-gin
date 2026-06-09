<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
import type { MailTestFormModel } from '../types/mail-page'

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
const formModel = defineModel<MailTestFormModel>('model', { required: true })
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
        title="测试邮箱"
        subtitle="发送一封测试邮件并记录测试结果。"
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
        <NFormItem label="收件人" path="to_text">
          <NInput
            v-model:value="formModel.to_text"
            type="textarea"
            :rows="2"
            placeholder="多个邮箱用逗号、分号或换行分隔"
          />
        </NFormItem>

        <NFormItem label="主题">
          <NInput v-model:value="formModel.subject" placeholder="测试邮件主题" />
        </NFormItem>

        <NFormItem label="正文">
          <NInput v-model:value="formModel.content" type="textarea" :rows="4" />
        </NFormItem>
      </NForm>
    </div>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton quaternary class="min-w-[92px]" @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" class="min-w-[92px]" :loading="saving" @click="$emit('submit')">
          发送测试
        </NButton>
      </div>
    </template>
  </NModal>
</template>
