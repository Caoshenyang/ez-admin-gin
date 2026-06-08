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
  background: #f8fafc;
  color: var(--ez-text-regular);
  font-size: 13px;
  font-weight: 700;
  line-height: 1.35;
  white-space: nowrap;
}

.ez-table-card-shell :deep(.n-data-table-td) {
  color: var(--ez-text-main);
  font-size: 13px;
}

.ez-table-card-shell :deep(.n-data-table-th--fixed-right),
.ez-table-card-shell :deep(.n-data-table-td--fixed-right) {
  box-shadow: -8px 0 14px -14px rgba(15, 23, 42, 0.32);
}

.ez-table-card-shell :deep(.n-button--tiny) {
  --n-padding: 0 9px;
}

.ez-table-card-shell :deep(.n-data-table-pagination) {
  padding: 9px 12px;
}
</style>
