<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

import { STATUS_FORM_OPTIONS } from '@/constants/status'
import FormModalHeader from '@/components/FormModalHeader.vue'
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
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="ez-modal-width-md"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增配置' : '编辑配置'"
        :subtitle="
          formMode === 'create'
            ? '填写配置分组、键和值，保存后立即生效。'
            : '修改配置名称和值，配置键在创建后保持只读。'
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
            <p>配置键只允许小写字母、数字、冒号、短横线和下划线。</p>
          </div>

          <div class="ez-form-grid ez-form-grid--2">
            <NFormItem label="分组" path="group_code">
              <NInput v-model:value="formModel.group_code" placeholder="例如 site" />
            </NFormItem>

            <NFormItem label="键" path="key">
              <NInput
                v-model:value="formModel.key"
                placeholder="例如 site_name"
                :disabled="formMode === 'edit'"
              />
            </NFormItem>

            <NFormItem label="名称" path="name">
              <NInput v-model:value="formModel.name" placeholder="例如 站点名称" />
            </NFormItem>

            <NFormItem label="排序">
              <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
            </NFormItem>
          </div>
        </section>

        <section class="ez-modal-section">
          <div class="ez-modal-section__head">
            <h3>配置值</h3>
            <p>配置值支持任意文本，启用后会被缓存到 Redis。</p>
          </div>

          <NFormItem label="值" path="value">
            <NInput
              v-model:value="formModel.value"
              type="textarea"
              :rows="3"
              placeholder="请输入配置值"
            />
          </NFormItem>
        </section>

        <section class="ez-modal-section">
          <div class="ez-form-grid ez-form-grid--2">
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
      <div class="ez-modal-footer">
        <NButton quaternary class="min-w-[92px]" @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" class="min-w-[92px]" :loading="saving" @click="$emit('submit')"
          >保存</NButton
        >
      </div>
    </template>
  </NModal>
</template>
