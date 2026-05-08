<script setup lang="ts">
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect, NSwitch } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
import { MenuStatus, MenuType } from '@/modules/iam/types/menu'

import type { MenuFormModel } from '../types/menu-page'

defineProps<{
  formMode: 'create' | 'edit'
  formTypeOptions: SelectOption[]
  parentOptions: SelectOption[]
  componentOptions: SelectOption[]
  rules: FormRules
  saving: boolean
  show: boolean
}>()

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()

const formRef = defineModel<FormInst | null>('formRef')
const formModel = defineModel<MenuFormModel>('model', { required: true })
</script>

<template>
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="menu-modal"
    style="width: 600px; max-width: calc(100vw - 32px)"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增菜单' : '编辑菜单'"
        :subtitle="formMode === 'create' ? '选择节点类型后填写对应字段。' : '权限标识保持只读，避免影响按钮权限判断。'"
        @close="$emit('update:show', false)"
      />
    </template>

    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      label-width="80"
    >
      <NFormItem label="菜单类型" path="type">
        <div class="type-segment" :class="{ 'is-disabled': formMode === 'edit' }">
          <button
            v-for="opt in formTypeOptions"
            :key="opt.value"
            type="button"
            class="type-segment__btn"
            :class="{ 'type-segment__btn--active': formModel.type === opt.value }"
            :disabled="formMode === 'edit'"
            @click="formModel.type = opt.value as MenuType"
          >
            {{ opt.label }}
          </button>
        </div>
      </NFormItem>

      <NFormItem label="父级节点" path="parent_id">
        <NSelect
          v-model:value="formModel.parent_id"
          filterable
          :options="parentOptions"
        />
      </NFormItem>

      <NFormItem label="菜单名称" path="title">
        <NInput v-model:value="formModel.title" placeholder="请输入菜单名称" />
      </NFormItem>

      <NFormItem label="权限标识" path="code">
        <NInput
          v-model:value="formModel.code"
          placeholder="system:example:list"
          :disabled="formMode === 'edit'"
        />
      </NFormItem>

      <NFormItem v-if="formModel.type !== MenuType.Button" label="路由地址" path="path">
        <NInput v-model:value="formModel.path" placeholder="/system/example" />
      </NFormItem>

      <NFormItem
        v-if="formModel.type === MenuType.Menu"
        label="组件路径"
        path="component"
      >
        <NSelect
          v-model:value="formModel.component"
          filterable
          tag
          :options="componentOptions"
          placeholder="system/UserView"
        />
      </NFormItem>

      <NFormItem label="图标 / 排序">
        <div class="grid w-full grid-cols-[1fr_120px] gap-2">
          <NInput v-model:value="formModel.icon" placeholder="setting / notification / layout-dashboard" />
          <NInputNumber v-model:value="formModel.sort" :min="0" />
        </div>
      </NFormItem>

      <NFormItem label="显示状态">
        <NSwitch
          :value="formModel.status === MenuStatus.Enabled"
          @update:value="
            (checked) => {
              formModel.status = checked ? MenuStatus.Enabled : MenuStatus.Disabled
            }
          "
        />
      </NFormItem>

      <NFormItem label="备注">
        <NInput
          v-model:value="formModel.remark"
          type="textarea"
          placeholder="请输入备注"
          :autosize="{ minRows: 3, maxRows: 5 }"
        />
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
.menu-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 20px;
  border: 1px solid #dfe9f5;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.menu-modal :deep(.n-card-header) {
  padding: 0;
}

.type-segment {
  display: flex;
  gap: 4px;
  padding: 4px;
  border-radius: 6px;
  background: #f3f4f6;
}

.type-segment.is-disabled {
  opacity: 0.6;
  pointer-events: none;
}

.type-segment__btn {
  padding: 4px 20px;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: #6b7280;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}

.type-segment__btn--active {
  background: #fff;
  color: #18a058;
  font-weight: 600;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}
</style>
