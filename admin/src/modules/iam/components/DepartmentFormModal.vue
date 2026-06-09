<script setup lang="ts">
import type { FormInst, FormRules, SelectOption, TreeSelectOption } from 'naive-ui'
import {
  NButton,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NSelect,
  NTreeSelect,
} from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
import UserSelect from '@/components/UserSelect.vue'
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
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="ez-modal-width-lg"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增部门' : '编辑部门'"
        subtitle="部门编码和层级会影响用户归属与数据权限，创建后建议保持稳定。"
        @close="$emit('update:show', false)"
      />
    </template>

    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="top"
      class="ez-modal-form"
    >
      <div class="ez-form-grid ez-form-grid--2">
        <NFormItem label="上级部门" path="parent_id">
          <NTreeSelect
            v-model:value="formModel.parent_id"
            :options="parentOptions"
            default-expand-all
            key-field="value"
            label-field="label"
            value-field="value"
            placeholder="请选择上级部门"
          />
        </NFormItem>

        <NFormItem label="负责人" path="leader_user_id">
          <UserSelect v-model:value="formModel.leader_user_id" placeholder="请选择负责人" />
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
      </div>

      <NFormItem label="备注" path="remark">
        <NInput
          v-model:value="formModel.remark"
          type="textarea"
          :rows="4"
          placeholder="补充记录这个部门的职责、边界或特殊说明"
        />
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
