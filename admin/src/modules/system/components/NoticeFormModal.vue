<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

import { STATUS_FORM_OPTIONS } from '@/constants/status'
import FormModalHeader from '@/components/FormModalHeader.vue'
import type { NoticeFormModel } from '../composables/notice-page.utils'

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
const formModel = defineModel<NoticeFormModel>('model', { required: true })
</script>

<template>
  <NModal :show="show" preset="card" :closable="false" class="ez-modal-width-md" @update:show="(value) => $emit('update:show', value)">
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增公告' : '编辑公告'"
        :subtitle="formMode === 'create' ? '填写公告标题和内容，保存后可立即展示。' : '修改标题和内容后，状态变更会立即生效。'"
        @close="$emit('update:show', false)"
      />
    </template>

    <div class="ez-modal-shell">
      <NForm ref="formRef" class="ez-modal-form" :model="formModel" :rules="rules" label-placement="left" label-width="76">
        <section class="ez-modal-section ez-modal-section--soft">
          <div class="ez-modal-section__head">
            <h3>公告信息</h3>
            <p>标题不超过 128 个字符，内容支持任意文本。</p>
          </div>

          <div class="grid gap-x-5 md:grid-cols-2">
            <NFormItem label="标题" path="title">
              <NInput v-model:value="formModel.title" placeholder="公告标题" />
            </NFormItem>

            <NFormItem label="排序">
              <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
            </NFormItem>
          </div>
        </section>

        <section class="ez-modal-section">
          <div class="ez-modal-section__head">
            <h3>公告内容</h3>
          </div>

          <NFormItem label="内容">
            <NInput v-model:value="formModel.content" type="textarea" :rows="4" placeholder="请输入公告内容" />
          </NFormItem>
        </section>

        <section class="ez-modal-section">
          <div class="grid gap-x-5 md:grid-cols-2">
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
        <NButton type="primary" class="min-w-[92px]" :loading="saving" @click="$emit('submit')">保存</NButton>
      </div>
    </template>
  </NModal>
</template>
