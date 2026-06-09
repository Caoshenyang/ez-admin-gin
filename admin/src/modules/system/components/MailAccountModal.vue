<script setup lang="ts">
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import {
  NButton,
  NCheckbox,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NSelect,
} from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
import { STATUS_FORM_OPTIONS } from '@/constants/status'
import type { MailAccountFormModel } from '../types/mail-page'

defineProps<{
  encryptionOptions: SelectOption[]
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
const formModel = defineModel<MailAccountFormModel>('model', { required: true })
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
        :title="formMode === 'create' ? '新增邮箱账号' : '编辑邮箱账号'"
        :subtitle="
          formMode === 'create'
            ? '配置 SMTP 服务器、发件人和默认状态。'
            : '密码留空表示沿用已有密码或授权码。'
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
        label-width="92"
      >
        <section class="ez-modal-section ez-modal-section--soft">
          <div class="ez-modal-section__head">
            <h3>SMTP 配置</h3>
            <p>支持无加密、SSL/TLS 和 STARTTLS 三种连接方式。</p>
          </div>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="邮箱名称" path="name">
              <NInput v-model:value="formModel.name" placeholder="例如 企业邮箱" />
            </NFormItem>

            <NFormItem label="SMTP 主机" path="host">
              <NInput v-model:value="formModel.host" placeholder="例如 smtp.example.com" />
            </NFormItem>

            <NFormItem label="端口">
              <NInputNumber v-model:value="formModel.port" :min="1" :max="65535" class="w-full" />
            </NFormItem>

            <NFormItem label="加密方式">
              <NSelect v-model:value="formModel.encryption" :options="encryptionOptions" />
            </NFormItem>

            <NFormItem label="用户名">
              <NInput v-model:value="formModel.username" placeholder="通常为邮箱地址" />
            </NFormItem>

            <NFormItem label="密码">
              <NInput
                v-model:value="formModel.password"
                type="password"
                show-password-on="click"
                placeholder="密码或授权码"
              />
            </NFormItem>
          </div>
        </section>

        <section class="ez-modal-section">
          <div class="ez-modal-section__head">
            <h3>发件信息</h3>
          </div>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="发件邮箱" path="from_email">
              <NInput v-model:value="formModel.from_email" placeholder="noreply@example.com" />
            </NFormItem>

            <NFormItem label="发件人名称">
              <NInput v-model:value="formModel.from_name" placeholder="例如 EZ Admin" />
            </NFormItem>

            <NFormItem label="状态">
              <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
            </NFormItem>

            <NFormItem label="默认账号">
              <NCheckbox v-model:checked="formModel.is_default">用于未指定账号的邮件</NCheckbox>
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
