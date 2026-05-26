<script setup lang="ts">
import { NButton, NDrawer, NDrawerContent } from 'naive-ui'

withDefaults(defineProps<{
  loading?: boolean
  show: boolean
  title: string
  width?: number | string
}>(), {
  loading: false,
  width: 640,
})

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()

defineSlots<{
  default?: () => unknown
  footer?: () => unknown
}>()
</script>

<template>
  <NDrawer
    :show="show"
    :width="width"
    placement="right"
    @update:show="(value) => $emit('update:show', value)"
  >
    <NDrawerContent :title="title" closable body-content-class="ez-drawer-form__body">
      <slot />
      <template #footer>
        <slot name="footer">
          <div class="ez-drawer-form__footer">
            <NButton @click="$emit('update:show', false)">取消</NButton>
            <NButton type="primary" :loading="loading" @click="$emit('submit')">保存</NButton>
          </div>
        </slot>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style>
.ez-drawer-form__body {
  padding: 20px 24px;
}

.ez-drawer-form__footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.ez-drawer-form__footer .n-button {
  min-width: 92px;
}
</style>
