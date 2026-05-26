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
  min-height: 0;
  flex: 1;
  overflow: hidden;
  border: 1px solid var(--ez-border);
  border-radius: 12px;
  background: var(--ez-card-bg);
  box-shadow: var(--ez-shadow-card);
}
</style>
