<script setup lang="ts">
import type { FormInst, FormRules, UploadFileInfo } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect, NUpload } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'
import { STATUS_FORM_OPTIONS } from '@/constants/status'
import type { AttachmentUploadFormModel } from '../types/attachment-page'

defineProps<{
  fileList: UploadFileInfo[]
  rules: FormRules
  saving: boolean
  show: boolean
}>()

defineEmits<{
  'update:fileList': [fileList: UploadFileInfo[]]
  'update:show': [value: boolean]
  submit: []
}>()

const formRef = defineModel<FormInst | null>('formRef')
const formModel = defineModel<AttachmentUploadFormModel>('model', { required: true })
</script>

<template>
  <NModal
    :show="show"
    preset="card"
    :closable="false"
    class="ez-modal-width-lg"
    :bordered="false"
    @update:show="(value) => $emit('update:show', value)"
  >
    <template #header>
      <FormModalHeader
        title="上传附件"
        subtitle="上传后会进入附件中心，后续可按分类、业务类型和状态继续检索复用。"
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
      <NFormItem label="附件名称" path="display_name">
        <NInput v-model:value="formModel.display_name" placeholder="可留空，默认使用原始文件名" />
      </NFormItem>

      <div class="ez-form-grid ez-form-grid--2">
        <NFormItem label="附件分类" path="category">
          <NInput v-model:value="formModel.category" placeholder="例如 contract / avatar" />
        </NFormItem>
        <NFormItem label="业务类型" path="biz_type">
          <NInput
            v-model:value="formModel.biz_type"
            placeholder="例如 customer / system-template"
          />
        </NFormItem>
      </div>

      <NFormItem label="状态" path="status">
        <NSelect v-model:value="formModel.status" :options="STATUS_FORM_OPTIONS" />
      </NFormItem>

      <NFormItem label="备注" path="remark">
        <NInput
          v-model:value="formModel.remark"
          type="textarea"
          :rows="3"
          placeholder="给这份附件补一段业务备注"
        />
      </NFormItem>

      <NFormItem label="上传文件">
        <NUpload
          :default-upload="false"
          :max="1"
          :file-list="fileList"
          @update:file-list="(fileList) => $emit('update:fileList', fileList)"
        >
          <NButton>选择文件</NButton>
        </NUpload>
      </NFormItem>

      <div class="ez-modal-footer">
        <NButton @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" :loading="saving" @click="$emit('submit')">上传并入库</NButton>
      </div>
    </NForm>
  </NModal>
</template>
