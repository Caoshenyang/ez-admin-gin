<script setup lang="ts">
import { NPopconfirm } from 'naive-ui'

import EzActionButton from './EzActionButton.vue'

withDefaults(
  defineProps<{
    buttonText?: string
    disabled?: boolean
    message?: string
    size?: 'tiny' | 'small' | 'medium' | 'large'
  }>(),
  {
    buttonText: '删除',
    disabled: false,
    message: '确认执行该危险操作？',
    size: 'small',
  },
)

defineEmits<{
  confirm: []
}>()
</script>

<template>
  <NPopconfirm :disabled="disabled" @positive-click="$emit('confirm')">
    <template #trigger>
      <EzActionButton
        kind="delete"
        :label="buttonText"
        :size="size"
        type="error"
        ghost
        :disabled="disabled"
      />
    </template>
    {{ message }}
  </NPopconfirm>
</template>
