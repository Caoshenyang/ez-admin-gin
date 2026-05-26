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
  min-height: 52px;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border-bottom: 1px solid var(--ez-border-light);
  padding: 10px 20px;
  color: var(--ez-text-secondary);
  font-size: 13px;
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
