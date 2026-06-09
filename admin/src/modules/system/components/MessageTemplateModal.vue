<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
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
import { MESSAGE_TEMPLATE_TYPE_OPTIONS } from '../composables/message-page.utils'
import type { MessageTemplateFormModel } from '../types/message-page'

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
const formModel = defineModel<MessageTemplateFormModel>('model', { required: true })
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
        :title="formMode === 'create' ? '新增消息模板' : '编辑消息模板'"
        :subtitle="
          formMode === 'create'
            ? '定义消息标题、内容和可替换变量。'
            : '修改模板内容后，后续提醒会按新模板生成。'
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
            <p>模板编码创建后保持只读，可在标题和内容中使用变量占位。</p>
          </div>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="编码" path="code">
              <NInput
                v-model:value="formModel.code"
                placeholder="例如 order:paid"
                :disabled="formMode === 'edit'"
              />
            </NFormItem>

            <NFormItem label="名称" path="name">
              <NInput v-model:value="formModel.name" placeholder="例如 订单支付提醒" />
            </NFormItem>

            <NFormItem label="类型">
              <NSelect v-model:value="formModel.type" :options="MESSAGE_TEMPLATE_TYPE_OPTIONS" />
            </NFormItem>

            <NFormItem label="排序">
              <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
            </NFormItem>
          </div>
        </section>

        <section class="ez-modal-section">
          <div class="ez-modal-section__head">
            <h3>消息内容</h3>
          </div>

          <NFormItem label="标题" path="title">
              <NInput v-model:value="formModel.title" placeholder="例如 订单已支付" />
          </NFormItem>

          <NFormItem label="内容" path="content">
            <NInput
              v-model:value="formModel.content"
              type="textarea"
              :rows="4"
              placeholder="请输入消息内容模板"
            />
          </NFormItem>

          <NFormItem label="变量说明">
            <NInput
              v-model:value="formModel.variables"
              type="textarea"
              :rows="3"
              placeholder="例如 order_no=订单号; paid_time=支付时间"
            />
          </NFormItem>
        </section>

        <section class="ez-modal-section">
          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="状态">
              <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
            </NFormItem>

            <NFormItem label="系统级">
              <NCheckbox v-model:checked="formModel.is_system">不可删除配置</NCheckbox>
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
        <NButton type="primary" class="min-w-[92px]" :loading="saving" @click="$emit('submit')"
          >保存</NButton
        >
      </div>
    </template>
  </NModal>
</template>
