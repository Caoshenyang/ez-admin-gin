<script setup lang="ts">
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'

import type { RoleFormModel } from '../composables/useRolePage'

defineProps<{
  formMode: 'create' | 'edit'
  rules: FormRules
  saving: boolean
  show: boolean
  statusOptions: SelectOption[]
}>()

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()

const formRef = defineModel<FormInst | null>('formRef')
const formModel = defineModel<RoleFormModel>('model', { required: true })
</script>

<template>
  <NModal :show="show" preset="card" :closable="false" class="role-modal" style="width: 560px; max-width: calc(100vw - 32px)" @update:show="(value) => $emit('update:show', value)">
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增角色' : '编辑角色'"
        :subtitle="formMode === 'create' ? '角色编码创建后会成为权限策略主体，建议使用稳定英文标识。' : '角色编码保持只读，避免影响已有权限策略。'"
        @close="$emit('update:show', false)"
      />
    </template>

    <NForm ref="formRef" :model="formModel" :rules="rules" label-placement="left" label-width="80">
      <NFormItem label="角色编码" path="code">
        <NInput v-model:value="formModel.code" placeholder="demo_operator" :disabled="formMode === 'edit'" />
      </NFormItem>
      <NFormItem label="角色名称" path="name">
        <NInput v-model:value="formModel.name" placeholder="请输入角色名称" />
      </NFormItem>
      <NFormItem label="排序" path="sort">
        <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
      </NFormItem>
      <NFormItem label="状态" path="status">
        <NSelect v-model:value="formModel.status" :options="statusOptions" />
      </NFormItem>
      <NFormItem label="备注" path="remark">
        <NInput v-model:value="formModel.remark" type="textarea" placeholder="请输入备注" :autosize="{ minRows: 3, maxRows: 5 }" />
      </NFormItem>
    </NForm>

    <template #footer>
      <div class="flex justify-end gap-3">
        <NButton @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" :loading="saving" @click="$emit('submit')">保存</NButton>
      </div>
    </template>
  </NModal>
</template>

<style scoped>
.role-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 20px;
  border: 1px solid #dfe9f5;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.role-modal :deep(.n-card-header) {
  padding: 0;
}
</style>
