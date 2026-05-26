<script setup lang="ts">
import { NButton, NPopconfirm } from 'naive-ui'

withDefaults(defineProps<{
  buttonText?: string
  disabled?: boolean
  message?: string
  size?: 'tiny' | 'small' | 'medium' | 'large'
}>(), {
  buttonText: '删除',
  disabled: false,
  message: '确认执行该危险操作？',
  size: 'small',
})

defineEmits<{
  confirm: []
}>()
</script>

<template>
  <NPopconfirm :disabled="disabled" @positive-click="$emit('confirm')">
    <template #trigger>
      <NButton :size="size" type="error" ghost :disabled="disabled">
        {{ buttonText }}
      </NButton>
    </template>
    {{ message }}
  </NPopconfirm>
</template>
