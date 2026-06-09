<script setup lang="ts">
import { NButton, NInput, NModal } from 'naive-ui'

import FormModalHeader from '@/components/FormModalHeader.vue'

defineProps<{
  content: string
  show: boolean
  subject: string
  title: string
}>()

defineEmits<{
  'update:show': [value: boolean]
}>()
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
        :title="title || '模板预览'"
        subtitle="按当前变量样例渲染后的邮件内容。"
        @close="$emit('update:show', false)"
      />
    </template>

    <div class="mail-preview">
      <div class="mail-preview__subject">{{ subject || '无主题' }}</div>
      <NInput :value="content" type="textarea" :rows="12" readonly />
    </div>

    <template #footer>
      <div class="ez-modal-footer">
        <NButton type="primary" class="min-w-[92px]" @click="$emit('update:show', false)">
          关闭
        </NButton>
      </div>
    </template>
  </NModal>
</template>

<style scoped>
.mail-preview {
  display: grid;
  gap: 12px;
}

.mail-preview__subject {
  overflow: hidden;
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-sm);
  background: var(--ez-surface-subtle);
  padding: 10px 12px;
  color: var(--ez-text-heading);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
