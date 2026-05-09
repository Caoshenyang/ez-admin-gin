<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

import { STATUS_FORM_OPTIONS } from '@/constants/status'
import FormModalHeader from '@/components/FormModalHeader.vue'
import type { PostFormModel } from '../types/post-page'

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
const formModel = defineModel<PostFormModel>('model', { required: true })
</script>

<template>
  <NModal :show="show" preset="card" :closable="false" class="ez-modal-width-lg" @update:show="(value) => $emit('update:show', value)">
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增岗位' : '编辑岗位'"
        :subtitle="formMode === 'create' ? '岗位编码建议稳定设计，便于在用户归属和业务协作中长期复用。' : '修改岗位信息时，优先保持岗位编码稳定，避免影响已有归属关系。'"
        @close="$emit('update:show', false)"
      />
    </template>

    <NForm ref="formRef" :model="formModel" :rules="rules" label-placement="top" class="grid gap-4 md:grid-cols-2">
      <NFormItem label="岗位编码" path="code">
        <NInput v-model:value="formModel.code" placeholder="请输入岗位编码" />
      </NFormItem>

      <NFormItem label="岗位名称" path="name">
        <NInput v-model:value="formModel.name" placeholder="请输入岗位名称" />
      </NFormItem>

      <NFormItem label="排序" path="sort">
        <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
      </NFormItem>

      <NFormItem label="状态" path="status">
        <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
      </NFormItem>

      <NFormItem label="备注" path="remark" class="md:col-span-2">
        <NInput v-model:value="formModel.remark" type="textarea" :rows="4" placeholder="补充记录岗位职责、适用部门或协作边界" />
      </NFormItem>
    </NForm>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton quaternary @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" :loading="saving" @click="$emit('submit')">保存</NButton>
      </div>
    </template>
  </NModal>
</template>
