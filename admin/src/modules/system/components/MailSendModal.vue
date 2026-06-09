<script setup lang="ts">
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { NButton, NCheckbox, NForm, NFormItem, NInput, NModal, NSelect } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
import type { MailSendFormModel } from '../types/mail-page'

defineProps<{
  accountOptions: SelectOption[]
  rules: FormRules
  saving: boolean
  show: boolean
  templateOptions: SelectOption[]
}>()

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()

const formRef = defineModel<FormInst | null>('formRef')
const formModel = defineModel<MailSendFormModel>('model', { required: true })
</script>

<template>
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="ez-modal-width-lg"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        title="发送邮件"
        subtitle="可选择模板发送，也可直接填写主题和正文。"
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
        label-width="86"
      >
        <section class="ez-modal-section ez-modal-section--soft">
          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="邮箱账号">
              <NSelect
                v-model:value="formModel.account_id"
                clearable
                filterable
                placeholder="默认邮箱"
                :options="accountOptions"
              />
            </NFormItem>

            <NFormItem label="邮件模板">
              <NSelect
                v-model:value="formModel.template_code"
                clearable
                filterable
                placeholder="不使用模板"
                :options="templateOptions"
              />
            </NFormItem>
          </div>

          <NFormItem label="收件人" path="to_text">
            <NInput
              v-model:value="formModel.to_text"
              type="textarea"
              :rows="2"
              placeholder="多个邮箱用逗号、分号或换行分隔"
            />
          </NFormItem>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="抄送">
              <NInput v-model:value="formModel.cc_text" placeholder="可选" />
            </NFormItem>

            <NFormItem label="密送">
              <NInput v-model:value="formModel.bcc_text" placeholder="可选" />
            </NFormItem>
          </div>
        </section>

        <section class="ez-modal-section">
          <NFormItem label="主题">
            <NInput
              v-model:value="formModel.subject"
              placeholder="使用模板时可留空"
              :disabled="!!formModel.template_code"
            />
          </NFormItem>

          <NFormItem label="正文">
            <NInput
              v-model:value="formModel.content"
              type="textarea"
              :rows="5"
              placeholder="使用模板时可留空"
              :disabled="!!formModel.template_code"
            />
          </NFormItem>

          <NFormItem label="变量">
            <NInput
              v-model:value="formModel.variables_text"
              type="textarea"
              :rows="4"
              placeholder="每行一个变量，例如 username=张三"
            />
          </NFormItem>

          <NFormItem label="HTML">
            <NCheckbox v-model:checked="formModel.is_html" :disabled="!!formModel.template_code">
              按 HTML 邮件发送
            </NCheckbox>
          </NFormItem>
        </section>
      </NForm>
    </div>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton quaternary class="min-w-[92px]" @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" class="min-w-[92px]" :loading="saving" @click="$emit('submit')">
          发送
        </NButton>
      </div>
    </template>
  </NModal>
</template>
