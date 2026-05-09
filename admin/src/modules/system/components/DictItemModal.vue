<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'

import { STATUS_FORM_OPTIONS } from '@/constants/status'
import FormModalHeader from '@/components/FormModalHeader.vue'
import type { DictItemFormModel } from '../types/dict-page'
import type { DictTypeItem } from '../types/dict'

defineProps<{
  formMode: 'create' | 'edit'
  rules: FormRules
  saving: boolean
  selectedType: DictTypeItem | null
  show: boolean
}>()

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()

const formRef = defineModel<FormInst | null>('formRef')
const formModel = defineModel<DictItemFormModel>('model', { required: true })
</script>

<template>
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="ez-modal-width-xl"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        :title="formMode === 'create' ? '新增字典项' : '编辑字典项'"
        :subtitle="selectedType ? `当前归属：${selectedType.name}（${selectedType.code}）` : '请选择字典类型后再维护字典项。'"
        @close="$emit('update:show', false)"
      />
    </template>

    <div class="ez-modal-shell">
      <NForm ref="formRef" class="ez-modal-form" :model="formModel" :rules="rules" label-placement="left" label-width="76">
        <section class="ez-modal-section ez-modal-section--soft">
          <div class="ez-modal-section__head">
            <h3>基础信息</h3>
            <p>字典项编码和显示值建议稳定设计，这样状态色、下拉项和表格标签都能长期复用。</p>
          </div>

          <div class="grid gap-x-5 md:grid-cols-2">
            <NFormItem label="字典项编码" path="item_key">
              <NInput v-model:value="formModel.item_key" placeholder="例如 enabled" :disabled="formMode === 'edit'" />
            </NFormItem>

            <NFormItem label="字典项名称" path="label">
              <NInput v-model:value="formModel.label" placeholder="例如 启用" />
            </NFormItem>

            <NFormItem label="字典项值" path="value">
              <NInput v-model:value="formModel.value" placeholder="例如 1" />
            </NFormItem>

            <NFormItem label="标签样式">
              <NInput v-model:value="formModel.tag_type" placeholder="例如 success / warning / info" />
            </NFormItem>

            <NFormItem label="排序">
              <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
            </NFormItem>

            <NFormItem label="状态">
              <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
            </NFormItem>
          </div>
        </section>

        <section class="ez-modal-section">
          <NFormItem label="备注">
            <NInput v-model:value="formModel.remark" type="textarea" :rows="3" placeholder="可填写这项配置值的展示说明或业务备注" />
          </NFormItem>
        </section>
      </NForm>
    </div>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton quaternary class="min-w-[92px]" @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" class="min-w-[92px]" :loading="saving" @click="$emit('submit')">
          保存
        </NButton>
      </div>
    </template>
  </NModal>
</template>
