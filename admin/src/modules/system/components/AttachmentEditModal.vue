<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect } from 'naive-ui'

import type { UpdateAttachmentPayload } from '../types/attachment'

interface EditFormModel extends UpdateAttachmentPayload {
  id: number
}

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
const formModel = defineModel<EditFormModel>('model', { required: true })
</script>

<template>
  <NModal :show="show" preset="card" title="编辑附件" class="max-w-[620px] rounded-2xl" :bordered="false" @update:show="(value) => $emit('update:show', value)">
    <NForm ref="formRef" :model="formModel" :rules="rules" label-placement="top">
      <NFormItem label="附件名称" path="display_name">
        <NInput v-model:value="formModel.display_name" placeholder="请输入附件名称" />
      </NFormItem>

      <div class="grid gap-4 md:grid-cols-2">
        <NFormItem label="附件分类" path="category">
          <NInput v-model:value="formModel.category" placeholder="附件分类" />
        </NFormItem>
        <NFormItem label="业务类型" path="biz_type">
          <NInput v-model:value="formModel.biz_type" placeholder="业务类型" />
        </NFormItem>
      </div>

      <NFormItem label="状态" path="status">
        <NSelect v-model:value="formModel.status" :options="[{ label: '启用', value: 1 }, { label: '禁用', value: 2 }]" />
      </NFormItem>

      <NFormItem label="备注" path="remark">
        <NInput v-model:value="formModel.remark" type="textarea" :rows="3" placeholder="附件备注" />
      </NFormItem>

      <div class="flex justify-end gap-3">
        <NButton @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" :loading="saving" @click="$emit('submit')">保存</NButton>
      </div>
    </NForm>
  </NModal>
</template>
