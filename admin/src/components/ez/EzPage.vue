<script setup lang="ts">
defineSlots<{
  default?: () => unknown
}>()
</script>

<template>
  <main class="ez-page">
    <slot />
  </main>
</template>

<style scoped>
.ez-page {
  display: flex;
  min-height: 0;
  height: 100%;
  flex: 1;
  flex-direction: column;
  gap: 16px;
  background: var(--ez-page-bg);
}
</style>
