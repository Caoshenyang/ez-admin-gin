<script setup lang="ts">
defineSlots<{
  default?: () => unknown
}>()
</script>

<template>
  <section class="ez-table-card-shell">
    <slot />
  </section>
</template>

<style scoped>
.ez-table-card-shell {
  display: flex;
  min-height: 0;
  flex: 1;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-card-bg);
  box-shadow: var(--ez-component-shadow);
}

.ez-table-card-shell :deep(.n-data-table-th) {
  background: var(--ez-surface-subtle);
  color: var(--ez-text-regular);
  font-size: var(--ez-text-xs);
}

.ez-table-card-shell :deep(.n-data-table-td) {
  font-size: var(--ez-text-sm);
}

.ez-table-card-shell :deep(.n-data-table-pagination) {
  padding: 9px 12px;
}
</style>
