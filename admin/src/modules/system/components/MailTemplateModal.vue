<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCheckbox, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
import { STATUS_FORM_OPTIONS } from '@/constants/status'
import type { MailTemplateFormModel } from '../types/mail-page'

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
const formModel = defineModel<MailTemplateFormModel>('model', { required: true })
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
        :title="formMode === 'create' ? '新增邮件模板' : '编辑邮件模板'"
        :subtitle="
          formMode === 'create'
            ? '定义主题、正文和可替换变量。'
            : '修改模板后，后续邮件会使用新内容。'
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
        label-width="86"
      >
        <section class="ez-modal-section ez-modal-section--soft">
          <div class="ez-modal-section__head">
            <h3>模板信息</h3>
            <p>可在主题和正文中使用 &#123;&#123; variable &#125;&#125; 形式的变量占位。</p>
          </div>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="编码" path="code">
              <NInput
                v-model:value="formModel.code"
                placeholder="例如 account:reset-password"
                :disabled="formMode === 'edit'"
              />
            </NFormItem>

            <NFormItem label="名称" path="name">
              <NInput v-model:value="formModel.name" placeholder="例如 重置密码邮件" />
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
          <div class="ez-modal-section__head">
            <h3>邮件内容</h3>
          </div>

          <NFormItem label="主题" path="subject">
            <NInput v-model:value="formModel.subject" placeholder="例如 {{username}}，请重置密码" />
          </NFormItem>

          <NFormItem label="正文" path="content">
            <NInput
              v-model:value="formModel.content"
              type="textarea"
              :rows="6"
              placeholder="请输入邮件正文模板"
            />
          </NFormItem>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="HTML">
              <NCheckbox v-model:checked="formModel.is_html">按 HTML 邮件发送</NCheckbox>
            </NFormItem>

            <NFormItem label="变量">
              <NInput
                v-model:value="formModel.variables_text"
                placeholder="例如 username, reset_link"
              />
            </NFormItem>
          </div>

          <NFormItem label="备注">
            <NInput v-model:value="formModel.remark" placeholder="可选" />
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
