<script setup lang="ts">
import type { FormInst, FormRules, UploadFileInfo } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect, NUpload } from 'naive-ui'

import type { AttachmentStatus } from '../types/attachment'

interface UploadFormModel {
  display_name: string
  category: string
  biz_type: string
  status: AttachmentStatus
  remark: string
}

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
const formModel = defineModel<UploadFormModel>('model', { required: true })
</script>

<template>
  <NModal
    :show="show"
    preset="card"
    title="上传附件"
    class="max-w-[620px] rounded-2xl"
    :bordered="false"
    @update:show="(value) => $emit('update:show', value)"
  >
    <NForm ref="formRef" :model="formModel" :rules="rules" label-placement="top">
      <NFormItem label="附件名称" path="display_name">
        <NInput v-model:value="formModel.display_name" placeholder="可留空，默认使用原始文件名" />
      </NFormItem>

      <div class="grid gap-4 md:grid-cols-2">
        <NFormItem label="附件分类" path="category">
          <NInput v-model:value="formModel.category" placeholder="例如 contract / avatar" />
        </NFormItem>
        <NFormItem label="业务类型" path="biz_type">
          <NInput v-model:value="formModel.biz_type" placeholder="例如 customer / system-template" />
        </NFormItem>
      </div>

      <NFormItem label="状态" path="status">
        <NSelect v-model:value="formModel.status" :options="[{ label: '启用', value: 1 }, { label: '禁用', value: 2 }]" />
      </NFormItem>

      <NFormItem label="备注" path="remark">
        <NInput v-model:value="formModel.remark" type="textarea" :rows="3" placeholder="给这份附件补一段业务备注" />
      </NFormItem>

      <NFormItem label="上传文件">
        <NUpload :default-upload="false" :max="1" :file-list="fileList" @update:file-list="(fileList) => $emit('update:fileList', fileList)">
          <NButton>选择文件</NButton>
        </NUpload>
      </NFormItem>

      <div class="flex justify-end gap-3">
        <NButton @click="$emit('update:show', false)">取消</NButton>
        <NButton type="primary" :loading="saving" @click="$emit('submit')">上传并入库</NButton>
      </div>
    </NForm>
  </NModal>
</template>
