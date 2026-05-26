<script setup lang="ts">
defineProps<{
  description?: string
  title: string
}>()

defineSlots<{
  actions?: () => unknown
}>()
</script>

<template>
  <header class="ez-page-header-shell">
    <div>
      <h1>{{ title }}</h1>
      <p v-if="description">{{ description }}</p>
    </div>
    <div v-if="$slots.actions" class="ez-page-header-shell__actions">
      <slot name="actions" />
    </div>
  </header>
</template>

<style scoped>
.ez-page-header-shell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.ez-page-header-shell h1 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--ez-text-main);
  line-height: 28px;
}

.ez-page-header-shell p {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--ez-text-secondary);
  line-height: 20px;
}

.ez-page-header-shell__actions {
  display: flex;
  flex-shrink: 0;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  justify-content: flex-end;
}

@media (max-width: 720px) {
  .ez-page-header-shell {
    align-items: flex-start;
    flex-direction: column;
  }

  .ez-page-header-shell__actions {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
