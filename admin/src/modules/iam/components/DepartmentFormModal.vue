<script setup lang="ts">
import type { FormInst, FormRules, SelectOption, TreeSelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect, NTreeSelect } from 'naive-ui'

import type { DepartmentFormModel } from '../types/department-page'

defineProps<{
  formMode: 'create' | 'edit'
  formStatusOptions: SelectOption[]
  parentOptions: TreeSelectOption[]
  rules: FormRules
  saving: boolean
  show: boolean
}>()

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()

const formRef = defineModel<FormInst | null>('formRef')
const formModel = defineModel<DepartmentFormModel>('model', { required: true })
</script>

<template>
  <NModal :show="show" preset="card" class="w-[680px]" :title="formMode === 'create' ? '新增部门' : '编辑部门'" @update:show="(value) => $emit('update:show', value)">
    <NForm ref="formRef" :model="formModel" :rules="rules" label-placement="top" class="grid gap-4 md:grid-cols-2">
      <NFormItem label="上级部门" path="parent_id">
        <NTreeSelect
          v-model:value="formModel.parent_id"
          :options="parentOptions"
          default-expand-all
          placeholder="请选择上级部门"
        />
      </NFormItem>

      <NFormItem label="负责人用户 ID" path="leader_user_id">
        <NInputNumber v-model:value="formModel.leader_user_id" :min="0" class="w-full" />
      </NFormItem>

      <NFormItem label="部门名称" path="name">
        <NInput v-model:value="formModel.name" placeholder="请输入部门名称" />
      </NFormItem>

      <NFormItem label="部门编码" path="code">
        <NInput v-model:value="formModel.code" placeholder="请输入部门编码" />
      </NFormItem>

      <NFormItem label="排序" path="sort">
        <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
      </NFormItem>

      <NFormItem label="状态" path="status">
        <NSelect v-model:value="formModel.status" :options="formStatusOptions" />
      </NFormItem>

      <NFormItem label="备注" path="remark" class="md:col-span-2">
        <NInput v-model:value="formModel.remark" type="textarea" :rows="4" placeholder="补充记录这个部门的职责、边界或特殊说明" />
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
