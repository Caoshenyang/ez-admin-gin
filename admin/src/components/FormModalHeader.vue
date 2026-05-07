<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'

defineProps<{
  subtitle?: string
  title: string
}>()

defineEmits<{ close: [] }>()
</script>

<template>
  <div class="modal-header">
    <h2>{{ title }}</h2>
    <p v-if="subtitle">{{ subtitle }}</p>
    <button type="button" class="modal-close" @click="$emit('close')">
      <NIcon :size="18">
        <CloseOutline />
      </NIcon>
    </button>
  </div>
</template>

<style scoped>
.modal-header {
  position: relative;
  padding: 24px 28px;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.modal-header h2 {
  font-size: 19px;
  font-weight: 700;
  color: #111827;
}

.modal-header p {
  margin-top: 8px;
  max-width: 420px;
  font-size: 13px;
  line-height: 1.6;
  color: #64748b;
}

.modal-close {
  position: absolute;
  top: 20px;
  right: 22px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.82);
  color: #64748b;
}
</style>
