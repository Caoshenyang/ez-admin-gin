<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NIcon, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

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
      <div class="modal-header">
        <h2>{{ formMode === 'create' ? '新增角色' : '编辑角色' }}</h2>
        <p>{{ formMode === 'create' ? '角色编码创建后会成为权限策略主体，建议使用稳定英文标识。' : '角色编码保持只读，避免影响已有权限策略。' }}</p>
        <button type="button" class="modal-close" @click="$emit('update:show', false)">
          <NIcon :size="18"><CloseOutline /></NIcon>
        </button>
      </div>
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

.modal-header {
  position: relative;
  padding: 24px 28px;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.modal-header h2 {
  font-size: 19px;
  font-weight: 700;
  color: #111827;
}

.modal-header p {
  margin-top: 8px;
  max-width: 420px;
  font-size: 13px;
  line-height: 1.6;
  color: #64748b;
}

.modal-close {
  position: absolute;
  top: 20px;
  right: 22px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.82);
  color: #64748b;
}
</style>
