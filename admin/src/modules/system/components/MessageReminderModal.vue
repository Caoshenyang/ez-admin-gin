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
import { computed } from 'vue'

import FormModalHeader from '@/components/FormModalHeader.vue'
import { STATUS_FORM_OPTIONS } from '@/constants/status'
import {
  MESSAGE_RECEIVER_TYPE_OPTIONS,
  receiverValuePlaceholder,
} from '../composables/message-page.utils'
import { MessageReceiverType } from '../types/message'
import type { MessageReminderFormModel } from '../types/message-page'

const props = defineProps<{
  formMode: 'create' | 'edit'
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
const formModel = defineModel<MessageReminderFormModel>('model', { required: true })

const receiverValuesDisabled = computed(() => {
  return (
    formModel.value.receiver_type === MessageReceiverType.Initiator ||
    formModel.value.receiver_type === MessageReceiverType.Assignee
  )
})

const receiverValuesPlaceholder = computed(() =>
  receiverValuePlaceholder(formModel.value.receiver_type),
)

const hasTemplateOptions = computed(() => props.templateOptions.length > 0)
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
        :title="formMode === 'create' ? '新增提醒规则' : '编辑提醒规则'"
        :subtitle="
          formMode === 'create'
            ? '配置触发事件、消息模板和接收人范围。'
            : '调整提醒规则后，后续触发会按新配置执行。'
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
            <h3>规则信息</h3>
            <p>触发事件使用稳定编码，例如 auth:login:success。</p>
          </div>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="编码" path="code">
              <NInput
                v-model:value="formModel.code"
                placeholder="例如 auth:login-success:admin"
                :disabled="formMode === 'edit'"
              />
            </NFormItem>

            <NFormItem label="名称" path="name">
              <NInput v-model:value="formModel.name" placeholder="例如 登录成功通知管理员" />
            </NFormItem>

            <NFormItem label="触发事件" path="trigger_event">
              <NInput v-model:value="formModel.trigger_event" placeholder="例如 auth:login:success" />
            </NFormItem>

            <NFormItem label="消息模板" path="template_id">
              <NSelect
                v-model:value="formModel.template_id"
                :options="templateOptions"
                :disabled="!hasTemplateOptions"
                filterable
                placeholder="请选择模板"
              />
            </NFormItem>
          </div>
        </section>

        <section class="ez-modal-section">
          <div class="ez-modal-section__head">
            <h3>提醒对象</h3>
          </div>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="渠道" path="channels">
              <NInput v-model:value="formModel.channels" placeholder="notification" />
            </NFormItem>

            <NFormItem label="接收人类型">
              <NSelect
                v-model:value="formModel.receiver_type"
                :options="MESSAGE_RECEIVER_TYPE_OPTIONS"
              />
            </NFormItem>
          </div>

          <NFormItem label="接收人">
            <NInput
              v-model:value="formModel.receiver_values"
              :disabled="receiverValuesDisabled"
              :placeholder="receiverValuesPlaceholder"
            />
          </NFormItem>
        </section>

        <section class="ez-modal-section">
          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="提前分钟">
              <NInputNumber
                v-model:value="formModel.advance_minutes"
                :min="0"
                :max="43200"
                class="w-full"
              />
            </NFormItem>

            <NFormItem label="排序">
              <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
            </NFormItem>

            <NFormItem label="跳转链接">
              <NInput v-model:value="formModel.link_url" placeholder="/audit/login-logs" />
            </NFormItem>

            <NFormItem label="状态">
              <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
            </NFormItem>
          </div>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="系统级">
              <NCheckbox v-model:checked="formModel.is_system">不可删除配置</NCheckbox>
            </NFormItem>

            <NFormItem label="备注">
              <NInput v-model:value="formModel.remark" placeholder="可选" />
            </NFormItem>
          </div>
        </section>
      </NForm>
    </div>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton quaternary class="min-w-[92px]" @click="$emit('update:show', false)">取消</NButton>
        <NButton
          type="primary"
          class="min-w-[92px]"
          :disabled="!hasTemplateOptions"
          :loading="saving"
          @click="$emit('submit')"
          >保存</NButton
        >
      </div>
    </template>
  </NModal>
</template>
