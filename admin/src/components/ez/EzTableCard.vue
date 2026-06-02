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
  border: 1px solid var(--ez-border);
  border-radius: var(--ez-radius-card);
  background: var(--ez-card-bg);
  box-shadow: var(--ez-shadow-card);
}
</style>
