<script setup lang="ts">
defineSlots<{
  actions?: () => unknown
  default?: () => unknown
}>()
</script>

<template>
  <div class="ez-toolbar-line">
    <div class="ez-toolbar-line__meta">
      <slot />
    </div>
    <div v-if="$slots.actions" class="ez-toolbar-line__actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
.ez-toolbar-line {
  display: flex;
  min-height: 42px;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  border-bottom: 1px solid var(--ez-border-light);
  background: var(--ez-surface-subtle);
  padding: 8px 12px;
  color: var(--ez-text-secondary);
  font-size: var(--ez-text-sm);
}

.ez-toolbar-line__meta,
.ez-toolbar-line__actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.ez-toolbar-line__actions {
  justify-content: flex-end;
}

@media (max-width: 720px) {
  .ez-toolbar-line {
    align-items: stretch;
    flex-direction: column;
  }

  .ez-toolbar-line__actions {
    justify-content: flex-start;
  }
}
</style>
