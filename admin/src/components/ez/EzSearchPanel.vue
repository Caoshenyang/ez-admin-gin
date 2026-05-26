<script setup lang="ts">
defineSlots<{
  actions?: () => unknown
  default?: () => unknown
}>()
</script>

<template>
  <section class="ez-search-panel">
    <div class="ez-search-panel__fields">
      <slot />
    </div>
    <div v-if="$slots.actions" class="ez-search-panel__actions">
      <slot name="actions" />
    </div>
  </section>
</template>

<style scoped>
.ez-search-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  border: 1px solid var(--ez-border);
  border-radius: 12px;
  background: var(--ez-card-bg);
  padding: 16px;
  box-shadow: var(--ez-shadow-sm);
}

.ez-search-panel__fields {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
}

.ez-search-panel__actions {
  display: flex;
  flex-shrink: 0;
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}

@media (max-width: 720px) {
  .ez-search-panel {
    align-items: stretch;
    flex-direction: column;
  }

  .ez-search-panel__actions {
    justify-content: flex-start;
  }
}
</style>
