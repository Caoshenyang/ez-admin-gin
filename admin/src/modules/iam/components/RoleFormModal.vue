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

import { RoleDataScope } from '../types/role'
import type { RoleFormModel } from '../types/role-page'

defineProps<{
  dataScopeOptions: SelectOption[]
  departmentTreeOptions: TreeSelectOption[]
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
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="ez-modal-width-md"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增角色' : '编辑角色'"
        :subtitle="
          formMode === 'create'
            ? '角色编码创建后会成为权限策略主体，建议使用稳定英文标识。'
            : '角色编码保持只读，避免影响已有权限策略。'
        "
        @close="$emit('update:show', false)"
      />
    </template>

    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      label-width="80"
      class="ez-modal-form"
    >
      <NFormItem label="角色编码" path="code">
        <NInput
          v-model:value="formModel.code"
          placeholder="demo_operator"
          :disabled="formMode === 'edit'"
        />
      </NFormItem>
      <NFormItem label="角色名称" path="name">
        <NInput v-model:value="formModel.name" placeholder="请输入角色名称" />
      </NFormItem>
      <NFormItem label="排序" path="sort">
        <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
      </NFormItem>
      <NFormItem label="数据范围" path="data_scope">
        <NSelect v-model:value="formModel.data_scope" :options="dataScopeOptions" />
      </NFormItem>
      <NFormItem
        v-if="formModel.data_scope === RoleDataScope.CustomDept"
        label="授权部门"
        path="custom_department_ids"
      >
        <NTreeSelect
          v-model:value="formModel.custom_department_ids"
          :options="departmentTreeOptions"
          multiple
          checkable
          default-expand-all
          placeholder="请选择授权部门"
        />
      </NFormItem>
      <NFormItem label="状态" path="status">
        <NSelect v-model:value="formModel.status" :options="statusOptions" />
      </NFormItem>
      <NFormItem label="备注" path="remark">
        <NInput
          v-model:value="formModel.remark"
          type="textarea"
          placeholder="请输入备注"
          :autosize="{ minRows: 3, maxRows: 5 }"
        />
      </NFormItem>
    </NForm>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" :loading="saving" @click="$emit('submit')">保存</NButton>
      </div>
    </template>
  </NModal>
</template>
