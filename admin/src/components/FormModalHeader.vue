<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import { NButton, NIcon } from 'naive-ui'

defineProps<{
  subtitle?: string
  title: string
}>()

defineEmits<{ close: [] }>()
</script>

<template>
  <div class="form-modal-header">
    <div class="form-modal-header__copy">
      <h2>{{ title }}</h2>
      <p v-if="subtitle">
        {{ subtitle }}
      </p>
    </div>
    <NButton
      quaternary
      circle
      class="form-modal-header__close"
      aria-label="关闭"
      @click="$emit('close')"
    >
      <NIcon :size="17">
        <CloseOutline />
      </NIcon>
    </NButton>
  </div>
</template>

<style scoped>
.form-modal-header {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: start;
  gap: 12px;
  min-width: 0;
  border-bottom: 1px solid var(--ez-border-light);
  background: var(--ez-card-bg);
  padding: 16px 20px 14px;
}

.form-modal-header__copy {
  min-width: 0;
}

.form-modal-header__copy h2 {
  margin: 0;
  color: var(--ez-text-main);
  font-size: var(--ez-text-xl);
  font-weight: 700;
  line-height: 26px;
}

.form-modal-header__copy p {
  max-width: 480px;
  margin: 4px 0 0;
  color: var(--ez-text-sub);
  font-size: var(--ez-text-sm);
  line-height: 21px;
}

.form-modal-header__close {
  width: 32px;
  height: 32px;
  justify-self: end;
  margin-top: -3px;
  border: 1px solid transparent;
  border-radius: var(--ez-radius-control);
  color: var(--ez-text-sub);
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease,
    color 0.2s ease,
    box-shadow 0.2s ease;
}

.form-modal-header__close:hover {
  border-color: var(--ez-border-light);
  background: var(--ez-surface-subtle);
  color: var(--ez-text-main);
}

.form-modal-header__close:focus-visible {
  box-shadow: 0 0 0 3px var(--ez-brand-muted);
}
</style>
